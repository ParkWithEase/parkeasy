package io.github.parkwithease.parkeasy.ui

import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.navigation.NavHostController
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.rememberNavController
import io.github.parkwithease.parkeasy.ui.cars.carsScreen
import io.github.parkwithease.parkeasy.ui.cars.navigateToCars
import io.github.parkwithease.parkeasy.ui.list.ListRoute
import io.github.parkwithease.parkeasy.ui.list.listScreen
import io.github.parkwithease.parkeasy.ui.login.loginScreen
import io.github.parkwithease.parkeasy.ui.login.navigateToLogin
import io.github.parkwithease.parkeasy.ui.map.mapScreen
import io.github.parkwithease.parkeasy.ui.profile.profileScreen

@Composable
fun MainNavGraph(
    modifier: Modifier = Modifier,
    navController: NavHostController = rememberNavController(),
) {
    NavHost(navController = navController, startDestination = ListRoute, modifier = modifier) {
        loginScreen(onLogin = navController::popBackStack)
        listScreen(
            onNavigateToLogin = navController::navigateToLogin,
            navController = navController,
        )
        mapScreen(onNavigateToLogin = navController::navigateToLogin, navController = navController)
        profileScreen(
            onNavigateToLogin = navController::navigateToLogin,
            onNavigateToCars = navController::navigateToCars,
            navController = navController,
        )
        carsScreen {}
    }
}
