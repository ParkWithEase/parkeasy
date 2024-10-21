package io.github.parkwithease.parkeasy.ui.login

import androidx.navigation.NavController
import androidx.navigation.NavGraphBuilder
import androidx.navigation.compose.composable

private const val LoginRoute = "login"

fun NavGraphBuilder.loginScreen(
    showSnackbar: suspend (String, String?) -> Boolean,
    onLogin: () -> Unit,
) {
    composable(LoginRoute) { LoginScreen(showSnackbar, onLogin) }
}

fun NavController.navigateToLogin() = this.navigate(LoginRoute)
