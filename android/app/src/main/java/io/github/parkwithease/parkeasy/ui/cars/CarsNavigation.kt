package io.github.parkwithease.parkeasy.ui.cars

import androidx.compose.animation.ExitTransition
import androidx.compose.animation.fadeIn
import androidx.navigation.NavController
import androidx.navigation.NavGraphBuilder
import androidx.navigation.compose.composable

private const val CarsRoute = "cars"

fun NavGraphBuilder.carsScreen() {
    composable(
        route = CarsRoute,
        enterTransition = { fadeIn() },
        exitTransition = { ExitTransition.None },
    ) {
        CarsScreen()
    }
}

fun NavController.navigateToCars() = this.navigate(CarsRoute)
