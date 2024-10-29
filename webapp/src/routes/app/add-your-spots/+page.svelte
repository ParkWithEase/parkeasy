<script lang="ts">
    import { goto } from '$app/navigation';
    import AvailabilityTable from '$lib/components/spot-component/availability-table.svelte';
    import {
        CREATE_WITH_EMPTY_AVAILABILITY_TABLE_ERROR,
        DAY_IN_A_WEEK,
        ERROR_MESSAGE_TIME_OUT,
        WAIT_TIME_BEFORE_AUTO_COMPLETE
    } from '$lib/constants';
    import { TimeSlotStatus } from '$lib/enum/timeslot-status';
    import { newClient } from '$lib/utils/client';
    import { getDateWithDayOffset, getMonday } from '$lib/utils/datetime-util';
    import { getErrorMessage } from '$lib/utils/error-handler';
    import { initializeAvailabilityTables, toTimeUnits } from '$lib/utils/time-table-util';

    import {
        Checkbox,
        Form,
        ProgressIndicator,
        ProgressStep,
        TextInput,
        Button,
        ToastNotification,
        Select,
        SelectItem
    } from 'carbon-components-svelte';
    import { ArrowLeft, ArrowRight } from 'carbon-icons-svelte';
    import { fade } from 'svelte/transition';

    let client = newClient();

    let today = new Date(Date.now());
    let currentMonday = getMonday(today);

    let availabilityTablesInitial = new Map<number, TimeSlotStatus[][]>();
    initializeAvailabilityTables(availabilityTablesInitial, currentMonday, today);
    let availabilityTables = structuredClone(availabilityTablesInitial);
    let availabilityTable: TimeSlotStatus[][];

    $: availabilityTable = availabilityTables.get(currentMonday.getTime()) || [];
    let next_monday: Date;
    $: next_monday = getDateWithDayOffset(currentMonday, DAY_IN_A_WEEK);

    //This contain spot info history to be submitted to the server
    let new_price_per_hour: number;
    let city: string;
    let country_code: string;
    let postal_code: string;
    let state: string;
    let street_address: string;

    let has_shelter: boolean;
    let has_plug_in: boolean;
    let has_charing_station: boolean;

    //control form flow
    let currentIndex: number = 0;
    let isLocationStepCompleted: boolean = false;
    let isAvailabilityStepCompleted: boolean = false;

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
        new_price_per_hour = 0;
    }

    function handleSubmitAvailability() {
        if (toTimeUnits(availabilityTables).length > 0) {
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
                        charging_station: has_charing_station,
                        plug_in: has_plug_in,
                        shelter: has_shelter
                    },
                    location: {
                        city: city,
                        country_code: country_code,
                        postal_code: postal_code,
                        state: state,
                        street_address: street_address
                    },
                    price_per_hour: new_price_per_hour,
                    availability: toTimeUnits(availabilityTables)
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
            <Form on:submit={handleSubmitLocation}>
                <TextInput
                    required
                    labelText="Street address"
                    name="street-address"
                    placeholder="Street Address"
                    readonly={currentIndex == 2}
                    bind:value={street_address}
                />
                <TextInput
                    required
                    labelText="City"
                    name="city"
                    placeholder="City"
                    readonly={currentIndex == 2}
                    bind:value={city}
                />
                <Select
                    style={currentIndex == 2 ? 'pointer-events: none;' : ' '}
                    labelText="Province"
                    bind:selected={state}
                >
                    <SelectItem value="MB" text="Manitoba" />
                    <SelectItem value="ON" text="Ontario" />
                    <SelectItem value="AB" text="Alberta" />
                    <SelectItem value="QC" text="Quebec" />
                    <SelectItem value="NS" text="Nova Scotia" />
                    <SelectItem value="BC" text="British Columbia" />
                    <SelectItem value="NL" text="Newfouundland and Labrador" />
                    <SelectItem value="PE" text="Prince Edward Island" />
                    <SelectItem value="SK" text="Saskatchewan" />
                    <SelectItem value="YT" text="Yukon" />
                    <SelectItem value="NU" text="Nunavut" />
                    <SelectItem value="NT" text="Northwest Territories" />
                </Select>

                <Select
                    style={currentIndex == 2 ? 'pointer-events: none;' : ' '}
                    labelText="Country Code"
                    bind:selected={country_code}
                >
                    <SelectItem value="CA" text="Canada" />
                </Select>
                <TextInput
                    required
                    labelText="Postal code"
                    name="postal-code"
                    placeholder="Postal code"
                    readonly={currentIndex == 2}
                    bind:value={postal_code}
                />
                <p>Utilities</p>
                <Checkbox
                    name="shelter"
                    style={currentIndex == 2 ? 'pointer-events: none;' : ' '}
                    labelText="shelter"
                    bind:checked={has_shelter}
                />
                <Checkbox
                    name="plug-in"
                    labelText="Plug-in"
                    style={currentIndex == 2 ? 'pointer-events: none;' : ' '}
                    bind:checked={has_plug_in}
                />
                <Checkbox
                    name="charging-station"
                    labelText="Charging Station"
                    style={currentIndex == 2 ? 'pointer-events: none;' : ' '}
                    bind:checked={has_charing_station}
                />
                {#if currentIndex == 0}
                    <Button type="submit">Submit</Button>
                {/if}
            </Form>
        {/if}
        {#if currentIndex == 1 || currentIndex == 2}
            <p class="spot-form-header">Availability</p>
            <div class="date-controller">
                <p>
                    From {currentMonday?.toString()} to {next_monday?.toString()}
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
            <div>
                <AvailabilityTable
                    bind:availability_table={availabilityTable}
                    on:edit={handleEdit}
                />
            </div>
            <Form on:submit={handleSubmitAvailability}>
                <div class="price-field">
                    <TextInput
                        labelText="Price per hour"
                        name="price-per-hour"
                        helperText="Price in CAD"
                        type="number"
                        required
                        readonly={currentIndex == 2}
                        min={0}
                        bind:value={new_price_per_hour}
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
