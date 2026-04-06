<script lang="ts">
	import { goto } from '$app/navigation';
	import { apiFetch, presignTripPhoto, uploadSignedPhoto } from '$lib/api';
	import { getDraft, clearDraft } from '$lib/tripDraft';
	import { onMount, onDestroy } from 'svelte';
	import type { TripCardData, User } from '$lib/types';
	import { getUserName, getUserInitials } from '$lib/format';
	import { setupMainButton, hideMainButton, setMainButtonState, setupBackButton } from '$lib/telegram';

	let description = $state('');
	let coverPhotoPreview = $state('');
	let coverPhotoFile: File | null = null;
	let fileInput: HTMLInputElement | null = null;
	let loading = $state(false);
	let error = $state('');
	let createdTripId = $state('');
	let invitesSent = $state(false);

	let selectedFriends = $state<User[]>([]);
	let friendSearch = $state('');
	let searchResults = $state<User[]>([]);
	let searchTimeout: ReturnType<typeof setTimeout>;

	const removeFriend = (id: string) => {
		selectedFriends = selectedFriends.filter((f) => f.id !== id);
	};

	const onSearchInput = () => {
		clearTimeout(searchTimeout);
		if (!friendSearch.trim()) {
			searchResults = [];
			return;
		}
		searchTimeout = setTimeout(async () => {
			try {
				const res = await apiFetch<{ items: User[] }>(`/v1/users/search?q=${encodeURIComponent(friendSearch.trim())}`);
				searchResults = (res.items || []).filter(u => !selectedFriends.some(sf => sf.id === u.id));
			} catch (e) {
				console.error('Search failed', e);
			}
		}, 300);
	};

	const addFriend = (user: User) => {
		selectedFriends = [...selectedFriends, user];
		friendSearch = '';
		searchResults = [];
	};

	$effect(() => {
		const buttonLabel = loading
			? 'Creating...'
			: createdTripId && coverPhotoFile
				? 'Retry Cover Upload'
				: 'Create Trip';
		setupMainButton(buttonLabel, submit, true, !loading);
	});

	onMount(() => {
		setupBackButton(() => history.back());
	});

	onDestroy(() => {
		hideMainButton();
		if (coverPhotoPreview) URL.revokeObjectURL(coverPhotoPreview);
	});

	const pickPhoto = () => fileInput?.click();

	const onFile = (event: Event) => {
		const target = event.currentTarget as HTMLInputElement;
		const file = target.files?.[0];
		if (!file) return;
		if (coverPhotoPreview) URL.revokeObjectURL(coverPhotoPreview);
		coverPhotoFile = file;
		coverPhotoPreview = URL.createObjectURL(file);
	};

	const submit = async () => {
		const draft = getDraft();
		if (!draft || !draft.title) {
			error = 'Trip Draft is missing title! Please start over.';
			return;
		}
		if (loading) return;
		loading = true;
		setMainButtonState(true);
		error = '';

		try {
			const tripId =
				createdTripId ||
				(
					await apiFetch<{ id: string }>('/v1/trips', {
						method: 'POST',
						body: JSON.stringify({
							title: draft.title,
							city: draft.destination,
							country_code: draft.country_code,
							description: description || undefined,
							start_date: draft.start_date || null,
							end_date: draft.end_date || null,
							status: 'planned',
							visibility: draft.visibility || 'public'
						})
					})
				).id;
			createdTripId = tripId;

			// 2. Upload Photo
			if (coverPhotoFile) {
				try {
					const presign = await presignTripPhoto(tripId, coverPhotoFile.name, coverPhotoFile.type);
					await uploadSignedPhoto(presign.signed_url, presign.token, coverPhotoFile);
					const updatedTrip = await apiFetch<TripCardData>(`/v1/trips/${tripId}`, {
						method: 'PATCH',
						body: JSON.stringify({ cover_photo_url: presign.path || presign.public_url })
					});
					if (!updatedTrip.cover_photo_url) {
						error =
							'Trip created, but cover photo was not attached yet. Tap the button again to retry the cover upload.';
						return;
					}
				} catch {
					error =
						'Trip created, but cover photo failed to save. Tap the button again to retry the cover upload.';
					return;
				}
			}

			// 2.5 Send Invites
			if (!invitesSent) {
				const failedInvites: string[] = [];
				for (const friend of selectedFriends) {
					try {
						await apiFetch(`/v1/trips/${tripId}/invite`, {
							method: 'POST',
							body: JSON.stringify({ user_id: friend.id })
						});
					} catch {
						failedInvites.push(friend.first_name || friend.username || 'someone');
					}
				}
				invitesSent = true;
				if (failedInvites.length > 0) {
					error = `Trip created! Some invites failed to send: ${failedInvites.join(', ')}`;
					await new Promise((r) => setTimeout(r, 1500));
				}
			}

			// 3. Clear draft and navigate
			clearDraft();
			createdTripId = '';
			invitesSent = false;
			goto(`/trips/${tripId}`);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to create trip';
		} finally {
			loading = false;
			setMainButtonState(false);
		}
	};
