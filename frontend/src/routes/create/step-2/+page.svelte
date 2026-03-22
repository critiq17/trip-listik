<script lang="ts">
	import { goto } from '$app/navigation';
	import { setDraft, getDraft } from '$lib/tripDraft';
	import { onMount, onDestroy } from 'svelte';
	import { setupMainButton, hideMainButton, setupBackButton } from '$lib/telegram';

	let draft = getDraft() || {};
	let visibility = $state<'public' | 'private'>(draft.visibility === 'private' ? 'private' : 'public');
	let startDate = $state<Date | null>(draft.start_date ? new Date(draft.start_date) : null);
	let endDate = $state<Date | null>(draft.end_date ? new Date(draft.end_date) : null);

	// Calendar state
	let calendarYear = $state(new Date().getFullYear());
	let calendarMonth = $state(new Date().getMonth()); // 0-indexed

	const MONTH_NAMES = [
		'January','February','March','April','May','June',
		'July','August','September','October','November','December'
	];
	const DAYS = ['SU','MO','TU','WE','TH','FR','SA'];

	// Format date for display
	const fmtDate = (d: Date | null) => {
		if (!d) return '—';
		return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' });
	};

	// Format date as yyyy-mm-dd string for draft
	const toISODate = (d: Date | null) => {
		if (!d) return '';
		return d.toISOString().split('T')[0];
	};

	// Build calendar grid
	function buildCalDays(year: number, month: number) {
		const firstDay = new Date(year, month, 1).getDay();
		const daysInMonth = new Date(year, month + 1, 0).getDate();
		const prevMonthDays = new Date(year, month, 0).getDate();

		const days: { day: number; month: 'prev' | 'cur' | 'next'; date: Date }[] = [];

		// Previous month fill
		for (let i = firstDay - 1; i >= 0; i--) {
			const d = prevMonthDays - i;
			days.push({ day: d, month: 'prev', date: new Date(year, month - 1, d) });
		}

		// Current month
		for (let i = 1; i <= daysInMonth; i++) {
			days.push({ day: i, month: 'cur', date: new Date(year, month, i) });
		}

		// Next month fill
		const remaining = 42 - days.length;
		for (let i = 1; i <= remaining; i++) {
			days.push({ day: i, month: 'next', date: new Date(year, month + 1, i) });
		}

		return days;
	}

	const calDays = $derived(buildCalDays(calendarYear, calendarMonth));

	const clickDay = (date: Date) => {
		if (!startDate || (startDate && endDate)) {
			// Start fresh selection
			startDate = date;
			endDate = null;
		} else {
			// Set end date (must be after start)
			if (date < startDate) {
				startDate = date;
				endDate = null;
			} else {
				endDate = date;
			}
		}
	};

	const isStart = (date: Date) => startDate && date.toDateString() === startDate.toDateString();
	const isEnd = (date: Date) => endDate && date.toDateString() === endDate.toDateString();
	const isInRange = (date: Date) => {
		if (!startDate || !endDate) return false;
		return date > startDate && date < endDate;
	};

	const prevMonth = () => {
		if (calendarMonth === 0) {
			calendarMonth = 11;
			calendarYear--;
		} else {
			calendarMonth--;
		}
	};

	const nextMonth = () => {
		if (calendarMonth === 11) {
			calendarMonth = 0;
			calendarYear++;
		} else {
			calendarMonth++;
		}
	};

	$effect(() => {
		if (startDate && endDate) {
			setupMainButton('Continue', next, true, true);
		} else {
			setupMainButton('Continue', next, true, false);
		}
	});

	onMount(() => {
		setupBackButton(() => history.back());
	});

	onDestroy(() => {
		hideMainButton();
	});

	const next = () => {
		if (!startDate || !endDate) return;
		setDraft({
			start_date: toISODate(startDate),
			end_date: toISODate(endDate),
			visibility: visibility
		});
		goto('/create/step-3');
	};
</script>

<svelte:head>
	<title>TripListik — Create Trip: Dates</title>
	<meta name="description" content="Set the dates for your upcoming trip." />
</svelte:head>

