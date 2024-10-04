<script lang="ts">
    import FormHeader from '$lib/components/form-header.svelte';
    import { BACKEND_SERVER, INTERNAL_SERVER_ERROR, PASSWORD_NOT_MATCH } from '$lib/constants';
    import { Form, PasswordInput } from 'carbon-components-svelte';

    export let resetToken: string | null;

    let errorMessage: string;
    let successMessage: string;

    let password: string = '';
    let passwordConfirm: string = '';

    let valid: boolean = false;
    let tokenUsed: boolean = false;
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

    async function resetPassword(event: Event) {
        // Default form behaviour will send the payload on to the url
        event.preventDefault();
        if (valid) {
            errorMessage = '';
            console.log(resetToken);
            let payload = {
                password_token: resetToken,
                new_password: password
            };
            try {
                const response = await fetch(`${BACKEND_SERVER}/auth/password:reset`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(payload)
                });

                if (response.ok) {
                    successMessage = 'Password changed success';
                    tokenUsed = true;
                } else {
                    const data = await response.json();
                    console.log(data);
                    errorMessage =
                        data.errors[0].message || "Hmm seems like your reset token doesn't work";
                    successMessage = '';
                }
            } catch {
                errorMessage = INTERNAL_SERVER_ERROR;
                successMessage = '';
            }
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
