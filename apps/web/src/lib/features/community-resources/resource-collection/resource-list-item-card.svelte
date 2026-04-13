<script lang="ts">
	import { CardHeader, Card } from '$lib/components/ui/card';

	let { children } = $props();
</script>

<Card class="transition-all duration-300 hover:-translate-y-0.5 hover:shadow-md">
	<CardHeader class="gap-3">
		{@render children()}
	</CardHeader>
</Card>
