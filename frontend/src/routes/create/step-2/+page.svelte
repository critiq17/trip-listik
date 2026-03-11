<script lang="ts">
	import { goto } from '$app/navigation';
	import { apiFetch } from '$lib/api';
	import { getDraft } from '$lib/tripDraft';

	let isPublic = true;
	let startDate: Date | null = new Date(2023, 9, 12);
	let endDate: Date | null = new Date(2023, 9, 18);
	let saving = false;
	let error = '';

	let viewMonth = new Date();
	viewMonth.setDate(1);

	const formatDisplay = (value: Date | null) => {
		if (!value) return '';
		return value.toLocaleDateString(undefined, {
			month: 'short',
			day: 'numeric',
			year: 'numeric'
		});
	};

	const toISO = (value: Date | null) => {
		if (!value) return null;
		return value.toISOString().split('T')[0];
	};

	const daysInMonth = (date: Date) => {
		return new Date(date.getFullYear(), date.getMonth() + 1, 0).getDate();
	};

	const startWeekday = (date: Date) => {
		return new Date(date.getFullYear(), date.getMonth(), 1).getDay();
	};

	const buildGrid = (date: Date) => {
		const totalDays = daysInMonth(date);
		const offset = startWeekday(date);
		const cells: Array<Date | null> = [];
		for (let i = 0; i < offset; i += 1) cells.push(null);
		for (let day = 1; day <= totalDays; day += 1) {
			cells.push(new Date(date.getFullYear(), date.getMonth(), day));
		}
		const rows: Array<Array<Date | null>> = [];
		for (let i = 0; i < cells.length; i += 7) {
			rows.push(cells.slice(i, i + 7));
		}
		return rows;
	};

	const isSameDay = (a: Date | null, b: Date | null) => {
		if (!a || !b) return false;
		return a.toDateString() === b.toDateString();
	};

	const isInRange = (day: Date) => {
		if (!startDate || !endDate) return false;
		const time = day.getTime();
		return time >= startDate.getTime() && time <= endDate.getTime();
	};

	const selectDay = (day: Date) => {
		if (!startDate || (startDate && endDate)) {
			startDate = day;
			endDate = null;
			return;
		}
		if (day.getTime() < startDate.getTime()) {
			endDate = startDate;
			startDate = day;
			return;
		}
		endDate = day;
	};

	const next = async () => {
		const draft = getDraft();
		if (!draft) {
			error = 'Draft not found';
			return;
		}
		saving = true;
		error = '';
		try {
			await apiFetch(`/v1/trips/${draft.id}`, {
				method: 'PATCH',
				body: JSON.stringify({
					start_date: toISO(startDate),
					end_date: toISO(endDate),
					visibility: isPublic ? 'public' : 'private'
				})
			});
			goto('/create/step-3');
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to save';
		} finally {
			saving = false;
		}
	};

	const back = () => history.back();

	const toggle = (value: boolean) => {
		isPublic = value;
	};
</script>

