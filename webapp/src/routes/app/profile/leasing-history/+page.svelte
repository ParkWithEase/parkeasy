<script lang="ts">
    import SpotTransactionDisplay from '$lib/components/transaction/spot-transaction-display.svelte';
    import { TransactionType } from '$lib/enum/transaction_type';
    import IntersectionObserver from 'svelte-intersection-observer';
    import type { PageData } from './$types';

    export let data: PageData;

    let loadLock = false;
    let canLoadMore = data.hasNext;
    let loadTrigger: HTMLElement | null = null;
    let intersecting: boolean;

    $: while (intersecting && canLoadMore && !loadLock) {
        loadLock = true;
        data.paging
            .next()
            .then(({ value: { data: leasing_transactions }, done }) => {
                if (leasing_transactions) {
                    data.leasing_transactions = [
                        ...data.leasing_transactions,
                        ...leasing_transactions
                    ];
                }
                canLoadMore = !done;
            })
            .finally(() => {
                loadLock = false;
            });
    }
</script>

<div class="list-container">
    {#key data.leasing_transactions}
        {#each data.leasing_transactions as transaction}
            <a href={`/app/transaction/${transaction.id}`} style="text-decoration: none;">
                <SpotTransactionDisplay {transaction} transaction_type={TransactionType.BOOK} />
            </a>
        {/each}
    {/key}
</div>

<IntersectionObserver element={loadTrigger} bind:intersecting>
    {#if canLoadMore}
        <div bind:this={loadTrigger}>Loading...</div>
    {/if}
</IntersectionObserver>
