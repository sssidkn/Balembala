<script lang="ts">
	import Logo from '$lib/components/Logo.svelte';
	let { children } = $props();
</script>

<div class="vstack">
	<Logo />
	<div class="card">
		{@render children()}
	</div>
</div>

<style>
	.vstack {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		margin-top: 5em;
		row-gap: 5pt;
	}
	:global(body) {
		background-color: var(--dark-background);
	}
</style>
