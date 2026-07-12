<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { apiFetch } from '$lib/api';
	import TripCard from '$lib/components/TripCard.svelte';
	import { getUserInitials, getUserName } from '$lib/format';
	import type { TripCardData, User, UserStats } from '$lib/types';

	type PublicProfileResponse = {
		user: User;
		stats: UserStats;
		world_explored_percent: number;
		public_trips: TripCardData[];
		wishlist: Array<{ id: string; country_code: string; city?: string; note?: string }>;
	};

	let data = $state<PublicProfileResponse | null>(null);
	let error = $state('');
	let ringEl = $state<SVGCircleElement | null>(null);

	onMount(async () => {
		const id = $page.params.id ?? '';
		if (!id) {
			error = 'Profile not found';
			return;
		}
		try {
			data = await apiFetch<PublicProfileResponse>(`/v1/users/${id}/profile`);
			if (ringEl) {
				const circumference = 2 * Math.PI * 44;
				const percent = data?.world_explored_percent ?? 0;
				const offset = circumference - (percent / 100) * circumference;
				ringEl.style.strokeDasharray = `${circumference}`;
				ringEl.style.strokeDashoffset = `${offset}`;
			}
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load profile';
		}
	});
</script>

<section class="public-profile">
	{#if error}
		<div class="state">{error}</div>
	{:else if data}
		<header class="hero">
			<div class="avatar-wrap">
				{#if data.user.photo_url}
					<div class="avatar" style={`background-image: url('${data.user.photo_url}')`}></div>
				{:else}
					<div class="avatar fallback">
						{getUserInitials(data.user.first_name, data.user.last_name, data.user.username)}
					</div>
				{/if}
			</div>
			<div>
				<h1>{getUserName(data.user)}</h1>
				<p class="handle">@{data.user.username ?? 'traveler'}</p>
				{#if data.user.bio}
					<p class="bio">{data.user.bio}</p>
				{/if}
			</div>
		</header>

		<section class="stats">
			<div class="stat">
				<p>{data.stats.total_trips}</p>
				<span>Trips</span>
			</div>
			<div class="stat">
				<p>{data.stats.countries_visited}</p>
				<span>Countries</span>
			</div>
			<div class="stat">
				<p>{data.stats.cities_visited}</p>
				<span>Cities</span>
			</div>
			<div class="stat">
				<p>{data.stats.trips_with_friends}</p>
				<span>Friends</span>
			</div>
		</section>

		<section class="map">
			<div class="map-head">
				<span>World explored</span>
				<strong>{data.world_explored_percent ?? 0}%</strong>
			</div>
			<div class="ring">
				<svg viewBox="0 0 100 100">
					<circle cx="50" cy="50" r="44" class="ring-bg"></circle>
					<circle cx="50" cy="50" r="44" class="ring-fill" bind:this={ringEl}></circle>
				</svg>
			</div>
		</section>

		<section class="section">
			<h3>Public trips</h3>
			<div class="list">
				{#if data.public_trips.length === 0}
					<div class="state">No public trips</div>
				{:else}
					{#each data.public_trips as trip}
						<TripCard {trip} variant="compact" />
					{/each}
				{/if}
			</div>
		</section>

		<section class="section">
			<h3>Wishlist</h3>
			<div class="list">
				{#if data.wishlist.length === 0}
					<div class="state">No wishlist items</div>
				{:else}
					{#each data.wishlist as item}
						<div class="wish">
							<strong>{item.country_code}</strong>
							{#if item.city}<span>{item.city}</span>{/if}
							{#if item.note}<p>{item.note}</p>{/if}
						</div>
					{/each}
				{/if}
			</div>
		</section>
	{/if}
</section>

<style>
	.public-profile {
		min-height: 100dvh;
		padding: 24px 16px 96px;
		background: var(--bg);
		color: var(--text);
		max-width: 480px;
		margin: 0 auto;
	}

	.hero {
		display: grid;
		grid-template-columns: 72px 1fr;
		gap: 16px;
		align-items: center;
		margin-bottom: 20px;
	}

	.avatar-wrap {
		border: 1px solid var(--border);
		border-radius: 50%;
		overflow: hidden;
		width: 72px;
		height: 72px;
	}

	.avatar {
		width: 100%;
		height: 100%;
		border-radius: 50%;
		background-size: cover;
		background-position: center;
	}

	.avatar.fallback {
		display: grid;
		place-items: center;
		background: var(--green-soft);
		color: var(--green);
		font-weight: 600;
		font-size: 20px;
	}

	h1 {
		font-size: 18px;
		font-weight: 700;
	}

	.handle {
		color: var(--text-sub);
		font-size: 14px;
	}

	.bio {
		margin-top: 4px;
		color: var(--text-sub);
		font-size: 13px;
	}

	.stats {
		display: grid;
		grid-template-columns: repeat(4, 1fr);
		padding: 12px 0;
		border-top: 1px solid var(--border);
		border-bottom: 1px solid var(--border);
		margin-bottom: 20px;
	}

	.stat {
		text-align: center;
	}

	.stat p {
		font-size: 17px;
		font-weight: 700;
		color: var(--text);
	}

	.stat span {
		font-size: 12px;
		color: var(--text-sub);
	}

	.map {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 12px;
		margin-bottom: 20px;
		background: var(--bg-card);
		border: 1px solid var(--border);
		border-radius: var(--radius-card);
		padding: 14px 16px;
	}

	.map-head {
		display: flex;
		flex-direction: column;
		gap: 2px;
		color: var(--text-sub);
		font-size: 13px;
	}

	.map-head strong {
		color: var(--green);
		font-size: 20px;
		font-weight: 700;
	}

	.ring svg {
		width: 56px;
		height: 56px;
		transform: rotate(-90deg);
	}

	.ring-bg {
		fill: none;
		stroke: var(--border);
		stroke-width: 8;
	}

	.ring-fill {
		fill: none;
		stroke: var(--green);
		stroke-width: 8;
		stroke-linecap: round;
	}

	.section {
		margin-bottom: 20px;
	}

	.section h3 {
		font-size: 13px;
		font-weight: 600;
		color: var(--text-sub);
		margin-bottom: 10px;
	}

	.list {
		display: grid;
		gap: 10px;
	}

	.wish {
		padding: 12px 14px;
		border-radius: var(--radius-card);
		background: var(--bg-card);
		border: 1px solid var(--border);
	}

	.wish strong {
		font-size: 14px;
		color: var(--text);
	}

	.wish span {
		margin-left: 6px;
		font-size: 13px;
		color: var(--text-sub);
	}

	.wish p {
		margin-top: 4px;
		font-size: 13px;
		color: var(--text-sub);
	}

	.state {
		padding: 16px;
		text-align: center;
		border-radius: var(--radius-card);
		background: var(--bg-subtle);
		color: var(--text-sub);
		font-size: 14px;
	}
</style>
