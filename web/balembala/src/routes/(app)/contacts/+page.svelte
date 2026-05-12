<script lang="ts">
	import FilteredList from '$lib/components/FilteredList.svelte';
	import Trash from '$lib/assets/icons/Trash.svg?component';
	import Plus from '$lib/assets/icons/Plus.svg?component';
	import Refresh from '$lib/assets/icons/Refresh.svg?component';
	import type { PageProps } from './$types';
	import type { Contact } from './+page.server';
	import { invalidateAll } from '$app/navigation';
	import { enhance } from '$app/forms';
	import { resolve } from '$app/paths';
	const props: PageProps = $props();
	let data = $derived(props.data.data);
	function isMatches(query: string, data: Contact): boolean {
		return data.email.toLowerCase().includes(query) || data.name.toLowerCase().includes(query);
	}
</script>

<div class="center">
	<div class="vstack" style:align-items="center">
		<FilteredList {data} {isMatches} pageSize={5}>
			{#snippet header(searchBar)}
				<form class="card entry" style:font-size="xx-large" style:margin-bottom="10pt">
					<div style:flex-grow="1"><b>Contacts</b></div>
					{@render searchBar('Search name or email...')}
					<button class="round" onclick={invalidateAll}><Refresh /></button>
					<a class="btn round" href={resolve('/contacts/new')}><Plus /></a>
				</form>
			{/snippet}

			{#snippet entry(d)}
				<form
					class="card entry"
					method="POST"
					action="?/delete"
					use:enhance={({ formData }) => {
						data = data.filter((d) => d.id != formData.get('id')?.toString());
						return ({ update }) => {
							update({ invalidateAll: false });
						};
					}}
				>
					<input value={d.id} type="hidden" name="id" />
					<a class="contact" href={resolve(`/contacts/${d.id}`)}>
						<div class="name"><b>{d.name}</b></div>
						<div class="email">{d.email}</div>
					</a>
					<button class="round"><Trash /></button>
				</form>
			{/snippet}
		</FilteredList>
	</div>
</div>

<style>
	a {
		text-decoration: none;
		color: inherit;
	}
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
