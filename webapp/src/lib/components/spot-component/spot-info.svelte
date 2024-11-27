<script lang="ts">
    import { Checkbox } from 'carbon-components-svelte';

    import type { components } from '$lib/sdk/schema';
    type ParkingSpot = components['schemas']['ParkingSpot'];
    export let spot: ParkingSpot | undefined;
</script>

<div class="spot-info-container">
    <div>
        <p class="spot-info-label">Street Address</p>
        <p class="spot-info-content">{spot?.location.street_address}</p>
    </div>
    <div>
        <p class="spot-info-label">City</p>
        <p class="spot-info-content">{spot?.location.city}</p>
    </div>
    <div>
        <p class="spot-info-label">State/Provice</p>
        <p class="spot-info-content">{spot?.location.state}</p>
    </div>
    <div>
        <p class="spot-info-label">Country Code</p>
        <p class="spot-info-content">{spot?.location.country_code}</p>
    </div>
    <div>
        <p class="spot-info-label">Postal Code</p>
        <p class="spot-info-content">{spot?.location.postal_code}</p>
    </div>

    <div>
        <p class="spot-info-label">Utilities</p>
        <Checkbox
            name="shelter"
            labelText="Shelter"
            checked={spot?.features?.shelter}
            style="pointer-events: none;"
            readonly
        />
        <Checkbox
            name="plug-in"
            labelText="Plug-in"
            checked={spot?.features?.plug_in}
            style="pointer-events: none;"
            readonly
        />
        <Checkbox
            name="charging-station"
            labelText="Charging Station"
            checked={spot?.features?.charging_station}
            style="pointer-events: none;"
            readonly
        />
    </div>
</div>

<style>
    .spot-info-label {
        font-size: 0.8rem;
    }

    .spot-info-content {
        font-size: 1.3rem;
    }

    .spot-info-container {
        display: flex;
        gap: 0.7rem;
        flex-direction: column;
        max-width: 20rem;
        margin-bottom: 1rem;
    }
</style>
