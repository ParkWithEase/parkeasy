package io.github.parkwithease.parkeasy.ui

import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Scaffold
import androidx.compose.material3.SnackbarDuration
import androidx.compose.material3.SnackbarHost
import androidx.compose.material3.SnackbarHostState
import androidx.compose.material3.SnackbarResult
import androidx.compose.runtime.Composable
import androidx.compose.runtime.remember
import androidx.compose.ui.Modifier
import androidx.navigation.NavHostController
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.rememberNavController
import io.github.parkwithease.parkeasy.ui.cars.carsScreen
import io.github.parkwithease.parkeasy.ui.cars.navigateToCars
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
    modifier: Modifier = Modifier,
    navController: NavHostController = rememberNavController(),
) {
    val snackbarHostState = remember { SnackbarHostState() }
    val showSnackbar: suspend (message: String, actionLabel: String?) -> Boolean =
        { message, actionLabel ->
            snackbarHostState.showSnackbar(
                message,
                actionLabel,
                duration = SnackbarDuration.Short,
            ) == SnackbarResult.ActionPerformed
        }
    Scaffold(
        bottomBar = {
            NavBar(
                navController::navigateToList,
                navController::navigateToMap,
                navController::navigateToProfile,
            )
        },
        snackbarHost = { SnackbarHost(hostState = snackbarHostState) },
        modifier = modifier,
    ) { innerPadding ->
        NavHost(
            navController = navController,
            startDestination = "login",
            modifier = Modifier.padding(innerPadding),
        ) {
            loginScreen(showSnackbar, navController::navigateToProfile)
            listScreen()
            mapScreen()
            profileScreen(
                showSnackbar,
                navController::navigateToCars,
                navController::navigateToLogin,
            )
            carsScreen {}
        }
    }
}
