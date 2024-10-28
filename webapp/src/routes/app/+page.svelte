<script lang="ts">
    import Navbar from '$lib/components/navbar.svelte';
    import { MapLibre, DefaultMarker, Popup, GeolocateControl } from 'svelte-maplibre';
    import ListingsListComponent from '$lib/components/spot-listings/listings-list-component.svelte';
    import { spots_data } from './your-spots/mock_data';
</script>

<Navbar />

<body>
    <div class="container">
        <div class="listings">
            {#key spots_data}
                {#each spots_data as listing}
                    <div class="booking-info-container">
                        <ListingsListComponent {listing} />
                    </div>
                {/each}
            {/key}
        </div>
        <div class="map-view">
            <MapLibre
                center={[-97.1, 49.9]}
                zoom={10}
                class="map"
                style="https://basemaps.cartocdn.com/gl/positron-gl-style/style.json"
            >
                <GeolocateControl position="top-left" fitBoundsOptions={{ maxZoom: 12 }} />
                {#each spots_data as { id, features, location, isListed }}
                    <DefaultMarker lngLat={[location.longitude, location.latitude]} draggable>
                        <Popup
                            offset={[0, -10]}
                        >
                            <div class="text-lg font-bold">{location.street_address}</div>
                            <h1>Listing data here</h1>
                        </Popup>
                    </DefaultMarker>
                {/each}
            </MapLibre>
        </div>
    </div>
</body>

<style>
    .container {
        display: flex;
        flex-direction: row;
        max-height: fit-content;
    }

    .listings {
        display: flex;
        width: 35%;
        max-height: 80vh;
        overflow-y: auto;
        flex-direction: column;
        padding: 0.5rem;
        background-color: rgb(186, 214, 183);
    }

    .map-view {
        width: 100%;
        background-color: white;
    }

    :global(.map) {
        height: 80vh;
    }
</style>
