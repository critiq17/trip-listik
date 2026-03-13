<script lang="ts">
	import { resolvePhotoUrl } from '$lib/photos';
	import { formatDateRange, getStatusLabel } from '$lib/format';
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
	// BUG FIX: was $derived(() => ...) which wraps in a function. Must be $derived(expr).
	const coverUrl = $derived(resolvePhotoUrl(trip.cover_photo_url));
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
			<span class="status-pill">{getStatusLabel(trip.status)}</span>
		{/if}
	</div>

	<div class="content">
		{#if variant === 'feed'}
			<div class="bottom-row">
				<div>
					<h3>{trip.title}</h3>
					<p class="meta">{formatDateRange(trip.start_date, trip.end_date)}</p>
				</div>
				<div class="votes">
					<span class="material-symbols-outlined fill-1">thumb_up</span>
					<span>{trip.vote_count ?? 0}</span>
				</div>
			</div>
		{:else if isCompact}
			<div class="compact-body">
				<h3>{trip.title}</h3>
				<p class="meta">{formatDateRange(trip.start_date, trip.end_date)}</p>
				<span class="pill">{trip.member_count ?? 1} travelers</span>
			</div>
		{:else if isHorizontal}
			<div class="h-body">
				<h3>{trip.title}</h3>
				<p class="location">{trip.city || trip.country_code || 'Trip'}</p>
				<span class="pill subtle">{formatDateRange(trip.start_date, trip.end_date)}</span>
			</div>
		{/if}
	</div>
</a>

<style>
	.card {
		position: relative;
		border-radius: var(--radius-card);
		overflow: hidden;
		display: block;
		box-shadow: var(--shadow-card);
		background: var(--bg-card);
		border: 1px solid var(--border);
		transition:
			transform 0.25s var(--transition-smooth),
			box-shadow 0.25s var(--transition-smooth),
			border-color 0.25s var(--transition-smooth);
		-webkit-tap-highlight-color: transparent;
	}

	.card:hover {
		transform: translateY(-2px);
		box-shadow: 0 24px 48px rgba(0, 0, 0, 0.55);
		border-color: rgba(61, 158, 95, 0.2);
	}

	.card:active {
		transform: scale(0.97);
		transition-duration: 0.1s;
	}

	/* ── Media ── */
	.media {
		aspect-ratio: 16 / 9;
		position: relative;
		overflow: hidden;
	}

	.card.compact .media {
		aspect-ratio: 16 / 10;
	}

	.card.horizontal {
		display: grid;
		grid-template-columns: 88px 1fr;
		align-items: stretch;
		padding: 0.6rem;
		gap: 0;
	}

	.card.horizontal .media {
		aspect-ratio: 1 / 1;
		border-radius: 10px;
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

	.trip-img.loaded { opacity: 1; }
	.trip-img.error  { display: none; }

	.placeholder {
		background: linear-gradient(135deg, var(--bg), var(--bg-elevated));
	}

	.overlay {
		position: absolute;
		inset: 0;
		background: linear-gradient(to top, rgba(0, 0, 0, 0.7), rgba(0, 0, 0, 0.04) 55%);
	}

	.status-pill {
		position: absolute;
		right: 0.75rem;
		top: 0.75rem;
		padding: 0.25rem 0.6rem;
		border-radius: var(--radius-pill);
		background: rgba(13, 31, 23, 0.7);
		backdrop-filter: blur(12px);
		color: var(--green-light);
		font-size: 0.58rem;
		letter-spacing: 0.12em;
		text-transform: uppercase;
		font-weight: 700;
		border: 1px solid rgba(61, 158, 95, 0.3);
	}

	/* ── Content: feed variant ── */
	.content {
		position: absolute;
		left: 0; right: 0; bottom: 0;
		padding: 0.75rem 0.85rem;
		z-index: 1;
	}

	.card.horizontal .content {
		position: static;
		display: flex;
		align-items: center;
		padding: 0 0.75rem;
	}

	.card.compact .content {
		position: static;
	}

	.compact-body {
		padding: 0.75rem 0.85rem;
		display: flex;
		flex-direction: column;
		gap: 0.2rem;
	}

	.bottom-row {
		display: flex;
		justify-content: space-between;
		align-items: flex-end;
		gap: 0.75rem;
	}

	h3 {
		font-size: 1rem;
		font-weight: 700;
		color: white;
		letter-spacing: -0.01em;
		line-height: 1.2;
	}

	.meta {
		font-size: 0.68rem;
		color: rgba(255, 255, 255, 0.7);
		margin-top: 0.2rem;
	}

	.votes {
		display: inline-flex;
		flex-shrink: 0;
		align-items: center;
		gap: 0.25rem;
		padding: 0.2rem 0.5rem;
		border-radius: var(--radius-pill);
		background: rgba(255, 255, 255, 0.14);
		backdrop-filter: blur(10px);
		color: white;
		font-weight: 700;
		font-size: 0.68rem;
	}

	.votes .material-symbols-outlined {
		font-size: 0.9rem;
		color: var(--green-light);
	}

	.pill {
		display: inline-flex;
		align-items: center;
		padding: 0.2rem 0.5rem;
		border-radius: var(--radius-pill);
		background: rgba(61, 158, 95, 0.15);
		color: var(--green-light);
		font-weight: 700;
		font-size: 0.65rem;
		margin-top: 0.3rem;
	}

	.pill.subtle {
		background: rgba(255, 255, 255, 0.08);
		color: var(--text-sub);
	}

	/* ── Horizontal variant ── */
	.h-body {
		padding-left: 0.75rem;
		display: flex;
		flex-direction: column;
		gap: 0.15rem;
	}

	.h-body h3 {
		font-size: 0.95rem;
	}

	.location {
		font-size: 0.75rem;
		color: var(--text-sub);
	}

	.card.horizontal .meta {
		color: var(--text-sub);
	}
</style>
