<script lang="ts">
	import { onMount } from 'svelte';
	import { fly } from 'svelte/transition';
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
	<!-- Header -->
	<header class="header">
		<div class="header-row">
			<button class="round-btn" aria-label="Back" onclick={() => history.back()}>
				<span class="material-symbols-outlined">arrow_back</span>
			</button>
			<button class="round-btn" aria-label="More options">
				<span class="material-symbols-outlined">more_horiz</span>
			</button>
		</div>
		<div class="header-copy">
			<p class="eyebrow">My Trips</p>
			<h1>Your travel board.</h1>
			<p class="subtitle">Manage your upcoming and past adventures.</p>
		</div>
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
		{#if !loading && !error}
			<div class="stats-badge">
				{items.length} visible · {active}
			</div>
		{/if}

		<div class="cards">
			{#if loading}
				<SkeletonCard ratio="16 / 10" />
			{:else if error}
				<div class="state">{error}</div>
			{:else if items.length === 0 && active === 'Upcoming'}
				<!-- Empty upcoming: show empty state -->
				<div class="empty-state">
					<div class="empty-icon">
						<span class="material-symbols-outlined">luggage</span>
					</div>
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
				{#each items as trip, i (trip.id)}
					<div in:fly={{ y: 24, duration: 250, delay: i * 60 }}>
						<TripCard {trip} variant="compact" />
					</div>
				{/each}
			{/if}
		</div>

		<!-- Quick Drafts Section -->
		{#if active === 'Upcoming' && !loading}
			<div class="drafts-section">
				<p class="section-label">Quick Drafts</p>
				<div class="draft-empty">
					<div class="draft-icon-wrap">
						<span class="material-symbols-outlined">edit_note</span>
					</div>
					<p class="draft-title">No active drafts</p>
					<span class="draft-sub">Start planning your next escape today.</span>
					<button class="draft-action" onclick={() => goto('/create')}>
						<span class="material-symbols-outlined">add</span>
						Create new trip
					</button>
				</div>
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

/* ── Header ─────────────────────────────────────────────── */
.header {
	padding: 24px 16px 0;
	max-width: 480px;
	margin: 0 auto;
}

.header-row {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 16px;
}

.round-btn {
	width: 40px;
	height: 40px;
	border-radius: 9999px;
	display: flex;
	align-items: center;
	justify-content: center;
	color: white;
	background: none;
	border: none;
	cursor: pointer;
}

.round-btn .material-symbols-outlined {
	font-size: 22px;
}

.header-copy {
	margin-top: 16px;
}

.eyebrow {
	font-size: 10px;
	color: #4d9d6d;
	font-weight: 700;
	text-transform: uppercase;
	letter-spacing: 0.12em;
	margin-bottom: 6px;
}

h1 {
	font-size: 30px;
	font-weight: 700;
	color: white;
	line-height: 1.15;
	margin-bottom: 6px;
}

.subtitle {
	font-size: 14px;
	color: #94a3b8;
}

/* ── Tab Bar ─────────────────────────────────────────────── */
.tab-bar {
	display: flex;
	gap: 32px;
	padding: 0 16px;
	max-width: 480px;
	margin: 16px auto 0;
	border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}

.tab {
	padding: 12px 0;
	background: none;
	border: none;
	border-bottom: 2px solid transparent;
	color: #64748b;
	font-size: 14px;
	font-weight: 600;
	cursor: pointer;
	transition: color 0.2s ease, border-color 0.2s ease;
	-webkit-tap-highlight-color: transparent;
}

.tab.active {
	color: #4d9d6d;
	border-bottom-color: #4d9d6d;
}

/* ── Content ─────────────────────────────────────────────── */
.content {
	padding: 16px 16px 0;
	max-width: 480px;
	margin: 0 auto;
}

.stats-badge {
	display: inline-block;
	padding: 4px 12px;
	border-radius: 9999px;
	background: rgba(77, 157, 109, 0.1);
	color: #4d9d6d;
	font-size: 12px;
	font-weight: 500;
	margin-bottom: 16px;
}

.cards {
	display: flex;
	flex-direction: column;
	gap: 16px;
}

/* ── Empty States ────────────────────────────────────────── */
.empty-state {
	display: flex;
	flex-direction: column;
	align-items: center;
	text-align: center;
	gap: 12px;
	padding: 32px 16px;
}

.empty-icon {
	width: 56px;
	height: 56px;
	border-radius: 9999px;
	background: rgba(77, 157, 109, 0.1);
	display: flex;
	align-items: center;
	justify-content: center;
	color: #64748b;
}

.empty-icon .material-symbols-outlined {
	font-size: 28px;
}

.empty-title {
	font-size: 16px;
	font-weight: 600;
	color: white;
}

.empty-sub {
	font-size: 13px;
	color: #94a3b8;
}

.empty-action {
	display: flex;
	align-items: center;
	gap: 6px;
	color: #4d9d6d;
	font-size: 14px;
	font-weight: 700;
	background: none;
	border: none;
	cursor: pointer;
}

.empty-action .material-symbols-outlined {
	font-size: 18px;
}

.empty-minimal {
	text-align: center;
	padding: 32px;
	color: #64748b;
	font-size: 14px;
}

.state {
	padding: 20px;
	text-align: center;
	border-radius: 12px;
	background: rgba(255, 255, 255, 0.04);
	color: #94a3b8;
	font-size: 14px;
}

/* ── Drafts section ──────────────────────────────────────── */
.drafts-section {
	margin-top: 24px;
}

.section-label {
	font-size: 12px;
	color: #64748b;
	font-weight: 700;
	text-transform: uppercase;
	letter-spacing: 0.12em;
	margin-bottom: 16px;
}

.draft-empty {
	border: 2px dashed rgba(255, 255, 255, 0.1);
	border-radius: 12px;
	padding: 32px;
	display: flex;
	flex-direction: column;
	align-items: center;
	text-align: center;
	gap: 12px;
}

.draft-icon-wrap {
	width: 48px;
	height: 48px;
	border-radius: 9999px;
	background: rgba(77, 157, 109, 0.1);
	display: flex;
	align-items: center;
	justify-content: center;
	color: #64748b;
}

.draft-icon-wrap .material-symbols-outlined {
	font-size: 24px;
}

.draft-title {
	font-size: 14px;
	font-weight: 600;
	color: white;
}

.draft-sub {
	font-size: 13px;
	color: #94a3b8;
}

.draft-action {
	display: flex;
	align-items: center;
	gap: 6px;
	color: #4d9d6d;
	font-size: 14px;
	font-weight: 700;
	background: none;
	border: none;
	cursor: pointer;
}

.draft-action .material-symbols-outlined {
	font-size: 18px;
}
</style>
