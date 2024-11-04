<script lang="ts">
    import { onDestroy } from 'svelte';
    import Navbar from '$lib/components/navbar.svelte';
    import { MapLibre, DefaultMarker, Popup, GeolocateControl } from 'svelte-maplibre';
    import SpotsListComponent from '$lib/components/spot-listings/spots-list-component.svelte';
    import { Button } from 'carbon-components-svelte';
    import { Search } from 'carbon-components-svelte';
    import type { LocationResult } from '$lib/types/spot/location-result';
    import BottomPanelOpen from 'carbon-icons-svelte/lib/BottomPanelOpen.svelte';
    import type { components } from '$lib/sdk/schema';
    import { newClient } from '$lib/utils/client';
    import { handleGetError } from '$lib/utils/error-handler';

    type ParkingSpot = components['schemas']['ParkingSpot'];

    const apiKey = import.meta.env.VITE_GEOCODING_API_KEY;
    const maxZoom: number = 12;
    let initZoom: number = 0;
    const selectedZoom: number = 13;
    const offset: [number, number] = [0, -10];

    let mapCenter: [number, number] = [0,0];
    let zoom: number = 10;

    let searchQuery = '';
    let results: LocationResult[] = [];
    let dropdownOpen: boolean = false;
    let selectedListingId: string | null = null;
    let showBackToTop: boolean = false;
    let listingsContainer: HTMLDivElement;
    let modalOpen : boolean = false;

    let debounceTimer: number | undefined;
    
    let numVisited : number = 0;

    let client = newClient();
    //export let data: PageData;
    let spotsData : ParkingSpot[] = [];

    const searchLocation = () => {
        const url = `https://api.geoapify.com/v1/geocode/search?text=${encodeURIComponent(searchQuery)}&apiKey=${apiKey}`;

        fetch(url)
            .then((response) => response.json())
            .then((data) => {
                results = data.features ?? [];
                dropdownOpen = true;
            })
            .catch((error) => {
                console.error('Error fetching data:', error);
            });
    };

    const handleInput = () => {
        dropdownOpen = true;
        clearTimeout(debounceTimer);
        debounceTimer = window.setTimeout(() => {
            if (searchQuery.length > 3) {
                searchLocation();
            }
        }, 500); // 0.5s delay
    };

    const handleSelect = (location: {
        properties: { formatted: string };
        geometry: { coordinates: number[] };
    }) => {
        searchQuery = location.properties.formatted;
        mapCenter = [location.geometry.coordinates[0], location.geometry.coordinates[1]];
        initZoom = selectedZoom;
        results = [];
        dropdownOpen = false;
        if (numVisited === 0) {
            numVisited += 1;
        }
        fetchSpots(location.geometry.coordinates);
    };

    const handleClickOutside = (event: Event) => {
        if (
            !(event.target as HTMLElement).closest('.dropdown') &&
            !(event.target as HTMLElement).closest('.search-input')
        ) {
            dropdownOpen = false;
        }
    };

    const handleMarkerClick = (listingId: string) => {
        selectedListingId = listingId;
        document
            .getElementById(`listing-${listingId}`)
            ?.scrollIntoView({ behavior: 'smooth', block: 'center' });
    };

    const scrollToTop = () => {
        if (listingsContainer) {
            listingsContainer.scrollTo({ top: 0, behavior: 'smooth' });
        }
    };

    const handleScroll = () => {
        showBackToTop = listingsContainer.scrollTop > 0;
    };

    const handleListingClick = () => {
        modalOpen = true;
    };

    const fetchSpots = async (coordinates: number[]) => {
        const { data: spots, error: errorSpots } = await client.GET('/spots', {
            params: {
                query: {
                    latitude: coordinates[1],
                    longitude: coordinates[0],
                    distance: 10000
                }
            }
        });
        handleGetError(errorSpots);
        spotsData = coalesceListings(spots??[]);
    }

    function coalesceListings(listings: ParkingSpot[]): ParkingSpot[] {
        const uniqueListings: ParkingSpot[] = [];
        const seenIds = new Set<string>();
    
        for (const listing of listings) {
            if (!seenIds.has(listing.id)) {
                uniqueListings.push(listing);
                seenIds.add(listing.id);
            }
        }
    
        return uniqueListings;
    }

    onDestroy(() => {
        clearTimeout(debounceTimer);
    });
