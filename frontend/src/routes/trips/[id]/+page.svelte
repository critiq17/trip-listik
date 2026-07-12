<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { apiFetch, deleteTrip, getMe } from '$lib/api';
	import { resolvePhotoUrl } from '$lib/photos';
	import { connectTripStream } from '$lib/realtime';
	import InviteModal from '$lib/components/InviteModal.svelte';
	import { hapticImpact, hapticNotification, setupBackButton, hideBackButton } from '$lib/telegram';
	import {
		formatDateRange,
		getStatusLabel,
		getTripLocation,
		getUserInitials,
		getUserName
	} from '$lib/format';
	import type { Member, TripCardData } from '$lib/types';

	type TripDetailResponse = {
		trip: TripCardData;
		member_count: number;
		photo_count: number;
		viewer_is_member: boolean;
	};

	let trip = $state<TripCardData | null>(null);
	let loading = $state(true);
	let error = $state('');
	let members = $state<Member[]>([]);
	let memberCount = $state(0);
	let membersExpanded = $state(false);
	let joinStatus = $state('');
	let joinLoading = $state(false);
	let meId = $state('');
	let inviteOpen = $state(false);
	let stream: EventSource | null = null;
	let membersCursor = $state<string | null>(null);
	let membersLoading = $state(false);
	let membersHasMore = $state(true);
	let coverUrl = $derived(resolvePhotoUrl(trip?.cover_photo_url));
	let viewerIsMember = $state(false);
	let deletingTrip = $state(false);

	const joinStatusLabel: Record<string, string> = {
		joined: 'Joined',
		already_member: 'Already a member',
		pending: 'Request sent'
	};

	async function loadTrip(id: string) {
		loading = true;
		error = '';
		try {
			const data = await apiFetch<TripDetailResponse>(`/v1/trips/${id}`);
			trip = data.trip;
			memberCount = data.member_count ?? 0;
			viewerIsMember = data.viewer_is_member ?? false;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load trip';
		} finally {
			loading = false;
		}
	}

	async function loadMembers(id: string, mode: 'reset' | 'append' = 'reset') {
		if (membersLoading) return;
		if (mode === 'append' && !membersHasMore) return;
		membersLoading = true;
		const limit = 50;
		const cursorQuery =
			mode === 'append' && membersCursor ? `&cursor=${encodeURIComponent(membersCursor)}` : '';
		try {
			const data = await apiFetch<{ items: Member[] }>(
				`/v1/trips/${id}/members?limit=${limit}${cursorQuery}`
			);
			const items = data.items ?? [];
			if (mode === 'append') {
				members = [...members, ...items];
			} else {
				members = items;
			}
			const last = members[members.length - 1];
			membersCursor = last?.joined_at ?? null;
			membersHasMore = items.length === limit;
			memberCount = Math.max(memberCount, members.length);
		} finally {
			membersLoading = false;
		}
	}

	async function loadMe() {
		const data = await getMe();
		meId = data.user?.id ?? '';
	}

	async function joinTrip() {
		if (!trip) return;
		hapticImpact('soft');
		joinLoading = true;
		try {
			const data = await apiFetch<{ status: string }>(`/v1/trips/${trip.id}/join`, {
				method: 'POST'
			});
			joinStatus = data.status;
			hapticNotification('success');
			await loadMembers(trip.id);
		} catch {
			hapticNotification('error');
		} finally {
			joinLoading = false;
		}
	}

	async function handleDeleteTrip() {
		if (!trip) return;
		if (!confirm(`Delete "${trip.title}"? This cannot be undone.`)) return;
		deletingTrip = true;
		try {
			await deleteTrip(trip.id);
			goto('/trips');
		} catch {
			deletingTrip = false;
		}
	}

	const isOwner = () => trip?.owner_id && meId && trip.owner_id === meId;

	const canInvite = () =>
		!!trip &&
		!!meId &&
		(trip.owner_id === meId || (trip.visibility === 'group' && viewerIsMember));

	onMount(() => {
		setupBackButton(() => history.back());
		const id = $page.params.id ?? '';
		if (!id) {
			error = 'Trip id is missing';
			loading = false;
			return;
		}
		loadTrip(id).then(() => {
			if (id) loadMembers(id);
		});
		loadMe();
		stream = connectTripStream(id);

		stream.onerror = () => {
			stream?.close();
			let reconnectDelay = 3000;
			const reconnect = () => {
				stream = connectTripStream(id);
				stream.onerror = () => {
					stream?.close();
					reconnectDelay = Math.min(reconnectDelay * 2, 30000);
					setTimeout(reconnect, reconnectDelay);
				};
			};
			setTimeout(reconnect, reconnectDelay);
		};
	});

	onDestroy(() => {
		stream?.close();
		hideBackButton();
	});
