<script lang="ts">
    import { BACKEND_SERVER } from '$lib/constants';
    import logo from '$lib/images/parkeasy-logo.png';
    import { TextInput, PasswordInput, Form } from 'carbon-components-svelte';
    import { goto } from '$app/navigation';
    import SubmitButton from '$lib/components/submit-button.svelte';

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

<div class="auth-form">
    <img src={logo} alt="logo" class="logo-medium" />
    <Form on:submit={login}>
        <TextInput labelText="Email" placeholder="Enter email..." required bind:value={email} />
        <PasswordInput
            required
            type="password"
            labelText="Password"
            placeholder="Enter password..."
            bind:value={password}
        />
        <SubmitButton buttonText={'Login'} />
        {#if loginFail}
            <div>
                <p style="color:red">Wrong password or Email</p>
            </div>
        {/if}
    </Form>

    <p>Forgot your password? <a href="/auth/password-forgot">Go here</a></p>
    <p>Don't have an account? <a href="/auth/signup">Create one</a></p>
</div>

<style>
    p {
        font-weight: 600;
    }
</style>
