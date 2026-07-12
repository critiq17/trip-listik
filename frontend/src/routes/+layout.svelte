<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import BottomNav from '$lib/components/BottomNav.svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { ensureAuth } from '$lib/auth';

	let { children } = $props();
	let authReady = $state(false);

	// Routes that hide the bottom nav
	const hideNavRoutes = ['/auth', '/create', '/trips/create'];

	const shouldShowNav = (pathname: string) =>
		!hideNavRoutes.some((r) => pathname === r || pathname.startsWith(r));

	onMount(() => {
		const tg = (window as any).Telegram?.WebApp;

		// ── CRITICAL: These MUST be synchronous, before any await ──
		if (tg) {
			tg.ready();                           // Tell Telegram the app is ready
			tg.expand();                          // Full-screen — synchronous, before any await
			tg.setHeaderColor?.('#ffffff');
			tg.setBackgroundColor?.('#ffffff');
			tg.setBottomBarColor?.('#ffffff');
		}

		// ── Deep links: t.me/...?startapp=<param> lands here as start_param ──
		// Supported: "notifications" / "inbox" → inbox, "profile_<id>" → profile,
		// "trip_<id>" → trip detail.
		const routeStartParam = () => {
			const raw: string =
				tg?.initDataUnsafe?.start_param ??
				new URLSearchParams(window.location.search).get('startapp') ??
				'';
			if (!raw) return;
			if (raw === 'notifications' || raw === 'inbox') {
				goto('/inbox');
				return;
			}
			const profileId = raw.startsWith('profile_') ? raw.slice('profile_'.length) : '';
			if (profileId) {
				goto(`/profile/${profileId}`);
				return;
			}
			const tripId = raw.startsWith('trip_') ? raw.slice('trip_'.length) : '';
			if (tripId) {
				goto(`/trips/${tripId}`);
			}
		};

		// ── Auth (async, after Telegram init) ──
		ensureAuth($page.url.pathname).then(() => {
			authReady = true;
			routeStartParam();
		});

		// ── Fallback: re-expand if minimized then re-opened ──
		const handleVisibilityChange = () => {
			const tgApp = (window as any).Telegram?.WebApp;
			if (document.visibilityState === 'visible' && tgApp && !tgApp.isExpanded) {
				tgApp.expand();
			}
		};
		document.addEventListener('visibilitychange', handleVisibilityChange);
		return () => document.removeEventListener('visibilitychange', handleVisibilityChange);
	});
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<meta name="theme-color" content="#ffffff" />
</svelte:head>

<div class="app">
	{#if !authReady}
		<div class="auth-loading" aria-hidden="true"></div>
	{:else}
		{@render children()}
		{#if shouldShowNav($page.url.pathname)}
			<BottomNav />
		{/if}
	{/if}
</div>

<style>
	.app {
		min-height: 100dvh;
		position: relative;
		background: var(--bg);
	}

	.auth-loading {
		min-height: 100dvh;
		background: var(--bg);
	}
</style>
