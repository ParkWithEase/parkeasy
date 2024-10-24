<script lang="ts">
    import { DAY_IN_A_WEEK, TOTAL_SEGMENTS_NUMBER } from '$lib/constants';
    import { TimeSlotStatus } from '$lib/enum/timeslot-status';

    export let availability_table: TimeSlotStatus[][] = Array.from(
        { length: TOTAL_SEGMENTS_NUMBER },
        () => Array(DAY_IN_A_WEEK).fill(0)
    );
</script>

<table class="availability-table">
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
    {#each Array(24) as _, segment (segment)}
        <tr class="odd-row">
            <th rowspan="2"> {segment}:00 </th>
            {#each Array(7) as _, day (day)}
                {@const status = availability_table[segment * 2][day]}
                {#if status == TimeSlotStatus.AVAILABLE}
                    <td class="available"></td>
                {:else if status == TimeSlotStatus.BOOKED}
                    <td class="booked"></td>
                {:else if status == TimeSlotStatus.AUCTIONED}
                    <td class="auctioned"></td>
                {:else if status == TimeSlotStatus.PASTDUE}
                    <td class="pastdue"></td>
                {:else}
                    <td></td>
                {/if}
            {/each}
        </tr>
        <tr class="even-row">
            {#each Array(7) as _, day (day)}
                {@const status = availability_table[segment * 2 + 1][day]}
                {#if status == TimeSlotStatus.AVAILABLE}
                    <td class="available"></td>
                {:else if status == TimeSlotStatus.BOOKED}
                    <td class="booked"></td>
                {:else if status == TimeSlotStatus.AUCTIONED}
                    <td class="auctioned"></td>
                {:else if status == TimeSlotStatus.PASTDUE}
                    <td class="pastdue"></td>
                {:else}
                    <td></td>
                {/if}
            {/each}
        </tr>
    {/each}
</table>

<style>
    td,
    th {
        border-left: 1px solid black;
        border-right: 1px solid black;
    }
    td {
        font-size: 1rem;
        padding: 0.3rem;
        vertical-align: middle;
        text-align: center;
    }

    th {
        font-size: 1.2rem;
        padding: 0.3rem;
        vertical-align: middle;
        text-align: center;
    }
    .odd-row {
        border: 1px solid black;
        border-bottom: none;
    }

    .even-row {
        border: 1px solid black;
        border-top: none;
    }

    .header-row {
        border: 1px solid black;
    }

    table {
        width: 100%;
        table-layout: fixed;
    }

    td.available {
        background-color: rgba(111, 220, 140, 0.5);
    }

    td.booked {
        background-color: rgba(218, 30, 40, 0.5);
    }

    td.auctioned {
        background-color: rgba(253, 220, 105, 0.5);
    }

    td.pastdue {
        background-color: rgba(185, 185, 185, 0.5);
    }
</style>
