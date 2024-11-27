<script lang="ts">
    import logo from '$lib/images/logo-text.png';
    import {
        Header,
        HeaderUtilities,
        HeaderGlobalAction,
        HeaderAction,
        HeaderPanelLinks,
        HeaderPanelLink,
        HeaderPanelDivider,
        TooltipIcon
    } from 'carbon-components-svelte';
    import { BACKEND_SERVER } from '$lib/constants';
    import { goto } from '$app/navigation';
    import { responsiveState } from '$lib/stores/responsive';
    import { Logout, User, Map, MobileAdd, Car } from 'carbon-icons-svelte';

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

    let shouldReponse: boolean = false;
    responsiveState.subscribe((value) => {
        shouldReponse = value;
    });

    let isOpen: boolean = false;
    $: isOpen = shouldReponse && false;
</script>

<Header expansionBreakpoint={0}>
    <a href="/#"><img src={logo} alt="parkeasy-logo" class="text-logo" /></a>
    <HeaderUtilities>
        {#if shouldReponse}
            <HeaderAction bind:isOpen>
                <HeaderPanelLinks>
                    <HeaderPanelDivider>Pages</HeaderPanelDivider>

                    <HeaderPanelLink href="/#" on:click={() => (isOpen = false)}>
                        <TooltipIcon icon={Map} tooltipText="Explore open spots" />
                        Explore open spots
                    </HeaderPanelLink>

                    <HeaderPanelLink href="/app/your-spots" on:click={() => (isOpen = false)}>
                        <TooltipIcon icon={MobileAdd} tooltipText="Your spots" />Your spots
                    </HeaderPanelLink>

                    <HeaderPanelLink href="/app/your-cars" on:click={() => (isOpen = false)}
                        ><TooltipIcon icon={Car} tooltipText="Your cars" />Your cars</HeaderPanelLink
                    >
                    <HeaderPanelLink
                        href="/app/profile/user-profile"
                        on:click={() => (isOpen = false)}
                        ><TooltipIcon icon={User} tooltipText="Profile" />Profile</HeaderPanelLink
                    >
                    <HeaderPanelDivider />
                    <HeaderPanelLink on:click={logout}
                        ><TooltipIcon icon={User} tooltipText="Profile" />Log out</HeaderPanelLink
                    >
                </HeaderPanelLinks>
            </HeaderAction>
        {:else}
            <HeaderGlobalAction
                iconDescription="Explore open spots"
                tooltipAlignment="start"
                icon={Map}
                href="/#"
            />
            <HeaderGlobalAction
                iconDescription="Your spots"
                icon={MobileAdd}
                href="/app/your-spots"
            />
            <HeaderGlobalAction iconDescription="Your cars" icon={Car} href="/app/your-cars" />
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
        {/if}
    </HeaderUtilities>
</Header>

<style>
    .text-logo {
        max-height: 4rem;
    }
</style>
