<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import Navbar from '$lib/components/navbar.svelte';
    import { MapLibre, DefaultMarker, Popup, GeolocateControl } from 'svelte-maplibre';
    import SpotsListComponent from '$lib/components/spot-listings/spots-list-component.svelte';
    import { Button } from "carbon-components-svelte";
    import { spots_data } from './your-spots/mock_data';
    import { Search } from 'carbon-components-svelte';
    import type { LocationResult } from '$lib/types/spot/location-result';
    import BottomPanelOpen from "carbon-icons-svelte/lib/BottomPanelOpen.svelte";

    const apiKey = import.meta.env.VITE_GEOCODING_API_KEY;
    const maxZoom : number = 12;
    const initZoom : number= 13;
    const offset : [number,number] = [0, -10];

    let mapCenter: [number, number] = [-97.1, 49.9];
    let zoom : number= 10;

    let searchQuery = '';
    let results: LocationResult[] = [];
    let dropdownOpen: boolean = false;
    let selectedListingId: string | null = null; 

    let showBackToTop : boolean = false;
    let listingsContainer: HTMLDivElement;

    let debounceTimer: number | undefined; 

    const searchLocation = () => {
        const url = `https://api.geoapify.com/v1/geocode/search?text=${encodeURIComponent(searchQuery)}&apiKey=${apiKey}`;

        fetch(url)
            .then(response => response.json())
            .then(data => {
                results = data.features ?? [];
                dropdownOpen = true;
            })
            .catch(error => {
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
        }, 750); // 0.75s delay
    };

    const handleSelect = (location: {
        properties: { formatted: string };
        geometry: { coordinates: number[] };
    }) => {
        searchQuery = location.properties.formatted;
        mapCenter = [location.geometry.coordinates[0], location.geometry.coordinates[1]];
        zoom = initZoom;
        results = [];
        dropdownOpen = false;
    };

    const handleClickOutside = (event: Event) => {
        if (
            !(event.target as HTMLElement).closest('.dropdown') &&
            !(event.target as HTMLElement).closest('.search-input')
        ) {
            dropdownOpen = false;
        }
    };

    onMount(() => {
        document.addEventListener('click', handleClickOutside);
        listingsContainer.addEventListener('scroll', handleScroll);
        return () => {
            document.removeEventListener('click', handleClickOutside);
        };
    });

    onDestroy(() => {
        document.removeEventListener('click', handleClickOutside);
        clearTimeout(debounceTimer);
    });

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
</script>

<Navbar />

<div class="container">
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

        {#each spots_data as listing}
            <div
                class="booking-info-container {listing.id === selectedListingId ? 'highlight' : ''}"
                id={`listing-${listing.id}`}
            >
                <SpotsListComponent {listing} />
            </div>
        {/each}

        {#if showBackToTop}
            <button class="back-to-top" on:click={scrollToTop}> ↑ Back to Top </button>
        {/if}
    </div>

    <div class="map-view">
        <MapLibre
            center={mapCenter}
            {zoom}
            style="https://basemaps.cartocdn.com/gl/positron-gl-style/style.json"
        >
            <GeolocateControl position="bottom-left" fitBoundsOptions={{ maxZoom: maxZoom }} />
            {#each spots_data as { id, location }}
                <DefaultMarker
                    lngLat={[location.longitude, location.latitude]}
                >
                    <Popup offset={offset}>
                        <div class = 'popup-container'>
                            <h2 class = 'popup-text'>
                                {location.street_address} <br>
                                {location.city},{location.state} <br>
                                {location.postal_code} <br>
                            </h2>
                            <Button kind="secondary" on:click={() => handleMarkerClick(id)} icon={BottomPanelOpen}>
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
        max-height: 82vh;
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

    .popup-container{
        display: flex;
        flex-direction: column;
        padding: 0.5rem
    }

    .popup-text{
        margin: 1rem;
        font-family: Fredoka, sans-serif;
    }
    
</style>
