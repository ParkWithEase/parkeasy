<script lang="ts">
    import { page } from '$app/stores';
    import { SideNav, SideNavDivider, SideNavItems, SideNavLink } from 'carbon-components-svelte';
    import { ChartLineSmooth, Information, Money } from 'carbon-icons-svelte';
    let currentSpotID: string = $page.params['id'];
    let spot_info_link = `/app/your-spots/${currentSpotID}/spot-info`;
    let spot_performance_link = `/app/your-spots/${currentSpotID}/spot-performance`;
    let leasing_history_link = `/app/your-spots/${currentSpotID}/spot-leasing-history`;
    console.log($page.url.pathname);
</script>

<SideNav isOpen fixed>
    <SideNavItems>
        <SideNavLink
            text="General Information"
            href={spot_info_link}
            isSelected={spot_info_link == $page.url.pathname}
            icon={Information}
        />
        <SideNavLink
            text="Performance "
            href={spot_performance_link}
            isSelected={spot_performance_link == $page.url.pathname}
            icon={ChartLineSmooth}
        />
        <SideNavLink
            text="Leasing History"
            href={leasing_history_link}
            isSelected={leasing_history_link == $page.url.pathname}
            icon={Money}
        />
        <SideNavDivider />
    </SideNavItems>
</SideNav>
<slot></slot>
