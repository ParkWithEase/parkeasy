package io.github.parkwithease.parkeasy.ui.login

import androidx.navigation.NavController
import androidx.navigation.NavGraphBuilder
import androidx.navigation.compose.composable

private const val LoginRoute = "login"

fun NavGraphBuilder.loginScreen(onExitApp: () -> Unit, onNavigateFromLogin: () -> Unit) {
    composable(LoginRoute) {
        LoginScreen(onExitApp = onExitApp, onNavigateFromLogin = onNavigateFromLogin)
    }
}

fun NavController.navigateToLogin() = this.navigate(LoginRoute)
