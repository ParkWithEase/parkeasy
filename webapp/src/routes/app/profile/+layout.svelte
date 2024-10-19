<script lang="ts">
    import userClipart from '$lib/images/user-clipart.png';
    import Navbar from '$lib/components/navbar.svelte';
    import { goto } from '$app/navigation';
    import { BACKEND_SERVER } from '$lib/constants';
    import { page, navigating } from '$app/stores';

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

<Navbar />

<div class="container">
    <div class="sidebar">
        <a href="/app/profile/user-profile"
            ><img src={userClipart} alt="user" class="logo-small" /></a
        >
        <a
            class="sidebar-link"
            href="/app/profile/user-profile"
            class:active={$page.url.pathname.includes('user-profile')}>User Profile</a
        >
        <a
            class="sidebar-link"
            href="/app/profile/booking-history"
            class:active={$page.url.pathname.includes('booking-history')}>Booking History</a
        >
        <a
            class="sidebar-link"
            href="/app/profile/leasing-history"
            class:active={$page.url.pathname.includes('leasing-history')}>Listing History</a
        >
        <a
            class="sidebar-link"
            href="/app/profile/preference-spots"
            class:active={$page.url.pathname.includes('preference-spots')}>Preference Cars</a
        >
        <a class="logout-link" href="/" on:click={logout}>Logout</a>
    </div>
    <div class="info-container">
        {#await $navigating?.complete}
            <div>Loading.</div>
        {:then}
            <slot></slot>
        {/await}
    </div>
</div>

<style>
    .container {
        display: flex;
        flex-direction: row;
        margin: 10rem 3rem;
        max-height: fit-content;
    }

    .container > div {
        margin-left: 1rem;
        border-radius: 20px;
        background-color: rgb(186, 214, 183);
    }

    .sidebar {
        display: flex;
        width: 30%;
        max-height: 70vh;
        flex-direction: column;
        font-size: 1.2rem;
        font-weight: bold;
        align-items: center;
        justify-content: center;
    }

    .sidebar > a {
        margin: 1rem;
        padding: 1rem 7rem;
        border-radius: 10px;
        text-decoration: none;
        color: #000000;
    }

    .sidebar-link:hover {
        color: #7ed957;
        background-color: rgb(0, 0, 0);
        transition: 0.3s;
    }

    .sidebar-link:active {
        color: #ffffff;
        background-color: rgb(0, 0, 0);
        transition: 0.3s;
    }

    a.active {
        color: #ffffff;
        background-color: #32683b;
        transition: 0.3s;
    }

    .logout-link:hover {
        color: #f83636;
        background-color: rgb(0, 0, 0);
        transition: 0.3s;
    }

    .info-container {
        position: relative;
        width: 70%;
        max-height: 70vh;
        background-color: aliceblue;
    }
</style>
