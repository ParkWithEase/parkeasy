<script lang="ts">
    import type { PageData } from './$types';
    export let data: PageData;
</script>

<div class="user-profile-container">
    <div>
        <p>Name: {data.full_name}</p>
    </div>
    <div>
        <p>Email: {data.email}</p>
    </div>
</div>

<style>
    .user-profile-container {
        display: flex;
        justify-items: flex-start;
        flex-direction: column;
    }

    p {
        font-size: 1.2rem;
        padding-top: 1rem;
        margin-left: 1rem;
    }
</style>
