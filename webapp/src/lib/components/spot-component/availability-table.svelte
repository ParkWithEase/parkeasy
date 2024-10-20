<script lang="ts">
    import { TimeSlotStatus } from '$lib/enum/timeslot-status';

    export const availability_table: TimeSlotStatus[][] = Array.from({ length: 48 }, () =>
        Array(7).fill(0)
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
    {#each Array(24) as _, index (index)}
        <tr class="odd-row">
            <th rowspan="2"> {index}:00 </th>
            <td>random data</td>
            <td>random data</td>
            <td>random data</td>
            <td>random data</td>
            <td>random data</td>
            <td>random data</td>
            <td>random data</td>
        </tr>
        <tr class="even-row">
            <td>random data</td>
            <td>random data</td>
            <td>random data</td>
            <td>random data</td>
            <td>random data</td>
            <td>random data</td>
            <td>random data</td>
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
</style>
