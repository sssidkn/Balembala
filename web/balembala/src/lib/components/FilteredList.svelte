<script lang="ts" generics="Data extends { id: string }">
	import type { Snippet } from 'svelte';
	import Search from '$lib/assets/icons/Search.svg?component';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { SvelteURLSearchParams } from 'svelte/reactivity';

	interface Props {
		data: Data[];
		pageSize: number;
		isMatches: (query: string, data: Data) => boolean;
		header: Snippet<[Snippet<[string]>]>;
		entry: Snippet<[Data]>;
	}

	const { data, pageSize, isMatches, header, entry }: Props = $props();

	let filters = $derived({
		page: page.url.searchParams.get('page') ?? '1',
		search: page.url.searchParams.get('search') ?? ''
	});
	let pageNumber = $derived(Number(filters.page));
	let filteredData = $derived.by(() => {
		const search = filters.search.toLowerCase();
		return data.filter((d) => isMatches(search, d));
	});
	let slicedData = $derived.by(() => {
		const start = (pageNumber - 1) * pageSize;
		return filteredData.slice(start, start + pageSize);
	});
	let lastPage = $derived(Math.ceil(filteredData.length / pageSize));

	function updatePage(nextPage: number) {
		const newParams = new SvelteURLSearchParams(filters);
		newParams.set('page', nextPage.toString());
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

{#snippet searchBar(placeholder: string)}
	<span class="search">
		<Search class="icon"></Search>
		<input type="search" {placeholder} onchange={(e) => updateFilter(e.currentTarget.value)} />
	</span>
{/snippet}
{@render header(searchBar)}
{#each slicedData as d (d.id)}
	{@render entry(d)}
{/each}

{#snippet pageSelector(num: number)}
	<button class="selector" onclick={() => updatePage(num)}>{num}</button>
{/snippet}

{#snippet middleSelector(num: number)}
	{#if num > 1 && num < lastPage}
		{@render pageSelector(num)}
	{/if}
{/snippet}

{#snippet elipsis()}
	<div class="selector" style:color="var(--gray-text)">...</div>
{/snippet}

<div>
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

	{#if lastPage > 1}
		{@render pageSelector(lastPage)}
	{/if}
</div>

<style>
	.search {
		display: inline-flex;
		align-items: center;
		position: relative;

		:global(.icon) {
			--size: 10pt;
			padding-left: 3pt;
			position: absolute;
			pointer-events: none;
		}
		input {
			background-color: inherit;
			padding: 3pt;
			padding-left: 14pt;
			border-radius: 5pt;
			border-width: 1pt;
			border-color: var(--border);
		}
	}
	.selector {
		width: 20pt;
		height: 20pt;
		display: inline-flex;
		align-items: center;
		justify-content: center;
	}

	button.selector {
		border: solid 1.5pt var(--border);
	}
</style>
