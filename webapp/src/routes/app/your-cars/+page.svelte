<script lang="ts">
    import type { Car } from '$lib/types/car/car';
    import { InlineNotification } from 'carbon-components-svelte';
    import CarDisplay from '$lib/components/car-component/car-display.svelte';
    import CarAddModal from '$lib/components/car-component/car-add-modal.svelte';
    import CarViewEditModal from '$lib/components/car-component/car-view-edit-modal.svelte';
    import { Button } from 'carbon-components-svelte';
    import Add from 'carbon-icons-svelte/lib/Add.svelte';
    import { BACKEND_SERVER } from '$lib/constants';

    import { getErrorMessage } from '$lib/utils/error-handler';
    import { intersectionObserver } from '@sveu/browser';

    const regex = /[?&]after=([^&?]+)/;

    let isViewEditModalOpen: boolean = false;
    let isEditModalOpen: boolean = false;
    let isAddModalOpen: boolean = false;
    let selectedCarID: string;
    let selectedCarInfo: Car;
    let errorMessage: string;
    let afterToken: string = '';
    let canLoadMore: boolean = true;
    let loadLock: boolean = false;
    let target: HTMLElement | null = null;
    let root: HTMLElement | null = null;

    let data = [];

    function selectCarIndex(index: string, CarInfo: Car) {
        isViewEditModalOpen = true;
        selectedCarID = index;
        selectedCarInfo = CarInfo;
    }
    root = this;

    $: target,
        intersectionObserver(
            target,
            ([{ isIntersecting }]) => {
                if (canLoadMore && !loadLock && isIntersecting) {
                    loadLock = true;
                    getPage();
                }
            },
            { root }
        );

    export async function getPage() {
        const response = await fetch(`${BACKEND_SERVER}/cars?after=${afterToken}&count=5`, {
            credentials: 'include'
        });
        try {
            if (response.ok) {
                data = [...data, ...(await response.json())];

                if (response.headers.get('link') != null) {
                    let param = response.headers.get('link')?.match(regex);
                    afterToken = param != null ? param[1] : '';
                } else {
                    canLoadMore = false;
                    afterToken = '';
                }
                loadLock = false;
            } else {
                const errorDetails = await response.json();
                errorMessage = getErrorMessage(errorDetails);
            }
        } catch (err) {
            errorMessage = 'Something wrong happen ' + err;
        }
    }

    export async function handleCreate(event: Event) {
        event.preventDefault();
        const formData = new FormData(event.target as HTMLFormElement);
        console.log(formData.get('color'));
        try {
            const response = await fetch(`${BACKEND_SERVER}/cars`, {
                method: 'POST',
                credentials: 'include',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    license_plate: formData.get('license-plate') as string,
                    color: formData.get('color') as string,
                    model: formData.get('model') as string,
                    make: formData.get('make') as string
                })
            });

            if (response.ok) {
                const new_car = await response.json();

                //temporary code to see update
                data = [...data, { details: new_car.details, id: new_car.id }];
                errorMessage = '';
            } else {
                const errorDetails = await response.json();
                errorMessage = getErrorMessage(errorDetails);
            }
        } catch (err) {
            errorMessage = 'Something wrong happen ' + err;
        }

        isAddModalOpen = false;
    }

    export async function handleDelete() {
        if (confirm('Are you sure you want to remove this car?')) {
            try {
                const response = await fetch(`${BACKEND_SERVER}/cars/${selectedCarID}`, {
                    method: 'DELETE',
                    credentials: 'include'
                });
                if (response.ok) {
                    errorMessage = '';
                    isViewEditModalOpen = false;
                    data = data.filter(function (item) {
                        return item.id !== selectedCarID;
                    });
                } else {
                    const errorDetails = await response.json();
                    errorMessage = getErrorMessage(errorDetails);
                }
            } catch (err) {
                errorMessage = 'Something wrong happen ' + err;
            }
        } else {
            isViewEditModalOpen = true;
        }
    }

    export async function handleEdit(event: Event) {
        event.preventDefault();
        const formData = new FormData(event.target as HTMLFormElement);
        try {
            const response = await fetch(`${BACKEND_SERVER}/cars/${selectedCarID}`, {
                method: 'PUT',
                credentials: 'include',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    license_plate: formData.get('license-plate') as string,
                    color: formData.get('color') as string,
                    model: formData.get('model') as string,
                    make: formData.get('make') as string
                })
            });
            if (response.ok) {
                errorMessage = '';
                isEditModalOpen = false;

                //temporary code to see update
                selectedCarInfo.license_plate = formData.get('license-plate') as string;
                selectedCarInfo.color = formData.get('color') as string;
                selectedCarInfo.model = formData.get('model') as string;
                selectedCarInfo.make = formData.get('make') as string;
                data = [...data];
            } else {
                isEditModalOpen = true;
                const errorDetails = await response.json();
                errorMessage = getErrorMessage(errorDetails);
            }
        } catch (err) {
            errorMessage = 'Something wrong happen ' + err;
        }
    }
</script>

{#if errorMessage}
    <InlineNotification title="Error:" bind:subtitle={errorMessage} />
{/if}

<div class="button-container" style="">
    <Button style="margin: 1rem;" icon={Add} on:click={() => (isAddModalOpen = true)}
        >New Car</Button
    >
</div>

<div>
    {#key data}
        {#each data as car}
            <!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
            <div on:click={() => selectCarIndex(car.id, car.details)}>
                <CarDisplay car={car.details}></CarDisplay>
            </div>
        {/each}
    {/key}

    {#if canLoadMore}
        <div bind:this={target}>Loading...</div>
    {/if}
</div>

<CarViewEditModal
    bind:openState={isViewEditModalOpen}
    bind:carInfo={selectedCarInfo}
    bind:isEdit={isEditModalOpen}
    on:submit={handleEdit}
    on:delete={handleDelete}
/>
<CarAddModal bind:openState={isAddModalOpen} on:submit={handleCreate} />

<style>
    .button-container {
        position: sticky;
        top: 0;
        background: white;
    }
</style>
