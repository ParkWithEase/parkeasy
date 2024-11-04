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
    let warningText : String = false;

    export let isReadOnly: boolean;

    export let houseNumber: string;
    export let streetAddress: string;
    export let city: string;
    export let state: string;
    export let country : string;
    export let countryCode: string;
    export let postalCode: string;
    export let hasShelter: boolean;
    export let hasPlugIn: boolean;
    export let hasChargingStation: boolean;

    const provinceData: Map<string, string> = new Map([
        ['MB', 'Manitoba'],
        ['ON', 'Ontario'],
        ['AB', 'Alberta'],
        ['QC', 'Quebec'],
        ['NS', 'Nova Scotia'],
        ['BC', 'British Columbia'],
        ['NL', 'Newfouundland and Labrador'],
        ['PE', 'Prince Edward Island'],
        ['SK', 'Saskatchewan'],
        ['YT', 'Yukon'],
        ['NU', 'Nunavut'],
        ['NT', 'Northwest Territories']
    ]);

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

    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    export let handleSubmit = (_: Event) => {};

    const handleInput = () => {
        dropdownOpen = true;
        clearTimeout(debounceTimer);
        debounceTimer = window.setTimeout(() => {
            if (searchQuery.length > 3) {
                searchAddress();
            }
        }, 500); // 0.5s delay
    };

    const handleSelect = (location: AddressResult) => {
        searchQuery = location.properties.formatted;
        results = [];
        dropdownOpen = false;
        if(location.properties.country_code !== 'ca'){
            warningModalOpen = true;
            warningText = 'This application supports locations in Canada only!'
        }
        houseNumber = location.properties.housenumber;
        streetAddress = location.properties.street;
        city = location.properties.city;
        state = location.properties.state_code;
        countryCode = location.properties.country_code;
        country = location.properties.country;
        postalCode = location.properties.postcode;
    };

    const handleClickOutside = (event: Event) => {
        if (
            !(event.target as HTMLElement).closest('.dropdown') &&
            !(event.target as HTMLElement).closest('.search-input')
        ) {
            dropdownOpen = false;
        }
    };

    // const verifyAddress = () => {
    //     fetch(`https://api.geoapify.com/v1/geocode/search?housenumber=${encodeURIComponent(houseNumber)}&street=${encodeURIComponent(streetAddress)}&postcode=${encodeURIComponent(postalCode)}
    //     &city=${encodeURIComponent(city)}&state=${encodeURIComponent(state)}&country=${encodeURIComponent(country)}&apiKey=${apiKey}`)
    //     .then(result => result.json()).then((result) => {
    //         let features = result.features || [];

    //         // To find a confidence level that works for you, try experimenting with different levels
    //         //Code referred from JSFiddle at: https://www.npmjs.com/package/@geoapify/geocoder-autocomplete
    //         const confidenceLevelToAccept = 0.25;
    //         features = features.filter(feature => feature.properties.rank.confidence >= confidenceLevelToAccept);

    //         if (features.length) {
    //             const foundAddress = features[0];
    //             if (foundAddress.properties.rank.confidence === 1) {
    //                 warningText = `We verified the address you entered. The formatted address is: ${foundAddress.properties.formatted}`;
    //             } else if (foundAddress.properties.rank.confidence > 0.5 && foundAddress.properties.rank.confidence_street_level === 1) {
    //                 warningText = `We have some doubts about the accuracy of the address: ${foundAddress.properties.formatted}`
    //             } else if (foundAddress.properties.rank.confidence_street_level === 1) {
    //                 warningText = `We can confirm the address up to street level: ${foundAddress.properties.formatted}`
    //             } else {
    //                 warningText = `We can only verify your address partially. The address we found is ${foundAddress.properties.formatted}`
    //             }
    //         } else {
    //             warningText = "We cannot find your address. Please check if you provided the correct address."
    //         }
    //         warningModalOpen = true;
    //     })
    // }
</script>

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions-->
<div on:click={handleClickOutside}>
    <Form on:submit={handleSubmit}>
        <Search
                type="text"
                bind:value={searchQuery}
                placeholder="Enter an address... [Atleast 3 characters required]"
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
        <Modal
            open={warningModalOpen}
            modalHeading = "Error - Please check Address"
            primaryButtonText="Re-enter address"
            on:close
            on:open
            on:submit={() => (warningModalOpen = false)}
        >{warningText}</Modal>
        <TextInput
            name="house-number"
            placeholder="House #"
            readonly={isReadOnly}
            bind:value={houseNumber}
        />
        <TextInput
            required
            name="street-address"
            placeholder="Street Address"
            readonly={isReadOnly}
            bind:value={streetAddress}
        />
        <TextInput
            required
            name="city"
            placeholder="City"
            readonly={isReadOnly}
            bind:value={city}
        />
        <TextInput
            required
            name="postal-code"
            placeholder="Postal code"
            readonly={isReadOnly}
            bind:value={postalCode}
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
        <!-- <Button on:click={verifyAddress()}>Verify Address</Button> -->
        {#if !isReadOnly}
            <Button type="submit">Submit</Button>
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
