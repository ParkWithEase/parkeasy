package io.github.parkwithease.parkeasy.ui.profile

import androidx.navigation.NavController
import androidx.navigation.NavGraphBuilder
import androidx.navigation.compose.composable

private const val ProfileRoute = "profile"

fun NavGraphBuilder.profileScreen(
    showSnackbar: suspend (String, String?) -> Boolean,
    onNavigateToMyCars: () -> Unit,
    onLogout: () -> Unit,
) {
    composable(ProfileRoute) { ProfileScreen(showSnackbar, onNavigateToMyCars, onLogout) }
}

fun NavController.navigateToProfile() = this.navigate(ProfileRoute)
