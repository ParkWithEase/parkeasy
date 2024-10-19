<script lang="ts">
    import type { Car } from '$lib/types/car/car';
    import {
        ComposedModal,
        ModalBody,
        ModalHeader,
        ModalFooter,
        Form,
        TextInput,
        InlineNotification
    } from 'carbon-components-svelte';

    import { createEventDispatcher } from 'svelte';

    const dispatch = createEventDispatcher();

    export let openState: boolean;
    export let isEdit: boolean = false;
    export let errorMessage: string;
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
        isEdit = false;
    }

    function submitForm() {
        let form: HTMLFormElement | null = document.getElementById(
            'car-edit-form'
        ) as HTMLFormElement | null;
        if (form) {
            form.requestSubmit();
        }
    }

    export let carInfo: Car;
</script>

{#if carInfo}
    <ComposedModal
        preventCloseOnClickOutside={true}
        bind:open={openState}
        on:click:button--primary={isEdit ? () => submitForm() : () => (isEdit = true)}
        on:close={() => {
            cleanUp();
            isEdit = false;
        }}
    >
        {#if isEdit}
            <ModalHeader title="Edit car" />
        {:else}
            <ModalHeader title="Car Info" />
        {/if}

        {#if isEdit}
            <ModalBody hasForm>
                {#if errorMessage}
                    <InlineNotification title="Error:" bind:subtitle={errorMessage} />
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

        {#if isEdit}
            <ModalFooter
                primaryButtonText="Confirm"
                secondaryButtonText="Cancel"
                on:click:button--secondary={() => {
                    cleanUp();
                    openState = true;
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
