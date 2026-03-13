<script lang="ts">
	import { onMount } from 'svelte';
	import { apiFetch } from '$lib/api';
	import { staggerList } from '$lib/actions/animate';
	import { formatRelativeDate, parseNotificationPayload } from '$lib/format';
	import type { InviteItem, NotificationItem } from '$lib/types';
	import { unreadCount } from '$lib/stores/notifications';

	let items = $state<NotificationItem[]>([]);
	let invites = $state<InviteItem[]>([]);
	let loading = $state(true);
	let error = $state('');

	const typeLabels: Record<string, string> = {
		join_request: 'Join request',
		join_approved: 'Request approved',
		comment_created: 'New comment',
		trip_invite: 'Trip invite',
		vote_item_created: 'New vote item'
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
			case 'vote_item_created':
				return tripTitle ? `New item to vote on in ${tripTitle}` : 'New item added';
			default:
				return 'Travel update';
		}
	};

	onMount(async () => {
		try {
			const data = await apiFetch<{ items: NotificationItem[]; invites?: InviteItem[] }>('/v1/inbox');
			items = data.items ?? [];
			// Only show pending invites
			invites = (data.invites ?? []).filter(i => i.status === 'pending');
			// Clear unread badge when viewing inbox
			unreadCount.set(0); 
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load inbox';
		} finally {
			loading = false;
		}
	});

	const respondInvite = async (inviteId: string, action: 'accept' | 'decline') => {
		// Optimistic UI hide
		invites = invites.filter((inv) => inv.id !== inviteId);
		try {
			await apiFetch(`/v1/invites/${inviteId}/respond`, {
				method: 'POST',
				body: JSON.stringify({ action })
			});
		} catch (err) {
			console.error(`Failed to ${action} invite:`, err);
			// Ideally we'd rollback, but for simplicity we rely on next reload
		}
	};
</script>

<svelte:head>
	<title>Inbox - TripListik</title>
</svelte:head>

<section class="inbox-page">
	<header class="header">
		<h1 class="serif-text">Notifications</h1>
		<p class="subtle">Invites, approvals, and trip updates.</p>
	</header>

	<div class="list" use:staggerList>
		{#if loading}
			<div class="skeleton-card"></div>
			<div class="skeleton-card"></div>
			<div class="skeleton-card"></div>
		{:else if error}
			<div class="card error-card">{error}</div>
		{:else}
			{#if invites.length > 0}
				<h3 class="section-title">Pending Invites</h3>
				{#each invites as inv (inv.id)}
					<div class="card invite-card" data-item>
						<div class="row">
							<div class="invite-info">
								<span class="material-symbols-outlined icon">mail</span>
								<strong>Trip Invitation</strong>
							</div>
							<span class="pill pending">Pending</span>
						</div>
						<p class="message">
							<strong>{inv.inviter_username ? `@${inv.inviter_username}` : 'Someone'}</strong> 
							invited you to 
							<strong>{inv.trip_title || 'a trip'}</strong>.
						</p>
						<div class="actions">
							<button class="btn accept" onclick={() => respondInvite(inv.id, 'accept')}>Accept</button>
							<button class="btn decline" onclick={() => respondInvite(inv.id, 'decline')}>Decline</button>
						</div>
					</div>
				{/each}
			{/if}

			{#if items.length > 0}
				<h3 class="section-title" style="margin-top: 1rem;">Recent Updates</h3>
				{#each items as item (item.id)}
					<div class="card notification-card" class:unread={!item.read_at} data-item>
						<div class="row">
							<strong>{typeLabels[item.type] ?? 'Update'}</strong>
							<span class="date">{formatRelativeDate(item.created_at)}</span>
						</div>
						<p class="message">{describe(item)}</p>
					</div>
				{/each}
			{/if}

			{#if items.length === 0 && invites.length === 0}
				<div class="empty-state">
					<div class="circle">
						<span class="material-symbols-outlined">notifications_off</span>
					</div>
					<h3>All caught up</h3>
					<p>You have no new notifications.</p>
				</div>
			{/if}
		{/if}
	</div>
</section>

<style>
	.inbox-page {
		padding: 1.5rem;
		padding-bottom: 6rem; /* space for bottom nav */
		min-height: 100dvh;
		background: var(--bg);
		color: var(--text);
	}

	.header {
		margin-bottom: 2rem;
	}

	.header h1 {
		font-size: 2.2rem;
		margin-bottom: 0.4rem;
	}

	.subtle {
		color: var(--text-sub);
		font-size: 0.9rem;
	}

	.section-title {
		font-size: 0.85rem;
		text-transform: uppercase;
		letter-spacing: 0.1em;
		color: var(--text-muted);
		margin-bottom: 0.8rem;
		font-weight: 700;
	}

	.list {
		display: flex;
		flex-direction: column;
		gap: 0.8rem;
	}

	.card {
		background: var(--bg-card);
		border: 1px solid var(--border);
		border-radius: var(--radius-card);
		padding: 1rem 1.25rem;
	}

	.row {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 0.6rem;
	}

	.invite-info {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.icon {
		color: var(--green);
		font-size: 1.2rem;
	}

	.message {
		font-size: 0.95rem;
		line-height: 1.4;
		color: var(--text);
	}

	.message strong {
		color: white;
	}

	.actions {
		display: flex;
		gap: 0.6rem;
		margin-top: 1rem;
	}

	.btn {
		flex: 1;
		padding: 0.6rem;
		border-radius: var(--radius-input);
		font-weight: 700;
		font-size: 0.9rem;
		border: none;
		cursor: pointer;
		-webkit-tap-highlight-color: transparent;
	}

	.accept {
		background: var(--green);
		color: var(--bg);
	}

	.accept:active {
		background: var(--green-light);
	}

	.decline {
		background: var(--bg-elevated);
		color: var(--text-sub);
		border: 1px solid var(--border);
	}

	.decline:active {
		background: rgba(255, 255, 255, 0.05);
	}

	.pill {
		padding: 0.25rem 0.6rem;
		border-radius: var(--radius-pill);
		font-size: 0.65rem;
		font-weight: 700;
		text-transform: uppercase;
		letter-spacing: 0.1em;
	}

	.pill.pending {
		background: rgba(249, 115, 22, 0.15); /* Warning orange */
		color: var(--warning);
	}

	.notification-card.unread {
		border-left: 3px solid var(--green);
		background: rgba(61, 158, 95, 0.05);
	}

	.date {
		font-size: 0.75rem;
		color: var(--text-muted);
	}

	.empty-state {
		text-align: center;
		padding: 4rem 1rem;
		color: var(--text-sub);
	}

	.circle {
		width: 4rem;
		height: 4rem;
		border-radius: 50%;
		background: var(--bg-elevated);
		display: grid;
		place-items: center;
		margin: 0 auto 1.5rem;
		color: var(--text-muted);
	}

	.circle .material-symbols-outlined {
		font-size: 2rem;
	}

	.empty-state h3 {
		font-size: 1.1rem;
		color: var(--text);
		margin-bottom: 0.4rem;
	}

	.skeleton-card {
		height: 80px;
		background: var(--bg-elevated);
		border-radius: var(--radius-card);
		animation: pulse 1.5s infinite;
	}

	@keyframes pulse {
		0% { opacity: 0.6; }
		50% { opacity: 0.3; }
		100% { opacity: 0.6; }
	}
</style>
