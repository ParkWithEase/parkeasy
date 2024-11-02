package io.github.parkwithease.parkeasy.ui.profile

import androidx.navigation.NavController
import androidx.navigation.NavGraphBuilder
import androidx.navigation.NavHostController
import androidx.navigation.compose.composable

private const val ProfileRoute = "profile"

fun NavGraphBuilder.profileScreen(
    showSnackbar: suspend (String, String?) -> Boolean,
    onNavigateToMyCars: () -> Unit,
    onLogout: () -> Unit,
    navController: NavHostController,
) {
    composable(ProfileRoute) {
        ProfileScreen(showSnackbar, onNavigateToMyCars, onLogout, navController)
    }
}

fun NavController.navigateToProfile() = this.navigate(ProfileRoute)
