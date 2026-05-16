<script lang="ts">
	import Refresh from '$lib/assets/icons/Refresh.svg?component';
	import Link from '$lib/assets/icons/Link.svg?component';
	import { invalidateAll } from '$app/navigation';
	import EditTemplate from '$lib/components/EditTemplate.svelte';
	import FilteredList from '$lib/components/FilteredList.svelte';
	import type { Contact } from '$lib/data';
	import type { PageProps } from './$types';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	const { data }: PageProps = $props();
	function isMatches(query: string, data: Contact): boolean {
		return data.email.toLowerCase().includes(query) || data.name.toLowerCase().includes(query);
	}
</script>

<EditTemplate data={data.template}></EditTemplate>

<div class="vstack" style:align-items="center">
	<FilteredList data={data.contacts} {isMatches} pageSize={5}>
		{#snippet header(searchBar)}
			<div class="card entry" style:font-size="xx-large" style:margin-bottom="10pt">
				<div style:flex-grow="1"><b>Bound Contacts</b></div>
				{@render searchBar('Search name or email...')}
				<button class="round" onclick={invalidateAll}><Refresh /></button>
				<a class="btn round" href={resolve(`/templates/${page.params.id}/bind`)}><Link /></a>
			</div>
		{/snippet}

		{#snippet entry(d)}
			<div class="entry card">
				<div class="contact">
					<div class="name"><b>{d.name}</b></div>
					<div class="email">{d.email}</div>
				</div>
			</div>
		{/snippet}
	</FilteredList>
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
