package io.github.parkwithease.parkeasy.ui.profile

import androidx.navigation.NavController
import androidx.navigation.NavGraphBuilder
import androidx.navigation.NavHostController
import androidx.navigation.compose.composable

const val ProfileRoute = "profile"

fun NavGraphBuilder.profileScreen(
    onNavigateToMyCars: () -> Unit,
    onLogout: () -> Unit,
    navController: NavHostController,
) {
    composable(ProfileRoute) { ProfileScreen(onNavigateToMyCars, onLogout, navController) }
}

fun NavController.navigateToProfile() = this.navigate(ProfileRoute)
