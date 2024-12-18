package io.github.parkwithease.parkeasy.ui.search.map

import androidx.compose.runtime.Composable
import androidx.navigation.NavController
import androidx.navigation.NavGraphBuilder
import androidx.navigation.compose.composable
import io.github.parkwithease.parkeasy.ui.common.enterAnimation
import io.github.parkwithease.parkeasy.ui.common.exitAnimation

const val MapRoute = "map"

fun NavGraphBuilder.mapScreen(onNavigateToLogin: () -> Unit, navBar: @Composable () -> Unit) {
    composable(
        route = MapRoute,
        enterTransition = { enterAnimation() },
        exitTransition = { exitAnimation() },
    ) {
        MapScreen(onNavigateToLogin = onNavigateToLogin, navBar = navBar)
    }
}

@Suppress("unused") fun NavController.navigateToMap() = this.navigate(MapRoute)
