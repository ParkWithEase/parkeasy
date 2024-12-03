<script lang="ts">
    import FormHeader from '$lib/components/form-header.svelte';
    import SubmitButton from '$lib/components/submit-button.svelte';
    import { INTERNAL_SERVER_ERROR, PASSWORD_RESET_TOKEN_GET_ERROR } from '$lib/constants';
    import { newClient } from '$lib/utils/client';
    import { Form, TextInput } from 'carbon-components-svelte';

    let email: string;
    let passwordToken: string;
    let tokenProduced: boolean = false;
    let errorMessage: string;

    let client = newClient();

    function getToken() {
        client
            .POST('/auth/password:forgot', { body: { email: email } })
            .then(({ data, error }) => {
                if (error || data.password_token == '') {
                    errorMessage = PASSWORD_RESET_TOKEN_GET_ERROR;
                } else {
                    passwordToken = data.password_token;
                    tokenProduced = true;
                    errorMessage = '';
                }
            })
            .catch(() => {
                errorMessage = INTERNAL_SERVER_ERROR;
            });
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
