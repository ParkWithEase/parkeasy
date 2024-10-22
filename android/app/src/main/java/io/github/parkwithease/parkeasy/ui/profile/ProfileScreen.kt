package io.github.parkwithease.parkeasy.ui.profile

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.Button
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.rememberUpdatedState
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.res.stringResource
import androidx.hilt.navigation.compose.hiltViewModel
import io.github.parkwithease.parkeasy.R
import io.github.parkwithease.parkeasy.model.Profile

@Composable
fun ProfileScreen(
    onLogout: () -> Unit,
    modifier: Modifier = Modifier,
    viewModel: ProfileViewModel = hiltViewModel<ProfileViewModel>(),
) {
    val loggedIn by viewModel.loggedIn.collectAsState(true)
    val latestOnLogout by rememberUpdatedState(onLogout)
    LaunchedEffect(loggedIn) {
        if (!loggedIn) {
            latestOnLogout()
        }
    }
    val profile by viewModel.profile.collectAsState()
    ProfileScreen(profile, viewModel::onLogoutClick, modifier)
}

@Composable
fun ProfileScreen(profile: Profile, onLogoutClick: () -> Unit, modifier: Modifier = Modifier) {
    Surface(modifier) {
        Column(
            modifier = Modifier.fillMaxSize(),
            verticalArrangement = Arrangement.Center,
            horizontalAlignment = Alignment.CenterHorizontally,
        ) {
            ProfileDetails(profile)
            Button(onClick = onLogoutClick) { Text("Logout 2") }
        }
    }
}

@Composable
fun ProfileDetails(profile: Profile, modifier: Modifier = Modifier) {
    Column(verticalArrangement = Arrangement.SpaceBetween, modifier = modifier) {
        OutlinedTextField(
            value = profile.name,
            onValueChange = {},
            label = { Text(stringResource(R.string.name)) },
        )
        OutlinedTextField(
            value = profile.email,
            onValueChange = {},
            label = { Text(stringResource(R.string.email)) },
        )
    }
}
