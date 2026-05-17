<script lang="ts">
	import Refresh from '$lib/assets/icons/Refresh.svg?component';
	import Check from '$lib/assets/icons/Check.svg?component';
	import Plus from '$lib/assets/icons/Plus.svg?component';
	import Minus from '$lib/assets/icons/Minus.svg?component';
	import { invalidateAll } from '$app/navigation';
	import FilteredList from '$lib/components/FilteredList.svelte';
	import type { Contact } from '$lib/data';
	import type { PageProps } from './$types';
	import { SvelteSet } from 'svelte/reactivity';
	import { enhance } from '$app/forms';
	const { data }: PageProps = $props();
	const checked: SvelteSet<number> = $state(new SvelteSet());
	function isMatches(query: string, data: Contact): boolean {
		return data.email.toLowerCase().includes(query) || data.name.toLowerCase().includes(query);
	}
</script>

<div class="vstack" style:align-items="center">
	<FilteredList data={Array.from(data.data)} {isMatches} pageSize={5}>
		{#snippet header(searchBar)}
			<div class="card entry" style:font-size="xx-large" style:margin-bottom="10pt">
				<div style:flex-grow="1"><b>Bind Contacts</b></div>
				{@render searchBar('Search name or email...')}
				<button class="round" onclick={invalidateAll}><Refresh /></button>
				<form
					method="POST"
					use:enhance={({ formData }) => {
						formData.append('contacts_id', JSON.stringify({ contacts_id: Array.from(checked) }));
					}}
				>
					<button class="round"><Check /></button>
				</form>
			</div>
		{/snippet}

		{#snippet entry(d)}
			<div class="entry card">
				<div class="contact">
					<div class="name"><b>{d.name}</b></div>
					<div class="email">{d.email}</div>
				</div>
				{#if !checked.has(d.id)}
					<button
						class="btn round"
						onclick={() => {
							checked.add(d.id);
						}}><Plus /></button
					>
				{:else}
					<button
						class="btn round"
						onclick={() => {
							checked.delete(d.id);
						}}><Minus /></button
					>
				{/if}
			</div>
		{/snippet}
	</FilteredList>
</div>

<style>
	.contact {
		flex-grow: 1;
		.name {
			font-size: x-large;
			color: var(--dark-text);
		}
		.email {
			font-size: small;
			color: var(--grey-text);
		}
	}
</style>
