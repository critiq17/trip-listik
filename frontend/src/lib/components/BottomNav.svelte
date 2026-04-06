<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { onMount, onDestroy } from 'svelte';
	import { apiFetch } from '$lib/api';
	import { unreadCount } from '$lib/stores/notifications';
	import { hapticImpact } from '$lib/telegram';

	const tabs = [
		{ id: 'feed',    href: '/',        icon: 'home',          label: 'Home'    },
		{ id: 'create',  href: '/create',  icon: 'add',           label: 'Create', isCreate: true },
		{ id: 'inbox',   href: '/inbox',   icon: 'notifications', label: 'Inbox',  badge: true },
		{ id: 'profile', href: '/profile', icon: 'person',        label: 'Profile' },
	];

	const isActive = (href: string) => {
		const path = $page.url.pathname;
		if (href === '/') return path === '/';
		if (href === '/create') return path.startsWith('/create');
		return path.startsWith(href);
	};

	const navigate = (href: string) => {
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
			class:create={tab.isCreate}
			onclick={() => navigate(tab.href)}
			aria-label={tab.label}
			aria-current={isActive(tab.href) ? 'page' : undefined}
		>
			{#if tab.isCreate}
				<span class="create-icon">
					<span class="material-symbols-outlined">add</span>
				</span>
			{:else}
				<span class="icon-wrap">
					<span
						class="material-symbols-outlined"
						class:fill-icon={isActive(tab.href)}
					>{tab.icon}</span>
					{#if tab.badge && $unreadCount > 0}
						<span class="badge">{$unreadCount > 99 ? '99+' : $unreadCount}</span>
					{/if}
				</span>
				<span class="nav-label">{tab.label}</span>
			{/if}
		</button>
	{/each}
</nav>

<style>
.bottom-nav {
	position: fixed;
	bottom: 0; left: 0; right: 0;
	z-index: 50;
	display: flex;
	justify-content: space-around;
	align-items: center;
	padding: 10px 8px env(safe-area-inset-bottom, 12px);
	background: rgba(22, 28, 24, 0.96);
	backdrop-filter: blur(20px);
	-webkit-backdrop-filter: blur(20px);
	border-top: 1px solid rgba(77, 157, 109, 0.12);
}

.nav-item {
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	gap: 4px;
	flex: 1;
	min-height: 52px;
	color: #64748b;
	background: none;
	border: none;
	cursor: pointer;
	border-radius: 10px;
	transition: color 0.18s ease;
	-webkit-tap-highlight-color: transparent;
	padding: 6px 4px;
}

.nav-item.active {
	color: #4d9d6d;
}

.icon-wrap {
	position: relative;
	display: flex;
	align-items: center;
	justify-content: center;
}

.nav-item .material-symbols-outlined {
	font-size: 26px;
	transition: font-variation-settings 0.18s ease;
}

.nav-item.active .material-symbols-outlined {
	font-variation-settings: 'FILL' 1;
}

.badge {
	position: absolute;
	top: -5px;
	right: -8px;
	background: #ef4444;
	color: white;
	font-size: 9px;
	font-weight: 700;
	line-height: 1;
	padding: 2px 4px;
	border-radius: 9999px;
	min-width: 16px;
	text-align: center;
}

.nav-label {
	font-size: 10px;
	font-weight: 600;
	letter-spacing: 0.04em;
	line-height: 1;
}

.nav-item.create {
	color: white;
}

.create-icon {
	width: 44px;
	height: 44px;
	display: flex;
	align-items: center;
	justify-content: center;
	background: #4d9d6d;
	border-radius: 14px;
	box-shadow: 0 4px 16px rgba(77, 157, 109, 0.4);
}

.create-icon .material-symbols-outlined {
	font-size: 26px;
	color: white;
}
</style>