</script>

<Navbar />

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions-->
<div class="container" on:click={handleClickOutside}>
    <div class="listings" bind:this={listingsContainer} on:scroll={handleScroll}>
        <Search
            class="search-input"
            type="text"
            bind:value={searchQuery}
            placeholder="Explore spots... [Atleast 3 characters required]"
            name="search"
            on:input={handleInput}
        />

        {#if results.length > 0 && dropdownOpen}
            <ul class="dropdown">
                {#each results as result}
                    <!-- svelte-ignore a11y-click-events-have-key-events a11y-no-noninteractive-element-interactions -->
                    <li on:click={() => handleSelect(result)}>
                        {result.properties.formatted}
                    </li>
                {/each}
            </ul>
        {/if}

        {#if spotsData.length > 0}
            {#each spotsData as listing}
            <!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
                <div
                    class="booking-info-container {listing.id === selectedListingId ? 'highlight' : ''}"
                    id={`listing-${listing.id}`}
                    on:click={() => handleListingClick()}
                >
                    <SpotsListComponent {listing}/>
                </div>
            {/each}
        {:else if numVisited === 0}
            <div class = "empty-container">
                <h2>Search for your destination ↑ </h2>
            </div>
        {:else}
            <div class = "empty-container">
                <h3>No listings Found!</h3>
            </div>
        {/if}

        {#if showBackToTop}
            <button class="back-to-top" on:click={scrollToTop}> ↑ Back to Top </button>
        {/if}
    </div>

    <div class="map-view">
        <MapLibre
            center={mapCenter}
            zoom = {initZoom}
            style="https://basemaps.cartocdn.com/gl/positron-gl-style/style.json"
        >
            <GeolocateControl position="bottom-left" fitBoundsOptions={{ maxZoom: maxZoom }} />
            {#each spotsData as { id, location }}
                <DefaultMarker lngLat={[location.longitude?? mapCenter[0], location.latitude?? mapCenter[1]]}>
                    <Popup {offset}>
                        <div class="popup-container">
                            <h2 class="popup-text">
                                {location.street_address} <br />
                                {location.city},{location.state} <br />
                                {location.postal_code} <br />
                            </h2>
                            <Button
                                kind="secondary"
                                on:click={() => handleMarkerClick(id)}
                                icon={BottomPanelOpen}
                            >
                                Go to listing
                            </Button>
                        </div>
                    </Popup>
                </DefaultMarker>
            {/each}
        </MapLibre>
    </div>
</div>

<style>
    .container {
        display: flex;
        flex-direction: row;
        max-height: fit-content;
    }

    .listings {
        display: flex;
        width: 25%;
        height: 85%;
        overflow-y: auto;
        flex-direction: column;
        padding: 0.5rem;
        background-color: #ffffff;
        position: fixed;
        z-index: 999;
        border: 2px solid rgb(0, 0, 0);
        border-radius: 5px;
    }

    .dropdown {
        position: absolute;
        background-color: white;
        border: 1px solid #868686;
        width: 100%;
        max-height: 10rem;
        overflow-y: auto;
        z-index: 1000;
        margin-top: 3rem;
    }
    .dropdown li {
        padding: 0.5rem;
        cursor: pointer;
    }
    .dropdown li:hover {
        background-color: #f0f0f0;
    }

    @keyframes smoothBlink {
        0%,
        100% {
            background-color: initial;
            opacity: 1;
        }
        50% {
            background-color: #007bff;
            opacity: 0.5;
        }
    }

    .highlight {
        animation: smoothBlink 0.6s linear 4;
    }

    .back-to-top {
        position: fixed;
        padding: 0.5rem 1rem;
        background-color: #007bff;
        color: white;
        border: none;
        border-radius: 5px;
        cursor: pointer;
        box-shadow: 0px 4px 8px rgba(0, 0, 0, 0.2);
        transition: background-color 0.8s ease;
    }

    .popup-container {
        display: flex;
        flex-direction: column;
        padding: 0.5rem;
    }

    .popup-text {
        margin: 1rem;
        font-family: Fredoka, sans-serif;
    }

    .empty-container{
        display: flex;
        justify-content: center;
        margin: 1rem;
    }
</style>