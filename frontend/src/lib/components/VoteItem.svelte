<script lang="ts">
	import { animate } from 'motion';
	import { hapticImpact } from '$lib/telegram';
	import { flipNumber } from '$lib/transitions';
	import { resolvePhotoUrl } from '$lib/photos';
	import { getUserName } from '$lib/format';
	import { apiFetch } from '$lib/api';

	/**
	 * VoteItem represents an option inside a trip category (e.g., a specific hotel)
	 * that travelers can vote up or down.
	 */
	type ItemProps = {
		id: string;
		trip_id: string;
		category: string;
		title: string;
		description: string;
		image_url: string;
		added_by: string; // user ID
		vote_count: number;
		has_voted: boolean;
	};

	let {
		item,
		onVoteChange
	}: {
		item: ItemProps;
		onVoteChange: (itemId: string, newCount: number, hasVoted: boolean) => void;
	} = $props();

	let voting = $state(false);

	let coverUrl = $derived(resolvePhotoUrl(item.image_url));

	const handleVote = async (e: MouseEvent) => {
		const target = e.currentTarget as HTMLElement;
		if (voting) return;

		hapticImpact('light');
		
		// Button scale animation
		animate(target, { scale: [1, 1.15, 0.95, 1] } as any, { duration: 0.3 } as any);

		voting = true;
		const method = item.has_voted ? 'DELETE' : 'POST';
		
		// Optimistic update
		const delta = item.has_voted ? -1 : 1;
		const newCount = Math.max(0, item.vote_count + delta);
		onVoteChange(item.id, newCount, !item.has_voted);

		try {
			const res = await apiFetch<{ item_id: string; vote_count: number }>(
				`/v1/trips/${item.trip_id}/vote-items/${item.id}/vote`,
				{ method }
			);
			// Sync with true server count just in case
			onVoteChange(item.id, res.vote_count, !item.has_voted);
		} catch (err) {
			console.error('Failed to vote:', err);
			// Rollback optimism
			onVoteChange(item.id, item.vote_count, item.has_voted);
		} finally {
			voting = false;
		}
	};
</script>

<article class="vote-item">
	{#if coverUrl}
		<div class="media">
			<img src={coverUrl} alt={item.title} />
		</div>
	{/if}
	
	<div class="body">
		<div class="top">
			<div class="info">
				<h4>{item.title}</h4>
				{#if item.description}
					<p>{item.description}</p>
				{/if}
			</div>

			<button 
				class="vote-btn" 
				class:voted={item.has_voted} 
				onclick={handleVote}
				disabled={voting}
				aria-label={item.has_voted ? "Remove vote" : "Add vote"}
			>
				<div class="vote-inner">
					<span class="material-symbols-outlined" class:fill-1={item.has_voted}>
						thumb_up
					</span>
					<span class="count-wrap">
						<span 
							class="count" 
							use:flipNumber={item.vote_count}
						>{item.vote_count}</span>
					</span>
				</div>
			</button>
		</div>
	</div>
</article>

<style>
	.vote-item {
		background: var(--bg-card);
		border: 1px solid var(--border);
		border-radius: var(--radius-xl);
		overflow: hidden;
		display: flex;
		flex-direction: column;
		box-shadow: var(--shadow-soft);
	}

	.media {
		aspect-ratio: 16 / 9;
		overflow: hidden;
		background: var(--bg-elevated);
	}

	.media img {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.body {
		padding: 1rem;
	}

	.top {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		gap: 1rem;
	}

	.info h4 {
		font-size: 1rem;
		font-weight: 700;
		color: var(--text);
		margin-bottom: 0.2rem;
		line-height: 1.3;
	}

	.info p {
		font-size: 0.82rem;
		color: var(--text-sub);
		line-height: 1.4;
	}

	.vote-btn {
		flex-shrink: 0;
		border: 1px solid var(--border);
		background: var(--bg-elevated);
		border-radius: var(--radius-pill);
		padding: 0.4rem 0.75rem;
		color: var(--text-sub);
		transition: all 0.2s;
		-webkit-tap-highlight-color: transparent;
	}

	.vote-inner {
		display: flex;
		align-items: center;
		gap: 0.4rem;
	}

	.vote-btn.voted {
		background: rgba(61, 158, 95, 0.15);
		border-color: rgba(61, 158, 95, 0.3);
		color: var(--green);
	}

	.material-symbols-outlined {
		font-size: 1.1rem;
		line-height: 1;
	}

	.count-wrap {
		position: relative;
		overflow: hidden;
		height: 1em;
		display: inline-flex;
		align-items: center;
	}

	.count {
		font-size: 0.9rem;
		font-weight: 700;
		line-height: 1;
		display: inline-block;
	}

	/* ── CSS Animation for FlipNumber ── */
	@keyframes -global-flipIn {
		0% {
			transform: translateY(-100%);
			opacity: 0;
		}
		100% {
			transform: translateY(0);
			opacity: 1;
		}
	}
	
	:global(.count[data-prev]::before) {
		content: attr(data-prev);
		position: absolute;
		top: 100%;
		left: 0;
		opacity: 0.5;
	}
</style>
