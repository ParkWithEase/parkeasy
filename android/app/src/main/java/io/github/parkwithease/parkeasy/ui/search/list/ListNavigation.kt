package io.github.parkwithease.parkeasy.ui.search.list

import androidx.compose.runtime.Composable
import androidx.navigation.NavController
import androidx.navigation.NavGraphBuilder
import androidx.navigation.compose.composable
import io.github.parkwithease.parkeasy.ui.common.enterAnimation
import io.github.parkwithease.parkeasy.ui.common.exitAnimation

const val ListRoute = "list"

fun NavGraphBuilder.listScreen(onNavigateToLogin: () -> Unit, navBar: @Composable () -> Unit) {
    composable(
        route = ListRoute,
        enterTransition = { enterAnimation() },
        exitTransition = { exitAnimation() },
    ) {
        ListScreen(onNavigateToLogin = onNavigateToLogin, navBar = navBar)
    }
}

@Suppress("unused") fun NavController.navigateToList() = this.navigate(ListRoute)
