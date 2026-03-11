<script lang="ts">
	import { formatDateRange } from '$lib/format';
	import type { TripCardData } from '$lib/types';
	let { trip }: { trip: TripCardData } = $props();
</script>

<a class="card" href={`/trips/${trip.id}`}>
	<div class="media">
		{#if trip.cover_photo_url}
			<img src={trip.cover_photo_url} alt={trip.title} loading="lazy" />
		{:else}
			<div class="placeholder"></div>
		{/if}
		<div class="overlay"></div>
	</div>
	<div class="content">
		<div>
			<h3>{trip.title}</h3>
			<p>{formatDateRange(trip.start_date, trip.end_date)}</p>
		</div>
		<div class="votes">
			<span class="material-symbols-outlined">thumb_up</span>
			<span>{trip.vote_count ?? 0}</span>
		</div>
	</div>
</a>

<style>
	.card {
		position: relative;
		border-radius: 12px;
		overflow: hidden;
		display: block;
		box-shadow: var(--shadow-card);
		background: #0f1411;
	}

	.media {
		aspect-ratio: 4 / 5;
		position: relative;
		overflow: hidden;
	}

	img,
	.placeholder {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.placeholder {
		background: linear-gradient(135deg, #0e1411, #1b2620);
	}

	.overlay {
		position: absolute;
		inset: 0;
		background: linear-gradient(to top, rgba(0, 0, 0, 0.75), transparent 60%);
	}

	h3 {
		font-size: 1.4rem;
		font-weight: 700;
		color: white;
		letter-spacing: -0.01em;
	}

	p {
		font-size: 0.78rem;
		color: rgba(255, 255, 255, 0.75);
		margin-top: 0.25rem;
	}

	.content {
		position: absolute;
		left: 0;
		right: 0;
		bottom: 0;
		padding: 1.25rem;
		display: flex;
		justify-content: space-between;
		align-items: flex-end;
		gap: 1rem;
		z-index: 1;
	}

	.votes {
		display: inline-flex;
		align-items: center;
		gap: 0.3rem;
		padding: 0.35rem 0.7rem;
		border-radius: var(--radius-pill);
		background: rgba(255, 255, 255, 0.18);
		backdrop-filter: blur(10px);
		color: white;
		font-weight: 700;
		font-size: 0.78rem;
	}

	.votes .material-symbols-outlined {
		font-size: 1rem;
	}
</style>
