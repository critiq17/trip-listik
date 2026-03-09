<script lang="ts">
	import { goto } from '$app/navigation';
	import { apiFetch } from '$lib/api';
	import { clearDraft, getDraft } from '$lib/tripDraft';

	let description = '';
	let invites = '';
	let saving = false;
	let error = '';

	const finish = async () => {
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
					description,
					status: 'planned'
				})
			});
			clearDraft();
			goto('/trips');
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
			<h1 class="headline">Describe the vibe.</h1>
			<p class="subtle">Friend invite search is not exposed by backend yet, so this step keeps the UI and saves the real trip data.</p>
		</div>
		<button class="back" on:click={back}>Back</button>
	</header>

	<div class="progress">
		<div class="progress-top">
			<span>Step 3 of 3</span>
			<strong>100%</strong>
		</div>
		<div class="bar"><span></span></div>
	</div>

	<div class="form panel">
		<label>
			<span>Description</span>
			<textarea bind:value={description} placeholder="What makes this trip special?"></textarea>
		</label>
		<label>
			<span>Invite people</span>
			<input bind:value={invites} placeholder="@anna, @max, @kate" />
		</label>
		<div class="invite">
			<div class="chip">{invites || 'No invite search yet'}</div>
		</div>
	</div>
	{#if error}
		<div class="error">{error}</div>
	{/if}

	<footer class="footer">
		<button class="cta" on:click={finish} disabled={saving}>
			{saving ? 'Creating…' : 'Create Trip'}
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
		width: 100%;
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

	input,
	textarea {
		background: rgba(255, 255, 255, 0.04);
		border: 1px solid rgba(255, 255, 255, 0.1);
		padding: 0.8rem;
		border-radius: var(--radius-lg);
		min-height: 120px;
		resize: none;
	}

	input {
		min-height: 52px;
	}

	.invite {
		display: flex;
		gap: 0.6rem;
	}

	.chip {
		padding: 0.7rem 1rem;
		border-radius: var(--radius-xl);
		background: rgba(13, 52, 96, 0.6);
		color: var(--text-secondary);
		font-weight: 600;
		font-size: 0.8rem;
		width: 100%;
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
