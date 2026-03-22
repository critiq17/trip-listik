<script lang="ts">
	import { onMount } from 'svelte';
	import { fly } from 'svelte/transition';
	import { apiFetch } from '$lib/api';
	import { resolvePhotoUrl } from '$lib/photos';

	let query = $state('');
	let country = $state('');
	let items = $state<Array<Record<string, any>>>([]);
	let loading = $state(false);
	let error = $state('');
	let timer: ReturnType<typeof setTimeout> | null = null;
	let ready = false;

	const countries = [
		{ code: '', label: 'All Regions' },
		{ code: 'IS', label: 'Iceland' },
		{ code: 'JP', label: 'Japan' },
		{ code: 'IT', label: 'Italy' },
		{ code: 'FR', label: 'France' },
		{ code: 'CH', label: 'Switzerland' },
		{ code: 'NO', label: 'Norway' },
		{ code: 'NZ', label: 'New Zealand' },
		{ code: 'TH', label: 'Thailand' },
		{ code: 'PT', label: 'Portugal' },
		{ code: 'ES', label: 'Spain' },
	];

	const search = async () => {
		loading = true;
		error = '';
		try {
			const data = await apiFetch<{ items?: Array<Record<string, any>> }>(
				`/v1/explore?q=${encodeURIComponent(query)}&country=${encodeURIComponent(country)}`
			);
			items = data.items ?? [];
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load explore results';
		} finally {
			loading = false;
		}
	};

	const debounceSearch = () => {
		if (timer) clearTimeout(timer);
		timer = setTimeout(() => search(), 300);
	};

	onMount(() => {
		ready = true;
		search();
	});

	$effect(() => {
		if (ready) {
			query;
			country;
			debounceSearch();
		}
	});

	let favoured = $state<Set<string>>(new Set());

	const toggleFav = (id: string) => {
		const next = new Set(favoured);
		if (next.has(id)) {
			next.delete(id);
		} else {
			next.add(id);
		}
		favoured = next;
	};
</script>

<svelte:head>
	<title>TripListik — Explore</title>
	<meta name="description" content="Search for destinations and discover recommended travel routes." />
</svelte:head>

