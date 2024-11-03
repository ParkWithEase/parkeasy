package io.github.parkwithease.parkeasy.ui.map

import androidx.compose.runtime.Composable
import androidx.navigation.NavController
import androidx.navigation.NavGraphBuilder
import androidx.navigation.compose.composable

const val MapRoute = "map"

fun NavGraphBuilder.mapScreen(onNavigateToLogin: () -> Unit, navBar: @Composable () -> Unit) {
    composable(MapRoute) { MapScreen(onNavigateToLogin = onNavigateToLogin, navBar = navBar) }
}

@Suppress("unused")
fun NavController.navigateToMap() = this.navigate(MapRoute)
