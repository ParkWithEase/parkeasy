<script lang="ts">
    import type { UserProfile } from '$lib/types/user/user';
    import { BACKEND_SERVER } from '$lib/constants';
    import { onMount } from 'svelte';
    import { goto } from '$app/navigation';

    let user: UserProfile;
    onMount(() => {
        async function getUser() {
            try {
                const response = await fetch(`${BACKEND_SERVER}/user`, {
                    credentials: 'include'
                });

                if (response.ok) {
                    user = await response.json();
                } else {
                    goto('/auth/login');
                }
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

{#if user === undefined}
    <h2>Loading User Profile...</h2>
{:else}
    <div class="profile-container">
        <h1>User Profile</h1>
        <label for="name">Name:</label>
        <p>{user.full_name}</p>

        <label for="email">Email:</label>
        <p>{user.email}</p>
    </div>
{/if}

<style>
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
