<script lang="ts">
	import FilteredList from '$lib/components/FilteredList.svelte';
	import Minus from '$lib/assets/icons/Minus.svg?component';
	import Check from '$lib/assets/icons/Check.svg?component';
	import Plus from '$lib/assets/icons/Plus.svg?component';
	import Refresh from '$lib/assets/icons/Refresh.svg?component';
	import { SvelteSet } from 'svelte/reactivity';
	import type { Contact } from '$lib/data';
	import { invalidateAll } from '$app/navigation';
	import type { ResolvedPathname } from '$app/types';
	import { enhance } from '$app/forms';

	interface Props {
		contacts: Contact[];
		toBind: SvelteSet<number>;
		action: ResolvedPathname;
	}

	const { contacts, toBind = $bindable(), action }: Props = $props();

	function isMatches(query: string, data: Contact): boolean {
		return data.email.toLowerCase().includes(query) || data.name.toLowerCase().includes(query);
	}
</script>

<div class="center">
	<div class="vstack" style:align-items="center">
		<FilteredList data={contacts} {isMatches} pageSize={5}>
			{#snippet header(searchBar)}
				<div class="card entry" style:font-size="xx-large" style:margin-bottom="10pt">
					<div style:flex-grow="1"><b>Contacts</b></div>
					{@render searchBar('Search name or email...')}
					<button class="round" onclick={invalidateAll}><Refresh /></button>
					<form
						method="POST"
						{action}
						use:enhance={({ formData }) => {
							formData.append('contacts', JSON.stringify(Array.from(toBind)));
							return async ({ update }) => {
								await update();
							};
						}}
					>
						<button class="btn round"><Check /></button>
					</form>
				</div>
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
