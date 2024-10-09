<script lang="ts">
    import type { Car } from '$lib/types/car/car';
    import parkingImg from '$lib/images/background.png';
    import addButton from '$lib/images/add-button.png';
    import editButton from '$lib/images/edit.png';
    import deleteButton from '$lib/images/bin.png';
    import cancelButton from '$lib/images/cancel.png';
    import { Form, Modal, TextInput } from 'carbon-components-svelte';
    import SubmitButton from '$lib/components/submit-button.svelte';

    let isViewModalOpen: boolean = false;
    let isEditModalOpen: boolean = false;
    let isAddModalOpen: boolean = false;
    let selectedCarIndex: number;
    let carList: Car[] = [
        {
            license_plate: '1234565',
            color: 'Red',
            model: 'nuddole',
            make: 'ferrari'
        },
        {
            license_plate: '78945613',
            color: 'Blue',
            model: 'Lefton',
            make: 'Rocky'
        },
        {
            license_plate: '78945613',
            color: 'Blue',
            model: 'Lefton',
            make: 'Rocky'
        },
        {
            license_plate: '78945613',
            color: 'Blue',
            model: 'Lefton',
            make: 'Rocky'
        },
        {
            license_plate: '78945613',
            color: 'Blue',
            model: 'Lefton',
            make: 'Rocky'
        },
        {
            license_plate: '78945613',
            color: 'Blue',
            model: 'Lefton',
            make: 'Rocky'
        }
    ];

    function selectCarIndex(index: number) {
        isViewModalOpen = true;
        selectedCarIndex = index;
    }

    function handleCreate(event: Event) {
        event.preventDefault();
        const formData = new FormData(event.target as HTMLFormElement);
        const newCar = {
            license_plate: formData.get('license-plate') as string,
            color: formData.get('color') as string,
            model: formData.get('model') as string,
            make: formData.get('make') as string
        };
        carList = [...carList, newCar];
        isAddModalOpen = false;
    }

    function handleEdit(event: Event) {
        event.preventDefault();
        const formData = new FormData(event.target as HTMLFormElement);
        carList[selectedCarIndex] = {
            license_plate: formData.get('license-plate') as string,
            color: formData.get('color') as string,
            model: formData.get('model') as string,
            make: formData.get('make') as string
        };
        isEditModalOpen = false;
    }

    function handleDelete() {
        if (confirm('Are you sure you want to remove this car?')) {
            carList.splice(selectedCarIndex, 1);
            carList = [...carList];
            isViewModalOpen = false;
        }
    }

    function resetForm(id: string) {
        let form = document.getElementById(id) as HTMLFormElement;
        if (form == undefined)
    {
        return;
    }
        form.reset();
    }
</script>

<div class="listing-container">
    {#key carList}
    {#each carList as car, index}
        <!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
        <div class="car-info-container" on:click={() => selectCarIndex(index)}>
            <img class="car-image" src={parkingImg} alt="a car" />
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
        </div>
    {/each}
{/key}
</div>

<input
    type="image"
    alt="add car"
    src={addButton}
    class="add-button"
    on:click={() => (isAddModalOpen = true)}
/>

<Modal
    bind:open={isViewModalOpen}
    modalHeading="Car information"
    passiveModal
    on:open
    on:close={() => {
        isEditModalOpen = false;
        resetForm('edit-form');
    }}
>
    {@const selectedCar = carList[selectedCarIndex]}
    {#if isEditModalOpen}
        {#if selectedCar != undefined}
            <Form id="edit-form" on:submit={handleEdit}>
                <input
                    type="image"
                    class="modal-icon"
                    alt="cancel edit button"
                    src={cancelButton}
                    on:click={() => {
                        isEditModalOpen = false;
                        resetForm('edit-form');
                    }}
                />
                <TextInput
                    required
                    labelText="License plate"
                    name="license-plate"
                    value={selectedCar.license_plate}
                />
                <TextInput required labelText="Color" name="color" value={selectedCar.color} />
                <TextInput required labelText="Model" name="model" value={selectedCar.model} />
                <TextInput required labelText="Make" name="make" value={selectedCar.make} />
                <SubmitButton buttonText={'Submit'} />
            </Form>
        {/if}
    {:else if selectedCar != undefined}
        <div class="car-info">
            <div>
                <p class="car-license-plate" style="border-bottom: 2px solid black;">
                    {selectedCar.license_plate}
                </p>
            </div>
            <div style="position: relative">
                <p class="car-data">Color: {selectedCar.color}</p>
                <p class="car-data">Model: {selectedCar.model}</p>
                <p class="car-data">Make: {selectedCar.make}</p>
            </div>
            <div>
                <input
                    type="image"
                    class="modal-icon"
                    src={editButton}
                    alt="edit button"
                    on:click={() => {
                        isEditModalOpen = true;
                    }}
                />
                <input
                    type="image"
                    class="modal-icon"
                    style="color: red;"
                    src={deleteButton}
                    alt="delete button"
                    on:click={handleDelete}
                />
            </div>
        </div>
    {/if}
</Modal>

<Modal
    bind:open={isAddModalOpen}
    modalHeading="Car information"
    passiveModal
    on:open
    on:close={() => resetForm('create-form')}
>
    <Form id="create-form" on:submit={handleCreate}>
        <TextInput
            required
            labelText="License plate"
            name="license-plate"
            placeholder="Your car license plate"
        />
        <TextInput required labelText="Color" name="color" placeholder="Car's color" />
        <TextInput required labelText="Model" name="model" placeholder="Car's model" />
        <TextInput required name="make" labelText="Make" placeholder="Car's make" />
        <SubmitButton buttonText={'Submit'} />
    </Form>
</Modal>

<style>
    .car-info-container {
        display: flex;
        flex-direction: row;
        margin: 1rem;
        padding: 1rem;
        gap: 1rem;
        justify-content: stretch;
        border-radius: 20px;
        background-color: rgba(0, 0, 0, 0.5);
        color: whitesmoke;
    }

    .car-info {
        flex: 1;
        display: flex;
        flex-direction: column;
    }

    .car-image {
        max-height: 12%;
        max-width: 16%;
    }

    .car-license-plate {
        font-size: 1.5rem;
        border-bottom: 2px solid white;
    }

    .car-data {
        font-size: 1rem;
    }

    .add-button {
        position: absolute;
        right: 5%;
        bottom: 5%;
        max-width: 3rem;
    }

    .modal-icon {
        max-width: 2rem;
    }

    .listing-container {
        max-height: inherit;
        overflow-y: scroll;
    }
</style>
