<script lang="ts">
    import NavigationBar from '$lib/components/navbar.svelte';
    import { Content } from 'carbon-components-svelte';
</script>

<NavigationBar />

<Content>
    <slot></slot>
</Content>
