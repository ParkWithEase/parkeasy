package io.github.parkwithease.parkeasy.ui.login

import androidx.compose.animation.ExitTransition
import androidx.compose.animation.fadeIn
import androidx.navigation.NavController
import androidx.navigation.NavGraphBuilder
import androidx.navigation.compose.composable

private const val LoginRoute = "login"

fun NavGraphBuilder.loginScreen(onExitApp: () -> Unit, onNavigateFromLogin: () -> Unit) {
    composable(
        route = LoginRoute,
        enterTransition = { fadeIn() },
        exitTransition = { ExitTransition.None },
    ) {
        LoginScreen(onExitApp = onExitApp, onNavigateFromLogin = onNavigateFromLogin)
    }
}

fun NavController.navigateToLogin() = this.navigate(LoginRoute)
