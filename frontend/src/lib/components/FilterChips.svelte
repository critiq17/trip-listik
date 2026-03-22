<script lang="ts">
	let {
		filters = ['All', 'Friends', 'Popular', 'Nearby'],
		active = $bindable('All')
	}: {
		filters?: string[];
		active?: string;
	} = $props();
</script>

<div class="filter-pills no-scrollbar">
	{#each filters as filter}
		<button
			class="pill"
			class:active={active === filter}
			onclick={() => (active = filter)}
		>
			{filter}
		</button>
	{/each}
</div>

<style>
.filter-pills {
	display: flex;
	flex-direction: row;
	gap: 24px;
	overflow-x: auto;
	padding: 16px 0;
}

.pill {
	flex-shrink: 0;
	background: none;
	border: none;
	font-size: 15px;
	font-weight: 500;
	color: #64748b;
	cursor: pointer;
	padding-bottom: 8px;
	border-bottom: 2px solid transparent;
	transition: color 0.2s ease, border-color 0.2s ease;
	white-space: nowrap;
	-webkit-tap-highlight-color: transparent;
}

.pill.active {
	color: white;
	font-weight: 600;
	border-bottom-color: #4d9d6d;
}
</style>
