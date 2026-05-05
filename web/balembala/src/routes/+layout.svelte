<script lang="ts">
	import favicon from '$lib/assets/favicon.svg';

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

{@render children()}

<style>
	:root {
		--light-text: #F9FAFB;
		--gray-text: #667085;
		--border: #667085;
		--dark-background: #101828;
		--soft-light-background: #F2F4F7;
		--hard-light-background: #FFFFFF;
		--button: #2563EB;
		--error: #EF4444;
	}
	:global {
		body {
        	font-family: Arial;
    	}
		button {
			border: none;
			border-radius: 5pt;
			padding: 5pt;
			background-color: var(--button);
			color: var(--light-text);
		}
		.card {
        	border-radius: 16pt;
        	background: var(--soft-light-background);
        	padding: 16pt;
    	}
		label.text-form {
			font-size: 0.8em;
			display: block;
			color: var(--gray-text);
			margin-left: 3pt;
		}
		input.text-form {
			display: block;
			border-radius: 5pt;
			border-width: 1pt;
			border-color: var(--border);
			padding: 3pt;
		}
	}
</style>
