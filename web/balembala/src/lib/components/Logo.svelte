<script>
	import { resolve } from '$app/paths';
</script>

<a href={resolve('/')}>
	<b>Balembala</b>
	<br />
	<small>Notification system</small>
</a>

<style>
	b {
		color: var(--light-text);
		font-size: xx-large;
	}
	small {
		color: var(--gray-text);
	}
	a {
		text-decoration: none;
		display: block;
		text-align: center;
		margin: 5pt;
	}
</style>
