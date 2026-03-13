<script lang="ts">
	import { onMount } from 'svelte';
	import FilterChips from '$lib/components/FilterChips.svelte';
	import TripCard from '$lib/components/TripCard.svelte';
	import SkeletonCard from '$lib/components/SkeletonCard.svelte';
	import { apiFetch } from '$lib/api';
	import { expandTelegram } from '$lib/telegram';
	import { staggerList, scalePress } from '$lib/actions/animate';
	import type { TripCardData } from '$lib/types';

	const filters = ['All', 'Friends', 'Popular', 'Nearby'];
	let active = $state('All');
	let searchQuery = $state('');
	let items = $state<TripCardData[]>([]);
	let loading = $state(true);
	let error = $state('');
	let ready = $state(false);
	let cursor = $state<string | null>(null);
	let loadingMore = $state(false);
	let hasMore = $state(true);
	let searchTimer: ReturnType<typeof setTimeout> | null = null;
	let sentinel = $state<HTMLDivElement | null>(null);

	const fetchFeed = async (mode: 'reset' | 'append' = 'reset') => {
		if (mode === 'append' && (loadingMore || !hasMore)) return;
		if (mode === 'append') loadingMore = true;
		loading = mode === 'reset';
		error = '';
		const filter = active.toLowerCase();
		const queryArgs = [];
		if (searchQuery.trim().length > 0) {
			queryArgs.push(`q=${encodeURIComponent(searchQuery)}`);
		} else if (filter !== 'all') {
			queryArgs.push(`filter=${filter}`);
		}
		
		if (mode === 'append' && cursor) {
			queryArgs.push(`cursor=${cursor}`);
		}
		const queryString = queryArgs.join('&');
		const endpoint = searchQuery.trim().length > 0 ? '/v1/explore' : '/v1/feed';

		try {
			const data = await apiFetch<{ items: TripCardData[]; next_cursor?: string; cursor?: string }>(
				`${endpoint}${queryString ? `?${queryString}` : ''}`
			);
			const next = data.next_cursor ?? data.cursor ?? null;
			if (mode === 'append') {
				items = [...items, ...(data.items ?? [])];
			} else {
				items = data.items ?? [];
			}
			cursor = next;
			hasMore = Boolean(next);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load feed';
		} finally {
			loading = false;
			loadingMore = false;
		}
	};

	const observe = () => {
		if (!sentinel) return;
		const io = new IntersectionObserver((entries) => {
			if (entries[0]?.isIntersecting) {
				fetchFeed('append');
			}
		});
		io.observe(sentinel);
		return () => io.disconnect();
	};

	onMount(() => {
		expandTelegram();
		ready = true;
		return observe();
	});

	const debounceSearch = () => {
		if (searchTimer) clearTimeout(searchTimer);
		searchTimer = setTimeout(() => {
			cursor = null;
			hasMore = true;
			fetchFeed('reset');
		}, 300);
	};

	$effect(() => {
		if (ready && active) {
			cursor = null;
			hasMore = true;
			fetchFeed('reset');
		}
	});

	$effect(() => {
		if (ready) {
			searchQuery;
			debounceSearch();
		}
	});
</script>

<section class="feed-page">

	<main class="content">
		<div class="headline">
			<p class="eyebrow">Public Feed</p>
			<h1 class="serif-text">Find trips worth joining.</h1>
			<p class="subtext">Curated journeys from the community.</p>
		</div>
		
		<div class="search-wrap">
			<span class="material-symbols-outlined search-icon">search</span>
			<input 
				type="text" 
				placeholder="Search destinations, tags..." 
				bind:value={searchQuery}
			/>
		</div>

		{#if !searchQuery.trim()}
			<FilterChips {filters} bind:active />

			<div class="trending">
				<div class="trend-left">
					<span class="material-symbols-outlined">trending_up</span>
					<p>Trending Now: <strong>#Iceland2026</strong></p>
				</div>
				<span class="trend-count">(1.2k)</span>
			</div>
		{/if}

		<div class="cards" use:staggerList>
			{#if loading}
				{#each Array(2) as _}
					<SkeletonCard ratio="3 / 4" />
				{/each}
			{:else if error}
				<div class="state error">{error}</div>
			{:else if items.length === 0}
				<div class="state empty">No trips yet</div>
			{:else}
				{#each items as trip (trip.id)}
					<div data-item use:scalePress>
						<TripCard {trip} />
					</div>
				{/each}
			{/if}
			<div bind:this={sentinel} class="sentinel" aria-hidden="true"></div>
		</div>
	</main>
</section>

<style>
	.feed-page {
		min-height: 100dvh;
		background: var(--bg);
		color: var(--text);
		padding-bottom: 5.5rem;
	}

	.content {
		padding: 1.5rem 1rem 1.25rem;
		max-width: 480px;
		margin: 0 auto;
	}

	.headline {
		margin-bottom: 0.9rem;
	}

	.eyebrow {
		color: var(--primary);
		font-size: 0.65rem;
		font-weight: 700;
		letter-spacing: 0.2em;
		text-transform: uppercase;
		margin-bottom: 0.5rem;
	}

	h1 {
		font-size: 2.25rem;
		line-height: 1.05;
		margin-bottom: 0.5rem;
	}

	.subtext {
		color: var(--text-secondary);
		font-size: 0.95rem;
	}

	.search-wrap {
		position: relative;
		margin-bottom: 2rem;
		margin-top: 1.5rem;
		width: 80%;
		transition: width 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
	}

	.search-wrap:focus-within {
		width: 100%;
	}

	.search-wrap input {
		width: 100%;
		padding: 1rem 1rem 1rem 3rem;
		border-radius: var(--radius-xl);
		background: rgba(255, 255, 255, 0.04);
		border: 1px solid rgba(255, 255, 255, 0.08);
		color: var(--text-primary);
		font-size: 1.05rem;
		outline: none;
		transition: border-color 0.2s;
	}

	.search-wrap input:focus {
		border-color: var(--primary);
	}

	.search-icon {
		position: absolute;
		left: 1rem;
		top: 50%;
		transform: translateY(-50%);
		color: var(--text-secondary);
		font-size: 1.2rem;
	}

	.trending {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 0.9rem 1rem;
		background: rgba(77, 157, 109, 0.12);
		border-radius: 12px;
		margin: 0.8rem 0 2rem;
	}

	.trend-left {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		color: var(--primary);
		font-weight: 600;
		font-size: 0.85rem;
	}

	.trend-count {
		color: rgba(77, 157, 109, 0.7);
		font-weight: 700;
		font-size: 0.75rem;
	}

	.cards {
		display: flex;
		flex-direction: column;
		gap: 2rem;
	}

	.state {
		border-radius: 12px;
		padding: 1.4rem;
		text-align: center;
		background: rgba(255, 255, 255, 0.04);
		color: var(--text-secondary);
	}

	.sentinel {
		height: 1px;
	}
</style>
