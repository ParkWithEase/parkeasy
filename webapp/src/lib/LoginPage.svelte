<script>
    import user from "../user";
    import {
        FluidForm,
        TextInput,
        PasswordInput,
        Button,
    } from "carbon-components-svelte";
    let username = "";
    let password = "";
    let show = false;

    let currentError = null;

    const login = () => {
        fetch("http://localhost:3030/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ username: username, password: password }),
        })
            .then((res) => {
                if (res.status < 299) return res.json();
                if (res.status > 299)
                    currentError = "Something not right with server response";
            })
            .then((data) => {
                if (data) user.update((val) => (val = { ...data }));
            })
            .catch((error) => {
                currentError = error;
                console.log("Error logging in: ", error);
            });
    };
</script>

<body class="loginForm">
    <form>
        <TextInput
            labelText="User name"
            placeholder="Enter user name..."
            required
            class="input-field"
        />
        <PasswordInput
            required
            type="password"
            labelText="Password"
            placeholder="Enter password..."
            class="input-field"
        />
        <Button type="submit">Submit</Button>
    </form>
</body>

<style>
    .loginForm {
        display: flex;
        flex-direction: column;
        align-items: center;
        width: 100%;
        height: 20vh;
    }

    :global(bx--btn.bx--btn--primary) {
        color: rgb(0, 248, 248);
        background-color: black;
        margin: 10px;
        padding: 5px;
        border: 1px solid white;
        border-radius: 5px;
    }

    :global(bx--form--fluid.fluidForm.bx--form) {
        border: 2px solid black;
    }

    .input-field {
        font-weight: 500;
        font-size: larger;
    }
</style>
