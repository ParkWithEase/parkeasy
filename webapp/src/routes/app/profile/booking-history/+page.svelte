<script lang="ts">
    import BookingHistoryDisplay from '$lib/components/history-list-component/history-list-display.svelte';
    let bookingList = [
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

<div class="list-container">
    {#key bookingList}
        {#each bookingList as booking}
            <div class="booking-info-container">
                <BookingHistoryDisplay {booking}></BookingHistoryDisplay>
            </div>
        {/each}
    {/key}
</div>

<style>
    :global(.bx--content) {
        padding: 0rem;
    }

    .booking-info-container {
        display: flex;
        flex-direction: row;
        color: rgb(0, 0, 0);
        border: 1px solid #cfcfcfcf;
    }
</style>
