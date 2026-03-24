<script lang="ts">
	import { onMount } from 'svelte';
	import { fly } from 'svelte/transition';
	import TripCard from '$lib/components/TripCard.svelte';
	import SkeletonCard from '$lib/components/SkeletonCard.svelte';
	import { apiFetch } from '$lib/api';
	import { expandTelegram } from '$lib/telegram';
	import type { TripCardData } from '$lib/types';

	// ── Public trips ────────────────────────────────────────
	let publicItems = $state<TripCardData[]>([]);
	let publicLoading = $state(true);
	let publicError = $state('');
	let publicCursor = $state<string | null>(null);
	let publicLoadingMore = $state(false);
	let publicHasMore = $state(true);
	let publicSentinel = $state<HTMLDivElement | null>(null);

	// ── Trending trips ──────────────────────────────────────
	let trendingItems = $state<TripCardData[]>([]);
	let trendingLoading = $state(true);
	let trendingError = $state('');

	// ── Friends' trips ──────────────────────────────────────
	let friendItems = $state<TripCardData[]>([]);
	let friendLoading = $state(true);
	let friendError = $state('');

	let ready = $state(false);

	const fetchPublicTrips = async (mode: 'reset' | 'append' = 'reset') => {
		if (mode === 'append' && (publicLoadingMore || !publicHasMore)) return;
		if (mode === 'append') publicLoadingMore = true;
		else publicLoading = true;
		publicError = '';
		const queryArgs = ['filter=all', 'limit=20'];
		if (mode === 'append' && publicCursor) {
			queryArgs.push(`cursor=${publicCursor}`);
		}
		try {
			const data = await apiFetch<{ items: TripCardData[]; next_cursor?: string; cursor?: string }>(
				`/v1/feed?${queryArgs.join('&')}`
			);
			const next = data.next_cursor ?? data.cursor ?? null;
			if (mode === 'append') {
				publicItems = [...publicItems, ...(data.items ?? [])];
			} else {
				publicItems = data.items ?? [];
			}
			publicCursor = next;
			publicHasMore = Boolean(next);
		} catch (err) {
			publicError = err instanceof Error ? err.message : 'Failed to load trips';
		} finally {
			publicLoading = false;
			publicLoadingMore = false;
		}
	};

	const fetchTrending = async () => {
		trendingLoading = true;
		trendingError = '';
		try {
			const data = await apiFetch<{ items: TripCardData[] }>(
				'/v1/feed?filter=popular&limit=6'
			);
			trendingItems = data.items ?? [];
		} catch (err) {
			trendingError = err instanceof Error ? err.message : 'Failed to load';
		} finally {
			trendingLoading = false;
		}
	};

	const fetchFriendTrips = async () => {
		friendLoading = true;
		friendError = '';
		try {
			const data = await apiFetch<{ items: TripCardData[] }>(
				'/v1/feed?filter=friends&limit=10'
			);
			friendItems = data.items ?? [];
		} catch (err) {
			const message = err instanceof Error ? err.message : 'Failed to load';
			if (message.toLowerCase().includes('unauthorized')) {
				friendError = 'connect';
			} else {
				friendError = message;
			}
		} finally {
			friendLoading = false;
		}
	};

	const observePublic = () => {
		if (!publicSentinel) return;
		const io = new IntersectionObserver((entries) => {
			if (entries[0]?.isIntersecting) fetchPublicTrips('append');
		});
		io.observe(publicSentinel);
		return () => io.disconnect();
	};

	onMount(() => {
		expandTelegram();
		ready = true;
		fetchTrending();
		fetchPublicTrips('reset');
		fetchFriendTrips();
		return observePublic();
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

		<!-- Trending -->
		<section class="feed-section">
			<p class="section-label">Trending</p>
			<div class="trending-row">
				{#if trendingLoading}
					{#each Array(3) as _}
						<div class="trending-skeleton"></div>
					{/each}
				{:else if trendingItems.length > 0}
					{#each trendingItems as trip (trip.id)}
						<div class="trending-card-wrap">
							<TripCard {trip} variant="compact" />
						</div>
					{/each}
				{/if}
			</div>
		</section>

		<!-- For You -->
		<section class="feed-section">
			<p class="section-label">For You</p>
			<div class="cards">
				{#if publicLoading}
					{#each Array(2) as _}
						<SkeletonCard ratio="3 / 4" />
					{/each}
				{:else if publicError}
					<div class="state error">{publicError}</div>
				{:else if publicItems.length === 0}
					<div class="state empty">No public trips yet — check back soon!</div>
				{:else}
					{#each publicItems as trip, i (trip.id)}
						<div in:fly={{ y: 40, duration: 300, delay: i * 80 }}>
							<TripCard {trip} />
						</div>
					{/each}
				{/if}
				<div bind:this={publicSentinel} class="sentinel" aria-hidden="true"></div>
				{#if publicLoadingMore}
					<SkeletonCard ratio="3 / 4" />
				{/if}
			</div>
		</section>

		<!-- Friends' Trips -->
		<section class="feed-section">
			<p class="section-label">Friends' Trips</p>
			<div class="cards">
				{#if friendLoading}
					{#each Array(1) as _}
						<SkeletonCard ratio="3 / 4" />
					{/each}
				{:else if friendError === 'connect' || friendItems.length === 0}
					<div class="state connect">
						<span class="material-symbols-outlined connect-icon">group</span>
						<p>Log in to see friends' trips</p>
					</div>
				{:else if friendError}
					<div class="state error">{friendError}</div>
				{:else}
					{#each friendItems as trip, i (trip.id)}
						<div in:fly={{ y: 40, duration: 300, delay: i * 80 }}>
							<TripCard {trip} />
						</div>
					{/each}
				{/if}
			</div>
		</section>
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

/* ── Feed sections ────────────────────────────────────────── */
.feed-section {
	margin-top: 32px;
}

.section-label {
	font-size: 11px;
	font-weight: 700;
	color: #4d9d6d;
	text-transform: uppercase;
	letter-spacing: 0.12em;
	margin-bottom: 16px;
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

.state.connect {
	display: flex;
	flex-direction: column;
	align-items: center;
	gap: 8px;
	padding: 32px 24px;
}

.connect-icon {
	font-size: 32px;
	color: #4d9d6d;
	opacity: 0.6;
}

.sentinel {
	height: 1px;
}

/* ── Trending row ────────────────────────────────────────── */
.trending-row {
	display: flex;
	gap: 12px;
	overflow-x: auto;
	padding-bottom: 8px;
	scrollbar-width: none;
}

.trending-row::-webkit-scrollbar {
	display: none;
}

.trending-card-wrap {
	flex: 0 0 160px;
}

.trending-skeleton {
	flex: 0 0 160px;
	aspect-ratio: 3 / 4;
	border-radius: 12px;
	background: rgba(255, 255, 255, 0.06);
	animation: shimmer 1.5s infinite;
}

@keyframes shimmer {
	0%, 100% { opacity: 0.6; }
	50% { opacity: 1; }
}
</style>
