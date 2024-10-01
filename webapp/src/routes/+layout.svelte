<script>
    import "../app.css";
    import {
        Footer,
        FooterBrand,
        FooterCopyright,
        FooterIcon,
        FooterLink,
        FooterLinkGroup,
    } from "flowbite-svelte";
</script>

<body class="app">
    <main class="mainLoginPage">
        <slot></slot>
    </main>
</body>

<Footer>
    <div class="footer">
        <FooterCopyright href="/" by="parkeasy™" year={2024} />
        <FooterLinkGroup
            ulClass="flex flex-wrap items-center mt-3 text-sm text-gray-500 dark:text-gray-400 sm:mt-0"
        ></FooterLinkGroup>
    </div>
</Footer>

<style>
    @import url("https://fonts.googleapis.com/css2?family=Fredoka:wght@300..700&display=swap");
    .app {
        display: flex;
        flex-direction: column;
        align-items: center;
        color: black;
        font-family: Fredoka, sans-serif;
        background-image: url("$lib/images/background.png");
        background-size: 100%;
    }

    .mainLoginPage {
        display: flex;
        flex-direction: row;
        box-sizing: border-box;
    }

    .footer {
        display: flex;
        flex-direction: row;
        justify-content: center;
        color: black;
        font-family: Fredoka, sans-serif;
        background-color: rgb(7, 7, 7);
        padding: 5px;
    }
</style>
