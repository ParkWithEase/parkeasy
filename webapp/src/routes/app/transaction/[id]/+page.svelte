<script lang="ts">
    import type { PageData } from './$types';
    import Background from '$lib/images/background.png';
    import { Content, NumberInput } from 'carbon-components-svelte';
    import { TimeSlotStatus } from '$lib/enum/timeslot-status';
    import { DAY_IN_A_WEEK } from '$lib/constants';
    import { getMonday, getDateWithDayOffset } from '$lib/utils/datetime-util';
    import { createEmptyTable } from '$lib/utils/time-table-util';
    import type { components } from '$lib/sdk/schema';
    import AvailabilitySection from '$lib/components/spot-component/availability-section.svelte';
    import SpotInfo from '$lib/components/spot-component/spot-info.svelte';

    export let data: PageData;
    type TimeUnit = components['schemas']['TimeUnit'];

    //Variable for location edit section
    let spot = data.spot;
    let car = data.car;
    //reminder to change these to now()
    let today: Date = new Date(Date.now());
    let currentMonday: Date;

    let nextMonday: Date;

    let bookedTableList = new Map<number, TimeSlotStatus[][]>();
    bookedTableList = createInitialBookedTable(data.transaction_info?.booked_times ?? [], today);
    currentMonday = new Date(
        bookedTableList.keys().reduce((min, current) => {
            return current < min ? current : min;
        })
    );

    let bookedTable: TimeSlotStatus[][];
    bookedTable = bookedTableList.get(currentMonday.getTime()) || [];
    let total: number = data.transaction_info?.paid_amount ?? 0;

    $: nextMonday = getDateWithDayOffset(currentMonday, DAY_IN_A_WEEK);
    $: bookedTable = bookedTableList.get(currentMonday.getTime()) || [];

    function createInitialBookedTable(timeUnits: TimeUnit[], today: Date) {
        let bookingTable = new Map<number, TimeSlotStatus[][]>();

        timeUnits.forEach((timeUnit) => {
            const date = new Date(timeUnit.start_time);
            const day = (date.getDay() || 7) - 1;
            const segment = date.getHours() * 2 + Math.ceil(date.getMinutes() / 30);
            const time_unit_week_monday = getMonday(date);
            if (!bookingTable.get(time_unit_week_monday.getTime())) {
                bookingTable.set(time_unit_week_monday.getTime(), createEmptyTable());
            }

            let currentTable = bookingTable.get(time_unit_week_monday.getTime()) ?? [];
            if (date.getTime() < today.getTime()) {
                currentTable[segment][day] = TimeSlotStatus.EXPIRED_BOOK;
            } else {
                currentTable[segment][day] = TimeSlotStatus.ACTIVE_BOOK;
            }
        });
        return bookingTable;
    }

    function toNextWeek() {
        const nextMonday = getDateWithDayOffset(currentMonday, DAY_IN_A_WEEK);
        if (!bookedTableList.get(nextMonday.getTime())) {
            return;
        } else {
            currentMonday = getDateWithDayOffset(currentMonday, DAY_IN_A_WEEK);
            bookedTable = bookedTableList.get(nextMonday.getTime()) ?? [];
        }
    }

    function toPrevWeek() {
        const lastMonday = getDateWithDayOffset(currentMonday, -DAY_IN_A_WEEK);
        if (!bookedTableList.get(lastMonday.getTime())) {
            return;
        } else {
            currentMonday = getDateWithDayOffset(currentMonday, -DAY_IN_A_WEEK);
            bookedTable = bookedTableList.get(lastMonday.getTime()) ?? [];
        }
    }
</script>

<Content>
    <div>
        <img src={Background} class="spot-info-image" alt="spot" />
        <p class="spot-info-header">Location</p>

        <SpotInfo bind:spot />

        <p class="spot-info-header">Parking Car Info</p>
        <p>License Plate: {car?.license_plate}</p>
        <p>Color: {car?.color}</p>
        <p>Model: {car?.model}</p>
        <p>Make: {car?.make}</p>

        <p class="spot-info-header">Booked Time slots</p>
        <div class="price-field">
            <NumberInput
                style="pointer-events:none"
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
        <AvailabilitySection
            bind:currentMonday
            bind:nextMonday
            bind:availabilityTable={bookedTable}
            {toPrevWeek}
            {toNextWeek}
        />
    </div></Content
>
