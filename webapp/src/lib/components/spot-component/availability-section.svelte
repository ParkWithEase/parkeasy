<script lang="ts">
    import { Button } from 'carbon-components-svelte';
    import AvailabilityTable from './availability-table.svelte';
    import { ArrowLeft, ArrowRight } from 'carbon-icons-svelte';
    import type { TimeSlotStatus } from '$lib/enum/timeslot-status';

    export let currentMonday: Date;
    export let nextMonday: Date;
    export let availabilityTable: TimeSlotStatus[][];

    export let toPrevWeek = () => {};
    export let toNextWeek = () => {};
    export let handleEdit = (_: CustomEvent) => {};
</script>

<div class="date-controller">
    <p>
        From {currentMonday?.toString()} to {nextMonday?.toString()}
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
    <AvailabilityTable bind:availability_table={availabilityTable} on:edit={handleEdit} />
</div>
