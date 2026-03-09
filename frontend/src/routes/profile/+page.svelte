<script lang="ts">
	import { onMount } from 'svelte';
	import { apiFetch } from '$lib/api';
	import {
		formatDateRange,
		getTripLocation,
		getUserInitials,
		getUserName,
		normalizeStats
	} from '$lib/format';
	import type { TripCardData, User, UserStats } from '$lib/types';

	let user: User | null = null;
	let stats: UserStats = normalizeStats(null);
	let world = 3.5;
	let history: TripCardData[] = [];
	let error = '';

	onMount(async () => {
		try {
			const data = await apiFetch<{
				user: User;
				stats: Record<string, unknown>;
				world_explored_percent: number;
			}>('/v1/me');
			user = data.user;
			stats = normalizeStats(data.stats);
			world = data.world_explored_percent;
			const trips = await apiFetch<{ items: TripCardData[] }>('/v1/trips?scope=mine');
			history = (trips.items ?? []).slice(0, 6);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load profile';
		}
	});
</script>

<section class="container">
	<header class="header">
		<div>
			<p class="eyebrow">Profile</p>
			<h1 class="headline">Your travel identity.</h1>
		</div>
		<button class="settings" aria-label="Settings">⚙</button>
	</header>

	{#if error}
		<div class="glass card">{error}</div>
	{:else}
		<div class="profile-card panel">
			<div class="avatar-ring avatar-shell">
				{#if user?.photo_url}
					<img src={user.photo_url} alt={getUserName(user)} />
				{:else}
					<div class="avatar-fallback">{getUserInitials(user?.first_name, user?.last_name, user?.username)}</div>
				{/if}
			</div>
			<h2>{getUserName(user ?? {})}</h2>
			<p class="muted">@{user?.username ?? 'traveler'}</p>
			<button class="edit ghost-button">Edit profile</button>
		</div>

		<div class="world panel">
			<div class="world-head">
				<div>
					<h3>World Explored</h3>
					<p class="muted">Based on completed trips in your real account data.</p>
				</div>
				<span>{world.toFixed(1)}%</span>
			</div>
			<div class="map">
				<div class="glow a"></div>
				<div class="glow b"></div>
				<div class="grid-line"></div>
			</div>
		</div>

		<div class="stats">
			<div class="stat glass">
				<span>{stats.total_trips}</span>
				<small>Trips</small>
			</div>
			<div class="stat glass">
				<span>{stats.countries_visited}</span>
				<small>Countries</small>
			</div>
			<div class="stat glass">
				<span>{stats.cities_visited}</span>
				<small>Cities</small>
			</div>
			<div class="stat glass">
				<span>{stats.trips_with_friends}</span>
				<small>Friends</small>
			</div>
		</div>

		<div class="history-head">
			<h3>Travel History</h3>
			<span class="muted">{history.length} recent</span>
		</div>

		<div class="history">
			{#if history.length === 0}
				<div class="card glass">No trips to show yet</div>
			{:else}
				{#each history as trip}
					<a class="history-item glass" href={`/trips/${trip.id}`}>
						{#if trip.cover_photo_url}
							<img src={trip.cover_photo_url} alt={trip.title} />
						{:else}
							<div class="thumb-placeholder"></div>
						{/if}
						<div class="history-copy">
							<strong>{trip.title}</strong>
							<p>{getTripLocation(trip)}</p>
							<small>{formatDateRange(trip.start_date, trip.end_date)}</small>
						</div>
						<span class="chevron">›</span>
					</a>
				{/each}
			{/if}
		</div>
	{/if}
</section>

<style>
	.header {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		gap: 1rem;
		margin-bottom: 1rem;
	}

	.settings {
		width: 2.8rem;
		height: 2.8rem;
		border-radius: 50%;
		background: rgba(255, 255, 255, 0.06);
		border: 1px solid rgba(255, 255, 255, 0.08);
	}

	.card {
		padding: 1.5rem;
		border-radius: var(--radius-xl);
		margin-bottom: 1rem;
		text-align: center;
	}

	.profile-card {
		padding: 1.6rem;
		border-radius: var(--radius-2xl);
		margin-bottom: 1rem;
		text-align: center;
	}

	.avatar-shell {
		width: 96px;
		height: 96px;
		margin: 0 auto 0.9rem;
	}

	h2 {
		font-size: 1.35rem;
		font-weight: 800;
		margin-bottom: 0.25rem;
	}

	.edit {
		margin-top: 1rem;
		width: 100%;
	}

	.world {
		padding: 1.35rem;
		border-radius: var(--radius-2xl);
		margin-bottom: 1rem;
	}

	.world-head {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		gap: 1rem;
		margin-bottom: 1rem;
	}

	.world span {
		font-size: 2rem;
		font-weight: 800;
		color: var(--accent-strong);
	}

	.map {
		position: relative;
		height: 144px;
		border-radius: 20px;
		background:
			linear-gradient(180deg, rgba(255, 255, 255, 0.04), transparent),
			radial-gradient(circle at 30% 40%, rgba(122, 234, 244, 0.25), transparent 12%),
			radial-gradient(circle at 58% 55%, rgba(32, 146, 186, 0.25), transparent 13%),
			radial-gradient(circle at 75% 38%, rgba(122, 234, 244, 0.22), transparent 10%),
			#071a31;
		overflow: hidden;
	}

	.glow {
		position: absolute;
		width: 14px;
		height: 14px;
		border-radius: 50%;
		background: var(--accent-strong);
		box-shadow: 0 0 0 8px rgba(122, 234, 244, 0.12);
		animation: pulse 2.2s infinite;
	}

	.glow.a {
		top: 34%;
		left: 28%;
	}

	.glow.b {
		top: 52%;
		left: 64%;
		animation-delay: 0.8s;
	}

	.grid-line {
		position: absolute;
		inset: 0;
		background-image: linear-gradient(rgba(255, 255, 255, 0.05) 1px, transparent 1px),
			linear-gradient(90deg, rgba(255, 255, 255, 0.05) 1px, transparent 1px);
		background-size: 32px 32px;
		opacity: 0.35;
	}

	.stats {
		display: grid;
		grid-template-columns: repeat(4, 1fr);
		gap: 0.75rem;
		margin-bottom: 1.1rem;
	}

	.stat {
		padding: 1rem 0.8rem;
		border-radius: var(--radius-lg);
		text-align: center;
	}

	.stat span {
		display: block;
		font-size: 1.45rem;
		font-weight: 800;
		color: var(--accent-strong);
		margin-bottom: 0.25rem;
	}

	small {
		color: var(--text-secondary);
	}

	.history-head {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 0.8rem;
	}

	.history {
		display: grid;
		gap: 0.8rem;
	}

	.history-item {
		display: grid;
		grid-template-columns: 68px minmax(0, 1fr) auto;
		align-items: center;
		gap: 0.85rem;
		padding: 0.75rem;
		border-radius: var(--radius-xl);
	}

	.history-item img,
	.thumb-placeholder {
		width: 68px;
		height: 68px;
		border-radius: 16px;
		object-fit: cover;
	}

	.thumb-placeholder {
		background: linear-gradient(135deg, rgba(32, 146, 186, 0.45), rgba(13, 52, 96, 0.9));
	}

	.history-copy {
		min-width: 0;
	}

	.history-copy strong,
	.history-copy p,
	.history-copy small {
		display: block;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.history-copy p {
		color: var(--text-secondary);
		margin: 0.2rem 0;
	}

	.chevron {
		font-size: 1.5rem;
		color: var(--text-muted);
	}

	@keyframes pulse {
		0%, 100% {
			transform: scale(1);
			opacity: 0.8;
		}
		50% {
			transform: scale(1.25);
			opacity: 1;
		}
	}

	@media (max-width: 640px) {
		.stats {
			grid-template-columns: repeat(2, 1fr);
		}
	}
</style>
