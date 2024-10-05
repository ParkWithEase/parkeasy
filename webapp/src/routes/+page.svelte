<script lang="ts">
    import { goto } from '$app/navigation';
    import { BACKEND_SERVER } from '$lib/constants';

    export async function logout() {
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

<section class="loginForm">
    <h1>Welcome back!</h1>
    <input type="button" value="logout" on:click={logout} />
</section>

<style>
    h1 {
        font-size: 25px;
        font-weight: bold;
    }
</style>
