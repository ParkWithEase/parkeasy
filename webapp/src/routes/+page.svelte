<script lang="ts">
    import Navbar from '$lib/components/navbar.svelte';
    import { MapLibre } from 'svelte-maplibre';
</script>

<Navbar />

<body>
    <div class="container">
        <div class="listings">
            <h2>Listings go here!</h2>
        </div>
        <div class="map-view">
            <MapLibre
                center={[-97.1, 49.9]}
                zoom={10}
                class="map"
                standardControls
                style="https://basemaps.cartocdn.com/gl/positron-gl-style/style.json"
            />
        </div>
    </div>
</body>

<style>
    .container {
        display: flex;
        flex-direction: row;
        margin: 8.5rem 3rem;
        max-height: fit-content;
    }

    .listings {
        display: flex;
        width: 30%;
        min-height: 100%;
        flex-direction: column;
        font-size: 1.2rem;
        font-weight: bold;
        align-items: center;
        justify-content: center;
        background-color: rgb(186, 214, 183);
    }

    .map-view {
        width: 75%;
        background-color: white;
    }

    :global(.map) {
        height: 75vh;
    }
</style>
