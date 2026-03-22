<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { onMount, onDestroy } from 'svelte';
	import { apiFetch } from '$lib/api';
	import { unreadCount } from '$lib/stores/notifications';
	import { hapticImpact } from '$lib/telegram';

	const tabs = [
		{ id: 'feed',    href: '/',        icon: 'rss_feed',   label: 'Feed'    },
		{ id: 'trips',   href: '/trips',   icon: 'map',        label: 'Trips'   },
		{ id: 'create',  href: '/create',  icon: 'add_circle', label: 'Create', fab: true },
		{ id: 'explore', href: '/explore', icon: 'explore',    label: 'Explore' },
		{ id: 'profile', href: '/profile', icon: 'person',     label: 'Profile' },
	];

	const isActive = (href: string) => {
		if (href === '/') return $page.url.pathname === '/';
		return $page.url.pathname.startsWith(href);
	};

	const navigate = (e: MouseEvent, href: string) => {
		hapticImpact('light');
		goto(href);
	};

	// ── Notification badge polling ──────────────────────────────────────
	let pollInterval: ReturnType<typeof setInterval>;

	const pollNotifications = async () => {
		try {
			const data = await apiFetch<{ unread: number }>('/v1/inbox/unread');
			unreadCount.set(data.unread ?? 0);
		} catch {
			// non-critical — silently ignore
		}
	};

	onMount(() => {
		pollNotifications();
		pollInterval = setInterval(pollNotifications, 30_000);
	});

	onDestroy(() => {
		clearInterval(pollInterval);
	});
</script>

<nav class="bottom-nav" aria-label="Main navigation">
	{#each tabs as tab}
		<button
			id="nav-{tab.id}"
			class="nav-item"
			class:active={isActive(tab.href)}
			class:fab={tab.fab}
			onclick={(e) => navigate(e, tab.href)}
			aria-label={tab.label}
			aria-current={isActive(tab.href) ? 'page' : undefined}
		>
			<span
				class="material-symbols-outlined"
				class:fill-icon={isActive(tab.href) && !tab.fab}
			>{tab.icon}</span>
			<span class="nav-label">{tab.label}</span>
		</button>
	{/each}
</nav>

<style>
.bottom-nav {
	position: fixed;
	bottom: 0; left: 0; right: 0;
	z-index: 50;
	display: flex;
	justify-content: space-between;
	align-items: center;
	padding: 12px 24px env(safe-area-inset-bottom, 16px);
	background: rgba(22, 28, 24, 0.92);
	backdrop-filter: blur(20px);
	-webkit-backdrop-filter: blur(20px);
	border-top: 1px solid rgba(77, 157, 109, 0.1);
}

.nav-item {
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	gap: 2px;
	min-width: 44px;
	min-height: 44px;
	color: #64748b;
	background: none;
	border: none;
	cursor: pointer;
	border-radius: 8px;
	transition: color 0.2s ease;
	-webkit-tap-highlight-color: transparent;
}

.nav-item.active {
	color: #4d9d6d;
}

.nav-item .material-symbols-outlined {
	font-size: 24px;
}

.nav-item.fab .material-symbols-outlined {
	font-size: 32px;
}

.nav-item.fab.active {
	color: #4d9d6d;
}

.nav-label {
	font-size: 10px;
	font-weight: 700;
	text-transform: uppercase;
	letter-spacing: 0.08em;
	line-height: 1;
}
</style>
