<script lang="ts">
    import { page } from '$app/stores';
    import { SideNav, SideNavDivider, SideNavItems, SideNavLink } from 'carbon-components-svelte';
    import { Information, Purchase } from 'carbon-icons-svelte';
    let currentSpotID: string = $page.params['id'];
    let spotBookingLink = `/app/booking/${currentSpotID}/spot-booking`;
    let spotBookingLinkHistory = `/app/booking/${currentSpotID}/booking-history`;
    import { responsiveState } from '$lib/stores/responsive';
    let shouldReponse: boolean = false;
    responsiveState.subscribe((value) => {
        shouldReponse = value;
    });
</script>

<SideNav isOpen={!shouldReponse} fixed={!shouldReponse} rail={shouldReponse}>
    <SideNavItems>
        <SideNavLink
            text="Make a booking"
            href={spotBookingLink}
            isSelected={spotBookingLink == $page.url.pathname}
            icon={Information}
        />
        <SideNavLink
            text="Booking History "
            href={spotBookingLinkHistory}
            isSelected={spotBookingLinkHistory == $page.url.pathname}
            icon={Purchase}
        />
        <SideNavDivider />
    </SideNavItems>
</SideNav>
<slot></slot>
