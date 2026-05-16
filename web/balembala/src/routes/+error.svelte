<script lang="ts">
	import { page } from '$app/state';
	import Logo from '$lib/components/Logo.svelte';
	import { resolve } from '$app/paths';
</script>

<div class="whole-bg" style:background-color="var(--dark-background)">
	<div class="vstack">
		<Logo />
		<div class="error-code"><b>{page.status}</b></div>
		<div class="description"><b>Something went wrong</b></div>
		<a href={resolve('/')}>Return to home page</a>
	</div>
</div>

<style>
	.vstack {
		margin: 15pt;
		display: flex;
		flex-direction: column;
		align-items: center;
		row-gap: 10pt;
		color: var(--light-text);
	}
	.error-code {
		font-size: 300pt;
	}
	.description {
		font-size: xx-large;
	}
</style>
