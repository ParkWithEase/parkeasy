<script lang="ts">
    import type { PageData } from './$types';
    export let data: PageData;
</script>

<header class="header">
    <h1>Your Profile</h1>
</header>
<div class="container">
    <div>
        <p>Name: {data.full_name}</p>
    </div>
    <div>
        <p>Email: {data.email}</p>
    </div>
</div>

<style>
    h1 {
        font-weight: bold;
    }

    .container {
        display: flex;
        justify-items: flex-start;
        flex-direction: column;
        border-top: 2px solid black;
    }

    p {
        padding-top: 1rem;
        margin-left: 1rem;
        font-size: 1.2rem;
    }
</style>
