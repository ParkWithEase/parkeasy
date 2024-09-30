<script>
    import logo from '$lib/images/parkeasy-logo.png';
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

<section class="loginForm">
    <img src={logo} alt="logo" class="logo"/>
    {#if isLoggedIn}
    <h2>Welcome back!</h2>
    <input type="button" value="logout" on:click={logout}>
    {:else}
    <h2>Login</h2>
    <LoginPage />
    {/if}
</section>

<style>
    h2{
        font-size: 25px;
        font-weight: bold ;
    }

    .loginForm {
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
        margin: 6rem;
        padding: 2rem;
        background-color: #95b18f;
        border-radius: 2rem;
        border: 0.2rem solid #070707;
    }

    .logo{
        border-radius: 2rem;
        width: 20rem;
    }

    
</style>
