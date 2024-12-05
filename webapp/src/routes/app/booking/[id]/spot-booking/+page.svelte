<script lang="ts">
    import type { PageData } from './$types';
    import Background from '$lib/images/background.png';
    import {
        Content,
        Button,
        TextInput,
        Form,
        ToastNotification,
        TooltipIcon,
        FormGroup,
        NumberInput
    } from 'carbon-components-svelte';
    import { TimeSlotStatus, TimeSlotStatusConverter } from '$lib/enum/timeslot-status';
    import {
        BOOK_WITHOUT_CAR,
        BOOK_WITHOUT_SLOTS_ERROR,
        DAY_IN_A_WEEK,
        ERROR_MESSAGE_TIME_OUT,
        SPOT_PREFFERED_ICON_SIZE,
        TOTAL_SEGMENTS_NUMBER
    } from '$lib/constants';
    import { getMonday, getDateWithDayOffset } from '$lib/utils/datetime-util';
    import { createEmptyTable, toTimeUnits } from '$lib/utils/time-table-util';
    import { newClient } from '$lib/utils/client';
    import type { components } from '$lib/sdk/schema';
    import { getErrorMessage } from '$lib/utils/error-handler';
    import { fade } from 'svelte/transition';
    import AvailabilitySection from '$lib/components/spot-component/availability-section.svelte';
    import SpotInfo from '$lib/components/spot-component/spot-info.svelte';
    import { ChevronSortDown, Favorite, FavoriteFilled } from 'carbon-icons-svelte';
    import IntersectionObserver from 'svelte-intersection-observer';

    export let data: PageData;
    type TimeUnit = components['schemas']['TimeUnit'];
    type Car = components['schemas']['Car'];

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
    let today: Date = new Date(Date.now());
    let currentMonday: Date = getMonday(today);
    let nextMonday: Date;

    let availabilityTablesInitial = new Map<number, TimeSlotStatus[][]>();
    availabilityTablesInitial.set(
        currentMonday.getTime(),
        createIntialAvailabilityTable(availabilityTimeSlot ?? [], today, currentMonday)
    );
    let availabilityTableBooking = structuredClone(availabilityTablesInitial);
    let availabilityTable: TimeSlotStatus[][];

    let canLoadMore = data.hasNext;
    let loadTrigger: HTMLElement | null = null;
    let intersecting: boolean;
    let loadLock = false;
    let selectedCar: Car | null = null;
    let showCarDropDown: boolean = false;

    $: while (intersecting && canLoadMore && !loadLock) {
        loadLock = true;
        data.paging
            .next()
            .then(({ value: { data: cars }, done }) => {
                if (cars) {
                    data.cars = [...data.cars, ...cars];
                }
                canLoadMore = !done;
            })
            .finally(() => {
                loadLock = false;
            });
    }

    let bookSlotsCount: number = 0;
    let total: number = 0;

    $: total = data.spot?.price_per_hour ? (bookSlotsCount * data.spot?.price_per_hour) / 2 : 0;
    $: nextMonday = getDateWithDayOffset(currentMonday, DAY_IN_A_WEEK);
    $: availabilityTable = availabilityTableBooking.get(currentMonday.getTime()) || [];

    function refreshAvailabilityTable() {
        currentMonday = getMonday(today);
        const nextMonday = getDateWithDayOffset(currentMonday, DAY_IN_A_WEEK);
        nextMonday.setMinutes(nextMonday.getMinutes() - 30);
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
            .then(({ data: availability, error: err }) => {
                if (err) {
                    errorMessage = getErrorMessage(err);
                    toastTimeOut = ERROR_MESSAGE_TIME_OUT;
                    return;
                }
                availabilityTablesInitial = new Map<number, TimeSlotStatus[][]>();
                availabilityTablesInitial.set(
                    currentMonday.getTime(),
                    createIntialAvailabilityTable(availability ?? [], today, currentMonday)
                );
                availabilityTableBooking = structuredClone(availabilityTablesInitial);
                availabilityTable = availabilityTableBooking.get(currentMonday.getTime()) ?? [];
                bookSlotsCount = 0;
                selectedCar = null;
            });
    }

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
                        availabilityTableBooking.set(
                            currentMonday.getTime(),
                            structuredClone(newTable)
                        );
                        availabilityTable =
                            availabilityTableBooking.get(currentMonday.getTime()) ?? [];
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

    function handleSubmitBooking(event: Event) {
        event.preventDefault();
        console.log('submit');
        let bookSlots = toTimeUnits(availabilityTableBooking, [TimeSlotStatus.BOOK_INTENT]);

        if (bookSlots.length == 0) {
            errorMessage = BOOK_WITHOUT_SLOTS_ERROR;
            toastTimeOut = ERROR_MESSAGE_TIME_OUT;
            return;
        }

        if (selectedCar == null) {
            errorMessage = BOOK_WITHOUT_CAR;
            toastTimeOut = ERROR_MESSAGE_TIME_OUT;
            return;
        }

        if (confirm(`Are you sure you want to book this spot? The total price will be ${total}`)) {
            client
                .POST('/spots/{id}/bookings', {
                    body: { booked_times: bookSlots, car_id: selectedCar.id },
                    params: { path: { id: spot?.id ?? '0' } }
                })
                .then(({ error: err }) => {
                    if (err) {
                        errorMessage = getErrorMessage(err);
                        toastTimeOut = ERROR_MESSAGE_TIME_OUT;
                    } else {
                        refreshAvailabilityTable();
                    }
                });
        }
    }

    function resetBookingEdit() {
        availabilityTableBooking = structuredClone(availabilityTablesInitial);
        bookSlotsCount = 0;
    }

    /*
    this function  remove the edit event if there is already one event at that time slot. Or 
    append the new event if no event happens to that time slot
    */
    function handleEdit(event: CustomEvent) {
        if ([TimeSlotStatus.BOOK_INTENT, TimeSlotStatus.AVAILABLE].includes(event.detail.status)) {
            availabilityTable[event.detail.segment] ??= [];
            if (event.detail.status === TimeSlotStatus.AVAILABLE) {
                bookSlotsCount += 1;
                availabilityTable[event.detail.segment][event.detail.day] =
                    TimeSlotStatus.BOOK_INTENT;
            } else {
                bookSlotsCount -= 1;
                availabilityTable[event.detail.segment][event.detail.day] =
                    TimeSlotStatus.AVAILABLE;
            }
            availabilityTableBooking.set(currentMonday.getTime(), availabilityTable);
            availabilityTableBooking = availabilityTableBooking;
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

    const handleClickOutside = (event: Event) => {
        if (!(event.target as HTMLElement).closest('.dropdown')) {
            showCarDropDown = false;
        }
    };
</script>

<Content>
    <!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
    <div on:click={handleClickOutside}>
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
                <Favorite
                    style="fill:black; position: absolute; "
                    size={SPOT_PREFFERED_ICON_SIZE}
                />
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
        <Form on:submit={handleSubmitBooking}>
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
                />
            </div>
            <div class="price-field">
                <NumberInput
                    label="Total price"
                    hideSteppers
                    step={0.01}
                    min={0}
                    name="total-per-hour"
                    helperText="Price in CAD"
                    required
                    bind:value={total}
                />
            </div>
            <FormGroup>
                <label for="cars">Choose a car:</label>
                <div class="dropdown">
                    {#if showCarDropDown}
                        <div class="dropdown-content">
                            {#key data.cars}
                                {#each data.cars as car}
                                    <!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
                                    <div
                                        on:click={() => {
                                            selectedCar = car;
                                            showCarDropDown = false;
                                        }}
                                    >
                                        {car.details.license_plate}
                                    </div>
                                {/each}
                            {/key}
                            <IntersectionObserver element={loadTrigger} bind:intersecting>
                                {#if canLoadMore}
                                    <div bind:this={loadTrigger}>Loading...</div>
                                {/if}
                            </IntersectionObserver>
                        </div>
                    {/if}

                    <!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
                    <div
                        class="dropdown-field"
                        on:click={() => {
                            showCarDropDown = true;
                        }}
                    >
                        <div>
                            {selectedCar ? selectedCar.details.license_plate : 'Select Your Car'}
                        </div>
                        <TooltipIcon
                            style="margin-left: auto; pointer-events: none;"
                            icon={ChevronSortDown}
                            tooltipText={'select your car'}
                        ></TooltipIcon>
                    </div>
                </div>
            </FormGroup>
            <Button kind="secondary" on:click={resetBookingEdit}>Reset</Button>
            <Button type="submit">Submit</Button>
        </Form>
    </div>
</Content>

<style>
    .header-with-preference-option {
        display: flex;
        gap: 1rem;
        align-self: center;
    }

    .dropdown {
        position: relative;
        display: inline-block;
    }

    /* Style for the button that triggers the dropdown */
    .dropdown-field {
        display: flex;
        flex-direction: row;
        padding: 10px;
        background-color: #f4f4f4;
        color: black;
        cursor: pointer;
        min-width: 5rem;
        font-size: 1rem;
        border-bottom: 1px solid black;
    }

    /* Style for the dropdown content (hidden by default) */
    .dropdown-content {
        display: block;
        position: absolute;
        background-color: #f9f9f9;
        min-width: 10rem;
        box-shadow: 0px 8px 16px rgba(0, 0, 0, 0.2);
        z-index: 1;
        max-height: 10rem; /* Limit the height */
        overflow-y: auto; /* Enable scrolling */
        border-radius: 5px;
    }

    /* Style for each item in the dropdown list */
    .dropdown-content div {
        color: black;
        padding: 12px 16px;
        text-decoration: none;
        display: block;
    }

    .dropdown-content div:hover {
        background-color: #ddd;
    }
</style>
