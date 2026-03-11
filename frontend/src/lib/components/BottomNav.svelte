<script lang="ts">
	import { page } from '$app/stores';
	import { animate } from 'motion';
	import { onMount } from 'svelte';

	const tabs = [
		{ href: '/', label: 'Feed', icon: 'rss_feed' },
		{ href: '/trips', label: 'Trips', icon: 'map' },
		{ href: '/create', label: 'Create', icon: 'add', fab: true },
		{ href: '/explore', label: 'Explore', icon: 'explore' },
		{ href: '/profile', label: 'Profile', icon: 'person' }
	];

	let dot: HTMLSpanElement | null = null;
	let tabRefs: HTMLElement[] = [];

	const isActive = (href: string) => {
		if (href === '/') return $page.url.pathname === '/';
		return $page.url.pathname.startsWith(href);
	};

	const bounce = (event: MouseEvent) => {
		const target = event.currentTarget as HTMLElement | null;
		if (!target) return;
		animate(target, { scale: [1, 0.85, 1.12, 1] } as any, {
			duration: 0.35,
			easing: [0.34, 1.56, 0.64, 1]
		} as any);
	};

	const moveDot = () => {
		if (!dot) return;
		const index = tabs.findIndex((tab) => isActive(tab.href));
		const target = tabRefs[index];
		if (!target) return;
		const rect = target.getBoundingClientRect();
		const left = rect.left + rect.width / 2 - 2;
		animate(dot, { x: [`${dot.offsetLeft}px`, `${left}px`] } as any, {
			duration: 0.25,
			easing: [0.34, 1.56, 0.64, 1]
		} as any);
		dot.style.transform = `translateX(${left}px)`;
	};

	onMount(() => {
		requestAnimationFrame(moveDot);
	});

	$effect(() => {
		moveDot();
	});
</script>

<nav class="nav">
	<span class="dot" bind:this={dot}></span>
	{#each tabs as tab, i}
		<a
			class="tab"
			class:active={isActive(tab.href)}
			class:fab={tab.fab}
			href={tab.href}
			bind:this={tabRefs[i]}
			onclick={bounce}
		>
			<span class="icon material-symbols-outlined">{tab.icon}</span>
			{#if isActive(tab.href)}
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
		justify-content: space-around;
		align-items: center;
		height: var(--nav-height);
		padding: 0 1.25rem calc(env(safe-area-inset-bottom) / 2);
		background: rgba(22, 28, 24, 0.75);
		backdrop-filter: blur(24px) saturate(1.6);
		border-top: 1px solid rgba(255, 255, 255, 0.06);
		z-index: 100;
	}

	.tab {
		position: relative;
		width: 48px;
		height: 48px;
		display: grid;
		place-items: center;
		color: var(--text-muted);
		text-decoration: none;
	}

	.tab.active {
		color: var(--primary);
	}

	.icon {
		font-size: 1.3rem;
		line-height: 1;
	}

	.tab.active .material-symbols-outlined {
		font-variation-settings: "FILL" 1, "wght" 500, "GRAD" 0, "opsz" 24;
	}

	.label {
		position: absolute;
		bottom: -2px;
		font-size: 0.56rem;
		letter-spacing: 0.18em;
		text-transform: uppercase;
		color: var(--primary);
		font-weight: 700;
		opacity: 0.9;
	}

	.dot {
		position: absolute;
		bottom: 6px;
		left: 0;
		width: 4px;
		height: 4px;
		border-radius: 999px;
		background: var(--primary);
		box-shadow: 0 0 8px rgba(77, 157, 109, 0.6);
		transform: translateX(0);
	}
</style>
