<script lang="ts">
    import { goto } from '$app/navigation';
    import { Modal } from 'carbon-components-svelte';
    export let listing;
    export let open = false;
    import { MapLibre, DefaultMarker } from 'svelte-maplibre';

    const defaultZoom: number = 13;
</script>

<Modal
    bind:open
    modalHeading={listing.location.street_address}
    primaryButtonText="Book now"
    on:open
    on:close
    on:submit={() => goto(`/app/booking/${listing.id}/spot-booking`)}
>
    <div class="modal-content">
        <div class="listing-content">
            <p><strong>Postal Code:</strong> {listing.location.postal_code}</p>
            <p><strong>City:</strong> {listing.location.city}</p>
            <p><strong>State:</strong> {listing.location.state}</p>
            <p><strong>Features:</strong></p>
            <ul>
                <li>Charging Station: {listing.features.charging_station ? 'Yes' : 'No'}</li>
                <li>Plug-In: {listing.features.plug_in ? 'Yes' : 'No'}</li>
                <li>Shelter: {listing.features.shelter ? 'Yes' : 'No'}</li>
            </ul>
        </div>
        <div>
            <MapLibre
                center={[listing.location.longitude, listing.location.latitude]}
                zoom={defaultZoom}
                style="https://basemaps.cartocdn.com/gl/positron-gl-style/style.json"
            >
                <DefaultMarker lngLat={[listing.location.longitude, listing.location.latitude]} />
            </MapLibre>
        </div>
    </div>
</Modal>

<style>
    .modal-content {
        display: flex;
        flex-direction: row;
    }
    .listing-content {
        padding: 1rem;
        z-index: 1000;
        width: 30%;
        background-color: rgb(209, 209, 209);
        border: 1px solid black;
        border-radius: 5px;
        margin-left: 2rem;
    }
</style>
