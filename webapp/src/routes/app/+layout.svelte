<script lang="ts">
    import NavigationBar from '$lib/components/navbar.svelte';
    import { Content } from 'carbon-components-svelte';

    let isSideNavOpen: boolean;
</script>

<NavigationBar />

<Content>
    <slot></slot>
</Content>
