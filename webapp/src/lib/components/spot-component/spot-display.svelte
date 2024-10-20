<script lang="ts">
    import BackgroundImage from '$lib/images/background.png';
    import { Checkbox, Tag } from 'carbon-components-svelte';
    export let spot;
</script>

<a href={`/spots/${spot.id}/spot-info`}>
    <div class="spot-info">
        <div>
            <img src={BackgroundImage} class="spot-image" alt="A parking spot" />
        </div>
        <div>
            <div>
                <p class="spot-title">
                    {spot.location.street_address}, {spot.location.postal_code}
                    {#if spot.isListed}
                        <Tag type="green" style="font-size: 1rem;">listed</Tag>
                    {/if}
                </p>
                <span style="font-size: 1rem"
                    >{spot.location.city}, {spot.location.country_code}</span
                >
            </div>
            <div>
                <Checkbox labelText="Shelter" checked={spot.features.shelter} readonly />
                <Checkbox labelText="Plug-in" checked={spot.features.plug_in} readonly />
                <Checkbox
                    labelText="Charing station"
                    checked={spot.features.charging_station}
                    readonly
                />
            </div>
        </div>
    </div>
</a>

<style>
    .spot-info {
        padding-left: 1rem;
        flex: 1;
        gap: 1rem;
        display: flex;
        align-items: center;
        flex-direction: row;
        border: 1px solid #cfcfcfcf;
    }

    .spot-image {
        max-height: 10rem;
    }

    .spot-title {
        line-height: 1;
        font-size: 1.5rem;
        padding-top: 0.5rem;
    }

    a {
        all: unset;
    }
</style>
