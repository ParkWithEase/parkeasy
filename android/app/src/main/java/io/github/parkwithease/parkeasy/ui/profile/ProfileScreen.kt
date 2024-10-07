package io.github.parkwithease.parkeasy.ui.profile

import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Button
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.ui.Modifier
import androidx.hilt.navigation.compose.hiltViewModel
import io.github.parkwithease.parkeasy.ui.navbar.NavBar

@Composable
fun ProfileScreen(
    onLogout: () -> Unit,
    modifier: Modifier = Modifier,
    viewModel: ProfileViewModel = hiltViewModel<ProfileViewModel>(),
) {
    Scaffold(bottomBar = { NavBar() }, modifier = modifier) { innerPadding ->
        Surface(Modifier.padding(innerPadding)) {
            Button({
                viewModel.onLogoutPress()
                onLogout()
            }) {
                Text("Logout")
            }
        }
    }
    val onLogoutEvent: () -> Unit = onLogout
    LaunchedEffect(viewModel.loggedIn) {
        if (!(viewModel.loggedIn.value)) {
            onLogoutEvent()
        }
    }
}
