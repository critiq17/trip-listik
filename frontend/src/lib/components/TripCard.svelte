<script lang="ts">
	import { formatDateRange, getStatusLabel } from '$lib/format';
	import { resolvePhotoUrl } from '$lib/photos';
	import type { TripCardData } from '$lib/types';
	let {
		trip,
		variant = 'feed'
	}: {
		trip: TripCardData;
		variant?: 'feed' | 'compact' | 'horizontal';
	} = $props();

	const isCompact = $derived(variant === 'compact');
	const isHorizontal = $derived(variant === 'horizontal');
	const coverUrl = $derived(() => resolvePhotoUrl(trip.cover_photo_url));
</script>

<a class={`card ${variant}`} href={`/trips/${trip.id}`}>
	<div class="media">
		{#if coverUrl}
			<div class="img-skeleton skeleton"></div>
			<img
				class="trip-img"
				src={coverUrl}
				alt={trip.title}
				loading="lazy"
				onload={(e) => e.currentTarget.classList.add('loaded')}
				onerror={(e) => e.currentTarget.classList.add('error')}
			/>
		{:else}
			<div class="placeholder"></div>
		{/if}
		<div class="overlay"></div>
		{#if isCompact}
			<div class="status-pill">{getStatusLabel(trip.status)}</div>
		{/if}
	</div>
	<div class="content">
		<div>
			<h3>{trip.title}</h3>
			<p class="meta">{formatDateRange(trip.start_date, trip.end_date)}</p>
		</div>
		{#if variant === 'feed'}
			<div class="votes">
				<span class="material-symbols-outlined">thumb_up</span>
				<span>{trip.vote_count ?? 0}</span>
			</div>
		{:else if isCompact}
			<div class="pill">{trip.member_count ?? 1} travelers</div>
		{/if}
	</div>
	{#if isHorizontal}
		<div class="row">
			<p class="location">{trip.city || trip.country_code || 'Trip'}</p>
			<span class="pill subtle">{formatDateRange(trip.start_date, trip.end_date)}</span>
		</div>
	{/if}
</a>

<style>
	.card {
		position: relative;
		border-radius: var(--radius-2xl);
		overflow: hidden;
		display: block;
		box-shadow: 0 12px 26px rgba(0, 0, 0, 0.32);
		background: var(--card-bg);
		border: 1px solid rgba(255, 255, 255, 0.06);
		transition:
			transform 0.25s var(--transition-smooth),
			box-shadow 0.25s var(--transition-smooth),
			border-color 0.25s var(--transition-smooth);
	}

	.card:hover {
		transform: translateY(-2px);
		box-shadow: 0 18px 36px rgba(0, 0, 0, 0.4);
		border-color: rgba(255, 255, 255, 0.12);
	}

	.card:active {
		transform: scale(0.97);
		transition-duration: 0.1s;
	}

	.media {
		aspect-ratio: 3 / 4;
		position: relative;
		overflow: hidden;
	}

	.card.compact .media {
		aspect-ratio: 16 / 10;
	}

	.card.horizontal {
		display: grid;
		grid-template-columns: 96px 1fr;
		align-items: stretch;
		gap: 0.85rem;
		padding: 0.7rem;
	}

	.card.horizontal .media {
		aspect-ratio: 1 / 1;
		border-radius: var(--radius-xl);
	}

	img,
	.placeholder {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.img-skeleton {
		position: absolute;
		inset: 0;
	}

	.trip-img {
		position: relative;
		opacity: 0;
		transition: opacity 0.4s ease;
	}

	.trip-img.loaded {
		opacity: 1;
	}

	.trip-img.error {
		display: none;
	}

	.placeholder {
		background: linear-gradient(135deg, #0e1411, #1b2620);
	}

	.overlay {
		position: absolute;
		inset: 0;
		background: linear-gradient(to top, rgba(0, 0, 0, 0.62), rgba(0, 0, 0, 0.05) 60%);
	}

	.status-pill {
		position: absolute;
		right: 0.85rem;
		top: 0.85rem;
		padding: 0.3rem 0.65rem;
		border-radius: var(--radius-pill);
		background: rgba(255, 255, 255, 0.18);
		backdrop-filter: blur(12px);
		color: white;
		font-size: 0.6rem;
		letter-spacing: 0.12em;
		text-transform: uppercase;
		font-weight: 700;
	}

	h3 {
		font-size: 1.05rem;
		font-weight: 700;
		color: white;
		letter-spacing: -0.01em;
	}

	.meta {
		font-size: 0.7rem;
		color: rgba(255, 255, 255, 0.75);
		margin-top: 0.25rem;
	}

	.content {
		position: absolute;
		left: 0;
		right: 0;
		bottom: 0;
		padding: 0.85rem 0.9rem;
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
		padding: 0.25rem 0.55rem;
		border-radius: var(--radius-pill);
		background: rgba(255, 255, 255, 0.18);
		backdrop-filter: blur(10px);
		color: white;
		font-weight: 700;
		font-size: 0.7rem;
	}

	.votes .material-symbols-outlined {
		font-size: 1rem;
	}

	.pill {
		display: inline-flex;
		align-items: center;
		gap: 0.35rem;
		padding: 0.25rem 0.55rem;
		border-radius: var(--radius-pill);
		background: rgba(255, 255, 255, 0.18);
		backdrop-filter: blur(12px);
		color: white;
		font-weight: 700;
		font-size: 0.68rem;
	}

	.pill.subtle {
		background: rgba(255, 255, 255, 0.08);
		color: var(--text-secondary);
	}

	.card.horizontal .content {
		position: static;
		padding: 0.2rem 0.4rem 0;
		flex-direction: column;
		align-items: flex-start;
		gap: 0.4rem;
	}

	.card.horizontal h3 {
		font-size: 1rem;
	}

	.card.horizontal .meta {
		font-size: 0.7rem;
	}

	.row {
		grid-column: 2 / span 1;
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding-right: 0.4rem;
		padding-bottom: 0.4rem;
	}

	.location {
		font-size: 0.78rem;
		font-weight: 600;
		color: var(--text-primary);
	}
</style>
