package io.github.parkwithease.parkeasy.ui.login

import androidx.navigation.NavController
import androidx.navigation.NavGraphBuilder
import androidx.navigation.compose.composable
import io.github.parkwithease.parkeasy.ui.common.enterAnimation
import io.github.parkwithease.parkeasy.ui.common.exitAnimation

private const val LoginRoute = "login"

fun NavGraphBuilder.loginScreen(onExitApp: () -> Unit, onNavigateFromLogin: () -> Unit) {
    composable(
        route = LoginRoute,
        enterTransition = { enterAnimation() },
        exitTransition = { exitAnimation() },
    ) {
        LoginScreen(onExitApp = onExitApp, onNavigateFromLogin = onNavigateFromLogin)
    }
}

fun NavController.navigateToLogin() = this.navigate(LoginRoute)
