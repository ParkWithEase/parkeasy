<script lang="ts">
    import type { PageData } from './$types';
    import type { Car } from '$lib/types/car/car';
    import CarDisplay from '$lib/components/car-component/car-display.svelte';
    import CarAddModal from '$lib/components/car-component/car-add-modal.svelte';
    import CarViewEditModal from '$lib/components/car-component/car-view-edit-modal.svelte';
    import { Button } from 'carbon-components-svelte';
    import Add from 'carbon-icons-svelte/lib/Add.svelte';

    import { getErrorMessage } from '$lib/utils/error-handler';
    import IntersectionObserver from 'svelte-intersection-observer';
    import { newClient } from '$lib/utils/client';

    let isViewEditModalOpen: boolean = false;
    let isEditModalOpen: boolean = false;
    let isAddModalOpen: boolean = false;
    let selectedCarID: string;
    let selectedCarInfo: Car;
    let errorMessage: string;
    let loadTrigger: HTMLElement | null = null;
    let intersecting: boolean;

    let client = newClient();
    export let data: PageData;
    let canLoadMore = data.hasNext;

    function selectCarIndex(index: string, CarInfo: Car) {
        isViewEditModalOpen = true;
        selectedCarID = index;
        selectedCarInfo = CarInfo;
    }

    let loadLock = false;

    $: while (intersecting && canLoadMore && !loadLock) {
        loadLock = true;
        data.paging.next().then(({ value: { data: cars, hasNext } }) => {
            if (cars) {
                data.cars = [...data.cars, ...cars];
            }
            canLoadMore = !!hasNext;
            loadLock = false;
        });
    }

    function handleCreate(event: Event) {
        event.preventDefault();
        const formData = new FormData(event.target as HTMLFormElement);
        console.log(formData.get('color'));

        client
            .POST('/cars', {
                body: {
                    license_plate: formData.get('license-plate') as string,
                    color: formData.get('color') as string,
                    model: formData.get('model') as string,
                    make: formData.get('make') as string
                }
            })
            .then(({ data: new_car, error }) => {
                if (new_car) {
                    data.cars = [...data.cars, { details: new_car.details, id: new_car.id }];
                    errorMessage = '';
                    isAddModalOpen = false;
                }
                if (error) {
                    errorMessage = getErrorMessage(error);
                }
            })
            .catch((err) => {
                errorMessage = err;
            });
    }

    function handleDelete() {
        if (confirm('Are you sure you want to remove this car?')) {
            client
                .DELETE('/cars/{id}', {
                    params: {
                        path: { id: selectedCarID }
                    }
                })
                .then(({ error }) => {
                    if (error) {
                        errorMessage = getErrorMessage(error);
                    } else {
                        data.cars = data.cars?.filter(function (item) {
                            return item.id !== selectedCarID;
                        });
                        errorMessage = '';
                    }
                })
                .catch((err) => {
                    errorMessage = err;
                })
                .finally(() => {
                    isViewEditModalOpen = false;
                });
        }
    }

    function handleEdit(event: Event) {
        event.preventDefault();
        const formData = new FormData(event.target as HTMLFormElement);
        client
            .PUT('/cars/{id}', {
                params: {
                    path: { id: selectedCarID }
                },
                body: {
                    license_plate: formData.get('license-plate') as string,
                    color: formData.get('color') as string,
                    model: formData.get('model') as string,
                    make: formData.get('make') as string
                }
            })
            .then(({ data: change_car, error }) => {
                if (change_car) {
                    isEditModalOpen = false;
                    data.cars = data.cars?.map((car) => {
                        if (car.id == change_car.id) {
                            return change_car;
                        } else {
                            return car;
                        }
                    });
                    selectedCarInfo = change_car.details;
                    errorMessage = '';
                }
                if (error) {
                    isEditModalOpen = true;
                    errorMessage = getErrorMessage(error);
                }
            })
            .catch((err) => {
                errorMessage = err;
            });
    }
</script>

<div class="button-container" style="">
    <Button style="margin: 1rem;" icon={Add} on:click={() => (isAddModalOpen = true)}
        >New Car</Button
    >
</div>

<div>
    {#key data.cars}
        {#each data?.cars as car}
            <!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
            <div on:click={() => selectCarIndex(car.id, car.details)}>
                <CarDisplay car={car.details}></CarDisplay>
            </div>
        {/each}
    {/key}
</div>
<IntersectionObserver element={loadTrigger} bind:intersecting>
    {#if canLoadMore}
        <div bind:this={loadTrigger}>Loading...</div>
    {/if}
</IntersectionObserver>

<CarViewEditModal
    bind:openState={isViewEditModalOpen}
    bind:carInfo={selectedCarInfo}
    bind:isEdit={isEditModalOpen}
    on:submit={handleEdit}
    on:delete={handleDelete}
    bind:errorMessage
/>
<CarAddModal bind:openState={isAddModalOpen} on:submit={handleCreate} bind:errorMessage />

<style>
    .button-container {
        position: sticky;
        top: 0;
        background: white;
    }
</style>
