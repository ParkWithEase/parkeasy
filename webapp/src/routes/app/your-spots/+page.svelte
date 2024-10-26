<script lang="ts">
    import { Button } from 'carbon-components-svelte';
    import type { PageData } from '../$types';
    import { Add } from 'carbon-icons-svelte';
    import SpotDisplay from '$lib/components/spot-component/spot-display.svelte';

    export let data: PageData;
</script>

<div>
    <Button style="margin: 1rem;" icon={Add} href="/app/add-your-spots">New Spot</Button>
    {#key data.spots}
        {#each data?.spots as spot}
            <SpotDisplay {spot} />
        {/each}
    {/key}
</div>
