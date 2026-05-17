<script lang="ts">
	import { enhance } from '$app/forms';
	import type { ResolvedPathname } from '$app/types';
	import Check from '$lib/assets/icons/Check.svg?component';
	import X from '$lib/assets/icons/X.svg?component';
	import type { Contact } from '$lib/data';

	interface Props {
		data: Contact;
		action?: ResolvedPathname;
	}

	const { data, action }: Props = $props();

	let name: string = $derived.by(() => {
		if (!changed) {
			return data.name;
		} else {
			return name;
		}
	});
	let email: string = $derived.by(() => {
		if (!changed) {
			return data.email;
		} else {
			return email;
		}
	});
	let changed = $state(false);
	const setChanged = () => (changed = true);
	const reset = () => {
		changed = false;
		name = data.name;
		email = data.email;
	};
</script>

<form
	class="card main vstack"
	method="POST"
	{action}
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
	<!-- {#if changed} -->
	<div class="control-buttons">
		<button type="submit" style:background-color="var(--success)"><Check /></button>
		<button style:background-color="var(--error)" onclick={reset}><X /></button>
	</div>
	<!-- {/if} -->
</form>

<style>
	.name-input {
		max-width: 100%;
		color: var(--dark-text);
		font-size: xx-large;
		font-weight: bold;
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
</style>
