package io.github.parkwithease.parkeasy.ui.spots

import androidx.compose.animation.ExitTransition
import androidx.compose.animation.fadeIn
import androidx.navigation.NavController
import androidx.navigation.NavGraphBuilder
import androidx.navigation.compose.composable

private const val SpotsRoute = "spots"

fun NavGraphBuilder.spotsScreen() {
    composable(
        route = SpotsRoute,
        enterTransition = { fadeIn() },
        exitTransition = { ExitTransition.None },
    ) {
        SpotsScreen()
    }
}

fun NavController.navigateToSpots() = this.navigate(SpotsRoute)
