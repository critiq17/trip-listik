<script lang="ts">
	import { onMount } from 'svelte';
	import { apiFetch } from '$lib/api';
	import { countUp } from '$lib/transitions';
	import { formatDateRange, getTripLocation, getUserInitials, getUserName } from '$lib/format';
	import type { TripCardData, User, UserStats } from '$lib/types';

	let user: User | null = null;
	let stats: UserStats | null = null;
	let world = 0;
	let history: TripCardData[] = [];
	let error = '';
	let statEls: Array<HTMLElement | null> = [];
	let worldEl: HTMLElement | null = null;
	let worldBar: HTMLDivElement | null = null;

	onMount(async () => {
		try {
			const [me, statsRes, worldRes] = await Promise.all([
				apiFetch<{ user: User }>('/v1/me'),
				apiFetch<UserStats>('/v1/me/stats'),
				apiFetch<{ world_explored_percent: number }>('/v1/me/world')
			]);
			user = me.user;
			stats = statsRes;
			world = worldRes.world_explored_percent ?? 0;
			const trips = await apiFetch<{ items: TripCardData[] }>('/v1/trips?scope=mine');
			history = (trips.items ?? []).slice(0, 6);

			if (stats && statEls.length) {
				const values = [stats.total_trips, stats.countries_visited, stats.cities_visited, stats.trips_with_friends];
				statEls.forEach((el, idx) => {
					if (el) countUp(el, values[idx] ?? 0);
				});
			}
			if (worldEl) countUp(worldEl, world, 1.1);
			if (worldBar) {
				worldBar.style.width = '0%';
				requestAnimationFrame(() => {
					worldBar && (worldBar.style.width = `${world}%`);
				});
			}
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load profile';
		}
	});
</script>

