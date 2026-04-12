<script lang="ts">
	import Badge from './badge.svelte';

	let {
		counter,
		pl,
		sl = pl.substring(0, pl.length - 1) // remove last letter (normally 's' from plural by default)
	}: {
		counter: number;
		sl?: string;
		pl: string;
	} = $props();
</script>

<Badge>{counter} {counter == 1 ? sl : pl}</Badge>
