<script lang="ts">
    import logo from '$lib/images/logo-text.png';
    import { Header, HeaderUtilities, HeaderGlobalAction } from 'carbon-components-svelte';
    import Logout from '../../../node_modules/carbon-icons-svelte/lib/Logout.svelte';
    import User from '../../../node_modules/carbon-icons-svelte/lib/User.svelte';
    import Map from '../../../node_modules/carbon-icons-svelte/lib/Map.svelte';
    import MobileAdd from 'carbon-icons-svelte/lib/MobileAdd.svelte';
    import Car from 'carbon-icons-svelte/lib/Car.svelte';
    import { BACKEND_SERVER } from '$lib/constants';
    import { goto } from '$app/navigation';

    async function logout() {
        try {
            const response = await fetch(`${BACKEND_SERVER}/auth`, {
                method: 'DELETE',
                credentials: 'include'
            });
            if (response.ok) {
                goto('/auth/login');
            } else {
                throw new Error("Can't log out for some reason");
            }
        } catch {
            throw new Error('Something went wrong');
        }
    }
</script>

<svelte:window />

<Header>
    <a href="/#"><img src={logo} alt="parkeasy-logo" class="text-logo" /></a>
    <HeaderUtilities>
        <HeaderGlobalAction
            iconDescription="Explore open spots"
            tooltipAlignment="start"
            icon={Map}
            href="/#"
        />
        <HeaderGlobalAction iconDescription="your spot" icon={MobileAdd} href="/app/your-spots" />
        <HeaderGlobalAction iconDescription="Your Cars" icon={Car} href="/app/your-cars" />
        <HeaderGlobalAction
            iconDescription="Profile"
            icon={User}
            href="/app/profile/user-profile"
        />
        <HeaderGlobalAction
            iconDescription="Log out"
            tooltipAlignment="end"
            icon={Logout}
            href="/"
            on:click={logout}
        />
    </HeaderUtilities>
</Header>

<style>
    .text-logo {
        max-height: 4rem;
    }
</style>
