<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import { page } from '$app/stores';
	import {
		apiFetch,
		getMe,
		getPublicPhotoURL,
		presignTripPhoto,
		uploadSignedPhoto
	} from '$lib/api';
	import { resolvePhotoUrl } from '$lib/photos';
	import { connectTripStream } from '$lib/realtime';
	import InviteModal from '$lib/components/InviteModal.svelte';
	import VoteItem from '$lib/components/VoteItem.svelte';
	import { staggerList, scalePress } from '$lib/actions/animate';
	import { animate } from 'motion';
	import { hapticImpact, hapticNotification, setupBackButton, hideBackButton } from '$lib/telegram';
	import {
		formatDateRange,
		formatLongDate,
		formatRelativeDate,
		getStatusLabel,
		getTripLocation,
		getUserInitials,
		getUserName
	} from '$lib/format';
	import type { Comment, InviteItem, JoinRequest, Member, Photo, TripCardData } from '$lib/types';

	type TripDetailResponse = {
		trip: TripCardData;
		member_count: number;
		comment_count: number;
		photo_count: number;
	};

	type VoteItemData = {
		id: string;
		trip_id: string;
		category: string;
		title: string;
		description: string;
		image_url: string;
		added_by: string;
		vote_count: number;
		has_voted: boolean;
	};

	let trip = $state<TripCardData | null>(null);
	let loading = $state(true);
	let error = $state('');
	let activeTab = $state('Details');
	let members = $state<Member[]>([]);
	let comments = $state<Comment[]>([]);
	let photos = $state<Photo[]>([]);
	let voteItems = $state<VoteItemData[]>([]);
	let memberCount = $state(0);
	let commentCount = $state(0);
	let photoCount = $state(0);
	let newComment = $state('');
	let uploading = $state(false);
	let uploadError = $state('');
	let joinStatus = $state('');

	const joinStatusLabel: Record<string, string> = {
		joined: 'Joined ✓',
		already_member: 'Already a member',
		pending: 'Request sent ✓'
	};
	let joinLoading = $state(false);
	let meId = $state('');
	let joinRequests = $state<JoinRequest[]>([]);
	let tripInvites = $state<InviteItem[]>([]);
	let inviteOpen = $state(false);
	let stream: EventSource | null = null;
	let viewerPhoto = $state<Photo | null>(null);
	let membersCursor = $state<string | null>(null);
	let membersLoading = $state(false);
	let membersHasMore = $state(true);
	let photosCursor = $state<string | null>(null);
	let photosLoading = $state(false);
	let photosHasMore = $state(true);
	let coverUrl = $derived(resolvePhotoUrl(trip?.cover_photo_url));
	let viewerIsMember = $state(false);

	const tabs = ['Details', 'Members', 'Photos', 'Discussion'];
	let loadedTabs = new Set<string>();

	async function loadTrip(id: string) {
		loading = true;
		error = '';
		try {
			const data = await apiFetch<TripDetailResponse>(`/v1/trips/${id}`);
			trip = data.trip;
			memberCount = data.member_count ?? 0;
			commentCount = data.comment_count ?? 0;
			photoCount = data.photo_count ?? 0;
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

	async function loadComments(id: string) {
		const data = await apiFetch<{ items: Comment[] }>(`/v1/trips/${id}/comments`);
		comments = data.items ?? [];
		commentCount = comments.length;
	}

	async function loadPhotos(id: string, mode: 'reset' | 'append' = 'reset') {
		if (photosLoading) return;
		if (mode === 'append' && !photosHasMore) return;
		photosLoading = true;
		const limit = 50;
		const cursorQuery =
			mode === 'append' && photosCursor ? `&cursor=${encodeURIComponent(photosCursor)}` : '';
		try {
			const data = await apiFetch<{ items: Photo[] }>(
				`/v1/trips/${id}/photos?limit=${limit}${cursorQuery}`
			);
			const items = data.items ?? [];
			if (mode === 'append') {
				photos = [...photos, ...items];
			} else {
				photos = items;
			}
			const last = photos[photos.length - 1];
			photosCursor = last?.created_at ?? null;
			photosHasMore = items.length === limit;
			photoCount = Math.max(photoCount, photos.length);
		} finally {
			photosLoading = false;
		}
	}

	async function loadVoteItems(id: string) {
		const data = await apiFetch<{ items: VoteItemData[] }>(`/v1/trips/${id}/vote-items`);
		voteItems = data.items ?? [];
	}

	const handleVoteChange = (itemId: string, newCount: number, hasVoted: boolean) => {
		voteItems = voteItems.map((item) => {
			if (item.id === itemId) {
				return { ...item, vote_count: newCount, has_voted: hasVoted };
			}
			return item;
		});
	};

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
		} catch (err) {
			hapticNotification('error');
		} finally {
			joinLoading = false;
		}
	}

	async function loadJoinRequests(id: string) {
		const data = await apiFetch<{ items: JoinRequest[] }>(
			`/v1/trips/${id}/join/requests?status=pending`
		);
		joinRequests = data.items ?? [];
	}

	async function loadInvites(id: string) {
		const data = await apiFetch<{ items: InviteItem[] }>(`/v1/trips/${id}/invites`);
		tripInvites = data.items ?? [];
	}

	async function approveJoin(userId: string) {
		if (!trip) return;
		await apiFetch(`/v1/trips/${trip.id}/join/approve`, {
			method: 'POST',
			body: JSON.stringify({ user_id: userId })
		});
		await loadJoinRequests(trip.id);
		await loadMembers(trip.id);
	}

	async function rejectJoin(userId: string) {
		if (!trip) return;
		await apiFetch(`/v1/trips/${trip.id}/join/reject`, {
			method: 'POST',
			body: JSON.stringify({ user_id: userId })
		});
		await loadJoinRequests(trip.id);
	}

	async function submitComment() {
		if (!trip || !newComment.trim()) return;
		await apiFetch(`/v1/trips/${trip.id}/comments`, {
			method: 'POST',
			body: JSON.stringify({ body: newComment })
		});
		newComment = '';
		await loadComments(trip.id);
	}

	async function handlePhotoUpload(event: Event) {
		if (!trip) return;
		const input = event.target as HTMLInputElement;
		const file = input.files?.[0];
		if (!file) return;

		uploading = true;
		uploadError = '';
		try {
			const presign = await presignTripPhoto(trip.id, file.name, file.type);
			await uploadSignedPhoto(presign.signed_url, presign.token, file);
			const publicUrl = getPublicPhotoURL(presign.path);
			await apiFetch(`/v1/trips/${trip.id}/photos`, {
				method: 'POST',
				body: JSON.stringify({ storage_path: presign.path, url: publicUrl })
			});
			await loadPhotos(trip.id);
		} catch (err) {
			uploadError = err instanceof Error ? err.message : 'Upload failed';
		} finally {
			uploading = false;
			input.value = '';
		}
	}

	const isOwner = () => trip?.owner_id && meId && trip.owner_id === meId;
	const canInvite = () =>
		!!trip &&
		!!meId &&
		(trip.owner_id === meId || (trip.visibility === 'group' && viewerIsMember));

	$effect(() => {
		if (!meId || members.length === 0) {
			viewerIsMember = false;
			return;
		}
		viewerIsMember = members.some((member) => member.user_id === meId);
	});

	onMount(() => {
		setupBackButton(() => history.back());
		const id = $page.params.id ?? '';
		if (!id) {
			error = 'Trip id is missing';
			loading = false;
			return;
		}
		loadTrip(id);
		loadMe();
		stream = connectTripStream(id);
		stream.addEventListener('comment_created', (event) => {
			const data = JSON.parse((event as MessageEvent).data) as Comment;
			comments = [data, ...comments.filter((item) => item.id !== data.id)];
			commentCount = comments.length;
			requestAnimationFrame(() => {
				const first = document.querySelector('.comments .comment:first-child');
				if (first) {
					import('motion').then(({ animate }) => {
						animate(first, { y: ['20px', '0px'], opacity: [0, 1] } as any, { duration: 0.3 } as any);
					});
				}
			});
		});
		stream.addEventListener('vote_item_vote', (event) => {
			// Realtime update for vote counts from others
			const data = JSON.parse((event as MessageEvent).data) as { item_id: string; vote_count: number };
			voteItems = voteItems.map((item) => {
				if (item.id === data.item_id) {
					// We only update count, not `has_voted` which is viewer-specific
					return { ...item, vote_count: data.vote_count };
				}
				return item;
			});
		});
		stream.addEventListener('vote_item_created', (event) => {
			const data = JSON.parse((event as MessageEvent).data) as VoteItemData;
			// Add default has_voted if missing and insert at top
			const newItem = { ...data, has_voted: false };
			voteItems = [newItem, ...voteItems.filter(v => v.id !== newItem.id)];
		});
		stream.addEventListener('photo_created', (event) => {
			const data = JSON.parse((event as MessageEvent).data) as Photo;
			const exists = photos.find((item) => item.id === data.id);
			photos = [data, ...photos.filter((item) => item.id !== data.id)];
			if (!exists) {
				photoCount += 1;
			}
		});

		stream.onerror = () => {
			stream?.close();
			setTimeout(() => {
				stream = connectTripStream(id);
			}, 3000);
		};
	});

	onDestroy(() => {
		stream?.close();
		hideBackButton();
	});

	$effect(() => {
		if (trip && activeTab === 'Members' && !loadedTabs.has('Members')) {
			loadedTabs.add('Members');
			members = [];
			membersCursor = null;
			membersHasMore = true;
			loadMembers(trip.id, 'reset');
			if (isOwner()) {
				loadJoinRequests(trip.id);
				loadInvites(trip.id);
			}
		}
	});

	$effect(() => {
		if (trip && activeTab === 'Photos' && !loadedTabs.has('Photos')) {
			loadedTabs.add('Photos');
			photos = [];
			photosCursor = null;
			photosHasMore = true;
			loadPhotos(trip.id, 'reset');
		}
	});

	$effect(() => {
		if (trip && activeTab === 'Discussion' && !loadedTabs.has('Discussion')) {
			loadedTabs.add('Discussion');
			loadComments(trip.id);
		}
	});

	$effect(() => {
		if (trip && activeTab === 'Details' && !loadedTabs.has('Details')) {
			loadedTabs.add('Details');
			loadVoteItems(trip.id);
		}
	});

	$effect(() => {
		if (!inviteOpen && trip && isOwner()) {
			loadInvites(trip.id);
		}
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
				<button class="fab-btn" aria-label="Share">↗</button>
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

		<div class="tabs-wrap">
			<div class="tabs">
				{#each tabs as tab}
					<button class:active={tab === activeTab} onclick={() => (activeTab = tab)}>
						{tab}
					</button>
				{/each}
			</div>
		</div>

		<div class="content container">
			{#if activeTab === 'Details'}
				<div class="panel detail-panel">
					<div class="stats-grid">
						<div class="stat glass">
							<strong>{trip.visibility === 'public' ? 'Public' : 'Invites'}</strong>
							<span>Visibility</span>
						</div>
						<div class="stat glass">
							<strong>{commentCount}</strong>
							<span>Discussion</span>
						</div>
						<div class="stat glass">
							<strong>{photoCount}</strong>
							<span>Photos</span>
						</div>
						<div class="stat glass">
							<strong>{memberCount}</strong>
							<span>Travelers</span>
						</div>
					</div>

					<div class="detail-block">
						<h2>Overview</h2>
						<p>{trip.description || 'No description yet. Use this trip as a living plan and update it before departure.'}</p>
					</div>

					<div class="detail-block">
						<h2>Dates</h2>
						<p>{formatLongDate(trip.start_date)} to {formatLongDate(trip.end_date)}</p>
					</div>

					<div class="detail-block map-block">
						<div>
							<h2>Route map</h2>
							<p>Live map integration is not wired yet, so location is shown from real trip fields.</p>
						</div>
						<div class="map-card">
							<strong>{getTripLocation(trip)}</strong>
							<span>{trip.visibility === 'private' ? 'Private plan' : 'Public plan'}</span>
						</div>
					</div>

					<div class="detail-block options-block">
						<h2>Options to vote on</h2>
						{#if voteItems.length === 0}
							<p class="muted">No options added yet. Group organizers will add hotels and activities here.</p>
						{:else}
							<div class="vote-item-list">
								{#each voteItems as item (item.id)}
									<VoteItem {item} onVoteChange={handleVoteChange} />
								{/each}
							</div>
						{/if}
					</div>

					<button class="cta-button" onclick={joinTrip} disabled={joinLoading || !!joinStatus} use:scalePress>
						{joinLoading ? 'Joining...' : (joinStatusLabel[joinStatus] ?? joinStatus) || 'Join Trip'}
					</button>
				</div>
			{:else if activeTab === 'Members'}
				<div class="panel member-panel">
					<div class="section-head">
						<h2>Members</h2>
						<span>{memberCount} total</span>
					</div>
					{#if canInvite()}
						<button class="invite-btn" onclick={() => (inviteOpen = true)}>Invite friends</button>
					{/if}
					{#if members.length === 0}
						<div class="empty glass">No members yet</div>
					{:else}
						<div class="member-list">
							{#each members as member}
								<div class="member-row glass">
									<div class="avatar-ring avatar">
										{#if member.photo_url}
											<img src={member.photo_url} alt={getUserName(member)} />
										{:else}
											<div class="avatar-fallback">
												{getUserInitials(member.first_name, member.last_name, member.username)}
											</div>
										{/if}
									</div>
									<div class="member-copy">
										<strong>{getUserName(member)}</strong>
										<small>@{member.username || 'traveler'}</small>
									</div>
									<span class="role">{member.role === 'owner' ? 'Crown' : 'Member'}</span>
								</div>
							{/each}
						</div>
						{#if membersHasMore}
							<button class="load-more" onclick={() => trip && loadMembers(trip.id, 'append')}>
								{membersLoading ? 'Loading...' : 'Load more'}
							</button>
						{/if}
					{/if}

					{#if isOwner()}
						<div class="requests">
							<div class="section-head">
								<h3>Pending requests</h3>
								<span>{joinRequests.length}</span>
							</div>
							{#if joinRequests.length === 0}
								<div class="empty glass">No pending requests</div>
							{:else}
								{#each joinRequests as req}
									<div class="request glass">
										<div>
											<strong>{getUserName(req)}</strong>
											<p class="muted">{formatRelativeDate(req.created_at)}</p>
										</div>
										<div class="actions">
											<button class="approve" onclick={() => approveJoin(req.user_id)}>Approve</button>
											<button class="reject" onclick={() => rejectJoin(req.user_id)}>Reject</button>
										</div>
									</div>
								{/each}
							{/if}
						</div>
						<div class="requests">
							<div class="section-head">
								<h3>Pending invites</h3>
								<span>{tripInvites.length}</span>
							</div>
							{#if tripInvites.length === 0}
								<div class="empty glass">No pending invites</div>
							{:else}
								{#each tripInvites as inv}
									<div class="request glass">
										<div>
											<strong>{inv.inviter_username ? `@${inv.inviter_username}` : 'Invite sent'}</strong>
											<p class="muted">{inv.trip_title}</p>
										</div>
										<span class="role">{inv.status}</span>
									</div>
								{/each}
							{/if}
						</div>
					{/if}
				</div>
			{:else if activeTab === 'Photos'}
				<div class="panel photo-panel">
					<div class="section-head">
						<div>
							<h2>Photos</h2>
							<p class="muted">Real uploads stored through the backend presign flow.</p>
						</div>
						<label class="upload-fab">
							<input type="file" accept="image/*" onchange={handlePhotoUpload} />
							<span>{uploading ? '...' : '+'}</span>
						</label>
					</div>

					{#if uploadError}
						<div class="error glass">{uploadError}</div>
					{/if}

					{#if photos.length === 0}
						<div class="empty glass">No photos yet</div>
					{:else}
						<div class="photo-grid">
							{#each photos as photo}
								<button class="photo" onclick={() => (viewerPhoto = photo)}>
									<div class="img-skeleton skeleton"></div>
									<img
										class="trip-img grid-img"
										src={resolvePhotoUrl(photo.url || photo.storage_path)}
										alt="Trip"
										loading="lazy"
										onload={(e) => e.currentTarget.classList.add('loaded')}
										onerror={(e) => e.currentTarget.classList.add('error')}
									/>
								</button>
							{/each}
						</div>
						{#if photosHasMore}
							<button class="load-more" onclick={() => trip && loadPhotos(trip.id, 'append')}>
								{photosLoading ? 'Loading...' : 'Load more'}
							</button>
						{/if}
					{/if}
				</div>
			{:else}
				<div class="panel discussion-panel">
					<div class="section-head">
						<h2>Discussion</h2>
						<span>{commentCount} messages</span>
					</div>
					<div class="comments">
						{#if comments.length === 0}
							<div class="empty glass">No comments yet</div>
						{:else}
							{#each comments as comment}
								<div class="comment glass">
									<div class="comment-head">
										<strong>{getUserName(comment)}</strong>
										<span>{formatRelativeDate(comment.created_at)}</span>
									</div>
									<p>{comment.body}</p>
								</div>
							{/each}
						{/if}
					</div>

					<div class="composer glass">
						<input
							bind:value={newComment}
							placeholder="Write a comment..."
							onkeydown={(e) => {
								if (e.key === 'Enter' && !e.shiftKey) {
									e.preventDefault();
									submitComment();
								}
							}}
						/>
						<button onclick={submitComment} use:scalePress>Send</button>
					</div>
				</div>
			{/if}
		</div>
	{/if}
</section>

<InviteModal bind:open={inviteOpen} tripId={trip?.id ?? ''} />

{#if viewerPhoto}
	<div class="viewer">
		<button class="viewer-backdrop" onclick={() => (viewerPhoto = null)} aria-label="Close photo viewer"></button>
		<button class="viewer-close" onclick={() => (viewerPhoto = null)}>×</button>
		<img src={resolvePhotoUrl(viewerPhoto.url || viewerPhoto.storage_path)} alt="Fullscreen trip" />
	</div>
{/if}

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
		position: relative;
		opacity: 0;
		transition: opacity 0.4s ease;
	}

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

	.tabs-wrap {
		position: sticky;
		top: 0;
		z-index: 4;
		padding: 0.85rem 1rem 0;
		background: linear-gradient(180deg, rgba(22, 28, 24, 0.98), rgba(22, 28, 24, 0.72));
		backdrop-filter: blur(18px);
	}

	.tabs {
		display: flex;
		gap: 0.55rem;
		overflow-x: auto;
		scrollbar-width: none;
	}

	.tabs::-webkit-scrollbar {
		display: none;
	}

	.tabs button {
		position: relative;
		padding: 0.8rem 0.2rem;
		margin-right: 1rem;
		color: var(--text-secondary);
		font-size: 0.88rem;
		font-weight: 700;
		white-space: nowrap;
	}

	.tabs button.active {
		color: white;
	}

	.tabs button.active::after {
		content: '';
		position: absolute;
		left: 0;
		right: 0;
		bottom: 0;
		height: 3px;
		border-radius: var(--radius-pill);
		background: var(--accent-grad);
	}

	.content {
		padding-top: 1rem;
	}

	.detail-panel,
	.member-panel,
	.photo-panel,
	.discussion-panel {
		display: grid;
		gap: 1rem;
		padding: 1.15rem;
		border-radius: var(--radius-2xl);
	}

	.stats-grid {
		display: grid;
		grid-template-columns: repeat(2, 1fr);
		gap: 0.75rem;
	}

	.stat {
		padding: 1rem;
		border-radius: var(--radius-xl);
	}

	.stat strong {
		display: block;
		font-size: 1.4rem;
		font-weight: 800;
		color: var(--accent-strong);
		margin-bottom: 0.3rem;
	}

	.stat span {
		color: var(--text-secondary);
		font-size: 0.8rem;
	}

	.detail-block h2,
	.section-head h2,
	.requests h3 {
		font-size: 1rem;
		font-weight: 800;
		margin-bottom: 0.35rem;
	}

	.invite-btn {
		padding: 0.6rem 0.9rem;
		border-radius: var(--radius-pill);
		background: var(--accent-grad);
		color: #06110b;
		font-weight: 800;
		letter-spacing: 0.08em;
		text-transform: uppercase;
		font-size: 0.65rem;
		justify-self: start;
	}

	.detail-block p,
	.map-card span {
		color: var(--text-secondary);
		line-height: 1.5;
	}

	.map-card {
		padding: 1rem;
		border-radius: var(--radius-xl);
		background:
			linear-gradient(180deg, rgba(255, 255, 255, 0.04), transparent),
			#0f1411;
		border: 1px solid rgba(255, 255, 255, 0.06);
		margin-top: 0.75rem;
	}

	.map-card strong {
		display: block;
		margin-bottom: 0.25rem;
	}

	.options-block {
		margin-top: 1rem;
	}

	.vote-item-list {
		display: grid;
		gap: 1rem;
		margin-top: 0.8rem;
	}

	.section-head {
		display: flex;
		justify-content: space-between;
		align-items: center;
		gap: 1rem;
	}

	.section-head span {
		color: var(--text-secondary);
		font-size: 0.82rem;
	}

	.member-list,
	.requests,
	.comments {
		display: grid;
		gap: 0.75rem;
	}

	.member-row,
	.request,
	.comment {
		display: flex;
		align-items: center;
		gap: 0.8rem;
		padding: 0.8rem 0.9rem;
		border-radius: var(--radius-xl);
	}

	.comment {
		display: block;
	}

	.avatar {
		width: 48px;
		height: 48px;
		flex-shrink: 0;
	}

	.member-copy {
		flex: 1;
		min-width: 0;
	}

	.member-copy small {
		color: var(--text-secondary);
	}

	.role {
		padding: 0.35rem 0.6rem;
		border-radius: var(--radius-pill);
		background: rgba(77, 157, 109, 0.16);
		color: var(--accent-strong);
		font-size: 0.72rem;
		font-weight: 700;
	}

	.actions {
		display: flex;
		gap: 0.5rem;
	}

	.actions button {
		padding: 0.7rem 0.9rem;
		border-radius: var(--radius-lg);
		font-weight: 700;
	}

	.approve {
		background: var(--accent-grad);
		color: white;
	}

	.reject {
		background: rgba(255, 255, 255, 0.06);
	}

	.upload-fab {
		position: relative;
		display: inline-flex;
		align-items: center;
		justify-content: center;
		width: 48px;
		height: 48px;
		border-radius: 50%;
		background: var(--accent-grad);
		color: white;
		font-size: 1.5rem;
		font-weight: 700;
		box-shadow: var(--shadow-glow);
	}

	.upload-fab input {
		position: absolute;
		inset: 0;
		opacity: 0;
		cursor: pointer;
	}

	.photo-grid {
		column-count: 2;
		column-gap: 0.75rem;
	}

	.photo {
		display: block;
		width: 100%;
		margin-bottom: 0.75rem;
		border-radius: 18px;
		overflow: hidden;
		background: transparent;
	}

	.photo img {
		width: 100%;
		border-radius: 18px;
	}

	.comment-head {
		display: flex;
		justify-content: space-between;
		align-items: baseline;
		gap: 1rem;
		margin-bottom: 0.45rem;
	}

	.comment-head span,
	.comment p {
		color: var(--text-secondary);
	}

	.composer {
		position: sticky;
		bottom: 5.7rem;
		display: grid;
		grid-template-columns: minmax(0, 1fr) auto;
		gap: 0.6rem;
		padding: 0.7rem;
		border-radius: var(--radius-xl);
	}

	.composer input {
		min-height: 46px;
		padding: 0 0.85rem;
		border-radius: var(--radius-lg);
		border: 1px solid rgba(255, 255, 255, 0.08);
		background: rgba(255, 255, 255, 0.04);
	}

	.composer button {
		min-width: 86px;
		border-radius: var(--radius-lg);
		background: var(--accent-grad);
		color: white;
		font-weight: 800;
	}

	.empty,
	.error {
		padding: 1rem;
		border-radius: var(--radius-xl);
		text-align: center;
		color: var(--text-secondary);
	}

	.viewer {
		position: fixed;
		inset: 0;
		z-index: 40;
		display: grid;
		place-items: center;
		padding: 1.5rem;
	}

	.viewer-backdrop {
		position: absolute;
		inset: 0;
		background: rgba(0, 0, 0, 0.92);
	}

	.viewer img {
		position: relative;
		max-height: calc(100vh - 3rem);
		border-radius: 24px;
	}

	.viewer-close {
		position: absolute;
		top: 1rem;
		right: 1rem;
		width: 44px;
		height: 44px;
		border-radius: 50%;
		background: rgba(255, 255, 255, 0.1);
		color: white;
		font-size: 1.7rem;
	}

	@media (max-width: 640px) {
		.photo-grid {
			column-count: 1;
		}
	}
</style>
