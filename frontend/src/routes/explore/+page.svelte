<script lang="ts">
	import { onMount } from 'svelte';
	import TripCard from '$lib/components/TripCard.svelte';
	import SkeletonCard from '$lib/components/SkeletonCard.svelte';
	import { apiFetch } from '$lib/api';
	import type { TripCardData } from '$lib/types';

	let query = '';
	let country = '';
	let items: TripCardData[] = [];
	let loading = true;
	let error = '';

	async function loadExplore() {
		loading = true;
		error = '';
		try {
			const params = new URLSearchParams();
			if (query) params.set('q', query);
			if (country) params.set('country', country);
			const suffix = params.size ? `?${params.toString()}` : '';
			const data = await apiFetch<{ items: TripCardData[] }>(`/v1/explore${suffix}`);
			items = data.items ?? [];
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load explore';
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		loadExplore();
	});
</script>

<section class="container">
	<header class="header">
		<p class="eyebrow">Discover</p>
		<h1 class="headline">Find the next route.</h1>
		<p class="subtle">Search public trips by destination or narrow the feed by country.</p>
	</header>

	<div class="search panel">
		<input
			placeholder="Search destinations or trip names"
			bind:value={query}
			on:keydown={(e) => e.key === 'Enter' && loadExplore()}
		/>
		<input
			class="country"
			placeholder="Country code"
			bind:value={country}
			maxlength="2"
			on:keydown={(e) => e.key === 'Enter' && loadExplore()}
		/>
		<button class="search-btn" on:click={loadExplore}>Search</button>
	</div>

	<div class="grid">
		{#if loading}
			{#each Array(3) as _}
				<SkeletonCard />
			{/each}
		{:else if error}
			<div class="error glass">{error}</div>
		{:else if items.length === 0}
			<div class="empty glass">No trips found</div>
		{:else}
			{#each items as trip}
				<TripCard {trip} />
			{/each}
		{/if}
	</div>
</section>

<style>
	.header {
		display: grid;
		gap: 0.45rem;
		margin-bottom: 1rem;
	}

	.search {
		display: grid;
		grid-template-columns: minmax(0, 1fr) 110px auto;
		gap: 0.75rem;
		padding: 1rem;
		border-radius: var(--radius-xl);
		margin-bottom: 1rem;
	}

	input {
		min-height: 52px;
		padding: 0.8rem 1rem;
		border-radius: var(--radius-lg);
		border: 1px solid rgba(255, 255, 255, 0.08);
		background: rgba(255, 255, 255, 0.04);
	}

	.country {
		text-transform: uppercase;
	}

	.search-btn {
		min-height: 52px;
		padding: 0 1.2rem;
		border-radius: var(--radius-lg);
		background: var(--accent-grad);
		font-weight: 800;
		color: white;
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

	@media (max-width: 640px) {
		.search {
			grid-template-columns: 1fr;
		}
	}
</style>
