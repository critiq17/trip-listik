<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { apiFetch } from '$lib/api';
	import { resolvePhotoUrl } from '$lib/photos';
	import type { InviteItem } from '$lib/types';

	let invite = $state<InviteItem | null>(null);
	let loading = $state(true);
	let error = $state('');
	let processing = $state(false);

	onMount(async () => {
		const code = ($page.params as Record<string, string>).code;
		if (!code) {
			error = 'Invalid invite link';
			loading = false;
			return;
		}

		try {
			// We fetch the invite by ID/code. Assuming backend will handle this lookup.
			const res = await apiFetch<{ invite: InviteItem }>(`/v1/invites/${code}`);
			invite = res.invite;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Invite not found or expired';
		} finally {
			loading = false;
		}
	});

	const acceptInvite = async () => {
		if (!invite || processing) return;
		processing = true;
		error = '';
		try {
			await apiFetch(`/v1/invites/${invite.id}/respond`, {
				method: 'POST',
				body: JSON.stringify({ action: 'accept' })
			});
			goto(`/trips/${invite.trip_id}`);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to accept invite';
			processing = false;
		}
	};

	const declineInvite = async () => {
		if (!invite || processing) return;
		processing = true;
		error = '';
		try {
			await apiFetch(`/v1/invites/${invite.id}/respond`, {
				method: 'POST',
				body: JSON.stringify({ action: 'decline' })
			});
			goto('/inbox');
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to decline invite';
			processing = false;
		}
	};
</script>

<svelte:head>
	<title>Trip Invite</title>
</svelte:head>

<main class="invite-page">
	{#if loading}
		<div class="center-content">
			<div class="loader">Loading...</div>
		</div>
	{:else if error}
		<div class="center-content error-state">
			<span class="material-symbols-outlined icon">error</span>
			<h2>Something went wrong</h2>
			<p>{error}</p>
			<button class="primary-btn" onclick={() => goto('/')}>Go Home</button>
		</div>
	{:else if invite}
		<section class="invite-card">
			<div class="header">
				{#if invite.inviter_photo_url}
					<img src={invite.inviter_photo_url} alt="Inviter Avatar" class="avatar" />
				{:else}
					<div class="avatar fallback">
						{invite.inviter_username?.charAt(0)?.toUpperCase() || '?'}
					</div>
				{/if}
				<div class="text">
					<strong>{invite.inviter_username ? `@${invite.inviter_username}` : 'A friend'}</strong>
					<p>invited you to join their trip</p>
				</div>
			</div>

			<div class="trip-preview">
				{#if invite.trip_cover_photo_url}
					<img src={resolvePhotoUrl(invite.trip_cover_photo_url)} alt="Cover" class="cover-image" />
				{:else}
					<div class="cover-placeholder">
						<span class="material-symbols-outlined">landscape</span>
					</div>
				{/if}
			</div>
			<h1 class="trip-title">{invite.trip_title}</h1>

			<div class="actions">
				<button class="accept-btn" onclick={acceptInvite} disabled={processing}>
					{processing ? 'Processing...' : 'Accept Invite'}
				</button>
				<button class="decline-btn" onclick={declineInvite} disabled={processing}>
					Decline
				</button>
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

	.trip-title {
		font-size: 20px;
		font-weight: 700;
		line-height: 1.2;
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

	button:disabled {
		opacity: 0.5;
	}
</style>
