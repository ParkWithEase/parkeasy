<script lang="ts">
    import type { PageData } from './$types';
    import Background from '$lib/images/background.png';
    import { Content, Checkbox, Button } from 'carbon-components-svelte';
    import { Edit } from 'carbon-icons-svelte';
    import SpotEditModal from '$lib/components/spot-component/spot-edit-modal.svelte';
    import AvailabilityTable from '$lib/components/spot-component/availability-table.svelte';
    import { TimeSlotStatus } from '$lib/enum/timeslot-status';
    import { DAY_IN_A_WEEK, TOTAL_SEGMENTS_NUMBER } from '$lib/constants';
    import { getMonday, getNextMonday, getPreviousMonday } from '$lib/utils/datetime-util';

    export let data: PageData;
    let is_edit_modal_open: boolean;
    let spot = data.spot;

    let availability_table: TimeSlotStatus[][] = Array.from({ length: TOTAL_SEGMENTS_NUMBER }, () =>
        Array(DAY_IN_A_WEEK).fill(0)
    );

    //reminder to change these to now()
    let today = new Date('2024-10-15T03:24:00');
    let currentMonday = getMonday(today);
    let nextMonday = getNextMonday(today);
    let time_slots = [];
    let time_slot_edit = [];

    extractRelevantTimeSlot();
    updateTimeTable();

    function extractRelevantTimeSlot() {
        time_slots = [];
        data.time_slots?.forEach((slot) => {
            if (slot.date >= currentMonday && slot.date < nextMonday) {
                time_slots = [...time_slots, slot];
            }
        });
    }

    function updateTimeTable() {
        let week_day = (today.getDay() || 7) - 1;
        let current_seg = today.getHours() * 2 + Math.ceil(today.getMinutes() / 30);

        for (let day = 0; day < DAY_IN_A_WEEK; day++) {
            let currentDate = new Date(currentMonday);
            currentDate.setDate(currentMonday.getDate() + day);
            for (let seg = 0; seg < TOTAL_SEGMENTS_NUMBER; seg++) {
                if (day < week_day || (day == week_day && seg <= current_seg)) {
                    availability_table[seg][day] = TimeSlotStatus.PASTDUE;
                } else {
                    availability_table[seg][day] = TimeSlotStatus.NONE;
                }
            }
        }

        time_slots.forEach((slot) => {
            let day = slot.date.getDay() || 7 - 1;
            availability_table[slot.segment][day] = slot.status;
        });

        time_slot_edit.forEach((slot) => {
            let day = slot.date.getDay() || 7 - 1;
            availability_table[slot.segment][day] = slot.status;
        });
    }

    function toNextWeek() {
        currentMonday = getNextMonday(currentMonday);
        nextMonday = getNextMonday(nextMonday);
        extractRelevantTimeSlot();
        updateTimeTable();
    }

    function toPrevWeek() {
        currentMonday = getPreviousMonday(currentMonday);
        nextMonday = getPreviousMonday(nextMonday);
        extractRelevantTimeSlot();
        updateTimeTable();
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

    <div>
        <p>From {currentMonday.toDateString()} at 00:00 to {nextMonday.toDateString()} at 00:00</p>
        <Button on:click={toPrevWeek}>Prev Week</Button>
        <Button on:click={toNextWeek}>Next Week</Button>
    </div>

    <AvailabilityTable bind:availability_table />
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
