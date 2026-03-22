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

	const isCompact    = $derived(variant === 'compact');
	const isHorizontal = $derived(variant === 'horizontal');
	const coverUrl     = $derived(resolvePhotoUrl(trip.cover_photo_url));
</script>

<a class={`card ${variant}`} href={`/trips/${trip.id}`}>

	{#if isHorizontal}
		<!-- ── Horizontal layout ── -->
		<div class="thumb-wrap">
			{#if coverUrl}
				<div class="img-skeleton skeleton"></div>
				<img
					class="thumb-img"
					src={coverUrl}
					alt={trip.title}
					loading="lazy"
					onload={(e) => e.currentTarget.classList.add('loaded')}
					onerror={(e) => e.currentTarget.classList.add('error')}
				/>
			{:else}
				<div class="thumb-placeholder"></div>
			{/if}
		</div>
		<div class="h-body">
			<p class="h-title">{trip.title}</p>
			<p class="h-meta">{formatDateRange(trip.start_date, trip.end_date)}</p>
		</div>

	{:else}
		<!-- ── Vertical card (feed / compact) ── -->
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

			{#if isCompact && trip.status}
				<div class="badge">{getStatusLabel(trip.status)}</div>
			{/if}
		</div>

		<div class="bottom">
			<div class="info">
				<h3>{trip.title}</h3>
				<p>{formatDateRange(trip.start_date, trip.end_date)}</p>
			</div>

			{#if !isCompact && (trip.vote_count ?? 0) > 0}
				<div class="like-pill">
					<span class="material-symbols-outlined fill-icon">thumb_up</span>
					{trip.vote_count ?? 0}
				</div>
			{/if}
		</div>
	{/if}
</a>

<style>
/* ── Base card ─────────────────────────────────────────── */
.card {
	display: block;
	position: relative;
	border-radius: 12px;
	overflow: hidden;
	cursor: pointer;
	box-shadow: 0 8px 30px rgba(0, 0, 0, 0.4);
	transition: transform 0.25s ease, box-shadow 0.25s ease;
	-webkit-tap-highlight-color: transparent;
}

.card:hover {
	transform: translateY(-2px);
	box-shadow: 0 16px 40px rgba(0, 0, 0, 0.55);
}

.card:active {
	transform: scale(0.97);
	transition-duration: 0.08s;
}

/* ── Feed variant ──────────────────────────────────────── */
.card.feed .media {
	aspect-ratio: 4 / 5;
	position: relative;
	overflow: hidden;
}

/* ── Compact variant (Trips page) ──────────────────────── */
.card.compact .media {
	aspect-ratio: 16 / 10;
	position: relative;
	overflow: hidden;
}

/* ── Shared media elements ─────────────────────────────── */
.media,
.placeholder {
	width: 100%;
	height: 100%;
}

.img-skeleton {
	position: absolute;
	inset: 0;
}

.trip-img {
	position: relative;
	width: 100%;
	height: 100%;
	object-fit: cover;
	display: block;
	opacity: 0;
	transition: opacity 0.4s ease, transform 500ms ease;
}

.trip-img.loaded {
	opacity: 1;
}

.trip-img.error {
	display: none;
}

.card:hover .trip-img {
	transform: scale(1.05);
}

.placeholder {
	background: linear-gradient(135deg, rgba(255,255,255,0.04), rgba(255,255,255,0.02));
}

.overlay {
	position: absolute;
	inset: 0;
	background: linear-gradient(to top, rgba(0, 0, 0, 0.82) 0%, transparent 60%);
	pointer-events: none;
}

/* ── Badge (compact) ───────────────────────────────────── */
.badge {
	position: absolute;
	top: 16px; right: 16px;
	background: rgba(255, 255, 255, 0.15);
	backdrop-filter: blur(8px);
	-webkit-backdrop-filter: blur(8px);
	color: white;
	font-size: 10px;
	font-weight: 700;
	text-transform: uppercase;
	letter-spacing: 0.08em;
	padding: 4px 8px;
	border-radius: 4px;
}

/* ── Bottom overlay content ────────────────────────────── */
.bottom {
	position: absolute;
	bottom: 0; left: 0; right: 0;
	padding: 24px;
	display: flex;
	justify-content: space-between;
	align-items: flex-end;
	gap: 12px;
}

.info h3 {
	font-size: 22px;
	font-weight: 700;
	color: white;
	margin: 0 0 4px;
	line-height: 1.2;
}

.info p {
	font-size: 14px;
	color: rgba(255, 255, 255, 0.75);
	font-weight: 500;
	margin: 0;
}

.like-pill {
	display: flex;
	align-items: center;
	gap: 6px;
	background: rgba(255, 255, 255, 0.15);
	backdrop-filter: blur(8px);
	-webkit-backdrop-filter: blur(8px);
	border-radius: 9999px;
	padding: 6px 12px;
	color: white;
	font-size: 14px;
	font-weight: 700;
	white-space: nowrap;
	flex-shrink: 0;
}

.like-pill .material-symbols-outlined {
	font-size: 16px;
}

/* ── Horizontal variant (Profile page) ─────────────────── */
.card.horizontal {
	display: flex;
	align-items: center;
	gap: 16px;
	padding: 0;
	box-shadow: none;
	background: transparent;
	border-bottom: 1px solid rgba(255, 255, 255, 0.06);
	border-radius: 0;
	padding-bottom: 16px;
}

.card.horizontal:hover {
	transform: none;
	box-shadow: none;
}

.thumb-wrap {
	position: relative;
	width: 48px;
	height: 48px;
	border-radius: 6px;
	overflow: hidden;
	flex-shrink: 0;
}

.thumb-img {
	position: relative;
	width: 100%;
	height: 100%;
	object-fit: cover;
	display: block;
	opacity: 0;
	transition: opacity 0.3s ease;
}

.thumb-img.loaded {
	opacity: 1;
}

.thumb-img.error {
	display: none;
}

.thumb-placeholder {
	width: 100%;
	height: 100%;
	background: rgba(255, 255, 255, 0.06);
	border-radius: 6px;
}

.h-body {
	flex: 1;
	min-width: 0;
}

.h-title {
	font-size: 15px;
	font-weight: 600;
	color: white;
	line-height: 1.3;
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
}

.h-meta {
	font-size: 13px;
	color: #94a3b8;
	margin-top: 2px;
}
</style>
