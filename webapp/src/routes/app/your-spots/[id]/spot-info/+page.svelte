<script lang="ts">
    import type { PageData } from './$types';
    import Background from '$lib/images/background.png';
    import { Content, Checkbox, Button, TextInput, Form } from 'carbon-components-svelte';
    import { Edit, ArrowLeft, ArrowRight } from 'carbon-icons-svelte';
    import SpotEditModal from '$lib/components/spot-component/spot-edit-modal.svelte';
    import AvailabilityTable from '$lib/components/spot-component/availability-table.svelte';
    import { TimeSlotStatus } from '$lib/enum/timeslot-status';
    import { DAY_IN_A_WEEK, TOTAL_SEGMENTS_NUMBER } from '$lib/constants';
    import { getMonday, getDateWithDayOffset } from '$lib/utils/datetime-util';
    import { getWeekAvailabilityTable } from '$lib/utils/time-table-util';

    export let data: PageData;

    //Variable for location edit section
    let is_edit_modal_open: boolean;
    let spot = data.spot;

    //These are Variables for availability edit section
    let availability_table: TimeSlotStatus[][];
    let new_price_per_hour: number | undefined = data.spot?.price_per_hour;
    //This array contain all edit history
    let time_slot_edit_records = [];

    //reminder to change these to now()
    let today: Date;
    let currentMonday: Date;
    let nextMonday: Date;

    $: isAvailablilityChanged =
        time_slot_edit_records.length != 0 || new_price_per_hour !== spot?.price_per_hour;

    today = new Date(Date.UTC(2024, 9, 15, 3, 24, 0));
    currentMonday = getMonday(today);
    nextMonday = getDateWithDayOffset(currentMonday, DAY_IN_A_WEEK);
    availability_table = getWeekAvailabilityTable(
        today,
        currentMonday,
        data.time_slots,
        time_slot_edit_records
    );

    function toNextWeek() {
        currentMonday = getDateWithDayOffset(currentMonday, DAY_IN_A_WEEK);
        nextMonday = getDateWithDayOffset(nextMonday, DAY_IN_A_WEEK);
        availability_table = getWeekAvailabilityTable(
            today,
            currentMonday,
            data.time_slots,
            time_slot_edit_records
        );
    }

    function toPrevWeek() {
        currentMonday = getDateWithDayOffset(currentMonday, -DAY_IN_A_WEEK);
        nextMonday = getDateWithDayOffset(nextMonday, -DAY_IN_A_WEEK);
        availability_table = getWeekAvailabilityTable(
            today,
            currentMonday,
            data.time_slots,
            time_slot_edit_records
        );
    }

    function handleSubmitLocation(event: Event) {
        event.preventDefault();
        const formData = new FormData(event.target as HTMLFormElement);
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
                state: formData.get('state'),
                latitude: 1,
                longitude: 1,
                postal_code: formData.get('postal-code'),
                street_address: formData.get('street-address')
            },
            isListed: false,
            price_per_hour: spot?.price_per_hour
        };
        is_edit_modal_open = false;
        spot = new_spot;
    }

    function handleSubmitAvailability(event: Event) {
        event.preventDefault();
        //TODO: change parking spot price + availability using the edit records.
        // Clear the edit records on success
        // Might need to check if anything actually change and
        const formData = new FormData(event.target as HTMLFormElement);
        console.log(formData.get('price-per-hour'));
        spot.price_per_hour = formData.get('price-per-hour');
        new_price_per_hour = spot?.price_per_hour;
    }

    function resetAvailabilityEdit() {
        time_slot_edit_records = [];
        new_price_per_hour = spot?.price_per_hour;
        availability_table = getWeekAvailabilityTable(
            today,
            currentMonday,
            data.time_slots,
            time_slot_edit_records
        );
    }

    /*
    this function should remove the edit event if there is already one event at that time slot. Or 
    append the new event if no event happens to that time slot
    */
    function handleEdit(event: Event) {
        if (
            [TimeSlotStatus.BOOKED, TimeSlotStatus.PASTDUE, TimeSlotStatus.AUCTIONED].includes(
                event.detail.status
            )
        ) {
            return;
        } else {
            let date = new Date(currentMonday);
            date.setUTCDate(currentMonday.getUTCDate() + event.detail.day);
            let new_time_slot = {
                date: date,
                segment: event.detail.segment,
                status: event.detail.status
            };
            if (event.detail.status === TimeSlotStatus.NONE) {
                new_time_slot.status = TimeSlotStatus.AVAILABLE;
            } else {
                new_time_slot.status = TimeSlotStatus.NONE;
            }

            let match_record_index = time_slot_edit_records.findIndex((slot) => {
                return (
                    slot.date.getTime() == new_time_slot.date.getTime() &&
                    slot.segment == new_time_slot.segment
                );
            });

            if (match_record_index !== -1) {
                time_slot_edit_records.splice(match_record_index, 1);
            } else {
                time_slot_edit_records.push(new_time_slot);
            }

            //trigger reactivity by reassigning object
            time_slot_edit_records = time_slot_edit_records;

            availability_table[event.detail.segment][event.detail.day] = new_time_slot.status;
        }
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
            <p class="spot-info-label">State/Provice</p>
            <p class="spot-info-content">{spot?.location.state}</p>
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
                style="pointer-events: none;"
                readonly
            />
            <Checkbox
                name="plug-in"
                labelText="Plug-in"
                checked={spot?.features.plug_in}
                style="pointer-events: none;"
                readonly
            />
            <Checkbox
                name="charging-station"
                labelText="Charging Station"
                checked={spot?.features.charging_station}
                style="pointer-events: none;"
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
        on:submit={handleSubmitLocation}
    />

    <p class="spot-info-header">Availability</p>

    <div class="date-controller">
        <p>
            From {currentMonday?.toUTCString()} to {nextMonday?.toUTCString()}
        </p>
        <Button
            kind="secondary"
            iconDescription="Last Week"
            size="small"
            on:click={toPrevWeek}
            icon={ArrowLeft}>Last Week</Button
        >
        <Button
            size="small"
            iconDescription="Next Week"
            kind="secondary"
            on:click={toNextWeek}
            icon={ArrowRight}>Next Week</Button
        >
    </div>

    <Form on:submit={handleSubmitAvailability}>
        <AvailabilityTable
            bind:availability_table
            on:edit={(e) => {
                handleEdit(e);
            }}
        />
        <div class="price-field">
            <TextInput
                style="max-width: 6rem;"
                labelText="Price per hour"
                name="price-per-hour"
                helperText="Price in CAD"
                type="number"
                required
                bind:value={new_price_per_hour}
            />
        </div>
        {#if isAvailablilityChanged}
            <Button kind="secondary" on:click={resetAvailabilityEdit}>Reset</Button>
            <Button type="submit">Submit</Button>
        {/if}
    </Form>
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
