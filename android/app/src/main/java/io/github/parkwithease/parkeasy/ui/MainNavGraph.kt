package io.github.parkwithease.parkeasy.ui

import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Scaffold
import androidx.compose.material3.SnackbarHost
import androidx.compose.material3.SnackbarHostState
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.navigation.NavHostController
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.rememberNavController
import io.github.parkwithease.parkeasy.ui.list.listScreen
import io.github.parkwithease.parkeasy.ui.list.navigateToList
import io.github.parkwithease.parkeasy.ui.login.loginScreen
import io.github.parkwithease.parkeasy.ui.login.navigateToLogin
import io.github.parkwithease.parkeasy.ui.map.mapScreen
import io.github.parkwithease.parkeasy.ui.map.navigateToMap
import io.github.parkwithease.parkeasy.ui.navbar.NavBar
import io.github.parkwithease.parkeasy.ui.profile.navigateToProfile
import io.github.parkwithease.parkeasy.ui.profile.profileScreen

@Composable
fun MainNavGraph(
    hostState: SnackbarHostState,
    modifier: Modifier = Modifier,
    navController: NavHostController = rememberNavController(),
) {
    Scaffold(
        bottomBar = {
            NavBar(
                navController::navigateToList,
                navController::navigateToMap,
                navController::navigateToProfile,
            )
        },
        snackbarHost = { SnackbarHost(hostState = hostState) },
        modifier = modifier,
    ) { innerPadding ->
        NavHost(
            navController = navController,
            startDestination = "login",
            modifier = Modifier.padding(innerPadding),
        ) {
            loginScreen { navController.navigateToProfile() }
            listScreen()
            mapScreen()
            profileScreen { navController.navigateToLogin() }
        }
    }
}
