<script>
    import '../app.css';
</script>

<slot></slot>

<footer class="footer">
    <p>By ParkEasy Team 2024</p>
</footer>

<style>
    /* @import url('https://fonts.googleapis.com/css2?family=Fredoka:wght@300..700&display=swap');
    .app {
        display: flex;
        flex-direction: column;
        align-items: center;
        color: black;
        font-family: Fredoka, sans-serif;
        background-image: url('$lib/images/background.png');
        background-size: 100%;
    } */

    .footer {
        position: absolute;
        text-align: center;
        width: 100%;
        bottom: 0px;
    }
</style>
