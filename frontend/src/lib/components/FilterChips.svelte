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
		gap: 0.6rem;
		overflow-x: auto;
		padding: 0.1rem 0 0.65rem;
		scrollbar-width: none;
	}

	.chips::-webkit-scrollbar {
		display: none;
	}

	button {
		padding: 0.6rem 1rem;
		border-radius: var(--radius-pill);
		background: rgba(255, 255, 255, 0.05);
		color: var(--text-secondary);
		font-weight: 700;
		font-size: 0.85rem;
		white-space: nowrap;
		transition: all 0.2s ease;
		border: 1px solid rgba(255, 255, 255, 0.06);
		backdrop-filter: blur(14px);
	}

	button.active {
		background: linear-gradient(135deg, rgba(32, 146, 186, 0.32), rgba(122, 234, 244, 0.14));
		color: white;
		border-color: rgba(122, 234, 244, 0.5);
		box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.12);
	}
</style>
