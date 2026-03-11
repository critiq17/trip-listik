<script lang="ts">
	import type { CountryVisit } from '$lib/types';

	let { countries = [] }: { countries?: CountryVisit[] } = $props();

	const dataByCode = $derived(
		countries.reduce<Record<string, CountryVisit>>((acc, item) => {
			acc[item.code.toUpperCase()] = item;
			return acc;
		}, {})
	);

	const palette = (count: number) => {
		if (count <= 0) return { fill: '#1a1f1c', opacity: 0.3 };
		if (count === 1) return { fill: '#2d6b44', opacity: 0.6 };
		if (count <= 3) return { fill: '#4d9d6d', opacity: 0.8 };
		return { fill: '#7fbf99', opacity: 1 };
	};

	type MapCountry = { code: string; name: string; path: string };

	// Minimal map set (core countries). Extend by adding more paths.
	const MAP: MapCountry[] = [
		{ code: 'US', name: 'United States', path: 'M120 180H320V240H120Z' },
		{ code: 'CA', name: 'Canada', path: 'M120 120H320V175H120Z' },
		{ code: 'MX', name: 'Mexico', path: 'M170 245H280V275H170Z' },
		{ code: 'BR', name: 'Brazil', path: 'M360 260H460V330H360Z' },
		{ code: 'AR', name: 'Argentina', path: 'M360 335H420V385H360Z' },
		{ code: 'GB', name: 'United Kingdom', path: 'M520 170H540V190H520Z' },
		{ code: 'FR', name: 'France', path: 'M540 195H565V215H540Z' },
		{ code: 'DE', name: 'Germany', path: 'M565 190H590V215H565Z' },
		{ code: 'IT', name: 'Italy', path: 'M570 220H600V245H570Z' },
		{ code: 'ES', name: 'Spain', path: 'M520 215H550V235H520Z' },
		{ code: 'RU', name: 'Russia', path: 'M600 120H820V185H600Z' },
		{ code: 'CN', name: 'China', path: 'M700 210H780V255H700Z' },
		{ code: 'JP', name: 'Japan', path: 'M800 220H820V245H800Z' },
		{ code: 'IN', name: 'India', path: 'M670 255H720V295H670Z' },
		{ code: 'AU', name: 'Australia', path: 'M800 330H880V380H800Z' },
		{ code: 'ZA', name: 'South Africa', path: 'M600 350H660V385H600Z' }
	];

	let hovered = $state<MapCountry | null>(null);
	let tooltip = $state({ x: 0, y: 0, show: false });
	let pinned = $state<MapCountry | null>(null);

	const onMove = (event: MouseEvent) => {
		tooltip.x = event.clientX + 12;
		tooltip.y = event.clientY + 12;
	};

	const onEnter = (item: MapCountry, event: MouseEvent) => {
		hovered = item;
		tooltip.show = true;
		onMove(event);
	};

	const onLeave = () => {
		hovered = null;
		if (!pinned) tooltip.show = false;
	};

	const onClick = (item: MapCountry, event: MouseEvent) => {
		pinned = pinned?.code === item.code ? null : item;
		tooltip.show = Boolean(pinned) || Boolean(hovered);
		onMove(event);
	};

	const onKey = (item: MapCountry, event: KeyboardEvent) => {
		if (event.key === 'Enter' || event.key === ' ') {
			event.preventDefault();
			const fakeEvent = new MouseEvent('click', { clientX: tooltip.x, clientY: tooltip.y });
			onClick(item, fakeEvent);
		}
	};

	const getCount = (code: string) => dataByCode[code]?.visit_count ?? 0;
	const getLabel = (item: MapCountry) => {
		const count = getCount(item.code);
		return `${item.name} — ${count} trip${count === 1 ? '' : 's'}`;
	};
</script>

<div class="map-wrap" role="application" onmousemove={onMove}>
	<svg viewBox="0 0 1000 500" role="img" aria-label="World map" class="map">
		<rect x="0" y="0" width="1000" height="500" class="ocean" />
		{#each MAP as item}
			{@const visit = getCount(item.code)}
			{@const style = palette(visit)}
			<path
				class="country"
				d={item.path}
				data-country={item.code}
				style={`fill: ${style.fill}; opacity: ${style.opacity}`}
				role="button"
				tabindex="0"
				aria-label={getLabel(item)}
				onmouseenter={(event) => onEnter(item, event)}
				onmouseleave={onLeave}
				onclick={(event) => onClick(item, event)}
				onkeydown={(event) => onKey(item, event)}
			/>
		{/each}
	</svg>

	{#if tooltip.show && (hovered || pinned)}
		<div class="tooltip" style={`left: ${tooltip.x}px; top: ${tooltip.y}px;`}>
			{getLabel(pinned ?? hovered!)}
		</div>
	{/if}
</div>

<style>
	.map-wrap {
		position: relative;
		width: 100%;
		border-radius: var(--radius-3xl);
		background: radial-gradient(circle at 20% 20%, rgba(77, 157, 109, 0.15), transparent 50%),
			linear-gradient(160deg, #0f1411, #111a16 55%, #0a100d);
		border: 1px solid rgba(255, 255, 255, 0.06);
		padding: 1rem;
		overflow: hidden;
	}

	.map {
		width: 100%;
		height: auto;
		display: block;
	}

	.ocean {
		fill: rgba(15, 20, 17, 0.9);
	}

	.country {
		stroke: rgba(255, 255, 255, 0.08);
		stroke-width: 1;
		transition: fill 0.3s var(--transition-smooth), opacity 0.3s var(--transition-smooth),
			transform 0.3s var(--transition-smooth);
	}

	.country:hover {
		transform: translateY(-2px);
		stroke: rgba(127, 191, 153, 0.6);
	}

	.tooltip {
		position: fixed;
		padding: 0.45rem 0.65rem;
		border-radius: var(--radius-pill);
		background: rgba(8, 13, 11, 0.9);
		color: white;
		font-size: 0.7rem;
		letter-spacing: 0.02em;
		box-shadow: 0 10px 24px rgba(0, 0, 0, 0.45);
		border: 1px solid rgba(255, 255, 255, 0.08);
		z-index: 10;
		pointer-events: none;
	}
</style>
