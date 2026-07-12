<script lang="ts">
	import { onMount } from 'svelte';
	import { apiFetch } from '$lib/api';
	import TripCard from '$lib/components/TripCard.svelte';
	import { getUserInitials, getUserName, normalizeStats, parseCalendarDate } from '$lib/format';
	import type { TripCardData, User, UserStats, UserSummary } from '$lib/types';

	let user = $state<User | null>(null);
	let stats = $state<UserStats | null>(null);
	let history = $state<TripCardData[]>([]);
	let friends = $state<UserSummary[]>([]);
	let historyError = $state('');
	let error = $state('');
	let contentTab = $state<'trips' | 'friends'>('trips');
	let historyTab = $state<'active' | 'past'>('active');
	let settingsOpen = $state(false);
	let referralCount = $state(0);
	let referralLink = $state('');
	let referralCopied = $state(false);

	onMount(async () => {
		try {
			const [meRes, statsRes, referralsRes] = await Promise.allSettled([
				apiFetch<{ user: User }>('/v1/me'),
				apiFetch<UserStats>('/v1/me/stats'),
				apiFetch<{ referral_count: number }>('/v1/me/referrals')
			]);

			if (meRes.status === 'fulfilled') {
				user = meRes.value.user;
				referralLink = `https://t.me/tripListikBot?startapp=profile_${user?.id}`;
			}

			if (statsRes.status === 'fulfilled') {
				stats = normalizeStats(statsRes.value as Record<string, unknown>);
			}

			if (referralsRes.status === 'fulfilled') {
				referralCount = referralsRes.value.referral_count ?? 0;
			}

			try {
				const [tripsRes, friendsRes] = await Promise.allSettled([
					apiFetch<{ items: TripCardData[] }>('/v1/trips?scope=mine'),
					apiFetch<{ items: UserSummary[] }>('/v1/me/friends')
				]);
				if (tripsRes.status === 'fulfilled') history = tripsRes.value.items ?? [];
				else historyError = tripsRes.reason instanceof Error ? tripsRes.reason.message : 'Failed to load trips';
				if (friendsRes.status === 'fulfilled') friends = friendsRes.value.items ?? [];
			} catch (e) {
				historyError = e instanceof Error ? e.message : 'Failed to load';
			}
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load profile';
		}
	});

	// A trip stays active through the whole last day in the user's timezone.
	const isPastTrip = (endDate?: string | null) => {
		const end = parseCalendarDate(endDate);
		if (!end) return false;
		end.setHours(23, 59, 59, 999);
		return end.getTime() < Date.now();
	};

	let activeHistory = $derived(history.filter((trip) => !isPastTrip(trip.end_date)));
	let pastHistory = $derived(history.filter((trip) => isPastTrip(trip.end_date)));

	const copyReferralLink = async () => {
		try {
			await navigator.clipboard.writeText(referralLink);
			referralCopied = true;
			setTimeout(() => (referralCopied = false), 2000);
		} catch {
			// fallback: do nothing
		}
	};
</script>

<svelte:head>
	<title>TripListik — Profile</title>
	<meta name="description" content="Your travel identity and journey history." />
</svelte:head>

