<script>
	import EditContact from '$lib/components/EditContact.svelte';

	const { data } = $props();
</script>

<EditContact data={data.data}></EditContact>
