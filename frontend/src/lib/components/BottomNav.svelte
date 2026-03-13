<script lang="ts">
	import { page } from '$app/stores';
	import { animate } from 'motion';
	import { onMount, onDestroy } from 'svelte';
	import { apiFetch } from '$lib/api';
	import { unreadCount } from '$lib/stores/notifications';
	import { goto } from '$app/navigation';
	import { hapticImpact } from '$lib/telegram';

	// ── 5-Tab layout: Feed | Notifications | + | Trips | Profile ──────────────
	const tabs = [
		{ href: '/',            label: 'Feed',     icon: 'home',          id: 'nav-feed' },
		{ href: '/inbox',       label: 'Inbox',    icon: 'notifications', id: 'nav-inbox', badge: true },
		{ href: '/create',      label: '',         icon: 'add',           id: 'nav-create', fab: true },
		{ href: '/trips',       label: 'Trips',    icon: 'map',           id: 'nav-trips' },
		{ href: '/profile',     label: 'Profile',  icon: 'person',        id: 'nav-profile' },
	];

	let tabRefs: HTMLElement[] = [];

	const isActive = (href: string) => {
		if (href === '/') return $page.url.pathname === '/';
		return $page.url.pathname.startsWith(href);
	};

	const bounce = (event: MouseEvent) => {
		const target = event.currentTarget as HTMLElement | null;
		if (!target) return;
		hapticImpact('light');
		animate(target, { scale: [1, 0.82, 1.1, 1] } as any, {
			duration: 0.3,
			easing: [0.34, 1.56, 0.64, 1]
		} as any);
	};

	const navigate = (event: MouseEvent, href: string) => {
		bounce(event);
		goto(href);
	};

	// ── Notification badge polling ─────────────────────────────────────────────
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
		pollNotifications(); // initial poll
		pollInterval = setInterval(pollNotifications, 30_000); // every 30s
	});

	onDestroy(() => {
		clearInterval(pollInterval);
	});
</script>

<nav class="nav" aria-label="Main navigation">
	{#each tabs as tab, i}
		<button
			id={tab.id}
			class="tab"
			class:active={isActive(tab.href)}
			class:fab={tab.fab}
			onclick={(e) => navigate(e, tab.href)}
			aria-label={tab.label || 'Create'}
			aria-current={isActive(tab.href) ? 'page' : undefined}
			bind:this={tabRefs[i]}
		>
			{#if tab.fab}
				<span class="fab-inner">
					<span class="material-symbols-outlined">add</span>
				</span>
			{:else}
				<span class="icon-wrap">
					<span class="material-symbols-outlined" class:filled={isActive(tab.href)}>
						{tab.icon}
					</span>
					{#if tab.badge && $unreadCount > 0}
						<span class="badge" aria-label="{$unreadCount} unread">
							{$unreadCount > 9 ? '9+' : $unreadCount}
						</span>
					{/if}
				</span>
				{#if isActive(tab.href) && tab.label}
					<span class="label">{tab.label}</span>
				{/if}
			{/if}
		</button>
	{/each}
</nav>

<style>
	.nav {
		position: fixed;
		left: 0;
		right: 0;
		bottom: 0;
		display: flex;
		justify-content: space-around;
		align-items: center;
		height: calc(var(--nav-height) + env(safe-area-inset-bottom));
		padding: 0 1rem calc(env(safe-area-inset-bottom) + 4px);
		background: rgba(13, 31, 23, 0.88);
		backdrop-filter: blur(24px) saturate(1.6);
		border-top: 1px solid var(--border);
		z-index: 100;
	}

	.tab {
		position: relative;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 2px;
		width: 52px;
		height: 52px;
		color: var(--text-muted);
		background: none;
		border: none;
		cursor: pointer;
		border-radius: 16px;
		transition: color 0.2s;
		-webkit-tap-highlight-color: transparent;
	}

	.tab.active {
		color: var(--green);
	}

	/* ── Center FAB ── */
	.tab.fab {
		color: var(--bg);
	}

	.fab-inner {
		width: 44px;
		height: 44px;
		border-radius: 14px;
		background: var(--green);
		display: grid;
		place-items: center;
		box-shadow: 0 4px 16px rgba(61, 158, 95, 0.5);
		transition: transform 0.2s var(--transition-spring), box-shadow 0.2s;
	}

	.tab.fab:active .fab-inner {
		transform: scale(0.92);
		box-shadow: 0 2px 8px rgba(61, 158, 95, 0.4);
	}

	/* ── Icon ── */
	.icon-wrap {
		position: relative;
		display: grid;
		place-items: center;
	}

	.material-symbols-outlined {
		font-size: 1.35rem;
		line-height: 1;
		font-variation-settings: 'FILL' 0, 'wght' 400, 'GRAD' 0, 'opsz' 24;
		transition: font-variation-settings 0.2s;
	}

	.material-symbols-outlined.filled {
		font-variation-settings: 'FILL' 1, 'wght' 500, 'GRAD' 0, 'opsz' 24;
	}

	/* ── Label ── */
	.label {
		font-size: 0.55rem;
		letter-spacing: 0.16em;
		text-transform: uppercase;
		font-weight: 700;
		color: var(--green);
		line-height: 1;
	}

	/* ── Unread badge ── */
	.badge {
		position: absolute;
		top: -4px;
		right: -4px;
		min-width: 16px;
		height: 16px;
		padding: 0 4px;
		border-radius: 999px;
		background: #e05555;
		color: white;
		font-size: 0.6rem;
		font-weight: 800;
		display: grid;
		place-items: center;
		line-height: 1;
		border: 2px solid var(--bg);
		animation: badgePop 0.3s var(--transition-spring) both;
	}

	@keyframes badgePop {
		from { transform: scale(0); }
		to   { transform: scale(1); }
	}
</style>
