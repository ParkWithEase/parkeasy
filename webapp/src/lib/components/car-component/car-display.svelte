<script lang="ts">
    import type { components } from '$lib/sdk/schema';

    type Car = components['schemas']['CarDetails'];
    export let car: Car;
</script>

<div class="car-info">
    <div>
        <p class="car-license-plate">{car.license_plate}</p>
    </div>
    <div>
        <p class="car-data">Color: {car.color}</p>
        <p class="car-data">Model: {car.model}</p>
        <p class="car-data">Make: {car.make}</p>
    </div>
</div>

<style>
    .car-info {
        padding-left: 1rem;
        flex: 1;
        display: flex;
        flex-direction: column;
        border: 1px solid #cfcfcfcf;
    }

    .car-data {
        font-size: 1rem;
    }

    .car-license-plate {
        font-size: 1.5rem;
    }
</style>
