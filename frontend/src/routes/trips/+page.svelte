<script lang="ts">
	import { onMount } from 'svelte';
	import SkeletonCard from '$lib/components/SkeletonCard.svelte';
	import TripCard from '$lib/components/TripCard.svelte';
	import { apiFetch } from '$lib/api';
	import { staggerList, scalePress } from '$lib/actions/animate';
	import { animate } from 'motion';
	import type { TripCardData } from '$lib/types';

	const tabs = ['Upcoming', 'Past', 'Drafts'];
	let active = $state('Upcoming');
	let items = $state<TripCardData[]>([]);
	let loading = $state(true);
	let error = $state('');
	let ready = $state(false);
	let tabBar = $state<HTMLElement | null>(null);
	let indicator = $state<HTMLDivElement | null>(null);
	let tabRefs = $state<HTMLButtonElement[]>([]);

	async function loadTrips() {
		loading = true;
		error = '';
		const status = active.toLowerCase();
		try {
			const data = await apiFetch<{ items: TripCardData[] }>(
				`/v1/trips?scope=mine&status=${status}`
			);
			items = data.items ?? [];
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load trips';
		} finally {
			loading = false;
		}
	}

	const moveIndicator = () => {
		if (!tabBar || !indicator) return;
		const index = tabs.findIndex((tab) => tab === active);
		const target = tabRefs[index];
		if (!target) return;
		const barRect = tabBar.getBoundingClientRect();
		const tabRect = target.getBoundingClientRect();
		const left = tabRect.left - barRect.left;
		animate(
			indicator,
			{
				x: [`${indicator.offsetLeft}px`, `${left}px`],
				width: [`${indicator.offsetWidth}px`, `${tabRect.width}px`]
			} as any,
			{ duration: 0.25, easing: [0.25, 0.46, 0.45, 0.94] } as any
		);
		indicator.style.transform = `translateX(${left}px)`;
		indicator.style.width = `${tabRect.width}px`;
	};

	onMount(() => {
		ready = true;
		setTimeout(moveIndicator, 0);
	});

	$effect(() => {
		if (ready && active) {
			loadTrips();
			setTimeout(moveIndicator, 0);
		}
	});
</script>

<section class="trips-page">
	<header class="header">
		<div class="header-row">
			<button class="round-btn" aria-label="Back">
				<span class="material-symbols-outlined">arrow_back</span>
			</button>
			<button class="round-btn" aria-label="More">
				<span class="material-symbols-outlined">more_horiz</span>
			</button>
		</div>
		<div class="header-copy">
			<span class="eyebrow">My Trips</span>
			<h1>Your travel board.</h1>
			<p>Manage your upcoming and past adventures.</p>
		</div>
	</header>

	<nav class="tabs" bind:this={tabBar}>
		{#each tabs as tab, i}
			<button bind:this={tabRefs[i]} class:active={tab === active} onclick={() => (active = tab)}>
				{tab}
			</button>
		{/each}
		<div class="indicator" bind:this={indicator}></div>
	</nav>

	<main class="content">
		<div class="summary">
			<span>{items.length} visible · {active}</span>
		</div>

		<div class="cards" use:staggerList>
			{#if loading}
				{#each Array(1) as _}
					<SkeletonCard ratio="16 / 10" />
				{/each}
			{:else if error}
				<div class="state">{error}</div>
			{:else if items.length === 0}
				<div class="empty">
					<div class="empty-icon">
						<span class="material-symbols-outlined">edit_note</span>
					</div>
					<div>
						<p>No active drafts</p>
						<span>Start planning your next escape today.</span>
					</div>
					<button class="empty-action">
						<span class="material-symbols-outlined">add</span>
						Create new trip
					</button>
				</div>
			{:else}
				{#each items as trip (trip.id)}
					<div data-item use:scalePress>
						<TripCard {trip} variant="compact" />
					</div>
				{/each}
			{/if}
		</div>

	</main>
</section>

<style>
	.trips-page {
		min-height: 100dvh;
		background: var(--background-dark);
		color: var(--text-primary);
		padding-bottom: 5.5rem;
	}

	.header {
		padding: 1.5rem 1rem 0.5rem;
		max-width: 480px;
		margin: 0 auto;
	}

	.header-row {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 1rem;
	}

	.round-btn {
		width: 2.6rem;
		height: 2.6rem;
		border-radius: 999px;
		display: grid;
		place-items: center;
		color: var(--text-primary);
	}

	.header-copy span {
		color: var(--primary);
		font-size: 0.6rem;
		letter-spacing: 0.2em;
		font-weight: 700;
		text-transform: uppercase;
	}

	.header-copy h1 {
		font-size: 1.9rem;
		font-weight: 700;
		margin: 0.4rem 0 0.35rem;
	}

	.header-copy p {
		color: var(--text-secondary);
		font-size: 0.88rem;
	}

	.tabs {
		position: relative;
		display: flex;
		gap: 1.5rem;
		padding: 0 1rem;
		border-bottom: 1px solid rgba(255, 255, 255, 0.06);
		max-width: 480px;
		margin: 0 auto;
	}

	.tabs button {
		padding: 0.85rem 0 0.9rem;
		background: none;
		border: none;
		color: var(--text-secondary);
		font-size: 0.85rem;
		font-weight: 600;
	}

	.tabs button.active {
		color: var(--primary);
	}

	.indicator {
		position: absolute;
		bottom: 0;
		left: 0;
		height: 2px;
		background: var(--primary);
		border-radius: 999px;
		width: 40px;
	}

	.content {
		padding: 1rem;
		max-width: 480px;
		margin: 0 auto;
	}

	.summary span {
		display: inline-block;
		padding: 0.35rem 0.75rem;
		border-radius: var(--radius-pill);
		background: rgba(77, 157, 109, 0.15);
		color: var(--primary);
		font-size: 0.72rem;
		font-weight: 600;
	}

	.cards {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
		margin-top: 1rem;
	}

	.state {
		padding: 1.5rem;
		text-align: center;
		background: rgba(255, 255, 255, 0.04);
		border-radius: 12px;
		color: var(--text-secondary);
	}

	.empty {
		border: 2px dashed rgba(255, 255, 255, 0.12);
		border-radius: 12px;
		padding: 2rem 1.5rem;
		text-align: center;
		color: var(--text-secondary);
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}

	.empty-icon {
		width: 3rem;
		height: 3rem;
		border-radius: 999px;
		background: rgba(255, 255, 255, 0.08);
		display: grid;
		place-items: center;
		margin: 0 auto;
		color: rgba(255, 255, 255, 0.6);
	}

	.empty p {
		font-weight: 600;
		color: var(--text-primary);
	}

	.empty span {
		font-size: 0.75rem;
		color: var(--text-secondary);
	}

	.empty-action {
		color: var(--primary);
		font-weight: 700;
		font-size: 0.8rem;
		display: inline-flex;
		align-items: center;
		gap: 0.3rem;
		justify-content: center;
	}

	.empty-action .material-symbols-outlined {
		font-size: 1rem;
	}

</style>
