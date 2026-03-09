<script lang="ts">
	import { onMount } from 'svelte';
	import FilterChips from '$lib/components/FilterChips.svelte';
	import TripCard from '$lib/components/TripCard.svelte';
	import SkeletonCard from '$lib/components/SkeletonCard.svelte';
	import { apiFetch } from '$lib/api';
	import { expandTelegram } from '$lib/telegram';
	import type { TripCardData } from '$lib/types';

	const filters = ['All', 'Friends', 'Popular', 'Nearby'];
	let active = 'All';
	let items: TripCardData[] = [];
	let loading = true;
	let error = '';
	let ready = false;

	async function loadFeed() {
		loading = true;
		error = '';
		const filter = active.toLowerCase();
		const query = filter === 'all' ? '' : `?filter=${filter}`;
		try {
			const data = await apiFetch<{ items: TripCardData[] }>(`/v1/feed${query}`);
			items = data.items ?? [];
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load feed';
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		expandTelegram();
		ready = true;
	});

	$: if (ready && active) {
		loadFeed();
	}
</script>

<section class="container">
	<header class="header">
		<div class="copy">
			<p class="eyebrow">Public Feed</p>
			<h1 class="headline">Find trips worth joining.</h1>
			<p class="subtle">Live public journeys, ranked by recency and engagement.</p>
		</div>
		<button class="icon-btn" aria-label="Search">
			<span>⌕</span>
		</button>
	</header>

	<FilterChips {filters} bind:active />

	<div class="hero-strip glass">
		<div>
			<p class="hero-label">Trending now</p>
			<strong>{items.length || 0} public trips ready to explore</strong>
		</div>
		<span class="hero-badge">{active}</span>
	</div>

	<div class="feed">
		{#if loading}
			{#each Array(3) as _}
				<SkeletonCard />
			{/each}
		{:else if error}
			<div class="error glass">{error}</div>
		{:else if items.length === 0}
			<div class="empty glass">No trips yet</div>
		{:else}
			{#each items as trip}
				<TripCard {trip} />
			{/each}
		{/if}
	</div>
</section>

<style>
	.header {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		gap: 1rem;
		margin-bottom: 1rem;
	}

	.icon-btn {
		width: 3rem;
		height: 3rem;
		border-radius: 50%;
		background: rgba(255, 255, 255, 0.06);
		font-size: 1.2rem;
		border: 1px solid rgba(255, 255, 255, 0.08);
		backdrop-filter: blur(16px);
		flex-shrink: 0;
	}

	.copy {
		display: grid;
		gap: 0.45rem;
	}

	.hero-strip {
		display: flex;
		justify-content: space-between;
		align-items: center;
		gap: 1rem;
		padding: 1rem 1.1rem;
		border-radius: var(--radius-xl);
		margin: 0.15rem 0 1.15rem;
	}

	.hero-label {
		font-size: 0.74rem;
		letter-spacing: 0.12em;
		text-transform: uppercase;
		color: var(--text-muted);
		margin-bottom: 0.25rem;
	}

	.hero-strip strong {
		font-size: 0.95rem;
	}

	.hero-badge {
		padding: 0.6rem 0.9rem;
		border-radius: var(--radius-pill);
		background: rgba(32, 146, 186, 0.16);
		border: 1px solid rgba(122, 234, 244, 0.24);
		color: var(--accent-strong);
		font-size: 0.75rem;
		font-weight: 800;
		text-transform: uppercase;
	}

	.feed {
		display: grid;
		gap: 1rem;
	}

	.empty,
	.error {
		padding: 1.5rem;
		border-radius: var(--radius-xl);
		color: var(--text-secondary);
		text-align: center;
	}
</style>
