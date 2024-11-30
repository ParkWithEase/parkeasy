<script lang="ts">
    import type { PageData } from './$types';
    import Background from '$lib/images/background.png';
    import {
        Content,
        Button,
        TextInput,
        Form,
        ToastNotification,
        TooltipIcon
    } from 'carbon-components-svelte';
    import { TimeSlotStatus, TimeSlotStatusConverter } from '$lib/enum/timeslot-status';
    import {
        DAY_IN_A_WEEK,
        ERROR_MESSAGE_TIME_OUT,
        SPOT_PREFFERED_ICON_SIZE,
        TOTAL_SEGMENTS_NUMBER
    } from '$lib/constants';
    import { getMonday, getDateWithDayOffset } from '$lib/utils/datetime-util';
    import { createEmptyTable } from '$lib/utils/time-table-util';
    import { newClient } from '$lib/utils/client';
    import type { components } from '$lib/sdk/schema';
    import { getErrorMessage } from '$lib/utils/error-handler';
    import { fade } from 'svelte/transition';
    import AvailabilitySection from '$lib/components/spot-component/availability-section.svelte';
    import SpotInfo from '$lib/components/spot-component/spot-info.svelte';
    import { Favorite, FavoriteFilled } from 'carbon-icons-svelte';

    export let data: PageData;
    type TimeUnit = components['schemas']['TimeUnit'];

    //Variable for location edit section
    let spot = data.spot;
    let availabilityTimeSlot = data.time_slots;
    let spotPreferred: boolean = data.isPreferred;
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

    function setPreference() {
        if (!spotPreferred) {
            console.log('add preference');
            client
                .POST('/spots/{id}/preference', { params: { path: { id: spot?.id ?? '0' } } })
                .then(({ error }) => {
                    if (error) {
                        errorMessage = getErrorMessage(error);
                        toastTimeOut = ERROR_MESSAGE_TIME_OUT;
                    } else {
                        spotPreferred = true;
                    }
                })
                .catch((err) => {
                    errorMessage = err;
                });
        } else {
            console.log('remove preference');
            client
                .DELETE('/spots/{id}/preference', { params: { path: { id: spot?.id ?? '0' } } })
                .then(({ error }) => {
                    if (error) {
                        errorMessage = getErrorMessage(error);
                        toastTimeOut = ERROR_MESSAGE_TIME_OUT;
                    } else {
                        spotPreferred = false;
                    }
                })
                .catch((err) => {
                    errorMessage = err;
                });
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
    <div class="header-with-preference-option">
        <p class="spot-info-header">Location</p>
        <TooltipIcon
            tooltipText={spotPreferred ? 'Remove from preferred list' : 'Add to preferred list'}
            on:click={() => setPreference()}
        >
            <FavoriteFilled
                style={`${spotPreferred ? 'fill:deeppink;' : 'fill:None;'} position:absolute; `}
                size={SPOT_PREFFERED_ICON_SIZE}
            />
            <Favorite style="fill:black; position: absolute; " size={SPOT_PREFFERED_ICON_SIZE} />
        </TooltipIcon>
    </div>

    <SpotInfo bind:spot />

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
    .header-with-preference-option {
        display: flex;
        gap: 1rem;
        align-self: center;
    }
</style>
