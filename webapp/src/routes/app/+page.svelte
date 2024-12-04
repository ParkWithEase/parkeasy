<script lang="ts">
    import { onDestroy } from 'svelte';
    import { MapLibre, DefaultMarker, Popup, GeolocateControl } from 'svelte-maplibre';
    import SpotsListComponent from '$lib/components/spot-listings/spots-list-component.svelte';
    import { Button } from 'carbon-components-svelte';
    import { Search, Slider } from 'carbon-components-svelte';
    import BottomPanelOpen from 'carbon-icons-svelte/lib/BottomPanelOpen.svelte';
    import type { components } from '$lib/sdk/schema';
    import { newClient } from '$lib/utils/client';
    import { handleGetError } from '$lib/utils/error-handler';
    import DetailModal from '$lib/components/spot-listings/detail-modal.svelte';
    import type { AddressResult } from '$lib/types/address/address';
    import {
        DEFAULT_DISTANCE,
        INIT_ZOOM,
        MAX_DISTANCE_RADIUS,
        MAX_ZOOM,
        MIN_DISTANCE_RADIUS,
        SELECTED_ZOOM,
        DISTANCE_RADIUS_STEP
    } from '$lib/constants';

    type ParkingSpot = components['schemas']['ParkingSpot'];

    const apiKey = import.meta.env.VITE_GEOCODING_API_KEY;
    const maxZoom = MAX_ZOOM;
    const defaultDistance = DEFAULT_DISTANCE;
    let distance = defaultDistance;

    let initZoom: number = INIT_ZOOM;
    const selectedZoom: number = SELECTED_ZOOM;
    const offset: [number, number] = [0, -10];

    let mapCenter: [number, number] = [0, 0];

    let searchQuery = '';
    let results: AddressResult[];
    let dropdownOpen: boolean = false;
    let selectedListingId: string | null = null;
    let showBackToTop: boolean = false;
    let listingsContainer: HTMLDivElement;

    let modalOpen: boolean = false;
    let selectedListing: ParkingSpot;

    let debounceTimer: number | undefined;

    let searchUsed: boolean = false;

    let client = newClient();
    let spotsData: ParkingSpot[] = [];

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

    const handleSelect = (location: AddressResult) => {
        searchQuery = location.properties.formatted;
        mapCenter = [location.geometry.coordinates[0], location.geometry.coordinates[1]];
        initZoom = selectedZoom;
        results = [];
        dropdownOpen = false;
        if (searchUsed === false) {
            searchUsed = true;
        }
        fetchSpots(location.geometry.coordinates, distance);
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

    const handleListingClick = (newSelectedListing: ParkingSpot) => {
        modalOpen = true;
        selectedListing = newSelectedListing;
    };

    const fetchSpots = async (coordinates: number[], radius: number) => {
        const { data: spots, error: errorSpots } = await client.GET('/spots', {
            params: {
                query: {
                    latitude: coordinates[1],
                    longitude: coordinates[0],
                    distance: radius
                }
            }
        });
        handleGetError(errorSpots);
        spotsData = spots ?? [];
    };

    const handleRadiusChange = (radius: number) => {
        distance = radius;
        if (mapCenter[0] !== 0 && mapCenter[1] !== 0) {
            fetchSpots(mapCenter, distance);
        }
    };

    onDestroy(() => {
        clearTimeout(debounceTimer);
    });
</script>

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

        {#if results?.length > 0 && dropdownOpen}
            <ul class="dropdown">
                {#each results as result}
                    <!-- svelte-ignore a11y-click-events-have-key-events a11y-no-noninteractive-element-interactions -->
                    <li on:click={() => handleSelect(result)}>
                        {result.properties.formatted}
                    </li>
                {/each}
            </ul>
        {/if}

        <div class="slider-container">
            <Slider
                value={distance}
                min={MIN_DISTANCE_RADIUS}
                max={MAX_DISTANCE_RADIUS}
                step={DISTANCE_RADIUS_STEP}
                minLabel="100m"
                maxLabel="5km"
                labelText="Distance Radius (metres)"
                on:change={(event) => handleRadiusChange(event.detail)}
            />
        </div>

        {#if spotsData.length > 0}
            {#each spotsData as listing}
                <!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
                <div
                    class="booking-info-container {listing.id === selectedListingId
                        ? 'highlight'
                        : ''}"
                    id={`listing-${listing.id}`}
                    on:click={() => handleListingClick(listing)}
                >
                    <SpotsListComponent {listing} />
                </div>
            {/each}
        {:else if searchUsed == false}
            <div class="empty-container">
                <h2>Search for your destination ↑</h2>
            </div>
        {:else}
            <div class="empty-container">
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
            zoom={initZoom}
            style="https://basemaps.cartocdn.com/gl/positron-gl-style/style.json"
        >
            <GeolocateControl position="bottom-left" fitBoundsOptions={{ maxZoom: maxZoom }} />
            {#each spotsData as { id, location }}
                <DefaultMarker
                    lngLat={[location.longitude ?? mapCenter[0], location.latitude ?? mapCenter[1]]}
                >
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

    {#if modalOpen}
        <DetailModal bind:open={modalOpen} bind:listing={selectedListing} />
    {/if}
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
        animation: smoothBlink 0.6s linear 3;
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

    .empty-container {
        display: flex;
        justify-content: center;
        margin: 1rem;
    }

    .slider-container :global(bx--form-item) {
        flex: 0 0 auto;
    }
</style>
