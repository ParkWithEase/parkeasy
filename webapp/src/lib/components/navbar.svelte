<script lang="ts">
    import logo from '$lib/images/logo-text.png';
    </script>
    
    <svelte:window/>
    <nav class="navbar">
        <a href="/#"
            ><img src={logo} alt="parkeasy-logo" class="text-logo"
        /></a>
        <a class="nav-link" href="/profile/user-profile">book a spot</a>
        <a class="nav-link" href="/profile/user-profile">list your spot</a>
        <a class="nav-link" href="/profile/user-profile">my profile</a>
        <a class="nav-link" href="/">about us</a>
    </nav>
    
    <style> 
    
        .navbar {   
            position: fixed;
            top: 0;
            width: 100%;  
            height: 7rem; 
            background-color: #000;   
            color: #FFF;   
            z-index: 999;  
            display: flex;   
            flex-direction: row;
            justify-content: space-evenly; 
            align-items: center; 
        } 
    
        .nav-link {
            background-color: white;
            padding: 1rem 2rem;
            border-radius: 20px;
            font-size: 1.3rem;
            font-weight: bold;
            color: black;
            text-decoration: none;
        }

        .nav-link:hover{
            color: #7ed957;
            background-color: rgb(54, 54, 54);
            transition: 0.3s;
        }

    
        .text-logo{
            max-height: 7rem;
        }
        
          
    </style> 