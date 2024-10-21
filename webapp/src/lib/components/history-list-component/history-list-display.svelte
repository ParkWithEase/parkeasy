<script lang="ts">
    import type { Booking } from '$lib/types/booking-history/booking';
    export let booking: Booking;
</script>

<div class="booking-info">
    <div>
        <p class="booking-head">{booking.address}</p>
        <p>{booking.postalCode}</p>
    </div>
    <div>
        <p class="booking-data">From: {booking.time_from}</p>
        <p class="booking-data">To: {booking.time_to}</p>
        <p class="booking-data">Total price: CA${booking.price}</p>
    </div>
</div>

<style>
    .booking-info {
        padding-left: 1rem;
        flex: 1;
        display: flex;
        flex-direction: column;
        border: 1px solid #cfcfcfcf;
    }
    .booking-head {
        font-size: 1.5rem;
        border-bottom: 2px solid white;
    }
    .booking-data {
        font-size: 1rem;
    }
</style>
