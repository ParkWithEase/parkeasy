<script lang="ts">
    import {
        Button,
        Checkbox,
        Form,
        Select,
        SelectItem,
        TextInput,
        Search,
        Modal
    } from 'carbon-components-svelte';
    import type { AddressResult } from '$lib/types/address/address';

    const apiKey = import.meta.env.VITE_GEOCODING_API_KEY;

    let searchQuery = '';
    let debounceTimer: number | undefined;
    let results: AddressResult[];
    let dropdownOpen: boolean = false;
    let warningModalOpen: boolean = false;
    let warningText: string = '';

    export let isReadOnly: boolean;
    export let isVerified: boolean;

    export let houseNumber: number;
    export let streetAddress: string;
    export let city: string;
    export let state: string;
    export let country: string;
    export let countryCode: string;
    export let postalCode: string;
    export let hasShelter: boolean;
    export let hasPlugIn: boolean;
    export let hasChargingStation: boolean;

    // Province data mapping
    const provinceData: Map<string, string> = new Map([
        ['MB', 'Manitoba'],
        ['ON', 'Ontario'],
        ['AB', 'Alberta'],
        ['QC', 'Quebec'],
        ['NS', 'Nova Scotia'],
        ['BC', 'British Columbia'],
        ['NL', 'Newfoundland and Labrador'],
        ['PE', 'Prince Edward Island'],
        ['SK', 'Saskatchewan'],
        ['YT', 'Yukon'],
        ['NU', 'Nunavut'],
        ['NT', 'Northwest Territories']
    ]);

    // Function to fetch address suggestions
    const searchAddress = () => {
        const url = `https://api.geoapify.com/v1/geocode/search?text=${encodeURIComponent(searchQuery)}&apiKey=${apiKey}`;

        fetch(url)
            .then((response) => response.json())
            .then((data) => {
                results = data.features ?? [];
                dropdownOpen = true;
            })
            .catch((error) => {
                console.error('Error fetching data:', error);
            });
    };

    // Handle input changes and debounce the search
    const handleInput = () => {
        dropdownOpen = true;
        isVerified = false; // Reset verification on input change
        clearTimeout(debounceTimer);
        debounceTimer = window.setTimeout(() => {
            if (searchQuery.length > 3) {
                searchAddress();
            }
        }, 500); // 0.5s delay
    };

    // Handle selection of an address from dropdown
    const handleSelect = (location: AddressResult) => {
        searchQuery = location.properties.formatted;
        results = [];
        dropdownOpen = false;

        if (location.properties.country_code !== 'ca') {
            warningModalOpen = true;
            warningText = 'This application supports locations in Canada only!';
        } else {
            houseNumber = location.properties.housenumber ?? '';
            streetAddress = `${location.properties.housenumber} ${location.properties.street}`;
            city = location.properties.city;
            state = location.properties.state_code;
            countryCode = location.properties.country_code?.toUpperCase();
            country = location.properties.country;
            postalCode = location.properties.postcode?.replace(/\s/g, '');
            isVerified = false; // Reset verification after selection
            isReadOnly = false;
        }
    };

    // Address verification using Geoapify API
    const verifyAddress = () => {
        const url = `https://api.geoapify.com/v1/geocode/search?text=${encodeURIComponent(
            `${houseNumber} ${streetAddress}, ${city}, ${state}, ${postalCode}, ${country}`
        )}&apiKey=${apiKey}`;

        fetch(url)
            .then((response) => response.json())
            .then((result) => {
                if (result.features.length === 0) {
                    warningModalOpen = true;
                    warningText = 'Address verification failed. Please check your inputs!';
                    return;
                }

                const matchedAddress = result.features[0].properties;
                // Validate all fields match exactly
                const isMatch =
                    matchedAddress.housenumber + ' ' +  matchedAddress.street?.toLowerCase() ===
                        streetAddress.toLowerCase() &&
                    matchedAddress.city?.toLowerCase() === city.toLowerCase() &&
                    matchedAddress.state_code?.toLowerCase() === state.toLowerCase() &&
                    matchedAddress.postcode?.replace(/\s/g, '').toLowerCase() ===
                        postalCode.toLowerCase();
                if (isMatch) {
                    isVerified = true;
                } else {
                    warningModalOpen = true;
                    warningText =
                        'Address verification failed. The entered address does not match the verified address!';
                }
            })
            .catch((error) => {
                console.log('error', error);
                warningModalOpen = true;
                warningText = 'An error occurred during verification. Please try again!';
            });
    };

    // Reset verification status if fields are edited
    const handleEdit = () => {
        if (isVerified) {
            isVerified = false; // Reset verification if user edits fields
        }
    };

    // Handle click outside dropdown to close it
    const handleClickOutside = (event: Event) => {
        if (
            !(event.target as HTMLElement).closest('.dropdown') &&
            !(event.target as HTMLElement).closest('.search-input')
        ) {
            dropdownOpen = false;
        }
    };

    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    export let handleSubmit = (_: Event) => {};