<section class="step-two">
	<header class="header">
		<button class="round-btn" onclick={back} aria-label="Back">
			<span class="material-symbols-outlined">arrow_back</span>
		</button>
		<h2>Create Trip</h2>
		<div class="spacer"></div>
	</header>

	<main>
		<div class="progress">
		<div class="progress-row">
			<p>Step 2 of 3</p>
			<span>66% Complete</span>
		</div>
			<div class="progress-bar">
				<div class="progress-fill"></div>
			</div>
		</div>

		<h1>When does the journey begin?</h1>

		<div class="dates">
			<div class="date-field">
				<span>From</span>
				<div class="date-input">
					<span class="material-symbols-outlined">calendar_today</span>
					<input type="text" readonly value={formatDisplay(startDate)} />
				</div>
			</div>
			<div class="date-field">
				<span>To</span>
				<div class="date-input">
					<span class="material-symbols-outlined">event</span>
					<input type="text" readonly value={formatDisplay(endDate)} />
				</div>
			</div>
		</div>

		<div class="calendar">
			<div class="calendar-head">
				<button
					class="material-symbols-outlined"
					aria-label="Previous month"
					onclick={() => (viewMonth = new Date(viewMonth.getFullYear(), viewMonth.getMonth() - 1, 1))}
				>
					chevron_left
				</button>
				<h3>
					{viewMonth.toLocaleDateString(undefined, { month: 'long', year: 'numeric' })}
				</h3>
				<button
					class="material-symbols-outlined"
					aria-label="Next month"
					onclick={() => (viewMonth = new Date(viewMonth.getFullYear(), viewMonth.getMonth() + 1, 1))}
				>
					chevron_right
				</button>
			</div>
			<div class="calendar-week">
				<span>SU</span><span>MO</span><span>TU</span><span>WE</span><span>TH</span><span>FR</span><span>SA</span>
			</div>
			<div class="calendar-grid">
				{#each buildGrid(viewMonth) as row}
					{#each row as day}
						{#if day}
							<button
								class:range={isInRange(day)}
								class:start={isSameDay(day, startDate)}
								class:end={isSameDay(day, endDate)}
								onclick={() => selectDay(day)}
							>
								{day.getDate()}
							</button>
						{:else}
							<span class="muted"></span>
						{/if}
					{/each}
				{/each}
			</div>
		</div>

		<div class="visibility">
			<div>
				<strong>Trip Visibility</strong>
				<p>Who can see your itinerary?</p>
			</div>
			<div class="toggle" class:private={!isPublic}>
				<div class="toggle-pill"></div>
				<button class:is-active={isPublic} onclick={() => toggle(true)}>Public</button>
				<button class:is-active={!isPublic} onclick={() => toggle(false)}>Private</button>
			</div>
		</div>

		<button class="continue" onclick={next} disabled={saving || !startDate || !endDate}>
			{saving ? 'Saving…' : 'Continue'}
			<span class="material-symbols-outlined">arrow_forward</span>
		</button>

		{#if error}
			<p class="error">{error}</p>
		{/if}
	</main>
</section>

<style>
	.step-two {
		min-height: 100dvh;
		background: var(--background-dark);
		color: var(--text-primary);
	}

	.header {
		position: sticky;
		top: 0;
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 0.9rem 1rem;
		border-bottom: 1px solid rgba(77, 157, 109, 0.15);
		background: var(--background-dark);
	}

	.header h2 {
		font-size: 1rem;
		font-weight: 700;
	}

	.spacer {
		width: 2.6rem;
	}

	.round-btn {
		width: 2.6rem;
		height: 2.6rem;
		border-radius: 999px;
		display: grid;
		place-items: center;
	}

	main {
		max-width: 480px;
		margin: 0 auto;
		padding: 1rem 1.5rem 3rem;
	}

	.progress {
		margin: 1rem 0 1.5rem;
	}

	.progress-row {
		display: flex;
		align-items: flex-end;
		justify-content: space-between;
		margin-bottom: 0.6rem;
	}

	.progress-row p {
		color: var(--primary);
		font-weight: 700;
		font-size: 0.7rem;
		letter-spacing: 0.2em;
		text-transform: uppercase;
	}

	.progress-row span {
		color: var(--text-secondary);
		font-size: 0.7rem;
	}

	.progress-bar {
		height: 8px;
		background: rgba(255, 255, 255, 0.12);
		border-radius: 999px;
		overflow: hidden;
	}

	.progress-fill {
		height: 100%;
		width: 66%;
		background: var(--primary);
		border-radius: inherit;
	}

	h1 {
		font-size: 1.4rem;
		font-weight: 700;
		margin-bottom: 1.5rem;
	}

	.dates {
		display: grid;
		grid-template-columns: repeat(2, 1fr);
		gap: 1rem;
		margin-bottom: 1.5rem;
	}

	.date-field span {
		font-size: 0.6rem;
		color: var(--text-secondary);
		text-transform: uppercase;
		letter-spacing: 0.1em;
		font-weight: 700;
	}

	.date-input {
		display: flex;
		align-items: center;
		gap: 0.4rem;
		border-bottom: 2px solid rgba(77, 157, 109, 0.4);
		padding: 0.4rem 0.2rem;
	}

	.date-input input {
		background: transparent;
		border: none;
		font-size: 0.9rem;
		font-weight: 600;
		color: var(--text-primary);
		width: 100%;
	}

	.calendar {
		background: rgba(255, 255, 255, 0.03);
		border: 1px solid rgba(77, 157, 109, 0.2);
		border-radius: 12px;
		padding: 1rem;
		margin-bottom: 1.5rem;
		box-shadow: var(--shadow-soft);
	}

	.calendar-head {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: 0.8rem;
	}

	.calendar-head h3 {
		font-size: 0.85rem;
		font-weight: 700;
	}

		.calendar-week,
		.calendar-grid {
			display: grid;
			grid-template-columns: repeat(7, 1fr);
			text-align: center;
			gap: 0.2rem;
		}

	.calendar-week span {
		font-size: 0.6rem;
		color: var(--text-secondary);
		font-weight: 700;
	}

		.calendar-grid span,
		.calendar-grid button {
			font-size: 0.7rem;
			padding: 0.3rem 0;
			border-radius: 6px;
			background: none;
			border: none;
			color: var(--text-primary);
		}

		.calendar-grid span.muted {
			color: var(--text-secondary);
		}

		.calendar-grid button.range {
			background: rgba(77, 157, 109, 0.2);
		}

		.calendar-grid button.start,
		.calendar-grid button.end {
			background: var(--primary);
			color: white;
		}

		.calendar-grid button.range {
			position: relative;
		}

		.calendar-grid button.start::after,
		.calendar-grid button.end::after {
			content: '';
			position: absolute;
			inset: 2px;
			border-radius: 6px;
			border: 1px solid rgba(255, 255, 255, 0.4);
		}

	.visibility {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 1rem;
		padding: 1rem;
		background: rgba(77, 157, 109, 0.1);
		border: 1px solid rgba(77, 157, 109, 0.2);
		border-radius: 12px;
		margin-bottom: 1.5rem;
	}

	.visibility strong {
		font-size: 0.85rem;
	}

	.visibility p {
		font-size: 0.7rem;
		color: var(--text-secondary);
	}

	.toggle {
		position: relative;
		display: grid;
		grid-template-columns: repeat(2, 1fr);
		background: rgba(0, 0, 0, 0.3);
		border-radius: 999px;
		padding: 0.2rem;
		width: 8rem;
		height: 2.2rem;
	}

	.toggle-pill {
		position: absolute;
		top: 0.2rem;
		left: 0.2rem;
		width: calc(50% - 0.2rem);
		height: calc(100% - 0.4rem);
		background: var(--primary);
		border-radius: 999px;
		transition: transform 0.25s ease;
	}

	.toggle.private .toggle-pill {
		transform: translateX(100%);
	}

	.toggle button {
		z-index: 1;
		font-size: 0.6rem;
		font-weight: 700;
		text-transform: uppercase;
		color: var(--text-secondary);
	}

	.toggle button.is-active {
		color: white;
	}

	.continue {
		width: 100%;
		padding: 1rem;
		border-radius: 12px;
		background: var(--primary);
		color: white;
		font-weight: 700;
		display: inline-flex;
		align-items: center;
		justify-content: center;
		gap: 0.4rem;
		box-shadow: 0 14px 30px rgba(77, 157, 109, 0.35);
		margin-bottom: 1rem;
	}

	.error {
		color: #e11d48;
		font-size: 0.8rem;
	}
</style>
