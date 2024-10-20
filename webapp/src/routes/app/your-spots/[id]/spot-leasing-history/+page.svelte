<script lang="ts">
    import SpotLeasingDisplay from '$lib/components/transaction/spot-leasing-display.svelte';
    import { Content } from 'carbon-components-svelte';
    import type { PageData } from './$types';

    export let data: PageData;
</script>

<Content>
    <div>
        {#key data.transactions}
            {#each data?.transactions as transaction}
                <SpotLeasingDisplay {transaction} />
            {/each}
        {/key}
    </div>
</Content>
