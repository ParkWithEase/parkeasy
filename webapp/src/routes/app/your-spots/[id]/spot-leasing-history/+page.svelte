<script lang="ts">
    import SpotTransactionDisplay from '$lib/components/transaction/spot-transaction-display.svelte';
    import { TransactionType } from '$lib/enum/transaction_type';
    import { Content } from 'carbon-components-svelte';
    import type { PageData } from './$types';
    import IntersectionObserver from 'svelte-intersection-observer';

    export let data: PageData;

    let loadLock = false;
    let canLoadMore = data.hasNext;
    let loadTrigger: HTMLElement | null = null;
    let intersecting: boolean;

    $: while (intersecting && canLoadMore && !loadLock) {
        loadLock = true;
        data.paging
            .next()
            .then(({ value: { data: leasing_transaction }, done }) => {
                if (leasing_transaction) {
                    data.leasing_transaction = [
                        ...data.leasing_transaction,
                        ...leasing_transaction
                    ];
                }
                canLoadMore = !done;
            })
            .finally(() => {
                loadLock = false;
            });
    }
</script>

<Content>
    <div class="list-container">
        {#key data.leasing_transaction}
            {#each data.leasing_transaction as transaction}
                <a href={`/app/transaction/${transaction.id}`} style="text-decoration: none;">
                    <SpotTransactionDisplay
                        {transaction}
                        transaction_type={TransactionType.LEASE}
                    />
                </a>
            {/each}
        {/key}
    </div>

    <IntersectionObserver element={loadTrigger} bind:intersecting>
        {#if canLoadMore}
            <div bind:this={loadTrigger}>Loading...</div>
        {/if}
    </IntersectionObserver>
</Content>
