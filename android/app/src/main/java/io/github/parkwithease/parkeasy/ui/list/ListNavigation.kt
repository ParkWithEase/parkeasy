package io.github.parkwithease.parkeasy.ui.list

import androidx.navigation.NavController
import androidx.navigation.NavGraphBuilder
import androidx.navigation.NavHostController
import androidx.navigation.compose.composable

const val ListRoute = "list"

fun NavGraphBuilder.listScreen(onNavigateToLogin: () -> Unit, navController: NavHostController) {
    composable(ListRoute) {
        ListScreen(onNavigateToLogin = onNavigateToLogin, navController = navController)
    }
}

fun NavController.navigateToList() = this.navigate(ListRoute)
