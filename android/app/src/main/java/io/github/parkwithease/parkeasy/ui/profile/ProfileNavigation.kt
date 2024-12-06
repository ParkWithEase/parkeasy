package io.github.parkwithease.parkeasy.ui.profile

import androidx.compose.runtime.Composable
import androidx.navigation.NavController
import androidx.navigation.NavGraphBuilder
import androidx.navigation.compose.composable
import io.github.parkwithease.parkeasy.ui.common.enterAnimation
import io.github.parkwithease.parkeasy.ui.common.exitAnimation

const val ProfileRoute = "profile"

fun NavGraphBuilder.profileScreen(
    onNavigateToLogin: () -> Unit,
    onNavigateToCars: () -> Unit,
    onNavigateToSpots: () -> Unit,
    onNavigateToBookings: () -> Unit,
    navBar: @Composable () -> Unit,
) {
    composable(
        route = ProfileRoute,
        enterTransition = { enterAnimation() },
        exitTransition = { exitAnimation() },
    ) {
        ProfileScreen(
            onNavigateToLogin = onNavigateToLogin,
            onNavigateToCars = onNavigateToCars,
            onNavigateToSpots = onNavigateToSpots,
            onNavigateToBookings = onNavigateToBookings,
            navBar = navBar,
        )
    }
}

@Suppress("unused") fun NavController.navigateToProfile() = this.navigate(ProfileRoute)
