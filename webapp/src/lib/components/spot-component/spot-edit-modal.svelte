<script lang="ts">
    import { Form, Modal, TextInput, Checkbox } from 'carbon-components-svelte';
    import type { components } from '$lib/sdk/schema';
    type ParkingSpot = components['schemas']['ParkingSpot'];
    export let openState: boolean;
    export let spotInfo: ParkingSpot;
    let form: HTMLFormElement | null;
</script>

<Modal
    bind:open={openState}
    modalHeading="Edit spots"
    primaryButtonText="Confirm"
    secondaryButtonText="Cancel"
    on:click:button--secondary={() => (openState = false)}
    on:click:button--primary={() => form?.requestSubmit()}
>
    <Form on:submit bind:ref={form}>
        <TextInput
            required
            labelText="Street address"
            name="street-address"
            placeholder="Street Address"
            value={spotInfo.location.street_address}
        />
        <TextInput
            required
            labelText="City"
            name="city"
            placeholder="City"
            value={spotInfo.location.city}
        />
        <TextInput
            required
            labelText="State/Province"
            name="state"
            placeholder="State/Province"
            value={spotInfo.location.state}
        />
        <TextInput
            required
            labelText="Country"
            name="country-code"
            placeholder="Country Code"
            value={spotInfo.location.country_code}
        />
        <TextInput
            required
            labelText="Postal code"
            name="postal-code"
            placeholder="Postal code"
            value={spotInfo?.location.postal_code}
        />
        <p>Utilities</p>
        <Checkbox name="shelter" labelText="shelter" checked={spotInfo.features?.shelter} />
        <Checkbox name="plug-in" labelText="Plug-in" checked={spotInfo.features?.plug_in} />
        <Checkbox
            name="charging-station"
            labelText="Charging Station"
            checked={spotInfo.features?.charging_station}
        />
    </Form>
</Modal>
