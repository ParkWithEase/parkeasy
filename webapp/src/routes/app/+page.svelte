<script lang="ts">
    import Navbar from '$lib/components/navbar.svelte';
    import { MapLibre, DefaultMarker, Popup, GeolocateControl } from 'svelte-maplibre';
    import SpotsListComponent from '$lib/components/spot-listings/spots-list-component.svelte';
    import { spots_data } from './your-spots/mock_data';
    import { Search } from "carbon-components-svelte";
    import { text } from 'stream/consumers';

    let searchText = '';
    let results = [];

    var requestOptions = {
        method: 'GET',
    };

    // Function to fetch data from Geoapify API
    async function fetchSpots() {
        if (searchText) {
            const apiKey = '4fc03e04196343548bf1d6d27e7bf6c0';
            const apiUrl = `https://api.geoapify.com/v1/geocode/search?text=${encodeURIComponent(searchText)}&apiKey=${apiKey}`;
            try {
                
                const response = await fetch(apiUrl, requestOptions)
                                .then(response => response.json())
                                .then(result => console.log(result))
                                .catch(error => console.log('error', error));
            } catch (error) {
                console.error('Error fetching data:', error);
            }
        } else {
            results = [];
        }
    }

    function handleSearchInput(event:any) {
        fetchSpots();
    }
</script>

<Navbar />

<body>
    <div class="container">
        <div class="listings">
            <Search placeholder = "Search for your destination... " bind:value={searchText} on:input={handleSearchInput} on:clear={() => console.log("clear")}/>
            {#key spots_data}
                {#each spots_data as listing}
                    <div class="booking-info-container">
                        <SpotsListComponent {listing} />
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