</script>

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions-->
<div on:click={handleClickOutside}>
    <Form on:submit={handleSubmit}>
        {#if !isReadOnly}
            <Search
                type="text"
                bind:value={searchQuery}
                placeholder="Enter an address... [At least 3 characters required]"
                name="search"
                on:input={handleInput}
            />
            {#if results?.length > 0 && dropdownOpen}
                <ul class="dropdown">
                    {#each results as result}
                        <!-- svelte-ignore a11y-click-events-have-key-events a11y-no-noninteractive-element-interactions -->
                        <li on:click={() => handleSelect(result)}>
                            {result.properties.formatted}
                        </li>
                    {/each}
                </ul>
            {/if}
        {/if}
        <Modal
            open={warningModalOpen}
            modalHeading="Error - Please check Address"
            primaryButtonText="Re-enter address"
            on:close
            on:open
            on:submit={() => (warningModalOpen = false)}
        >
            {warningText}
        </Modal>
        <TextInput
            required
            name="street-address"
            placeholder="House # & Street Address"
            readonly={isReadOnly}
            bind:value={streetAddress}
            on:input={handleEdit}
            on:change={handleEdit}
        />
        <TextInput
            required
            name="city"
            placeholder="City"
            readonly={isReadOnly}
            bind:value={city}
            on:change={handleEdit}
        />
        <TextInput
            required
            name="postal-code"
            placeholder="Postal code"
            readonly={isReadOnly}
            bind:value={postalCode}
            on:change={handleEdit}
        />
        <Select
            style={isReadOnly ? 'pointer-events: none;' : ' '}
            labelText="Province"
            bind:selected={state}
        >
            {#each provinceData.keys() as key}
                <SelectItem value={key} text={provinceData.get(key)} />
            {/each}
        </Select>

        <Select
            style={isReadOnly ? 'pointer-events: none;' : ' '}
            labelText="Country"
            bind:selected={countryCode}
        >
            <SelectItem value="CA" text="Canada" />
        </Select>
        <p>Utilities</p>
        <Checkbox
            name="shelter"
            style={isReadOnly ? 'pointer-events: none;' : ' '}
            labelText="shelter"
            bind:checked={hasShelter}
        />
        <Checkbox
            name="plug-in"
            labelText="Plug-in"
            style={isReadOnly ? 'pointer-events: none;' : ' '}
            bind:checked={hasPlugIn}
        />
        <Checkbox
            name="charging-station"
            labelText="Charging Station"
            style={isReadOnly ? 'pointer-events: none;' : ' '}
            bind:checked={hasChargingStation}
        />

        {#if !isReadOnly}
            <!-- Verify button -->
            <Button type="button" on:click={verifyAddress} disabled={isVerified}>
                Verify Address
            </Button>
            <!-- Submit button -->
            <Button type="submit" disabled={!isVerified}>Submit</Button>
        {/if}
    </Form>
</div>

<style>
    .dropdown {
        position: absolute;
        background-color: white;
        border: 1px solid #868686;
        width: 100%;
        max-height: 10rem;
        overflow-y: auto;
        z-index: 1000;
    }
    .dropdown li {
        padding: 0.5rem;
        cursor: pointer;
    }
    .dropdown li:hover {
        background-color: #f0f0f0;
    }
</style>
