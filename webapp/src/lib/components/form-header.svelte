<script>
    import logo from '$lib/images/parkeasy-logo.png';
    export let headerText;
</script>

<div class="form-header">
    <img src={logo} alt="logo" class="logo-small" />
    <h1 class="header-text">{headerText}</h1>
</div>

<style>
    .form-header {
        display: flex;
        direction: row;
        padding: 1rem;
        gap: 1rem;
    }

    .header-text {
        font-weight: bolder;
        max-width: 20rem;
        align-self: center;
        text-align: left;
    }
</style>
