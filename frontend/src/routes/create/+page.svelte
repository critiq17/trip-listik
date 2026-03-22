<script lang="ts">
	import { goto } from '$app/navigation';
	import { setDraft, getDraft } from '$lib/tripDraft';
	import { onMount } from 'svelte';
	import { setupMainButton, setupBackButton } from '$lib/telegram';
	import CityAutocomplete from '$lib/components/CityAutocomplete.svelte';

	let draft = getDraft() || {};
	let title = $state(draft.title || '');
	let destination = $state(draft.destination || '');
	let countryCode = $state(draft.country_code || '');
	let lat = $state(draft.lat || 0);
	let lng = $state(draft.lon || 0);

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
	<title>TripListik — Create Trip</title>
	<meta name="description" content="Start planning your next trip adventure." />
</svelte:head>

<div class="page">
	<!-- Progress Bar -->
	<div class="progress-track">
		<div class="progress-fill" style="width: 33.33%; transition: width 500ms cubic-bezier(0.4, 0, 0.2, 1)"></div>
	</div>

	<main class="content">
		<p class="step-label">Step 1 of 3</p>
		<h1 class="serif">Where are you headed?</h1>

		<div class="form">
			<!-- Trip Name -->
			<div class="field">
				<label class="field-label" for="trip-name">Trip Name</label>
				<input
					id="trip-name"
					class="big-input"
					type="text"
					placeholder="e.g. Summer in Tuscany"
					bind:value={title}
					autocomplete="off"
				/>
			</div>

			<!-- Destination -->
			<div class="field autocomplete-field">
				<label class="field-label" for="trip-destination">Destination</label>
				<div class="input-icon-wrap">
					<CityAutocomplete
						bind:value={destination}
						bind:countryCode={countryCode}
						bind:lat={lat}
						bind:lng={lng}
						id="trip-destination"
						placeholder="Where to?"
					/>
					<span class="material-symbols-outlined input-suffix-icon">map</span>
				</div>
			</div>

			<!-- Cover Photo Upload Zone -->
			<div class="field">
				<label class="field-label" for="cover-photo-hint">Cover Photo</label>
				<div class="upload-zone" role="button" tabindex="0" aria-label="Add a cover image">
					<div class="upload-icon-circle">
						<span class="material-symbols-outlined">add_a_photo</span>
					</div>
					<p class="upload-hint">Add a cover image</p>
					<p class="upload-sub">You can add a photo in the final step</p>
				</div>
			</div>
		</div>

		<!-- Continue Button -->
		<button
			class="continue-btn"
			onclick={next}
			disabled={!title.trim()}
		>
			Continue
			<span class="material-symbols-outlined">arrow_forward</span>
		</button>
	</main>
</div>

<style>
.page {
	min-height: 100dvh;
	background: #161c18;
	color: #f1f5f9;
}

/* ── Progress Bar ────────────────────────────────────────── */
.progress-track {
	width: 100%;
	height: 4px;
	background: rgba(77, 157, 109, 0.2);
}

.progress-fill {
	height: 100%;
	background: #4d9d6d;
}

/* ── Content ─────────────────────────────────────────────── */
.content {
	max-width: 480px;
	margin: 0 auto;
	padding: 32px 24px 96px;
}

.step-label {
	font-size: 11px;
	color: #4d9d6d;
	font-weight: 700;
	text-transform: uppercase;
	letter-spacing: 0.15em;
	margin-bottom: 16px;
}

h1 {
	font-size: 38px;
	line-height: 1.1;
	color: white;
	margin-bottom: 40px;
}

/* ── Form ───────────────────────────────────────────────── */
.form {
	display: flex;
	flex-direction: column;
	gap: 48px;
	margin-bottom: 48px;
}

.field {
	display: flex;
	flex-direction: column;
	gap: 4px;
}

.field-label {
	font-size: 14px;
	color: #94a3b8;
	font-weight: 500;
	margin-bottom: 4px;
}

/* ── Text Inputs ─────────────────────────────────────────── */
.big-input {
	width: 100%;
	background: transparent;
	border: none;
	border-bottom: 2px solid rgba(255, 255, 255, 0.1);
	border-radius: 0;
	font-size: 24px;
	font-weight: 500;
	color: white;
	padding: 8px 0;
	outline: none;
	transition: border-bottom-color 0.2s ease;
}

.big-input::placeholder {
	color: rgba(255, 255, 255, 0.2);
}

.big-input:focus {
	border-bottom-color: #4d9d6d;
}

/* Destination field with map icon */
.input-icon-wrap {
	position: relative;
}

.input-suffix-icon {
	position: absolute;
	right: 0;
	top: 50%;
	transform: translateY(-50%);
	font-size: 24px;
	color: #64748b;
	pointer-events: none;
}

/* ── Upload Zone ─────────────────────────────────────────── */
.upload-zone {
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
	transition: border-color 0.2s ease, background 0.2s ease;
}

.upload-zone:hover {
	border-color: rgba(77, 157, 109, 0.4);
	background: rgba(77, 157, 109, 0.04);
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

.upload-sub {
	font-size: 11px;
	color: #4a5568;
}

/* ── Continue Button ─────────────────────────────────────── */
.continue-btn {
	width: 100%;
	display: flex;
	align-items: center;
	justify-content: center;
	gap: 8px;
	padding: 16px;
	border-radius: 12px;
	background: #4d9d6d;
	color: white;
	font-size: 16px;
	font-weight: 700;
	border: none;
	cursor: pointer;
	box-shadow: 0 8px 25px rgba(77, 157, 109, 0.3);
	transition: opacity 0.15s ease, transform 0.1s ease;
}

.continue-btn:hover:not(:disabled) {
	opacity: 0.92;
}

.continue-btn:active:not(:disabled) {
	transform: scale(0.98);
}

.continue-btn:disabled {
	opacity: 0.4;
	cursor: not-allowed;
}

.continue-btn .material-symbols-outlined {
	font-size: 20px;
}
</style>
