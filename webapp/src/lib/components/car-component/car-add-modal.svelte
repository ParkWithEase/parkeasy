<script lang="ts">
    import { Modal, TextInput, Form, InlineNotification } from 'carbon-components-svelte';
    import { CarModalState } from '$lib/enum/car-model';

    export let state: CarModalState;
    export let errorMessage: string;

    let form: HTMLFormElement | null;
</script>

<Modal
    open={state == CarModalState.ADD}
    modalHeading="New Car"
    primaryButtonText="Confirm"
    secondaryButtonText="Cancel"
    on:click:button--secondary={() => (state = CarModalState.NONE)}
    on:close={() => {
        form?.reset();
        errorMessage = '';
        state = CarModalState.NONE;
    }}
    on:open
    on:click:button--primary={() => form?.requestSubmit()}
>
    {#if errorMessage}
        <InlineNotification
            title="Error:"
            bind:subtitle={errorMessage}
            on:close={() => (errorMessage = '')}
        />
    {/if}
    <Form on:submit bind:ref={form}>
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
