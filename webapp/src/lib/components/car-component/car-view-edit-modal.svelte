<script lang="ts">
    import {
        ComposedModal,
        ModalBody,
        ModalHeader,
        ModalFooter,
        Form,
        TextInput,
        InlineNotification
    } from 'carbon-components-svelte';
    import { CarModalState } from '$lib/enum/car-model';

    import { createEventDispatcher } from 'svelte';
    import type { components } from '$lib/sdk/schema';

    const dispatch = createEventDispatcher();

    type Car = components['schemas']['CarDetails'];
    export let state: CarModalState;
    export let errorMessage: string;
    export let carInfo: Car | null;
    function deleteCar() {
        dispatch('delete', {
            text: 'Delete'
        });
    }

    function cleanUp() {
        let form: HTMLFormElement | null = document.getElementById(
            'car-edit-form'
        ) as HTMLFormElement | null;
        if (form) {
            form.reset();
        }
        state = CarModalState.VIEW;
    }

    function submitForm() {
        let form: HTMLFormElement | null = document.getElementById(
            'car-edit-form'
        ) as HTMLFormElement | null;
        if (form) {
            form.requestSubmit();
        }
    }
</script>

{#if carInfo}
    <ComposedModal
        preventCloseOnClickOutside={true}
        open={state == CarModalState.EDIT || state == CarModalState.VIEW}
        on:click:button--primary={state == CarModalState.EDIT
            ? () => submitForm()
            : () => (state = CarModalState.EDIT)}
        on:close={() => {
            cleanUp();
            errorMessage = '';
            state = CarModalState.NONE;
        }}
    >
        {#if state == CarModalState.EDIT}
            <ModalHeader title="Edit car" />
        {:else}
            <ModalHeader title="Car Info" />
        {/if}

        {#if state == CarModalState.EDIT}
            <ModalBody hasForm>
                {#if errorMessage}
                    <InlineNotification
                        title="Error:"
                        bind:subtitle={errorMessage}
                        on:close={() => (errorMessage = '')}
                    />
                {/if}
                <Form on:submit id="car-edit-form">
                    <TextInput
                        required
                        labelText="License plate"
                        name="license-plate"
                        placeholder="Your car license plate"
                        value={carInfo.license_plate}
                    />
                    <TextInput
                        required
                        labelText="Color"
                        name="color"
                        placeholder="Car's color"
                        value={carInfo.color}
                    />
                    <TextInput
                        required
                        labelText="Model"
                        name="model"
                        placeholder="Car's model"
                        value={carInfo.model}
                    />
                    <TextInput
                        required
                        name="make"
                        labelText="Make"
                        placeholder="Car's make"
                        value={carInfo.make}
                    />
                </Form>
            </ModalBody>
        {:else}
            <ModalBody>
                <div>
                    <p
                        class="car-license-plate"
                        style="border-bottom: 2px solid black; font-size: 1.5rem"
                    >
                        {carInfo.license_plate}
                    </p>
                </div>
                <div style="position: relative">
                    <p class="car-data">Color: {carInfo.color}</p>
                    <p class="car-data">Model: {carInfo.model}</p>
                    <p class="car-data">Make: {carInfo.make}</p>
                </div>
            </ModalBody>
        {/if}

        {#if state == CarModalState.EDIT}
            <ModalFooter
                primaryButtonText="Confirm"
                secondaryButtonText="Cancel"
                on:click:button--secondary={() => {
                    cleanUp();
                    state = CarModalState.VIEW;
                }}
            />
        {:else}
            <ModalFooter
                primaryButtonText="Edit"
                secondaryButtonText="Delete"
                on:click:button--secondary={() => {
                    deleteCar();
                }}
            />
        {/if}
    </ComposedModal>
{/if}
