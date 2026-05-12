<script>
	import { enhance } from '$app/forms';
	import Check from '$lib/assets/icons/Check.svg?component';
	import X from '$lib/assets/icons/X.svg?component';

	const { data } = $props();
	let name = $derived(data.name);
	let email = $derived(data.email);
	let changed = $state(false);
	const setChanged = () => (changed = true);
	const reset = () => {
		changed = false;
		name = data.name;
		email = data.email;
	};
</script>

<div class="center">
	<form
		class="card main vstack"
		method="POST"
		use:enhance={() =>
			({ update }) => {
				changed = false;
				update();
			}}
	>
		<input type="text" name="name" class="name-input" bind:value={name} oninput={setChanged} />
		<div>
			<label for="email" class="text-form">Email</label>
			<input
				type="email"
				name="email"
				class="text-form"
				style:width="100%"
				bind:value={email}
				oninput={setChanged}
			/>
		</div>
		{#if changed}
			<div class="control-buttons">
				<button type="submit" style:background-color="var(--success)"><Check /></button>
				<button style:background-color="var(--error)" onclick={reset}><X /></button>
			</div>
		{/if}
	</form>
</div>

<style>
	.name-input {
		max-width: 100%;
		font-size: x-large;
		background: inherit;
		border: none;
	}
	.name-input:focus {
		outline: none;
	}
	.main {
		row-gap: 10pt;
		border: solid 1.5pt var(--border);
		width: 200pt;
	}
	.control-buttons {
		display: flex;
		justify-content: space-between;
		button {
			border-radius: 100%;
			color: var(--light-text);
			aspect-ratio: 1;
			display: flex;
			align-items: center;
			justify-content: center;
		}
	}
</style>
