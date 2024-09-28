<script lang="ts">
    import type { UserProfile } from '$lib/types/user/user';
    import { BACKEND_SERVER } from '$lib/constants';
    import { onMount } from 'svelte';

    let user: UserProfile;
    onMount(() => {
        async function getUser() {
            try {
                user = await fetch(`${BACKEND_SERVER}/user`, {
                    credentials: 'include'
                }).then((x) => x.json());
                console.log(user);
            } catch (err: unknown) {
                if (typeof err === 'string') {
                    console.log(err.toUpperCase());
                } else if (err instanceof Error) {
                    console.log(err.message);
                }
            }
        }

        getUser();
    });
</script>

<html lang="en">
    {#if user === undefined}
        <h2>Loading User Profile...</h2>
    {:else}
        <body>
            <div class="profile-container">
                <h1>User Profile</h1>
                <label for="name">Name:</label>
                <p>{user.full_name}</p>

                <label for="email">Email:</label>
                <p>{user.email}</p>
            </div>
        </body>
    {/if}
</html>

<style>
    body {
        font-family: Arial, sans-serif;
        background-color: #f0f8f0;
        color: #333;
        display: flex;
        justify-content: center;
        align-items: flex-start;
        height: 100vh;
        margin: 0;
        padding-top: 50px;
    }
    .profile-container {
        background-color: #d9f8d9;
        border-radius: 10px;
        padding: 20px;
        box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
        width: 300px;
    }
    h1 {
        text-align: center;
        color: #2c5f2d;
    }
    label {
        display: block;
        margin: 10px 0 5px;
        font-weight: bold;
    }
</style>
