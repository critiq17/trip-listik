<script lang="ts">
	import { page } from '$app/stores';
	import { animate } from 'motion';

	const tabs = [
		{ href: '/', label: 'Feed', icon: 'rss_feed' },
		{ href: '/trips', label: 'Trips', icon: 'map' },
		{ href: '/create', label: 'Create', icon: 'add_circle', fab: true },
		{ href: '/explore', label: 'Explore', icon: 'explore' },
		{ href: '/profile', label: 'Profile', icon: 'person' }
	];

	const isActive = (href: string) => {
		if (href === '/') return $page.url.pathname === '/';
		return $page.url.pathname.startsWith(href);
	};

	const bounce = (event: MouseEvent) => {
		const target = event.currentTarget as HTMLElement | null;
		if (!target) return;
		animate(target, { y: ['0px', '-5px', '0px'] }, { duration: 0.28 });
	};
</script>

<nav class="nav">
	{#each tabs as tab}
		<a class:active={isActive(tab.href)} class:fab={tab.fab} href={tab.href} on:click={bounce}>
			{#if tab.fab}
				<span class="icon material-symbols-outlined fab-icon">{tab.icon}</span>
				<span class="label">{tab.label}</span>
			{:else}
				<span class="icon material-symbols-outlined">{tab.icon}</span>
				<span class="label">{tab.label}</span>
			{/if}
		</a>
	{/each}
</nav>

<style>
	.nav {
		position: fixed;
		left: 0;
		right: 0;
		bottom: 0;
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 0.75rem 1.5rem calc(0.9rem + env(safe-area-inset-bottom));
		background: rgba(22, 28, 24, 0.85);
		backdrop-filter: blur(18px);
		border-top: 1px solid rgba(77, 157, 109, 0.15);
		z-index: 100;
	}

	a {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.25rem;
		color: var(--text-muted);
		font-size: 0.62rem;
		font-weight: 700;
		text-transform: uppercase;
		letter-spacing: 0.12em;
		flex: 1;
	}

	a.active {
		color: var(--primary);
	}

	.fab {
		transform: translateY(-6px);
	}

	.icon {
		font-size: 1.2rem;
		line-height: 1;
	}

	.fab-icon {
		font-size: 2rem;
	}

	a.active .material-symbols-outlined {
		font-variation-settings: "FILL" 1, "wght" 500, "GRAD" 0, "opsz" 24;
	}

	.label {
		line-height: 1;
	}
</style>
