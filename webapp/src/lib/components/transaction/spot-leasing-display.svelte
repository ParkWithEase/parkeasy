<script lang="ts">
    export let transaction;
</script>

<div class="transaction-info">
    <div>
        <p class="transaction-title">
            {transaction.spot.street_address}, {transaction.spot.postal_code}
        </p>
        <span style="font-size: 1rem"
            >{transaction.spot.city}, {transaction.spot.state} {transaction.spot.country_code}</span
        >
    </div>
    <div>
        <p>
            Car Detail: {transaction.car.license_plate}, {transaction.car.color}, {transaction.car
                .model}, {transaction.car.make}
        </p>
        <p>
            Leased on: {transaction.time.date} from {transaction.time.start} to {transaction.time
                .end}
        </p>
        <p>Collected: {transaction.paid_amount}</p>
    </div>
</div>

<style>
    .transaction-info {
        padding-left: 1rem;
        flex: 1;
        gap: 1rem;
        display: flex;
        flex-direction: column;
        border: 1px solid #cfcfcfcf;
    }

    .transaction-title {
        line-height: 1;
        font-size: 1.5rem;
        padding-top: 0.5rem;
    }
</style>
