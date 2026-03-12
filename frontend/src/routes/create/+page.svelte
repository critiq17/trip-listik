<script lang="ts">
	import { goto } from '$app/navigation';
	import { setDraft } from '$lib/tripDraft';
	import { scalePress } from '$lib/actions/animate';
	import { apiFetch, presignTripPhoto, uploadSignedPhoto, getPublicPhotoURL } from '$lib/api';
	import { onDestroy, onMount } from 'svelte';
	import { setupMainButton, hideMainButton, setMainButtonState, setupBackButton } from '$lib/telegram';
	import CityAutocomplete from '$lib/components/CityAutocomplete.svelte';

	let title = $state('');
	let destination = $state('');
	let countryCode = $state('');
	let coverPhotoPreview = $state('');
	let coverPhotoFile: File | null = null;
	let creating = $state(false);
	let error = $state('');
	let fileInput: HTMLInputElement | null = null;

	$effect(() => {
		if (title.trim() && !creating) {
			setupMainButton('Continue', next, true, true);
		} else {
			setupMainButton('Continue', next, true, false);
		}
	});

	onMount(() => {
		setupBackButton(() => history.back());
	});

	const pickPhoto = () => {
		fileInput?.click();
	};

	const onFile = (event: Event) => {
		const target = event.currentTarget as HTMLInputElement;
		const file = target.files?.[0];
		if (!file) return;
		if (coverPhotoPreview) URL.revokeObjectURL(coverPhotoPreview);
		coverPhotoFile = file;
		coverPhotoPreview = URL.createObjectURL(file);
	};

	onDestroy(() => {
		if (coverPhotoPreview) URL.revokeObjectURL(coverPhotoPreview);
	});

	const next = async () => {
		if (!title.trim() || creating) return;
		creating = true;
		setMainButtonState(true);
		error = '';
		try {
			// Step 1: create trip draft first (we need its ID for presign)
			const trip = await apiFetch<{ id: string }>('/v1/trips', {
				method: 'POST',
				body: JSON.stringify({
					title,
					city: destination,
					country_code: countryCode,
					cover_photo_url: null,
					status: 'draft',
					visibility: 'public'
				})
			});

			// Step 2: upload cover photo if file was selected
			if (coverPhotoFile) {
				try {
					const presign = await presignTripPhoto(trip.id, coverPhotoFile.name, coverPhotoFile.type);
					await uploadSignedPhoto(presign.signed_url, presign.token, coverPhotoFile);
					const publicUrl = getPublicPhotoURL(presign.path);
					await apiFetch(`/v1/trips/${trip.id}`, {
						method: 'PATCH',
						body: JSON.stringify({ cover_photo_url: publicUrl })
					});
				} catch {
					// Non-fatal: continue without cover photo
					console.warn('Cover photo upload failed, continuing without it');
				}
			}

			setDraft(trip.id);
			goto('/create/step-2');
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to create trip';
		} finally {
			creating = false;
			setMainButtonState(false);
		}
	};
</script>

<section class="create-step">
	<div class="progress-bar">
		<div class="progress-fill"></div>
	</div>

	<main>
		<p class="eyebrow">Step 1 of 3</p>
		<h1 class="serif-text">Where are you headed?</h1>

		<div class="form">
			<div class="field">
				<label for="trip-name">Trip Name</label>
				<input id="trip-name" bind:value={title} placeholder="e.g. Summer in Tuscany" />
			</div>

			<div class="field autocomplete-field">
				<label for="trip-destination">Destination</label>
				<CityAutocomplete 
					bind:value={destination} 
					bind:countryCode={countryCode}
					id="trip-destination" 
					placeholder="Where to?" 
				/>
			</div>

			<div class="field">
				<label for="cover-photo">Cover Photo</label>
				<button type="button" class="upload" onclick={pickPhoto}>
					{#if coverPhotoPreview}
						<img src={coverPhotoPreview} alt="Cover preview" />
					{:else}
						<div class="upload-icon">
							<span class="material-symbols-outlined">add_a_photo</span>
						</div>
						<span>Add a cover image</span>
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
			{#if error}
				<p class="error">{error}</p>
			{/if}
		</div>
	</main>
</section>

<style>
	.create-step {
		min-height: 100dvh;
		background: var(--background-dark);
		color: var(--text-primary);
		padding-bottom: 2rem;
	}

	.progress-bar {
		width: 100%;
		height: 4px;
		background: rgba(77, 157, 109, 0.2);
	}

	.progress-fill {
		height: 100%;
		width: 33.33%;
		background: var(--primary);
	}

	main {
		max-width: 480px;
		margin: 0 auto;
		padding: 2rem 1.5rem 3rem;
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}

	.eyebrow {
		color: var(--primary);
		font-size: 0.65rem;
		font-weight: 700;
		letter-spacing: 0.2em;
		text-transform: uppercase;
	}

	h1 {
		font-size: 2.4rem;
		line-height: 1.05;
	}

	.form {
		display: flex;
		flex-direction: column;
		gap: 2rem;
	}

	.field label {
		display: block;
		font-size: 0.85rem;
		color: var(--text-secondary);
		margin-bottom: 0.4rem;
	}

	.field input {
		width: 100%;
		background: transparent;
		border: none;
		border-bottom: 2px solid rgba(255, 255, 255, 0.1);
		padding: 0.6rem 0;
		font-size: 1.4rem;
		font-weight: 500;
		color: var(--text-primary);
	}


	.upload {
		width: 100%;
		aspect-ratio: 16 / 9;
		border-radius: 12px;
		border: 2px dashed rgba(255, 255, 255, 0.12);
		background: rgba(255, 255, 255, 0.02);
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 0.4rem;
		cursor: pointer;
		overflow: hidden;
		text-align: center;
	}

	.upload img {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.upload-icon {
		width: 3rem;
		height: 3rem;
		border-radius: 999px;
		background: rgba(77, 157, 109, 0.12);
		display: grid;
		place-items: center;
		color: var(--primary);
	}

	.upload span {
		font-size: 0.7rem;
		color: var(--text-secondary);
		font-weight: 600;
	}

	.file-input {
		display: none;
	}


	.error {
		color: #e11d48;
		font-size: 0.8rem;
	}
</style>
