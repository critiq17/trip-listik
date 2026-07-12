<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import SkeletonCard from '$lib/components/SkeletonCard.svelte';
	import TripCard from '$lib/components/TripCard.svelte';
	import { apiFetch } from '$lib/api';
	import type { TripCardData } from '$lib/types';

	const tabs = ['Upcoming', 'Past', 'Drafts'] as const;
	type Tab = typeof tabs[number];

	let active = $state<Tab>('Upcoming');
	let items = $state<TripCardData[]>([]);
	let loading = $state(true);
	let error = $state('');
	let ready = $state(false);

	async function loadTrips() {
		loading = true;
		error = '';
		const status = active === 'Drafts' ? 'draft' : active.toLowerCase();
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

	$effect(() => {
		if (ready && active) {
			loadTrips();
		}
	});
</script>

<svelte:head>
	<title>TripListik — My Trips</title>
	<meta name="description" content="Manage your upcoming, past, and draft travel trips." />
</svelte:head>

<div class="page">
	<header class="top-bar">
		<span class="page-title">My Trips</span>
		<a href="/create" class="create-btn" aria-label="Create trip">
			<span class="material-symbols-outlined">add</span>
		</a>
	</header>

	<!-- Tab Bar -->
	<nav class="tab-bar">
		{#each tabs as tab}
			<button
				class="tab"
				class:active={tab === active}
				onclick={() => (active = tab)}
			>
				{tab}
			</button>
		{/each}
	</nav>

	<!-- Content -->
	<main class="content">
		<div class="cards">
			{#if loading}
				<SkeletonCard ratio="16 / 10" />
			{:else if error}
				<div class="state">{error}</div>
			{:else if items.length === 0 && active === 'Upcoming'}
				<div class="empty-state">
					<span class="material-symbols-outlined empty-icon">luggage</span>
					<p class="empty-title">No upcoming trips</p>
					<span class="empty-sub">Plan your next adventure and it'll appear here.</span>
					<button class="empty-action" onclick={() => goto('/create')}>
						<span class="material-symbols-outlined">add</span>
						Create new trip
					</button>
				</div>
			{:else if items.length === 0}
				<div class="empty-minimal">No {active.toLowerCase()} trips found.</div>
			{:else}
				{#each items as trip (trip.id)}
					<TripCard {trip} variant="compact" />
				{/each}
			{/if}
		</div>
	</main>
</div>

<style>
.page {
	min-height: 100dvh;
	background: var(--bg);
	color: var(--text);
	padding-bottom: 96px;
}

/* ── Header ─────────────────────────────────────────────── */
.top-bar {
	position: sticky;
	top: 0;
	z-index: 50;
	display: flex;
	align-items: center;
	justify-content: space-between;
	height: 52px;
	padding: 0 16px;
	background: var(--bg);
	border-bottom: 1px solid var(--border);
	max-width: 480px;
	margin: 0 auto;
}

.page-title {
	font-size: 18px;
	font-weight: 700;
	color: var(--text);
}

.create-btn {
	width: 36px;
	height: 36px;
	display: flex;
	align-items: center;
	justify-content: center;
	background: var(--green-soft);
	border-radius: 8px;
	color: var(--green);
	text-decoration: none;
}

.create-btn .material-symbols-outlined {
	font-size: 22px;
}

/* ── Tab Bar ─────────────────────────────────────────────── */
.tab-bar {
	display: flex;
	gap: 24px;
	padding: 0 16px;
	max-width: 480px;
	margin: 0 auto;
	border-bottom: 1px solid var(--border);
}

.tab {
	padding: 12px 0;
	background: none;
	border: none;
	border-bottom: 2px solid transparent;
	color: var(--text-sub);
	font-size: 14px;
	font-weight: 600;
	cursor: pointer;
	transition: color 0.15s ease, border-color 0.15s ease;
	-webkit-tap-highlight-color: transparent;
}

.tab.active {
	color: var(--green);
	border-bottom-color: var(--green);
}

/* ── Content ─────────────────────────────────────────────── */
.content {
	padding: 12px 16px 0;
	max-width: 480px;
	margin: 0 auto;
}

.cards {
	display: flex;
	flex-direction: column;
	gap: 12px;
}

/* ── Empty States ────────────────────────────────────────── */
.empty-state {
	display: flex;
	flex-direction: column;
	align-items: center;
	text-align: center;
	gap: 8px;
	padding: 40px 16px;
}

.empty-icon {
	font-size: 32px;
	color: var(--text-muted);
	margin-bottom: 4px;
}

.empty-title {
	font-size: 16px;
	font-weight: 600;
	color: var(--text);
}

.empty-sub {
	font-size: 13px;
	color: var(--text-sub);
}

.empty-action {
	display: flex;
	align-items: center;
	gap: 6px;
	margin-top: 8px;
	padding: 8px 20px;
	background: var(--green-soft);
	border-radius: 8px;
	color: var(--green);
	font-size: 14px;
	font-weight: 600;
	border: none;
	cursor: pointer;
}

.empty-action .material-symbols-outlined {
	font-size: 18px;
}

.empty-minimal {
	text-align: center;
	padding: 32px;
	color: var(--text-sub);
	font-size: 14px;
}

.state {
	padding: 20px;
	text-align: center;
	border-radius: var(--radius-card);
	background: var(--danger-soft);
	color: var(--danger);
	font-size: 14px;
}
</style>
