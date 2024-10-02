<script lang="ts">
    import { BACKEND_SERVER } from '$lib/constants';
    import logo from '$lib/images/parkeasy-logo.png';
    import { TextInput, PasswordInput, Button, Form } from 'carbon-components-svelte';
    import { goto } from '$app/navigation';

    let email: string = '';
    let password: string = '';

    let loginFail: boolean = false;

    async function login(event: Event) {
        event.preventDefault();
        try {
            const response = await fetch(`${BACKEND_SERVER}/auth`, {
                method: 'POST',
                credentials: 'include',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ email: email, password: password })
            });
            if (response.ok) {
                goto('/');
            } else {
                loginFail = true;
            }
        } catch (err) {
            console.log(err);
        }
    }
</script>

<div class="loginSection">
    <div class="loginForm">
        <img src={logo} alt="logo" class="logo" />
        <Form>
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
        </Form>
        {#if loginFail}
            <div>
                <p style="color:red">Wrong password or Email</p>
            </div>
        {/if}
        <Button type="submit" on:click={login}>Login</Button>
    </div>
</div>

<style>
    .loginForm {
        display: flex;
        flex-direction: column;
        align-items: self-start;
        width: auto;
        height: auto;
        padding: 2rem;
        background-color: #c3e9bb;
        border-radius: 2rem;
        border: 0.2rem solid #070707;
    }

    .loginSection {
        position: absolute;
        top: 50%;
        left: 50%;
        margin-right: -50%;
        transform: translate(-50%, -50%);
    }

    .logo {
        border-radius: 2rem;
        max-width: 15rem;
    }
</style>
