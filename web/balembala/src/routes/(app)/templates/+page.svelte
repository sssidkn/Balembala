<script lang="ts">
	import FilteredList from '$lib/components/FilteredList.svelte';
	import Trash from '$lib/assets/icons/Trash.svg?component';
	import Send from '$lib/assets/icons/Send.svg?component';
	import Plus from '$lib/assets/icons/Plus.svg?component';
	import Refresh from '$lib/assets/icons/Refresh.svg?component';
	import type { PageProps } from './$types';
	import type { Template } from '$lib/data';
	import { invalidateAll } from '$app/navigation';
	import { enhance } from '$app/forms';
	import { resolve } from '$app/paths';
	const props: PageProps = $props();
	let data = $derived(props.data.data);
	function isMatches(query: string, data: Template): boolean {
		return data.title.toLowerCase().includes(query);
	}
</script>

<div class="vstack" style:align-items="center">
	<FilteredList {data} {isMatches} pageSize={5}>
		{#snippet header(searchBar)}
			<div class="card entry" style:font-size="xx-large" style:margin-bottom="10pt">
				<div style:flex-grow="1"><b>Templates</b></div>
				{@render searchBar('Search title...')}
				<button class="round" onclick={invalidateAll}><Refresh /></button>
				<a class="btn round" href={resolve('/templates/new')}><Plus /></a>
			</div>
		{/snippet}

		{#snippet entry(d)}
			<form
				class="card entry"
				method="POST"
				use:enhance={({ formData, action }) => {
					if (action.search === '?/delete') {
						data = data.filter((d) => d.id != Number(formData.get('id')));
					}
					return async ({ update }) => {
						await update({ invalidateAll: false });
					};
				}}
			>
				<input value={d.id} type="hidden" name="id" />
				<a class="contact" href={resolve(`/templates/${d.id}`)}>
					<div class="title"><b>{d.title}</b></div>
				</a>
				<button class="round" formaction="?/send"><Send /></button>
				<button class="round" formaction="?/delete"><Trash /></button>
			</form>
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
		.title {
			font-size: x-large;
			color: var(--dark-text);
		}
	}
</style>
