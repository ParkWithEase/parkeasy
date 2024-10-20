<script lang="ts">
    import { Form, Modal, TextInput, Checkbox } from 'carbon-components-svelte';
    export let openState: boolean;
</script>

<Modal
    bind:open={openState}
    modalHeading="Add spots"
    primaryButtonText="Confirm"
    secondaryButtonText="Cancel"
    on:click:button--secondary={() => (openState = false)}
    on:close={() => document.getElementById('spot-create-form').reset()}
    on:click:button--primary={() => document.getElementById('spot-create-form').requestSubmit()}
>
    <Form on:submit id="spot-create-form">
        <TextInput
            required
            labelText="Street address"
            name="street-address"
            placeholder="Street Address"
        />
        <TextInput required labelText="City" name="city" placeholder="City" />
        <TextInput required labelText="Country" name="country-code" placeholder="Country Code" />
        <TextInput required labelText="Postal code" name="postal-code" placeholder="Postal code" />
        <p>Utilities</p>
        <Checkbox name="shelter" labelText="shelter" />
        <Checkbox name="charging-station" labelText="Charging Station" />
        <Checkbox name="plug-in" labelText="Plug-in" />
    </Form>
</Modal>
