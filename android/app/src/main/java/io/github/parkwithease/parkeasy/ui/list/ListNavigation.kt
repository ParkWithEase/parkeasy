package io.github.parkwithease.parkeasy.ui.list

import androidx.navigation.NavController
import androidx.navigation.NavGraphBuilder
import androidx.navigation.NavHostController
import androidx.navigation.compose.composable

private const val ListRoute = "list"

fun NavGraphBuilder.listScreen(navController: NavHostController) {
    composable(ListRoute) { ListScreen(navController = navController) }
}

fun NavController.navigateToList() = this.navigate(ListRoute)
