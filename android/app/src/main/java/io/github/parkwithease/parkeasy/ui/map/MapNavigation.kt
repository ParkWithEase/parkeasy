package io.github.parkwithease.parkeasy.ui.map

import androidx.navigation.NavController
import androidx.navigation.NavGraphBuilder
import androidx.navigation.compose.composable

private const val MapRoute = "map"

fun NavGraphBuilder.mapScreen() {
    composable(MapRoute) { MapScreen() }
}

fun NavController.navigateToMap() = this.navigate(MapRoute)
