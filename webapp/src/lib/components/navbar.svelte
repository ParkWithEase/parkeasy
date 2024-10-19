<script lang="ts">
    import logo from '$lib/images/logo-text.png';
    import {
    Header,
    HeaderUtilities,
    HeaderGlobalAction,
    SkipToContent
    } from "carbon-components-svelte";
    import Logout from "../../../node_modules/carbon-icons-svelte/lib/Logout.svelte";
    import User from "../../../node_modules/carbon-icons-svelte/lib/User.svelte";
    import Map from "../../../node_modules/carbon-icons-svelte/lib/Map.svelte";
    import MobileAdd from "carbon-icons-svelte/lib/MobileAdd.svelte";
</script>

<svelte:window />

<Header>
    <a href="/#"><img src={logo} alt="parkeasy-logo" class="text-logo" /></a>
    <svelte:fragment slot="skip-to-content">
      <SkipToContent />
    </svelte:fragment>
    <HeaderUtilities>
      <HeaderGlobalAction
        iconDescription="Explore open spots"
        tooltipAlignment="start"
        icon={Map}
        href="/#"
      />
      <HeaderGlobalAction iconDescription="Add your spot" icon={MobileAdd} href="/app/profile/user-profile"/>
      <HeaderGlobalAction iconDescription="Profile" icon={User} href="/app/profile/user-profile"/>
      <HeaderGlobalAction
        iconDescription="Log out"
        tooltipAlignment="end"
        icon={Logout}
      />
    </HeaderUtilities>
</Header>

<style>

    .text-logo {
        max-height: 4rem;
    }
</style>
