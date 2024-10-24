<script lang="ts">
    import type { PageData } from './$types';
    import Background from '$lib/images/background.png';
    import { Content, Checkbox, Button } from 'carbon-components-svelte';
    import { Edit, ArrowLeft, ArrowRight } from 'carbon-icons-svelte';
    import SpotEditModal from '$lib/components/spot-component/spot-edit-modal.svelte';
    import AvailabilityTable from '$lib/components/spot-component/availability-table.svelte';
    import { TimeSlotStatus } from '$lib/enum/timeslot-status';
    import { DAY_IN_A_WEEK, TOTAL_SEGMENTS_NUMBER } from '$lib/constants';
    import { getMonday, getDateWithDayOffset } from '$lib/utils/datetime-util';

    export let data: PageData;
    let is_edit_modal_open: boolean;
    let spot = data.spot;

    let availability_table: TimeSlotStatus[][] = Array.from({ length: TOTAL_SEGMENTS_NUMBER }, () =>
        Array(DAY_IN_A_WEEK).fill(0)
    );

    //reminder to change these to now()
    let today: Date;
    let currentMonday: Date;
    let nextMonday: Date;

    //This array contain all edit history
    let time_slot_edit_records = [];

    //These array contain exactly 7 days of data we want to display
    let time_slots_display = [];
    let time_slot_edit_display = [];

    today = new Date(Date.UTC(2024, 9, 15, 3, 24, 0));
    currentMonday = getMonday(today);
    nextMonday = getDateWithDayOffset(currentMonday, DAY_IN_A_WEEK);
    time_slots_display = extractRelevantTimeSlot(data.time_slots, currentMonday, nextMonday);
    updateTimeTable();

    function extractRelevantTimeSlot(full_time_slot: [], startDate: Date, endDate: Date) {
        let relevant_time_slots = [];
        full_time_slot?.forEach((slot) => {
            if (slot.date >= startDate && slot.date < endDate) {
                relevant_time_slots = [...relevant_time_slots, slot];
            }
        });
        return relevant_time_slots;
    }

    function updateTimeTable() {
        let current_seg = today.getUTCHours() * 2 + Math.ceil(today.getUTCMinutes() / 30);
        let cutoff_date = new Date(
            Date.UTC(today.getFullYear(), today.getUTCMonth(), today.getUTCDate())
        );

        console.log('cutoff: ' + cutoff_date.toUTCString());

        for (let day = 0; day < DAY_IN_A_WEEK; day++) {
            let currentDate = new Date(currentMonday);
            currentDate.setUTCDate(currentMonday.getUTCDate() + day);
            currentDate.setUTCHours(0, 0, 0, 0);
            for (let seg = 0; seg < TOTAL_SEGMENTS_NUMBER; seg++) {
                if (
                    currentDate.getTime() < cutoff_date.getTime() ||
                    (currentDate.getTime() == cutoff_date.getTime() && seg < current_seg)
                ) {
                    availability_table[seg][day] = TimeSlotStatus.PASTDUE;
                } else {
                    availability_table[seg][day] = TimeSlotStatus.NONE;
                }
            }
        }

        time_slots_display.forEach((slot) => {
            let day = (slot.date.getUTCDay() || 7) - 1;
            availability_table[slot.segment][day] = slot.status;
        });

        time_slot_edit_display.forEach((slot) => {
            let day = (slot.date.getUTCDay() || 7) - 1;
            availability_table[slot.segment][day] = slot.status;
        });
    }

    function toNextWeek() {
        currentMonday = getDateWithDayOffset(currentMonday, DAY_IN_A_WEEK);
        nextMonday = getDateWithDayOffset(nextMonday, DAY_IN_A_WEEK);
        time_slots_display = extractRelevantTimeSlot(data.time_slots, currentMonday, nextMonday);
        time_slot_edit_display = extractRelevantTimeSlot(
            time_slot_edit_records,
            currentMonday,
            nextMonday
        );
        console.log(time_slot_edit_display);
        updateTimeTable();
    }

    function toPrevWeek() {
        currentMonday = getDateWithDayOffset(currentMonday, -DAY_IN_A_WEEK);
        nextMonday = getDateWithDayOffset(nextMonday, -DAY_IN_A_WEEK);
        time_slots_display = extractRelevantTimeSlot(data.time_slots, currentMonday, nextMonday);
        time_slot_edit_display = extractRelevantTimeSlot(
            time_slot_edit_records,
            currentMonday,
            nextMonday
        );
        console.log(time_slot_edit_display);
        updateTimeTable();
    }

    function handleSubmit(event: Event) {
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
            console.log('Can not modify this slot');
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
            let match_display_index = time_slot_edit_display.findIndex((slot) => {
                return (
                    slot.date.getTime() == new_time_slot.date.getTime() &&
                    slot.segment == new_time_slot.segment
                );
            });
            if (match_display_index !== -1) {
                time_slot_edit_display = [
                    ...time_slot_edit_display.slice(0, match_display_index),
                    ...time_slot_edit_display.slice(match_display_index + 1)
                ];

                let match_record_index = time_slot_edit_records.findIndex((slot) => {
                    return (
                        slot.date.getTime() == new_time_slot.date.getTime() &&
                        slot.segment == new_time_slot.segment
                    );
                });
                time_slot_edit_records = [
                    ...time_slot_edit_records.slice(0, match_record_index),
                    ...time_slot_edit_records.slice(match_record_index + 1)
                ];
            } else {
                time_slot_edit_display = [...time_slot_edit_display, new_time_slot];
                time_slot_edit_records = [...time_slot_edit_records, new_time_slot];
            }

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

    <div class="date-controller">
        <p>
            From {currentMonday.toUTCString()} to {nextMonday.toUTCString()}
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

    <AvailabilityTable bind:availability_table on:edit={(e) => handleEdit(e)} />
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

    .date-controller {
        margin: 1rem;
    }
</style>
