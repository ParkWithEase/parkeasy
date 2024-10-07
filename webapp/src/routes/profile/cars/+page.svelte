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
        }
    ];

    function selectCarIndex(index: number) {
        isViewModalOpen = true;
        selectedCarIndex = index;
    }

    let carLicensePlateInput: string;
    let carColorInput: string;
    let carModelInput: string;
    let carMakeInput: string;

    function handleCreate(event: Event) {
        event.preventDefault();
        let newCar = {
            license_plate: carLicensePlateInput,
            color: carColorInput,
            model: carModelInput,
            make: carMakeInput
        };

        carList = [...carList, newCar];
        isAddModalOpen = false;
        clearUp();
    }

    function handleEdit(event: Event) {
        event.preventDefault();
        carList[selectedCarIndex] = {
            license_plate: carLicensePlateInput,
            color: carColorInput,
            model: carModelInput,
            make: carMakeInput
        };
        clearUpEditMode();
    }

    function clearUpEditMode() {
        isEditModalOpen = false;
        clearUp();
    }
    function clearUp() {
        carLicensePlateInput = '';
        carColorInput = '';
        carModelInput = '';
        carMakeInput = '';
    }

    function loadFormInput(car: Car) {
        carLicensePlateInput = car.license_plate;
        carColorInput = car.color;
        carModelInput = car.model;
        carMakeInput = car.make;
    }

    function handleDelete() {
        if (confirm('Are you sure you want to remove this car?')) {
            carList.splice(selectedCarIndex, 1);
            carList = [...carList];
            isViewModalOpen = false;
        }
    }
</script>

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
    on:close={() => clearUpEditMode()}
>
    {@const selectedCar = carList[selectedCarIndex]}
    {#if isEditModalOpen}
        {#if selectedCar != undefined}
            <Form on:submit={handleEdit}>
                <input
                    type="image"
                    class="modal-icon"
                    alt="cancel edit button"
                    src={cancelButton}
                    on:click={() => clearUpEditMode()}
                />
                <TextInput
                    required
                    labelText="License plate"
                    bind:value={carLicensePlateInput}
                    placeholder={selectedCar.license_plate}
                />
                <TextInput
                    required
                    labelText="Color"
                    bind:value={carColorInput}
                    placeholder={selectedCar.color}
                />
                <TextInput
                    required
                    labelText="Model"
                    bind:value={carModelInput}
                    placeholder={selectedCar.model}
                />
                <TextInput
                    required
                    labelText="Make"
                    bind:value={carMakeInput}
                    placeholder={selectedCar.make}
                />
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
                        loadFormInput(selectedCar);
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
    on:close={clearUp}
>
    <Form on:submit={handleCreate}>
        <TextInput
            required
            labelText="License plate"
            bind:value={carLicensePlateInput}
            placeholder="Your car license plate"
        />
        <TextInput
            required
            labelText="Color"
            bind:value={carColorInput}
            placeholder="Car's color"
        />
        <TextInput
            required
            labelText="Model"
            bind:value={carModelInput}
            placeholder="Car's model"
        />
        <TextInput required labelText="Make" bind:value={carMakeInput} placeholder="Car's make" />
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
</style>
