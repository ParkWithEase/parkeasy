<script lang="ts">
    import { DAY_IN_A_WEEK, TOTAL_SEGMENTS_NUMBER } from '$lib/constants';
    import { TimeSlotStatus } from '$lib/enum/timeslot-status';
    import { createEventDispatcher } from 'svelte';

    let dispatcher = createEventDispatcher();
    export let availability_table: TimeSlotStatus[][] = Array.from(
        { length: TOTAL_SEGMENTS_NUMBER },
        () => Array(DAY_IN_A_WEEK).fill(0)
    );

    function handleEdit(event: Event, day: number, segment: number, status: TimeSlotStatus) {
        event.preventDefault();
        if (event.buttons == 1 || (event.buttons == 0 && event.button == 0)) {
            dispatcher('edit', {
                day: day,
                segment: segment,
                status: status
            });
        }
    }

    let segmentIndexArray: Array<number> = [...Array(Math.floor(TOTAL_SEGMENTS_NUMBER / 2)).keys()];
    let dayIndexArray: Array<number> = [...Array(DAY_IN_A_WEEK).keys()];
</script>

<table class="availability-table" draggable="false">
    <tr class="header-row">
        <th>Time</th>
        <th>Mon</th>
        <th>Tues</th>
        <th>Wed</th>
        <th>Thu</th>
        <th>Fri</th>
        <th>Sat</th>
        <th>Sun</th>
    </tr>

    {#each segmentIndexArray as segment (segment)}
        <tr class="odd-row" draggable="false">
            <th rowspan="2" draggable="false"> {segment}:00 </th>
            {#each dayIndexArray as day (day)}
                {@const status = availability_table[segment * 2][day]}
                <td
                    draggable="false"
                    on:pointerenter={(e) => handleEdit(e, day, segment * 2, status)}
                    on:mousedown={(e) => handleEdit(e, day, segment * 2, status)}
                >
                    {#if status == TimeSlotStatus.AVAILABLE}
                        <div class="available"></div>
                    {:else if status == TimeSlotStatus.BOOKED}
                        <div class="booked"></div>
                    {:else if status == TimeSlotStatus.AUCTIONED}
                        <div class="auctioned"></div>
                    {:else if status == TimeSlotStatus.PASTDUE}
                        <div class="pastdue"></div>
                    {:else}
                        <div></div>
                    {/if}
                </td>
            {/each}
        </tr>
        <tr draggable="false" class="even-row">
            {#each dayIndexArray as day (day)}
                {@const status = availability_table[segment * 2 + 1][day]}
                <td
                    draggable="false"
                    on:pointerenter={(e) => handleEdit(e, day, segment * 2 + 1, status)}
                    on:mousedown={(e) => handleEdit(e, day, segment * 2 + 1, status)}
                >
                    {#if status == TimeSlotStatus.AVAILABLE}
                        <div class="available"></div>
                    {:else if status == TimeSlotStatus.BOOKED}
                        <div class="booked"></div>
                    {:else if status == TimeSlotStatus.AUCTIONED}
                        <div class="auctioned"></div>
                    {:else if status == TimeSlotStatus.PASTDUE}
                        <div class="pastdue"></div>
                    {:else}
                        <div></div>
                    {/if}
                </td>
            {/each}
        </tr>
    {/each}
</table>

<style>
    td,
    th {
        border-left: 2px solid black;
        border-right: 2px solid black;
    }
    td {
        height: 100%;
        width: 100%;
        font-size: 1rem;
        vertical-align: middle;
        text-align: center;
    }

    td:hover {
        border: 2px solid rgb(0, 102, 255);
    }
    tr {
        height: 100%;
    }

    th {
        font-size: 1.2rem;
        padding: 0.3rem;
        vertical-align: middle;
        text-align: center;
    }
    .odd-row {
        border: 2px solid black;
        border-bottom: none;
    }

    .even-row {
        border: 2px solid black;
        border-top: none;
    }

    .header-row {
        border: 2px solid black;
    }

    table {
        width: 100%;
        height: 1px;
        table-layout: fixed;
    }

    td > div {
        height: 100%;
        pointer-events: none;
    }
    div.available {
        background-color: rgba(111, 220, 140, 0.5);
    }

    div.booked {
        background-color: rgba(218, 30, 40, 0.5);
    }

    div.auctioned {
        background-color: rgba(253, 220, 105, 0.5);
    }

    div.pastdue {
        background-color: rgba(185, 185, 185, 0.5);
    }
</style>
