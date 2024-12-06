<script lang="ts">
    import logo from '$lib/images/logo_stacked_outlined.svg';
    import { TextInput, PasswordInput, Form } from 'carbon-components-svelte';
    import { goto } from '$app/navigation';
    import SubmitButton from '$lib/components/submit-button.svelte';
    import { newClient } from '$lib/utils/client';

    let email: string = '';
    let password: string = '';

    let loginFail: boolean = false;

    let client = newClient();

    function login(event: Event) {
        event.preventDefault();
        client
            .POST('/auth', { body: { email: email, password: password, persist: false } })
            .then(({ error }) => {
                if (error) {
                    loginFail = true;
                } else {
                    goto('/app/');
                }
            })
            .catch((err) => {
                console.log(err);
            });
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
