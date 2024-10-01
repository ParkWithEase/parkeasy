<script lang="ts">
    import { BACKEND_SERVER } from '$lib/constants';
    import logo from '$lib/images/parkeasy-logo.png';
    import { TextInput, PasswordInput, Button } from 'carbon-components-svelte';
    import { goto } from '$app/navigation';
    import isLogIn from '../../../loginData';

    let email: string = '';
    let password: string = '';

    let currentError = null;

    const login = () => {
        fetch(`${BACKEND_SERVER}/auth`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ email: email, password: password })
        })
            .then((res) => {
                if (res.status < 299) $isLogIn = true;
                goto('/');
                if (res.status > 299) currentError = 'Something not right with server response';
            })
            .catch((error) => {
                currentError = error;
                console.log('Error logging in: ', currentError);
            });
    };
</script>

<html lang="en">
    <section class="loginSection">
        <body class="loginForm">
            <img src={logo} alt="logo" class="logo" />
            <form on:submit={login}>
                <TextInput
                    class="input-field"
                    labelText="Email"
                    placeholder="Enter email..."
                    required
                    bind:value={email}
                />
                <PasswordInput
                    required
                    type="password"
                    labelText="Password"
                    placeholder="Enter password..."
                    class="input-field"
                    bind:value={password}
                />

                <Button type="submit">Submit</Button>
            </form>
        </body>
    </section>
</html>

<style>
    .loginForm {
        display: flex;
        flex-direction: column;
        align-items: center;
        width: 100%;
        height: 20vh;
    }

    .loginSection {
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
        margin: 5rem 2rem 5rem;
        height: 100%;
        padding: 2rem;
        background-color: #c3e9bb;
        border-radius: 2rem;
        border: 0.2rem solid #070707;
    }

    .logo {
        border-radius: 2rem;
        width: 20rem;
    }
</style>
