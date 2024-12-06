package io.github.parkwithease.parkeasy.ui.cars

import androidx.navigation.NavController
import androidx.navigation.NavGraphBuilder
import androidx.navigation.compose.composable
import io.github.parkwithease.parkeasy.ui.common.enterAnimation
import io.github.parkwithease.parkeasy.ui.common.exitAnimation

private const val CarsRoute = "cars"

fun NavGraphBuilder.carsScreen() {
    composable(
        route = CarsRoute,
        enterTransition = { enterAnimation() },
        exitTransition = { exitAnimation() },
    ) {
        CarsScreen()
    }
}

fun NavController.navigateToCars() = this.navigate(CarsRoute)
