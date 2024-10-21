<script lang="ts">
    import type { Booking } from '$lib/types/booking-history/booking';
    import LeasingHistoryDisplay from '$lib/components/history-list-component/history-list-display.svelte';
    let leasingList: Booking[] = [
        {
            booking_id: '000101010',
            address: 'Chancellor\'s Circle',
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
        },
        {
            booking_id: '000101010',
            address: 'Chancellor\'s Circle',
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
        },
        {
            booking_id: '000101010',
            address: 'Chancellor\'s Circle',
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
    {#key leasingList}
        {#each leasingList as leasing}
            <!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
            <div class="leasing-info-container">
                <LeasingHistoryDisplay booking = {leasing}></LeasingHistoryDisplay>
            </div>
        {/each}
        
    {/key}
</div>



<style>
    .leasing-info-container {
        display: flex;
        flex-direction: row;
        padding: 1rem;
        gap: 2rem;
        color: rgb(0, 0, 0);
    }

</style>