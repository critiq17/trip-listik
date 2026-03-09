<script lang="ts">
	import { goto } from '$app/navigation';
	import { apiFetch } from '$lib/api';
	import { setDraft } from '$lib/tripDraft';

	let title = '';
	let destination = '';
	let coverPhotoURL = '';
	let creating = false;
	let error = '';

	const next = async () => {
		if (!title.trim()) {
			error = 'Title is required';
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
					cover_photo_url: coverPhotoURL,
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

<section class="container">
	<header class="header">
		<p class="eyebrow">Create Trip</p>
		<h1 class="headline">Set the core idea.</h1>
		<p class="subtle">Trip name, destination and cover image URL for the first version.</p>
	</header>

	<div class="progress">
		<div class="progress-top">
			<span>Step 1 of 3</span>
			<strong>33%</strong>
		</div>
		<div class="bar"><span></span></div>
	</div>

	<div class="form panel">
		<label>
			<span>Trip name</span>
			<input bind:value={title} placeholder="Summer in Lisbon" />
		</label>
		<label>
			<span>Destination</span>
			<div class="destination">
				<span>✈</span>
				<input bind:value={destination} placeholder="Lisbon, PT" />
			</div>
		</label>
		<label>
			<span>Cover image URL</span>
			<input bind:value={coverPhotoURL} placeholder="https://images.example/trip-cover.jpg" />
		</label>
		<div class="upload-box" class:has-image={Boolean(coverPhotoURL)}>
			{#if coverPhotoURL}
				<img src={coverPhotoURL} alt="Cover preview" />
			{:else}
				<div>
					<strong>Add a cover preview</strong>
					<p>Paste a real image URL now. Photo upload for trip galleries lives in trip detail.</p>
				</div>
			{/if}
		</div>
		{#if error}
			<div class="error">{error}</div>
		{/if}
	</div>

	<footer class="footer">
		<button class="cta" on:click={next} disabled={creating}>
			{creating ? 'Saving…' : 'Next Step'}
		</button>
	</footer>
</section>

<style>
	.header {
		display: grid;
		gap: 0.45rem;
		margin-bottom: 1rem;
	}

	.progress {
		margin-bottom: 1rem;
	}

	.progress-top {
		display: flex;
		justify-content: space-between;
		align-items: center;
		color: var(--text-secondary);
		font-size: 0.84rem;
		margin-bottom: 0.55rem;
	}

	.bar {
		height: 8px;
		border-radius: var(--radius-pill);
		background: rgba(255, 255, 255, 0.08);
		overflow: hidden;
	}

	.bar span {
		display: block;
		width: 33%;
		height: 100%;
		border-radius: inherit;
		background: var(--accent-grad);
	}

	.form {
		display: grid;
		gap: 1.25rem;
		padding: 1.2rem;
		border-radius: var(--radius-2xl);
	}

	label {
		display: grid;
		gap: 0.5rem;
		font-size: 0.85rem;
		color: var(--text-secondary);
	}

	input {
		background: transparent;
		border: none;
		border-bottom: 1px solid rgba(255, 255, 255, 0.16);
		padding: 0.7rem 0;
		font-size: 1rem;
	}

	.destination {
		display: flex;
		align-items: center;
		gap: 0.7rem;
	}

	.upload-box {
		min-height: 180px;
		border: 1px dashed rgba(255, 255, 255, 0.22);
		border-radius: var(--radius-xl);
		display: flex;
		align-items: center;
		justify-content: center;
		color: var(--text-secondary);
		text-align: center;
		padding: 1rem;
		overflow: hidden;
		background: rgba(255, 255, 255, 0.03);
	}

	.upload-box img {
		width: 100%;
		height: 100%;
		object-fit: cover;
		border-radius: 18px;
	}

	.upload-box.has-image {
		padding: 0.4rem;
	}

	.error {
		color: var(--danger);
		font-size: 0.85rem;
	}

	.footer {
		position: sticky;
		bottom: 0;
		margin-top: 2rem;
		padding-bottom: 5.8rem;
	}

	.cta {
		width: 100%;
		padding: 1rem;
		border-radius: var(--radius-xl);
		background: var(--accent-grad);
		font-weight: 800;
		color: #fff;
		box-shadow: var(--shadow-glow);
	}
</style>
