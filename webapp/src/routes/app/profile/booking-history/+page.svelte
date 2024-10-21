<script lang="ts">
    import type { Booking } from '$lib/types/booking-history/booking';
    import bookingImg from '$lib/images/parkeasy-logo.png';
    let bookingList: Booking[] = [
        {
            booking_id: '000101010',
            address: "Chancellor's Circle",
            postalCode: 'X1X1X1',
            price: 10.26,
            time_from: 800,
            time_to: 1000,
            duration: 2
        },
        {
            booking_id: '000101099',
            address: 'St. Vital',
            postalCode: 'ABC123',
            price: 31.48,
            time_from: 1600,
            time_to: 2030,
            duration: 4.5
        }
    ];
</script>

<h1>Booking History</h1>

<div class="list-container">
    {#key bookingList}
        {#each bookingList as booking, index}
            <!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
            <div class="booking-info-container">
                <img class="booking-image" src={bookingImg} alt="a booking" />
                <div class="booking-info">
                    <div>
                        <p class="booking-license-plate">{booking.address}</p>
                        <p>{booking.postalCode}</p>
                    </div>
                    <div>
                        <p class="booking-data">From: {booking.time_from}</p>
                        <p class="booking-data">To: {booking.time_to}</p>
                        <p class="booking-data">Total price: CA${booking.price}</p>
                    </div>
                </div>
            </div>
            <hr class="separator" />
        {/each}
    {/key}
</div>

<style>
    h1 {
        font-weight: bold;
    }

    .separator {
        width: 90%;
    }
    .list-container {
        width: 80%;
        border-top: 2px solid black;
    }
    .booking-info-container {
        display: flex;
        flex-direction: row;
        margin-left: 1rem;
        padding: 1rem;
        gap: 2rem;
        color: rgb(0, 0, 0);
    }
    .booking-info {
        flex: 1;
        display: flex;
        flex-direction: column;
    }
    .booking-image {
        max-height: 10%;
        max-width: 13%;
    }
    .booking-license-plate {
        font-size: 1.5rem;
        border-bottom: 2px solid white;
    }
    .booking-data {
        font-size: 1rem;
    }
</style>
