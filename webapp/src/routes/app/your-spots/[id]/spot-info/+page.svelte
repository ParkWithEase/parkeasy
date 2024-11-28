<script lang="ts">
    import type { PageData } from './$types';
    import Background from '$lib/images/background.png';
    import {
        Content,
        Checkbox,
        Button,
        TextInput,
        Form,
        ToastNotification
    } from 'carbon-components-svelte';
    import { TimeSlotStatus, TimeSlotStatusConverter } from '$lib/enum/timeslot-status';
    import { DAY_IN_A_WEEK, ERROR_MESSAGE_TIME_OUT, TOTAL_SEGMENTS_NUMBER } from '$lib/constants';
    import { getMonday, getDateWithDayOffset } from '$lib/utils/datetime-util';
    import { createEmptyTable } from '$lib/utils/time-table-util';
    import { newClient } from '$lib/utils/client';
    import type { components } from '$lib/sdk/schema';
    import { getErrorMessage } from '$lib/utils/error-handler';
    import { fade } from 'svelte/transition';
    import AvailabilitySection from '$lib/components/spot-component/availability-section.svelte';

    export let data: PageData;
    type TimeUnit = components['schemas']['TimeUnit'];

    //Variable for location edit section
    let spot = data.spot;
    let availabilityTimeSlot = data.time_slots;

    //temporarily use id from params

    //These are Variables for availability edit section

    let newPricePerHour: number = spot?.price_per_hour || 0;
    //This array contain all edit history

    let client = newClient();

    //For ToastMessage
    let toastTimeOut: number = 0;
    let errorMessage: string = '';
    $: showToast = toastTimeOut !== 0;

    //reminder to change these to now()
    let today = new Date(Date.now());
    let currentMonday = getMonday(today);
    let nextMonday: Date;

    let availabilityTablesInitial = new Map<number, TimeSlotStatus[][]>();
    availabilityTablesInitial.set(
        currentMonday.getTime(),
        createIntialAvailabilityTable(availabilityTimeSlot ?? [], today, currentMonday)
    );
    let availabilityTableEdit = structuredClone(availabilityTablesInitial);
    let availabilityTable: TimeSlotStatus[][];

    $: nextMonday = getDateWithDayOffset(currentMonday, DAY_IN_A_WEEK);
    $: availabilityTable = availabilityTableEdit.get(currentMonday.getTime()) || [];

    function createIntialAvailabilityTable(
        timeUnits: TimeUnit[],
        today: Date,
        currentMonday: Date
    ) {
        let newTable = createEmptyTable();
        timeUnits.forEach((timeUnit) => {
            const date = new Date(timeUnit.start_time);
            const day = (date.getDay() || 7) - 1;
            const segment = date.getHours() * 2 + Math.ceil(date.getMinutes() / 30);
            newTable[segment][day] = TimeSlotStatusConverter[timeUnit.status || 'none'];
        });

        const current_seg = today.getHours() * 2 + Math.floor(today.getMinutes() / 30);
        const cutoff_date = new Date(
            new Date(today.getFullYear(), today.getMonth(), today.getDate())
        );

        for (let day = 0; day < DAY_IN_A_WEEK; day++) {
            const currentDate = new Date(currentMonday);
            currentDate.setDate(currentMonday.getDate() + day);
            currentDate.setHours(0, 0, 0, 0);
            for (let seg = 0; seg < TOTAL_SEGMENTS_NUMBER; seg++) {
                if (
                    currentDate.getTime() < cutoff_date.getTime() ||
                    (currentDate.getTime() == cutoff_date.getTime() && seg <= current_seg)
                ) {
                    newTable[seg][day] = TimeSlotStatus.PASTDUE;
                }
            }
        }
        return newTable;
    }

    function createAvailabilityTable(timeUnits: TimeUnit[]) {
        let newTable = createEmptyTable();
        timeUnits.forEach((timeUnit) => {
            const date = new Date(timeUnit.start_time);
            const day = (date.getDay() || 7) - 1;
            const segment = date.getHours() * 2 + Math.ceil(date.getMinutes() / 30);
            newTable[segment][day] = TimeSlotStatusConverter[timeUnit.status || 'none'];
        });
        return newTable;
    }

    function toNextWeek() {
        currentMonday = getDateWithDayOffset(currentMonday, DAY_IN_A_WEEK);
        const nextMonday = getDateWithDayOffset(currentMonday, DAY_IN_A_WEEK);
        nextMonday.setMinutes(nextMonday.getMinutes() - 30);
        if (!availabilityTablesInitial.get(currentMonday.getTime())) {
            client
                .GET('/spots/{id}/availability', {
                    params: {
                        path: {
                            id: spot?.id ?? '0'
                        },
                        query: {
                            availability_start: currentMonday.toISOString(),
                            availability_end: nextMonday.toISOString()
                        }
                    }
                })
                .then(({ data: availability, error }) => {
                    if (availability) {
                        let newTable = createAvailabilityTable(availability);
                        availabilityTablesInitial.set(currentMonday.getTime(), newTable);
                        availabilityTableEdit.set(
                            currentMonday.getTime(),
                            structuredClone(newTable)
                        );
                        availabilityTable =
                            availabilityTableEdit.get(currentMonday.getTime()) ?? [];
                    }
                    if (error) {
                        errorMessage = getErrorMessage(error);
                        toastTimeOut = ERROR_MESSAGE_TIME_OUT;
                    }
                })
                .catch((err) => {
                    errorMessage = err;
                });
        }
    }

    function toPrevWeek() {
        if (availabilityTable?.[0]?.[0] === TimeSlotStatus.PASTDUE) {
            return;
        }
        currentMonday = getDateWithDayOffset(currentMonday, -DAY_IN_A_WEEK);
    }

    function checkAvailabilityChange() {
        for (let [key, initialTable] of availabilityTablesInitial) {
            let editTable = availabilityTableEdit.get(key);
            for (let seg = 0; seg < TOTAL_SEGMENTS_NUMBER; seg++) {
                for (let day = 0; day < DAY_IN_A_WEEK; day++) {
                    if (initialTable[seg][day] != editTable?.[seg][day]) {
                        return true;
                    }
                }
            }
        }
        return false;
    }

    function handleSubmitAvailability(event: Event) {
        event.preventDefault();
        //TODO: change parking spot price + availability using the edit records.
        // Clear the edit records on success
        // Might need to check if anything actually change and
        console.log(checkAvailabilityChange());
    }

    function resetAvailabilityEdit() {
        newPricePerHour = spot?.price_per_hour || 0;
        availabilityTableEdit = structuredClone(availabilityTablesInitial);
    }

    /*
    this function  remove the edit event if there is already one event at that time slot. Or 
    append the new event if no event happens to that time slot
    */
    function handleEdit(event: CustomEvent) {
        if (
            [TimeSlotStatus.BOOKED, TimeSlotStatus.PASTDUE, TimeSlotStatus.AUCTIONED].includes(
                event.detail.status
            )
        ) {
            return;
        } else {
            availabilityTable[event.detail.segment] ??= [];
            availabilityTable[event.detail.segment][event.detail.day] =
                event.detail.status === TimeSlotStatus.NONE
                    ? TimeSlotStatus.AVAILABLE
                    : TimeSlotStatus.NONE;
            availabilityTableEdit.set(currentMonday.getTime(), availabilityTable);
            availabilityTableEdit = availabilityTableEdit;
        }
    }
