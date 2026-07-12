<script lang="ts">
	import { apiFetch } from '$lib/api';
	import { onDestroy, untrack } from 'svelte';

	let {
		value = $bindable(''),
		countryCode = $bindable(''),
		lat = $bindable(0),
		lng = $bindable(0),
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
	let controller: AbortController | null = null;  // ← AbortController for cancellation
	let wrapper: HTMLElement | null = null;
	let suppressSearch = false;

	// Sync external value changes (e.g. form reset from parent).
	// Use untrack on query so this effect only re-runs when `value` (the prop) changes,
	// NOT when the user types (which mutates query directly). Without untrack, every
	// keystroke triggers this effect, which reads the stale `value` and resets query.
	$effect(() => {
		const v = value;
		untrack(() => {
			if (v !== query && !isOpen) {
				query = v;
			}
		});
	});

	const search = async (q: string) => {
		if (q.length < 2) {
			results = [];
			isOpen = false;
			return;
		}

		// Cancel previous in-flight request
		controller?.abort();
		controller = new AbortController();

		loading = true;
		error = '';
		try {
			const data = await apiFetch<{ items: GeocodeItem[] }>(
				`/v1/geocode/search?q=${encodeURIComponent(q)}&limit=5`,
				{ signal: controller.signal } as RequestInit
			);
			results = (data.items ?? []).slice(0, 5);
			isOpen = results.length > 0;
		} catch (err: unknown) {
			// Ignore abort errors — they're intentional
			if (err instanceof Error && err.name === 'AbortError') return;
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
			controller?.abort();
			return;
		}
		if (timer) window.clearTimeout(timer);
		// Exactly 300ms debounce per spec
		timer = window.setTimeout(() => search(val), 300);
	};

	const select = (item: GeocodeItem) => {
		suppressSearch = true;
		value = item.city || item.description.split(',')[0];
		query = value;
		countryCode = item.country_code ? item.country_code.toUpperCase() : '';
		lat = parseFloat(item.lat) || 0;
		lng = parseFloat(item.lon) || 0;
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

	// Sync query back to the bound value when user types
	$effect(() => {
		if (value !== query) {
			value = query;
		}
	});

	// Trigger debounced search whenever query changes
	$effect(() => {
		const q = query;
		if (suppressSearch) {
			suppressSearch = false;
			return;
		}
		debounceSearch(q.trim());
	});

	onDestroy(() => {
		if (timer) window.clearTimeout(timer);
		controller?.abort();
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
		<ul class="dropdown" role="listbox">
			{#each results as item}
				<li role="option" aria-selected="false">
					<button type="button" onclick={() => select(item)}>
						<span class="material-symbols-outlined icon">location_on</span>
						<div class="item-text">
							<strong>{item.city || item.description.split(',')[0]}</strong>
							<small>{item.country}{item.country_code ? ` (${item.country_code.toUpperCase()})` : ''}</small>
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
		border-radius: var(--radius-input);
		background: var(--bg-input);
		border: 1px solid var(--border);
		color: var(--text);
		font-family: inherit;
		outline: none;
		transition: border-color 0.2s;
	}

	input:focus {
		border-color: var(--green);
	}

	.spinner {
		position: absolute;
		right: 1rem;
		width: 1rem;
		height: 1rem;
		border: 2px solid var(--border);
		border-top-color: var(--green);
		border-radius: 50%;
		animation: spin 0.8s linear infinite;
	}

	@keyframes spin {
		to { transform: rotate(360deg); }
	}

	.dropdown {
		position: absolute;
		top: calc(100% + 0.4rem);
		left: 0;
		right: 0;
		background: var(--bg-elevated);
		border: 1px solid var(--border);
		border-radius: var(--radius-card);
		list-style: none;
		margin: 0;
		padding: 0.4rem;
		z-index: 50;
		box-shadow: 0 4px 16px rgba(17, 20, 24, 0.08);
		max-height: 250px;
		overflow-y: auto;
		animation: dropIn 0.15s ease-out;
	}

	@keyframes dropIn {
		from { opacity: 0; transform: translateY(-4px); }
		to   { opacity: 1; transform: translateY(0); }
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
		color: var(--text);
		text-align: left;
		border-radius: var(--radius-input);
		cursor: pointer;
		transition: background 0.15s;
	}

	.dropdown button:hover,
	.dropdown button:focus {
		background: var(--bg-subtle);
		outline: none;
	}

	.icon {
		color: var(--text-sub);
		font-size: 1.2rem;
		flex-shrink: 0;
	}

	.item-text {
		display: flex;
		flex-direction: column;
		gap: 0.1rem;
	}

	.item-text strong {
		font-size: 0.95rem;
		font-weight: 600;
	}

	.item-text small {
		font-size: 0.75rem;
		color: var(--text-sub);
	}
</style>
