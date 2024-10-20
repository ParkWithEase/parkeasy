<script lang="ts">
    import { Button } from 'carbon-components-svelte';
    import type { PageData } from '../$types';
    import { Add } from 'carbon-icons-svelte';
    import SpotDisplay from '$lib/components/spot-component/spot-display.svelte';
    import SpotAddModal from '$lib/components/spot-component/spot-add-modal.svelte';

    let isAddModalOpen: boolean;
    console.log();
    export let data: PageData;

    function handleSubmit(event: Event) {
        event.preventDefault();
        const formData = new FormData(event.target as HTMLFormElement);
        console.log(formData.get('shelter'));
        let new_spot = {
            id: 'random',
            features: {
                charging_station: formData.get('charging-station') == null ? false : true,
                plug_in: formData.get('plug-in') == null ? false : true,
                shelter: formData.get('shelter') == null ? false : true
            },
            location: {
                city: formData.get('city'),
                country_code: formData.get('country-code'),
                latitude: 1,
                longitude: 1,
                postal_code: formData.get('postal-code'),
                street_address: formData.get('street-address')
            },
            isListed: false
        };
        data.spots = [...data.spots, new_spot];
        isAddModalOpen = false;
    }
</script>

<Button style="margin: 1rem;" icon={Add} on:click={() => (isAddModalOpen = true)}>New Spot</Button>
{#key data.spots}
    {#each data?.spots as spot}
        <SpotDisplay {spot} />
    {/each}
{/key}

{#if isAddModalOpen}
    <SpotAddModal bind:openState={isAddModalOpen} on:submit={handleSubmit} />
{/if}
