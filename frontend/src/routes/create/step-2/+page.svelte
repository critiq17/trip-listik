<script lang="ts">
	import { goto } from '$app/navigation';
	import { setDraft, getDraft } from '$lib/tripDraft';
	import { onMount, onDestroy } from 'svelte';
	import { setupMainButton, hideMainButton, setupBackButton, isTelegramEnv } from '$lib/telegram';
	import { parseCalendarDate, toCalendarDateString } from '$lib/format';

	let showFallbackButton = $state(false);
	let draft = getDraft() || {};
	let visibility = $state<'public' | 'private'>(draft.visibility === 'private' ? 'private' : 'public');
	let startDate = $state<Date | null>(parseCalendarDate(draft.start_date));
	let endDate = $state<Date | null>(parseCalendarDate(draft.end_date));

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
		showFallbackButton = !isTelegramEnv();
	});

	onDestroy(() => {
		hideMainButton();
	});

	const next = () => {
		if (!startDate || !endDate) return;
		setDraft({
			start_date: toCalendarDateString(startDate),
			end_date: toCalendarDateString(endDate),
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
	<header class="top-bar">
		<button class="back-btn" aria-label="Go back" onclick={() => history.back()}>
			<span class="material-symbols-outlined">arrow_back</span>
		</button>
		<span class="header-title">New Trip</span>
		<div class="header-spacer"></div>
	</header>

	<div class="progress-track">
		<div class="progress-fill" style="width: 66.66%"></div>
	</div>

	<main class="content">
		<p class="step-label">Step 2 of 3</p>
		<h1>When does the journey begin?</h1>

		<!-- Date Inputs -->
		<div class="date-grid">
			<div class="date-field" class:active={startDate !== null}>
				<p class="date-field-label">From</p>
				<p class="date-value">{fmtDate(startDate)}</p>
			</div>
			<div class="date-field" class:active={endDate !== null}>
				<p class="date-field-label">To</p>
				<p class="date-value">{fmtDate(endDate)}</p>
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
				<p class="vis-title">Trip visibility</p>
				<p class="vis-sub">Who can see your itinerary?</p>
			</div>
			<div class="toggle-track">
				<div class="toggle-pill" class:private={visibility === 'private'}></div>
				<button
					class="toggle-option"
					class:active={visibility === 'public'}
					onclick={() => (visibility = 'public')}
				>Public</button>
				<button
					class="toggle-option"
					class:active={visibility === 'private'}
					onclick={() => (visibility = 'private')}
				>Private</button>
			</div>
		</div>

		<!-- Fallback for plain-browser sessions; Telegram shows the native MainButton -->
		{#if showFallbackButton}
			<button
				class="continue-btn"
				onclick={next}
				disabled={!startDate || !endDate}
			>
				Continue
			</button>
		{/if}
	</main>
</div>

<style>
.page {
	min-height: 100dvh;
	background: var(--bg);
	color: var(--text);
}

/* ── Header ─────────────────────────────────────────────── */
.top-bar {
	position: sticky;
	top: 0;
	z-index: 50;
	display: flex;
	align-items: center;
	height: 52px;
	padding: 0 8px;
	background: var(--bg);
	border-bottom: 1px solid var(--border);
}

.back-btn {
	width: 44px;
	height: 44px;
	display: flex;
	align-items: center;
	justify-content: center;
	color: var(--text);
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
	font-size: 16px;
	font-weight: 600;
	color: var(--text);
}

.header-spacer {
	width: 44px;
}

/* ── Progress ────────────────────────────────────────────── */
.progress-track {
	width: 100%;
	height: 3px;
	background: var(--border);
}

.progress-fill {
	height: 100%;
	background: var(--green);
	transition: width 0.2s ease;
}

/* ── Content ─────────────────────────────────────────────── */
.content {
	padding: 24px 16px 96px;
	max-width: 480px;
	margin: 0 auto;
}

.step-label {
	font-size: 13px;
	color: var(--text-sub);
	font-weight: 500;
	margin-bottom: 8px;
}

h1 {
	font-size: 24px;
	font-weight: 700;
	color: var(--text);
	margin-bottom: 24px;
	line-height: 1.2;
}

/* ── Date Inputs ─────────────────────────────────────────── */
.date-grid {
	display: grid;
	grid-template-columns: 1fr 1fr;
	gap: 12px;
	margin-bottom: 20px;
}

.date-field {
	background: var(--bg-input);
	border: 1px solid var(--border);
	border-radius: var(--radius-input);
	padding: 10px 14px;
	transition: border-color 0.15s ease;
}

.date-field.active {
	border-color: var(--green);
}

.date-field-label {
	font-size: 12px;
	color: var(--text-sub);
	margin-bottom: 2px;
}

.date-value {
	font-size: 15px;
	font-weight: 600;
	color: var(--text);
}

/* ── Calendar ────────────────────────────────────────────── */
.calendar {
	background: var(--bg-card);
	border: 1px solid var(--border);
	border-radius: var(--radius-card);
	padding: 16px;
	margin-bottom: 20px;
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
	color: var(--text);
	cursor: pointer;
	border-radius: 6px;
}

.cal-nav .material-symbols-outlined {
	font-size: 20px;
}

.cal-title {
	font-size: 14px;
	font-weight: 600;
	color: var(--text);
}

.cal-days-header {
	display: grid;
	grid-template-columns: repeat(7, 1fr);
	margin-bottom: 8px;
}

.cal-days-header span {
	text-align: center;
	font-size: 10px;
	color: var(--text-muted);
	font-weight: 600;
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
	color: var(--text);
	background: none;
	border: none;
	cursor: pointer;
	border-radius: 8px;
	transition: background 0.15s ease, color 0.15s ease;
}

.cal-day.other-month {
	color: var(--text-muted);
}

.cal-day.in-range {
	background: var(--green-soft);
	border-radius: 0;
}

.cal-day.is-start,
.cal-day.is-end {
	background: var(--green);
	color: #fff;
	font-weight: 600;
}

.cal-day.is-start {
	border-radius: 8px 0 0 8px;
}

.cal-day.is-end {
	border-radius: 0 8px 8px 0;
}

.cal-day.is-start.is-end {
	border-radius: 8px;
}

/* ── Visibility Toggle ───────────────────────────────────── */
.visibility-card {
	background: var(--bg-card);
	border: 1px solid var(--border);
	border-radius: var(--radius-card);
	padding: 14px 16px;
	display: flex;
	align-items: center;
	justify-content: space-between;
	gap: 16px;
	margin-bottom: 24px;
}

.vis-left {
	flex: 1;
}

.vis-title {
	font-size: 14px;
	font-weight: 600;
	color: var(--text);
	margin-bottom: 2px;
}

.vis-sub {
	font-size: 13px;
	color: var(--text-sub);
}

.toggle-track {
	position: relative;
	display: flex;
	background: var(--bg-input);
	border: 1px solid var(--border);
	border-radius: var(--radius-pill);
	width: 148px;
	height: 36px;
	padding: 3px;
	overflow: hidden;
	flex-shrink: 0;
}

.toggle-pill {
	position: absolute;
	top: 3px;
	left: 3px;
	width: calc(50% - 3px);
	height: 28px;
	background: var(--green);
	border-radius: var(--radius-pill);
	transition: transform 0.15s ease;
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
	font-size: 12px;
	font-weight: 600;
	cursor: pointer;
	color: var(--text-sub);
	border-radius: var(--radius-pill);
	transition: color 0.15s ease;
}

.toggle-option.active {
	color: #fff;
}

/* ── Continue Button ─────────────────────────────────────── */
.continue-btn {
	width: 100%;
	padding: 14px;
	border-radius: var(--radius-input);
	background: var(--green);
	color: #fff;
	font-size: 16px;
	font-weight: 600;
	border: none;
	cursor: pointer;
}

.continue-btn:disabled {
	opacity: 0.4;
	cursor: not-allowed;
}
</style>