<section class="profile-page">
	<header class="header">
		<span class="eyebrow">Profile</span>
		<button class="settings" aria-label="Settings">
			<span class="material-symbols-outlined">settings</span>
		</button>
	</header>

	<main class="content hide-scrollbar">
		<h1>Your travel identity.</h1>

		{#if error}
			<div class="state">{error}</div>
		{:else}
			<div class="profile-row">
				<div class="avatar-wrap">
					{#if user?.photo_url}
						<div class="avatar" style={`background-image: url('${user.photo_url}')`}></div>
					{:else}
						<div class="avatar fallback">{getUserInitials(user?.first_name, user?.last_name, user?.username)}</div>
					{/if}
				</div>
				<div>
					<h2>{getUserName(user ?? {})}</h2>
					<p class="handle">@{user?.username ?? 'traveler'}</p>
				</div>
			</div>

			<button class="edit-btn">Edit profile</button>

			<div class="world">
				<div class="world-head">
					<span>World explored</span>
					<strong bind:this={worldEl}>0%</strong>
				</div>
				<div class="world-bar">
					<div class="world-fill" bind:this={worldBar}></div>
				</div>
			</div>

			<div class="stats">
				<div class="stat">
					<p bind:this={statEls[0]}>0</p>
					<span>Trips</span>
				</div>
				<div class="divider"></div>
				<div class="stat">
					<p bind:this={statEls[1]}>0</p>
					<span>Countries</span>
				</div>
				<div class="divider"></div>
				<div class="stat">
					<p bind:this={statEls[2]}>0</p>
					<span>Cities</span>
				</div>
				<div class="divider"></div>
				<div class="stat">
					<p bind:this={statEls[3]}>0</p>
					<span>Friends</span>
				</div>
			</div>

			<section class="history">
				<div class="history-head">
					<h3>Travel history</h3>
					<button class="map-btn">View Map</button>
				</div>
				<div class="history-list">
					{#if history.length === 0}
						<div class="state">No trips to show yet.</div>
					{:else}
						{#each history as trip}
							<a class="history-item" href={`/trips/${trip.id}`}>
								<div
									class="history-thumb"
									style={`background-image: url('${trip.cover_photo_url ?? ''}')`}
								></div>
								<div class="history-copy">
									<p>{getTripLocation(trip)}</p>
									<span>{formatDateRange(trip.start_date, trip.end_date)}</span>
								</div>
							</a>
						{/each}
					{/if}
				</div>
			</section>
		{/if}
	</main>
</section>

<style>
	.profile-page {
		min-height: 100dvh;
		background: var(--background-dark);
		color: var(--text-primary);
		padding-bottom: 5.5rem;
	}

	.header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 2rem 1.5rem 0.5rem;
		max-width: 480px;
		margin: 0 auto;
	}

	.eyebrow {
		color: var(--primary);
		font-size: 0.6rem;
		letter-spacing: 0.2em;
		font-weight: 700;
		text-transform: uppercase;
	}

	.settings {
		width: 2.5rem;
		height: 2.5rem;
		border-radius: 999px;
		display: grid;
		place-items: center;
		color: var(--text-secondary);
	}

	.content {
		padding: 0 1.5rem 2rem;
		overflow-y: auto;
		max-width: 480px;
		margin: 0 auto;
	}

	h1 {
		font-size: 1.9rem;
		font-weight: 700;
		margin-bottom: 2rem;
	}

	.profile-row {
		display: flex;
		align-items: center;
		gap: 1.2rem;
		margin-bottom: 1.5rem;
	}

	.avatar-wrap {
		border: 1px solid rgba(77, 157, 109, 0.4);
		border-radius: 999px;
		padding: 1px;
	}

	.avatar {
		width: 6rem;
		height: 6rem;
		border-radius: 999px;
		background-size: cover;
		background-position: center;
		border: 2px solid var(--background-dark);
	}

	.avatar.fallback {
		display: grid;
		place-items: center;
		background: rgba(255, 255, 255, 0.08);
		color: var(--text-primary);
		font-weight: 700;
	}

	h2 {
		font-size: 1.2rem;
		font-weight: 700;
	}

	.handle {
		color: var(--text-secondary);
		font-size: 0.85rem;
	}

	.edit-btn {
		width: 100%;
		padding: 0.65rem;
		border-radius: 8px;
		border: 1px solid rgba(255, 255, 255, 0.12);
		color: var(--text-primary);
		font-weight: 600;
		margin-bottom: 2rem;
	}

	.world {
		margin-bottom: 2rem;
	}

	.world-head {
		display: flex;
		align-items: flex-end;
		justify-content: space-between;
		margin-bottom: 0.6rem;
	}

	.world-head span {
		font-size: 0.65rem;
		letter-spacing: 0.2em;
		color: var(--text-secondary);
		text-transform: uppercase;
		font-weight: 700;
	}

	.world-head strong {
		font-size: 1.5rem;
		color: var(--primary);
	}

	.world-bar {
		width: 100%;
		height: 1px;
		background: rgba(255, 255, 255, 0.12);
	}

	.world-fill {
		height: 100%;
		background: var(--primary);
		width: 0%;
		transition: width 0.9s ease;
	}

	.stats {
		display: grid;
		grid-template-columns: repeat(4, 1fr);
		align-items: center;
		gap: 0.4rem;
		padding: 1rem 0;
		border-top: 1px solid rgba(255, 255, 255, 0.08);
		border-bottom: 1px solid rgba(255, 255, 255, 0.08);
		margin-bottom: 2rem;
	}

	.stat {
		text-align: center;
	}

	.stat p {
		font-size: 1.05rem;
		font-weight: 700;
	}

	.stat span {
		font-size: 0.55rem;
		text-transform: uppercase;
		letter-spacing: 0.2em;
		color: var(--text-secondary);
	}

	.divider {
		width: 1px;
		height: 2rem;
		background: rgba(255, 255, 255, 0.08);
	}

	.history-head {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: 1.5rem;
	}

	.history-head h3 {
		font-size: 0.8rem;
		letter-spacing: 0.12em;
		text-transform: uppercase;
		font-weight: 700;
	}

	.map-btn {
		font-size: 0.75rem;
		color: var(--primary);
		font-weight: 600;
	}

	.history-list {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
		padding-bottom: 2rem;
	}

	.history-item {
		display: flex;
		align-items: center;
		gap: 1rem;
	}

	.history-thumb {
		width: 3rem;
		height: 3rem;
		border-radius: 0.5rem;
		background-size: cover;
		background-position: center;
		background-color: rgba(255, 255, 255, 0.08);
	}

	.history-copy p {
		font-weight: 600;
	}

	.history-copy span {
		font-size: 0.75rem;
		color: var(--text-secondary);
	}

	.state {
		padding: 1rem;
		background: rgba(255, 255, 255, 0.04);
		border-radius: 12px;
		text-align: center;
		color: var(--text-secondary);
	}
</style>
