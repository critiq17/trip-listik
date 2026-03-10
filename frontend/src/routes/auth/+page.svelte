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
	<div class="orb"></div>
	<div class="logo">TripListik</div>
	<p class="tag">Telegram Mini App</p>
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
		gap: 0.75rem;
		text-align: center;
		padding: 1.5rem;
	}

	.orb {
		width: 90px;
		height: 90px;
		border-radius: 50%;
		margin: 0 auto 0.2rem;
		background:
			radial-gradient(circle at 30% 30%, rgba(255, 255, 255, 0.9), transparent 25%),
			var(--primary);
		box-shadow: 0 12px 30px rgba(77, 157, 109, 0.45);
	}

	.logo {
		font-size: 2.6rem;
		font-weight: 800;
		letter-spacing: -0.05em;
	}

	.tag {
		color: var(--primary);
		font-size: 0.78rem;
		font-weight: 700;
		letter-spacing: 0.12em;
		text-transform: uppercase;
	}

	.error {
		color: #e11d48;
	}
</style>
