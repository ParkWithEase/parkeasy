<script lang="ts">
    import { navigating, page } from '$app/stores';
    import { SideNav, SideNavItems, SideNavLink, Content } from 'carbon-components-svelte';
    import { responsiveState } from '$lib/stores/responsive';
    import { Currency, LocationHeart, Purchase, UserAvatarFilledAlt } from 'carbon-icons-svelte';

    let user_profile_link = '/app/profile/user-profile';
    let booking_history_link = '/app/profile/booking-history';
    let leasing_history_link = '/app/profile/leasing-history';
    let preferred_spots_link = '/app/profile/preferred-spots';

    let shouldReponse: boolean = false;
    responsiveState.subscribe((value) => {
        shouldReponse = value;
    });
</script>

<SideNav isOpen={!shouldReponse} fixed={!shouldReponse} rail={shouldReponse}>
    <SideNavItems>
        <SideNavLink
            text="Your Profile"
            icon={UserAvatarFilledAlt}
            href={user_profile_link}
            isSelected={user_profile_link == $page.url.pathname}
        />
        <SideNavLink
            text="Booking History"
            icon={Purchase}
            href={booking_history_link}
            isSelected={booking_history_link == $page.url.pathname}
        />
        <SideNavLink
            text="Leasing History"
            icon={Currency}
            href={leasing_history_link}
            isSelected={leasing_history_link == $page.url.pathname}
        />
        <SideNavLink
            text="Preferred Spots"
            icon={LocationHeart}
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
