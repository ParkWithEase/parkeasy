<script>
    import transparentLogo from '$lib/images/parkeasy-logo.png';
</script>

<svelte:head>
    <title>parkeasy</title>
</svelte:head>

<section>
    <img src={transparentLogo} alt="Welcome" />
</section>

<style>
    section {
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
        flex: 0.6;
    }
</style>
