<script lang="ts">
    import {
        Button,
        Checkbox,
        Form,
        Select,
        SelectItem,
        TextInput
    } from 'carbon-components-svelte';

    export let isReadOnly: boolean;
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    export let handleSubmit = (_: Event) => {};

    export let streetAddress: string;
    export let city: string;
    export let state: string;
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
</script>

<Form on:submit={handleSubmit}>
    <TextInput
        required
        labelText="Street address"
        name="street-address"
        placeholder="Street Address"
        readonly={isReadOnly}
        bind:value={streetAddress}
    />
    <TextInput
        required
        labelText="City"
        name="city"
        placeholder="City"
        readonly={isReadOnly}
        bind:value={city}
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
        labelText="Country Code"
        bind:selected={countryCode}
    >
        <SelectItem value="CA" text="Canada" />
    </Select>
    <TextInput
        required
        labelText="Postal code"
        name="postal-code"
        placeholder="Postal code"
        readonly={isReadOnly}
        bind:value={postalCode}
    />
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
        <Button type="submit">Submit</Button>
    {/if}
</Form>
