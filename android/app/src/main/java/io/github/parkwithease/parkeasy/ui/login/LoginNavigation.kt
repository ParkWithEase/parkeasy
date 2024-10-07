package io.github.parkwithease.parkeasy.ui.login

import androidx.navigation.NavGraphBuilder
import androidx.navigation.compose.composable

private const val LoginRoute = "login"

fun NavGraphBuilder.loginScreen(onLogin: () -> Unit) {
    composable(LoginRoute) { LoginScreen(onLogin) }
}
