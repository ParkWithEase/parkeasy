package io.github.parkwithease.parkeasy.ui.map

import androidx.navigation.NavController
import androidx.navigation.NavGraphBuilder
import androidx.navigation.NavHostController
import androidx.navigation.compose.composable

const val MapRoute = "map"

fun NavGraphBuilder.mapScreen(navController: NavHostController) {
    composable(MapRoute) { MapScreen(navController = navController) }
}

fun NavController.navigateToMap() = this.navigate(MapRoute)
