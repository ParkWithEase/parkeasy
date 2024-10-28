<script lang="ts">
    import { Modal } from 'carbon-components-svelte';
    export let listing;
    export let open = false;
    import { MapLibre, DefaultMarker } from 'svelte-maplibre';
</script>

<Modal
    bind:open
    modalHeading={listing.location.street_address}
    primaryButtonText="Book now"
    on:open
    on:close
    on:submit
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
        <MapLibre
            center={[listing.location.longitude, listing.location.latitude]}
            zoom={13}
            class="map"
            style="https://basemaps.cartocdn.com/gl/positron-gl-style/style.json"
        >
            <DefaultMarker lngLat={[listing.location.longitude, listing.location.latitude]} draggable/>
        </MapLibre>
    </div>
</Modal>

<style>
    .modal-content{
        display: flex;
        flex-direction: row;
    }
    .listing-content {
        padding: 1rem;
    }
</style>