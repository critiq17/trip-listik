<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { authenticate } from '$lib/auth';

	let status = 'Connecting…';
	let error = '';

	onMount(async () => {
		try {
			await authenticate();
			status = 'Ready';
			await goto('/');
		} catch (err) {
			error = err instanceof Error ? err.message : 'Auth failed';
		}
	});
</script>

<section class="auth">
	<div class="mark">
		<span class="material-symbols-outlined">travel_explore</span>
	</div>
	<div class="logo">TripListik</div>
	<p class="muted">{status}</p>
	{#if error}
		<p class="error">{error}</p>
	{/if}
</section>

<style>
	.auth {
		min-height: 100dvh;
		display: grid;
		place-content: center;
		gap: 8px;
		text-align: center;
		padding: 16px;
		background: var(--bg);
		color: var(--text);
	}

	.mark {
		width: 64px;
		height: 64px;
		border-radius: 16px;
		margin: 0 auto 8px;
		background: var(--green-soft);
		color: var(--green);
		display: grid;
		place-items: center;
	}

	.mark .material-symbols-outlined {
		font-size: 32px;
	}

	.logo {
		font-size: 28px;
		font-weight: 700;
		letter-spacing: -0.02em;
	}

	.muted {
		color: var(--text-sub);
		font-size: 14px;
	}

	.error {
		color: var(--danger);
		font-size: 14px;
	}
</style>
