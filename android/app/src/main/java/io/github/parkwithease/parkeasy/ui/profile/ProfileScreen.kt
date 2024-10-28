package io.github.parkwithease.parkeasy.ui.profile

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.width
import androidx.compose.material3.Button
import androidx.compose.material3.ButtonColors
import androidx.compose.material3.MaterialTheme
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
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import io.github.parkwithease.parkeasy.R
import io.github.parkwithease.parkeasy.model.Profile

@Composable
fun ProfileScreen(
    showSnackbar: suspend (String, String?) -> Boolean,
    onLogout: () -> Unit,
    modifier: Modifier = Modifier,
    viewModel: ProfileViewModel =
        hiltViewModel<ProfileViewModel, ProfileViewModel.Factory>(
            creationCallback = { factory -> factory.create(showSnackbar = showSnackbar) }
        ),
) {
    val loggedIn by viewModel.loggedIn.collectAsState(true)
    val latestOnLogout by rememberUpdatedState(onLogout)
    LaunchedEffect(loggedIn) {
        if (!loggedIn) {
            latestOnLogout()
        }
    }
    LaunchedEffect(Unit) { viewModel.refresh() }
    val profile by viewModel.profile.collectAsState()
    ProfileScreen(profile, viewModel::onLogoutClick, modifier)
}

@Composable
fun ProfileScreen(profile: Profile, onLogoutClick: () -> Unit, modifier: Modifier = Modifier) {
    Surface(modifier) {
        Column(
            modifier = Modifier.fillMaxSize(),
            verticalArrangement = Arrangement.spacedBy(4.dp, Alignment.CenterVertically),
            horizontalAlignment = Alignment.CenterHorizontally,
        ) {
            ProfileDetails(profile)
            LogoutButton(onLogoutClick)
        }
    }
}

@Composable
fun ProfileDetails(profile: Profile, modifier: Modifier = Modifier) {
    Column(modifier = modifier, verticalArrangement = Arrangement.spacedBy(4.dp)) {
        OutlinedTextField(
            value = profile.name,
            onValueChange = {},
            modifier = Modifier.width(280.dp),
            readOnly = true,
            label = { Text(stringResource(R.string.name)) },
            singleLine = true,
        )
        OutlinedTextField(
            value = profile.email,
            onValueChange = {},
            modifier = Modifier.width(280.dp),
            readOnly = true,
            label = { Text(stringResource(R.string.email)) },
            singleLine = true,
        )
    }
}

@Composable
fun LogoutButton(onLogoutClick: () -> Unit, modifier: Modifier = Modifier) {
    Button(
        onClick = onLogoutClick,
        modifier = modifier.width(280.dp),
        colors =
            ButtonColors(
                containerColor = MaterialTheme.colorScheme.error,
                contentColor = MaterialTheme.colorScheme.onError,
                disabledContainerColor = MaterialTheme.colorScheme.errorContainer,
                disabledContentColor = MaterialTheme.colorScheme.onErrorContainer,
            ),
        content = { Text(stringResource(R.string.logout)) },
    )
}
