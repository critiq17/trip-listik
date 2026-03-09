<script lang="ts">
	import { onMount } from 'svelte';
	import TripCard from '$lib/components/TripCard.svelte';
	import SkeletonCard from '$lib/components/SkeletonCard.svelte';
	import { apiFetch } from '$lib/api';
	import type { TripCardData } from '$lib/types';

	const tabs = ['Upcoming', 'Past', 'Drafts'];
	let active = 'Upcoming';
	let items: TripCardData[] = [];
	let loading = true;
	let error = '';
	let ready = false;

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

	onMount(() => {
		ready = true;
	});

	$: if (ready && active) {
		loadTrips();
	}
</script>

<section class="container">
	<header class="header">
		<div class="copy">
			<p class="eyebrow">My Trips</p>
			<h1 class="headline">Your travel board.</h1>
			<p class="subtle">Own trips, joined plans and drafts in one place.</p>
		</div>
	</header>

	<div class="tabs">
		{#each tabs as tab}
			<button class:active={tab === active} on:click={() => (active = tab)}>{tab}</button>
		{/each}
	</div>

	<div class="summary glass">
		<div>
			<p class="label">Visible now</p>
			<strong>{items.length}</strong>
		</div>
		<div>
			<p class="label">Section</p>
			<strong>{active}</strong>
		</div>
	</div>

	<div class="grid">
		{#if loading}
			{#each Array(2) as _}
				<SkeletonCard />
			{/each}
		{:else if error}
			<div class="error glass">{error}</div>
		{:else if items.length === 0}
			<div class="empty glass">No trips in this section</div>
		{:else}
			{#each items as trip}
				<TripCard {trip} />
			{/each}
		{/if}
	</div>
</section>

<style>
	.header {
		margin-bottom: 1rem;
	}

	.tabs {
		display: flex;
		gap: 0.6rem;
		margin-bottom: 1rem;
		overflow-x: auto;
		scrollbar-width: none;
	}

	.tabs::-webkit-scrollbar {
		display: none;
	}

	button {
		padding: 0.6rem 1rem;
		border-radius: var(--radius-pill);
		background: rgba(255, 255, 255, 0.05);
		color: var(--text-secondary);
		font-weight: 700;
		font-size: 0.8rem;
		border: 1px solid rgba(255, 255, 255, 0.06);
		white-space: nowrap;
	}

	button.active {
		background: rgba(32, 146, 186, 0.18);
		color: white;
		border: 1px solid rgba(122, 234, 244, 0.4);
	}

	.summary {
		display: grid;
		grid-template-columns: repeat(2, 1fr);
		gap: 0.75rem;
		padding: 1rem 1.15rem;
		border-radius: var(--radius-xl);
		margin-bottom: 1rem;
	}

	.label {
		color: var(--text-muted);
		font-size: 0.72rem;
		text-transform: uppercase;
		letter-spacing: 0.1em;
		margin-bottom: 0.25rem;
	}

	.summary strong {
		font-size: 1.1rem;
	}

	.grid {
		display: grid;
		gap: 1rem;
	}

	.empty,
	.error {
		padding: 1.5rem;
		border-radius: var(--radius-xl);
		text-align: center;
		color: var(--text-secondary);
	}
</style>
