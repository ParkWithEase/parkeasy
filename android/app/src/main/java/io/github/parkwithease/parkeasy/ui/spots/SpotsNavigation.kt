package io.github.parkwithease.parkeasy.ui.spots

import androidx.navigation.NavController
import androidx.navigation.NavGraphBuilder
import androidx.navigation.compose.composable
import io.github.parkwithease.parkeasy.model.Spot

private const val SpotsRoute = "spots"

fun NavGraphBuilder.spotsScreen(onSpotClick: (Spot) -> Unit) {
    composable(SpotsRoute) { SpotsScreen(onSpotClick) }
}

fun NavController.navigateToSpots() = this.navigate(SpotsRoute)
