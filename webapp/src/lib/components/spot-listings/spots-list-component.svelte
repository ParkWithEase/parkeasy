<script lang="ts">
    import { ChargingStation, PlugFilled, Home } from "carbon-icons-svelte";
    import { TooltipIcon } from 'carbon-components-svelte';
    import DetailModal from "./detail-modal.svelte";
    export let listing;

    let modalOpen = false;

    const handleClick = () => {
        modalOpen = true;
    };

</script>

<div class="listing-info" on:click={() => handleClick()}
                            on:keyup={() => {}}
                            role="button"
                            tabindex="0">
    <div class="address">
        <p class="listing-head">{listing.location.street_address}</p>
        <p>{listing.location.postal_code}</p>
    </div>
    <div class="features">
        <div class="feature">
            {#if listing.features.charging_station}
                <TooltipIcon 
                    tooltipText="Charging Station Available"
                    direction="right"
                    align="end">
                    <i><ChargingStation size={20}/></i>
                </TooltipIcon>
            {/if}
        </div>
        <div class="feature">
            {#if listing.features.plug_in}
                <TooltipIcon 
                    tooltipText="Plug-in Available"
                    direction="right"
                    align="end">
                    <i class="icon"><PlugFilled size={20}/></i>
                </TooltipIcon>
                
            {/if}
        </div>
        <div class="feature">
            {#if listing.features.shelter}
            <TooltipIcon 
                tooltipText="Shelter Available"
                direction="right"
                align="end">
                <i class="icon"><Home size={20}/></i>
            </TooltipIcon>
            {/if}
        </div>
    </div>
</div>
<DetailModal bind:open={modalOpen} listing={listing} />
<style>
    .listing-info {
        flex: 1;
        display: flex;
        flex-direction: column;
        padding: 10px;
        border: 1px solid #ddd;
        margin: 5px;
        border-radius: 4px;
        background-color: #f8f9fa;
    }
    .listing-head {
        font-size: 1.5rem;
        border-bottom: 2px solid #333;
        margin-bottom: 5px;
    }
    .address {
        margin-bottom: 10px;
    }
    .features {
        display: flex;
        flex-direction: row;
        margin-left: 0.5rem;
    }
    .feature {
        display: flex;
        align-items: center;
        font-size: 1rem;
        color: #1ac24c;
        margin: 0rem 0.2rem;
    }
    .icon {
        color: #1ac24c; /* Adjust color as desired */
    }
</style>