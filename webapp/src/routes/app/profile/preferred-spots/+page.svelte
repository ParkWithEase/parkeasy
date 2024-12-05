<script lang="ts">
    import SpotDisplay from '$lib/components/spot-component/spot-display.svelte';
    import IntersectionObserver from 'svelte-intersection-observer';
    import type { PageData } from './$types';

    export let data: PageData;

    let intersecting: boolean;
    let loadTrigger: HTMLElement | null = null;
    let canLoadMore = data.hasNext;
    let loadLock = false;

    $: while (intersecting && canLoadMore && !loadLock) {
        loadLock = true;
        data.paging
            .next()
            .then(({ value: { data: spots }, done }) => {
                if (spots) {
                    data.spots = [...data.spots, ...spots];
                }
                canLoadMore = !done;
            })
            .finally(() => {
                loadLock = false;
            });
    }
</script>

<div>
    {#key data.spots}
        {#each data.spots as spot}
            <a href={`/app/booking/${spot.id}/spot-booking`} style="text-decoration: none;">
                <SpotDisplay {spot} />
            </a>
        {/each}
    {/key}
</div>
<IntersectionObserver element={loadTrigger} bind:intersecting>
    {#if canLoadMore}
        <div bind:this={loadTrigger}>Loading...</div>
    {/if}
</IntersectionObserver>
