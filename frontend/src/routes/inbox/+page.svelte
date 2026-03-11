<script lang="ts">
	import { onMount } from 'svelte';
	import { apiFetch } from '$lib/api';
	import { formatRelativeDate, parseNotificationPayload } from '$lib/format';
	import type { InviteItem, NotificationItem } from '$lib/types';

	let items: NotificationItem[] = [];
	let invites: InviteItem[] = [];
	let loading = true;
	let error = '';

	const typeLabels: Record<string, string> = {
		join_request: 'Join request',
		join_approved: 'Request approved',
		comment_created: 'New comment',
		trip_invite: 'Trip invite'
	};

	const describe = (item: NotificationItem) => {
		const payload = parseNotificationPayload(item.payload);
		const tripTitle = typeof payload.trip_title === 'string' ? payload.trip_title : '';
		const actorName = typeof payload.actor_name === 'string' ? payload.actor_name : '';
		switch (item.type) {
			case 'join_request':
				return actorName ? `${actorName} asked to join ${tripTitle || 'your trip'}` : 'New join request';
			case 'join_approved':
				return tripTitle ? `You were approved for ${tripTitle}` : 'Your request was approved';
			case 'comment_created':
				return tripTitle ? `Fresh discussion on ${tripTitle}` : 'New comment on a trip';
			case 'trip_invite':
				return tripTitle ? `Invitation waiting for ${tripTitle}` : 'You received a trip invite';
			default:
				return 'Travel update';
		}
	};

	onMount(async () => {
		try {
			const data = await apiFetch<{ items: NotificationItem[]; invites?: InviteItem[] }>('/v1/inbox');
			items = data.items ?? [];
			invites = data.invites ?? [];
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load inbox';
		} finally {
			loading = false;
		}
	});

	const respondInvite = async (inviteId: string, action: 'accept' | 'decline', comment?: string) => {
		await apiFetch(`/v1/invites/${inviteId}/respond`, {
			method: 'POST',
			body: JSON.stringify({ action, comment })
		});
		invites = invites.filter((inv) => inv.id !== inviteId);
	};
</script>

<section class="container">
	<header class="header">
		<p class="eyebrow">Inbox</p>
		<h1 class="headline">Requests and updates.</h1>
		<p class="subtle">Real notifications from the backend, ordered by latest activity.</p>
	</header>

	<div class="list">
		{#if loading}
			<div class="card glass">Loading updates...</div>
		{:else if error}
			<div class="card glass">{error}</div>
		{:else}
			{#if invites.length > 0}
				{#each invites as inv}
					<div class="card glass invite-card">
						<div class="row">
							<strong>Invite</strong>
							<span class="pill">{inv.status}</span>
						</div>
						<p class="message">
							{inv.inviter_username ? `@${inv.inviter_username}` : 'A friend'} invited you to
							{inv.trip_title ? ` ${inv.trip_title}` : ' a trip'}
						</p>
						<div class="actions">
							<button class="accept" onclick={() => respondInvite(inv.id, 'accept')}>Accept</button>
							<button class="decline" onclick={() => respondInvite(inv.id, 'decline')}>Decline</button>
						</div>
					</div>
				{/each}
			{/if}
			{#if items.length === 0 && invites.length === 0}
				<div class="card glass">No notifications yet</div>
			{:else}
				{#each items as item}
					<div class="card glass">
						<div class="row">
							<strong>{typeLabels[item.type] ?? item.type}</strong>
							<span class:unread={!item.read_at}>{item.read_at ? 'Read' : 'New'}</span>
						</div>
						<p class="message">{describe(item)}</p>
						<p class="muted">{formatRelativeDate(item.created_at)}</p>
					</div>
				{/each}
			{/if}
		{/if}
	</div>
</section>

<style>
	.header {
		display: grid;
		gap: 0.45rem;
		margin-bottom: 1rem;
	}

	.list {
		display: grid;
		gap: 0.75rem;
	}

	.card {
		padding: 1rem 1.1rem;
		border-radius: var(--radius-xl);
	}

	.row {
		display: flex;
		justify-content: space-between;
		align-items: center;
		gap: 1rem;
		margin-bottom: 0.45rem;
	}

	.message {
		color: white;
		line-height: 1.45;
		margin-bottom: 0.35rem;
	}

	span {
		padding: 0.35rem 0.6rem;
		border-radius: var(--radius-pill);
		background: rgba(255, 255, 255, 0.05);
		color: var(--text-secondary);
		font-size: 0.7rem;
		font-weight: 700;
		text-transform: uppercase;
		letter-spacing: 0.08em;
	}

	.invite-card {
		display: grid;
		gap: 0.6rem;
	}

	.actions {
		display: flex;
		gap: 0.6rem;
	}

	.actions button {
		flex: 1;
		padding: 0.5rem 0.75rem;
		border-radius: var(--radius-pill);
		font-weight: 700;
	}

	.accept {
		background: var(--primary);
		color: #0b120f;
	}

	.decline {
		background: rgba(255, 255, 255, 0.08);
		color: var(--text-secondary);
	}

	.pill {
		background: rgba(77, 157, 109, 0.15);
		color: var(--primary);
	}

	.unread {
		background: rgba(32, 146, 186, 0.16);
		color: var(--accent-strong);
	}
</style>