<div class="page">
	<!-- Sticky Header -->
	<header class="top-bar">
		<button class="back-btn" aria-label="Go back" onclick={() => history.back()}>
			<span class="material-symbols-outlined">arrow_back</span>
		</button>
		<span class="header-title">Create Trip</span>
		<div class="header-spacer"></div>
	</header>

	<!-- Progress -->
	<div class="progress-track">
		<div class="progress-fill" style="width: 66.66%"></div>
	</div>

	<main class="content">
		<div class="progress-meta">
			<span class="step-label">Step 2 of 3</span>
			<span class="pct-label">66% Complete</span>
		</div>

		<h1>When does the journey begin?</h1>

		<!-- Date Inputs -->
		<div class="date-grid">
			<div class="date-field" class:active={startDate !== null}>
				<p class="date-field-label">FROM</p>
				<div class="date-value">
					<span class="material-symbols-outlined date-icon">calendar_today</span>
					<span>{fmtDate(startDate)}</span>
				</div>
			</div>
			<div class="date-field" class:active={endDate !== null}>
				<p class="date-field-label">TO</p>
				<div class="date-value">
					<span class="material-symbols-outlined date-icon">event</span>
					<span>{fmtDate(endDate)}</span>
				</div>
			</div>
		</div>

		<!-- Calendar Widget -->
		<div class="calendar">
			<!-- Month header -->
			<div class="cal-header">
				<button class="cal-nav" aria-label="Previous month" onclick={prevMonth}>
					<span class="material-symbols-outlined">chevron_left</span>
				</button>
				<p class="cal-title">{MONTH_NAMES[calendarMonth]} {calendarYear}</p>
				<button class="cal-nav" aria-label="Next month" onclick={nextMonth}>
					<span class="material-symbols-outlined">chevron_right</span>
				</button>
			</div>

			<!-- Day labels -->
			<div class="cal-days-header">
				{#each DAYS as d}
					<span>{d}</span>
				{/each}
			</div>

			<!-- Day grid -->
			<div class="cal-grid">
				{#each calDays as entry}
					<button
						class="cal-day"
						class:other-month={entry.month !== 'cur'}
						class:is-start={isStart(entry.date)}
						class:is-end={isEnd(entry.date)}
						class:in-range={isInRange(entry.date)}
						onclick={() => clickDay(entry.date)}
						aria-label="Select {entry.date.toDateString()}"
					>
						{entry.day}
					</button>
				{/each}
			</div>
		</div>

		<!-- Visibility Toggle -->
		<div class="visibility-card">
			<div class="vis-left">
				<p class="vis-title">Trip Visibility</p>
				<p class="vis-sub">Who can see your itinerary?</p>
			</div>
			<div class="vis-toggle">
				<div class="toggle-track">
					<div class="toggle-pill" class:private={visibility === 'private'}></div>
					<button
						class="toggle-option"
						class:active={visibility === 'public'}
						onclick={() => (visibility = 'public')}
					>PUBLIC</button>
					<button
						class="toggle-option"
						class:active={visibility === 'private'}
						onclick={() => (visibility = 'private')}
					>PRIVATE</button>
				</div>
			</div>
		</div>

		<!-- Continue Button -->
		<button
			class="continue-btn"
			onclick={next}
			disabled={!startDate || !endDate}
		>
			Continue
			<span class="material-symbols-outlined">arrow_forward</span>
		</button>
	</main>
</div>

<style>
.page {
	min-height: 100dvh;
	background: #161c18;
	color: #f1f5f9;
}

/* ── Sticky Header ─────────────────────────────────────── */
.top-bar {
	position: sticky;
	top: 0;
	z-index: 50;
	display: flex;
	align-items: center;
	height: 56px;
	padding: 0 16px;
	background: #161c18;
	border-bottom: 1px solid rgba(77, 157, 109, 0.1);
}

.back-btn {
	width: 44px;
	height: 44px;
	display: flex;
	align-items: center;
	justify-content: center;
	color: white;
	background: none;
	border: none;
	cursor: pointer;
	flex-shrink: 0;
}

.back-btn .material-symbols-outlined {
	font-size: 22px;
}

.header-title {
	flex: 1;
	text-align: center;
	font-size: 18px;
	font-weight: 700;
	color: white;
}

.header-spacer {
	width: 44px;
}

/* ── Progress ────────────────────────────────────────────── */
.progress-track {
	width: 100%;
	height: 8px;
	background: rgba(77, 157, 109, 0.15);
	border-radius: 9999px;
}

.progress-fill {
	height: 100%;
	background: #4d9d6d;
	border-radius: 9999px;
	transition: width 500ms cubic-bezier(0.4, 0, 0.2, 1);
}

/* ── Content ─────────────────────────────────────────────── */
.content {
	padding: 16px 16px 96px;
	max-width: 480px;
	margin: 0 auto;
}

.progress-meta {
	display: flex;
	justify-content: space-between;
	align-items: center;
	padding-top: 8px;
	margin-bottom: 16px;
}

.step-label {
	font-size: 13px;
	color: #4d9d6d;
	font-weight: 600;
}

.pct-label {
	font-size: 13px;
	color: #94a3b8;
}

h1 {
	font-size: 24px;
	font-weight: 700;
	color: white;
	margin-bottom: 24px;
	line-height: 1.2;
}

/* ── Date Inputs ─────────────────────────────────────────── */
.date-grid {
	display: grid;
	grid-template-columns: 1fr 1fr;
	gap: 16px;
	margin-bottom: 32px;
}

.date-field {
	border-bottom: 2px solid rgba(77, 157, 109, 0.3);
	padding-bottom: 8px;
	transition: border-color 0.2s ease;
}

.date-field.active {
	border-bottom-color: #4d9d6d;
}

.date-field-label {
	font-size: 10px;
	color: #94a3b8;
	text-transform: uppercase;
	letter-spacing: 0.1em;
	margin-bottom: 4px;
}

.date-value {
	display: flex;
	align-items: center;
	gap: 8px;
	padding: 8px 0;
}

.date-icon {
	font-size: 20px;
	color: #4d9d6d;
}

.date-value span:last-child {
	font-size: 16px;
	font-weight: 500;
	color: white;
}

/* ── Calendar ────────────────────────────────────────────── */
.calendar {
	background: rgba(255, 255, 255, 0.03);
	border: 1px solid rgba(77, 157, 109, 0.15);
	border-radius: 12px;
	padding: 16px;
	margin-bottom: 32px;
}

.cal-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 16px;
}

.cal-nav {
	width: 32px;
	height: 32px;
	display: flex;
	align-items: center;
	justify-content: center;
	background: none;
	border: none;
	color: white;
	cursor: pointer;
	border-radius: 6px;
	transition: background 0.15s ease;
}

.cal-nav:hover {
	background: rgba(255, 255, 255, 0.06);
}

.cal-nav .material-symbols-outlined {
	font-size: 20px;
}

.cal-title {
	font-size: 14px;
	font-weight: 700;
	color: white;
}

.cal-days-header {
	display: grid;
	grid-template-columns: repeat(7, 1fr);
	margin-bottom: 8px;
}

.cal-days-header span {
	text-align: center;
	font-size: 10px;
	color: #94a3b8;
	font-weight: 700;
}

.cal-grid {
	display: grid;
	grid-template-columns: repeat(7, 1fr);
	row-gap: 4px;
}

.cal-day {
	display: flex;
	align-items: center;
	justify-content: center;
	height: 36px;
	font-size: 13px;
	color: #e2e8f0;
	background: none;
	border: none;
	cursor: pointer;
	border-radius: 8px;
	transition: background 0.15s ease, color 0.15s ease;
}

.cal-day.other-month {
	color: #4a5568;
}

.cal-day.in-range {
	background: rgba(77, 157, 109, 0.15);
	border-radius: 0;
}

.cal-day.is-start,
.cal-day.is-end {
	background: #4d9d6d;
	color: white;
	font-weight: 700;
	border-radius: 8px;
}

.cal-day.is-start {
	border-radius: 8px 0 0 8px;
}

.cal-day.is-end {
	border-radius: 0 8px 8px 0;
}

.cal-day.is-start:not(.is-end),
.cal-day.is-end:not(.is-start) {
	position: relative;
	z-index: 1;
}

.cal-day:hover:not(.is-start):not(.is-end) {
	background: rgba(77, 157, 109, 0.2);
}

/* ── Visibility Toggle ───────────────────────────────────── */
.visibility-card {
	background: rgba(77, 157, 109, 0.06);
	border: 1px solid rgba(77, 157, 109, 0.2);
	border-radius: 12px;
	padding: 16px;
	display: flex;
	align-items: center;
	justify-content: space-between;
	gap: 16px;
	margin-bottom: 32px;
}

.vis-left {
	flex: 1;
}

.vis-title {
	font-size: 14px;
	font-weight: 700;
	color: white;
	margin-bottom: 4px;
}

.vis-sub {
	font-size: 13px;
	color: #94a3b8;
}

/* Pill toggle */
.vis-toggle {
	flex-shrink: 0;
}

.toggle-track {
	position: relative;
	display: flex;
	background: #161c18;
	border-radius: 9999px;
	width: 140px;
	height: 40px;
	padding: 4px;
	overflow: hidden;
}

.toggle-pill {
	position: absolute;
	top: 4px;
	left: 4px;
	width: calc(50% - 4px);
	height: 32px;
	background: #4d9d6d;
	border-radius: 9999px;
	transition: transform 0.2s ease;
	pointer-events: none;
}

.toggle-pill.private {
	transform: translateX(100%);
}

.toggle-option {
	flex: 1;
	z-index: 1;
	background: none;
	border: none;
	font-size: 11px;
	font-weight: 700;
	cursor: pointer;
	color: #64748b;
	border-radius: 9999px;
	transition: color 0.2s ease;
}

.toggle-option.active {
	color: white;
}

/* ── Continue Button ─────────────────────────────────────── */
.continue-btn {
	width: 100%;
	display: flex;
	align-items: center;
	justify-content: center;
	gap: 8px;
	padding: 16px;
	border-radius: 12px;
	background: #4d9d6d;
	color: white;
	font-size: 16px;
	font-weight: 700;
	border: none;
	cursor: pointer;
	box-shadow: 0 8px 25px rgba(77, 157, 109, 0.3);
	transition: opacity 0.15s ease, transform 0.1s ease;
}

.continue-btn:hover:not(:disabled) {
	opacity: 0.92;
}

.continue-btn:active:not(:disabled) {
	transform: scale(0.98);
}

.continue-btn:disabled {
	opacity: 0.4;
	cursor: not-allowed;
}

.continue-btn .material-symbols-outlined {
	font-size: 20px;
}
</style>
