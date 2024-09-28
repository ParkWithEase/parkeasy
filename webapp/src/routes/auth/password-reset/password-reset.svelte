<script lang="ts">
    import { BACKEND_SERVER } from '$lib/constants';

    export let resetToken: string | null;

    let errorMessage: string;
    let successMessage: string;

    let password: string = '';
    let passwordConfirm: string = '';

    let valid: boolean = false;
    $: {
        if (resetToken == null) {
            valid = false;
        }

        if (password !== passwordConfirm && password.length > 0) {
            errorMessage = "password doesn't match";
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
            console.log(JSON.stringify(payload));
            try {
                const response = await fetch(`${BACKEND_SERVER}/auth/password:reset`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(payload)
                });

                const data = await response.json();
                console.log(data);
                if (response.ok) {
                    successMessage = 'Password changed success';
                } else {
                    errorMessage = data.message || 'Failure';
                }
            } catch {
                errorMessage = 'Something went wrong';
                successMessage = '';
            }
        }
    }
</script>

<html lang="en">
    <h1>Reset your password</h1>
    <body>
        <div class="form-container">
            <form on:submit={resetPassword}>
                <div>
                    <label for="new-password"> New password</label>
                    <input type="password" id="new-password" bind:value={password} />
                </div>

                <div>
                    <label for="confirm-password">Confirm password</label>
                    <input type="password" id="confirm-password" bind:value={passwordConfirm} />
                </div>

                <div>
                    <input type="submit" />
                </div>

                <div>
                    {#if successMessage}
                        <p style="color:green">{successMessage}</p>
                    {/if}

                    {#if errorMessage}
                        <p style="color:red">{errorMessage}</p>
                    {/if}
                </div>
            </form>
        </div>
    </body>
</html>

<style>
    body {
        font-family: Arial, sans-serif;
        background-color: #f0f8f0;
        color: #333;
        display: flex;
        justify-content: center;
        align-items: flex-start;
        height: 100vh;
        margin: 0;
    }

    .form-container {
        background-color: #d9f8d9;
        display: flex;
        justify-content: center;
        align-items: center;
        border-radius: 10px;
        padding: 20px;
        box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
        width: auto;
        margin-top: 20px;
    }

    h1 {
        text-align: center;
        color: #2c5f2d;
    }

    input[type='password'] {
        width: 90%;
        padding: 10px;
        margin: 10px 0;
        border: 1px solid #2c5f2d;
        border-radius: 5px;
    }

    input[type='submit'] {
        background-color: #2c5f2d;
        color: white;
        border: none;
        padding: 10px;
        border-radius: 5px;
        align-self: center;
        cursor: pointer;
        width: 50%;
    }

    input[type='submit']:hover {
        background-color: #245424;
    }
</style>
