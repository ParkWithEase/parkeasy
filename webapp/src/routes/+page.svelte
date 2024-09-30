<script>
    import transparentLogo from '$lib/images/parkeasy-logo.png';
    import user from '../user';
    import LoginPage from '$lib/LoginPage.svelte';

    $: isLoggedIn = $user === null? false: true;

    const logout = () =>{
        user.update(val => val = null);
    }
</script>

<svelte:head>
    <title>parkeasy</title>
</svelte:head>

<section>
    <img src={transparentLogo} alt="Welcome" />
    {#if isLoggedIn}
    <h2>Welcome back!</h2>
    <input type="button" value="logout" on:click={logout}>
    {:else}
    <h2>Login</h2>
    <LoginPage />
    {/if}
</section>

<style>
    section {
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
        flex: 0.6;
    }
</style>
