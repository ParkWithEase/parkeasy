package io.github.parkwithease.parkeasy.ui.cars

import androidx.navigation.NavController
import androidx.navigation.NavGraphBuilder
import androidx.navigation.compose.composable
import io.github.parkwithease.parkeasy.model.Car

private const val CarsRoute = "cars"

fun NavGraphBuilder.carsScreen(onCarClick: (Car) -> Unit) {
    composable(CarsRoute) { CarsScreen(onCarClick) }
}

fun NavController.navigateToCars() = this.navigate(CarsRoute)
