<script lang="ts">
	import { goto } from '$app/navigation';
	import { setDraft, getDraft } from '$lib/tripDraft';
	import { onMount, onDestroy } from 'svelte';
	import { setupMainButton, hideMainButton, setupBackButton } from '$lib/telegram';
	import DateRangePicker from '$lib/components/DateRangePicker.svelte';

	let draft = getDraft() || {};
	let visibility = $state(draft.visibility || 'public');
	let startDate = $state(draft.start_date || '');
	let endDate = $state(draft.end_date || '');

	$effect(() => {
		if (startDate && endDate) {
			setupMainButton('Continue', next, true, true);
		} else {
			setupMainButton('Continue', next, true, false);
		}
	});

	onMount(() => {
		setupBackButton(() => history.back());
	});

	onDestroy(() => {
		hideMainButton();
	});

	const next = () => {
		if (!startDate || !endDate) return;
		setDraft({
			start_date: startDate,
			end_date: endDate,
			visibility: visibility
		});
		goto('/create/step-3');
	};

	const selectVisibility = (value: string) => {
		visibility = value;
	};
</script>

<svelte:head>
	<title>Create Trip - Step 2</title>
</svelte:head>

<section class="step-two">
	<div class="progress-bar">
		<div class="progress-fill" style="width: 66.66%"></div>
	</div>

	<main>
		<p class="eyebrow">Step 2 of 3</p>
		<h1 class="serif-text">When does the journey begin?</h1>

		<DateRangePicker bind:startDate bind:endDate />

		<div class="visibility-section">
			<div class="section-head">
				<strong>Trip Type</strong>
				<p>Who can join and see your itinerary?</p>
			</div>
			
			<div class="trip-types">
				<button 
					type="button"
					class="type-card" 
					class:active={visibility === 'public'} 
					onclick={() => selectVisibility('public')}
				>
					<span class="material-symbols-outlined icon">public</span>
					<div class="info">
						<strong>Public Trip</strong>
						<p>Anyone can view and join. Best for open invitations.</p>
					</div>
				</button>
				<button 
					type="button"
					class="type-card" 
					class:active={visibility === 'private'} 
					onclick={() => selectVisibility('private')}
				>
					<span class="material-symbols-outlined icon">lock</span>
					<div class="info">
						<strong>Private Trip</strong>
						<p>Only you and approved members can view. Invite only.</p>
					</div>
				</button>
				<button 
					type="button"
					class="type-card" 
					class:active={visibility === 'group'} 
					onclick={() => selectVisibility('group')}
				>
					<span class="material-symbols-outlined icon">group</span>
					<div class="info">
						<strong>Group Trip</strong>
						<p>Private trip for a group of friends. Members can invite others.</p>
					</div>
				</button>
				<button 
					type="button"
					class="type-card" 
					class:active={visibility === 'tour'} 
					onclick={() => selectVisibility('tour')}
				>
					<span class="material-symbols-outlined icon">tour</span>
					<div class="info">
						<strong>Public Tour</strong>
						<p>Publicly visible itinerary, but members cannot edit. Managed by organizer.</p>
					</div>
				</button>
			</div>
		</div>
	</main>
</section>

<style>
	.step-two {
		min-height: 100dvh;
		background: var(--bg);
		color: var(--text);
		padding-bottom: 3rem;
	}

	.progress-bar {
		width: 100%;
		height: 4px;
		background: var(--bg-elevated);
	}

	.progress-fill {
		height: 100%;
		background: var(--green);
		transition: width 0.3s ease;
	}

	main {
		max-width: 480px;
		margin: 0 auto;
		padding: 2rem 1.5rem;
	}

	.eyebrow {
		color: var(--green);
		font-size: 0.65rem;
		font-weight: 700;
		letter-spacing: 0.2em;
		text-transform: uppercase;
	}

	h1 {
		font-size: 2.2rem;
		line-height: 1.1;
		margin-bottom: 2rem;
	}

	.visibility-section {
		margin-top: 2rem;
		padding-top: 2rem;
		border-top: 1px solid var(--border);
	}

	.section-head {
		margin-bottom: 1.2rem;
	}

	.section-head strong {
		font-size: 1.05rem;
		display: block;
		margin-bottom: 0.2rem;
		color: var(--text);
	}

	.section-head p {
		color: var(--text-sub);
		font-size: 0.85rem;
	}

	.trip-types {
		display: flex;
		flex-direction: column;
		gap: 0.8rem;
	}

	.type-card {
		display: flex;
		align-items: center;
		gap: 1rem;
		padding: 1rem;
		border-radius: var(--radius-card);
		background: var(--bg-elevated);
		border: 1px solid var(--border);
		text-align: left;
		transition: all 0.2s;
		cursor: pointer;
		-webkit-tap-highlight-color: transparent;
	}

	.type-card:hover {
		border-color: rgba(61, 158, 95, 0.4);
	}

	.type-card .icon {
		color: var(--text-sub);
		background: var(--bg-card);
		padding: 0.6rem;
		border-radius: var(--radius-input);
	}

	.type-card .info {
		flex: 1;
	}

	.type-card .info strong {
		display: block;
		font-size: 0.95rem;
		color: var(--text);
		margin-bottom: 0.15rem;
	}

	.type-card .info p {
		font-size: 0.75rem;
		color: var(--text-sub);
		line-height: 1.3;
	}

	.type-card.active {
		background: rgba(61, 158, 95, 0.1);
		border-color: var(--green);
	}

	.type-card.active .icon {
		background: var(--green);
		color: var(--bg);
	}
</style>
