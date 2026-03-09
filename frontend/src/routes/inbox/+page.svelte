<script lang="ts">
	import { onMount } from 'svelte';
	import { apiFetch } from '$lib/api';
	import { formatRelativeDate, parseNotificationPayload } from '$lib/format';
	import type { NotificationItem } from '$lib/types';

	let items: NotificationItem[] = [];
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
			const data = await apiFetch<{ items: NotificationItem[] }>('/v1/inbox');
			items = data.items ?? [];
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load inbox';
		} finally {
			loading = false;
		}
	});
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
		{:else if items.length === 0}
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

	.unread {
		background: rgba(32, 146, 186, 0.16);
		color: var(--accent-strong);
	}
</style>
