<script lang="ts">
    import FormHeader from '$lib/components/form-header.svelte';
    import {
        INTERNAL_SERVER_ERROR,
        PASSWORD_NOT_MATCH,
        PASSWORD_RESET_SUCCESS_MESSAGE
    } from '$lib/constants';
    import { newClient } from '$lib/utils/client';
    import { getErrorMessage } from '$lib/utils/error-handler';
    import { Form, PasswordInput } from 'carbon-components-svelte';

    export let resetToken: string | null;

    let errorMessage: string;
    let successMessage: string;

    let password: string = '';
    let passwordConfirm: string = '';

    let valid: boolean = false;
    let tokenUsed: boolean = false;

    let client = newClient();

    $: {
        if (resetToken == null) {
            valid = false;
        }

        if (password !== passwordConfirm && password.length > 0) {
            errorMessage = PASSWORD_NOT_MATCH;
            valid = false;
        } else {
            errorMessage = '';
            valid = true;
        }
    }

    function resetPassword(event: Event) {
        event.preventDefault();
        if (valid) {
            errorMessage = '';
            client
                .POST('/auth/password:reset', {
                    body: { new_password: password, password_token: resetToken ?? '' }
                })
                .then(({ error }) => {
                    if (error) {
                        errorMessage = getErrorMessage(error);
                        successMessage = '';
                    } else {
                        successMessage = PASSWORD_RESET_SUCCESS_MESSAGE;
                        tokenUsed = true;
                    }
                })
                .catch(() => {
                    errorMessage = INTERNAL_SERVER_ERROR;
                    successMessage = '';
                });
        }
    }
</script>

<div>
    <div class="auth-form" on:submit={resetPassword}>
        <FormHeader headerText={'Reset Your password'} />

        <Form>
            {#if !tokenUsed}
                <PasswordInput
                    required
                    class="input-field"
                    type="password"
                    labelText="New Password"
                    bind:value={password}
                />
                <PasswordInput
                    required
                    class="input-field"
                    type="password"
                    labelText="Confirm Password"
                    placeholder="confirm your password..."
                    bind:value={passwordConfirm}
                />
            {/if}

            {#if successMessage}
                <p style="color:green">{successMessage}</p>
            {:else if errorMessage}
                <p style="color:red">{errorMessage}</p>
            {/if}

            {#if !tokenUsed}
                <div class="submission">
                    <button class="submit-button" type="submit">Submit</button>
                    <span style="align-self:center;">
                        <a href="/auth/login">Back to login?</a>
                    </span>
                </div>
            {:else}
                <span style="align-self:center;"> <a href="/auth/login">Back to login?</a> </span>
            {/if}
        </Form>
    </div>
</div>

<style>
    .submission {
        display: flex;
        width: auto;
        justify-content: space-between;
        flex-direction: row;
        justify-items: center;
    }
</style>
