<script lang="ts">
    import LeasingHistoryDisplay from '$lib/components/history-list-component/history-list-display.svelte';
    import logo from '$lib/images/parkeasy-logo.png';
    let leasingList = [
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
        },
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
        },
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
    {#key leasingList}
        {#each leasingList as leasing}
            <!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
            <div class="leasing-info-container">
                <img class="booking-image" src={logo} alt="a car" />
                <LeasingHistoryDisplay booking={leasing}></LeasingHistoryDisplay>
            </div>
        {/each}
    {/key}
</div>

<style>
    :global(.bx--content){
        padding: 0rem;
    }

    .booking-image{
        max-height: 9%;
        max-width: 12%;
    }

    .leasing-info-container {
        display: flex;
        flex-direction: row;
        color: rgb(0, 0, 0);
        border: 1px solid #cfcfcfcf;
    }
</style>
