<script lang="ts">
	import { apiFetch } from '$lib/api';
	import { getUserInitials, getUserName } from '$lib/format';
	import { buildJoinLink, shareLinkViaTelegram } from '$lib/links';
	import type { InviteItem, UserSummary } from '$lib/types';

	let {
		open = $bindable(false),
		tripId,
		tripTitle = ''
	}: {
		open?: boolean;
		tripId: string;
		tripTitle?: string;
	} = $props();

	let query = $state('');
	let results = $state<UserSummary[]>([]);
	let pending = $state<InviteItem[]>([]);
	let loading = $state(false);
	let error = $state('');
	let inviteLoading = $state<string | null>(null);
	let shareLoading = $state(false);

	// Share a join link into any Telegram chat — works for people who never
	// opened the bot: the link itself launches the mini app on the join screen.
	const shareInviteLink = async () => {
		if (!tripId || shareLoading) return;
		shareLoading = true;
		error = '';
		try {
			const res = await apiFetch<{ token: string }>(`/v1/trips/${tripId}/invite-link`, {
				method: 'POST'
			});
			const text = tripTitle
				? `Join my trip "${tripTitle}" on TripListik!`
				: 'Join my trip on TripListik!';
			shareLinkViaTelegram(buildJoinLink(res.token), text);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to create invite link';
		} finally {
			shareLoading = false;
		}
	};

	const close = () => {
		open = false;
		query = '';
		results = [];
		error = '';
	};

	const loadInvites = async () => {
		if (!tripId) return;
		const data = await apiFetch<{ items: InviteItem[] }>(`/v1/trips/${tripId}/invites`);
		pending = data.items ?? [];
	};

	const search = async (q: string) => {
		if (!q || q.length < 2) {
			results = [];
			return;
		}
		loading = true;
		error = '';
		try {
			const data = await apiFetch<{ items: UserSummary[] }>(
				`/v1/users/search?q=${encodeURIComponent(q)}`
			);
			results = data.items ?? [];
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to search users';
		} finally {
			loading = false;
		}
	};

	const invite = async (user: UserSummary) => {
		if (!tripId) return;
		inviteLoading = user.id;
		try {
			await apiFetch(`/v1/trips/${tripId}/invite`, {
				method: 'POST',
				body: JSON.stringify({ user_id: user.id })
			});
			await loadInvites();
		} finally {
			inviteLoading = null;
		}
	};

	let timer: number | null = null;
	$effect(() => {
		if (!open) return;
		const trimmedQuery = query.trim();
		if (timer) window.clearTimeout(timer);
		timer = window.setTimeout(() => search(trimmedQuery), 250);
		return () => timer && window.clearTimeout(timer);
	});

	$effect(() => {
		if (open) loadInvites();
	});
</script>

