<script lang="ts">
    import userClipart from '$lib/images/user-clipart.png';
    import Navbar from '$lib/components/navbar.svelte';
    import { goto } from '$app/navigation';
    import { BACKEND_SERVER } from '$lib/constants';
    import { onMount } from 'svelte';
    import { page } from '$app/stores';

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
        <a href="/profile/user-profile"><img src={userClipart} alt="user" class="logo-small" /></a>
        <a class="sidebar-link" href="/profile/user-profile" class:active={$page.url.pathname.includes('user-profile')}>User Profile</a>
        <a class="sidebar-link" href="/profile/booking-history" class:active={$page.url.pathname.includes('booking-history')}>Booking History</a>
        <a class="sidebar-link" href="/profile/listing-history" class:active={$page.url.pathname.includes('listing-history')}>Listing History</a>
        <a class="sidebar-link" href="/profile/cars" class:active={$page.url.pathname.includes('cars')}>Preference Cars</a>
        <a class="logout-link" href='/' on:click={logout}>Logout</a>
    </div>
    <div class="listing-container">
        <slot></slot>
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
        min-height: 100%;
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
        color:#ffffff;
        background-color: #090909;
        border: 1px solid #fcfcfc;
        transition: 0.5s;
    }

    .logout-link:hover {
        color: #f83636;
        background-color: rgb(0, 0, 0);
        transition: 0.3s;
    }

    .listing-container {
        position: relative;
        width: 70%;
        min-height: 100%;
        background-color: aliceblue;
        overflow-y: auto;
    }
</style>
