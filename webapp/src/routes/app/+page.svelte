<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import Navbar from '$lib/components/navbar.svelte';
    import { MapLibre, DefaultMarker, Popup, GeolocateControl } from 'svelte-maplibre';
    import SpotsListComponent from '$lib/components/spot-listings/spots-list-component.svelte';
    import { spots_data } from './your-spots/mock_data';
    import { Search } from 'carbon-components-svelte';

    const apiKey = import.meta.env.VITE_API_KEY;

    let mapCenter: [number, number] = [-97.1, 49.9];
    let zoom = 10;

    let searchQuery = '';
    let results = [];
    let dropdownOpen = false;
    let selectedListingId: string | null = null; // to track the selected listing for highlighting

    let showBackToTop = false;
    let listingsContainer: HTMLDivElement;

    const searchLocation = async () => {
        const url = `https://api.geoapify.com/v1/geocode/search?text=${encodeURIComponent(searchQuery)}&apiKey=${apiKey}`;

        try {
            const response = await fetch(url);
            const data = await response.json();
            results = data.features || [];
            dropdownOpen = true;
        } catch (error) {
            console.error('Error fetching data:', error);
        }
    };

    const handleSelect = (location: {
        properties: { formatted: string };
        geometry: { coordinates: number[] };
    }) => {
        searchQuery = location.properties.formatted;
        mapCenter = [location.geometry.coordinates[0], location.geometry.coordinates[1]];
        zoom = 13;
        results = [];
        dropdownOpen = false;
    };

    $: if (searchQuery && searchQuery.length > 3) {
        searchLocation();
    }

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
            listingsContainer.removeEventListener('scroll', handleScroll);
        };
    });

    onDestroy(() => {
        document.removeEventListener('click', handleClickOutside);
        listingsContainer.removeEventListener('scroll', handleScroll);
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
    <div class="listings" bind:this={listingsContainer}>
        <Search
            class="search-input"
            type="text"
            bind:value={searchQuery}
            placeholder="Explore spots around this area..."
            name="search"
            on:input={() => {
                dropdownOpen = true;
            }}
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
            <GeolocateControl position="bottom-left" fitBoundsOptions={{ maxZoom: 12 }} />
            {#each spots_data as { id, location }}
                <DefaultMarker
                    lngLat={[location.longitude, location.latitude]}
                    on:click={() => handleMarkerClick(id)}
                >
                    <Popup offset={[0, -10]}>
                        <div class="text-lg font-bold">
                            {location.street_address}, {location.city}
                        </div>
                        <button on:click={() => handleMarkerClick(id)}>Go to listing</button>
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
        background-color: rgb(186, 214, 183);
        position: fixed;
        z-index: 999;
        border: 1px solid rgb(104, 104, 104);
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
            background-color: #49da55;
            opacity: 0.5;
        }
    }

    .highlight {
        animation: smoothBlink 0.75s linear 2;
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
</style>
