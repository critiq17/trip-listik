<script lang="ts">
	import { goto } from '$app/navigation';
	import { apiFetch } from '$lib/api';
	import { getDraft } from '$lib/tripDraft';

	let isPublic = true;
	let startDate = '';
	let endDate = '';
	let saving = false;
	let error = '';

	const next = async () => {
		const draft = getDraft();
		if (!draft) {
			error = 'Draft not found';
			return;
		}
		saving = true;
		error = '';
		try {
			await apiFetch(`/v1/trips/${draft.id}`, {
				method: 'PATCH',
				body: JSON.stringify({
					start_date: startDate || null,
					end_date: endDate || null,
					visibility: isPublic ? 'public' : 'private'
				})
			});
			goto('/create/step-3');
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to save';
		} finally {
			saving = false;
		}
	};
	const back = () => history.back();
</script>

<section class="container">
	<header class="header">
		<div class="copy">
			<p class="eyebrow">Create Trip</p>
			<h1 class="headline">Choose dates and privacy.</h1>
		</div>
		<button class="back" on:click={back}>Back</button>
	</header>

	<div class="progress">
		<div class="progress-top">
			<span>Step 2 of 3</span>
			<strong>66%</strong>
		</div>
		<div class="bar"><span></span></div>
	</div>

	<div class="form panel">
		<label>
			<span>Start date</span>
			<input type="date" bind:value={startDate} />
		</label>
		<label>
			<span>End date</span>
			<input type="date" bind:value={endDate} />
		</label>
		<div class="toggle">
			<span>Visibility</span>
			<div class="switch">
				<button class:active={isPublic} on:click={() => (isPublic = true)}>Public</button>
				<button class:active={!isPublic} on:click={() => (isPublic = false)}>Private</button>
			</div>
		</div>
	</div>
	{#if error}
		<div class="error">{error}</div>
	{/if}

	<footer class="footer">
		<button class="cta" on:click={next} disabled={saving}>
			{saving ? 'Saving…' : 'Next Step'}
		</button>
	</footer>
</section>

<style>
	.header {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		margin-bottom: 1rem;
		gap: 1rem;
	}

	.back {
		font-size: 0.8rem;
		color: var(--text-secondary);
	}

	.progress {
		margin-bottom: 1rem;
	}

	.progress-top {
		display: flex;
		justify-content: space-between;
		align-items: center;
		color: var(--text-secondary);
		font-size: 0.84rem;
		margin-bottom: 0.55rem;
	}

	.bar {
		height: 8px;
		border-radius: var(--radius-pill);
		background: rgba(255, 255, 255, 0.08);
		overflow: hidden;
	}

	.bar span {
		display: block;
		width: 66%;
		height: 100%;
		border-radius: inherit;
		background: var(--accent-grad);
	}

	.form {
		display: grid;
		gap: 1.25rem;
		padding: 1.2rem;
		border-radius: var(--radius-2xl);
	}

	label {
		display: grid;
		gap: 0.5rem;
		font-size: 0.85rem;
		color: var(--text-secondary);
	}

	input {
		background: rgba(255, 255, 255, 0.04);
		border: 1px solid rgba(255, 255, 255, 0.1);
		padding: 0.85rem 0.95rem;
		border-radius: var(--radius-lg);
	}

	.toggle {
		display: flex;
		justify-content: space-between;
		align-items: center;
		color: var(--text-secondary);
		font-size: 0.9rem;
	}

	.switch {
		display: flex;
		padding: 0.25rem;
		border-radius: var(--radius-pill);
		background: rgba(255, 255, 255, 0.05);
		border: 1px solid rgba(255, 255, 255, 0.08);
	}

	.toggle button {
		padding: 0.55rem 1rem;
		border-radius: var(--radius-pill);
		color: var(--text-secondary);
		font-weight: 700;
	}

	.toggle button.active {
		background: var(--accent-grad);
		color: white;
	}

	.footer {
		position: sticky;
		bottom: 0;
		margin-top: 2rem;
		padding-bottom: 5rem;
	}

	.cta {
		width: 100%;
		padding: 1rem;
		border-radius: var(--radius-xl);
		background: var(--accent-grad);
		font-weight: 800;
		color: #fff;
	}

	.error {
		margin-top: 1rem;
		color: var(--danger);
	}
</style>
