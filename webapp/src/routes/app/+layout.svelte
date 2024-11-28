<script lang="ts">
    import NavigationBar from '$lib/components/navbar.svelte';
    import { RESPONSE_WIDTH } from '$lib/constants';
    import { Content } from 'carbon-components-svelte';
    import { responsiveState } from '$lib/stores/responsive';
    import { navigating } from '$app/stores';
    let innerWidth: number = 0;
    $: responsiveState.set(innerWidth < RESPONSE_WIDTH);
</script>

<svelte:window bind:innerWidth />

<NavigationBar />
<Content>
    {#await $navigating?.complete}
        <div>Loading.</div>
    {:then}
        <slot></slot>
    {/await}
</Content>
