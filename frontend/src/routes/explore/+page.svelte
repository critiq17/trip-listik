<script lang="ts">
	import { onMount } from 'svelte';
	import { apiFetch } from '$lib/api';
	import { staggerList } from '$lib/actions/animate';

	let query = $state('');
	let country = $state('');
	let items = $state<Array<Record<string, any>>>([]);
	let loading = $state(false);
	let error = $state('');
	let timer: ReturnType<typeof setTimeout> | null = null;
	let ready = false;

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
		timer = setTimeout(() => {
			search();
		}, 300);
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
</script>

<section class="explore-page">
	<header class="header">
		<div class="header-row">
			<span class="eyebrow">Explore</span>
			<button class="icon-btn" aria-label="Menu">
				<span class="material-symbols-outlined">menu</span>
			</button>
		</div>
		<h1>Find the next route.</h1>
		<p>Search for destinations and itineraries.</p>
	</header>

	<main class="content">
		<div class="form">
			<div class="field">
				<label for="explore-destination">Destination</label>
				<input
					id="explore-destination"
					type="text"
					placeholder="Where to?"
					bind:value={query}
				/>
			</div>
			<div class="field select">
				<label for="explore-country">Country Filter</label>
				<select id="explore-country" bind:value={country}>
					<option value="">All Regions</option>
					<option value="europe">Europe</option>
					<option value="asia">Asia</option>
					<option value="north-america">North America</option>
				</select>
				<span class="material-symbols-outlined">expand_more</span>
			</div>
			<button
				class="search-btn"
				onclick={(event) => {
					event.preventDefault();
					search();
				}}
			>
				Search Destinations
			</button>
		</div>

		<section class="results">
			<div class="results-header">
				<h3>Recommended Routes</h3>
			</div>
			{#if loading}
				<div class="state">Searching…</div>
			{:else if error}
				<div class="state">{error}</div>
			{:else if items.length === 0}
				<div class="state">No routes found.</div>
			{:else}
				<div class="cards" use:staggerList>
					{#each items as item (item.id ?? item.title)}
						<article class="route-card" data-item>
							<div class="route-image-wrap">
								<div class="img-skeleton skeleton"></div>
								<img
									class="trip-img route-image-element"
									src={item.cover_photo_url ?? item.image ?? ''}
									alt={item.title ?? item.name ?? ''}
									onload={(e) => e.currentTarget.classList.add('loaded')}
									onerror={(e) => e.currentTarget.classList.add('error')}
								/>
							</div>
							<div class="route-body">
								<div class="route-top">
									<h4>{item.title ?? item.name ?? 'Untitled Route'}</h4>
									<span class="price">{item.price ?? item.cost ?? ''}</span>
								</div>
								<div class="route-meta">
									<span class="material-symbols-outlined">schedule</span>
									<span>{item.duration ?? item.days ?? '—'}</span>
									<span class="dot">•</span>
									<span class="material-symbols-outlined">distance</span>
									<span>{item.distance ?? item.length ?? '—'}</span>
								</div>
							</div>
							<button class="fav-btn" aria-label="Save">
								<span class="material-symbols-outlined">favorite</span>
							</button>
						</article>
					{/each}
				</div>
			{/if}
		</section>
	</main>
</section>

<style>
	.explore-page {
		min-height: 100dvh;
		background: var(--background-dark);
		color: var(--text-primary);
		padding-bottom: 5.5rem;
	}

	.header {
		padding: 2rem 1.5rem 1rem;
		max-width: 480px;
		margin: 0 auto;
	}

	.header-row {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: 1.5rem;
	}

	.eyebrow {
		color: var(--primary);
		font-size: 0.65rem;
		font-weight: 700;
		letter-spacing: 0.2em;
		text-transform: uppercase;
	}

	.icon-btn {
		width: 2.5rem;
		height: 2.5rem;
		display: grid;
		place-items: center;
	}

	h1 {
		font-size: 1.85rem;
		font-weight: 700;
		margin-bottom: 0.4rem;
	}

	p {
		color: var(--text-secondary);
		font-size: 0.85rem;
	}

	.content {
		padding: 0 1.5rem 1.5rem;
		display: flex;
		flex-direction: column;
		gap: 2rem;
		max-width: 480px;
		margin: 0 auto;
	}

	.form {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}

	.field label {
		display: block;
		font-size: 0.65rem;
		font-weight: 600;
		letter-spacing: 0.18em;
		text-transform: uppercase;
		color: var(--text-secondary);
		margin-bottom: 0.3rem;
	}

	.field input,
	.field select {
		width: 100%;
		background: transparent;
		border: none;
		border-bottom: 1px solid rgba(255, 255, 255, 0.18);
		padding: 0.6rem 0;
		font-size: 1.05rem;
		color: var(--text-primary);
	}

	.field select {
		appearance: none;
	}

	.field.select {
		position: relative;
	}

	.field.select span {
		position: absolute;
		right: 0;
		bottom: 0.55rem;
		color: var(--text-secondary);
	}

	.search-btn {
		width: 100%;
		padding: 0.95rem;
		border-radius: 10px;
		background: var(--primary);
		color: white;
		font-weight: 600;
		font-size: 1rem;
		box-shadow: 0 12px 24px rgba(77, 157, 109, 0.3);
	}

	.results-header h3 {
		color: var(--text-secondary);
		font-size: 0.7rem;
		letter-spacing: 0.18em;
		text-transform: uppercase;
		font-weight: 700;
		margin-bottom: 1rem;
	}

	.cards {
		display: flex;
		flex-direction: column;
		gap: 1rem;
		padding-bottom: 2rem;
	}

	.route-card {
		position: relative;
		background: rgba(255, 255, 255, 0.03);
		border-radius: 12px;
		overflow: hidden;
		border: 1px solid rgba(255, 255, 255, 0.08);
		box-shadow: var(--shadow-card);
	}

	.route-body {
		padding: 1rem;
	}

	.route-top {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		margin-bottom: 0.4rem;
	}

	.route-top h4 {
		font-size: 1.1rem;
		font-weight: 700;
	}

	.price {
		color: var(--primary);
		font-weight: 700;
	}

	.route-meta {
		display: flex;
		align-items: center;
		gap: 0.35rem;
		color: var(--text-secondary);
		font-size: 0.75rem;
	}

	.route-meta .material-symbols-outlined {
		font-size: 1rem;
	}

	.dot {
		margin: 0 0.2rem;
	}

	.fav-btn {
		position: absolute;
		top: 1rem;
		right: 1rem;
		width: 2rem;
		height: 2rem;
		border-radius: 999px;
		background: rgba(255, 255, 255, 0.2);
		backdrop-filter: blur(10px);
		color: white;
		display: grid;
		place-items: center;
	}

	.state {
		padding: 1.2rem;
		background: rgba(255, 255, 255, 0.04);
		border-radius: 12px;
		text-align: center;
		color: var(--text-secondary);
	}
</style>
