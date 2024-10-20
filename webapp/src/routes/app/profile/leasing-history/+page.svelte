<script lang="ts">
    import type { Booking } from '$lib/types/booking-history/booking';
    import leasingImg from '$lib/images/parkeasy-logo.png';
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

<h1>Leasing History</h1>

<div class="list-container">
    {#key leasingList}
        {#each leasingList as leasing, index}
            <!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
            <div class="leasing-info-container">
                <img class="leasing-image" src={leasingImg} alt="a leasing" />
                <div class="leasing-info">
                    <div>
                        <p class="leasing-license-plate">{leasing.address}</p>
                        <p>{leasing.postalCode}</p>
                    </div>
                    <div>
                        <p class="leasing-data">From: {leasing.time_from}</p>
                        <p class="leasing-data">To: {leasing.time_to}</p>
                        <p class="leasing-data">Total price: CA${leasing.price}</p>
                    </div>
                </div>
            </div>
            <hr class="separator">
        {/each}
        
    {/key}
</div>



<style>
    h1{
        font-weight: bold;
    }

    .separator{
        width: 90%;
    }
    .list-container{
        width: 80%;
        border-top: 2px solid black;
    }
    .leasing-info-container {
        display: flex;
        flex-direction: row;
        margin-left: 1rem;
        padding: 1rem;
        gap: 2rem;
        color: rgb(0, 0, 0);
    }
    .leasing-info {
        flex: 1;
        display: flex;
        flex-direction: column
    }
    .leasing-image {
        max-height: 10%;
        max-width: 13%;
    }
    .leasing-license-plate {
        font-size: 1.5rem;
        border-bottom: 2px solid white;
    }
    .leasing-data {
        font-size: 1rem;
    }
</style>