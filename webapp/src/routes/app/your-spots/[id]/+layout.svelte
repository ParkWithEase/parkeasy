<script lang="ts">
    import { page } from '$app/stores';
    import { SideNav, SideNavDivider, SideNavItems, SideNavLink } from 'carbon-components-svelte';
    import { ChartLineSmooth, Information, Money } from 'carbon-icons-svelte';
    let currentSpotID: string = $page.params['id'];
    let spotInfoLink = `/app/your-spots/${currentSpotID}/spot-info`;
    let spotPerformanceLink = `/app/your-spots/${currentSpotID}/spot-performance`;
    let leasingHistoryLink = `/app/your-spots/${currentSpotID}/spot-leasing-history`;
</script>

<SideNav isOpen fixed>
    <SideNavItems>
        <SideNavLink
            text="General Information"
            href={spotInfoLink}
            isSelected={spotInfoLink == $page.url.pathname}
            icon={Information}
        />
        <SideNavLink
            text="Performance "
            href={spotPerformanceLink}
            isSelected={spotPerformanceLink == $page.url.pathname}
            icon={ChartLineSmooth}
        />
        <SideNavLink
            text="Leasing History"
            href={leasingHistoryLink}
            isSelected={leasingHistoryLink == $page.url.pathname}
            icon={Money}
        />
        <SideNavDivider />
    </SideNavItems>
</SideNav>
<slot></slot>
