package io.github.parkwithease.parkeasy.ui

import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.navigation.NavHostController
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.rememberNavController
import io.github.parkwithease.parkeasy.ui.login.loginScreen
import io.github.parkwithease.parkeasy.ui.login.navigateToLogin
import io.github.parkwithease.parkeasy.ui.profile.navigateToProfile
import io.github.parkwithease.parkeasy.ui.profile.profileScreen

@Composable
fun MainNavGraph(
    modifier: Modifier = Modifier,
    navController: NavHostController = rememberNavController(),
) {
    NavHost(navController = navController, startDestination = "login", modifier = modifier) {
        loginScreen { navController.navigateToProfile() }
        profileScreen { navController.navigateToLogin() }
    }
}
