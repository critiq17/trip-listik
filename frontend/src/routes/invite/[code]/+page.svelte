<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { apiFetch } from '$lib/api';
	import { expandTelegram } from '$lib/telegram';
	import { scalePress } from '$lib/actions/animate';
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

<main class="invite-page container">
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
					<img src={invite.trip_cover_photo_url} alt="Cover" class="cover-image" />
				{:else}
					<div class="cover-placeholder">
						<span class="material-symbols-outlined">landscape</span>
					</div>
				{/if}
				<div class="trip-info">
					<h1>{invite.trip_title}</h1>
				</div>
			</div>

			<div class="actions">
				<button class="accept-btn" onclick={acceptInvite} disabled={processing} use:scalePress>
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
		padding: 1.5rem;
		background: var(--background-dark);
	}

	.center-content {
		text-align: center;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 1rem;
		height: 50vh;
	}

	.loader {
		color: var(--primary);
	}

	.error-state .icon {
		font-size: 3rem;
		color: var(--error);
		margin-bottom: 1rem;
	}

	.primary-btn {
		margin-top: 1rem;
		padding: 0.8rem 1.5rem;
		border-radius: var(--radius-pill);
		background: var(--primary);
		color: #000;
		font-weight: 700;
		border: none;
	}

	.invite-card {
		background: rgba(255, 255, 255, 0.04);
		border: 1px solid rgba(255, 255, 255, 0.08);
		border-radius: var(--radius-2xl);
		padding: 1.5rem;
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}

	.header {
		display: flex;
		align-items: center;
		gap: 1rem;
	}

	.avatar {
		width: 48px;
		height: 48px;
		border-radius: 999px;
		object-fit: cover;
	}

	.avatar.fallback {
		background: rgba(77, 157, 109, 0.2);
		color: var(--primary);
		display: grid;
		place-items: center;
		font-weight: 700;
		font-size: 1.2rem;
	}

	.text strong {
		display: block;
		font-size: 1.1rem;
	}

	.text p {
		color: var(--text-secondary);
		font-size: 0.85rem;
		margin-top: 0.2rem;
	}

	.trip-preview {
		border-radius: var(--radius-xl);
		overflow: hidden;
		position: relative;
		height: 200px;
	}

	.cover-image {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.cover-placeholder {
		width: 100%;
		height: 100%;
		background: rgba(0, 0, 0, 0.3);
		display: grid;
		place-items: center;
		color: rgba(255, 255, 255, 0.2);
	}

	.cover-placeholder span {
		font-size: 3rem;
	}

	.trip-info {
		position: absolute;
		bottom: 0;
		left: 0;
		right: 0;
		padding: 2rem 1rem 1rem;
		background: linear-gradient(to top, rgba(0,0,0,0.8), transparent);
	}

	.trip-info h1 {
		font-size: 1.4rem;
		font-weight: 800;
		text-shadow: 0 2px 4px rgba(0,0,0,0.5);
	}

	.actions {
		display: flex;
		flex-direction: column;
		gap: 0.8rem;
		margin-top: 0.5rem;
	}

	.accept-btn {
		padding: 1rem;
		border-radius: var(--radius-xl);
		background: var(--primary);
		color: #000;
		font-weight: 700;
		font-size: 1rem;
		border: none;
	}

	.decline-btn {
		padding: 1rem;
		border-radius: var(--radius-xl);
		background: transparent;
		color: var(--text-secondary);
		font-weight: 600;
		border: none;
	}

	button:disabled {
		opacity: 0.5;
	}
</style>
