<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import BottomNav from '$lib/components/BottomNav.svelte';
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { ensureAuth } from '$lib/auth';
	import { pageEnter } from '$lib/transitions';

	let { children } = $props();
	let authReady = $state(false);

	onMount(async () => {
		await ensureAuth($page.url.pathname);
		authReady = true;
	});
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<meta name="theme-color" content="#161c18" />
</svelte:head>

<div class="app">
	{#if !authReady}
		<div class="auth-loading"></div>
	{:else}
		{#key $page.url.pathname}
			<div use:pageEnter>
				{@render children()}
			</div>
		{/key}
		{#if $page.url.pathname !== '/auth' && !$page.url.pathname.startsWith('/create') && !$page.url.pathname.startsWith('/invite')}
			<BottomNav />
		{/if}
	{/if}
</div>

<style>
	.app {
		min-height: 100dvh;
		position: relative;
	}

	.auth-loading {
		min-height: 100dvh;
		background: var(--background-dark, #161c18);
	}
</style>
