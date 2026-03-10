<script lang="ts">
	import { goto } from '$app/navigation';
	import { apiFetch } from '$lib/api';
	import { setDraft } from '$lib/tripDraft';
	import { scalePress } from '$lib/actions/animate';

	let title = '';
	let destination = '';
	let coverPhotoPreview = '';
	let creating = false;
	let error = '';
	let fileInput: HTMLInputElement | null = null;

	const pickPhoto = () => {
		fileInput?.click();
	};

	const onFile = (event: Event) => {
		const target = event.currentTarget as HTMLInputElement;
		const file = target.files?.[0];
		if (!file) return;
		coverPhotoPreview = URL.createObjectURL(file);
	};

	const next = async () => {
		if (!title.trim()) {
			error = 'Trip name is required';
			return;
		}
		creating = true;
		error = '';
		try {
			const trip = await apiFetch<{ id: string }>('/v1/trips', {
				method: 'POST',
				body: JSON.stringify({
					title,
					city: destination,
					cover_photo_url: null,
					status: 'draft',
					visibility: 'public'
				})
			});
			setDraft(trip.id);
			goto('/create/step-2');
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to create trip';
		} finally {
			creating = false;
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

			<div class="field">
				<label for="trip-destination">Destination</label>
				<div class="field-row">
					<input id="trip-destination" bind:value={destination} placeholder="Where to?" />
					<span class="material-symbols-outlined">map</span>
				</div>
			</div>

			<div class="field">
				<label for="cover-photo">Cover Photo</label>
				<button type="button" class="upload" on:click={pickPhoto}>
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
					on:change={onFile}
				/>
			</div>
			{#if error}
				<p class="error">{error}</p>
			{/if}
		</div>

		<div class="actions">
			<button use:scalePress on:click={next} disabled={creating}>
				{creating ? 'Saving…' : 'Continue'}
				<span class="material-symbols-outlined">arrow_forward</span>
			</button>
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

	.field-row {
		display: flex;
		align-items: center;
		position: relative;
	}

	.field-row span {
		position: absolute;
		right: 0;
		color: var(--text-secondary);
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

	.actions button {
		width: 100%;
		padding: 1rem;
		border-radius: 12px;
		background: var(--primary);
		color: white;
		font-weight: 700;
		display: inline-flex;
		align-items: center;
		justify-content: center;
		gap: 0.4rem;
		box-shadow: 0 14px 30px rgba(77, 157, 109, 0.35);
	}

	.error {
		color: #e11d48;
		font-size: 0.8rem;
	}
</style>
