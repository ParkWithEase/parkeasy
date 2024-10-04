<script lang="ts">
    import FormHeader from '$lib/components/form-header.svelte';
    import SubmitButton from '$lib/components/submit-button.svelte';
    import { BACKEND_SERVER } from '$lib/constants';
    import { Form, TextInput } from 'carbon-components-svelte';

    let email: string;
    let passwordToken: string;
    let tokenProduced: boolean = false;
    let errorMessage: string;

    async function getToken() {
        try {
            const response = await fetch(`${BACKEND_SERVER}/auth/password:forgot`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ email: email })
            });

            if (response.ok) {
                const data = await response.json();
                passwordToken = data.password_token;
                if (passwordToken !== '') {
                    tokenProduced = true;
                    errorMessage = '';
                } else {
                    errorMessage =
                        "We shouldn't be doing this but for demo sake, your email doesn't exist";
                }
            } else {
                const data = await response.json();
                errorMessage = data.errors[0].message;
            }
        } catch (err) {
            console.log(err);
            errorMessage = 'Something went wrong in the server';
        }
    }
</script>

<div class="auth-form">
    <FormHeader headerText={'Forgot password?'}></FormHeader>

    <Form on:submit={getToken}>
        {#if !tokenProduced}
            <TextInput
                style={'min-width:20rem'}
                labelText="Email"
                placeholder="Enter your email..."
                bind:value={email}
                required
            />
            {#if errorMessage}
                <p style="color:red">{errorMessage}</p>
            {/if}
            <SubmitButton buttonText={'Submit'} />
        {/if}
    </Form>

    {#if passwordToken}
        <p>Here is the password reset link.</p>
        <a
            aria-label="password reset link"
            href={`/auth/password-reset?password_reset_token=${passwordToken}`}
            >Click here for the password reset link</a
        >
    {/if}
    <p><a href="/auth/login"> Back to login? </a></p>
</div>
