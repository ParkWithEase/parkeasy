<script lang="ts">
    import { navigating, page } from '$app/stores';
    import { SideNav, SideNavItems, SideNavLink, Content } from 'carbon-components-svelte';

    let user_profile_link = '/app/profile/user-profile';
    let booking_history_link = '/app/profile/booking-history';
    let leasing_history_link = '/app/profile/leasing-history';
    let preferred_spots_link = '/app/profile/preferred-spots';

    let is_side_nav_open: boolean = true;
</script>

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
    {#await $navigating?.complete}
        <div>Loading.</div>
    {:then}
        <slot></slot>
    {/await}
</Content>

<style>
</style>
