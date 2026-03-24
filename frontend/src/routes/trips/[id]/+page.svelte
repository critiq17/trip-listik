<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import {
		apiFetch,
		deleteTrip,
		getMe,
		getPublicPhotoURL,
		presignTripPhoto,
		uploadSignedPhoto
	} from '$lib/api';
	import { resolvePhotoUrl } from '$lib/photos';
	import { connectTripStream } from '$lib/realtime';
	import InviteModal from '$lib/components/InviteModal.svelte';
	import { scalePress } from '$lib/actions/animate';
	import { hapticImpact, hapticNotification, setupBackButton, hideBackButton } from '$lib/telegram';
	import {
		formatDateRange,
		formatLongDate,
		getStatusLabel,
		getTripLocation,
		getUserInitials,
		getUserName
	} from '$lib/format';
	import type { InviteItem, JoinRequest, Member, Photo, TripCardData, VoteItemData } from '$lib/types';

	type TripDetailResponse = {
		trip: TripCardData;
		member_count: number;
		comment_count: number;
		photo_count: number;
		viewer_is_member: boolean;
	};

	let trip = $state<TripCardData | null>(null);
	let loading = $state(true);
	let error = $state('');
	let activeTab = $state('Details');
	let members = $state<Member[]>([]);
	let photos = $state<Photo[]>([]);
	let memberCount = $state(0);
	let commentCount = $state(0);
	let photoCount = $state(0);
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

	// Members bottom sheet
	let membersSheetOpen = $state(false);
	let sheetMembers = $state<Member[]>([]);
	let sheetMembersLoading = $state(false);

	// Vote items
	let voteItems = $state<VoteItemData[]>([]);
	let voteItemsLoading = $state(false);
	let addOptionOpen = $state(false);
	// Add option form
	let newOptionCategory = $state('ACTIVITY');
	let newOptionTitle = $state('');
	let newOptionDesc = $state('');
	let newOptionSubmitting = $state(false);
	let newOptionError = $state('');

	const tabs = ['Details', 'Photos'];
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

	let deletingTrip = $state(false);

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

	const canInvite = () =>
		!!trip &&
		!!meId &&
		(trip.owner_id === meId || (trip.visibility === 'group' && viewerIsMember));

	async function openMembersSheet() {
		membersSheetOpen = true;
		if (sheetMembers.length > 0) return;
		sheetMembersLoading = true;
		try {
			const data = await apiFetch<{ items: Member[] }>(
				`/v1/trips/${trip!.id}/members?limit=50`
			);
			sheetMembers = data.items ?? [];
		} catch {
			// non-critical — sheet stays open showing empty state
		} finally {
			sheetMembersLoading = false;
		}
	}

	async function loadVoteItems(id: string) {
		if (!viewerIsMember) return;
		voteItemsLoading = true;
		try {
			const data = await apiFetch<{ items: VoteItemData[] }>(`/v1/trips/${id}/vote-items`);
			voteItems = data.items ?? [];
		} catch {
			// non-critical — silently ignore
		} finally {
			voteItemsLoading = false;
		}
	}

	async function handleVote(item: VoteItemData) {
		if (!trip) return;
		// Optimistic update
		const prev = { vote_count: item.vote_count, has_voted: item.has_voted };
		if (item.has_voted) {
			item.vote_count = Math.max(0, item.vote_count - 1);
			item.has_voted = false;
		} else {
			item.vote_count += 1;
			item.has_voted = true;
		}
		voteItems = [...voteItems]; // trigger reactivity
		try {
			const method = prev.has_voted ? 'DELETE' : 'POST';
			const res = await apiFetch<{ item_id: string; vote_count: number }>(
				`/v1/trips/${trip.id}/vote-items/${item.id}/vote`,
				{ method }
			);
			// Sync server count
			voteItems = voteItems.map(v => v.id === item.id ? { ...v, vote_count: res.vote_count } : v);
		} catch {
			// Roll back
			voteItems = voteItems.map(v => v.id === item.id ? { ...v, ...prev } : v);
		}
	}

	async function submitOption() {
		if (!trip || !newOptionTitle.trim()) return;
		newOptionSubmitting = true;
		newOptionError = '';
		try {
			const item = await apiFetch<VoteItemData>(`/v1/trips/${trip.id}/vote-items`, {
				method: 'POST',
				body: JSON.stringify({
					category: newOptionCategory,
					title: newOptionTitle.trim(),
					description: newOptionDesc.trim() || undefined
				})
			});
			voteItems = [...voteItems, { ...item, vote_count: 0, has_voted: false }];
			newOptionTitle = '';
			newOptionDesc = '';
			newOptionCategory = 'ACTIVITY';
			addOptionOpen = false;
		} catch (err) {
			newOptionError = err instanceof Error ? err.message : 'Failed to add option';
		} finally {
			newOptionSubmitting = false;
		}
	}

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
		stream.addEventListener('photo_created', (event) => {
			const data = JSON.parse((event as MessageEvent).data) as Photo;
			const exists = photos.find((item) => item.id === data.id);
			photos = [data, ...photos.filter((item) => item.id !== data.id)];
			if (!exists) {
				photoCount += 1;
			}
		});

		stream.addEventListener('vote_item_updated', (event) => {
			const data = JSON.parse((event as MessageEvent).data) as { item_id: string; vote_count: number };
			voteItems = voteItems.map(v => v.id === data.item_id ? { ...v, vote_count: data.vote_count } : v);
		});

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
		if (trip && activeTab === 'Details' && !loadedTabs.has('Details') && viewerIsMember) {
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
					<button class="members-chip" onclick={openMembersSheet}>
						{memberCount} members
					</button>
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

					<!-- Vote Items -->
					{#if viewerIsMember}
						<div class="detail-block vote-section">
							<div class="vote-header">
								<h2>Options</h2>
								<button class="add-option-btn" onclick={() => (addOptionOpen = true)}>
									<span class="material-symbols-outlined">add</span>
									Add option
								</button>
							</div>

							{#if voteItemsLoading}
								<div class="vote-skeleton"></div>
							{:else if voteItems.length === 0}
								<p class="muted">No options yet. Add hotels, airlines, or activities to vote on.</p>
							{:else}
								<!-- Group by category -->
								{#each ['HOTEL', 'AIRLINE', 'ACTIVITY', 'OTHER'] as cat}
									{#if voteItems.some(v => v.category === cat)}
										<div class="vote-group">
											<p class="vote-cat-label">{cat}</p>
											{#each voteItems.filter(v => v.category === cat) as item (item.id)}
												<div class="vote-item-row">
													<div class="vote-item-info">
														<p class="vote-item-title">{item.title}</p>
														{#if item.description}
															<p class="vote-item-desc">{item.description}</p>
														{/if}
													</div>
													<button
														class="vote-btn"
														class:voted={item.has_voted}
														onclick={() => handleVote(item)}
														aria-label="Vote for {item.title}"
													>
														<span class="material-symbols-outlined">thumb_up</span>
														<span class="vote-count">{item.vote_count}</span>
													</button>
												</div>
											{/each}
										</div>
									{/if}
								{/each}
							{/if}
						</div>
					{/if}

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

					{#if isOwner()}
						<div class="owner-actions">
							<a class="owner-btn owner-btn--edit" href="/trips/{trip.id}/edit">
								<span class="material-symbols-outlined">edit</span>
								Edit trip
							</a>
							<button class="owner-btn owner-btn--delete" onclick={handleDeleteTrip} disabled={deletingTrip}>
								<span class="material-symbols-outlined">delete</span>
								{deletingTrip ? 'Deleting…' : 'Delete trip'}
							</button>
						</div>
					{:else}
						<button class="cta-button" onclick={joinTrip} disabled={joinLoading || !!joinStatus} use:scalePress>
							{joinLoading ? 'Joining...' : (joinStatusLabel[joinStatus] ?? joinStatus) || 'Join Trip'}
						</button>
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
			{/if}
		</div>
	{/if}
</section>

<!-- Members bottom sheet -->
{#if membersSheetOpen}
	<button
		class="sheet-backdrop"
		onclick={() => (membersSheetOpen = false)}
		aria-label="Close members panel"
	></button>
	<div class="sheet" role="dialog" aria-label="Trip members">
		<div class="sheet-handle"></div>
		<div class="sheet-header">
			<h2>Members</h2>
			<button class="sheet-close" onclick={() => (membersSheetOpen = false)} aria-label="Close">×</button>
		</div>
		{#if sheetMembersLoading}
			<div class="sheet-state">Loading…</div>
		{:else if sheetMembers.length === 0}
			<div class="sheet-state">No members yet</div>
		{:else}
			<div class="sheet-list">
				{#each sheetMembers as member}
					<div class="sheet-member">
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
		{/if}
	</div>
{/if}

<!-- Add Option Bottom Sheet -->
{#if addOptionOpen}
	<div class="sheet-backdrop" onclick={() => (addOptionOpen = false)} role="presentation"></div>
	<div class="bottom-sheet" role="dialog" aria-label="Add option">
		<div class="sheet-handle"></div>
		<div class="sheet-header">
			<h3>Add Option</h3>
			<button class="sheet-close" onclick={() => (addOptionOpen = false)} aria-label="Close">×</button>
		</div>

		<div class="sheet-body">
			<div class="sheet-field">
				<label class="sheet-label" for="option-category">Category</label>
				<select id="option-category" class="sheet-select" bind:value={newOptionCategory}>
					<option value="HOTEL">Hotel</option>
					<option value="AIRLINE">Airline</option>
					<option value="ACTIVITY">Activity</option>
					<option value="OTHER">Other</option>
				</select>
			</div>

			<div class="sheet-field">
				<label class="sheet-label" for="option-title">Title *</label>
				<input
					id="option-title"
					class="sheet-input"
					type="text"
					placeholder="e.g. Marriott Rome"
					bind:value={newOptionTitle}
				/>
			</div>

			<div class="sheet-field">
				<label class="sheet-label" for="option-desc">Description (optional)</label>
				<textarea
					id="option-desc"
					class="sheet-textarea"
					rows="2"
					placeholder="Add details..."
					bind:value={newOptionDesc}
				></textarea>
			</div>

			{#if newOptionError}
				<p class="sheet-error">{newOptionError}</p>
			{/if}

			<button
				class="sheet-submit"
				onclick={submitOption}
				disabled={!newOptionTitle.trim() || newOptionSubmitting}
			>
				{newOptionSubmitting ? 'Adding...' : 'Add Option'}
			</button>
		</div>
	</div>
{/if}

<!-- Escape key closes photo viewer -->
<svelte:window onkeydown={(e) => { if (e.key === 'Escape') { viewerPhoto = null; membersSheetOpen = false; addOptionOpen = false; } }} />

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

	/* Grid gallery thumbnails need relative positioning + intrinsic size */
	.grid-img {
		position: relative;
		width: 100%;
		height: auto;
		opacity: 0;
		transition: opacity 0.3s ease;
		aspect-ratio: 1;
		object-fit: cover;
	}
	:global(.grid-img.loaded) { opacity: 1; }
	:global(.grid-img.error)  { display: none; }

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

	.members-chip {
		padding: 0.45rem 0.7rem;
		border-radius: var(--radius-pill);
		background: rgba(255, 255, 255, 0.08);
		backdrop-filter: blur(16px);
		font-size: 0.8rem;
		color: white;
		border: none;
		cursor: pointer;
		transition: background 0.15s ease;
	}

	.members-chip:hover {
		background: rgba(255, 255, 255, 0.14);
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
	.photo-panel {
		display: grid;
		gap: 1rem;
		padding: 1.15rem;
		border-radius: var(--radius-2xl);
	}

	.stats-grid {
		display: grid;
		grid-template-columns: repeat(3, 1fr);
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
	.section-head h2 {
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

	.owner-actions {
		display: flex;
		gap: 12px;
		margin-top: 8px;
	}

	.owner-btn {
		flex: 1;
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 6px;
		padding: 13px 16px;
		border-radius: 12px;
		font-size: 14px;
		font-weight: 600;
		cursor: pointer;
		border: none;
		text-decoration: none;
		transition: opacity 0.15s ease;
		-webkit-tap-highlight-color: transparent;
	}

	.owner-btn:disabled {
		opacity: 0.5;
		cursor: default;
	}

	.owner-btn--edit {
		background: rgba(77, 157, 109, 0.15);
		color: #4d9d6d;
	}

	.owner-btn--delete {
		background: rgba(239, 68, 68, 0.1);
		color: #ef4444;
	}

	.owner-btn .material-symbols-outlined {
		font-size: 18px;
	}

	/* ── Members bottom sheet ─────────────────────────────── */
	.sheet-backdrop {
		position: fixed;
		inset: 0;
		z-index: 30;
		background: rgba(0, 0, 0, 0.6);
		backdrop-filter: blur(4px);
		border: none;
		cursor: pointer;
	}

	.sheet {
		position: fixed;
		bottom: 0;
		left: 0;
		right: 0;
		z-index: 31;
		background: #1a2420;
		border-radius: 20px 20px 0 0;
		padding: 12px 20px 40px;
		max-height: 70vh;
		overflow-y: auto;
	}

	.sheet-handle {
		width: 36px;
		height: 4px;
		border-radius: 9999px;
		background: rgba(255, 255, 255, 0.12);
		margin: 0 auto 16px;
	}

	.sheet-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: 16px;
	}

	.sheet-header h2 {
		font-size: 16px;
		font-weight: 700;
		color: white;
	}

	.sheet-close {
		width: 32px;
		height: 32px;
		border-radius: 50%;
		background: rgba(255, 255, 255, 0.08);
		color: white;
		font-size: 1.4rem;
		display: grid;
		place-items: center;
		border: none;
		cursor: pointer;
	}

	.sheet-state {
		text-align: center;
		color: var(--text-secondary);
		font-size: 14px;
		padding: 24px 0;
	}

	.sheet-list {
		display: grid;
		gap: 0.75rem;
	}

	.sheet-member {
		display: flex;
		align-items: center;
		gap: 0.8rem;
		padding: 0.8rem 0.9rem;
		border-radius: var(--radius-xl);
		background: rgba(255, 255, 255, 0.04);
		border: 1px solid rgba(255, 255, 255, 0.06);
	}

	.avatar {
		width: 48px;
		height: 48px;
		flex-shrink: 0;
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
		font-size: 1rem;
		font-weight: 700;
		border-radius: 50%;
	}

	.member-copy {
		flex: 1;
		min-width: 0;
	}

	.member-copy strong {
		display: block;
		color: white;
		font-size: 0.9rem;
	}

	.member-copy small {
		color: var(--text-secondary);
		font-size: 0.8rem;
	}

	.role {
		padding: 0.35rem 0.6rem;
		border-radius: var(--radius-pill);
		background: rgba(77, 157, 109, 0.16);
		color: var(--accent-strong);
		font-size: 0.72rem;
		font-weight: 700;
	}

	.load-more {
		padding: 0.7rem 1rem;
		border-radius: var(--radius-pill);
		background: rgba(255, 255, 255, 0.06);
		color: var(--text-secondary);
		font-size: 0.85rem;
		font-weight: 600;
		border: none;
		cursor: pointer;
		width: 100%;
	}

	.muted {
		color: var(--text-secondary);
	}

	/* ── Vote Items ─────────────────────────────────────────── */
	.vote-section { gap: 0; }

	.vote-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: 12px;
	}

	.add-option-btn {
		display: flex;
		align-items: center;
		gap: 4px;
		background: rgba(77,157,109,0.12);
		border: 1px solid rgba(77,157,109,0.3);
		color: #4d9d6d;
		border-radius: 8px;
		padding: 6px 12px;
		font-size: 13px;
		font-weight: 600;
		cursor: pointer;
	}

	.add-option-btn .material-symbols-outlined { font-size: 16px; }

	.vote-group { margin-bottom: 16px; }

	.vote-cat-label {
		font-size: 10px;
		font-weight: 700;
		letter-spacing: 0.1em;
		text-transform: uppercase;
		color: #64748b;
		margin-bottom: 8px;
	}

	.vote-item-row {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 10px 0;
		border-bottom: 1px solid rgba(255,255,255,0.06);
	}

	.vote-item-info { flex: 1; min-width: 0; }

	.vote-item-title {
		font-size: 15px;
		color: #f1f5f9;
		font-weight: 500;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.vote-item-desc {
		font-size: 12px;
		color: #64748b;
		margin-top: 2px;
	}

	.vote-btn {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 2px;
		background: rgba(255,255,255,0.06);
		border: 1px solid rgba(255,255,255,0.1);
		border-radius: 10px;
		padding: 6px 10px;
		color: #94a3b8;
		cursor: pointer;
		min-width: 48px;
		transition: background 0.15s, color 0.15s;
		flex-shrink: 0;
		margin-left: 12px;
	}

	.vote-btn.voted {
		background: rgba(77,157,109,0.15);
		border-color: rgba(77,157,109,0.4);
		color: #4d9d6d;
	}

	.vote-btn .material-symbols-outlined { font-size: 18px; }

	.vote-count { font-size: 12px; font-weight: 600; }

	.vote-skeleton {
		height: 60px;
		border-radius: 8px;
		background: rgba(255,255,255,0.06);
		animation: pulse 1.5s infinite;
	}

	/* ── Add Option bottom sheet ─── */
	.bottom-sheet {
		position: fixed;
		bottom: 0;
		left: 0;
		right: 0;
		z-index: 31;
		background: #1a2420;
		border-radius: 20px 20px 0 0;
		padding: 12px 20px 40px;
		max-height: 80vh;
		overflow-y: auto;
	}

	.bottom-sheet .sheet-header h3 {
		font-size: 16px;
		font-weight: 700;
		color: white;
	}

	.sheet-body { padding-top: 8px; }

	/* ── Add Option sheet fields ─── */
	.sheet-field { margin-bottom: 16px; }

	.sheet-label {
		display: block;
		font-size: 12px;
		color: #94a3b8;
		font-weight: 500;
		margin-bottom: 6px;
	}

	.sheet-select,
	.sheet-input,
	.sheet-textarea {
		width: 100%;
		background: rgba(255,255,255,0.06);
		border: 1px solid rgba(255,255,255,0.1);
		border-radius: 8px;
		padding: 10px 12px;
		font-size: 15px;
		color: #f1f5f9;
		outline: none;
		box-sizing: border-box;
	}

	.sheet-select option { background: #1e2a22; }

	.sheet-textarea { resize: none; }

	.sheet-error {
		color: #f87171;
		font-size: 13px;
		margin-bottom: 12px;
	}

	.sheet-submit {
		width: 100%;
		padding: 14px;
		background: #4d9d6d;
		border: none;
		border-radius: 10px;
		color: white;
		font-size: 15px;
		font-weight: 700;
		cursor: pointer;
	}

	.sheet-submit:disabled { opacity: 0.4; cursor: not-allowed; }
</style>
