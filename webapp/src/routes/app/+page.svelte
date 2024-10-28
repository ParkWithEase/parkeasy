<script lang="ts">
    import Navbar from '$lib/components/navbar.svelte';
    import { MapLibre, DefaultMarker, Popup, GeolocateControl } from 'svelte-maplibre';
    import ListingsListComponent from '$lib/components/spot-listings/listings-list-component.svelte';
    import { MapLibreSearchControl } from "@stadiamaps/maplibre-search-box";

    const coordsList: number[][] = [
        [-97.133290, 49.808856],
        [-97.199808, 49.811791]
    ];

    const lngLatList: [number, number][] = coordsList as [number, number][];
    let listings= [
        {
        lngLat: lngLatList[0],
        label: 'Engineering Building',
        name: 'Engineering Building UofM',
        },
        {
        lngLat: lngLatList[1],
        label: 'Whyte Ridge',
        name: 'Whyte Ridge',
        }
    ];

</script>

<Navbar />

<body>
    <div class="container">
        <div class="listings">
            {#key listings}
                {#each listings as listing}
                    <div class="booking-info-container">
                        <ListingsListComponent {listing}></ListingsListComponent>
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
                {#each listings as { lngLat, label, name }}
                    <DefaultMarker {lngLat} draggable>
                        <Popup
                            offset={[0, -10]}
                            on:open={async () => {
                                const resp = await fetch(`/examples/popup_remote/${name}`);
                                const result = await resp.json();
                            }}
                        >
                            <div class="text-lg font-bold">{name}</div>
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
        width: 30%;
        min-height: 100%;
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
