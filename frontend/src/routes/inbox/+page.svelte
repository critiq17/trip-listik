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
			// Clear unread badge locally and on the server, otherwise it
			// reappears on the next app launch.
			unreadCount.set(0);
			apiFetch('/v1/inbox/read-all', { method: 'POST' }).catch(() => {});
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load inbox';
		} finally {
			loading = false;
		}
	});

	let actionError = $state('');

	let declineSheetOpen = $state(false);
	let decliningInviteId = $state('');
	let declineReason = $state('');
	let declineSuggestDate = $state('');
	let declineSubmitting = $state(false);

	const respondInvite = async (inviteId: string, action: 'accept' | 'decline', reason?: string, suggestDate?: string) => {
		// Optimistic UI: remove the card now, restore it if the request fails.
		const removed = invites.find((inv) => inv.id === inviteId);
		invites = invites.filter((inv) => inv.id !== inviteId);
		try {
			const body: Record<string, unknown> = { action };
			if (reason) body.comment = reason;
			if (suggestDate) body.alternative_date = suggestDate;
			await apiFetch(`/v1/invites/${inviteId}/respond`, {
				method: 'POST',
				body: JSON.stringify(body)
			});
		} catch (err) {
			if (removed) invites = [removed, ...invites];
			actionError = `Failed to ${action} the invite. Please try again.`;
			console.error(`Failed to ${action} invite:`, err);
		}
	};

	const openDeclineSheet = (inviteId: string) => {
		decliningInviteId = inviteId;
		declineReason = '';
		declineSuggestDate = '';
		declineSheetOpen = true;
	};

	const confirmDecline = async () => {
		if (!decliningInviteId) return;
		declineSubmitting = true;
		declineSheetOpen = false;
		await respondInvite(decliningInviteId, 'decline', declineReason || undefined, declineSuggestDate || undefined);
		declineSubmitting = false;
		decliningInviteId = '';
	};
</script>

<svelte:head>
	<title>TripListik — Inbox</title>
	<meta name="description" content="Your pending trip invites and travel notifications." />
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
			{#if actionError}
				<div class="card error-card">{actionError}</div>
			{/if}
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
							<button class="btn decline" onclick={() => openDeclineSheet(inv.id)}>Decline</button>
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

<!-- Decline bottom sheet -->
{#if declineSheetOpen}
	<div
		class="sheet-backdrop"
		onclick={() => (declineSheetOpen = false)}
		role="presentation"
	></div>
	<div class="decline-sheet" role="dialog" aria-label="Decline invite">
		<div class="sheet-handle"></div>
		<h3 class="sheet-title">Decline Invite</h3>

		<div class="sheet-field">
			<label class="sheet-label" for="decline-reason">Reason (optional)</label>
			<textarea
				id="decline-reason"
				class="sheet-textarea"
				rows="3"
				placeholder="Let them know why..."
				bind:value={declineReason}
			></textarea>
		</div>

		<div class="sheet-field">
			<label class="sheet-label" for="decline-date">Suggest another date (optional)</label>
			<input
				id="decline-date"
				class="sheet-input"
				type="date"
				bind:value={declineSuggestDate}
			/>
		</div>

		<button
			class="decline-confirm-btn"
			onclick={confirmDecline}
			disabled={declineSubmitting}
		>
			{declineSubmitting ? 'Declining...' : 'Decline Invite'}
		</button>
	</div>
{/if}

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

	.sheet-backdrop {
		position: fixed;
		inset: 0;
		background: rgba(0,0,0,0.5);
		z-index: 100;
	}

	.decline-sheet {
		position: fixed;
		bottom: 0;
		left: 0;
		right: 0;
		z-index: 101;
		background: #1e2a22;
		border-radius: 20px 20px 0 0;
		padding: 12px 20px 32px;
		max-width: 480px;
		margin: 0 auto;
	}

	.sheet-handle {
		width: 40px;
		height: 4px;
		background: rgba(255,255,255,0.2);
		border-radius: 2px;
		margin: 0 auto 16px;
	}

	.sheet-title {
		font-size: 17px;
		font-weight: 700;
		color: #f1f5f9;
		margin-bottom: 20px;
	}

	.sheet-field {
		margin-bottom: 16px;
	}

	.sheet-label {
		display: block;
		font-size: 12px;
		color: #94a3b8;
		font-weight: 500;
		margin-bottom: 6px;
	}

	.sheet-textarea,
	.sheet-input {
		width: 100%;
		background: rgba(255,255,255,0.06);
		border: 1px solid rgba(255,255,255,0.1);
		border-radius: 8px;
		padding: 10px 12px;
		font-size: 15px;
		color: #f1f5f9;
		outline: none;
		box-sizing: border-box;
		resize: none;
	}

	.decline-confirm-btn {
		width: 100%;
		padding: 14px;
		background: #ef4444;
		border: none;
		border-radius: 10px;
		color: white;
		font-size: 15px;
		font-weight: 700;
		cursor: pointer;
		margin-top: 8px;
	}

	.decline-confirm-btn:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}
</style>
