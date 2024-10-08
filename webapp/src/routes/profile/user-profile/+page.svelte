<script lang="ts">
    import { goto } from '$app/navigation';
    import { BACKEND_SERVER } from '$lib/constants';

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

<div class="container">
    <h1>This is user profile</h1>
</div>

<style> 
    .container{
        display: flex;
        flex-direction: row;
    }


</style> 