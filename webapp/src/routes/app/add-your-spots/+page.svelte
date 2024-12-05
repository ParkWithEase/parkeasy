<script lang="ts">
    import { goto } from '$app/navigation';
    import AvailabilitySection from '$lib/components/spot-component/availability-section.svelte';
    import SpotLocationCreateForm from '$lib/components/spot-component/spot-location-create-form.svelte';
    import {
        CREATE_WITH_EMPTY_AVAILABILITY_TABLE_ERROR,
        DAY_IN_A_WEEK,
        ERROR_MESSAGE_TIME_OUT
    } from '$lib/constants';
    import { TimeSlotStatus } from '$lib/enum/timeslot-status';
    import { newClient } from '$lib/utils/client';
    import { getDateWithDayOffset, getMonday } from '$lib/utils/datetime-util';
    import { getErrorMessage } from '$lib/utils/error-handler';
    import { initializeAvailabilityTables, toTimeUnits } from '$lib/utils/time-table-util';

    import {
        Form,
        ProgressIndicator,
        ProgressStep,
        TextInput,
        Button,
        ToastNotification,
        NumberInput
    } from 'carbon-components-svelte';
    import { fade } from 'svelte/transition';

    let client = newClient();

    let today = new Date(Date.now());
    let currentMonday = getMonday(today);

    let availabilityTablesInitial = new Map<number, TimeSlotStatus[][]>();
    initializeAvailabilityTables(availabilityTablesInitial, currentMonday, today);
    let availabilityTables = structuredClone(availabilityTablesInitial);
    let availabilityTable: TimeSlotStatus[][];

    $: availabilityTable = availabilityTables.get(currentMonday.getTime()) || [];
    let nextMonday: Date;
    $: nextMonday = getDateWithDayOffset(currentMonday, DAY_IN_A_WEEK);

    //This contain spot info history to be submitted to the server
    let newPricePerHour: number;
    let city: string;
    let countryCode: string;
    let postalCode: string;
    let state: string;
    let streetAddress: string;

    let hasShelter: boolean;
    let hasPlugIn: boolean;
    let hasChargingStation: boolean;

    //control form flow
    let currentIndex: number = 0;
    let isLocationStepCompleted: boolean = false;
    let isAvailabilityStepCompleted: boolean = false;

    $: isLocationReadOnly = currentIndex == 2;
    //Error message toast
    let toastTimeOut: number = 0;
    let errorMessage: string = '';
    $: showToast = toastTimeOut !== 0;

    function toNextWeek() {
        currentMonday = getDateWithDayOffset(currentMonday, DAY_IN_A_WEEK);
    }

    function toPrevWeek() {
        if (availabilityTable?.[0]?.[0] === TimeSlotStatus.PASTDUE) {
            return;
        }
        currentMonday = getDateWithDayOffset(currentMonday, -DAY_IN_A_WEEK);
    }

    function handleSubmitLocation(event: Event) {
        event.preventDefault();
        currentIndex += 1;
        isLocationStepCompleted = true;
    }

    function clearEditRecords() {
        availabilityTables = structuredClone(availabilityTablesInitial);
        newPricePerHour = 0;
    }

    function handleSubmitAvailability(event: Event) {
        event.preventDefault();
        if (toTimeUnits(availabilityTables, [TimeSlotStatus.AVAILABLE]).length > 0) {
            currentIndex += 1;
            isAvailabilityStepCompleted = true;
        } else {
            isAvailabilityStepCompleted = false;
            errorMessage = CREATE_WITH_EMPTY_AVAILABILITY_TABLE_ERROR;
            toastTimeOut = ERROR_MESSAGE_TIME_OUT;
        }
    }

    function handleSubmitAll() {
        client
            .POST('/spots', {
                body: {
                    features: {
                        charging_station: hasChargingStation,
                        plug_in: hasPlugIn,
                        shelter: hasShelter
                    },
                    location: {
                        city: city,
                        country_code: countryCode,
                        postal_code: postalCode,
                        state: state,
                        street_address: streetAddress
                    },
                    price_per_hour: newPricePerHour,
                    availability: toTimeUnits(availabilityTables, [TimeSlotStatus.AVAILABLE])
                }
            })
            .then(({ data, error: err }) => {
                if (err) {
                    errorMessage = getErrorMessage(err);
                    toastTimeOut = ERROR_MESSAGE_TIME_OUT;
                } else {
                    goto(`/app/your-spots/${data.id}/spot-info`, { replaceState: true });
                }
            });
    }

    function handleEdit(event: CustomEvent) {
        if (event.detail.status == TimeSlotStatus.PASTDUE || currentIndex == 2) {
            return;
        } else {
            availabilityTable[event.detail.segment] ??= [];
            availabilityTable[event.detail.segment][event.detail.day] =
                event.detail.status === TimeSlotStatus.NONE
                    ? TimeSlotStatus.AVAILABLE
                    : TimeSlotStatus.NONE;
            availabilityTables.set(currentMonday.getTime(), availabilityTable);
            availabilityTables = availabilityTables;
        }
    }
</script>

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

<div class="add-page-layout">
    <div style="position:sticky; top: 4rem;">
        <ProgressIndicator vertical bind:currentIndex>
            <ProgressStep
                label="Parking spot location"
                bind:complete={isLocationStepCompleted}
                description="Info about your parking spot location"
            ></ProgressStep>
            <ProgressStep
                label="Availability Table"
                bind:complete={isAvailabilityStepCompleted}
                description="Info about your parking spot availability"
            ></ProgressStep>
            <ProgressStep label="Review" description="Review before submission"></ProgressStep>
        </ProgressIndicator>
    </div>

    <div class="spot-form-container">
        {#if currentIndex == 0 || currentIndex == 2}
            <p class="spot-form-header">Location and Utilities</p>
            <SpotLocationCreateForm
                bind:isReadOnly={isLocationReadOnly}
                bind:streetAddress
                bind:city
                bind:state
                bind:countryCode
                bind:postalCode
                bind:hasShelter
                bind:hasPlugIn
                bind:hasChargingStation
                handleSubmit={handleSubmitLocation}
            />
        {/if}
        {#if currentIndex == 1 || currentIndex == 2}
            <p class="spot-form-header">Availability</p>
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
                {#if currentIndex == 1}
                    <Button kind="secondary" on:click={clearEditRecords}>Clear</Button>
                    <Button type="submit">Submit</Button>
                {/if}
            </Form>
        {/if}

        {#if currentIndex == 2}
            <Button on:click={handleSubmitAll}>Confirm</Button>
        {/if}
    </div>
</div>

<style>
    .spot-form-container {
        flex-grow: 8;
    }

    .add-page-layout {
        display: flex;
        flex-direction: row;
        align-items: start;
        gap: 2rem;
    }

    .spot-form-header {
        font-size: 2rem;
    }

    .error-message {
        position: sticky;
        top: 4rem;
        width: 100%;
        background-color: black;
        z-index: 2;
    }
</style>
