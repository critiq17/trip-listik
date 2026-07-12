<script lang="ts">
	import { onMount } from 'svelte';
	import { apiFetch } from '$lib/api';
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

<div class="page">
	<header class="top-bar">
		<span class="page-title">Notifications</span>
	</header>

	<main class="list">
		{#if loading}
			<div class="skeleton-card skeleton"></div>
			<div class="skeleton-card skeleton"></div>
			<div class="skeleton-card skeleton"></div>
		{:else if error}
			<div class="card error-card">{error}</div>
		{:else}
			{#if actionError}
				<div class="card error-card">{actionError}</div>
			{/if}
			{#if invites.length > 0}
				<h3 class="section-title">Pending invites</h3>
				{#each invites as inv (inv.id)}
					<div class="card">
						<div class="row">
							<div class="invite-info">
								<span class="material-symbols-outlined icon">mail</span>
								<strong>Trip invitation</strong>
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
				<h3 class="section-title">Recent updates</h3>
				{#each items as item (item.id)}
					<div class="card" class:unread={!item.read_at}>
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
					<span class="material-symbols-outlined empty-icon">notifications_off</span>
					<h3>All caught up</h3>
					<p>You have no new notifications.</p>
				</div>
			{/if}
		{/if}
	</main>
</div>

<!-- Decline bottom sheet -->
{#if declineSheetOpen}
	<div
		class="sheet-backdrop"
		onclick={() => (declineSheetOpen = false)}
		role="presentation"
	></div>
	<div class="decline-sheet" role="dialog" aria-label="Decline invite">
		<div class="sheet-handle"></div>
		<h3 class="sheet-title">Decline invite</h3>

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
			{declineSubmitting ? 'Declining...' : 'Decline invite'}
		</button>
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

	.section-title {
		font-size: 13px;
		color: var(--text-sub);
		margin: 8px 0;
		font-weight: 600;
	}

	.list {
		display: flex;
		flex-direction: column;
		gap: 10px;
		padding: 12px 16px 0;
		max-width: 480px;
		margin: 0 auto;
	}

	.card {
		background: var(--bg-card);
		border: 1px solid var(--border);
		border-radius: var(--radius-card);
		padding: 14px 16px;
	}

	.card.unread {
		border-left: 3px solid var(--green);
	}

	.error-card {
		background: var(--danger-soft);
		border-color: transparent;
		color: var(--danger);
		text-align: center;
		font-size: 14px;
	}

	.row {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 6px;
	}

	.row strong {
		font-size: 14px;
		color: var(--text);
	}

	.invite-info {
		display: flex;
		align-items: center;
		gap: 8px;
	}

	.icon {
		color: var(--green);
		font-size: 18px;
	}

	.message {
		font-size: 14px;
		line-height: 1.4;
		color: var(--text-sub);
	}

	.message strong {
		color: var(--text);
	}

	.actions {
		display: flex;
		gap: 8px;
		margin-top: 12px;
	}

	.btn {
		flex: 1;
		padding: 10px;
		border-radius: var(--radius-input);
		font-weight: 600;
		font-size: 14px;
		border: none;
		cursor: pointer;
		-webkit-tap-highlight-color: transparent;
	}

	.accept {
		background: var(--green);
		color: #fff;
	}

	.decline {
		background: var(--bg-subtle);
		color: var(--text-sub);
	}

	.pill {
		padding: 3px 8px;
		border-radius: var(--radius-pill);
		font-size: 11px;
		font-weight: 600;
	}

	.pill.pending {
		background: var(--warning-soft);
		color: var(--warning);
	}

	.date {
		font-size: 12px;
		color: var(--text-muted);
	}

	.empty-state {
		text-align: center;
		padding: 48px 16px;
		color: var(--text-sub);
	}

	.empty-icon {
		font-size: 32px;
		color: var(--text-muted);
		margin-bottom: 12px;
	}

	.empty-state h3 {
		font-size: 16px;
		color: var(--text);
		margin-bottom: 4px;
	}

	.empty-state p {
		font-size: 14px;
	}

	.skeleton-card {
		height: 80px;
	}

	/* ── Decline sheet ──────────────────────────────────────── */
	.sheet-backdrop {
		position: fixed;
		inset: 0;
		background: rgba(17, 20, 24, 0.4);
		z-index: 100;
	}

	.decline-sheet {
		position: fixed;
		bottom: 0;
		left: 0;
		right: 0;
		z-index: 101;
		background: var(--bg-card);
		border-radius: 16px 16px 0 0;
		padding: 12px 20px 32px;
		max-width: 480px;
		margin: 0 auto;
	}

	.sheet-handle {
		width: 40px;
		height: 4px;
		background: var(--border);
		border-radius: 2px;
		margin: 0 auto 16px;
	}

	.sheet-title {
		font-size: 16px;
		font-weight: 700;
		color: var(--text);
		margin-bottom: 16px;
	}

	.sheet-field {
		margin-bottom: 14px;
	}

	.sheet-label {
		display: block;
		font-size: 13px;
		color: var(--text-sub);
		font-weight: 500;
		margin-bottom: 6px;
	}

	.sheet-textarea,
	.sheet-input {
		width: 100%;
		background: var(--bg-input);
		border: 1px solid var(--border);
		border-radius: var(--radius-input);
		padding: 10px 12px;
		font-size: 15px;
		color: var(--text);
		outline: none;
		box-sizing: border-box;
		resize: none;
	}

	.sheet-textarea:focus,
	.sheet-input:focus {
		border-color: var(--green);
	}

	.decline-confirm-btn {
		width: 100%;
		padding: 13px;
		background: var(--danger);
		border: none;
		border-radius: var(--radius-input);
		color: #fff;
		font-size: 15px;
		font-weight: 600;
		cursor: pointer;
		margin-top: 4px;
	}

	.decline-confirm-btn:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}
</style>
