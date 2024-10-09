<script lang="ts">
    import { goto } from '$app/navigation';
    import { BACKEND_SERVER, INTERNAL_SERVER_ERROR } from '$lib/constants';
    import type { UserProfile } from '$lib/types/user/user';
    import { SkeletonText } from 'carbon-components-svelte';
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
        } catch {
            throw new Error(INTERNAL_SERVER_ERROR);
        }
    }

    let promise = getUser();
</script>

{#await promise}
    <SkeletonText paragraph lines={2} />
{:then user}
    {#if user != undefined}
        <div class="container">
            <div>
                <p>Name: {user.full_name}</p>
            </div>
            <div>
                <p>Email: {user.email}</p>
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

    p {
        font-size: 2rem;
    }
</style>
