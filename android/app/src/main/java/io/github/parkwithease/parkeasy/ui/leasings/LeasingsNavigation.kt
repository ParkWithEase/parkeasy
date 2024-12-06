package io.github.parkwithease.parkeasy.ui.leasings

import androidx.navigation.NavController
import androidx.navigation.NavGraphBuilder
import androidx.navigation.compose.composable
import io.github.parkwithease.parkeasy.ui.common.enterAnimation
import io.github.parkwithease.parkeasy.ui.common.exitAnimation

private const val LeasingsRoute = "leasings"

fun NavGraphBuilder.leasingsScreen() {
    composable(
        route = LeasingsRoute,
        enterTransition = { enterAnimation() },
        exitTransition = { exitAnimation() },
    ) {
        LeasingsScreen()
    }
}

fun NavController.navigateToLeasings() = this.navigate(LeasingsRoute)
