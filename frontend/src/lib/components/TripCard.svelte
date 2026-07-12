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

			{#if isCompact && trip.status}
				<div class="badge">{getStatusLabel(trip.status)}</div>
			{/if}
		</div>

		<div class="bottom">
			<div class="info">
				<h3>{trip.title}</h3>
				<p>{formatDateRange(trip.start_date, trip.end_date)}</p>
			</div>

			{#if (trip.member_count ?? 0) > 0}
				<div class="count-pill">
					<span class="material-symbols-outlined">group</span>
					{trip.member_count}
				</div>
			{:else if !isCompact && (trip.vote_count ?? 0) > 0}
				<div class="count-pill">
					<span class="material-symbols-outlined">thumb_up</span>
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
	border-radius: var(--radius-card);
	overflow: hidden;
	cursor: pointer;
	background: var(--bg-card);
	border: 1px solid var(--border);
	-webkit-tap-highlight-color: transparent;
}

/* ── Feed variant ──────────────────────────────────────── */
.card.feed .media {
	aspect-ratio: 3 / 4;
	position: relative;
	overflow: hidden;
}

/* ── Compact variant (Feed page) ───────────────────────── */
.card.compact .media {
	aspect-ratio: 16 / 9;
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
	border-radius: 0;
}

.trip-img {
	position: relative;
	width: 100%;
	height: 100%;
	object-fit: cover;
	display: block;
	opacity: 0;
	transition: opacity 0.3s ease;
}

.trip-img.loaded {
	opacity: 1;
}

.trip-img.error {
	display: none;
}

.placeholder {
	background: var(--bg-subtle);
}

/* ── Badge (compact) ───────────────────────────────────── */
.badge {
	position: absolute;
	top: 12px; right: 12px;
	background: rgba(255, 255, 255, 0.92);
	color: var(--text);
	font-size: 11px;
	font-weight: 600;
	padding: 4px 8px;
	border-radius: 6px;
}

/* ── Bottom content ────────────────────────────────────── */
.bottom {
	padding: 12px 14px;
	display: flex;
	justify-content: space-between;
	align-items: center;
	gap: 12px;
}

.info h3 {
	font-size: 16px;
	font-weight: 600;
	color: var(--text);
	margin: 0 0 2px;
	line-height: 1.3;
}

.info p {
	font-size: 13px;
	color: var(--text-sub);
	margin: 0;
}

.count-pill {
	display: flex;
	align-items: center;
	gap: 5px;
	color: var(--text-sub);
	font-size: 13px;
	font-weight: 600;
	white-space: nowrap;
	flex-shrink: 0;
}

.count-pill .material-symbols-outlined {
	font-size: 18px;
}

/* ── Horizontal variant (Profile page) ─────────────────── */
.card.horizontal {
	display: flex;
	align-items: center;
	gap: 12px;
	padding: 12px 0;
	background: transparent;
	border: none;
	border-bottom: 1px solid var(--border);
	border-radius: 0;
}

.thumb-wrap {
	position: relative;
	width: 48px;
	height: 48px;
	border-radius: 8px;
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
	background: var(--bg-subtle);
	border-radius: 8px;
}

.h-body {
	flex: 1;
	min-width: 0;
}

.h-title {
	font-size: 15px;
	font-weight: 600;
	color: var(--text);
	line-height: 1.3;
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
}

.h-meta {
	font-size: 13px;
	color: var(--text-sub);
	margin-top: 2px;
}
</style>
