package io.github.parkwithease.parkeasy.ui.spots

import androidx.navigation.NavController
import androidx.navigation.NavGraphBuilder
import androidx.navigation.compose.composable
import io.github.parkwithease.parkeasy.ui.common.enterAnimation
import io.github.parkwithease.parkeasy.ui.common.exitAnimation

private const val SpotsRoute = "spots"

fun NavGraphBuilder.spotsScreen() {
    composable(
        route = SpotsRoute,
        enterTransition = { enterAnimation() },
        exitTransition = { exitAnimation() },
    ) {
        SpotsScreen()
    }
}

fun NavController.navigateToSpots() = this.navigate(SpotsRoute)
