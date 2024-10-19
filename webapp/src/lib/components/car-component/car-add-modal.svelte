<script lang="ts">
    import { Modal, TextInput, Form, InlineNotification } from 'carbon-components-svelte';

    export let openState: boolean;
    export let errorMessage: string;
</script>

<Modal
    bind:open={openState}
    modalHeading="New Car"
    primaryButtonText="Confirm"
    secondaryButtonText="Cancel"
    on:click:button--secondary={() => (openState = false)}
    on:close={() => document.getElementById('car-create-form').reset()}
    on:open
    on:click:button--primary={() => document.getElementById('car-create-form').requestSubmit()}
>
    {#if errorMessage}
        <InlineNotification title="Error:" bind:subtitle={errorMessage} />
    {/if}
    <Form on:submit id="car-create-form">
        <TextInput
            required
            labelText="License plate"
            name="license-plate"
            placeholder="Your car license plate"
        />
        <TextInput required labelText="Color" name="color" placeholder="Car's color" />
        <TextInput required labelText="Model" name="model" placeholder="Car's model" />
        <TextInput required name="make" labelText="Make" placeholder="Car's make" />
    </Form>
</Modal>
