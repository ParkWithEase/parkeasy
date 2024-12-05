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
            .then(({ value: { data: booking_transactions }, done }) => {
                if (booking_transactions) {
                    data.booking_transactions = [
                        ...data.booking_transactions,
                        ...booking_transactions
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
        {#key data.booking_transactions}
            {#each data.booking_transactions as transaction}
                <a
                    href={`/app/booking-transaction/${transaction.id}`}
                    style="text-decoration: none;"
                >
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
</Content>
