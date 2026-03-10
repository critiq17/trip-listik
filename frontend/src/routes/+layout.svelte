<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import BottomNav from '$lib/components/BottomNav.svelte';
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { ensureAuth } from '$lib/auth';

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
	{@render children()}
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
