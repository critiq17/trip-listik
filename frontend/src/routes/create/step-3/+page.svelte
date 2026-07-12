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
	<header class="top-bar">
		<button class="back-btn" aria-label="Go back" onclick={() => history.back()}>
			<span class="material-symbols-outlined">arrow_back</span>
		</button>
		<span class="header-title">New Trip</span>
		<div class="header-spacer"></div>
	</header>

	<div class="progress-track">
		<div class="progress-fill" style="width: 100%"></div>
	</div>

	<main class="content">
		<p class="step-label">Step 3 of 3</p>
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
			<span class="material-symbols-outlined search-icon">search</span>
			<input
				type="text"
				class="search-input"
				placeholder="Search friends..."
				bind:value={friendSearch}
				oninput={onSearchInput}
				autocomplete="off"
			/>

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
		<div class="section">
			<label class="section-label" for="cover-file">Cover photo</label>
			<button type="button" class="upload-zone" onclick={pickPhoto}>
				{#if coverPhotoPreview}
					<img class="preview-img" src={coverPhotoPreview} alt="Cover preview" />
				{:else}
					<span class="material-symbols-outlined upload-icon">add_a_photo</span>
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
		<div class="section">
			<label class="section-label" for="trip-desc">Description <span class="optional">(optional)</span></label>
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
			{loading ? 'Creating...' : 'Create Trip'}
		</button>
	</main>
</div>

<style>
.page {
	min-height: 100dvh;
	background: var(--bg);
	color: var(--text);
}

/* ── Header ─────────────────────────────────────────────── */
.top-bar {
	position: sticky;
	top: 0;
	z-index: 50;
	display: flex;
	align-items: center;
	height: 52px;
	padding: 0 8px;
	background: var(--bg);
	border-bottom: 1px solid var(--border);
}

.back-btn {
	width: 44px;
	height: 44px;
	display: flex;
	align-items: center;
	justify-content: center;
	color: var(--text);
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
	font-size: 16px;
	font-weight: 600;
	color: var(--text);
}

.header-spacer {
	width: 44px;
}

/* ── Progress ────────────────────────────────────────────── */
.progress-track {
	width: 100%;
	height: 3px;
	background: var(--border);
}

.progress-fill {
	height: 100%;
	background: var(--green);
	transition: width 0.2s ease;
}

/* ── Content ─────────────────────────────────────────────── */
.content {
	padding: 24px 16px 96px;
	max-width: 480px;
	margin: 0 auto;
}

.step-label {
	font-size: 13px;
	color: var(--text-sub);
	font-weight: 500;
	margin-bottom: 8px;
}

h1 {
	font-size: 24px;
	font-weight: 700;
	color: var(--text);
	margin-bottom: 24px;
	line-height: 1.2;
}

/* ── Friend Chips ────────────────────────────────────────── */
.chips-wrap {
	display: flex;
	flex-wrap: wrap;
	gap: 8px;
	margin-bottom: 12px;
}

.chip {
	display: flex;
	align-items: center;
	gap: 6px;
	background: var(--green-soft);
	border-radius: var(--radius-pill);
	padding: 4px 8px 4px 4px;
}

.chip-avatar {
	width: 28px;
	height: 28px;
	border-radius: var(--radius-pill);
	object-fit: cover;
	flex-shrink: 0;
}

.chip-avatar.fallback {
	display: flex;
	align-items: center;
	justify-content: center;
	background: var(--bg-card);
	color: var(--green);
	font-size: 11px;
	font-weight: 600;
}

.chip-name {
	font-size: 14px;
	font-weight: 500;
	color: var(--text);
}

.chip-remove {
	width: 20px;
	height: 20px;
	display: flex;
	align-items: center;
	justify-content: center;
	background: none;
	border: none;
	color: var(--text-sub);
	cursor: pointer;
	padding: 0;
}

.chip-remove .material-symbols-outlined {
	font-size: 16px;
}

/* ── Friend Search ───────────────────────────────────────── */
.friend-search {
	position: relative;
	margin-bottom: 24px;
}

.search-input {
	width: 100%;
	background: var(--bg-input);
	border: 1px solid var(--border);
	border-radius: var(--radius-input);
	font-size: 15px;
	color: var(--text);
	padding: 12px 14px 12px 40px;
	outline: none;
	transition: border-color 0.15s ease;
}

.search-input:focus {
	border-color: var(--green);
}

.search-icon {
	position: absolute;
	left: 12px;
	top: 22px;
	transform: translateY(-50%);
	font-size: 20px;
	color: var(--text-muted);
	pointer-events: none;
}

.search-dropdown {
	position: absolute;
	top: calc(100% + 4px);
	left: 0;
	right: 0;
	background: var(--bg-card);
	border: 1px solid var(--border);
	border-radius: var(--radius-card);
	overflow: hidden;
	z-index: 10;
	box-shadow: 0 4px 16px rgba(17, 20, 24, 0.08);
}

.search-item {
	display: flex;
	align-items: center;
	gap: 10px;
	width: 100%;
	padding: 10px 14px;
	background: none;
	border: none;
	color: var(--text);
	font-size: 14px;
	cursor: pointer;
	text-align: left;
}

.search-item:hover {
	background: var(--bg-subtle);
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
	background: var(--green-soft);
	color: var(--green);
	font-size: 13px;
	font-weight: 600;
	display: flex;
	align-items: center;
	justify-content: center;
	flex-shrink: 0;
}

/* ── Sections ────────────────────────────────────────────── */
.section {
	margin-bottom: 24px;
}

.section-label {
	display: block;
	font-size: 13px;
	color: var(--text-sub);
	font-weight: 500;
	margin-bottom: 8px;
}

.optional {
	font-weight: 400;
	color: var(--text-muted);
}

.upload-zone {
	width: 100%;
	aspect-ratio: 16 / 9;
	border-radius: var(--radius-card);
	border: 1px dashed var(--border);
	background: var(--bg-subtle);
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	gap: 8px;
	cursor: pointer;
	overflow: hidden;
	transition: border-color 0.15s ease;
}

.upload-zone:hover {
	border-color: var(--green);
}

.upload-icon {
	font-size: 28px;
	color: var(--green);
}

.upload-hint {
	font-size: 13px;
	color: var(--text-sub);
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
.desc-area {
	width: 100%;
	background: var(--bg-input);
	border: 1px solid var(--border);
	border-radius: var(--radius-input);
	font-size: 15px;
	color: var(--text);
	padding: 12px 14px;
	resize: none;
	font-family: inherit;
	outline: none;
	transition: border-color 0.15s ease;
}

.desc-area:focus {
	border-color: var(--green);
}

/* ── Error ───────────────────────────────────────────────── */
.error-msg {
	color: var(--danger);
	font-size: 14px;
	text-align: center;
	padding: 12px;
	background: var(--danger-soft);
	border-radius: var(--radius-input);
	margin-bottom: 16px;
}

/* ── Create Button ───────────────────────────────────────── */
.create-btn {
	width: 100%;
	padding: 14px;
	border-radius: var(--radius-input);
	background: var(--green);
	color: #fff;
	font-size: 16px;
	font-weight: 600;
	border: none;
	cursor: pointer;
}

.create-btn:disabled {
	opacity: 0.6;
	cursor: not-allowed;
}
</style>
