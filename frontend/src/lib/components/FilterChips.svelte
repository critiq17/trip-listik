<script lang="ts">
	export let filters: string[] = [];
	export let active: string;
</script>

<div class="chips">
	{#each filters as filter}
		<button class:active={filter === active} on:click={() => (active = filter)}>
			{filter}
		</button>
	{/each}
</div>

<style>
	.chips {
		display: flex;
		gap: 1.5rem;
		overflow-x: auto;
		padding: 0.5rem 0 0.9rem;
		scrollbar-width: none;
	}

	.chips::-webkit-scrollbar {
		display: none;
	}

	button {
		position: relative;
		padding-bottom: 0.55rem;
		color: var(--text-secondary);
		font-weight: 600;
		font-size: 0.85rem;
		white-space: nowrap;
		transition: color 0.2s ease;
	}

	button.active {
		color: var(--text-primary);
		font-weight: 700;
	}

	button.active::after {
		content: '';
		position: absolute;
		left: 0;
		right: 0;
		bottom: 0;
		height: 2px;
		background: var(--primary);
		border-radius: 999px;
	}
</style>
