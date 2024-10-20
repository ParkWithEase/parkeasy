<script lang="ts">
    import Navbar from '$lib/components/navbar.svelte';
    
    
    import { navigating } from '$app/stores';
    import {
    SideNav,
    SideNavItems,
    SideNavLink
  } from "carbon-components-svelte";

    let isSideNavOpen = true;
</script>

<Navbar />

<SideNav bind:isOpen={isSideNavOpen}>
    <SideNavItems>
      <SideNavLink text="Your Profile" href="/app/profile/user-profile"/>
      <SideNavLink text="Booking History" href="/app/profile/booking-history"/>
      <SideNavLink text="Leasing History" href="/app/profile/leasing-history"/>
      <SideNavLink text="Preferred Spots" href="/app/profile/preferred-spots"/>
    </SideNavItems>
</SideNav>

<div class="info-container">
    {#await $navigating?.complete}
        <div>Loading.</div>
    {:then}
        <slot></slot>
    {/await}
</div>


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
        margin-left: 16rem;
        margin-top: 3rem;
        padding: 1rem;
    }
</style>
