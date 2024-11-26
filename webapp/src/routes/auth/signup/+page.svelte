<script lang="ts">
    import FormHeader from '$lib/components/form-header.svelte';
    import SubmitButton from '$lib/components/submit-button.svelte';
    import {
        ACCOUNT_CREATE_SUCCESS,
        DEFAULT_ACCOUNT_CREATION_ERROR,
        PASSWORD_NOT_MATCH
    } from '$lib/constants';
    import { newClient } from '$lib/utils/client';
    import { getErrorMessage } from '$lib/utils/error-handler';
    import { Form, TextInput, PasswordInput } from 'carbon-components-svelte';
    let firstName: string;
    let lastName: string;
    let email: string;

    let password: string;
    let passwordConfirm: string;

    let errorMessage: string;
    let successMessage: string;
    let accountCreated: boolean = false;
    let valid: boolean = false;

    let client = newClient();

    $: {
        if (password !== passwordConfirm) {
            errorMessage = PASSWORD_NOT_MATCH;
            valid = false;
        } else {
            errorMessage = '';
            valid = true;
        }
    }

    function signup() {
        if (valid) {
            client
                .POST('/user', {
                    body: {
                        full_name: `${firstName} ${lastName}`,
                        email: email,
                        password: password
                    }
                })
                .then(({ error }) => {
                    if (error) {
                        errorMessage = error.errors
                            ? (error.errors[0].message ?? getErrorMessage(error))
                            : DEFAULT_ACCOUNT_CREATION_ERROR;
                    } else {
                        successMessage = ACCOUNT_CREATE_SUCCESS;

                        accountCreated = true;
                        errorMessage = '';
                    }
                })
                .catch((err) => {
                    errorMessage = err;
                    successMessage = '';
                });
        }
    }
</script>

<div class="auth-form">
    <FormHeader headerText={'Signup for ParkEasy'}></FormHeader>
    <Form on:submit={signup}>
        {#if !accountCreated}
            <div class="name-section">
                <TextInput
                    required
                    labelText="First Name"
                    placeholder="Your first name"
                    bind:value={firstName}
                />
                <TextInput
                    required
                    labelText="Last Name"
                    placeholder="Your last name"
                    bind:value={lastName}
                />
            </div>
            <TextInput required labelText="Email" placeholder="Your email " bind:value={email} />
            <PasswordInput
                required
                type="password"
                labelText="Password"
                placeholder="Enter password..."
                bind:value={password}
            />
            <PasswordInput
                required
                type="password"
                labelText="Conform Password"
                placeholder="Enter password..."
                bind:value={passwordConfirm}
            />
        {/if}

        {#if errorMessage}
            <p style="color:red">{errorMessage}</p>
        {/if}

        {#if successMessage}
            <p aria-label="sucess message" style="color:green">{successMessage}</p>
        {/if}

        {#if !accountCreated}
            <SubmitButton buttonText={'Sign Up'}></SubmitButton>
        {/if}
    </Form>

    <a href="/auth/login">Back to Login?</a>
</div>

<style>
    .name-section {
        display: flex;
        flex-direction: row;
        justify-content: space-around;
        gap: 1rem;
    }
</style>
