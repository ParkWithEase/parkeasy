package io.github.parkwithease.parkeasy.ui.spots

import androidx.navigation.NavController
import androidx.navigation.NavGraphBuilder
import androidx.navigation.compose.composable

private const val SpotsRoute = "spots"

fun NavGraphBuilder.spotsScreen() {
    composable(SpotsRoute) { SpotsScreen() }
}

fun NavController.navigateToSpots() = this.navigate(SpotsRoute)
