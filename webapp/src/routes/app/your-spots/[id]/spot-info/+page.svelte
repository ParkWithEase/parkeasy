<script lang="ts">
    import type { PageData } from './$types';
    import Background from '$lib/images/background.png';
    import { Content, Checkbox, Button } from 'carbon-components-svelte';
    import { Edit } from 'carbon-icons-svelte';
    import SpotEditModal from '$lib/components/spot-component/spot-edit-modal.svelte';
    import AvailabilityTable from '$lib/components/spot-component/availability-table.svelte';
    import type { TimeSlotStatus } from '$lib/enum/timeslot-status';
    import { DAY_IN_A_WEEK, TOTAL_SEGMENTS_NUMBER } from '$lib/constants';

    export let data: PageData;
    let is_edit_modal_open: boolean;
    let spot = data.spot;
    let time_slots = data.time_slots;
    let time_slot_edit = [];
    let status_table: TimeSlotStatus[][] = Array.from({ length: TOTAL_SEGMENTS_NUMBER }, () =>
        Array(DAY_IN_A_WEEK).fill(0)
    );

    let currentMonday = getMonday(new Date(Date.now()));
    console.log(currentMonday);
    console.log(getNextMonday(currentMonday));

    function getMonday(d: Date) {
        let this_week_monday = new Date(d);
        this_week_monday.setSeconds(0);
        this_week_monday.setMinutes(0);
        this_week_monday.setHours(0);
        let day = this_week_monday.getDay() || 7;
        let diff = this_week_monday.getDate() - day + 1;
        this_week_monday.setDate(diff);
        return this_week_monday;
    }

    function getNextMonday(d: Date) {
        let next_week_monday = getMonday(d);
        next_week_monday.setDate(next_week_monday.getDate() + 7);
        return next_week_monday;
    }

    function handleSubmit(event: Event) {
        event.preventDefault();
        const formData = new FormData(event.target as HTMLFormElement);
        console.log(formData.get('shelter'));
        let new_spot = {
            id: spot.id,
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
        is_edit_modal_open = false;
        spot = new_spot;
    }
</script>

<Content>
    <img src={Background} class="spot-info-image" alt="spot" />
    <div class="spot-info-container">
        <p class="spot-info-header">Location</p>
        <div>
            <p class="spot-info-label">Street Address</p>
            <p class="spot-info-content">{spot?.location.street_address}</p>
        </div>
        <div>
            <p class="spot-info-label">City</p>
            <p class="spot-info-content">{spot?.location.city}</p>
        </div>
        <div>
            <p class="spot-info-label">Country Code</p>
            <p class="spot-info-content">{spot?.location.country_code}</p>
        </div>
        <div>
            <p class="spot-info-label">Postal Code</p>
            <p class="spot-info-content">{spot?.location.postal_code}</p>
        </div>

        <div>
            <p class="spot-info-label">Utilities</p>
            <Checkbox
                name="shelter"
                labelText="Shelter"
                checked={spot?.features.shelter}
                readonly
            />
            <Checkbox
                name="plug-in"
                labelText="Plug-in"
                checked={spot?.features.plug_in}
                readonly
            />
            <Checkbox
                name="charging-station"
                labelText="Charging Station"
                checked={spot?.features.charging_station}
                readonly
            />
        </div>
        <Button style="max-width: 4rem;" icon={Edit} on:click={() => (is_edit_modal_open = true)}
            >Edit</Button
        >
    </div>
    <SpotEditModal
        bind:openState={is_edit_modal_open}
        bind:spotInfo={spot}
        on:submit={handleSubmit}
    />

    <p class="spot-info-header">Availability</p>
    <AvailabilityTable bind:availability_table={status_table} />
</Content>

<style>
    .spot-info-image {
        max-height: 20rem;
    }

    .spot-info-header {
        font-size: 2rem;
    }

    .spot-info-label {
        font-size: 0.8rem;
    }

    .spot-info-content {
        font-size: 1.3rem;
    }

    .spot-info-container {
        display: flex;
        gap: 0.7rem;
        flex-direction: column;
        max-width: 20rem;
        margin-bottom: 1rem;
    }
</style>
