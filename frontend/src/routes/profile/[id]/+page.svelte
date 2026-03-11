<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { apiFetch } from '$lib/api';
	import TripCard from '$lib/components/TripCard.svelte';
	import { countUp } from '$lib/transitions';
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
	let statEls = $state<Array<HTMLElement | null>>([]);
	let worldEl = $state<HTMLElement | null>(null);
	let ringEl = $state<SVGCircleElement | null>(null);

	onMount(async () => {
		const id = $page.params.id ?? '';
		if (!id) {
			error = 'Profile not found';
			return;
		}
		try {
			data = await apiFetch<PublicProfileResponse>(`/v1/users/${id}/profile`);
			if (data?.stats && statEls.length) {
				const values = [
					data.stats.total_trips,
					data.stats.countries_visited,
					data.stats.cities_visited,
					data.stats.trips_with_friends
				];
				statEls.forEach((el, idx) => el && countUp(el, values[idx] ?? 0));
			}
			if (worldEl) countUp(worldEl, data?.world_explored_percent ?? 0, 1.1);
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
				<p bind:this={statEls[0]}>0</p>
				<span>Trips</span>
			</div>
			<div class="stat">
				<p bind:this={statEls[1]}>0</p>
				<span>Countries</span>
			</div>
			<div class="stat">
				<p bind:this={statEls[2]}>0</p>
				<span>Cities</span>
			</div>
			<div class="stat">
				<p bind:this={statEls[3]}>0</p>
				<span>Friends</span>
			</div>
		</section>

		<section class="map">
			<div class="map-head">
				<span>World explored</span>
				<strong bind:this={worldEl}>0%</strong>
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
						<div class="wish glass">
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
		padding: 2rem 1.5rem 5.5rem;
		background: var(--background-dark);
		color: var(--text-primary);
		max-width: 560px;
		margin: 0 auto;
	}

	.hero {
		display: grid;
		grid-template-columns: 88px 1fr;
		gap: 1.2rem;
		align-items: center;
		margin-bottom: 2rem;
	}

	.avatar-wrap {
		border: 1px solid rgba(77, 157, 109, 0.4);
		border-radius: 999px;
		padding: 2px;
	}

	.avatar {
		width: 84px;
		height: 84px;
		border-radius: 999px;
		background-size: cover;
		background-position: center;
		border: 2px solid var(--background-dark);
	}

	.avatar.fallback {
		display: grid;
		place-items: center;
		background: rgba(255, 255, 255, 0.08);
		font-weight: 700;
	}

	.handle {
		color: var(--text-secondary);
	}

	.bio {
		margin-top: 0.4rem;
		color: var(--text-secondary);
	}

	.stats {
		display: grid;
		grid-template-columns: repeat(4, 1fr);
		gap: 0.6rem;
		margin-bottom: 2rem;
	}

	.stat {
		text-align: center;
	}

	.stat p {
		font-size: 1.1rem;
		font-weight: 800;
	}

	.stat span {
		font-size: 0.6rem;
		text-transform: uppercase;
		letter-spacing: 0.2em;
		color: var(--text-secondary);
	}

	.map {
		display: grid;
		gap: 0.6rem;
		margin-bottom: 2rem;
	}

	.map-head {
		display: flex;
		justify-content: space-between;
		align-items: baseline;
		color: var(--text-secondary);
		font-size: 0.7rem;
		letter-spacing: 0.2em;
		text-transform: uppercase;
	}

	.map-head strong {
		color: var(--primary);
		font-size: 1.3rem;
		letter-spacing: 0;
	}

	.ring svg {
		width: 70px;
		height: 70px;
	}

	.section h3 {
		font-size: 0.7rem;
		letter-spacing: 0.2em;
		text-transform: uppercase;
		color: var(--text-secondary);
		margin-bottom: 0.8rem;
	}

	.list {
		display: grid;
		gap: 0.8rem;
	}

	.wish {
		padding: 0.8rem 1rem;
		border-radius: var(--radius-xl);
	}

	.state {
		padding: 1rem;
		text-align: center;
		border-radius: var(--radius-xl);
		background: rgba(255, 255, 255, 0.04);
		color: var(--text-secondary);
	}
</style>
