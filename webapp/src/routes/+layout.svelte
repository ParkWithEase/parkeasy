<script>import '../app.css';</script>

<div class="app">
    <main>
        <slot></slot>
    </main>
</div>


<style>
    .app {
        display: flex;
        flex-direction: column;
        align-items: center;
        min-height: 100vh;
        font-family:'Trebuchet MS', 'Lucida Sans Unicode', 'Lucida Grande', 'Lucida Sans', Arial, sans-serif;
        color: black;
        background-image: url('../lib/images/background.jpg');
    }

    main {
        flex: 1;
        display: flex;
        flex-direction: column;
        padding: 1rem;
        width: 100%;
        max-width: 64rem;
        margin: 0 auto;
        box-sizing: border-box;
    }
</style>
