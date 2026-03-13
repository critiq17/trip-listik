<script lang="ts">
	/**
	 * DateRangePicker — custom calendar per spec (no native <input type="date">).
	 * Features: month navigation, range highlight, past dates greyed out.
	 */

	let {
		startDate = $bindable(''),
		endDate   = $bindable('')
	}: {
		startDate?: string;
		endDate?: string;
	} = $props();

	const today = new Date();
	today.setHours(0, 0, 0, 0);

	let viewYear  = $state(today.getFullYear());
	let viewMonth = $state(today.getMonth());   // 0–11
	let selecting: 'start' | 'end' = 'start';
	let hoverDate = $state<string>('');

	const DAYS   = ['S', 'M', 'T', 'W', 'T', 'F', 'S'];
	const MONTHS = [
		'January', 'February', 'March', 'April', 'May', 'June',
		'July', 'August', 'September', 'October', 'November', 'December'
	];

	// Generate all cells for the visible month (including leading blanks)
	const cells = $derived((() => {
		const firstDay = new Date(viewYear, viewMonth, 1).getDay(); // 0=Sun
		const daysInMonth = new Date(viewYear, viewMonth + 1, 0).getDate();
		const result: Array<{ date: string; day: number } | null> = [];
		for (let i = 0; i < firstDay; i++) result.push(null);
		for (let d = 1; d <= daysInMonth; d++) {
			const date = iso(viewYear, viewMonth + 1, d);
			result.push({ date, day: d });
		}
		return result;
	})());

	function iso(y: number, m: number, d: number): string {
		return `${y}-${String(m).padStart(2, '0')}-${String(d).padStart(2, '0')}`;
	}

	function isPast(date: string): boolean {
		return date < iso(today.getFullYear(), today.getMonth() + 1, today.getDate());
	}

	function isStart(date: string) { return date === startDate; }
	function isEnd(date: string)   { return date === endDate; }

	function inRange(date: string): boolean {
		const s = startDate;
		const e = endDate || hoverDate;
		if (!s) return false;
		if (!e) return false;
		const lo = s < e ? s : e;
		const hi = s < e ? e : s;
		return date > lo && date < hi;
	}

	function select(date: string) {
		if (isPast(date)) return;
		if (selecting === 'start' || !startDate) {
			startDate = date;
			endDate   = '';
			selecting = 'end';
		} else {
			if (date < startDate) {
				endDate   = startDate;
				startDate = date;
			} else {
				endDate = date;
			}
			selecting = 'start';
		}
	}

	function prevMonth() {
		if (viewMonth === 0) { viewMonth = 11; viewYear--; }
		else viewMonth--;
	}

	function nextMonth() {
		if (viewMonth === 11) { viewMonth = 0; viewYear++; }
		else viewMonth++;
	}

	const label = $derived(
		startDate && endDate
			? `${startDate} → ${endDate}`
			: startDate
			? `${startDate} → pick end date`
			: 'Pick start date'
	);
</script>

<div class="picker">
	<!-- Month header -->
	<div class="header">
		<button type="button" class="nav" onclick={prevMonth} aria-label="Previous month">
			<span class="material-symbols-outlined">chevron_left</span>
		</button>
		<span class="month-label">{MONTHS[viewMonth]} {viewYear}</span>
		<button type="button" class="nav" onclick={nextMonth} aria-label="Next month">
			<span class="material-symbols-outlined">chevron_right</span>
		</button>
	</div>

	<!-- Day-of-week headers -->
	<div class="grid days-header">
		{#each DAYS as d}
			<span class="dow">{d}</span>
		{/each}
	</div>

	<!-- Date cells -->
	<div class="grid" role="grid" aria-label="Calendar">
		{#each cells as cell}
			{#if cell === null}
				<div></div>
			{:else}
				<button
					type="button"
					class="cell"
					class:past={isPast(cell.date)}
					class:is-start={isStart(cell.date)}
					class:is-end={isEnd(cell.date)}
					class:in-range={inRange(cell.date)}
					class:today={cell.date === iso(today.getFullYear(), today.getMonth() + 1, today.getDate())}
					onclick={() => select(cell.date)}
					onmouseenter={() => (hoverDate = cell.date)}
					onmouseleave={() => (hoverDate = '')}
					disabled={isPast(cell.date)}
					aria-label={cell.date}
					aria-selected={isStart(cell.date) || isEnd(cell.date)}
					role="gridcell"
				>
					{cell.day}
				</button>
			{/if}
		{/each}
	</div>

	<p class="range-label">{label}</p>
</div>

<style>
	.picker {
		background: var(--bg-elevated);
		border: 1px solid var(--border);
		border-radius: var(--radius-card);
		padding: 1rem;
		user-select: none;
	}

	.header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: 0.75rem;
	}

	.month-label {
		font-size: 0.95rem;
		font-weight: 700;
		color: var(--text);
	}

	.nav {
		width: 2rem;
		height: 2rem;
		border-radius: 8px;
		display: grid;
		place-items: center;
		color: var(--text-sub);
		transition: background 0.15s;
	}

	.nav:hover {
		background: rgba(255, 255, 255, 0.06);
	}

	.nav .material-symbols-outlined {
		font-size: 1.25rem;
	}

	.grid {
		display: grid;
		grid-template-columns: repeat(7, 1fr);
		gap: 2px;
	}

	.days-header {
		margin-bottom: 4px;
	}

	.dow {
		text-align: center;
		font-size: 0.65rem;
		color: var(--text-muted);
		font-weight: 600;
		padding: 0.25rem 0;
	}

	.cell {
		aspect-ratio: 1;
		border-radius: 8px;
		font-size: 0.82rem;
		font-weight: 500;
		color: var(--text);
		display: grid;
		place-items: center;
		transition: background 0.12s, color 0.12s;
		cursor: pointer;
	}

	.cell:hover:not(.past):not(.is-start):not(.is-end) {
		background: rgba(61, 158, 95, 0.12);
	}

	.cell.past {
		color: var(--text-muted);
		cursor: default;
		opacity: 0.4;
	}

	.cell.today {
		border: 1px solid var(--border);
	}

	.cell.in-range {
		background: rgba(61, 158, 95, 0.12);
		border-radius: 0;
	}

	.cell.is-start,
	.cell.is-end {
		background: var(--green);
		color: var(--bg);
		font-weight: 700;
		border-radius: 8px;
	}

	.range-label {
		margin-top: 0.75rem;
		text-align: center;
		font-size: 0.78rem;
		color: var(--text-sub);
		min-height: 1.2em;
	}
</style>
