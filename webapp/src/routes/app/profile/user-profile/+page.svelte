<script lang="ts">
    import type { PageData } from './$types';
    export let data: PageData;
</script>

<div class="container">
    <div>
        <p>Name: {data.full_name}</p>
    </div>
    <div>
        <p>Email: {data.email}</p>
    </div>
</div>

<style>
    .container {
        display: flex;
        justify-items: flex-start;
        flex-direction: column;
    }

    p {
        padding-top: 1rem;
        margin-left: 1rem;
        font-size: 1.2rem;
    }
</style>
