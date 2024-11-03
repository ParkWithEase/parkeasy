package io.github.parkwithease.parkeasy.ui.list

import androidx.compose.runtime.Composable
import androidx.navigation.NavController
import androidx.navigation.NavGraphBuilder
import androidx.navigation.compose.composable

const val ListRoute = "list"

fun NavGraphBuilder.listScreen(onNavigateToLogin: () -> Unit, navBar: @Composable () -> Unit) {
    composable(ListRoute) { ListScreen(onNavigateToLogin = onNavigateToLogin, navBar = navBar) }
}

@Suppress("unused") fun NavController.navigateToList() = this.navigate(ListRoute)
