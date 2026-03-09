<script lang="ts">
	import { page } from '$app/stores';

	const tabs = [
		{ href: '/', label: 'Feed', icon: '⌂' },
		{ href: '/trips', label: 'Trips', icon: '◫' },
		{ href: '/create', label: 'Create', icon: '+' , fab: true },
		{ href: '/explore', label: 'Explore', icon: '◎' },
		{ href: '/profile', label: 'Profile', icon: '◌' }
	];

	const isActive = (href: string) => {
		if (href === '/') return $page.url.pathname === '/';
		return $page.url.pathname.startsWith(href);
	};
</script>

<nav class="nav">
	{#each tabs as tab}
		<a class:active={isActive(tab.href)} class:fab={tab.fab} href={tab.href}>
			{#if tab.fab}
				<div class="fab-circle">
					<span class="icon">{tab.icon}</span>
				</div>
				<span class="label">{tab.label}</span>
			{:else}
				<span class="icon">{tab.icon}</span>
				<span class="label">{tab.label}</span>
			{/if}
		</a>
	{/each}
</nav>

<style>
	.nav {
		position: fixed;
		left: 0.75rem;
		right: 0.75rem;
		bottom: 0;
		display: flex;
		justify-content: space-around;
		align-items: flex-end;
		padding: 0.75rem 0.5rem calc(0.75rem + env(safe-area-inset-bottom));
		background: rgba(4, 36, 68, 0.82);
		backdrop-filter: blur(18px);
		border: 1px solid rgba(255, 255, 255, 0.08);
		border-bottom: none;
		border-radius: 24px 24px 0 0;
		box-shadow: 0 -16px 40px rgba(3, 18, 35, 0.4);
		z-index: 100;
	}

	a {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.2rem;
		color: var(--text-muted);
		font-size: 0.68rem;
		font-weight: 600;
		text-transform: uppercase;
		letter-spacing: 0.08em;
		flex: 1;
	}

	a.active {
		color: white;
	}

	.fab {
		transform: translateY(-18px);
	}

	.fab-circle {
		width: 56px;
		height: 56px;
		border-radius: 50%;
		background: var(--accent-grad);
		display: flex;
		align-items: center;
		justify-content: center;
		box-shadow: var(--shadow-glow);
	}

	.icon {
		font-size: 1.15rem;
		font-weight: 700;
		line-height: 1;
	}

	.label {
		line-height: 1;
	}

	a.active .icon {
		transform: translateY(-2px);
	}
</style>
