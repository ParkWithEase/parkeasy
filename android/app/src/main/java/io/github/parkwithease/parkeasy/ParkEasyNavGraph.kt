package io.github.parkwithease.parkeasy

import androidx.compose.runtime.Composable
import androidx.navigation.NavHostController
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import io.github.parkwithease.parkeasy.ui.login.LoginScreen

@Composable
fun ParkEasyNavGraph(navController: NavHostController = rememberNavController()) {
    NavHost(navController = navController, startDestination = "Login") {
        composable(route = "Login") { LoginScreen() }
    }
}
