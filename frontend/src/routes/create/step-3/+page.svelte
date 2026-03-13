<script lang="ts">
	import { goto } from '$app/navigation';
	import { apiFetch, presignTripPhoto, uploadSignedPhoto, getPublicPhotoURL } from '$lib/api';
	import { getDraft, clearDraft } from '$lib/tripDraft';
	import { onMount, onDestroy } from 'svelte';
	import { setupMainButton, hideMainButton, setMainButtonState, setupBackButton } from '$lib/telegram';

	let description = $state('');
	let coverPhotoPreview = $state('');
	let coverPhotoFile: File | null = null;
	let fileInput: HTMLInputElement | null = null;
	let loading = $state(false);
	let error = $state('');

	$effect(() => {
		setupMainButton(loading ? 'Creating...' : 'Create Trip', submit, true, !loading);
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
			// 1. Create Trip
			const trip = await apiFetch<{ id: string }>('/v1/trips', {
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
			});

			// 2. Upload Photo
			if (coverPhotoFile) {
				try {
					const presign = await presignTripPhoto(trip.id, coverPhotoFile.name, coverPhotoFile.type);
					await uploadSignedPhoto(presign.signed_url, presign.token, coverPhotoFile);
					const publicUrl = getPublicPhotoURL(presign.path);
					await apiFetch(`/v1/trips/${trip.id}`, {
						method: 'PATCH',
						body: JSON.stringify({ cover_photo_url: publicUrl })
					});
				} catch (uploadErr) {
					console.warn('Cover photo upload failed, continuing without it', uploadErr);
				}
			}

			// 3. Clear draft and go to trip
			clearDraft();
			goto(`/trips/${trip.id}`);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to create trip';
		} finally {
			loading = false;
			setMainButtonState(false);
		}
	};
</script>

<svelte:head>
	<title>Create Trip - Step 3</title>
</svelte:head>

<section class="step-three">
	<div class="progress-bar">
		<div class="progress-fill" style="width: 100%"></div>
	</div>

	<main>
		<p class="eyebrow">Step 3 of 3</p>
		<h1 class="serif-text">Final touches.</h1>

		<div class="form">
			<div class="field">
				<label for="cover-photo">Cover Photo</label>
				<button type="button" class="upload" onclick={pickPhoto}>
					{#if coverPhotoPreview}
						<img src={coverPhotoPreview} alt="Cover preview" />
					{:else}
						<div class="upload-icon">
							<span class="material-symbols-outlined">add_a_photo</span>
						</div>
						<span>Add an inspiring image</span>
					{/if}
				</button>
				<input
					id="cover-photo"
					class="file-input"
					type="file"
					accept="image/*"
					bind:this={fileInput}
					onchange={onFile}
				/>
			</div>

			<div class="field textarea">
				<label for="trip-description">Trip Description (Optional)</label>
				<textarea
					id="trip-description"
					rows="4"
					bind:value={description}
					placeholder="Tell everyone what the plan is..."
				></textarea>
			</div>

			{#if error}
				<p class="error">{error}</p>
			{/if}
		</div>
	</main>
</section>

<style>
	.step-three {
		min-height: 100dvh;
		background: var(--bg);
		color: var(--text);
		padding-bottom: 3rem;
	}

	.progress-bar {
		width: 100%;
		height: 4px;
		background: var(--bg-elevated);
	}

	.progress-fill {
		height: 100%;
		background: var(--green);
		transition: width 0.3s ease;
	}

	main {
		max-width: 480px;
		margin: 0 auto;
		padding: 2rem 1.5rem;
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}

	.eyebrow {
		color: var(--green);
		font-size: 0.65rem;
		font-weight: 700;
		letter-spacing: 0.2em;
		text-transform: uppercase;
	}

	h1 {
		font-size: 2.2rem;
		line-height: 1.1;
		margin-bottom: 0.5rem;
	}

	.form {
		display: flex;
		flex-direction: column;
		gap: 2rem;
	}

	.field label {
		display: block;
		font-size: 0.85rem;
		color: var(--text-sub);
		margin-bottom: 0.4rem;
	}

	.upload {
		width: 100%;
		aspect-ratio: 16 / 9;
		border-radius: var(--radius-card);
		border: 2px dashed var(--border);
		background: var(--bg-elevated);
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 0.8rem;
		cursor: pointer;
		overflow: hidden;
		text-align: center;
		transition: all 0.2s;
	}

	.upload:hover {
		border-color: rgba(61, 158, 95, 0.4);
		background: rgba(61, 158, 95, 0.05);
	}

	.upload img {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.upload-icon {
		width: 3.5rem;
		height: 3.5rem;
		border-radius: 999px;
		background: rgba(61, 158, 95, 0.15);
		display: grid;
		place-items: center;
		color: var(--green);
	}

	.upload-icon .material-symbols-outlined {
		font-size: 1.5rem;
	}

	.upload span {
		font-size: 0.85rem;
		color: var(--text-sub);
		font-weight: 500;
	}

	.file-input {
		display: none;
	}

	.field.textarea textarea {
		width: 100%;
		background: transparent;
		border: 1px solid var(--border);
		border-radius: var(--radius-input);
		padding: 0.8rem;
		font-size: 0.95rem;
		color: var(--text);
		resize: none;
		font-family: inherit;
		transition: border-color 0.2s;
	}

	.field.textarea textarea:focus {
		outline: none;
		border-color: var(--green);
	}

	.error {
		color: #e05555;
		font-size: 0.85rem;
		text-align: center;
		padding: 0.5rem;
		background: rgba(224, 85, 85, 0.1);
		border-radius: 8px;
	}
</style>
