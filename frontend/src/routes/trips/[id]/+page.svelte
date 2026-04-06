<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { apiFetch, deleteTrip, getMe } from '$lib/api';
	import { resolvePhotoUrl } from '$lib/photos';
	import { connectTripStream } from '$lib/realtime';
	import InviteModal from '$lib/components/InviteModal.svelte';
	import { scalePress } from '$lib/actions/animate';
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
			<div class="error glass">{error}</div>
		</div>
	{:else if trip}
		<div class="hero">
			{#if coverUrl}
				<div class="img-skeleton skeleton"></div>
				<img
					class="trip-img top-image"
					src={coverUrl}
					alt={trip.title}
					onload={(e) => e.currentTarget.classList.add('loaded')}
					onerror={(e) => e.currentTarget.classList.add('error')}
				/>
			{:else}
				<div class="placeholder"></div>
			{/if}
			<div class="overlay"></div>
			<div class="floating top">
				<a class="fab-btn" href="/trips">←</a>
			</div>
			<div class="hero-content">
				<p class="badge">{getStatusLabel(trip.status)}</p>
				<h1>{trip.title}</h1>
				<p>{getTripLocation(trip)}</p>
				<div class="meta-row">
					<span>{formatDateRange(trip.start_date, trip.end_date)}</span>
					<span>{memberCount} members</span>
				</div>
			</div>
		</div>

		<div class="content container">
			{#if trip.description}
				<div class="detail-block glass">
					<p>{trip.description}</p>
				</div>
			{/if}

			<div class="detail-block glass members-block">
				<div class="members-header">
					<h2>{memberCount} Members</h2>
					<button class="toggle-btn" onclick={() => (membersExpanded = !membersExpanded)}>
						{membersExpanded ? 'Hide' : 'Show all'}
					</button>
				</div>
				<div class="avatar-row">
					{#each members.slice(0, 5) as member (member.user_id ?? member.username)}
						<div class="avatar-ring avatar">
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
								<div class="avatar-ring avatar avatar--sm">
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
						use:scalePress
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
	}

	.hero {
		position: relative;
		height: 45vh;
		min-height: 360px;
		max-height: 560px;
		overflow: hidden;
	}

	.hero img,
	.placeholder,
	.skeleton,
	.img-skeleton {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.img-skeleton {
		position: absolute;
		inset: 0;
	}

	.trip-img {
		position: absolute;
		inset: 0;
		width: 100%;
		height: 100%;
		object-fit: cover;
		opacity: 0;
		transition: opacity 0.4s ease;
	}

	:global(.trip-img.loaded) { opacity: 1; }
	:global(.trip-img.error)  { display: none; }

	.placeholder,
	.skeleton {
		background:
			radial-gradient(circle at top right, rgba(127, 191, 153, 0.12), transparent 30%),
			linear-gradient(135deg, rgba(30, 38, 32, 1), rgba(38, 48, 40, 1) 60%, rgba(30, 38, 32, 1));
	}

	.overlay {
		position: absolute;
		inset: 0;
		background: linear-gradient(180deg, rgba(3, 18, 35, 0.15) 0%, rgba(3, 18, 35, 0.88) 100%);
	}

	.floating {
		position: absolute;
		left: 1rem;
		right: 1rem;
		top: 1rem;
		z-index: 2;
		display: flex;
		justify-content: space-between;
	}

	.fab-btn {
		width: 46px;
		height: 46px;
		border-radius: 50%;
		display: grid;
		place-items: center;
		background: rgba(255, 255, 255, 0.12);
		backdrop-filter: blur(18px);
		border: 1px solid rgba(255, 255, 255, 0.14);
		color: white;
		text-decoration: none;
	}

	.hero-content {
		position: absolute;
		left: 1.25rem;
		right: 1.25rem;
		bottom: 1.25rem;
		z-index: 2;
	}

	.badge {
		display: inline-flex;
		padding: 0.45rem 0.8rem;
		border-radius: var(--radius-pill);
		background: rgba(32, 146, 186, 0.18);
		border: 1px solid rgba(122, 234, 244, 0.24);
		color: white;
		font-size: 0.72rem;
		font-weight: 800;
		letter-spacing: 0.1em;
		text-transform: uppercase;
		margin-bottom: 0.7rem;
	}

	h1 {
		font-size: clamp(2rem, 8vw, 3.2rem);
		font-weight: 800;
		line-height: 0.98;
		letter-spacing: -0.06em;
		margin-bottom: 0.45rem;
	}

	.hero-content p {
		color: var(--text-secondary);
	}

	.meta-row {
		display: flex;
		flex-wrap: wrap;
		gap: 0.55rem;
		margin-top: 0.85rem;
	}

	.meta-row span {
		padding: 0.45rem 0.7rem;
		border-radius: var(--radius-pill);
		background: rgba(255, 255, 255, 0.08);
		backdrop-filter: blur(16px);
		font-size: 0.8rem;
	}

	.container {
		max-width: 640px;
		margin: 0 auto;
		padding: 0 1rem;
	}

	.content {
		padding-top: 1rem;
	}

	.glass {
		background: rgba(255, 255, 255, 0.03);
		border: 1px solid rgba(255, 255, 255, 0.07);
		backdrop-filter: blur(12px);
	}

	.detail-block {
		padding: 1.15rem;
		border-radius: var(--radius-2xl);
		margin-bottom: 1rem;
	}

	.detail-block p {
		color: var(--text-secondary);
		line-height: 1.5;
	}

	.error {
		padding: 1rem;
		border-radius: var(--radius-xl);
		text-align: center;
		color: var(--text-secondary);
	}

	/* ── Members block ───────────────────────────────────────── */
	.members-block {
		display: flex;
		flex-direction: column;
		gap: 0.85rem;
	}

	.members-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
	}

	.members-header h2 {
		font-size: 1rem;
		font-weight: 800;
		color: white;
	}

	.toggle-btn {
		background: none;
		border: none;
		color: var(--accent-strong);
		font-size: 0.82rem;
		font-weight: 600;
		cursor: pointer;
		padding: 0;
	}

	.avatar-row {
		display: flex;
		gap: 0.4rem;
		align-items: center;
	}

	.avatar {
		width: 40px;
		height: 40px;
		flex-shrink: 0;
	}

	.avatar--sm {
		width: 36px;
		height: 36px;
	}

	.avatar-ring {
		border-radius: 50%;
		overflow: hidden;
	}

	.avatar-ring img {
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
		background: rgba(77, 157, 109, 0.18);
		color: var(--accent-strong);
		font-size: 0.9rem;
		font-weight: 700;
		border-radius: 50%;
	}

	.avatar-more {
		width: 40px;
		height: 40px;
		border-radius: 50%;
		background: rgba(255, 255, 255, 0.08);
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 0.75rem;
		font-weight: 700;
		color: var(--text-secondary);
		flex-shrink: 0;
	}

	.member-list {
		display: flex;
		flex-direction: column;
		gap: 0.6rem;
		padding-top: 0.5rem;
		border-top: 1px solid rgba(255, 255, 255, 0.06);
	}

	.member-row {
		display: flex;
		align-items: center;
		gap: 0.75rem;
	}

	.member-info {
		flex: 1;
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.member-info strong {
		color: white;
		font-size: 0.9rem;
	}

	.role-badge {
		padding: 0.2rem 0.5rem;
		border-radius: var(--radius-pill);
		background: rgba(77, 157, 109, 0.16);
		color: var(--accent-strong);
		font-size: 0.68rem;
		font-weight: 700;
	}

	/* ── Action buttons ──────────────────────────────────────── */
	.actions {
		display: flex;
		gap: 0.75rem;
		margin-bottom: 1rem;
	}

	.action-btn {
		flex: 1;
		padding: 14px 16px;
		border-radius: 12px;
		font-size: 15px;
		font-weight: 700;
		border: none;
		cursor: pointer;
		transition: opacity 0.15s ease;
		-webkit-tap-highlight-color: transparent;
	}

	.action-btn:disabled {
		opacity: 0.5;
		cursor: default;
	}

	.action-btn--invite {
		background: rgba(77, 157, 109, 0.2);
		color: #4d9d6d;
	}

	.action-btn--delete {
		background: rgba(239, 68, 68, 0.1);
		color: #ef4444;
	}

	.action-btn--join {
		background: var(--accent-grad);
		color: white;
		box-shadow: var(--shadow-glow);
	}
</style>
