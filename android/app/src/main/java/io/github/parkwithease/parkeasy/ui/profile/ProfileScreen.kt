package io.github.parkwithease.parkeasy.ui.profile

import androidx.annotation.StringRes
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
    Surface(modifier) {
        Column(
            modifier = Modifier.fillMaxSize(),
            verticalArrangement = Arrangement.Center,
            horizontalAlignment = Alignment.CenterHorizontally,
        ) {
            val profile by viewModel.profile.collectAsState()
            ProfileText(profile.name, R.string.name)
            ProfileText(profile.email, R.string.email)
            ProfileButton({}, R.string.booking_history)
            ProfileButton({}, R.string.listing_history)
            ProfileButton({}, R.string.preferences_cars)
            ProfileButton(viewModel::onLogoutPress, R.string.logout)
        }
    }
}

@Composable
fun ProfileText(value: String, @StringRes labelId: Int, modifier: Modifier = Modifier) {
    OutlinedTextField(
        value = value,
        label = { Text(stringResource(labelId)) },
        onValueChange = {},
        readOnly = true,
        singleLine = true,
        modifier = modifier,
    )
}

@Composable
fun ProfileButton(onLogout: () -> Unit, @StringRes labelId: Int, modifier: Modifier = Modifier) {
    Button(onClick = onLogout, modifier = modifier) { Text(stringResource(labelId)) }
}
