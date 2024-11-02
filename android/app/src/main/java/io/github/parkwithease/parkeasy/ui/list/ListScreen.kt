package io.github.parkwithease.parkeasy.ui.list

import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.rememberUpdatedState
import androidx.compose.ui.Modifier
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.navigation.NavHostController
import io.github.parkwithease.parkeasy.ui.navbar.NavBar

@Composable
fun ListScreen(
    onNavigateToLogin: () -> Unit,
    navController: NavHostController,
    modifier: Modifier = Modifier,
    viewModel: ListViewModel = hiltViewModel<ListViewModel>(),
) {
    val latestOnNavigateToLogin by rememberUpdatedState(onNavigateToLogin)
    val loggedIn by viewModel.loggedIn.collectAsState(true)

    if (!loggedIn) {
        LaunchedEffect(Unit) { latestOnNavigateToLogin() }
    } else {
        Scaffold(modifier = modifier, bottomBar = { NavBar(navController = navController) }) {
            innerPadding ->
            Surface(modifier = Modifier.padding(innerPadding)) { Text("List") }
        }
    }
}
