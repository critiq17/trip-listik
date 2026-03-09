<script lang="ts">
	import { formatDateRange, getStatusLabel, getTripLocation } from '$lib/format';
	import type { TripCardData } from '$lib/types';
	export let trip: TripCardData;
</script>

<a class="card" href={`/trips/${trip.id}`}>
	{#if trip.cover_photo_url}
		<img src={trip.cover_photo_url} alt={trip.title} loading="lazy" />
	{:else}
		<div class="placeholder"></div>
	{/if}
	<div class="overlay"></div>
	<div class="status-row">
		<span class="status">{getStatusLabel(trip.status)}</span>
		{#if trip.country_code}
			<span class="country">{trip.country_code}</span>
		{/if}
	</div>
	<footer>
		<div class="info">
			<h3>{trip.title}</h3>
			<p>{formatDateRange(trip.start_date, trip.end_date)}</p>
			<p>{getTripLocation(trip)} · {trip.member_count ?? 0} members</p>
		</div>
		<span class="votes">👍 {trip.vote_average?.toFixed?.(1) ?? '0.0'}</span>
	</footer>
</a>

<style>
	.card {
		position: relative;
		border-radius: 24px;
		overflow: hidden;
		aspect-ratio: 16 / 10;
		display: block;
		box-shadow: var(--shadow-card);
		transform: translateZ(0);
	}

	img,
	.placeholder {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.placeholder {
		background:
			radial-gradient(circle at top right, rgba(122, 234, 244, 0.15), transparent 35%),
			linear-gradient(120deg, #0b2b4d, #0f3d6e);
	}

	.overlay {
		position: absolute;
		inset: 0;
		background: linear-gradient(to top, rgba(4, 20, 38, 0.98), transparent 62%);
	}

	.status-row {
		position: absolute;
		top: 1rem;
		left: 1rem;
		right: 1rem;
		display: flex;
		justify-content: space-between;
		align-items: center;
		gap: 0.5rem;
		z-index: 1;
	}

	.status,
	.country {
		padding: 0.45rem 0.75rem;
		border-radius: var(--radius-pill);
		backdrop-filter: blur(14px);
		background: rgba(4, 36, 68, 0.5);
		border: 1px solid rgba(255, 255, 255, 0.12);
		font-size: 0.72rem;
		font-weight: 700;
		letter-spacing: 0.08em;
		text-transform: uppercase;
	}

	.status {
		color: white;
	}

	footer {
		position: absolute;
		left: 0;
		right: 0;
		bottom: 0;
		padding: 1rem 1.25rem;
		display: flex;
		justify-content: space-between;
		align-items: flex-end;
		gap: 1rem;
		z-index: 1;
	}

	h3 {
		font-size: 1.5rem;
		font-weight: 800;
		color: white;
		letter-spacing: -0.03em;
	}

	p {
		font-size: 0.8rem;
		color: var(--text-secondary);
		margin-top: 0.22rem;
	}

	.votes {
		background: linear-gradient(135deg, rgba(32, 146, 186, 0.9), rgba(26, 122, 156, 0.95));
		border: 1px solid rgba(122, 234, 244, 0.2);
		color: white;
		padding: 0.45rem 0.85rem;
		border-radius: var(--radius-pill);
		font-weight: 800;
		font-size: 0.8rem;
		box-shadow: 0 10px 24px rgba(32, 146, 186, 0.25);
	}
</style>
