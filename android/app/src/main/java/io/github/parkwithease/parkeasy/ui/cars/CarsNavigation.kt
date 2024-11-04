package io.github.parkwithease.parkeasy.ui.cars

import androidx.navigation.NavController
import androidx.navigation.NavGraphBuilder
import androidx.navigation.compose.composable

private const val CarsRoute = "cars"

fun NavGraphBuilder.carsScreen() {
    composable(CarsRoute) { CarsScreen() }
}

fun NavController.navigateToCars() = this.navigate(CarsRoute)
