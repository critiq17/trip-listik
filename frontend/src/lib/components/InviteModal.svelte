<script lang="ts">
	import { apiFetch } from '$lib/api';
	import { getUserInitials, getUserName } from '$lib/format';
	import type { InviteItem, UserSummary } from '$lib/types';

	let {
		open = $bindable(false),
		tripId
	}: {
		open?: boolean;
		tripId: string;
	} = $props();

	let query = $state('');
	let results = $state<UserSummary[]>([]);
	let pending = $state<InviteItem[]>([]);
	let loading = $state(false);
	let error = $state('');
	let inviteLoading = $state<string | null>(null);

	import { animate } from 'motion';

	let sheetEl: HTMLElement | null = null;
	let backdropEl: HTMLButtonElement | null = null;

	const close = () => {
		if (sheetEl && backdropEl) {
			animate(backdropEl, { opacity: [1, 0] }, { duration: 0.2 });
			animate(
				sheetEl,
				// @ts-ignore
				{ transform: ['translateY(0%)', 'translateY(100%)'], opacity: [1, 0] },
				{ duration: 0.3 }
			);
			setTimeout(() => {
				open = false;
				query = '';
				results = [];
				error = '';
			}, 300);
		} else {
			open = false;
			query = '';
			results = [];
			error = '';
		}
	};

	$effect(() => {
		if (open && sheetEl && backdropEl) {
			animate(backdropEl, { opacity: [0, 1] }, { duration: 0.2 });
			animate(
				sheetEl,
				// @ts-ignore
				{ transform: ['translateY(100%)', 'translateY(0%)'], opacity: [0.5, 1] },
				{ duration: 0.35, easing: [0.16, 1, 0.3, 1] }
			);
		}
	});

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
		background: rgba(5, 8, 7, 0.6);
		backdrop-filter: blur(4px);
		z-index: 20;
		border: none;
		padding: 0;
	}

	.sheet {
		position: fixed;
		left: 0;
		right: 0;
		bottom: 0;
		background: #0f1512;
		border-radius: var(--radius-3xl) var(--radius-3xl) 0 0;
		border-top: 1px solid rgba(255, 255, 255, 0.06);
		padding: 1.2rem 1.2rem 2rem;
		z-index: 21;
		box-shadow: 0 -24px 40px rgba(0, 0, 0, 0.35);
		max-height: 80vh;
		overflow-y: auto;
	}

	header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: 1rem;
	}

	h3 {
		font-size: 1.1rem;
	}

	.close {
		width: 32px;
		height: 32px;
		border-radius: 999px;
		background: rgba(255, 255, 255, 0.08);
		color: white;
		font-size: 1.2rem;
	}

	.search input {
		width: 100%;
		padding: 0.75rem 0.9rem;
		border-radius: var(--radius-xl);
		background: rgba(255, 255, 255, 0.06);
		border: 1px solid rgba(255, 255, 255, 0.08);
		color: white;
	}

	.hint {
		display: block;
		margin-top: 0.4rem;
		color: var(--text-secondary);
		font-size: 0.75rem;
	}

	.error {
		display: block;
		margin-top: 0.4rem;
		color: var(--error);
		font-size: 0.75rem;
	}

	.results {
		margin-top: 1rem;
		display: flex;
		flex-direction: column;
		gap: 0.8rem;
	}

	.row {
		display: grid;
		grid-template-columns: 40px 1fr auto;
		gap: 0.8rem;
		align-items: center;
		background: rgba(255, 255, 255, 0.04);
		padding: 0.6rem 0.8rem;
		border-radius: var(--radius-xl);
	}

	.avatar {
		width: 40px;
		height: 40px;
		border-radius: 999px;
		overflow: hidden;
		background: rgba(255, 255, 255, 0.08);
		display: grid;
		place-items: center;
		color: white;
		font-weight: 700;
	}

	.avatar img {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.meta strong {
		display: block;
	}

	.meta small {
		color: var(--text-secondary);
	}

	.invite-btn {
		padding: 0.4rem 0.8rem;
		border-radius: var(--radius-pill);
		background: var(--primary);
		color: #0b120f;
		font-weight: 700;
		font-size: 0.75rem;
	}

	.pending {
		margin-top: 1.2rem;
	}

	.pending-head {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: 0.6rem;
	}

	.pending-head h4 {
		font-size: 0.65rem;
		letter-spacing: 0.2em;
		text-transform: uppercase;
		color: var(--text-secondary);
	}

	.pending-row {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 0.6rem 0.8rem;
		border-radius: var(--radius-xl);
		background: rgba(255, 255, 255, 0.04);
		margin-bottom: 0.6rem;
	}

	.pending-row small {
		color: var(--text-secondary);
		font-size: 0.7rem;
	}

	.status {
		font-size: 0.7rem;
		color: var(--text-secondary);
		text-transform: uppercase;
		letter-spacing: 0.12em;
	}

	.empty {
		margin-top: 0.8rem;
		color: var(--text-secondary);
		font-size: 0.85rem;
	}
</style>
