<script lang="ts">
	import { onMount } from 'svelte';
	import SkeletonCard from '$lib/components/SkeletonCard.svelte';
	import { apiFetch } from '$lib/api';
	import { formatDateRange, getStatusLabel } from '$lib/format';
	import { staggerList, scalePress } from '$lib/actions/animate';
	import { animate } from 'motion';
	import type { TripCardData } from '$lib/types';

	const tabs = ['Upcoming', 'Past', 'Drafts'];
	let active = 'Upcoming';
	let items: TripCardData[] = [];
	let loading = true;
	let error = '';
	let ready = false;
	let tabBar: HTMLDivElement | null = null;
	let indicator: HTMLDivElement | null = null;
	let tabRefs: HTMLButtonElement[] = [];

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
			{ x: [`${indicator.offsetLeft}px`, `${left}px`], width: [`${indicator.offsetWidth}px`, `${tabRect.width}px`] },
			{ duration: 0.25, easing: [0.25, 0.46, 0.45, 0.94] }
		);
		indicator.style.transform = `translateX(${left}px)`;
		indicator.style.width = `${tabRect.width}px`;
	};

	onMount(() => {
		ready = true;
		setTimeout(moveIndicator, 0);
	});

	$: if (ready && active) {
		loadTrips();
		setTimeout(moveIndicator, 0);
	}
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
			<button bind:this={tabRefs[i]} class:active={tab === active} on:click={() => (active = tab)}>
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
					<a class="trip-card" data-item use:scalePress href={`/trips/${trip.id}`}>
						<div class="trip-media">
							{#if trip.cover_photo_url}
								<img src={trip.cover_photo_url} alt={trip.title} loading="lazy" />
							{:else}
								<div class="trip-placeholder"></div>
							{/if}
							<div class="trip-overlay"></div>
							<div class="trip-status">{getStatusLabel(trip.status)}</div>
							<div class="trip-content">
								<div>
									<h3>{trip.title}</h3>
									<p>
										<span class="material-symbols-outlined">calendar_today</span>
										{formatDateRange(trip.start_date, trip.end_date)}
									</p>
								</div>
								<div class="avatars">
									<div class="avatar"></div>
									<div class="avatar count">+{Math.max((trip.member_count ?? 1) - 1, 1)}</div>
								</div>
							</div>
						</div>
					</a>
				{/each}
			{/if}
		</div>

		{#if !loading && !error}
			<div class="drafts">
				<div class="drafts-head">
					<h4>Quick Drafts</h4>
				</div>
				<div class="drafts-box">
					<div class="drafts-icon">
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
			</div>
		{/if}
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

	.trip-card {
		border-radius: 12px;
		overflow: hidden;
		box-shadow: var(--shadow-card);
	}

	.trip-media {
		position: relative;
		aspect-ratio: 16 / 10;
	}

	.trip-media img,
	.trip-placeholder {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.trip-placeholder {
		background: linear-gradient(135deg, #0f1411, #1a251f);
	}

	.trip-overlay {
		position: absolute;
		inset: 0;
		background: linear-gradient(to top, rgba(0, 0, 0, 0.75), rgba(0, 0, 0, 0.2), transparent);
	}

	.trip-status {
		position: absolute;
		top: 1rem;
		right: 1rem;
		padding: 0.3rem 0.6rem;
		border-radius: 0.4rem;
		background: rgba(255, 255, 255, 0.2);
		color: white;
		font-size: 0.6rem;
		text-transform: uppercase;
		letter-spacing: 0.1em;
		font-weight: 700;
		backdrop-filter: blur(10px);
	}

	.trip-content {
		position: absolute;
		left: 0;
		right: 0;
		bottom: 0;
		padding: 1.1rem;
		display: flex;
		justify-content: space-between;
		align-items: flex-end;
		color: white;
	}

	.trip-content h3 {
		font-size: 1.2rem;
		font-weight: 700;
		margin-bottom: 0.25rem;
	}

	.trip-content p {
		display: flex;
		align-items: center;
		gap: 0.3rem;
		font-size: 0.78rem;
		color: rgba(255, 255, 255, 0.8);
	}

	.trip-content .material-symbols-outlined {
		font-size: 0.95rem;
	}

	.avatars {
		display: flex;
		align-items: center;
	}

	.avatar {
		width: 2rem;
		height: 2rem;
		border-radius: 999px;
		border: 2px solid white;
		background: #6b7280;
		margin-left: -0.5rem;
		display: grid;
		place-items: center;
		color: white;
		font-size: 0.6rem;
		font-weight: 700;
	}

	.avatar:first-child {
		margin-left: 0;
	}

	.avatar.count {
		background: var(--primary);
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

	.drafts {
		margin-top: 2rem;
	}

	.drafts-head h4 {
		font-size: 0.65rem;
		letter-spacing: 0.2em;
		text-transform: uppercase;
		color: var(--text-secondary);
		margin-bottom: 1rem;
	}

	.drafts-box {
		border: 2px dashed rgba(255, 255, 255, 0.12);
		border-radius: 12px;
		padding: 2rem 1.5rem;
		text-align: center;
		display: flex;
		flex-direction: column;
		gap: 0.8rem;
	}

	.drafts-icon {
		width: 3rem;
		height: 3rem;
		border-radius: 999px;
		background: rgba(255, 255, 255, 0.08);
		display: grid;
		place-items: center;
		margin: 0 auto;
		color: rgba(255, 255, 255, 0.6);
	}

	.drafts-box p {
		font-weight: 600;
		color: var(--text-primary);
	}

	.drafts-box span {
		font-size: 0.75rem;
		color: var(--text-secondary);
	}
</style>