{#if open}
	<button class="backdrop" type="button" aria-label="Close invite modal" onclick={close}></button>
	<section class="sheet">
		<header>
			<h3>Invite friends</h3>
			<button class="close" onclick={close} aria-label="Close">×</button>
		</header>

		<button class="share-link-btn" onclick={shareInviteLink} disabled={shareLoading}>
			<span class="material-symbols-outlined">send</span>
			{shareLoading ? 'Creating link...' : 'Send invite link to any chat'}
		</button>
		<p class="share-hint">Works even if they haven't opened the bot yet.</p>

		<div class="search">
			<input bind:value={query} placeholder="Search by username or name" />
			{#if loading}
				<span class="hint">Searching...</span>
			{/if}
			{#if error}
				<span class="error">{error}</span>
			{/if}
		</div>

		{#if results.length > 0}
			<div class="results">
				{#each results as user}
					<div class="row">
						<div class="avatar">
							{#if user.photo_url}
								<img src={user.photo_url} alt={getUserName(user)} />
							{:else}
								<span>{getUserInitials(user.first_name, user.last_name, user.username)}</span>
							{/if}
						</div>
						<div class="meta">
							<strong>{getUserName(user)}</strong>
							<small>@{user.username ?? 'traveler'}</small>
						</div>
						<button class="invite-btn" onclick={() => invite(user)} disabled={inviteLoading === user.id}>
							{inviteLoading === user.id ? 'Inviting...' : 'Invite'}
						</button>
					</div>
				{/each}
			</div>
		{:else if query.length >= 2 && !loading}
			<div class="empty">No users found</div>
		{/if}

		<div class="pending">
			<div class="pending-head">
				<h4>Pending invites</h4>
				<span>{pending.length}</span>
			</div>
			{#if pending.length === 0}
				<div class="empty">No pending invites</div>
			{:else}
				{#each pending as item}
					<div class="pending-row">
						<div>
							<strong>{item.inviter_username ? `@${item.inviter_username}` : 'Invite sent'}</strong>
							<small>{item.trip_title}</small>
						</div>
						<span class="status">{item.status}</span>
					</div>
				{/each}
			{/if}
		</div>
	</section>
{/if}

<style>
	.backdrop {
		position: fixed;
		inset: 0;
		background: rgba(17, 20, 24, 0.4);
		z-index: 20;
		border: none;
		padding: 0;
	}

	.sheet {
		position: fixed;
		left: 0;
		right: 0;
		bottom: 0;
		background: var(--bg-card);
		border-radius: 16px 16px 0 0;
		border-top: 1px solid var(--border);
		padding: 16px 16px 32px;
		z-index: 21;
		max-height: 80vh;
		overflow-y: auto;
		max-width: 480px;
		margin: 0 auto;
		animation: sheet-in 0.2s ease-out;
	}

	@keyframes sheet-in {
		from { transform: translateY(24px); opacity: 0; }
		to   { transform: translateY(0); opacity: 1; }
	}

	header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: 12px;
	}

	h3 {
		font-size: 16px;
		font-weight: 700;
		color: var(--text);
	}

	.close {
		width: 32px;
		height: 32px;
		border-radius: 50%;
		background: var(--bg-subtle);
		color: var(--text);
		font-size: 1.2rem;
		display: grid;
		place-items: center;
	}

	.share-link-btn {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 8px;
		width: 100%;
		padding: 12px;
		border-radius: var(--radius-input);
		background: var(--green);
		color: #fff;
		font-weight: 600;
		font-size: 14px;
		border: none;
	}

	.share-link-btn:disabled {
		opacity: 0.6;
	}

	.share-link-btn .material-symbols-outlined {
		font-size: 18px;
	}

	.share-hint {
		font-size: 12px;
		color: var(--text-sub);
		text-align: center;
		margin: 6px 0 14px;
	}

	.search input {
		width: 100%;
		padding: 12px 14px;
		border-radius: var(--radius-input);
		background: var(--bg-input);
		border: 1px solid var(--border);
		color: var(--text);
		outline: none;
		transition: border-color 0.15s ease;
	}

	.search input:focus {
		border-color: var(--green);
	}

	.hint {
		display: block;
		margin-top: 6px;
		color: var(--text-sub);
		font-size: 12px;
	}

	.error {
		display: block;
		margin-top: 6px;
		color: var(--danger);
		font-size: 12px;
	}

	.results {
		margin-top: 12px;
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	.row {
		display: grid;
		grid-template-columns: 40px 1fr auto;
		gap: 10px;
		align-items: center;
		background: var(--bg-subtle);
		padding: 10px 12px;
		border-radius: var(--radius-input);
	}

	.avatar {
		width: 40px;
		height: 40px;
		border-radius: 50%;
		overflow: hidden;
		background: var(--green-soft);
		display: grid;
		place-items: center;
		color: var(--green);
		font-weight: 600;
		font-size: 13px;
	}

	.avatar img {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.meta strong {
		display: block;
		font-size: 14px;
		font-weight: 600;
		color: var(--text);
	}

	.meta small {
		color: var(--text-sub);
		font-size: 12px;
	}

	.invite-btn {
		padding: 7px 14px;
		border-radius: var(--radius-pill);
		background: var(--green);
		color: #fff;
		font-weight: 600;
		font-size: 13px;
	}

	.invite-btn:disabled {
		opacity: 0.5;
	}

	.pending {
		margin-top: 16px;
	}

	.pending-head {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: 8px;
	}

	.pending-head h4 {
		font-size: 13px;
		font-weight: 600;
		color: var(--text-sub);
	}

	.pending-head span {
		font-size: 12px;
		color: var(--text-sub);
	}

	.pending-row {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 10px 12px;
		border-radius: var(--radius-input);
		background: var(--bg-subtle);
		margin-bottom: 8px;
	}

	.pending-row strong {
		display: block;
		font-size: 14px;
		font-weight: 600;
		color: var(--text);
	}

	.pending-row small {
		color: var(--text-sub);
		font-size: 12px;
	}

	.status {
		font-size: 12px;
		color: var(--text-sub);
	}

	.empty {
		margin-top: 8px;
		color: var(--text-sub);
		font-size: 13px;
	}
</style>
