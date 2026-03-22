<script lang="ts">
	import { onMount } from 'svelte';
	import { apiFetch } from '$lib/api';
	import TripCard from '$lib/components/TripCard.svelte';
	import { countUp } from '$lib/transitions';
	import { getUserInitials, getUserName } from '$lib/format';
	import { hapticNotification } from '$lib/telegram';
	import type { TripCardData, User, UserStats } from '$lib/types';

	let user = $state<User | null>(null);
	let stats = $state<UserStats | null>(null);
	let history = $state<TripCardData[]>([]);
	let statEls = $state<Array<HTMLElement | null>>([]);
	let wishlist = $state<Array<{ id: string; country_code: string; city?: string; note?: string }>>([]);
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
	let error = $state('');
	let historyTab = $state<'active' | 'past'>('active');

	onMount(async () => {
		try {
			const [meRes, statsRes] = await Promise.allSettled([
				apiFetch<{ user: User }>('/v1/me'),
				apiFetch<UserStats>('/v1/me/stats')
			]);

			if (meRes.status === 'fulfilled') {
				user = meRes.value.user;
				bio = user?.bio ?? '';
				isPublic = user?.is_public ?? true;
				if (user?.id) {
					const tg = (window as any).Telegram?.WebApp;
					const botUsername = import.meta.env.PUBLIC_BOT_USERNAME || 'triplistikbot';
					if (tg?.initData) {
						shareURL = `https://t.me/${botUsername}?startapp=profile_${user.id}`;
					} else {
						shareURL = `${window.location.origin}/profile/${user.id}`;
					}
				}
			}

			if (statsRes.status === 'fulfilled') {
				stats = statsRes.value;
			}

			try {
				const wishlistRes = await apiFetch<{
					items: Array<{ id: string; country_code: string; city?: string; note?: string }>;
				}>('/v1/me/wishlist');
				wishlist = wishlistRes.items ?? [];
			} catch (e) {
				wishlistError = e instanceof Error ? e.message : 'Failed to load wishlist';
			}

			try {
				const trips = await apiFetch<{ items: TripCardData[] }>('/v1/trips?scope=mine');
				history = trips.items ?? [];
			} catch (e) {
				historyError = e instanceof Error ? e.message : 'Failed to load trips';
			}

			if (statsRes.status === 'fulfilled' && statsRes.value && statEls.length) {
				const values = [
					statsRes.value.total_trips,
					statsRes.value.countries_visited,
					statsRes.value.cities_visited,
					statsRes.value.trips_with_friends
				];
				statEls.forEach((el, idx) => {
					if (el) countUp(el, values[idx] ?? 0);
				});
			}
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load profile';
		}
	});

	const copyShare = async () => {
		if (!shareURL) return;
		try {
			await navigator.clipboard.writeText(shareURL);
			hapticNotification('success');
			const tg = (window as any).Telegram?.WebApp;
			if (tg?.showPopup) {
				tg.showPopup({ title: 'Copied', message: 'Referral link copied to clipboard!' });
			} else {
				alert('Link copied to clipboard!');
			}
		} catch (err) {
			console.error('Failed to copy', err);
		}
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

	let activeHistory = $derived(history.filter((trip) => {
		if (!trip.end_date) return true;
		return new Date(trip.end_date) >= new Date();
	}));

	let pastHistory = $derived(history.filter((trip) => {
		if (!trip.end_date) return false;
		return new Date(trip.end_date) < new Date();
	}));

	// World explored percentage (from stats or default)
	const worldPct = $derived(
		stats?.countries_visited ? Math.min((stats.countries_visited / 195) * 100, 100).toFixed(1) : '3.5'
	);
</script>

<svelte:head>
	<title>TripListik — Profile</title>
	<meta name="description" content="Your travel identity and journey history." />
</svelte:head>

<div class="page">
	<!-- Header -->
	<header class="header">
		<p class="eyebrow">Profile</p>
		<button class="settings-btn" aria-label="Settings">
			<span class="material-symbols-outlined">settings</span>
		</button>
	</header>

	<main class="content hide-scrollbar">
		<h1>Your travel identity.</h1>

		{#if error}
			<div class="state">{error}</div>
		{:else}
			<!-- Avatar + User Info -->
			<div class="profile-row">
				<div class="avatar-wrap">
					{#if user?.photo_url}
						<img class="avatar-img" src={user.photo_url} alt="Profile photo" loading="lazy" />
					{:else}
						<div class="avatar-fallback">{getUserInitials(user?.first_name, user?.last_name, user?.username)}</div>
					{/if}
				</div>
				<div class="user-info">
					<p class="user-name">{getUserName(user ?? {})}</p>
					<p class="user-handle">@{user?.username ?? 'traveler'}</p>
				</div>
			</div>

			<!-- Edit Profile Button -->
			<div class="edit-row">
				<button class="edit-btn" onclick={() => (editing = !editing)}>
					{editing ? 'Cancel' : 'Edit profile'}
				</button>
			</div>

			<!-- Edit Panel -->
			{#if editing}
				<div class="edit-panel glass">
					<label class="edit-label">
						Bio
						<textarea rows="3" bind:value={bio} placeholder="Tell the world about your travel style"></textarea>
					</label>
					<label class="toggle-row">
						<span>Public profile</span>
						<input type="checkbox" bind:checked={isPublic} />
					</label>
					{#if saveError}
						<div class="save-error">{saveError}</div>
					{/if}
					<button class="save-btn" onclick={saveProfile}>Save changes</button>
				</div>
			{/if}

			<!-- World Explored -->
			<div class="world-row">
				<p class="world-label">World Explored</p>
				<p class="world-pct">{worldPct}%</p>
			</div>
			<div class="world-bar-track">
				<div class="world-bar-fill" style="width: {worldPct}%"></div>
			</div>

			<!-- Stats Row -->
			<div class="stats-row">
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

			<!-- Travel History -->
			<section class="history">
				<div class="history-head">
					<p class="history-title">Travel History</p>
					<button class="view-map-btn">View Map</button>
				</div>

				{#if historyError}
					<div class="state">{historyError}</div>
				{:else}
					<div class="history-tabs">
						<button
							class="htab"
							class:active={historyTab === 'active'}
							onclick={() => (historyTab = 'active')}
						>Active / Upcoming</button>
						<button
							class="htab"
							class:active={historyTab === 'past'}
							onclick={() => (historyTab = 'past')}
						>Past</button>
					</div>

					<div class="history-list">
						{#if historyTab === 'active'}
							{#if activeHistory.length === 0}
								<div class="state small">No active or upcoming trips.</div>
							{:else}
								{#each activeHistory as trip}
									<TripCard {trip} variant="horizontal" />
								{/each}
							{/if}
						{:else}
							{#if pastHistory.length === 0}
								<div class="state small">No past trips yet.</div>
							{:else}
								{#each pastHistory as trip}
									<TripCard {trip} variant="horizontal" />
								{/each}
							{/if}
						{/if}
					</div>
				{/if}
			</section>

			<!-- Referral Box -->
			{#if shareURL}
				<section class="referral-section">
					<div class="referral-box">
						<p class="ref-label">Your Referral Link</p>
						<p class="ref-sub">Invite friends and earn limits</p>
						<div class="ref-row">
							<input type="text" readonly value={shareURL} class="ref-input" />
							<button class="copy-btn" onclick={copyShare} aria-label="Copy referral link">
								<span class="material-symbols-outlined">content_copy</span>
							</button>
						</div>
					</div>
				</section>
			{/if}
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
	display: flex;
	align-items: center;
	justify-content: space-between;
	padding: 32px 24px 0;
	max-width: 480px;
	margin: 0 auto;
}

.eyebrow {
	font-size: 10px;
	color: #4d9d6d;
	font-weight: 700;
	text-transform: uppercase;
	letter-spacing: 0.12em;
}

.settings-btn {
	width: 32px;
	height: 32px;
	display: flex;
	align-items: center;
	justify-content: center;
	color: #94a3b8;
	background: none;
	border: none;
	cursor: pointer;
}

.settings-btn .material-symbols-outlined {
	font-size: 22px;
}

/* ── Content ─────────────────────────────────────────────── */
.content {
	padding: 0 24px 2rem;
	max-width: 480px;
	margin: 0 auto;
}

h1 {
	font-size: 30px;
	font-weight: 700;
	color: white;
	margin: 8px 0 32px;
	line-height: 1.15;
}

/* ── Profile Row ─────────────────────────────────────────── */
.profile-row {
	display: flex;
	align-items: center;
	gap: 20px;
	margin-bottom: 24px;
}

.avatar-wrap {
	width: 96px;
	height: 96px;
	border-radius: 9999px;
	outline: 1px solid rgba(77, 157, 109, 0.4);
	outline-offset: 3px;
	flex-shrink: 0;
	overflow: hidden;
}

.avatar-img {
	width: 100%;
	height: 100%;
	object-fit: cover;
	border-radius: 9999px;
}

.avatar-fallback {
	width: 100%;
	height: 100%;
	display: flex;
	align-items: center;
	justify-content: center;
	background: rgba(255, 255, 255, 0.08);
	color: white;
	font-size: 24px;
	font-weight: 700;
	border-radius: 9999px;
}

.user-info {
	display: flex;
	flex-direction: column;
	gap: 4px;
}

.user-name {
	font-size: 20px;
	font-weight: 700;
	color: white;
}

.user-handle {
	font-size: 14px;
	color: #94a3b8;
}

/* ── Edit ───────────────────────────────────────────────── */
.edit-row {
	margin-bottom: 40px;
}

.edit-btn {
	width: 100%;
	padding: 10px;
	border: 1px solid rgba(255, 255, 255, 0.1);
	border-radius: 8px;
	background: transparent;
	color: white;
	font-size: 14px;
	font-weight: 600;
	cursor: pointer;
	transition: background 0.2s ease;
}

.edit-btn:hover {
	background: rgba(255, 255, 255, 0.04);
}

.edit-panel {
	padding: 16px;
	border-radius: 12px;
	display: flex;
	flex-direction: column;
	gap: 12px;
	margin-bottom: 24px;
}

.edit-label {
	display: flex;
	flex-direction: column;
	gap: 6px;
	font-size: 12px;
	color: #94a3b8;
	text-transform: uppercase;
	letter-spacing: 0.1em;
}

.edit-label textarea {
	background: rgba(255, 255, 255, 0.06);
	border: 1px solid rgba(255, 255, 255, 0.08);
	border-radius: 8px;
	color: white;
	padding: 10px 12px;
	font-size: 14px;
	resize: none;
	font-family: inherit;
}

.toggle-row {
	display: flex;
	justify-content: space-between;
	align-items: center;
	font-size: 14px;
	color: #94a3b8;
}

.toggle-row input {
	width: 42px;
	height: 22px;
	appearance: none;
	border-radius: 999px;
	background: rgba(255, 255, 255, 0.12);
	position: relative;
	cursor: pointer;
}

.toggle-row input::after {
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

.toggle-row input:checked {
	background: #4d9d6d;
}

.toggle-row input:checked::after {
	transform: translateX(20px);
}

.save-error {
	color: #f87171;
	font-size: 13px;
}

.save-btn {
	padding: 10px;
	border-radius: 9999px;
	background: #4d9d6d;
	color: white;
	font-size: 14px;
	font-weight: 700;
	border: none;
	cursor: pointer;
}

/* ── World Explored ──────────────────────────────────────── */
.world-row {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 8px;
}

.world-label {
	font-size: 12px;
	color: #94a3b8;
	font-weight: 600;
	text-transform: uppercase;
	letter-spacing: 0.1em;
}

.world-pct {
	font-size: 24px;
	font-weight: 700;
	color: #4d9d6d;
}

.world-bar-track {
	height: 1px;
	background: rgba(255, 255, 255, 0.08);
	margin-bottom: 40px;
}

.world-bar-fill {
	height: 100%;
	background: #4d9d6d;
	transition: width 0.8s cubic-bezier(0.4, 0, 0.2, 1);
}

/* ── Stats Row ───────────────────────────────────────────── */
.stats-row {
	display: grid;
	grid-template-columns: 1fr auto 1fr auto 1fr auto 1fr;
	align-items: center;
	padding: 16px 0;
	border-top: 1px solid rgba(255, 255, 255, 0.06);
	border-bottom: 1px solid rgba(255, 255, 255, 0.06);
	margin-bottom: 32px;
}

.stat {
	text-align: center;
}

.stat p {
	font-size: 18px;
	font-weight: 700;
	color: white;
	margin-bottom: 4px;
}

.stat span {
	font-size: 10px;
	text-transform: uppercase;
	letter-spacing: 0.1em;
	color: #64748b;
}

.divider {
	width: 1px;
	height: 32px;
	background: rgba(255, 255, 255, 0.06);
	align-self: center;
}

/* ── Travel History ──────────────────────────────────────── */
.history {
	margin-bottom: 32px;
}

.history-head {
	display: flex;
	align-items: center;
	justify-content: space-between;
	margin-bottom: 16px;
}

.history-title {
	font-size: 13px;
	font-weight: 700;
	color: white;
	text-transform: uppercase;
	letter-spacing: 0.08em;
}

.view-map-btn {
	font-size: 13px;
	font-weight: 500;
	color: #4d9d6d;
	background: none;
	border: none;
	cursor: pointer;
}

.history-tabs {
	display: flex;
	gap: 16px;
	margin-bottom: 16px;
}

.htab {
	font-size: 13px;
	font-weight: 500;
	color: #64748b;
	background: none;
	border: none;
	border-bottom: 2px solid transparent;
	padding: 4px 0;
	cursor: pointer;
	transition: color 0.2s ease, border-color 0.2s ease;
}

.htab.active {
	color: white;
	border-bottom-color: #4d9d6d;
}

.history-list {
	display: flex;
	flex-direction: column;
	gap: 0;
}

/* ── Referral ────────────────────────────────────────────── */
.referral-section {
	margin-top: 8px;
}

.referral-box {
	padding: 20px;
	border-radius: 12px;
	background: rgba(77, 157, 109, 0.05);
	border: 1px solid rgba(77, 157, 109, 0.2);
}

.ref-label {
	font-size: 15px;
	font-weight: 700;
	color: white;
	margin-bottom: 4px;
}

.ref-sub {
	font-size: 13px;
	color: #94a3b8;
	margin-bottom: 12px;
}

.ref-row {
	display: flex;
	gap: 8px;
	align-items: center;
}

.ref-input {
	flex: 1;
	background: rgba(255, 255, 255, 0.06);
	border: 1px solid rgba(255, 255, 255, 0.08);
	border-radius: 8px;
	padding: 10px 12px;
	color: white;
	font-size: 13px;
	font-family: monospace;
}

.copy-btn {
	width: 42px;
	height: 42px;
	border-radius: 8px;
	background: #4d9d6d;
	color: white;
	display: flex;
	align-items: center;
	justify-content: center;
	border: none;
	cursor: pointer;
	flex-shrink: 0;
}

.copy-btn .material-symbols-outlined {
	font-size: 20px;
}

/* ── Misc ────────────────────────────────────────────────── */
.state {
	padding: 20px;
	border-radius: 12px;
	background: rgba(255, 255, 255, 0.04);
	text-align: center;
	color: #94a3b8;
	font-size: 14px;
}

.state.small {
	padding: 16px;
	font-size: 13px;
}
</style>
