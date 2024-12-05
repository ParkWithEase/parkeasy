<script lang="ts">
    import { TransactionType } from '$lib/enum/transaction_type';

    export let transaction;

    export let transaction_type: TransactionType = TransactionType.BOOK;
</script>

<div class="transaction-info">
    <div class="transaction-header">
        <!-- <p class="transaction-title">
            {transaction.spot.street_address}, {transaction.spot.postal_code}
        </p>
        <span class="transaction-subtitle"
            >{transaction.spot.city}, {transaction.spot.state}
            {transaction.spot.country_code}</span
        > -->

        <p class="transaction-title">Street number, postal code</p>
        <span class="transaction-subtitle">city, state, country</span>
    </div>
    <div>
        <!-- <p>
            Car Detail: {transaction?.car.license_plate}, {transaction?.car.color}, {transaction
                ?.car.model}, {transaction?.car.make}
        </p> -->
        {#if transaction_type == TransactionType.LEASE}
            <p>
                Leased on: {new Date(transaction.booking_time)}
            </p>
            <p>Collected: {transaction.paid_amount} CAD</p>
        {:else if transaction_type == TransactionType.BOOK}
            <p>
                Booked on: {new Date(transaction.booking_time)}
            </p>
            <p>Paid: {transaction.paid_amount} CAD</p>
        {/if}
    </div>
</div>

<style>
    .transaction-info {
        color: black;
        padding: 0 1rem 0.5rem 1rem;
        flex: 1;
        gap: 1rem;
        display: flex;
        flex-direction: column;
        border: 1px solid #cfcfcfcf;
    }

    .transaction-header {
        border-bottom: 1px solid rgba(0, 0, 0, 0.5);
    }
    .transaction-title {
        line-height: 1;
        font-size: 1.5rem;
        padding-top: 0.5rem;
    }

    .transaction-subtitle {
        font-size: 1rem;
    }
</style>
