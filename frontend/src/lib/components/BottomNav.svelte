<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { hapticImpact } from '$lib/telegram';
	import { unreadCount } from '$lib/stores/notifications';

	const tabs = [
		{ id: 'home',    href: '/',        icon: 'home' },
		{ id: 'explore', href: '/explore', icon: 'search' },
		{ id: 'create',  href: '/create',  icon: 'add_box' },
		{ id: 'inbox',   href: '/inbox',   icon: 'notifications' },
		{ id: 'profile', href: '/profile', icon: 'person' }
	];

	const isActive = (href: string) => {
		const path = $page.url.pathname;
		if (href === '/') return path === '/';
		return path.startsWith(href);
	};

	const navigate = (href: string) => {
		hapticImpact('light');
		goto(href);
	};
</script>

<nav class="bottom-nav" aria-label="Main navigation">
	{#each tabs as tab}
		<button
			id="nav-{tab.id}"
			class="nav-item"
			class:active={isActive(tab.href)}
			onclick={() => navigate(tab.href)}
			aria-label={tab.id}
			aria-current={isActive(tab.href) ? 'page' : undefined}
		>
			<span class="icon-wrap">
				<span
					class="material-symbols-outlined"
					class:fill-icon={isActive(tab.href)}
				>{tab.icon}</span>
				{#if tab.id === 'inbox' && $unreadCount > 0}
					<span class="badge">{$unreadCount > 9 ? '9+' : $unreadCount}</span>
				{/if}
			</span>
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
	padding: 6px 8px env(safe-area-inset-bottom, 8px);
	background: var(--bg);
	border-top: 1px solid var(--border);
}

.nav-item {
	display: flex;
	align-items: center;
	justify-content: center;
	flex: 1;
	min-height: 48px;
	color: var(--text-sub);
	background: none;
	border: none;
	cursor: pointer;
	-webkit-tap-highlight-color: transparent;
	padding: 6px 4px;
	transition: color 0.15s ease;
}

.nav-item.active {
	color: var(--text);
}

.icon-wrap {
	position: relative;
	display: flex;
	align-items: center;
	justify-content: center;
}

.nav-item .material-symbols-outlined {
	font-size: 26px;
}

.badge {
	position: absolute;
	top: -4px;
	right: -8px;
	min-width: 16px;
	height: 16px;
	padding: 0 4px;
	border-radius: var(--radius-pill);
	background: var(--green);
	color: #fff;
	font-size: 10px;
	font-weight: 700;
	display: flex;
	align-items: center;
	justify-content: center;
}
</style>
