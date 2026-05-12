<script lang="ts">
	import EditTemplate from '$lib/components/EditTemplate.svelte';
	import type { PageProps } from './$types';
	const { data }: PageProps = $props();
</script>

<EditTemplate {data}></EditTemplate>
