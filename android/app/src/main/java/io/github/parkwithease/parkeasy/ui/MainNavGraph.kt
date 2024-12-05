package io.github.parkwithease.parkeasy.ui

import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.navigation.NavHostController
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.rememberNavController
import io.github.parkwithease.parkeasy.ui.cars.carsScreen
import io.github.parkwithease.parkeasy.ui.cars.navigateToCars
import io.github.parkwithease.parkeasy.ui.login.loginScreen
import io.github.parkwithease.parkeasy.ui.login.navigateToLogin
import io.github.parkwithease.parkeasy.ui.navbar.NavBar
import io.github.parkwithease.parkeasy.ui.profile.profileScreen
import io.github.parkwithease.parkeasy.ui.search.list.ListRoute
import io.github.parkwithease.parkeasy.ui.search.list.listScreen
import io.github.parkwithease.parkeasy.ui.search.map.mapScreen
import io.github.parkwithease.parkeasy.ui.spots.navigateToSpots
import io.github.parkwithease.parkeasy.ui.spots.spotsScreen

@Composable
fun MainNavGraph(
    onExitApp: () -> Unit,
    modifier: Modifier = Modifier,
    navController: NavHostController = rememberNavController(),
) {
    val navBar = @Composable { NavBar(navController = navController) }
    NavHost(navController = navController, startDestination = ListRoute, modifier = modifier) {
        loginScreen(onExitApp = onExitApp, onNavigateFromLogin = navController::popBackStack)
        listScreen(onNavigateToLogin = navController::navigateToLogin, navBar = navBar)
        mapScreen(onNavigateToLogin = navController::navigateToLogin, navBar = navBar)
        profileScreen(
            onNavigateToLogin = navController::navigateToLogin,
            onNavigateToCars = navController::navigateToCars,
            onNavigateToSpots = navController::navigateToSpots,
            navBar = navBar,
        )
        carsScreen()
        spotsScreen()
    }
}
