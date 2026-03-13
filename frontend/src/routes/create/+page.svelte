<script lang="ts">
	import { goto } from '$app/navigation';
	import { setDraft, getDraft } from '$lib/tripDraft';
	import { onDestroy, onMount } from 'svelte';
	import { setupMainButton, hideMainButton, setMainButtonState, setupBackButton } from '$lib/telegram';
	import CityAutocomplete from '$lib/components/CityAutocomplete.svelte';

	let draft = getDraft() || {};
	let title = $state(draft.title || '');
	let destination = $state(draft.destination || '');
	let countryCode = $state(draft.country_code || '');
	let lat = $state(draft.lat || 0);
	let lng = $state(draft.lon || 0);

	// Note: Cover photo handling happens at the end (Step 3) now
	// to avoid uploading photos for abandoned trips.

	$effect(() => {
		if (title.trim()) {
			setupMainButton('Continue', next, true, true);
		} else {
			setupMainButton('Continue', next, true, false);
		}
	});

	onMount(() => {
		setupBackButton(() => history.back());
	});

	const next = () => {
		if (!title.trim()) return;
		setDraft({
			title,
			destination,
			country_code: countryCode,
			lat,
			lon: lng
		});
		goto('/create/step-2');
	};
</script>

<svelte:head>
	<title>Create Trip - Step 1</title>
</svelte:head>

<section class="create-step">
	<div class="progress-bar">
		<div class="progress-fill" style="width: 33.33%"></div>
	</div>

	<main>
		<p class="eyebrow">Step 1 of 3</p>
		<h1 class="serif-text">Where are you headed?</h1>

		<div class="form">
			<div class="field">
				<label for="trip-name">Trip Name</label>
				<input id="trip-name" bind:value={title} placeholder="e.g. Summer in Tuscany" autocomplete="off" />
			</div>

			<div class="field autocomplete-field">
				<label for="trip-destination">Destination</label>
				<CityAutocomplete 
					bind:value={destination} 
					bind:countryCode={countryCode}
					bind:lat={lat}
					bind:lng={lng}
					id="trip-destination" 
					placeholder="Where to?" 
				/>
			</div>
		</div>
	</main>
</section>

<style>
	.create-step {
		min-height: 100dvh;
		background: var(--bg);
		color: var(--text);
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
		padding: 2rem 1.5rem 3rem;
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

	.field input {
		width: 100%;
		background: transparent;
		border: none;
		border-bottom: 2px solid var(--border);
		padding: 0.6rem 0;
		font-size: 1.4rem;
		font-weight: 500;
		color: var(--text);
		border-radius: 0; /* iOS fix */
	}

	.field input:focus {
		outline: none;
		border-bottom-color: var(--green);
	}
</style>
