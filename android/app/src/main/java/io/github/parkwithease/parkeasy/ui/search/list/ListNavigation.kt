package io.github.parkwithease.parkeasy.ui.search.list

import androidx.compose.animation.ExitTransition
import androidx.compose.animation.fadeIn
import androidx.compose.runtime.Composable
import androidx.navigation.NavController
import androidx.navigation.NavGraphBuilder
import androidx.navigation.compose.composable

const val ListRoute = "list"

fun NavGraphBuilder.listScreen(onNavigateToLogin: () -> Unit, navBar: @Composable () -> Unit) {
    composable(
        route = ListRoute,
        enterTransition = { fadeIn() },
        exitTransition = { ExitTransition.None },
    ) {
        ListScreen(onNavigateToLogin = onNavigateToLogin, navBar = navBar)
    }
}

@Suppress("unused") fun NavController.navigateToList() = this.navigate(ListRoute)
