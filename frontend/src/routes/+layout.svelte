<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import BottomNav from '$lib/components/BottomNav.svelte';
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { ensureAuth } from '$lib/auth';
	import { pageEnter } from '$lib/transitions';

	let { children } = $props();

	onMount(() => {
		ensureAuth($page.url.pathname);
	});
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<meta name="theme-color" content="#161c18" />
</svelte:head>

<div class="app">
	{#key $page.url.pathname}
		<div use:pageEnter>
			{@render children()}
		</div>
	{/key}
	{#if $page.url.pathname !== '/auth'}
		{#if !$page.url.pathname.startsWith('/create')}
			<BottomNav />
		{/if}
	{/if}
</div>

<style>
	.app {
		min-height: 100dvh;
		position: relative;
	}
</style>
