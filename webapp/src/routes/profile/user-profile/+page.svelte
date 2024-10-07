<script lang="ts">
    import { goto } from '$app/navigation';
    import { BACKEND_SERVER } from '$lib/constants';
    import userClipart from '$lib/images/user-clipart.png';
    import Navbar from '$lib/components/navbar.svelte';

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
        <a href="/profile/user-profile"
            ><img src={userClipart} alt="user" class="logo-extra-small"
        /></a>
        <a class="sidebar-link" href="/profile/user-profile">User Profile</a>
        <a class="sidebar-link" href="/profile/user-profile">Booking History</a>
        <a class="sidebar-link" href="/">Listing History</a>
        <a class="sidebar-link" href="/">Preference Cars</a>
        <a class="logout-link" href='/' on:click={logout}>Logout</a>
    </div>
    <div class="listing-container">

    </div>
</div>

<style> 
    .container{
        display: flex;
        flex-direction: row;
        margin: 10rem 3rem;
        min-height: 70vh;
    }

    .container > div{
        margin-left: 1rem;
        border-radius: 20px;
        background-color: rgb(255, 255, 255);
    }

    .sidebar{
        display: flex;
        width: 25%;
        flex-direction: column;
        font-size: 1.2rem;
        font-weight: bold;
        align-items: center;
        justify-content: center;
    }

    .sidebar > a {
        margin: 1rem;
        padding: 1rem 5rem;
        border-radius: 10px;
        text-decoration: none;
        color: #000000;
    }

    .sidebar-link:hover{
        color: #7ed957;
        background-color: rgb(0, 0, 0);
        transition: 0.3s;
    }

    .sidebar-link:active{
        color: #7ed957;
        background-color: rgb(0, 0, 0);
        transition: 0.3s;
    }

    .logout-link:hover{
        color: #f83636;
        background-color: rgb(0, 0, 0);
        transition: 0.3s;
    }

    .listing-container{
        width: 75%;
        background-color: aliceblue;
    }

</style> 