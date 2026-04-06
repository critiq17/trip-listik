<script lang="ts">
	import { onMount } from 'svelte';
	import { fly } from 'svelte/transition';
	import TripCard from '$lib/components/TripCard.svelte';
	import SkeletonCard from '$lib/components/SkeletonCard.svelte';
	import { apiFetch } from '$lib/api';
	import { expandTelegram } from '$lib/telegram';
	import type { TripCardData } from '$lib/types';

	type Tab = 'all' | 'mine';

	let activeTab = $state<Tab>('all');

	// ── All trips (public feed) ──────────────────────────────
	let allItems = $state<TripCardData[]>([]);
	let allLoading = $state(true);
	let allError = $state('');
	let allCursor = $state<string | null>(null);
	let allLoadingMore = $state(false);
	let allHasMore = $state(true);
	let allSentinel = $state<HTMLDivElement | null>(null);

	// ── My trips ─────────────────────────────────────────────
	let myItems = $state<TripCardData[]>([]);
	let myLoading = $state(true);
	let myError = $state('');

	let ready = $state(false);

	const fetchAllTrips = async (mode: 'reset' | 'append' = 'reset') => {
		if (mode === 'append' && (allLoading || allLoadingMore || !allHasMore || !allCursor)) return;
		if (mode === 'append') allLoadingMore = true;
		else {
			allLoading = true;
			allCursor = null;
			allHasMore = true;
		}
		allError = '';
		const params = ['filter=all', 'limit=20'];
		if (mode === 'append' && allCursor) params.push(`cursor=${allCursor}`);
		try {
			const data = await apiFetch<{ items: TripCardData[]; next_cursor?: string }>(`/v1/feed?${params.join('&')}`);
			const next = data.next_cursor ?? null;
			allItems = mode === 'append' ? [...allItems, ...(data.items ?? [])] : (data.items ?? []);
			allCursor = next;
			allHasMore = Boolean(next);
		} catch (err) {
			allError = err instanceof Error ? err.message : 'Failed to load trips';
		} finally {
			allLoading = false;
			allLoadingMore = false;
		}
	};

	const fetchMyTrips = async () => {
		myLoading = true;
		myError = '';
		try {
			const data = await apiFetch<{ items: TripCardData[] }>('/v1/trips?scope=mine');
			myItems = data.items ?? [];
		} catch (err) {
			const msg = err instanceof Error ? err.message : 'Failed to load';
			myError = msg.toLowerCase().includes('unauthorized') ? 'auth' : msg;
		} finally {
			myLoading = false;
		}
	};

	const observeAll = () => {
		if (!allSentinel) return;
		const io = new IntersectionObserver((entries) => {
			if (entries[0]?.isIntersecting) fetchAllTrips('append');
		});
		io.observe(allSentinel);
		return () => io.disconnect();
	};

	onMount(() => {
		expandTelegram();
		ready = true;
		fetchAllTrips('reset');
		fetchMyTrips();
		return observeAll();
	});
</script>

<svelte:head>
	<title>TripListik — Feed</title>
	<meta name="description" content="Browse curated trips and travel journeys from the community." />
</svelte:head>

