<script lang="ts">
	import { enhance } from '$app/forms';
	import type { ResolvedPathname } from '$app/types';
	import Check from '$lib/assets/icons/Check.svg?component';
	import X from '$lib/assets/icons/X.svg?component';
	import type { Template } from '$lib/data';

	interface Props {
		data: Template;
		action?: ResolvedPathname;
	}

	const { data = $bindable(), action }: Props = $props();
	let title: string = $derived.by(() => {
		if (!changed) {
			return data.title;
		} else {
			return title;
		}
	});
	let message: string = $derived.by(() => {
		if (!changed) {
			return data.message;
		} else {
			return message;
		}
	});
	let changed = $state(false);
	const setChanged = () => (changed = true);
	const reset = () => {
		changed = false;
		title = data.title;
		message = data.message;
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
	<input type="hidden" value={data.id} />
	<input type="text" name="title" class="name-input" bind:value={title} oninput={setChanged} />
	<div>
		<label for="message" class="text-form">Message</label>
		<textarea
			rows="10"
			name="message"
			class="text-form"
			style:width="100%"
			bind:value={message}
			oninput={setChanged}
		></textarea>
	</div>
	{#if changed}
		<div class="control-buttons">
			<button type="submit" style:background-color="var(--success)"><Check /></button>
			<button style:background-color="var(--error)" onclick={reset}><X /></button>
		</div>
	{/if}
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
	textarea {
		box-sizing: border-box;
		resize: none;
	}
</style>
