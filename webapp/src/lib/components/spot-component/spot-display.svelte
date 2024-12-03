<script lang="ts">
    import BackgroundImage from '$lib/images/background.png';
    import { Checkbox } from 'carbon-components-svelte';
    import type { components } from '$lib/sdk/schema';
    type ParkingSpot = components['schemas']['ParkingSpot'];
    export let spot: ParkingSpot;
</script>

{#if spot}
    <div class="spot-info">
        <div>
            <img src={BackgroundImage} class="spot-image" alt="A parking spot" />
        </div>
        <div>
            <div>
                <p class="spot-title">
                    {spot.location.street_address}, {spot.location.postal_code}
                </p>
                <span style="font-size: 1rem"
                    >{spot.location.city}, {spot.location.state} {spot.location.country_code}</span
                >
            </div>
            <div>
                <Checkbox labelText="Shelter" checked={spot.features?.shelter} readonly />
                <Checkbox labelText="Plug-in" checked={spot.features?.plug_in} readonly />
                <Checkbox
                    labelText="Charging station"
                    checked={spot.features?.charging_station}
                    readonly
                />
            </div>
        </div>
    </div>
{/if}

<style>
    .spot-info {
        padding-left: 1rem;
        flex: 1;
        gap: 1rem;
        display: flex;
        align-items: center;
        flex-direction: row;
        border: 1px solid #cfcfcfcf;
        color: black;
    }

    .spot-image {
        max-height: 10rem;
    }

    .spot-title {
        line-height: 1;
        font-size: 1.5rem;
        padding-top: 0.5rem;
    }
</style>