<div class="page">
	<header class="top-bar">
		<span class="logo">TripListik</span>
		<a href="/create" class="create-btn" aria-label="Create trip">
			<span class="material-symbols-outlined">add</span>
		</a>
	</header>

	<div class="tabs-wrap">
		<div class="tabs">
			<button class="tab" class:active={activeTab === 'all'} onclick={() => (activeTab = 'all')}>
				All Trips
			</button>
			<button class="tab" class:active={activeTab === 'mine'} onclick={() => (activeTab = 'mine')}>
				My Trips
			</button>
		</div>
	</div>

	<main class="content">
		{#if activeTab === 'all'}
			<div class="cards">
				{#if allLoading}
					{#each Array(3) as _}
						<SkeletonCard ratio="16 / 9" />
					{/each}
				{:else if allError}
					<div class="state error">{allError}</div>
				{:else if allItems.length === 0}
					<div class="state empty">No public trips yet — check back soon!</div>
				{:else}
					{#each allItems as trip, i (trip.id)}
						<div in:fly={{ y: 30, duration: 250, delay: i < 6 ? i * 60 : 0 }}>
							<TripCard {trip} variant="compact" />
						</div>
					{/each}
				{/if}
				<div bind:this={allSentinel} class="sentinel" aria-hidden="true"></div>
				{#if allLoadingMore}
					<SkeletonCard ratio="16 / 9" />
				{/if}
			</div>
		{:else}
			<div class="cards">
				{#if myLoading}
					{#each Array(3) as _}
						<SkeletonCard ratio="16 / 9" />
					{/each}
				{:else if myError === 'auth'}
					<div class="state connect">
						<span class="material-symbols-outlined connect-icon">lock</span>
						<p>Log in to see your trips</p>
					</div>
				{:else if myError}
					<div class="state error">{myError}</div>
				{:else if myItems.length === 0}
					<div class="state empty">
						<p>No trips yet.</p>
						<a href="/create" class="empty-cta">Plan your first trip</a>
					</div>
				{:else}
					{#each myItems as trip, i (trip.id)}
						<div in:fly={{ y: 30, duration: 250, delay: i < 6 ? i * 60 : 0 }}>
							<TripCard {trip} variant="compact" />
						</div>
					{/each}
				{/if}
			</div>
		{/if}
	</main>
</div>

<style>
.page {
	min-height: 100dvh;
	background: #161c18;
	color: #f1f5f9;
	padding-bottom: 96px;
}

/* ── Sticky Header ───────────────────────────────────────── */
.top-bar {
	position: sticky;
	top: 0;
	z-index: 50;
	display: flex;
	align-items: center;
	justify-content: space-between;
	height: 52px;
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

.create-btn {
	width: 36px;
	height: 36px;
	display: flex;
	align-items: center;
	justify-content: center;
	background: rgba(77, 157, 109, 0.15);
	border: 1px solid rgba(77, 157, 109, 0.3);
	border-radius: 8px;
	color: #4d9d6d;
	text-decoration: none;
}

.create-btn .material-symbols-outlined {
	font-size: 22px;
}

/* ── Tabs ────────────────────────────────────────────────── */
.tabs-wrap {
	position: sticky;
	top: 52px;
	z-index: 40;
	background: #161c18;
	border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.tabs {
	display: flex;
	padding: 0 16px;
	max-width: 480px;
	margin: 0 auto;
	gap: 24px;
}

.tab {
	padding: 12px 0;
	font-size: 14px;
	font-weight: 600;
	color: #64748b;
	background: none;
	border: none;
	border-bottom: 2px solid transparent;
	cursor: pointer;
	transition: color 0.15s ease, border-color 0.15s ease;
	-webkit-tap-highlight-color: transparent;
}

.tab.active {
	color: #4d9d6d;
	border-bottom-color: #4d9d6d;
}

/* ── Content ─────────────────────────────────────────────── */
.content {
	padding: 12px 16px 0;
	max-width: 480px;
	margin: 0 auto;
}

/* ── Cards ───────────────────────────────────────────────── */
.cards {
	display: flex;
	flex-direction: column;
	gap: 12px;
}

.state {
	text-align: center;
	padding: 32px 24px;
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
}

.connect-icon {
	font-size: 32px;
	color: #4d9d6d;
	opacity: 0.6;
}

.empty-cta {
	display: inline-block;
	margin-top: 12px;
	padding: 8px 20px;
	background: rgba(77, 157, 109, 0.15);
	border: 1px solid rgba(77, 157, 109, 0.3);
	border-radius: 8px;
	color: #4d9d6d;
	font-size: 14px;
	font-weight: 600;
	text-decoration: none;
}

.sentinel {
	height: 1px;
}
</style>
