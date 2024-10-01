<script lang="ts">
    import isLogIn from '../loginData';
    import { goto } from '$app/navigation';
    import { onMount } from 'svelte';

    onMount(() => {
        if (!$isLogIn) {
            goto('/auth/login');
        }
    });
    const logout = () => {
        $isLogIn = false;
        goto('/auth/login');
    };
</script>

<section class="loginForm">
    <h2>Welcome back!</h2>
    <input type="button" value="logout" on:click={logout} />
</section>

<style>
    h2 {
        font-size: 25px;
        font-weight: bold;
    }
</style>
