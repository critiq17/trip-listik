<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { apiFetch } from '$lib/api';
	import { resolvePhotoUrl } from '$lib/photos';
	import { hapticNotification } from '$lib/telegram';
	import { parseCalendarDate } from '$lib/format';
	import type { InviteLinkPreview } from '$lib/types';

	let preview = $state<InviteLinkPreview | null>(null);
	let loading = $state(true);
	let error = $state('');
	let processing = $state(false);

	const token = () => ($page.params as Record<string, string>).token ?? '';

	onMount(async () => {
		if (!token()) {
			error = 'Invalid invite link';
			loading = false;
			return;
		}
		try {
			const res = await apiFetch<{ link: InviteLinkPreview }>(`/v1/invite-links/${token()}`);
			preview = res.link;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Invite link not found or expired';
		} finally {
			loading = false;
		}
	});

	const fmtDate = (raw?: string | null) => {
		const d = parseCalendarDate(raw);
		if (!d) return '';
		return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
	};

	const dates = $derived.by(() => {
		if (!preview) return '';
		const from = fmtDate(preview.start_date);
		const to = fmtDate(preview.end_date);
		if (from && to) return `${from} – ${to}`;
		return from || to || '';
	});

	const accept = async () => {
		if (!preview || processing) return;
		processing = true;
		error = '';
		try {
			const res = await apiFetch<{ trip_id: string }>(`/v1/invite-links/${token()}/accept`, {
				method: 'POST'
			});
			hapticNotification('success');
			goto(`/trips/${res.trip_id}`);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to join the trip';
			processing = false;
		}
	};
</script>

<svelte:head>
	<title>TripListik — Trip Invite</title>
	<meta name="description" content="You have been invited to join a trip." />
</svelte:head>

<main class="invite-page">
	{#if loading}
		<div class="center-content">
			<div class="loader">Loading...</div>
		</div>
	{:else if error && !preview}
		<div class="center-content error-state">
			<span class="material-symbols-outlined icon">link_off</span>
			<h2>Invite unavailable</h2>
			<p>{error}</p>
			<button class="primary-btn" onclick={() => goto('/')}>Go Home</button>
		</div>
	{:else if preview}
		<section class="invite-card">
			<div class="header">
				{#if preview.inviter_photo_url}
					<img src={preview.inviter_photo_url} alt="Inviter avatar" class="avatar" />
				{:else}
					<div class="avatar fallback">
						{preview.inviter_name?.charAt(0)?.toUpperCase() || '?'}
					</div>
				{/if}
				<div class="text">
					<strong>{preview.inviter_name || 'A friend'}</strong>
					<p>invited you to join their trip</p>
				</div>
			</div>

			<div class="trip-preview">
				{#if preview.cover_photo_url}
					<img src={resolvePhotoUrl(preview.cover_photo_url)} alt="Cover" class="cover-image" />
				{:else}
					<div class="cover-placeholder">
						<span class="material-symbols-outlined">landscape</span>
					</div>
				{/if}
			</div>

			<div class="trip-info">
				<h1 class="trip-title">{preview.trip_title}</h1>
				<p class="trip-meta">
					{#if preview.city}{preview.city}{/if}
					{#if preview.city && dates}&nbsp;·&nbsp;{/if}
					{#if dates}{dates}{/if}
				</p>
				{#if preview.member_count > 0}
					<p class="trip-members">
						<span class="material-symbols-outlined">group</span>
						{preview.member_count} {preview.member_count === 1 ? 'traveler' : 'travelers'}
					</p>
				{/if}
			</div>

			{#if error}
				<p class="error-msg">{error}</p>
			{/if}

			<div class="actions">
				{#if preview.viewer_is_member}
					<button class="accept-btn" onclick={() => goto(`/trips/${preview?.trip_id}`)}>
						Open Trip
					</button>
					<p class="already-member">You are already on this trip.</p>
				{:else}
					<button class="accept-btn" onclick={accept} disabled={processing}>
						{processing ? 'Joining...' : 'Join Trip'}
					</button>
					<button class="decline-btn" onclick={() => goto('/')} disabled={processing}>
						Not now
					</button>
				{/if}
			</div>
		</section>
	{/if}
</main>

<style>
	.invite-page {
		min-height: 100dvh;
		display: flex;
		flex-direction: column;
		justify-content: center;
		padding: 16px;
		background: var(--bg);
		color: var(--text);
		max-width: 480px;
		margin: 0 auto;
	}

	.center-content {
		text-align: center;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 8px;
		height: 50vh;
	}

	.loader {
		color: var(--text-sub);
	}

	.error-state .icon {
		font-size: 40px;
		color: var(--danger);
		margin-bottom: 8px;
	}

	.error-state h2 {
		font-size: 18px;
		font-weight: 700;
	}

	.error-state p {
		color: var(--text-sub);
		font-size: 14px;
	}

	.primary-btn {
		margin-top: 12px;
		padding: 12px 24px;
		border-radius: var(--radius-input);
		background: var(--green);
		color: #fff;
		font-weight: 600;
		border: none;
	}

	.invite-card {
		background: var(--bg-card);
		border: 1px solid var(--border);
		border-radius: var(--radius-card);
		padding: 20px;
		display: flex;
		flex-direction: column;
		gap: 16px;
	}

	.header {
		display: flex;
		align-items: center;
		gap: 12px;
	}

	.avatar {
		width: 44px;
		height: 44px;
		border-radius: 50%;
		object-fit: cover;
	}

	.avatar.fallback {
		background: var(--green-soft);
		color: var(--green);
		display: grid;
		place-items: center;
		font-weight: 600;
		font-size: 17px;
	}

	.text strong {
		display: block;
		font-size: 15px;
		font-weight: 600;
	}

	.text p {
		color: var(--text-sub);
		font-size: 13px;
		margin-top: 2px;
	}

	.trip-preview {
		border-radius: var(--radius-input);
		overflow: hidden;
		height: 180px;
		background: var(--bg-subtle);
	}

	.cover-image {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.cover-placeholder {
		width: 100%;
		height: 100%;
		display: grid;
		place-items: center;
		color: var(--text-muted);
	}

	.cover-placeholder span {
		font-size: 40px;
	}

	.trip-info {
		display: flex;
		flex-direction: column;
		gap: 4px;
	}

	.trip-title {
		font-size: 20px;
		font-weight: 700;
		line-height: 1.2;
	}

	.trip-meta {
		font-size: 14px;
		color: var(--text-sub);
	}

	.trip-members {
		display: flex;
		align-items: center;
		gap: 6px;
		font-size: 13px;
		color: var(--text-sub);
		margin-top: 4px;
	}

	.trip-members .material-symbols-outlined {
		font-size: 18px;
	}

	.error-msg {
		color: var(--danger);
		font-size: 13px;
		text-align: center;
		padding: 10px;
		background: var(--danger-soft);
		border-radius: var(--radius-input);
	}

	.actions {
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	.accept-btn {
		padding: 13px;
		border-radius: var(--radius-input);
		background: var(--green);
		color: #fff;
		font-weight: 600;
		font-size: 15px;
		border: none;
	}

	.decline-btn {
		padding: 13px;
		border-radius: var(--radius-input);
		background: var(--bg-subtle);
		color: var(--text-sub);
		font-weight: 600;
		font-size: 15px;
		border: none;
	}

	.already-member {
		text-align: center;
		font-size: 13px;
		color: var(--text-sub);
	}

	button:disabled {
		opacity: 0.5;
	}
</style>
