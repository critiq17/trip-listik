<script lang="ts">
	let { ratio = '4 / 5' } = $props<{ ratio?: string }>();
</script>

<div class="skeleton-card skeleton" style={`aspect-ratio: ${ratio};`}></div>

<style>
	.skeleton-card {
		width: 100%;
		border-radius: var(--radius-card);
	}
</style>
