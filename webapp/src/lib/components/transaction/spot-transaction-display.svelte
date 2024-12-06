<script lang="ts">
    import { TransactionType } from '$lib/enum/transaction_type';

    export let transaction;

    export let transaction_type: TransactionType = TransactionType.BOOK;
</script>

<div class="transaction-info">
    <div class="transaction-header">
        <p class="transaction-title">
            {transaction.parkingspot_location.street_address}, {transaction.parkingspot_location
                .postal_code}
        </p>
        <span class="transaction-subtitle"
            >{transaction.parkingspot_location.city}, {transaction.parkingspot_location.state}
            {transaction.parkingspot_location.country_code}</span
        >
    </div>
    <div>
        <div class="car-detail-section">
            <p>Car Detail:</p>
            <div class="car-detail">
                <p>License plate: {transaction?.car_details.license_plate}</p>
                <p>Color: {transaction?.car_details.color}</p>
                <p>Model: {transaction?.car_details.model}</p>
                <p>Make: {transaction?.car_details.make}</p>
            </div>
        </div>
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
        border: 1px solid #241c1ca6;
        margin: 0.5rem;
        border-radius: 5px;
    }

    .transaction-header {
        border-bottom: 1px solid rgba(0, 0, 0, 0.5);
        padding-bottom: 0.5rem;
    }
    .transaction-title {
        line-height: 1;
        font-size: 1.5rem;
        padding-top: 0.5rem;
    }

    .transaction-subtitle {
        font-size: 1rem;
    }

    .car-detail-section {
        display: flex;
        flex-direction: column;
    }

    .car-detail {
        padding-left: 1rem;
    }
</style>