<div class="page">
	<!-- Header -->
	<header class="header">
		<div class="header-top">
			<p class="eyebrow">Explore</p>
			<button class="icon-btn" aria-label="Menu">
				<span class="material-symbols-outlined">menu</span>
			</button>
		</div>
		<h1>Find the next route.</h1>
		<p class="subtitle">Search for destinations and itineraries.</p>
	</header>

	<main class="content">
		<!-- Search Form -->
		<div class="form">
			<!-- Destination -->
			<div class="field">
				<label class="field-label" for="explore-dest">Destination</label>
				<input
					id="explore-dest"
					type="text"
					class="underline-input"
					placeholder="Where to?"
					bind:value={query}
					autocomplete="off"
				/>
			</div>

			<!-- Country Filter -->
			<div class="field">
				<label class="field-label" for="explore-country">Country Filter</label>
				<div class="select-wrap">
					<select
						id="explore-country"
						class="underline-input"
						bind:value={country}
					>
						{#each countries as c}
							<option value={c.code}>{c.label}</option>
						{/each}
					</select>
					<span class="material-symbols-outlined select-chevron">expand_more</span>
				</div>
			</div>

			<!-- Search Button -->
			<button
				class="search-btn"
				onclick={(e) => { e.preventDefault(); search(); }}
			>
				Search Destinations
			</button>
		</div>

		<!-- Recommended Routes -->
		<section class="results">
			<p class="section-label">Recommended Routes</p>

			{#if loading}
				<div class="state">Searching…</div>
			{:else if error}
				<div class="state">{error}</div>
			{:else if items.length === 0}
				<div class="state">No routes found. Try a different search.</div>
			{:else}
				<div class="route-list">
					{#each items as item, i (item.id ?? item.title)}
						<article
							class="route-card"
							in:fly={{ y: 24, duration: 250, delay: i * 60 }}
						>
							<div class="route-img-wrap">
								{#if item.cover_photo_url || item.image}
									<div class="img-skeleton skeleton"></div>
									<img
										class="route-img"
										src={resolvePhotoUrl(item.cover_photo_url ?? item.image ?? '')}
										alt={item.title ?? item.name ?? ''}
										loading="lazy"
										onload={(e) => e.currentTarget.classList.add('loaded')}
										onerror={(e) => e.currentTarget.classList.add('error')}
									/>
								{:else}
									<div class="route-img-placeholder"></div>
								{/if}

								<button
									class="fav-btn"
									aria-label={favoured.has(item.id) ? 'Remove from saved' : 'Save route'}
									onclick={() => toggleFav(item.id)}
								>
									<span
										class="material-symbols-outlined"
										class:fill-icon={favoured.has(item.id)}
									>favorite</span>
								</button>
							</div>

							<div class="route-body">
								<div class="route-top">
									<h4>{item.title ?? item.name ?? 'Untitled Route'}</h4>
									{#if item.price ?? item.cost}
										<span class="price">{item.price ?? item.cost}</span>
									{/if}
								</div>
								<div class="route-meta">
									{#if item.duration ?? item.days}
										<span class="material-symbols-outlined meta-icon">schedule</span>
										<span>{item.duration ?? item.days} Days</span>
										<span class="dot">•</span>
									{/if}
									{#if item.distance ?? item.length ?? item.difficulty}
										<span class="material-symbols-outlined meta-icon">distance</span>
										<span>{item.distance ?? item.length ?? item.difficulty}</span>
									{/if}
								</div>
							</div>
						</article>
					{/each}
				</div>
			{/if}
		</section>
	</main>
</div>

<style>
.page {
	min-height: 100dvh;
	background: #161c18;
	color: #f1f5f9;
	padding-bottom: 96px;
}

/* ── Header ─────────────────────────────────────────────── */
.header {
	padding: 32px 24px 0;
	max-width: 480px;
	margin: 0 auto;
}

.header-top {
	display: flex;
	align-items: center;
	justify-content: space-between;
	margin-bottom: 12px;
}

.eyebrow {
	font-size: 10px;
	color: #4d9d6d;
	font-weight: 700;
	text-transform: uppercase;
	letter-spacing: 0.12em;
}

.icon-btn {
	width: 40px;
	height: 40px;
	display: flex;
	align-items: center;
	justify-content: center;
	color: white;
	background: none;
	border: none;
	cursor: pointer;
}

h1 {
	font-size: 30px;
	font-weight: 700;
	color: white;
	margin-bottom: 8px;
	line-height: 1.15;
}

.subtitle {
	font-size: 14px;
	color: #94a3b8;
	margin-bottom: 0;
}

/* ── Content ─────────────────────────────────────────────── */
.content {
	padding: 0 24px;
	max-width: 480px;
	margin: 0 auto;
	display: flex;
	flex-direction: column;
	gap: 24px;
}

/* ── Form ───────────────────────────────────────────────── */
.form {
	display: flex;
	flex-direction: column;
	gap: 24px;
	padding-top: 24px;
}

.field {
	display: flex;
	flex-direction: column;
	gap: 4px;
}

.field-label {
	font-size: 10px;
	color: #94a3b8;
	font-weight: 500;
	text-transform: uppercase;
	letter-spacing: 0.1em;
}

.underline-input {
	width: 100%;
	background: transparent;
	border: none;
	border-bottom: 1px solid rgba(255, 255, 255, 0.12);
	border-radius: 0;
	font-size: 18px;
	color: white;
	padding: 8px 0;
	outline: none;
	transition: border-bottom-color 0.2s ease;
	appearance: none;
	-webkit-appearance: none;
}

.underline-input::placeholder {
	color: #4a5568;
}

.underline-input:focus {
	border-bottom-color: #4d9d6d;
}

/* Select wrapper */
.select-wrap {
	position: relative;
}

.select-chevron {
	position: absolute;
	right: 0;
	top: 50%;
	transform: translateY(-50%);
	font-size: 20px;
	color: #64748b;
	pointer-events: none;
}

/* Search Button */
.search-btn {
	width: 100%;
	padding: 16px;
	border-radius: 8px;
	background: #4d9d6d;
	color: white;
	font-size: 18px;
	font-weight: 400;
	border: none;
	cursor: pointer;
	box-shadow: 0 8px 20px rgba(77, 157, 109, 0.25);
	transition: opacity 0.15s ease;
}

.search-btn:hover {
	opacity: 0.9;
}

/* ── Results ─────────────────────────────────────────────── */
.results {
	padding-bottom: 8px;
}

.section-label {
	font-size: 12px;
	color: #64748b;
	font-weight: 700;
	text-transform: uppercase;
	letter-spacing: 0.12em;
	margin-bottom: 16px;
}

.state {
	text-align: center;
	padding: 24px;
	border-radius: 12px;
	background: rgba(255, 255, 255, 0.04);
	color: #94a3b8;
	font-size: 14px;
}

.route-list {
	display: flex;
	flex-direction: column;
	gap: 16px;
}

/* ── Route Card ──────────────────────────────────────────── */
.route-card {
	position: relative;
	border-radius: 12px;
	overflow: hidden;
	background: rgba(255, 255, 255, 0.04);
	border: 1px solid rgba(255, 255, 255, 0.06);
	transition: transform 0.2s ease;
}

.route-card:hover {
	transform: scale(1.01);
}

.route-img-wrap {
	position: relative;
	height: 192px;
	overflow: hidden;
}

.img-skeleton {
	position: absolute;
	inset: 0;
}

.route-img {
	width: 100%;
	height: 100%;
	object-fit: cover;
	display: block;
	opacity: 0;
	transition: opacity 0.3s ease;
}

.route-img.loaded {
	opacity: 1;
}

.route-img.error {
	display: none;
}

.route-img-placeholder {
	width: 100%;
	height: 100%;
	background: rgba(255, 255, 255, 0.06);
}

.fav-btn {
	position: absolute;
	top: 16px;
	right: 16px;
	width: 36px;
	height: 36px;
	border-radius: 9999px;
	background: rgba(255, 255, 255, 0.15);
	backdrop-filter: blur(8px);
	-webkit-backdrop-filter: blur(8px);
	display: flex;
	align-items: center;
	justify-content: center;
	color: white;
	border: none;
	cursor: pointer;
	transition: background 0.2s ease;
}

.fav-btn .material-symbols-outlined {
	font-size: 20px;
}

.route-body {
	padding: 16px;
}

.route-top {
	display: flex;
	justify-content: space-between;
	align-items: flex-start;
	margin-bottom: 8px;
}

.route-top h4 {
	font-size: 18px;
	font-weight: 700;
	color: white;
}

.price {
	font-size: 18px;
	font-weight: 700;
	color: #4d9d6d;
	flex-shrink: 0;
}

.route-meta {
	display: flex;
	align-items: center;
	gap: 6px;
	color: #94a3b8;
	font-size: 12px;
}

.meta-icon {
	font-size: 14px;
	color: #94a3b8;
}

.dot {
	color: #94a3b8;
}
</style>
