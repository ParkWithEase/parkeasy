<script lang="ts">
    import { goto } from '$app/navigation';
    import { BACKEND_SERVER, INTERNAL_SERVER_ERROR } from '$lib/constants';
    import type { UserProfile } from '$lib/types/user/user';
    import { onMount } from 'svelte';

    async function getUser(): Promise<UserProfile | null> {
        try {
            const response = await fetch(`${BACKEND_SERVER}/user`, {
                credentials: 'include'
            });

            if (response.ok) {
                console.log(response);
                let user: UserProfile = (await response.json()) as UserProfile;
                return user;
            } else {
                goto('/auth/login');
                return null;
            }
        } catch (err: unknown) {
            throw new Error(INTERNAL_SERVER_ERROR);
        }
    }

    let promise = getUser();

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

{#await promise}
    <p>waiting...</p>
{:then user}
    {#if user != undefined}
        <div class="container">
            <div>
                <h1>Name: {user.full_name}</h1>
            </div>
            <div>
                <h1>Email: {user.email}</h1>
            </div>
        </div>
    {:else}
        <div>
            {INTERNAL_SERVER_ERROR}
        </div>
    {/if}
{/await}

<style>
    .container {
        display: flex;
        justify-items: flex-start;
        flex-direction: column;
    }

    h1 {
        text-align: left;
    }
</style>
