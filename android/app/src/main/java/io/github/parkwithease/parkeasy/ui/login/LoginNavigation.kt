package io.github.parkwithease.parkeasy.ui.login

import androidx.hilt.navigation.compose.hiltViewModel
import androidx.navigation.NavController
import androidx.navigation.NavGraphBuilder
import androidx.navigation.compose.composable

private const val LoginRoute = "login"

fun NavGraphBuilder.loginScreen(
    showSnackbar: suspend (String, String?) -> Boolean,
    onLogin: () -> Unit,
) {
    composable(LoginRoute) {
        val viewModel =
            hiltViewModel<LoginViewModel, LoginViewModel.Factory> { factory ->
                factory.create(showSnackbar = showSnackbar)
            }
        LoginScreen(onLogin, viewModel)
    }
}

fun NavController.navigateToLogin() = this.navigate(LoginRoute)
