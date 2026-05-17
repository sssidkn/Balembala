<script lang="ts">
	import { page } from '$app/state';
	import Logo from '$lib/components/Logo.svelte';
	import { resolve } from '$app/paths';
</script>

<div class="whole-bg error-page">
	<div class="vstack">
		<Logo />
		<div class="error-code"><b>{page.status}</b></div>
		<div class="description"><b>Something went wrong: {page.error?.message}</b></div>
		<a href={resolve('/')}>Return to home page</a>
	</div>
</div>

<style>
	.error-page {
		color: var(--light-text);
		background-color: var(--dark-background);
		.vstack {
			margin: 0;
			align-items: center;
		}
	}
	.error-code {
		font-size: 300pt;
	}
	.description {
		font-size: xx-large;
	}
</style>
