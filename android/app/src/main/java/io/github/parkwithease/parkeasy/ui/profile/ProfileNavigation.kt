package io.github.parkwithease.parkeasy.ui.profile

import androidx.compose.runtime.Composable
import androidx.navigation.NavController
import androidx.navigation.NavGraphBuilder
import androidx.navigation.compose.composable

const val ProfileRoute = "profile"

fun NavGraphBuilder.profileScreen(
    onNavigateToLogin: () -> Unit,
    onNavigateToCars: () -> Unit,
    navBar: @Composable () -> Unit,
) {
    composable(ProfileRoute) {
        ProfileScreen(
            onNavigateToLogin = onNavigateToLogin,
            onNavigateToCars = onNavigateToCars,
            navBar = navBar,
        )
    }
}

@Suppress("unused")
fun NavController.navigateToProfile() = this.navigate(ProfileRoute)
