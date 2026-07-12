<script lang="ts">
	import { onMount } from 'svelte';
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
	<header class="top-bar">
		<span class="page-title">Explore</span>
	</header>

	<main class="content">
		<!-- Search Form -->
		<div class="form">
			<div class="search-wrap">
				<span class="material-symbols-outlined search-icon">search</span>
				<input
					id="explore-dest"
					type="text"
					class="search-input"
					placeholder="Where to?"
					bind:value={query}
					autocomplete="off"
				/>
			</div>

			<div class="select-wrap">
				<select
					id="explore-country"
					class="country-select"
					bind:value={country}
					aria-label="Country filter"
				>
					{#each countries as c}
						<option value={c.code}>{c.label}</option>
					{/each}
				</select>
				<span class="material-symbols-outlined select-chevron">expand_more</span>
			</div>
		</div>

		<!-- Recommended Routes -->
		<section class="results">
			<p class="section-label">Recommended routes</p>

			{#if loading}
				<div class="state">Searching…</div>
			{:else if error}
				<div class="state error">{error}</div>
			{:else if items.length === 0}
				<div class="state">No routes found. Try a different search.</div>
			{:else}
				<div class="route-list">
					{#each items as item (item.id ?? item.title)}
						<article class="route-card">
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
	background: var(--bg);
	color: var(--text);
	padding-bottom: 96px;
}

/* ── Header ─────────────────────────────────────────────── */
.top-bar {
	position: sticky;
	top: 0;
	z-index: 50;
	display: flex;
	align-items: center;
	height: 52px;
	padding: 0 16px;
	background: var(--bg);
	border-bottom: 1px solid var(--border);
	max-width: 480px;
	margin: 0 auto;
}

.page-title {
	font-size: 18px;
	font-weight: 700;
	color: var(--text);
}

/* ── Content ─────────────────────────────────────────────── */
.content {
	padding: 16px;
	max-width: 480px;
	margin: 0 auto;
	display: flex;
	flex-direction: column;
	gap: 20px;
}

/* ── Form ───────────────────────────────────────────────── */
.form {
	display: flex;
	flex-direction: column;
	gap: 10px;
}

.search-wrap {
	position: relative;
}

.search-icon {
	position: absolute;
	left: 12px;
	top: 50%;
	transform: translateY(-50%);
	font-size: 20px;
	color: var(--text-muted);
	pointer-events: none;
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

/* Select */
.select-wrap {
	position: relative;
}

.country-select {
	width: 100%;
	background: var(--bg-input);
	border: 1px solid var(--border);
	border-radius: var(--radius-input);
	font-size: 15px;
	color: var(--text);
	padding: 12px 40px 12px 14px;
	outline: none;
	appearance: none;
	-webkit-appearance: none;
	transition: border-color 0.15s ease;
}

.country-select:focus {
	border-color: var(--green);
}

.select-chevron {
	position: absolute;
	right: 12px;
	top: 50%;
	transform: translateY(-50%);
	font-size: 20px;
	color: var(--text-muted);
	pointer-events: none;
}

/* ── Results ─────────────────────────────────────────────── */
.section-label {
	font-size: 13px;
	color: var(--text-sub);
	font-weight: 600;
	margin-bottom: 12px;
}

.state {
	text-align: center;
	padding: 24px;
	border-radius: var(--radius-card);
	background: var(--bg-subtle);
	color: var(--text-sub);
	font-size: 14px;
}

.state.error {
	background: var(--danger-soft);
	color: var(--danger);
}

.route-list {
	display: flex;
	flex-direction: column;
	gap: 12px;
}

/* ── Route Card ──────────────────────────────────────────── */
.route-card {
	position: relative;
	border-radius: var(--radius-card);
	overflow: hidden;
	background: var(--bg-card);
	border: 1px solid var(--border);
}

.route-img-wrap {
	position: relative;
	height: 180px;
	overflow: hidden;
}

.img-skeleton {
	position: absolute;
	inset: 0;
	border-radius: 0;
}

.route-img {
	width: 100%;
	height: 100%;
	object-fit: cover;
	display: block;
	opacity: 0;
	transition: opacity 0.3s ease;
}

:global(.route-img.loaded) {
	opacity: 1;
}

:global(.route-img.error) {
	display: none;
}

.route-img-placeholder {
	width: 100%;
	height: 100%;
	background: var(--bg-subtle);
}

.fav-btn {
	position: absolute;
	top: 12px;
	right: 12px;
	width: 36px;
	height: 36px;
	border-radius: 50%;
	background: rgba(255, 255, 255, 0.92);
	display: flex;
	align-items: center;
	justify-content: center;
	color: var(--text);
	border: none;
	cursor: pointer;
}

.fav-btn .fill-icon {
	color: var(--danger);
}

.fav-btn .material-symbols-outlined {
	font-size: 20px;
}

.route-body {
	padding: 12px 14px;
}

.route-top {
	display: flex;
	justify-content: space-between;
	align-items: flex-start;
	gap: 8px;
	margin-bottom: 4px;
}

.route-top h4 {
	font-size: 16px;
	font-weight: 600;
	color: var(--text);
}

.price {
	font-size: 15px;
	font-weight: 700;
	color: var(--green);
	flex-shrink: 0;
}

.route-meta {
	display: flex;
	align-items: center;
	gap: 6px;
	color: var(--text-sub);
	font-size: 13px;
}

.meta-icon {
	font-size: 15px;
	color: var(--text-sub);
}

.dot {
	color: var(--text-muted);
}
</style>
