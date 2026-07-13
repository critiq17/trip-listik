<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import BottomNav from '$lib/components/BottomNav.svelte';
	import { page, updated } from '$app/stores';
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
			tg.disableVerticalSwipes?.();         // Swipe on content scrolls; app closes only via header
			tg.setHeaderColor?.('#ffffff');
			tg.setBackgroundColor?.('#ffffff');
			tg.setBottomBarColor?.('#ffffff');
		}

		// iOS pinch bypasses the viewport meta in some WebView builds
		const preventGesture = (e: Event) => e.preventDefault();
		document.addEventListener('gesturestart', preventGesture);

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

		// ── Fallback: re-expand and re-apply colors if minimized then re-opened ──
		// Telegram restores a suspended WebView with its own theme colors, so a
		// dark-theme user gets a black background until we paint it white again.
		const handleVisibilityChange = () => {
			const tgApp = (window as any).Telegram?.WebApp;
			if (document.visibilityState === 'visible' && tgApp) {
				if (!tgApp.isExpanded) tgApp.expand();
				tgApp.disableVerticalSwipes?.();
				tgApp.setHeaderColor?.('#ffffff');
				tgApp.setBackgroundColor?.('#ffffff');
				tgApp.setBottomBarColor?.('#ffffff');
				// A restored session may predate the latest deploy; reload it
				// instead of letting stale JS mix with fresh assets.
				updated.check().then((stale) => {
					if (stale) location.reload();
				});
			}
		};
		document.addEventListener('visibilitychange', handleVisibilityChange);
		return () => {
			document.removeEventListener('visibilitychange', handleVisibilityChange);
			document.removeEventListener('gesturestart', preventGesture);
		};
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
