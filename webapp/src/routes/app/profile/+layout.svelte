<script lang="ts">
    import Navbar from '$lib/components/navbar.svelte';

    import { navigating, page } from '$app/stores';
    import { SideNav, SideNavItems, SideNavLink, Content } from 'carbon-components-svelte';

    let user_profile_link = '/app/profile/user-profile';
    let booking_history_link = '/app/profile/booking-history';
    let leasing_history_link = '/app/profile/leasing-history';
    let preferred_spots_link = '/app/profile/preferred-spots';

    let is_side_nav_open: boolean = true;
</script>

<Navbar />

<SideNav bind:isOpen={is_side_nav_open}>
    <SideNavItems>
        <SideNavLink
            text="Your Profile"
            href={user_profile_link}
            isSelected={user_profile_link == $page.url.pathname}
        />
        <SideNavLink
            text="Booking History"
            href={booking_history_link}
            isSelected={booking_history_link == $page.url.pathname}
        />
        <SideNavLink
            text="Leasing History"
            href={leasing_history_link}
            isSelected={leasing_history_link == $page.url.pathname}
        />
        <SideNavLink
            text="Preferred Spots"
            href={preferred_spots_link}
            isSelected={preferred_spots_link == $page.url.pathname}
        />
    </SideNavItems>
</SideNav>

<Content>
    <div class="info-container">
        {#await $navigating?.complete}
            <div>Loading.</div>
        {:then}
            <slot></slot>
        {/await}
    </div>
</Content>

<style>
    :global(.bx--side-nav__items) {
        height: 50%;
        border-right: 1px black solid;
        font-size: 1rem;
    }

    :global(.bx--side-nav__link) {
        border-radius: 10px;
        text-decoration: none;
        color: #000000;
    }

    :global(.bx--side-nav__link):hover {
        color: #7ed957;
        background-color: rgb(0, 0, 0);
        transition: 0.3s;
    }

    :global(.bx--side-nav__link):active {
        color: #ffffff;
        background-color: rgb(0, 0, 0);
        transition: 0.3s;
    }

    .info-container {
        position: relative;
        width: 100%;
        height: 100vh;
    }
</style>
