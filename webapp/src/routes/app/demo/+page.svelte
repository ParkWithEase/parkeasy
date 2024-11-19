<script lang="ts">
    import type { PageData } from './$types';

    export let data: PageData;
</script>

<div class="user-profile-container">
    <iframe
        width="1120"
        height="630"
        title=""
        src={data.message}
        frameborder="0"
        allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
        referrerpolicy="strict-origin-when-cross-origin"
        allowfullscreen
    ></iframe>
</div>

<style>
    .user-profile-container {
        display: flex;
        justify-items: flex-start;
        flex-direction: column;
    }
</style>
