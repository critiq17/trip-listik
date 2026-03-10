<script lang="ts">
	import { goto } from '$app/navigation';
	import { apiFetch } from '$lib/api';
	import { getDraft, clearDraft } from '$lib/tripDraft';
	import { scalePress } from '$lib/actions/animate';

	let description = '';
	let loading = false;
	let error = '';

	const submit = async () => {
		const draft = getDraft();
		if (!draft) {
			error = 'Draft not found';
			return;
		}
		loading = true;
		error = '';
		try {
			await apiFetch(`/v1/trips/${draft.id}`, {
				method: 'PATCH',
				body: JSON.stringify({
					description,
					status: 'planned'
				})
			});
			clearDraft();
			goto('/trips');
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to create trip';
		} finally {
			loading = false;
		}
	};

	const back = () => history.back();
</script>

<section class="step-three">
	<header class="header">
		<button class="round-btn" on:click={back} aria-label="Back">
			<span class="material-symbols-outlined">arrow_back</span>
		</button>
		<h2>Create Trip</h2>
		<div class="spacer"></div>
	</header>

	<div class="progress">
		<div class="progress-row">
			<p>Step 3 of 3</p>
			<span>100%</span>
		</div>
		<div class="progress-bar">
			<div class="progress-fill"></div>
		</div>
	</div>

	<main>
		<h1>Who's coming along?</h1>

		<div class="chips">
			<div class="chip">
				<div class="chip-avatar"></div>
				<span>Alex Rivera</span>
				<button>
					<span class="material-symbols-outlined">close</span>
				</button>
			</div>
			<div class="chip">
				<div class="chip-avatar"></div>
				<span>Sarah Chen</span>
				<button>
					<span class="material-symbols-outlined">close</span>
				</button>
			</div>
		</div>

		<div class="field">
			<input placeholder="Search friends..." />
			<span class="material-symbols-outlined">search</span>
		</div>

		<div class="field textarea">
			<label for="trip-description">Trip Description (Optional)</label>
			<textarea
				id="trip-description"
				rows="3"
				bind:value={description}
				placeholder="Tell everyone what the plan is..."
			></textarea>
		</div>

		<div class="actions">
			<button use:scalePress on:click={submit} disabled={loading}>
				{loading ? 'Creating…' : 'Create Trip'}
			</button>
			{#if error}
				<p class="error">{error}</p>
			{/if}
		</div>
	</main>
</section>

<style>
	.step-three {
		min-height: 100dvh;
		background: var(--background-dark);
		color: var(--text-primary);
	}

	.header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 0.9rem 1rem;
		border-bottom: 1px solid rgba(255, 255, 255, 0.08);
	}

	.header h2 {
		font-size: 1rem;
		font-weight: 700;
	}

	.round-btn {
		width: 2.6rem;
		height: 2.6rem;
		border-radius: 999px;
		display: grid;
		place-items: center;
	}

	.spacer {
		width: 2.6rem;
	}

	.progress {
		padding: 1.2rem 1.5rem 0;
	}

	.progress-row {
		display: flex;
		align-items: flex-end;
		justify-content: space-between;
		margin-bottom: 0.6rem;
	}

	.progress-row p {
		color: var(--primary);
		font-weight: 700;
		font-size: 0.65rem;
		letter-spacing: 0.2em;
		text-transform: uppercase;
	}

	.progress-row span {
		color: var(--primary);
		font-weight: 600;
		font-size: 0.8rem;
	}

	.progress-bar {
		height: 8px;
		background: rgba(255, 255, 255, 0.12);
		border-radius: 999px;
		overflow: hidden;
	}

	.progress-fill {
		height: 100%;
		width: 100%;
		background: var(--primary);
		border-radius: inherit;
	}

	main {
		padding: 1.5rem;
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}

	h1 {
		font-size: 1.4rem;
		font-weight: 700;
	}

	.chips {
		display: flex;
		flex-wrap: wrap;
		gap: 0.6rem;
	}

	.chip {
		display: inline-flex;
		align-items: center;
		gap: 0.4rem;
		padding: 0.3rem 0.6rem 0.3rem 0.3rem;
		border-radius: 999px;
		background: rgba(77, 157, 109, 0.2);
		border: 1px solid rgba(77, 157, 109, 0.3);
		font-size: 0.8rem;
	}

	.chip-avatar {
		width: 1.75rem;
		height: 1.75rem;
		border-radius: 999px;
		background: rgba(255, 255, 255, 0.2);
	}

	.chip button {
		display: grid;
		place-items: center;
	}

	.field {
		position: relative;
	}

	.field input {
		width: 100%;
		background: transparent;
		border: none;
		border-bottom: 2px solid rgba(255, 255, 255, 0.12);
		padding: 0.6rem 0;
		font-size: 1rem;
	}

	.field span {
		position: absolute;
		right: 0;
		top: 50%;
		transform: translateY(-50%);
		color: var(--text-secondary);
	}

	.field.textarea label {
		display: block;
		font-size: 0.8rem;
		color: var(--text-secondary);
		margin-bottom: 0.4rem;
	}

	.field.textarea textarea {
		width: 100%;
		background: transparent;
		border: none;
		border-bottom: 2px solid rgba(255, 255, 255, 0.12);
		padding: 0.4rem 0;
		font-size: 0.95rem;
		color: var(--text-primary);
		resize: none;
	}

	.actions button {
		width: 100%;
		padding: 1rem;
		border-radius: 12px;
		background: var(--primary);
		color: white;
		font-weight: 700;
		box-shadow: 0 14px 30px rgba(77, 157, 109, 0.35);
	}

	.error {
		margin-top: 0.75rem;
		color: #e11d48;
		font-size: 0.8rem;
		text-align: center;
	}
</style>
