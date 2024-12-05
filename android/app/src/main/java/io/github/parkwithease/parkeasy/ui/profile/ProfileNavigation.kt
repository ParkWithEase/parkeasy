package io.github.parkwithease.parkeasy.ui.profile

import androidx.compose.animation.ExitTransition
import androidx.compose.animation.fadeIn
import androidx.compose.runtime.Composable
import androidx.navigation.NavController
import androidx.navigation.NavGraphBuilder
import androidx.navigation.compose.composable

const val ProfileRoute = "profile"

fun NavGraphBuilder.profileScreen(
    onNavigateToLogin: () -> Unit,
    onNavigateToCars: () -> Unit,
    onNavigateToSpots: () -> Unit,
    navBar: @Composable () -> Unit,
) {
    composable(
        route = ProfileRoute,
        enterTransition = { fadeIn() },
        exitTransition = { ExitTransition.None },
    ) {
        ProfileScreen(
            onNavigateToLogin = onNavigateToLogin,
            onNavigateToCars = onNavigateToCars,
            onNavigateToSpots = onNavigateToSpots,
            navBar = navBar,
        )
    }
}

@Suppress("unused") fun NavController.navigateToProfile() = this.navigate(ProfileRoute)
