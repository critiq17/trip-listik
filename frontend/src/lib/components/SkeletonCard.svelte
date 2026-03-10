<script lang="ts">
	export let ratio = '4 / 5';
</script>

<div class="skeleton-card skeleton" style={`aspect-ratio: ${ratio};`}></div>

<style>
	.skeleton-card {
		width: 100%;
		border-radius: 12px;
		box-shadow: var(--shadow-card);
	}
</style>