</script>

<svelte:head>
	<title>TripListik — Create Trip: Travelers</title>
	<meta name="description" content="Choose your travel companions and add a description." />
</svelte:head>

<div class="page">
	<!-- Sticky Header -->
	<header class="top-bar">
		<button class="back-btn" aria-label="Go back" onclick={() => history.back()}>
			<span class="material-symbols-outlined">arrow_back</span>
		</button>
		<span class="header-title">Create Trip</span>
		<div class="header-spacer"></div>
	</header>

	<!-- Progress -->
	<div class="progress-track">
		<div class="progress-fill"></div>
	</div>

	<main class="content">
		<div class="progress-meta">
			<span class="step-label">Step 3 of 3</span>
			<span class="pct-label">100% Complete</span>
		</div>

		<h1>Who's coming along?</h1>

		<!-- Friend Chips -->
		{#if selectedFriends.length > 0}
			<div class="chips-wrap">
				{#each selectedFriends as friend (friend.id)}
					<div class="chip">
						{#if friend.photo_url}
							<img class="chip-avatar" src={friend.photo_url} alt={getUserName(friend)} loading="lazy" />
						{:else}
							<div class="chip-avatar fallback">{getUserInitials(friend.first_name, friend.last_name, friend.username)}</div>
						{/if}
						<span class="chip-name">{getUserName(friend)}</span>
						<button
							class="chip-remove"
							aria-label="Remove {getUserName(friend)}"
							onclick={() => removeFriend(friend.id)}
						>
							<span class="material-symbols-outlined">close</span>
						</button>
					</div>
				{/each}
			</div>
		{/if}

		<!-- Friend Search -->
		<div class="friend-search">
			<input
				type="text"
				class="search-input"
				placeholder="Search friends..."
				bind:value={friendSearch}
				oninput={onSearchInput}
				autocomplete="off"
			/>
			<span class="material-symbols-outlined search-icon">search</span>
			
			{#if searchResults.length > 0}
				<div class="search-dropdown">
					{#each searchResults as user (user.id)}
						<button class="search-item" onclick={() => addFriend(user)}>
							{#if user.photo_url}
								<img src={user.photo_url} alt={getUserName(user)} loading="lazy" />
							{:else}
								<div class="fallback-avatar">{getUserInitials(user.first_name, user.last_name, user.username)}</div>
							{/if}
							<span>{getUserName(user)}</span>
						</button>
					{/each}
				</div>
			{/if}
		</div>

		<!-- Cover Photo -->
		<div class="cover-section">
			<label class="section-label" for="cover-file">Cover Photo</label>
			<button type="button" class="upload-zone" onclick={pickPhoto}>
				{#if coverPhotoPreview}
					<img class="preview-img" src={coverPhotoPreview} alt="Cover preview" />
				{:else}
					<div class="upload-icon-circle">
						<span class="material-symbols-outlined">add_a_photo</span>
					</div>
					<p class="upload-hint">Add a cover image</p>
				{/if}
			</button>
			<input
				id="cover-file"
				class="file-input"
				type="file"
				accept="image/*"
				bind:this={fileInput}
				onchange={onFile}
			/>
		</div>

		<!-- Description -->
		<div class="desc-section">
			<label class="section-label" for="trip-desc">Trip Description <span class="optional">(Optional)</span></label>
			<textarea
				id="trip-desc"
				class="desc-area"
				rows="3"
				bind:value={description}
				placeholder="Tell everyone what the plan is..."
			></textarea>
		</div>

		{#if error}
			<p class="error-msg">{error}</p>
		{/if}

		<!-- Create Trip Button -->
		<button
			class="create-btn"
			onclick={submit}
			disabled={loading}
		>
			{#if loading}
				<span class="material-symbols-outlined spin-icon">refresh</span>
				Creating...
			{:else}
				Create Trip
			{/if}
		</button>
	</main>
</div>

<style>
.page {
	min-height: 100dvh;
	background: #161c18;
	color: #f1f5f9;
}

/* ── Sticky Header ─────────────────────────────────────── */
.top-bar {
	position: sticky;
	top: 0;
	z-index: 50;
	display: flex;
	align-items: center;
	height: 56px;
	padding: 0 16px;
	background: #161c18;
	border-bottom: 1px solid rgba(77, 157, 109, 0.1);
}

.back-btn {
	width: 44px;
	height: 44px;
	display: flex;
	align-items: center;
	justify-content: center;
	color: white;
	background: none;
	border: none;
	cursor: pointer;
	flex-shrink: 0;
}

.back-btn .material-symbols-outlined {
	font-size: 22px;
}

.header-title {
	flex: 1;
	text-align: center;
	font-size: 18px;
	font-weight: 700;
	color: white;
}

.header-spacer {
	width: 44px;
}

/* ── Progress ────────────────────────────────────────────── */
.progress-track {
	width: 100%;
	height: 8px;
	background: rgba(77, 157, 109, 0.15);
	border-radius: 9999px;
}

.progress-fill {
	width: 100%;
	height: 100%;
	background: #4d9d6d;
	border-radius: 9999px;
	transition: width 500ms cubic-bezier(0.4, 0, 0.2, 1);
}

/* ── Content ─────────────────────────────────────────────── */
.content {
	padding: 16px 16px 96px;
	max-width: 480px;
	margin: 0 auto;
}

.progress-meta {
	display: flex;
	justify-content: space-between;
	align-items: center;
	padding-top: 8px;
	margin-bottom: 16px;
}

.step-label {
	font-size: 13px;
	color: #4d9d6d;
	font-weight: 600;
}

.pct-label {
	font-size: 13px;
	color: #94a3b8;
}

h1 {
	font-size: 24px;
	font-weight: 700;
	color: white;
	margin-bottom: 24px;
	line-height: 1.2;
}

/* ── Friend Chips ────────────────────────────────────────── */
.chips-wrap {
	display: flex;
	flex-wrap: wrap;
	gap: 8px;
	margin-bottom: 16px;
}

.chip {
	display: flex;
	align-items: center;
	gap: 6px;
	background: rgba(77, 157, 109, 0.15);
	border: 1px solid rgba(77, 157, 109, 0.25);
	border-radius: 9999px;
	padding: 4px 8px 4px 4px;
}

.chip-avatar {
	width: 28px;
	height: 28px;
	border-radius: 9999px;
	object-fit: cover;
	flex-shrink: 0;
}

.chip-name {
	font-size: 14px;
	font-weight: 500;
	color: white;
}

.chip-remove {
	width: 20px;
	height: 20px;
	display: flex;
	align-items: center;
	justify-content: center;
	background: none;
	border: none;
	color: #94a3b8;
	cursor: pointer;
	padding: 0;
}

.chip-remove .material-symbols-outlined {
	font-size: 16px;
}

/* ── Friend Search ───────────────────────────────────────── */
.friend-search {
	position: relative;
	margin-bottom: 32px;
}

.search-input {
	width: 100%;
	background: transparent;
	border: none;
	border-bottom: 1px solid rgba(255, 255, 255, 0.12);
	border-radius: 0;
	font-size: 18px;
	color: white;
	padding: 8px 36px 8px 0;
	outline: none;
	transition: border-bottom-color 0.2s ease;
}

.search-input:focus {
	border-bottom-color: #4d9d6d;
}

.search-icon {
	position: absolute;
	right: 0;
	top: 50%;
	transform: translateY(-50%);
	font-size: 24px;
	color: #94a3b8;
	pointer-events: none;
}

/* ── Cover Photo ─────────────────────────────────────────── */
.cover-section {
	margin-bottom: 32px;
}

.section-label {
	display: block;
	font-size: 14px;
	color: #94a3b8;
	font-weight: 600;
	margin-bottom: 12px;
}

.optional {
	font-weight: 400;
	color: #64748b;
}

.upload-zone {
	width: 100%;
	aspect-ratio: 16 / 9;
	border-radius: 12px;
	border: 2px dashed rgba(255, 255, 255, 0.1);
	background: rgba(255, 255, 255, 0.03);
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	gap: 8px;
	cursor: pointer;
	overflow: hidden;
	transition: border-color 0.2s ease, background 0.2s ease;
}

.upload-zone:hover {
	border-color: rgba(77, 157, 109, 0.4);
}

.upload-icon-circle {
	width: 48px;
	height: 48px;
	border-radius: 9999px;
	background: rgba(77, 157, 109, 0.1);
	display: flex;
	align-items: center;
	justify-content: center;
}

.upload-icon-circle .material-symbols-outlined {
	font-size: 24px;
	color: #4d9d6d;
}

.upload-hint {
	font-size: 12px;
	color: #64748b;
	font-weight: 500;
}

.preview-img {
	width: 100%;
	height: 100%;
	object-fit: cover;
}

.file-input {
	display: none;
}

/* ── Description ─────────────────────────────────────────── */
.desc-section {
	margin-bottom: 32px;
}

.desc-area {
	width: 100%;
	background: transparent;
	border: none;
	border-bottom: 2px solid rgba(255, 255, 255, 0.1);
	border-radius: 0;
	font-size: 16px;
	color: white;
	padding: 8px 0;
	resize: none;
	font-family: inherit;
	outline: none;
	transition: border-bottom-color 0.2s ease;
}

.desc-area::placeholder {
	color: rgba(255, 255, 255, 0.25);
}

.desc-area:focus {
	border-bottom-color: #4d9d6d;
}

/* ── Error ───────────────────────────────────────────────── */
.error-msg {
	color: #f87171;
	font-size: 14px;
	text-align: center;
	padding: 12px;
	background: rgba(248, 113, 113, 0.08);
	border-radius: 8px;
	margin-bottom: 16px;
}

/* ── Create Button ───────────────────────────────────────── */
.create-btn {
	width: 100%;
	display: flex;
	align-items: center;
	justify-content: center;
	gap: 8px;
	padding: 16px;
	border-radius: 12px;
	background: #4d9d6d;
	color: white;
	font-size: 18px;
	font-weight: 700;
	border: none;
	cursor: pointer;
	box-shadow: 0 0 24px rgba(77, 157, 109, 0.35);
	transition: opacity 0.15s ease, transform 0.1s ease;
}

.create-btn:hover:not(:disabled) {
	opacity: 0.92;
}

.create-btn:active:not(:disabled) {
	transform: scale(0.98);
}

.create-btn:disabled {
	opacity: 0.6;
	cursor: not-allowed;
}

@keyframes spin {
	from { transform: rotate(0deg); }
	to { transform: rotate(360deg); }
}

.spin-icon {
	animation: spin 1s linear infinite;
	font-size: 20px;
}

/* ── Search Dropdown ─────────────────────────────────── */
.search-dropdown {
	position: absolute;
	top: calc(100% + 4px);
	left: 0;
	right: 0;
	background: #1e2820;
	border: 1px solid rgba(77, 157, 109, 0.2);
	border-radius: 12px;
	overflow: hidden;
	z-index: 10;
	box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
}

.search-item {
	display: flex;
	align-items: center;
	gap: 10px;
	width: 100%;
	padding: 10px 14px;
	background: none;
	border: none;
	color: white;
	font-size: 14px;
	cursor: pointer;
	text-align: left;
	transition: background 0.15s ease;
}

.search-item:hover {
	background: rgba(77, 157, 109, 0.1);
}

.search-item img {
	width: 32px;
	height: 32px;
	border-radius: 50%;
	object-fit: cover;
	flex-shrink: 0;
}

.fallback-avatar {
	width: 32px;
	height: 32px;
	border-radius: 50%;
	background: rgba(77, 157, 109, 0.2);
	color: #4d9d6d;
	font-size: 13px;
	font-weight: 700;
	display: flex;
	align-items: center;
	justify-content: center;
	flex-shrink: 0;
}
</style>
