<script lang="ts">
	import FilteredList from '$lib/components/FilteredList.svelte';
	import type { PageProps } from './$types';
	const { data }: PageProps = $props();
</script>

<FilteredList data={data.data} pageSize={1}>
	{#snippet header(searchBar)}
		{@render searchBar()}
	{/snippet}

	{#snippet entry(d)}
		{d.id}
	{/snippet}
</FilteredList>
