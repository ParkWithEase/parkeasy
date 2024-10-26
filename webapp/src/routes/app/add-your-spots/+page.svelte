<script lang="ts">
    import AvailabilityTable from '$lib/components/spot-component/availability-table.svelte';
    import {
        CREATE_WITH_EMPTY_AVAILABILITY_TABLE_ERROR,
        DAY_IN_A_WEEK,
        ERROR_MESSAGE_TIME_OUT,
        WAIT_TIME_BEFORE_AUTO_COMPLETE
    } from '$lib/constants';
    import { TimeSlotStatus } from '$lib/enum/timeslot-status';
    import { convertDateToUTC, getDateWithDayOffset, getMonday } from '$lib/utils/datetime-util';
    import { getWeekAvailabilityTable } from '$lib/utils/time-table-util';

    import {
        Checkbox,
        Form,
        ProgressIndicator,
        ProgressStep,
        TextInput,
        Button,
        ToastNotification
    } from 'carbon-components-svelte';
    import { ArrowLeft, ArrowRight } from 'carbon-icons-svelte';
    import { fade } from 'svelte/transition';
    let availability_table: TimeSlotStatus[][];

    //reminder to change these to now()
    let today: Date;
    let currentMonday: Date;
    let nextMonday: Date;

    //This array contain all edit history
    let time_slot_edit_records = [];
    //This contain spot info history to be submitted to the server
    let spotInfo;

    //control form flow
    let currentIndex: number = 0;
    let isLocationStepCompleted: boolean = false;
    let isAvailablityStepCompleted: boolean = false;

    //Error message toast

    let toastTimeOut: number = 0;
    let errorMessage: string = '';
    $: showToast = toastTimeOut !== 0;

    today = convertDateToUTC(new Date(Date.now()));
    currentMonday = getMonday(today);
    nextMonday = getDateWithDayOffset(currentMonday, DAY_IN_A_WEEK);
    availability_table = getWeekAvailabilityTable(today, currentMonday, [], time_slot_edit_records);

    function toNextWeek() {
        currentMonday = getDateWithDayOffset(currentMonday, DAY_IN_A_WEEK);
        nextMonday = getDateWithDayOffset(nextMonday, DAY_IN_A_WEEK);
        console.log(time_slot_edit_records);
        availability_table = getWeekAvailabilityTable(
            today,
            currentMonday,
            [],
            time_slot_edit_records
        );
    }

    function toPrevWeek() {
        currentMonday = getDateWithDayOffset(currentMonday, -DAY_IN_A_WEEK);
        nextMonday = getDateWithDayOffset(nextMonday, -DAY_IN_A_WEEK);
        availability_table = getWeekAvailabilityTable(
            today,
            currentMonday,
            [],
            time_slot_edit_records
        );
    }

    function handleSubmitLocation(event: Event) {
        event.preventDefault();
        const formData = new FormData(event.target as HTMLFormElement);
        spotInfo = {
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
            }
        };
        currentIndex += 1;
        isLocationStepCompleted = true;
    }

    function clearEditRecords() {
        time_slot_edit_records = [];
        availability_table = getWeekAvailabilityTable(
            today,
            currentMonday,
            [],
            time_slot_edit_records
        );
    }

    function handleSubmitAvailability() {
        if (time_slot_edit_records.length !== 0) {
            currentIndex += 1;
            isAvailablityStepCompleted = true;
        } else {
            isAvailablityStepCompleted = false;
            errorMessage = CREATE_WITH_EMPTY_AVAILABILITY_TABLE_ERROR;
            toastTimeOut = ERROR_MESSAGE_TIME_OUT;
        }
    }

    function handleSubmitAll() {
        console.log('Submit everything and go to the other page after the request is done');
    }

    let autoCompleteTimer: NodeJS.Timeout;
    function getAutoCompleteAddress() {
        clearTimeout(autoCompleteTimer);
        autoCompleteTimer = setTimeout(() => console.log('fetch'), WAIT_TIME_BEFORE_AUTO_COMPLETE);
    }

    function handleEdit(event: Event) {
        if (event.detail.status == TimeSlotStatus.PASTDUE || currentIndex == 2) {
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
            availability_table[event.detail.segment][event.detail.day] = new_time_slot.status;
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
                bind:complete={isAvailablityStepCompleted}
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
                    on:keyup={getAutoCompleteAddress}
                    required
                    labelText="Street address"
                    name="street-address"
                    placeholder="Street Address"
                    readonly={currentIndex == 2}
                    value={spotInfo?.location.street_address}
                />
                <TextInput
                    required
                    labelText="City"
                    name="city"
                    placeholder="City"
                    readonly={currentIndex == 2}
                    value={spotInfo?.location.city}
                />
                <TextInput
                    required
                    labelText="State/Province"
                    name="state"
                    placeholder="State/Province"
                    readonly={currentIndex == 2}
                    value={spotInfo?.location.state}
                />
                <TextInput
                    required
                    labelText="Country"
                    name="country-code"
                    placeholder="Country Code"
                    readonly={currentIndex == 2}
                    value={spotInfo?.location.country_code}
                />
                <TextInput
                    required
                    labelText="Postal code"
                    name="postal-code"
                    placeholder="Postal code"
                    readonly={currentIndex == 2}
                    value={spotInfo?.location.postal_code}
                />
                <p>Utilities</p>
                <Checkbox
                    name="shelter"
                    style={currentIndex == 2 ? 'pointer-events: none;' : ' '}
                    labelText="shelter"
                    checked={spotInfo?.features.shelter}
                />
                <Checkbox
                    name="plug-in"
                    labelText="Plug-in"
                    style={currentIndex == 2 ? 'pointer-events: none;' : ' '}
                    checked={spotInfo?.features.plug_in}
                />
                <Checkbox
                    name="charging-station"
                    labelText="Charging Station"
                    style={currentIndex == 2 ? 'pointer-events: none;' : ' '}
                    checked={spotInfo?.features.charging_station}
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
            <div style="margin: 1rem">
                <AvailabilityTable bind:availability_table on:edit={handleEdit} />
            </div>

            {#if currentIndex == 1}
                <Button kind="secondary" on:click={clearEditRecords}>Clear</Button>
                <Button on:click={handleSubmitAvailability}>Submit</Button>
            {/if}
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
