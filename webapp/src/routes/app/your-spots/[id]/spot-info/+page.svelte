<script lang="ts">
    import type { PageData } from './$types';
    import Background from '$lib/images/background.png';
    import {
        Content,
        Button,
        Form,
        ToastNotification,
        NumberInput,
        Checkbox
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
    import SpotInfoEditable from '$lib/components/spot-component/spot-info-editable.svelte';

    export let data: PageData;
    type TimeUnit = components['schemas']['TimeUnit'];

    //Variable for location edit section
    let spot = data.spot;
    let availabilityTimeSlot = data.time_slots;

    //temporarily use id from params

    //These are Variables for availability edit section

    let newPricePerHour: number = spot?.price_per_hour || 0;
    let spotChargingStation: boolean = spot?.features?.charging_station ?? false;
    let spotPlugIn: boolean = spot?.features?.plug_in ?? false;
    let spotShelter: boolean = spot?.features?.shelter ?? false;
    //This array contain all edit history

    let client = newClient();

    //For ToastMessage
    let errorTimeOut: number = 0;
    let successTimeOut: number = 0;
    let errorMessage: string = '';
    $: showToast = errorTimeOut !== 0;
    $: showSuccess = successTimeOut !== 0;

    //reminder to change these to now()
    let today = new Date(Date.now());
    let currentMonday = getMonday(today);
    let nextMonday: Date;

    //updates
    let priceUtilitiesUpdated: boolean = false;

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
                        errorTimeOut = ERROR_MESSAGE_TIME_OUT;
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

    async function updatePriceUtilities(event:Event) {
        event.preventDefault();
        priceUtilitiesUpdated = true;
    }

    async function handleUpdateAvailability(event: Event) {
        event.preventDefault();

        if (!checkAvailabilityChange() && !priceUtilitiesUpdated) {
            errorMessage = "No changes detected.";
            errorTimeOut = ERROR_MESSAGE_TIME_OUT;
            return;
        }

        if(priceUtilitiesUpdated){
            client
                .PUT('/spots/{id}', {
                    params: {
                        path: {
                            id: spot?.id ?? '0'
                        }
                    },
                    headers: { 'Content-Type': 'application/json' },
                    body: {
                        features: {"charging_station": spotChargingStation,
                                    "plug_in": spotPlugIn,
                                    "shelter": spotShelter},
                        price_per_hour: newPricePerHour
                    },
                })
                .then(({ data, error }) => {
                    if (data) {
                        successTimeOut = ERROR_MESSAGE_TIME_OUT;
                        priceUtilitiesUpdated = false;
                    }
                    if (error) {
                        errorMessage = getErrorMessage(error);
                        errorTimeOut = ERROR_MESSAGE_TIME_OUT;
                    }
                })
                .catch((err) => {
                    errorMessage = "An error occurred while updating.";
                    errorTimeOut = ERROR_MESSAGE_TIME_OUT;
                });
        }

        let addAvailability: Array<{ end_time: string, start_time: string }> = [];
        let removeAvailability: Array<{ end_time: string, start_time: string }> = [];

        // Compare the initial and edited tables
        for (let [key, initialTable] of availabilityTablesInitial) {
            let editTable = availabilityTableEdit.get(key);
            const startOfWeek = new Date(key);

            for (let seg = 0; seg < TOTAL_SEGMENTS_NUMBER; seg++) {
                for (let day = 0; day < DAY_IN_A_WEEK; day++) {
                    const initialStatus = initialTable[seg][day];
                    const editedStatus = editTable?.[seg]?.[day] ?? TimeSlotStatus.NONE;

                    if (initialStatus !== editedStatus) {
                        const startDateTime = new Date(startOfWeek);
                        startDateTime.setDate(startOfWeek.getDate() + day);
                        startDateTime.setHours(Math.floor(seg / 2), (seg % 2) * 30, 0, 0);

                        const endDateTime = new Date(startDateTime);
                        endDateTime.setMinutes(endDateTime.getMinutes() + 30);

                        if (editedStatus === TimeSlotStatus.AVAILABLE) {
                            addAvailability.push({
                                start_time: startDateTime.toISOString(),
                                end_time: endDateTime.toISOString(),
                            });
                        } else if (initialStatus === TimeSlotStatus.AVAILABLE) {
                            removeAvailability.push({
                                start_time: startDateTime.toISOString(),
                                end_time: endDateTime.toISOString(),
                            });
                        }
                    }
                }
            }
        }

        client
            .PUT('/spots/{id}/availability', {
                params: {
                    path: {
                        id: spot?.id ?? '0'
                    }
                },
                headers: { 'Content-Type': 'application/json' },
                body: {
                    add_availability: addAvailability,
                    remove_availability: removeAvailability
                },
            })
            .then(({ data, error }) => {
                if (data) {
                    availabilityTablesInitial = structuredClone(availabilityTableEdit);
                    successTimeOut = ERROR_MESSAGE_TIME_OUT;
                }
                if (error) {
                    errorMessage = getErrorMessage(error);
                    errorTimeOut = ERROR_MESSAGE_TIME_OUT;
                }
            })
            .catch((err) => {
                errorMessage = "An error occurred while updating availability.";
                errorTimeOut = ERROR_MESSAGE_TIME_OUT;
            });
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

    <img src={Background} class="spot-info-image" alt="spot" />
    <p class="spot-info-header">Location</p>
    <SpotInfoEditable bind:spot />

    <div>
        <p class="spot-info-label">Utilities</p>
        <Checkbox
            name="shelter"
            labelText="Shelter"
            bind:checked={spotShelter}
            on:change={updatePriceUtilities}
        />
        <Checkbox
            name="plug-in"
            labelText="Plug-in"
            bind:checked={spotPlugIn}
            on:change={updatePriceUtilities}
        />
        <Checkbox
            name="charging-station"
            labelText="Charging Station"
            bind:checked={spotChargingStation}
            on:change={updatePriceUtilities}
        />
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

    <Form on:submit={handleUpdateAvailability}>
        <div class="price-field">
            <NumberInput
                label="Price per hour"
                hideSteppers
                step={0.01}
                min={0}
                name="price-per-hour"
                helperText="Price in CAD"
                required
                bind:value={newPricePerHour}
                on:change={updatePriceUtilities}
            />
        </div>
        <Button kind="secondary" on:click={resetAvailabilityEdit}>Reset</Button>
        <Button type="submit">Update</Button>
        {#if showToast}
            <div transition:fade class="error-message">
                <ToastNotification
                    bind:timeout={errorTimeOut}
                    kind="error"
                    fullWidth
                    title="Error"
                    subtitle={errorMessage}
                    on:close={() => {
                        errorTimeOut = 0;
                    }}
                />
            </div>
        {/if}
        {#if showSuccess}
            <div transition:fade class="success-message">
                <ToastNotification
                    bind:timeout={successTimeOut}
                    kind="success"
                    fullWidth
                    title="Update Successful!"
                    on:close={() => {
                        successTimeOut = 0;
                    }}
                />
            </div>
        {/if}
    </Form>
</Content>

<style>
    .spot-info-label {
        font-size: 0.8rem;
    }
</style>
