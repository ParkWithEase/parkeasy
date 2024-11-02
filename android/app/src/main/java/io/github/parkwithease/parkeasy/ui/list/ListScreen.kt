package io.github.parkwithease.parkeasy.ui.list

import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.navigation.NavHostController
import io.github.parkwithease.parkeasy.ui.navbar.NavBar

@Composable
fun ListScreen(
    navController: NavHostController,
    modifier: Modifier = Modifier,
    @Suppress("unused") viewModel: ListViewModel = hiltViewModel<ListViewModel>(),
) {
    Scaffold(modifier = modifier, bottomBar = { NavBar(navController = navController) }) {
        innerPadding ->
        Surface(modifier = Modifier.padding(innerPadding)) { Text("List") }
    }
}