</script>

<section class="trip">
	{#if loading}
		<div class="hero skeleton"></div>
	{:else if error}
		<div class="container">
			<div class="error">{error}</div>
		</div>
	{:else if trip}
		<div class="hero">
			{#if coverUrl}
				<div class="img-skeleton skeleton"></div>
				<img
					class="trip-img"
					src={coverUrl}
					alt={trip.title}
					onload={(e) => e.currentTarget.classList.add('loaded')}
					onerror={(e) => e.currentTarget.classList.add('error')}
				/>
			{:else}
				<div class="placeholder"></div>
			{/if}
			<a class="back-btn" href="/trips" aria-label="Back to trips">
				<span class="material-symbols-outlined">arrow_back</span>
			</a>
		</div>

		<div class="content container">
			<div class="head">
				<span class="badge">{getStatusLabel(trip.status)}</span>
				<h1>{trip.title}</h1>
				<p class="location">{getTripLocation(trip)}</p>
				<div class="meta-row">
					<span>{formatDateRange(trip.start_date, trip.end_date)}</span>
					<span class="dot">·</span>
					<span>{memberCount} members</span>
				</div>
			</div>

			{#if trip.description}
				<div class="block">
					<p>{trip.description}</p>
				</div>
			{/if}

			<div class="block members-block">
				<div class="members-header">
					<h2>{memberCount} Members</h2>
					<button class="toggle-btn" onclick={() => (membersExpanded = !membersExpanded)}>
						{membersExpanded ? 'Hide' : 'Show all'}
					</button>
				</div>
				<div class="avatar-row">
					{#each members.slice(0, 5) as member (member.user_id ?? member.username)}
						<div class="avatar">
							{#if member.photo_url}
								<img src={member.photo_url} alt={getUserName(member)} />
							{:else}
								<div class="avatar-fallback">
									{getUserInitials(member.first_name, member.last_name, member.username)}
								</div>
							{/if}
						</div>
					{/each}
					{#if memberCount > 5}
						<div class="avatar-more">+{memberCount - 5}</div>
					{/if}
				</div>
				{#if membersExpanded}
					<div class="member-list">
						{#each members as member (member.user_id ?? member.username)}
							<div class="member-row">
								<div class="avatar avatar--sm">
									{#if member.photo_url}
										<img src={member.photo_url} alt={getUserName(member)} />
									{:else}
										<div class="avatar-fallback">
											{getUserInitials(member.first_name, member.last_name, member.username)}
										</div>
									{/if}
								</div>
								<div class="member-info">
									<strong>{getUserName(member)}</strong>
									{#if member.role === 'owner'}
										<span class="role-badge">Owner</span>
									{/if}
								</div>
							</div>
						{/each}
					</div>
				{/if}
			</div>

			<div class="actions">
				{#if canInvite()}
					<button class="action-btn action-btn--invite" onclick={() => (inviteOpen = true)}>
						Invite
					</button>
				{/if}
				{#if isOwner()}
					<button
						class="action-btn action-btn--delete"
						onclick={handleDeleteTrip}
						disabled={deletingTrip}
					>
						{deletingTrip ? 'Deleting…' : 'Delete'}
					</button>
				{:else if trip.visibility === 'public' && !viewerIsMember}
					<button
						class="action-btn action-btn--join"
						onclick={joinTrip}
						disabled={joinLoading || !!joinStatus}
					>
						{joinLoading ? 'Joining...' : (joinStatusLabel[joinStatus] ?? joinStatus) || 'Join Trip'}
					</button>
				{/if}
			</div>
		</div>
	{/if}
</section>

<InviteModal bind:open={inviteOpen} tripId={trip?.id ?? ''} />

<style>
	.trip {
		padding-bottom: 6rem;
		background: var(--bg);
		min-height: 100dvh;
	}

	.hero {
		position: relative;
		aspect-ratio: 16 / 10;
		max-height: 360px;
		width: 100%;
		overflow: hidden;
		background: var(--bg-subtle);
	}

	.img-skeleton {
		position: absolute;
		inset: 0;
		border-radius: 0;
	}

	.hero.skeleton {
		border-radius: 0;
	}

	.trip-img {
		position: absolute;
		inset: 0;
		width: 100%;
		height: 100%;
		object-fit: cover;
		opacity: 0;
		transition: opacity 0.3s ease;
	}

	:global(.trip-img.loaded) { opacity: 1; }
	:global(.trip-img.error)  { display: none; }

	.placeholder {
		width: 100%;
		height: 100%;
		background: var(--bg-subtle);
	}

	.back-btn {
		position: absolute;
		top: 12px;
		left: 12px;
		width: 40px;
		height: 40px;
		border-radius: 50%;
		display: grid;
		place-items: center;
		background: rgba(255, 255, 255, 0.92);
		color: var(--text);
		text-decoration: none;
	}

	.back-btn .material-symbols-outlined {
		font-size: 20px;
	}

	.container {
		max-width: 480px;
		margin: 0 auto;
		padding: 0 16px;
	}

	.content {
		padding-top: 16px;
	}

	/* ── Head ────────────────────────────────────────────────── */
	.head {
		margin-bottom: 16px;
	}

	.badge {
		display: inline-block;
		padding: 3px 8px;
		border-radius: 6px;
		background: var(--green-soft);
		color: var(--green);
		font-size: 12px;
		font-weight: 600;
		margin-bottom: 8px;
	}

	h1 {
		font-size: 24px;
		font-weight: 700;
		line-height: 1.2;
		color: var(--text);
		margin-bottom: 4px;
	}

	.location {
		color: var(--text-sub);
		font-size: 14px;
	}

	.meta-row {
		display: flex;
		align-items: center;
		gap: 6px;
		margin-top: 8px;
		color: var(--text-sub);
		font-size: 13px;
	}

	.dot {
		color: var(--text-muted);
	}

	/* ── Blocks ──────────────────────────────────────────────── */
	.block {
		background: var(--bg-card);
		border: 1px solid var(--border);
		border-radius: var(--radius-card);
		padding: 16px;
		margin-bottom: 12px;
	}

	.block p {
		color: var(--text-sub);
		line-height: 1.5;
		font-size: 14px;
	}

	.error {
		padding: 16px;
		border-radius: var(--radius-card);
		text-align: center;
		color: var(--danger);
		background: var(--danger-soft);
		margin-top: 16px;
	}

	/* ── Members block ───────────────────────────────────────── */
	.members-block {
		display: flex;
		flex-direction: column;
		gap: 12px;
	}

	.members-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
	}

	.members-header h2 {
		font-size: 15px;
		font-weight: 600;
		color: var(--text);
	}

	.toggle-btn {
		background: none;
		border: none;
		color: var(--green);
		font-size: 13px;
		font-weight: 600;
		cursor: pointer;
		padding: 0;
	}

	.avatar-row {
		display: flex;
		gap: 6px;
		align-items: center;
	}

	.avatar {
		width: 40px;
		height: 40px;
		flex-shrink: 0;
		border-radius: 50%;
		overflow: hidden;
	}

	.avatar--sm {
		width: 36px;
		height: 36px;
	}

	.avatar img {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.avatar-fallback {
		width: 100%;
		height: 100%;
		display: flex;
		align-items: center;
		justify-content: center;
		background: var(--green-soft);
		color: var(--green);
		font-size: 13px;
		font-weight: 600;
		border-radius: 50%;
	}

	.avatar-more {
		width: 40px;
		height: 40px;
		border-radius: 50%;
		background: var(--bg-subtle);
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 12px;
		font-weight: 600;
		color: var(--text-sub);
		flex-shrink: 0;
	}

	.member-list {
		display: flex;
		flex-direction: column;
		gap: 10px;
		padding-top: 10px;
		border-top: 1px solid var(--border);
	}

	.member-row {
		display: flex;
		align-items: center;
		gap: 10px;
	}

	.member-info {
		flex: 1;
		display: flex;
		align-items: center;
		gap: 8px;
	}

	.member-info strong {
		color: var(--text);
		font-size: 14px;
		font-weight: 600;
	}

	.role-badge {
		padding: 2px 8px;
		border-radius: var(--radius-pill);
		background: var(--green-soft);
		color: var(--green);
		font-size: 11px;
		font-weight: 600;
	}

	/* ── Action buttons ──────────────────────────────────────── */
	.actions {
		display: flex;
		gap: 10px;
		margin-bottom: 16px;
	}

	.action-btn {
		flex: 1;
		padding: 13px 16px;
		border-radius: var(--radius-input);
		font-size: 15px;
		font-weight: 600;
		border: none;
		cursor: pointer;
		-webkit-tap-highlight-color: transparent;
	}

	.action-btn:disabled {
		opacity: 0.5;
		cursor: default;
	}

	.action-btn--invite {
		background: var(--green-soft);
		color: var(--green);
	}

	.action-btn--delete {
		background: var(--danger-soft);
		color: var(--danger);
	}

	.action-btn--join {
		background: var(--green);
		color: #fff;
	}
</style>
