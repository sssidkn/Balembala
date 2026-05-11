<script lang="ts" generics="Data extends { id: number, name: string }">
	import type { Snippet } from 'svelte';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { SvelteURLSearchParams } from 'svelte/reactivity';

	interface Props {
		data: Data[];
		pageSize: number;
		header: Snippet<[Snippet<[]>]>;
		entry: Snippet<[Data]>;
	}

	const { data, pageSize, header, entry }: Props = $props();

	let filters = $derived({
		page: page.url.searchParams.get('page') ?? '1',
		search: page.url.searchParams.get('search') ?? ''
	});
	let pageNumber = $derived(Number(filters.page));
	let filteredData = $derived.by(() => {
		const search = filters.search.toLowerCase();
		const filtered = data.filter((d) => d.name.toLowerCase().includes(search));
		const start = (pageNumber - 1) * pageSize;
		return filtered.slice(start, start + pageSize);
	});
	let lastPage = $derived(Math.ceil(data.length / pageSize));

	function updatePage(f: (old: number) => number) {
		const newParams = new SvelteURLSearchParams(filters);
		newParams.set('page', f(pageNumber).toString());
		// eslint-disable-next-line svelte/no-navigation-without-resolve
		goto('?' + newParams.toString());
	}

	function updateFilter(newFilter: string) {
		const newParams = new SvelteURLSearchParams(filters);
		newParams.set('search', newFilter);
		// eslint-disable-next-line svelte/no-navigation-without-resolve
		goto('?' + newParams.toString());
	}
</script>

{#snippet searchBar()}
	<input type="text" onchange={(e) => updateFilter(e.currentTarget.value)} />
{/snippet}
{@render header(searchBar)}
{#each filteredData as d (d.id)}
	{@render entry(d)}
{/each}
{#snippet pageSelector(num: number)}
	<button onclick={() => updatePage(() => num)}>{num}</button>
{/snippet}

{#snippet middleSelector(num: number)}
	{#if num > 1 && num < lastPage}
		{@render pageSelector(num)}
	{/if}
{/snippet}

{#snippet elipsis()}
	...
{/snippet}

{@render pageSelector(1)}

{#if pageNumber - 2 > 1}
	{@render elipsis()}
{/if}

{@render middleSelector(pageNumber - 1)}
{@render middleSelector(pageNumber)}
{@render middleSelector(pageNumber + 1)}

{#if pageNumber + 2 < lastPage}
	{@render elipsis()}
{/if}

{#if lastPage != 1}
	{@render pageSelector(lastPage)}
{/if}