<div class="page">
	<header class="top-bar">
		<span class="page-title">Profile</span>
		<button class="settings-btn" aria-label="Settings" onclick={() => (settingsOpen = true)}>
			<span class="material-symbols-outlined">settings</span>
		</button>
	</header>

	<main class="content">
		{#if error}
			<div class="state">{error}</div>
		{:else}
			<!-- Avatar + User Info -->
			<div class="profile-row">
				<div class="avatar-wrap">
					{#if user?.photo_url}
						<img class="avatar-img" src={user.photo_url} alt="Profile" loading="lazy" />
					{:else}
						<div class="avatar-fallback">{getUserInitials(user?.first_name, user?.last_name, user?.username)}</div>
					{/if}
				</div>
				<div class="user-info">
					<p class="user-name">{getUserName(user ?? {})}</p>
					<p class="user-handle">@{user?.username ?? 'traveler'}</p>
				</div>
			</div>

			<!-- Stats Row -->
			<div class="stats-row">
				<div class="stat">
					<p>{stats?.total_trips ?? 0}</p>
					<span>Trips</span>
				</div>
				<div class="stat">
					<p>{stats?.countries_visited ?? 0}</p>
					<span>Countries</span>
				</div>
				<div class="stat">
					<p>{stats?.trips_with_friends ?? 0}</p>
					<span>Friends</span>
				</div>
				<div class="stat">
					<p>{referralCount}</p>
					<span>Referrals</span>
				</div>
			</div>

			<!-- Referral link row -->
			{#if referralLink}
				<div class="referral-row">
					<input class="referral-input" type="text" readonly value={referralLink} />
					<button class="copy-btn" onclick={copyReferralLink} aria-label="Copy referral link">
						<span class="material-symbols-outlined">
							{referralCopied ? 'check' : 'content_copy'}
						</span>
					</button>
				</div>
			{/if}

			<!-- Content tabs: Trips / Friends -->
			<div class="content-tabs">
				<button class="ctab" class:active={contentTab === 'trips'} onclick={() => (contentTab = 'trips')}>
					Trips
				</button>
				<button class="ctab" class:active={contentTab === 'friends'} onclick={() => (contentTab = 'friends')}>
					Friends {#if friends.length > 0}<span class="tab-count">{friends.length}</span>{/if}
				</button>
			</div>

			{#if contentTab === 'trips'}
				<!-- Travel History -->
				<section class="history">
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
			{:else}
				<!-- Friends list -->
				<section class="friends-list">
					{#if friends.length === 0}
						<div class="state small">No friends yet. Join or create trips with others to connect.</div>
					{:else}
						{#each friends as friend (friend.id)}
							<div class="friend-row">
								<div class="friend-avatar">
									{#if friend.photo_url}
										<img src={friend.photo_url} alt={getUserName(friend)} loading="lazy" />
									{:else}
										<span>{getUserInitials(friend.first_name, friend.last_name, friend.username)}</span>
									{/if}
								</div>
								<div class="friend-info">
									<p class="friend-name">{getUserName(friend)}</p>
									<p class="friend-handle">@{friend.username ?? 'traveler'}</p>
								</div>
							</div>
						{/each}
					{/if}
				</section>
			{/if}
		{/if}
	</main>
</div>

<!-- Settings bottom sheet -->
{#if settingsOpen}
	<button
		class="sheet-backdrop"
		onclick={() => (settingsOpen = false)}
		aria-label="Close settings"
	></button>
	<div class="sheet" role="dialog" aria-label="Settings">
		<div class="sheet-handle"></div>
		<div class="sheet-header">
			<h2>Settings</h2>
			<button class="sheet-close" onclick={() => (settingsOpen = false)} aria-label="Close">×</button>
		</div>
		<div class="sheet-body">
			<p class="sheet-section-label">Account</p>
			<p class="sheet-section-sub">Manage your profile settings and preferences.</p>
		</div>
	</div>
{/if}

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

.settings-btn {
	width: 36px;
	height: 36px;
	display: flex;
	align-items: center;
	justify-content: center;
	color: var(--text-sub);
	background: none;
	border: none;
	cursor: pointer;
}

.settings-btn .material-symbols-outlined {
	font-size: 22px;
}

/* ── Content ─────────────────────────────────────────────── */
.content {
	padding: 0 16px 2rem;
	max-width: 480px;
	margin: 0 auto;
}

/* ── Profile Row ─────────────────────────────────────────── */
.profile-row {
	display: flex;
	align-items: center;
	gap: 16px;
	margin-top: 20px;
	margin-bottom: 20px;
}

.avatar-wrap {
	width: 72px;
	height: 72px;
	border-radius: 50%;
	flex-shrink: 0;
	overflow: hidden;
	border: 1px solid var(--border);
}

.avatar-img {
	width: 100%;
	height: 100%;
	object-fit: cover;
	border-radius: 50%;
}

.avatar-fallback {
	width: 100%;
	height: 100%;
	display: flex;
	align-items: center;
	justify-content: center;
	background: var(--green-soft);
	color: var(--green);
	font-size: 22px;
	font-weight: 600;
	border-radius: 50%;
}

.user-info {
	display: flex;
	flex-direction: column;
	gap: 2px;
}

.user-name {
	font-size: 18px;
	font-weight: 700;
	color: var(--text);
}

.user-handle {
	font-size: 14px;
	color: var(--text-sub);
}

/* ── Stats Row ───────────────────────────────────────────── */
.stats-row {
	display: grid;
	grid-template-columns: repeat(4, 1fr);
	padding: 12px 0;
	border-top: 1px solid var(--border);
	border-bottom: 1px solid var(--border);
	margin-bottom: 16px;
}

.stat {
	text-align: center;
}

.stat p {
	font-size: 17px;
	font-weight: 700;
	color: var(--text);
	margin-bottom: 2px;
}

.stat span {
	font-size: 12px;
	color: var(--text-sub);
}

/* ── Referral row ────────────────────────────────────────── */
.referral-row {
	display: flex;
	gap: 8px;
	align-items: center;
	margin-bottom: 20px;
}

.referral-input {
	flex: 1;
	background: var(--bg-input);
	border: 1px solid var(--border);
	border-radius: var(--radius-input);
	padding: 10px 12px;
	font-size: 12px;
	color: var(--text-sub);
	outline: none;
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
}

.copy-btn {
	width: 40px;
	height: 40px;
	display: flex;
	align-items: center;
	justify-content: center;
	background: var(--green-soft);
	border-radius: var(--radius-input);
	color: var(--green);
	cursor: pointer;
	flex-shrink: 0;
	border: none;
}

.copy-btn .material-symbols-outlined { font-size: 18px; }

/* ── Content tabs ────────────────────────────────────────── */
.content-tabs {
	display: flex;
	gap: 20px;
	margin-bottom: 12px;
	border-bottom: 1px solid var(--border);
}

.ctab {
	display: flex;
	align-items: center;
	gap: 6px;
	padding: 8px 0;
	font-size: 14px;
	font-weight: 600;
	color: var(--text-sub);
	background: none;
	border: none;
	border-bottom: 2px solid transparent;
	margin-bottom: -1px;
	cursor: pointer;
	transition: color 0.15s ease, border-color 0.15s ease;
}

.ctab.active {
	color: var(--text);
	border-bottom-color: var(--green);
}

.tab-count {
	font-size: 11px;
	background: var(--green-soft);
	color: var(--green);
	padding: 1px 6px;
	border-radius: var(--radius-pill);
}

/* ── Travel History ──────────────────────────────────────── */
.history {
	margin-bottom: 32px;
}

.history-tabs {
	display: flex;
	gap: 16px;
	margin-bottom: 8px;
}

.htab {
	font-size: 13px;
	font-weight: 500;
	color: var(--text-sub);
	background: none;
	border: none;
	border-bottom: 2px solid transparent;
	padding: 4px 0;
	cursor: pointer;
	transition: color 0.15s ease, border-color 0.15s ease;
}

.htab.active {
	color: var(--text);
	border-bottom-color: var(--green);
}

.history-list {
	display: flex;
	flex-direction: column;
}

/* ── Friends list ────────────────────────────────────────── */
.friends-list {
	display: flex;
	flex-direction: column;
}

.friend-row {
	display: flex;
	align-items: center;
	gap: 12px;
	padding: 12px 0;
	border-bottom: 1px solid var(--border);
}

.friend-avatar {
	width: 44px;
	height: 44px;
	border-radius: 50%;
	overflow: hidden;
	background: var(--green-soft);
	display: flex;
	align-items: center;
	justify-content: center;
	flex-shrink: 0;
	font-size: 15px;
	font-weight: 600;
	color: var(--green);
}

.friend-avatar img {
	width: 100%;
	height: 100%;
	object-fit: cover;
}

.friend-info {
	display: flex;
	flex-direction: column;
	gap: 2px;
}

.friend-name {
	font-size: 15px;
	font-weight: 600;
	color: var(--text);
}

.friend-handle {
	font-size: 13px;
	color: var(--text-sub);
}

/* ── Misc ────────────────────────────────────────────────── */
.state {
	padding: 20px;
	border-radius: var(--radius-card);
	background: var(--bg-subtle);
	text-align: center;
	color: var(--text-sub);
	font-size: 14px;
}

.state.small {
	padding: 16px;
	font-size: 13px;
}

/* ── Settings sheet ──────────────────────────────────────── */
.sheet-backdrop {
	position: fixed;
	inset: 0;
	z-index: 30;
	background: rgba(17, 20, 24, 0.4);
	border: none;
	cursor: pointer;
}

.sheet {
	position: fixed;
	bottom: 0;
	left: 0;
	right: 0;
	z-index: 31;
	background: var(--bg-card);
	border-radius: 16px 16px 0 0;
	padding: 12px 20px 48px;
	max-width: 480px;
	margin: 0 auto;
}

.sheet-handle {
	width: 36px;
	height: 4px;
	border-radius: var(--radius-pill);
	background: var(--border);
	margin: 0 auto 16px;
}

.sheet-header {
	display: flex;
	align-items: center;
	justify-content: space-between;
	margin-bottom: 24px;
}

.sheet-header h2 {
	font-size: 16px;
	font-weight: 700;
	color: var(--text);
}

.sheet-close {
	width: 32px;
	height: 32px;
	border-radius: 50%;
	background: var(--bg-subtle);
	color: var(--text);
	font-size: 1.4rem;
	display: grid;
	place-items: center;
	border: none;
	cursor: pointer;
}

.sheet-body {
	padding: 0 4px;
}

.sheet-section-label {
	font-size: 13px;
	font-weight: 600;
	color: var(--text);
	margin-bottom: 4px;
}

.sheet-section-sub {
	font-size: 12px;
	color: var(--text-sub);
}
</style>
