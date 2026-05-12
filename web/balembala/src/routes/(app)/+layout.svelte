<script lang="ts">
	import './common.css';
	import { resolve } from '$app/paths';
	import Logout from '$lib/assets/icons/Logout.svg?component';
	import Logo from '$lib/components/Logo.svelte';

	let { children } = $props();
</script>

<div class="top-bar">
	<Logo />
	<nav>
		<b>
			<a href={resolve('/templates')}>Templates</a>
			|
			<a href={resolve('/contacts')}>Contacts</a>
		</b>
	</nav>
	<form method="POST" action="/?/logout">
		<button class="icon">
			<Logout />Log out
		</button>
	</form>
</div>
<div class="center">{@render children()}</div>

<style>
	.top-bar {
		color: var(--light-text);
		display: flex;
		align-items: center;
		margin: 0;
		background-color: var(--dark-background);
	}
	button {
		margin-right: 10pt;
	}
	nav {
		flex-grow: 1;
		font-size: x-large;
		margin-left: 10pt;
	}
	a {
		text-decoration: none;
		color: inherit;
	}

	:global(body) {
		background-color: var(--hard-light-background);
	}
</style>
