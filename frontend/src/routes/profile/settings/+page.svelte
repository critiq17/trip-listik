<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { apiFetch } from '$lib/api';
	import { setupBackButton, hideBackButton, hapticImpact } from '$lib/telegram';
	import { getUserInitials, getUserName } from '$lib/format';
	import { buildProfileLink } from '$lib/links';
	import type { User } from '$lib/types';

	let user = $state<User | null>(null);
	let error = $state('');
	let referralCopied = $state(false);

	let referralLink = $derived(user ? buildProfileLink(user.id) : '');

	onMount(async () => {
		setupBackButton(() => history.back());
		try {
			const res = await apiFetch<{ user: User }>('/v1/me');
			user = res.user;
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load account';
		}
	});

	onDestroy(() => {
		hideBackButton();
	});

	const copyReferralLink = async () => {
		if (!referralLink) return;
		try {
			await navigator.clipboard.writeText(referralLink);
			hapticImpact('light');
			referralCopied = true;
			setTimeout(() => (referralCopied = false), 2000);
		} catch {
			// clipboard unavailable: leave the row unchanged
		}
	};

	const openBot = () => {
		const tg = (window as any).Telegram?.WebApp;
		const url = 'https://t.me/tripListikBot';
		if (tg?.openTelegramLink) {
			tg.openTelegramLink(url);
		} else {
			window.open(url, '_blank');
		}
	};
</script>

<svelte:head>
	<title>TripListik — Settings</title>
	<meta name="description" content="Manage your TripListik account and preferences." />
</svelte:head>

<div class="page">
	<header class="top-bar">
		<button class="back-btn" aria-label="Go back" onclick={() => history.back()}>
			<span class="material-symbols-outlined">arrow_back</span>
		</button>
		<span class="header-title">Settings</span>
		<div class="header-spacer"></div>
	</header>

	<main class="content">
		{#if error}
			<div class="state">{error}</div>
		{:else}
			<!-- Account -->
			<p class="section-label">Account</p>
			<div class="account-card">
				<div class="avatar-wrap">
					{#if user?.photo_url}
						<img class="avatar-img" src={user.photo_url} alt="Profile" loading="lazy" />
					{:else}
						<div class="avatar-fallback">{getUserInitials(user?.first_name, user?.last_name, user?.username)}</div>
					{/if}
				</div>
				<div class="account-info">
					<p class="account-name">{getUserName(user ?? {})}</p>
					<p class="account-handle">@{user?.username ?? 'traveler'}</p>
				</div>
			</div>
			<p class="section-hint">Name, username and photo sync from your Telegram account.</p>

			<!-- Sharing -->
			<p class="section-label">Sharing</p>
			<div class="rows-card">
				<button class="row" onclick={copyReferralLink} disabled={!referralLink}>
					<span class="material-symbols-outlined row-icon">
						{referralCopied ? 'check' : 'content_copy'}
					</span>
					<span class="row-text">{referralCopied ? 'Link copied' : 'Copy referral link'}</span>
					<span class="material-symbols-outlined row-chevron">chevron_right</span>
				</button>
			</div>

			<!-- Help -->
			<p class="section-label">Help</p>
			<div class="rows-card">
				<button class="row" onclick={openBot}>
					<span class="material-symbols-outlined row-icon">smart_toy</span>
					<span class="row-text">Open TripListik bot</span>
					<span class="material-symbols-outlined row-chevron">chevron_right</span>
				</button>
			</div>
		{/if}
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
	height: 52px;
	padding: 0 8px;
	background: var(--bg);
	border-bottom: 1px solid var(--border);
	max-width: 480px;
	margin: 0 auto;
}

.back-btn {
	width: 44px;
	height: 44px;
	display: flex;
	align-items: center;
	justify-content: center;
	color: var(--text);
	background: none;
	border: none;
	cursor: pointer;
	flex-shrink: 0;
}

.back-btn .material-symbols-outlined {
	font-size: 22px;
}

.header-title {
	flex: 1;
	text-align: center;
	font-size: 16px;
	font-weight: 600;
	color: var(--text);
}

.header-spacer {
	width: 44px;
}

/* ── Content ─────────────────────────────────────────────── */
.content {
	padding: 20px 16px 2rem;
	max-width: 480px;
	margin: 0 auto;
}

.section-label {
	font-size: 13px;
	font-weight: 600;
	color: var(--text-sub);
	margin-bottom: 8px;
}

.section-hint {
	font-size: 12px;
	color: var(--text-muted);
	margin: 8px 2px 24px;
}

/* ── Account card ────────────────────────────────────────── */
.account-card {
	display: flex;
	align-items: center;
	gap: 14px;
	background: var(--bg-card);
	border: 1px solid var(--border);
	border-radius: var(--radius-card);
	padding: 14px 16px;
}

.avatar-wrap {
	width: 52px;
	height: 52px;
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
	font-size: 17px;
	font-weight: 600;
	border-radius: 50%;
}

.account-info {
	display: flex;
	flex-direction: column;
	gap: 2px;
	min-width: 0;
}

.account-name {
	font-size: 16px;
	font-weight: 700;
	color: var(--text);
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
}

.account-handle {
	font-size: 13px;
	color: var(--text-sub);
}

/* ── Setting rows ────────────────────────────────────────── */
.rows-card {
	background: var(--bg-card);
	border: 1px solid var(--border);
	border-radius: var(--radius-card);
	overflow: hidden;
	margin-bottom: 24px;
}

.row {
	display: flex;
	align-items: center;
	gap: 12px;
	width: 100%;
	min-height: 52px;
	padding: 0 16px;
	background: none;
	border: none;
	cursor: pointer;
	text-align: left;
	color: var(--text);
}

.row:disabled {
	opacity: 0.5;
	cursor: default;
}

.row-icon {
	font-size: 20px;
	color: var(--green);
	flex-shrink: 0;
}

.row-text {
	flex: 1;
	font-size: 14px;
	font-weight: 500;
}

.row-chevron {
	font-size: 20px;
	color: var(--text-muted);
	flex-shrink: 0;
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
</style>
