<script lang="ts">
	import { onMount } from 'svelte';
	import { apiFetch } from '$lib/api';
	import { countUp } from '$lib/transitions';
	import TripCard from '$lib/components/TripCard.svelte';
	import WorldMap from '$lib/components/WorldMap.svelte';
	import { getUserInitials, getUserName } from '$lib/format';
	import type { CountryVisit, TripCardData, User, UserStats } from '$lib/types';

	let user = $state<User | null>(null);
	let stats = $state<UserStats | null>(null);
	let world = $state(0);
	let history = $state<TripCardData[]>([]);
	let mapCountries = $state<CountryVisit[]>([]);
	let error = $state('');
	let statEls = $state<Array<HTMLElement | null>>([]);
	let worldEl = $state<HTMLElement | null>(null);
	let worldBar = $state<HTMLDivElement | null>(null);
	let ringEl = $state<SVGCircleElement | null>(null);
	let totalCountries = $state(0);
	let wishlist = $state<Array<{ id: string; country_code: string; city?: string; note?: string }>>(
		[]
	);
	let wishCountry = $state('');
	let wishCity = $state('');
	let wishNote = $state('');
	let shareURL = $state('');
	let editing = $state(false);
	let bio = $state('');
	let isPublic = $state(true);
	let saveError = $state('');
	let wishlistError = $state('');
	let historyError = $state('');

	onMount(async () => {
		try {
			const [me, statsRes, worldRes, mapRes] = await Promise.all([
				apiFetch<{ user: User }>('/v1/me'),
				apiFetch<UserStats>('/v1/me/stats'),
				apiFetch<{ world_explored_percent: number }>('/v1/me/world'),
				apiFetch<{ countries: CountryVisit[]; total_countries: number; world_explored_percent: number }>(
					'/v1/me/map'
				)
			]);
			user = me.user;
			bio = user?.bio ?? '';
			isPublic = user?.is_public ?? true;
			stats = statsRes;
			world = worldRes.world_explored_percent ?? 0;
			mapCountries = mapRes.countries ?? [];
			totalCountries = mapRes.total_countries ?? 0;
			if (user?.id) {
				shareURL = `${window.location.origin}/profile/${user.id}`;
			}

			try {
				const wishlistRes = await apiFetch<{
					items: Array<{ id: string; country_code: string; city?: string; note?: string }>;
				}>('/v1/me/wishlist');
				wishlist = wishlistRes.items ?? [];
			} catch (err) {
				wishlistError = err instanceof Error ? err.message : 'Failed to load wishlist';
			}

			try {
				const trips = await apiFetch<{ items: TripCardData[] }>('/v1/trips?scope=mine');
				history = (trips.items ?? []).slice(0, 6);
			} catch (err) {
				historyError = err instanceof Error ? err.message : 'Failed to load trips';
			}

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
			if (ringEl) {
				const circumference = 2 * Math.PI * 44;
				const offset = circumference - (world / 100) * circumference;
				ringEl.style.strokeDasharray = `${circumference}`;
				ringEl.style.strokeDashoffset = `${offset}`;
			}
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load profile';
		}
	});

	const copyShare = async () => {
		if (!shareURL) return;
		await navigator.clipboard.writeText(shareURL);
	};

	const saveProfile = async () => {
		saveError = '';
		try {
			const updated = await apiFetch<User>('/v1/me', {
				method: 'PATCH',
				body: JSON.stringify({ bio, is_public: isPublic })
			});
			user = updated;
			editing = false;
		} catch (err) {
			saveError = err instanceof Error ? err.message : 'Failed to update profile';
		}
	};

	const addWishlist = async () => {
		if (!wishCountry.trim()) return;
		const item = await apiFetch<{ id: string; country_code: string; city?: string; note?: string }>(
			'/v1/me/wishlist',
			{
				method: 'POST',
				body: JSON.stringify({
					country_code: wishCountry.trim().toUpperCase(),
					city: wishCity.trim(),
					note: wishNote.trim()
				})
			}
		);
		wishlist = [item, ...wishlist];
		wishCountry = '';
		wishCity = '';
		wishNote = '';
	};

	const removeWishlist = async (id: string) => {
		await apiFetch(`/v1/me/wishlist/${id}`, { method: 'DELETE' });
		wishlist = wishlist.filter((item) => item.id !== id);
	};
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

			<div class="edit-row">
				<button class="edit-btn" onclick={() => (editing = !editing)}>
					{editing ? 'Cancel' : 'Edit profile'}
				</button>
				<button class="share-btn" onclick={copyShare} disabled={!user?.is_public}>
					Share profile
				</button>
			</div>

			{#if editing}
				<div class="edit-panel glass">
					<label>
						Bio
						<textarea rows="3" bind:value={bio} placeholder="Tell the world about your travel style"></textarea>
					</label>
					<label class="toggle">
						<span>Public profile</span>
						<input type="checkbox" bind:checked={isPublic} />
					</label>
					{#if saveError}
						<div class="error">{saveError}</div>
					{/if}
					<button class="save-btn" onclick={saveProfile}>Save changes</button>
				</div>
			{/if}

			<div class="world">
				<div class="world-head">
					<span>World explored</span>
					<strong bind:this={worldEl}>0%</strong>
				</div>
				<div class="world-ring">
					<svg viewBox="0 0 100 100">
						<circle cx="50" cy="50" r="44" class="ring-bg"></circle>
						<circle cx="50" cy="50" r="44" class="ring-fill" bind:this={ringEl}></circle>
					</svg>
					<div class="ring-meta">
						<strong>{totalCountries}</strong>
						<span>countries</span>
					</div>
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
					<button class="map-btn">World Map</button>
				</div>
				<WorldMap countries={mapCountries} />
				<div class="top-countries">
					<h4>Top countries</h4>
					<div class="chips">
						{#each mapCountries.slice(0, 5) as item}
							<span class="chip">{item.code} · {item.visit_count}</span>
						{/each}
					</div>
				</div>
				<div class="history-list">
					{#if historyError}
						<div class="state">{historyError}</div>
					{:else if history.length === 0}
						<div class="state">No trips to show yet.</div>
					{:else}
						{#each history as trip}
							<TripCard {trip} variant="horizontal" />
						{/each}
					{/if}
				</div>
			</section>

			<section class="wishlist">
				<div class="history-head">
					<h3>Wishlist</h3>
				</div>
				<div class="wish-form glass">
					<input placeholder="Country code (e.g. JP)" bind:value={wishCountry} maxlength="2" />
					<input placeholder="City (optional)" bind:value={wishCity} />
					<input placeholder="Note (optional)" bind:value={wishNote} />
					<button class="save-btn" onclick={addWishlist}>Add</button>
				</div>
				<div class="history-list">
					{#if wishlistError}
						<div class="state">{wishlistError}</div>
					{:else if wishlist.length === 0}
						<div class="state">No wishlist items yet.</div>
					{:else}
						{#each wishlist as item}
							<div class="wish glass">
								<div>
									<strong>{item.country_code}</strong>
									{#if item.city}<span>{item.city}</span>{/if}
									{#if item.note}<p>{item.note}</p>{/if}
								</div>
								<button class="remove" onclick={() => removeWishlist(item.id)}>Remove</button>
							</div>
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

	.edit-row {
		display: grid;
		gap: 0.6rem;
		margin-bottom: 2rem;
	}

	.edit-btn,
	.share-btn {
		width: 100%;
		padding: 0.65rem;
		border-radius: 8px;
		border: 1px solid rgba(255, 255, 255, 0.12);
		color: var(--text-primary);
		font-weight: 600;
	}

	.share-btn {
		color: var(--primary);
		border-color: rgba(255, 255, 255, 0.08);
		background: rgba(77, 157, 109, 0.08);
	}

	.edit-panel {
		padding: 1rem;
		border-radius: var(--radius-2xl);
		display: grid;
		gap: 0.8rem;
		margin-bottom: 2rem;
	}

	.edit-panel label {
		display: grid;
		gap: 0.4rem;
		color: var(--text-secondary);
		font-size: 0.75rem;
		text-transform: uppercase;
		letter-spacing: 0.12em;
	}

	.edit-panel textarea {
		background: rgba(255, 255, 255, 0.06);
		border: 1px solid rgba(255, 255, 255, 0.08);
		color: white;
		padding: 0.6rem 0.75rem;
		border-radius: var(--radius-xl);
		resize: vertical;
		font-size: 0.85rem;
	}

	.toggle {
		display: flex;
		justify-content: space-between;
		align-items: center;
		gap: 1rem;
	}

	.toggle input {
		width: 42px;
		height: 22px;
		appearance: none;
		border-radius: 999px;
		background: rgba(255, 255, 255, 0.12);
		position: relative;
		outline: none;
		cursor: pointer;
	}

	.toggle input::after {
		content: '';
		position: absolute;
		top: 2px;
		left: 2px;
		width: 18px;
		height: 18px;
		border-radius: 999px;
		background: white;
		transition: transform 0.2s ease;
	}

	.toggle input:checked {
		background: var(--primary);
	}

	.toggle input:checked::after {
		transform: translateX(20px);
	}

	.save-btn {
		padding: 0.65rem;
		border-radius: var(--radius-pill);
		background: var(--accent-grad);
		color: #07120c;
		font-weight: 700;
	}

	.error {
		color: #f97316;
		font-size: 0.8rem;
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

	.top-countries {
		margin-top: 1rem;
		margin-bottom: 1.5rem;
	}

	.top-countries h4 {
		font-size: 0.65rem;
		letter-spacing: 0.2em;
		text-transform: uppercase;
		color: var(--text-secondary);
		margin-bottom: 0.5rem;
	}

	.chips {
		display: flex;
		flex-wrap: wrap;
		gap: 0.5rem;
	}

	.chip {
		padding: 0.35rem 0.6rem;
		border-radius: var(--radius-pill);
		background: rgba(255, 255, 255, 0.08);
		font-size: 0.7rem;
		color: var(--text-secondary);
	}

	.history-list {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
		padding-bottom: 2rem;
	}

	.world-ring {
		display: grid;
		grid-template-columns: auto 1fr;
		gap: 1rem;
		align-items: center;
		margin-bottom: 0.6rem;
	}

	.world-ring svg {
		width: 88px;
		height: 88px;
	}

	.ring-bg {
		fill: none;
		stroke: rgba(255, 255, 255, 0.08);
		stroke-width: 8;
	}

	.ring-fill {
		fill: none;
		stroke: var(--primary);
		stroke-width: 8;
		transform: rotate(-90deg);
		transform-origin: 50% 50%;
		transition: stroke-dashoffset 0.9s ease;
	}

	.ring-meta strong {
		font-size: 1.2rem;
		display: block;
	}

	.ring-meta span {
		font-size: 0.65rem;
		letter-spacing: 0.2em;
		text-transform: uppercase;
		color: var(--text-secondary);
	}

	.state {
		padding: 1rem;
		background: rgba(255, 255, 255, 0.04);
		border-radius: 12px;
		text-align: center;
		color: var(--text-secondary);
	}

	.wishlist {
		margin-top: 2rem;
	}

	.wish-form {
		display: grid;
		gap: 0.6rem;
		padding: 1rem;
		border-radius: var(--radius-2xl);
		margin-bottom: 1.2rem;
	}

	.wish-form input {
		background: rgba(255, 255, 255, 0.06);
		border: 1px solid rgba(255, 255, 255, 0.08);
		color: white;
		padding: 0.6rem 0.75rem;
		border-radius: var(--radius-xl);
	}

	.save-btn {
		padding: 0.6rem;
		border-radius: var(--radius-pill);
		background: var(--accent-grad);
		color: #07120c;
		font-weight: 700;
		font-size: 0.75rem;
	}

	.wish {
		padding: 0.9rem 1rem;
		border-radius: var(--radius-xl);
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 1rem;
	}

	.wish span {
		display: inline-block;
		margin-left: 0.5rem;
		color: var(--text-secondary);
		font-size: 0.8rem;
	}

	.wish p {
		margin-top: 0.3rem;
		color: var(--text-secondary);
		font-size: 0.8rem;
	}

	.remove {
		background: rgba(255, 255, 255, 0.08);
		border-radius: var(--radius-pill);
		padding: 0.4rem 0.7rem;
		color: var(--text-secondary);
		font-size: 0.7rem;
	}
</style>