</script>

<Content>
    {#if showToast}
        <div transition:fade class="error-message">
            <ToastNotification
                bind:timeout={toastTimeOut}
                kind="error"
                fullWidth
                title="Error"
                subtitle={errorMessage}
                on:close={() => {
                    toastTimeOut = 0;
                }}
            />
        </div>
    {/if}

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
                checked={spot?.features?.shelter}
                style="pointer-events: none;"
                readonly
            />
            <Checkbox
                name="plug-in"
                labelText="Plug-in"
                checked={spot?.features?.plug_in}
                style="pointer-events: none;"
                readonly
            />
            <Checkbox
                name="charging-station"
                labelText="Charging Station"
                checked={spot?.features?.charging_station}
                style="pointer-events: none;"
                readonly
            />
        </div>
    </div>

    <p class="spot-info-header">Availability</p>

    <AvailabilitySection
        bind:currentMonday
        bind:nextMonday
        bind:availabilityTable
        {toPrevWeek}
        {toNextWeek}
        {handleEdit}
    />

    <Form on:submit={handleSubmitAvailability}>
        <div class="price-field">
            <TextInput
                labelText="Price per hour"
                name="price-per-hour"
                helperText="Price in CAD"
                type="number"
                required
                bind:value={newPricePerHour}
            />
        </div>
        <Button kind="secondary" on:click={resetAvailabilityEdit}>Reset</Button>
        <Button type="submit">Submit</Button>
    </Form>
</Content>

<style>
    .spot-info-image {
        max-height: 20rem;
        max-width: 100%;
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
