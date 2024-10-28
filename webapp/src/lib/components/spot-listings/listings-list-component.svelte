<script lang="ts">
    export let listing;
</script>

<div class="listing-info">
    <div>
        <p class="listing-head">{listing.address}</p>
        <p>{listing.postalCode}</p>
    </div>
    <div>
        <p class="listing-data">From: {listing.time_from}</p>
        <p class="listing-data">To: {listing.time_to}</p>
    </div>
</div>

<style>
    .listing-info {
        flex: 1;
        display: flex;
        flex-direction: column;
    }
    .listing-head {
        font-size: 1.5rem;
        border-bottom: 2px solid white;
    }
    .listing-data {
        font-size: 1rem;
    }
</style>
