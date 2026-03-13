<script lang="ts">
	import { apiFetch } from '$lib/api';
	import { onDestroy } from 'svelte';

	let {
		value = $bindable(''),
		countryCode = $bindable(''),
		placeholder = 'Search cities...',
		required = false,
		id = ''
	} = $props();

	type GeocodeItem = {
		id: string;
		description: string;
		city: string;
		country: string;
		country_code: string;
		lat: string;
		lon: string;
	};

	let query = $state(value);
	let results = $state<GeocodeItem[]>([]);
	let loading = $state(false);
	let error = $state('');
	let isOpen = $state(false);
	let timer: number | null = null;
	let wrapper: HTMLElement | null = null;
	let suppressSearch = false;

	// Watch external value changes (e.g. reset form)
	$effect(() => {
		if (value !== query && !isOpen) {
			query = value;
		}
	});

	const search = async (q: string) => {
		if (q.length < 2) {
			results = [];
			isOpen = false;
			return;
		}
		loading = true;
		error = '';
		try {
			const data = await apiFetch<{ items: GeocodeItem[] }>(
				`/v1/geocode/search?q=${encodeURIComponent(q)}`
			);
			results = data.items ?? [];
			isOpen = results.length > 0;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Search failed';
			isOpen = false;
		} finally {
			loading = false;
		}
	};

	const debounceSearch = (val: string) => {
		if (val.length < 2) {
			results = [];
			isOpen = false;
			if (timer) window.clearTimeout(timer);
			return;
		}
		if (timer) window.clearTimeout(timer);
		timer = window.setTimeout(() => search(val), 350);
	};

	const select = (item: GeocodeItem) => {
		suppressSearch = true;
		value = item.city || item.description.split(',')[0];
		query = value;
		countryCode = item.country_code ? item.country_code.toUpperCase() : '';
		isOpen = false;
		results = [];
	};

	const onClickOutside = (e: MouseEvent) => {
		if (wrapper && !wrapper.contains(e.target as Node)) {
			isOpen = false;
		}
	};

	$effect(() => {
		if (typeof window !== 'undefined') {
			window.addEventListener('click', onClickOutside);
			return () => window.removeEventListener('click', onClickOutside);
		}
	});

	$effect(() => {
		query;
		if (value !== query) {
			value = query;
		}
		if (suppressSearch) {
			suppressSearch = false;
			return;
		}
		debounceSearch(query.trim());
	});

	onDestroy(() => {
		if (timer) window.clearTimeout(timer);
	});
</script>

<div class="autocomplete" bind:this={wrapper}>
	<div class="input-wrap">
		<input
			{id}
			type="search"
			{placeholder}
			{required}
			bind:value={query}
			inputmode="search"
			autocapitalize="none"
			onfocus={() => { if (results.length > 0) isOpen = true; }}
			autocomplete="off"
		/>
		{#if loading}
			<div class="spinner"></div>
		{/if}
	</div>

	{#if isOpen}
		<ul class="dropdown">
			{#each results as item}
				<li>
					<button type="button" onclick={() => select(item)}>
						<span class="material-symbols-outlined icon">location_on</span>
						<div class="item-text">
							<strong>{item.city || item.description.split(',')[0]}</strong>
							<small>{item.country} {item.country_code ? `(${item.country_code.toUpperCase()})` : ''}</small>
						</div>
					</button>
				</li>
			{/each}
		</ul>
	{/if}
</div>

<style>
	.autocomplete {
		position: relative;
		width: 100%;
	}

	.input-wrap {
		position: relative;
		display: flex;
		align-items: center;
	}

	input {
		width: 100%;
		padding: 0.9rem 1rem;
		border-radius: var(--radius-xl);
		background: rgba(255, 255, 255, 0.04);
		border: 1px solid rgba(255, 255, 255, 0.08);
		color: var(--text-primary);
		font-family: inherit;
		outline: none;
		transition: border-color 0.2s;
	}

	input:focus {
		border-color: var(--primary);
	}

	.spinner {
		position: absolute;
		right: 1rem;
		width: 1rem;
		height: 1rem;
		border: 2px solid rgba(255, 255, 255, 0.2);
		border-top-color: var(--primary);
		border-radius: 50%;
		animation: spin 0.8s linear infinite;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}

	.dropdown {
		position: absolute;
		top: calc(100% + 0.4rem);
		left: 0;
		right: 0;
		background: #1c231f;
		border: 1px solid rgba(255, 255, 255, 0.08);
		border-radius: var(--radius-xl);
		list-style: none;
		margin: 0;
		padding: 0.4rem;
		z-index: 50;
		box-shadow: 0 10px 30px rgba(0, 0, 0, 0.5);
		max-height: 250px;
		overflow-y: auto;
	}

	.dropdown li {
		margin: 0;
		padding: 0;
	}

	.dropdown button {
		width: 100%;
		display: flex;
		align-items: center;
		gap: 0.8rem;
		padding: 0.8rem;
		background: transparent;
		border: none;
		color: var(--text-primary);
		text-align: left;
		border-radius: var(--radius-md);
		cursor: pointer;
		transition: background 0.2s;
	}

	.dropdown button:hover,
	.dropdown button:focus {
		background: rgba(255, 255, 255, 0.06);
		outline: none;
	}

	.icon {
		color: var(--text-secondary);
		font-size: 1.2rem;
	}

	.item-text {
		display: flex;
		flex-direction: column;
	}

	.item-text strong {
		font-size: 0.95rem;
		font-weight: 600;
	}

	.item-text small {
		font-size: 0.75rem;
		color: var(--text-secondary);
	}
</style>
