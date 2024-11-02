package io.github.parkwithease.parkeasy.ui.profile

import androidx.navigation.NavController
import androidx.navigation.NavGraphBuilder
import androidx.navigation.NavHostController
import androidx.navigation.compose.composable

const val ProfileRoute = "profile"

fun NavGraphBuilder.profileScreen(
    onNavigateToLogin: () -> Unit,
    onNavigateToCars: () -> Unit,
    navController: NavHostController,
) {
    composable(ProfileRoute) { ProfileScreen(onNavigateToLogin, onNavigateToCars, navController) }
}

fun NavController.navigateToProfile() = this.navigate(ProfileRoute)
