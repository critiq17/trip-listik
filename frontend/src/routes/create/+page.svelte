<script lang="ts">
	import { goto } from '$app/navigation';
	import { setDraft, getDraft } from '$lib/tripDraft';
	import { onMount, onDestroy } from 'svelte';
	import { setupMainButton, hideMainButton, setupBackButton, isTelegramEnv } from '$lib/telegram';
	import CityAutocomplete from '$lib/components/CityAutocomplete.svelte';

	let draft = getDraft() || {};
	let title = $state(draft.title || '');
	let destination = $state(draft.destination || '');
	let countryCode = $state(draft.country_code || '');
	let lat = $state(draft.lat || 0);
	let lng = $state(draft.lon || 0);
	// Telegram renders the native MainButton at the bottom; the inline button
	// is only a fallback for plain-browser sessions (local dev).
	let showFallbackButton = $state(false);

	$effect(() => {
		if (title.trim()) {
			setupMainButton('Continue', next, true, true);
		} else {
			setupMainButton('Continue', next, true, false);
		}
	});

	onMount(() => {
		setupBackButton(() => history.back());
		showFallbackButton = !isTelegramEnv();
	});

	onDestroy(() => {
		hideMainButton();
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
	<header class="top-bar">
		<button class="back-btn" aria-label="Go back" onclick={() => history.back()}>
			<span class="material-symbols-outlined">arrow_back</span>
		</button>
		<span class="header-title">New Trip</span>
		<div class="header-spacer"></div>
	</header>

	<div class="progress-track">
		<div class="progress-fill" style="width: 33.33%"></div>
	</div>

	<main class="content">
		<p class="step-label">Step 1 of 3</p>
		<h1>Where are you headed?</h1>

		<div class="form">
			<div class="field">
				<label class="field-label" for="trip-name">Trip name</label>
				<input
					id="trip-name"
					class="text-input"
					type="text"
					placeholder="e.g. Summer in Tuscany"
					bind:value={title}
					autocomplete="off"
				/>
			</div>

			<div class="field">
				<label class="field-label" for="trip-destination">Destination</label>
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

		{#if showFallbackButton}
			<button
				class="continue-btn"
				onclick={next}
				disabled={!title.trim()}
			>
				Continue
			</button>
		{/if}
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
	max-width: 480px;
	margin: 0 auto;
	padding: 24px 16px 96px;
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
	line-height: 1.2;
	color: var(--text);
	margin-bottom: 24px;
}

/* ── Form ───────────────────────────────────────────────── */
.form {
	display: flex;
	flex-direction: column;
	gap: 20px;
	margin-bottom: 32px;
}

.field {
	display: flex;
	flex-direction: column;
	gap: 6px;
}

.field-label {
	font-size: 13px;
	color: var(--text-sub);
	font-weight: 500;
}

.text-input {
	width: 100%;
	background: var(--bg-input);
	border: 1px solid var(--border);
	border-radius: var(--radius-input);
	font-size: 16px;
	color: var(--text);
	padding: 12px 14px;
	outline: none;
	transition: border-color 0.15s ease;
}

.text-input:focus {
	border-color: var(--green);
}

/* ── Continue Button ─────────────────────────────────────── */
.continue-btn {
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

.continue-btn:disabled {
	opacity: 0.4;
	cursor: not-allowed;
}
</style>
