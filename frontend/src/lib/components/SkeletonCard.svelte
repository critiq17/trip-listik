<div class="skeleton"></div>

<style>
	.skeleton {
		width: 100%;
		aspect-ratio: 16 / 10;
		border-radius: 22px;
		background: linear-gradient(90deg, rgba(13, 52, 96, 0.7) 0%, rgba(28, 87, 133, 0.85) 40%, rgba(13, 52, 96, 0.7) 80%);
		background-size: 200% 100%;
		animation: shimmer 1.6s infinite;
		box-shadow: var(--shadow-card);
	}

	@keyframes shimmer {
		0% {
			background-position: -200% 0;
		}
		100% {
			background-position: 200% 0;
		}
	}
</style>
