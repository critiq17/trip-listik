<script lang="ts">
	import { onMount } from 'svelte';
	import { fly } from 'svelte/transition';
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
			if (entries[0]?.isIntersecting) fetchFeed('append');
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

<svelte:head>
	<title>TripListik — Feed</title>
	<meta name="description" content="Browse curated trips and travel journeys from the community." />
</svelte:head>

<div class="page">
	<!-- Sticky Header -->
	<header class="top-bar">
		<button class="icon-btn" aria-label="Open menu">
			<span class="material-symbols-outlined">menu</span>
		</button>
		<span class="logo">TripListik</span>
		<button class="icon-btn" aria-label="Notifications">
			<span class="material-symbols-outlined">notifications</span>
		</button>
	</header>

	<main class="content">
		<!-- Hero -->
		<section class="hero">
			<p class="eyebrow">Public Feed</p>
			<h1 class="serif">Find trips worth joining.</h1>
			<p class="subtitle">Curated journeys from the community.</p>
		</section>

		<!-- Filter Pills -->
		<FilterChips {filters} bind:active />

		<!-- Trending Bar -->
		{#if !searchQuery.trim()}
			<div class="trending">
				<div class="trend-left">
					<span class="material-symbols-outlined trend-icon">trending_up</span>
					<span class="trend-text"><span class="trend-label">Trending Now:</span> <strong>#Iceland2026</strong></span>
				</div>
				<span class="trend-count">(1.2k)</span>
			</div>
		{/if}

		<!-- Cards -->
		<div class="cards">
			{#if loading}
				{#each Array(2) as _}
					<SkeletonCard ratio="4 / 5" />
				{/each}
			{:else if error}
				<div class="state error">{error}</div>
			{:else if items.length === 0}
				<div class="state empty">No trips yet—check back soon!</div>
			{:else}
				{#each items as trip, i (trip.id)}
					<div in:fly={{ y: 40, duration: 300, delay: i * 80 }}>
						<TripCard {trip} />
					</div>
				{/each}
			{/if}
			<div bind:this={sentinel} class="sentinel" aria-hidden="true"></div>

			{#if loadingMore}
				<SkeletonCard ratio="4 / 5" />
			{/if}
		</div>
	</main>
</div>

<style>
.page {
	min-height: 100dvh;
	background: #161c18;
	color: #f1f5f9;
	padding-bottom: 96px;
}

/* ── Sticky Header ─────────────────────────────────────── */
.top-bar {
	position: sticky;
	top: 0;
	z-index: 50;
	display: flex;
	align-items: center;
	justify-content: space-between;
	height: 56px;
	padding: 0 16px;
	background: #161c18;
	border-bottom: 1px solid rgba(77, 157, 109, 0.1);
}

.logo {
	font-size: 18px;
	font-weight: 700;
	color: white;
	letter-spacing: -0.01em;
}

.icon-btn {
	width: 48px;
	height: 48px;
	display: flex;
	align-items: center;
	justify-content: center;
	color: white;
	border-radius: 8px;
	background: none;
	border: none;
	cursor: pointer;
}

.icon-btn .material-symbols-outlined {
	font-size: 24px;
}

/* ── Content ────────────────────────────────────────────── */
.content {
	padding: 0 16px;
	max-width: 480px;
	margin: 0 auto;
}

/* ── Hero ────────────────────────────────────────────────── */
.hero {
	padding-top: 32px;
	padding-bottom: 0;
	margin-bottom: 0;
}

.eyebrow {
	font-size: 11px;
	color: #4d9d6d;
	font-weight: 700;
	text-transform: uppercase;
	letter-spacing: 0.15em;
	margin-bottom: 8px;
}

h1 {
	font-size: 36px;
	line-height: 1.1;
	color: white;
	margin-bottom: 8px;
}

.subtitle {
	font-size: 15px;
	color: #94a3b8;
	font-weight: 400;
}

/* ── Trending bar ────────────────────────────────────────── */
.trending {
	display: flex;
	align-items: center;
	justify-content: space-between;
	background: rgba(77, 157, 109, 0.08);
	border-radius: 12px;
	padding: 14px 16px;
	margin-bottom: 32px;
}

.trend-left {
	display: flex;
	align-items: center;
	gap: 8px;
}

.trend-icon {
	font-size: 16px;
	color: #4d9d6d;
}

.trend-text {
	font-size: 14px;
	color: #4d9d6d;
	font-weight: 500;
}

.trend-text strong {
	font-weight: 700;
}

.trend-label {
	font-weight: 500;
}

.trend-count {
	font-size: 12px;
	color: #4d9d6d;
	font-weight: 700;
	opacity: 0.8;
}

/* ── Cards ───────────────────────────────────────────────── */
.cards {
	display: flex;
	flex-direction: column;
	gap: 32px;
}

.state {
	text-align: center;
	padding: 24px;
	border-radius: 12px;
	background: rgba(255, 255, 255, 0.04);
	color: #94a3b8;
	font-size: 14px;
}

.state.error {
	color: #f87171;
	background: rgba(248, 113, 113, 0.08);
}

.sentinel {
	height: 1px;
}
</style>
