package io.github.parkwithease.parkeasy.ui.profile

import androidx.annotation.StringRes
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.width
import androidx.compose.material3.Button
import androidx.compose.material3.ButtonColors
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.SnackbarHost
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
import io.github.parkwithease.parkeasy.model.FieldState
import io.github.parkwithease.parkeasy.model.Profile
import io.github.parkwithease.parkeasy.ui.common.ParkEasyTextField
import io.github.parkwithease.parkeasy.ui.common.PreviewAll
import io.github.parkwithease.parkeasy.ui.navbar.NavBar
import io.github.parkwithease.parkeasy.ui.theme.ParkEasyTheme

@Composable
fun ProfileScreen(
    onNavigateToLogin: () -> Unit,
    onNavigateToCars: () -> Unit,
    onNavigateToSpots: () -> Unit,
    navBar: @Composable () -> Unit,
    modifier: Modifier = Modifier,
    viewModel: ProfileViewModel = hiltViewModel(),
) {
    val loggedIn by viewModel.loggedIn.collectAsState(true)
    val latestOnNavigateToLogin by rememberUpdatedState(onNavigateToLogin)

    if (!loggedIn) {
        LaunchedEffect(Unit) { latestOnNavigateToLogin() }
    } else {
        val profile by viewModel.profile.collectAsState()
        val profileRoutes =
            listOf(
                ProfileButtonRoute(R.string.cars, onNavigateToCars),
                ProfileButtonRoute(R.string.spots, onNavigateToSpots),
            )

        ProfileScreen(
            profile = profile,
            profileRoutes = profileRoutes,
            onLogoutClick = viewModel::onLogoutClick,
            navBar = navBar,
            snackbarHost = { SnackbarHost(hostState = viewModel.snackbarState) },
            modifier = modifier,
        )
    }
}

@Composable
fun ProfileScreen(
    profile: Profile,
    profileRoutes: List<ProfileButtonRoute>,
    onLogoutClick: () -> Unit,
    navBar: @Composable (() -> Unit),
    snackbarHost: @Composable (() -> Unit),
    modifier: Modifier = Modifier,
) {
    Scaffold(modifier = modifier, bottomBar = navBar, snackbarHost = snackbarHost) { innerPadding ->
        Surface(Modifier.padding(innerPadding)) {
            Column(
                modifier = Modifier.fillMaxSize(),
                verticalArrangement = Arrangement.spacedBy(4.dp, Alignment.CenterVertically),
                horizontalAlignment = Alignment.CenterHorizontally,
            ) {
                ProfileDetails(profile = profile, modifier = Modifier.width(280.dp))
                profileRoutes.forEach { ProfileButton(it.name, it.onNavigateTo) }
                LogoutButton(onLogoutClick)
            }
        }
    }
}

@Composable
fun ProfileDetails(profile: Profile, modifier: Modifier = Modifier) {
    Column(modifier = modifier, verticalArrangement = Arrangement.spacedBy(4.dp)) {
        ParkEasyTextField(
            state = FieldState(profile.name),
            onValueChange = {},
            enabled = false,
            visuallyEnabled = true,
            readOnly = true,
            labelId = R.string.name,
        )
        ParkEasyTextField(
            state = FieldState(profile.email),
            onValueChange = {},
            enabled = false,
            visuallyEnabled = true,
            readOnly = true,
            labelId = R.string.email,
        )
    }
}

@Composable
fun ProfileButton(@StringRes id: Int, onClick: () -> Unit, modifier: Modifier = Modifier) {
    Button(onClick = onClick, modifier = modifier.width(280.dp)) { Text(stringResource(id)) }
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
    ) {
        Text(stringResource(R.string.logout))
    }
}

data class ProfileButtonRoute(@StringRes val name: Int, val onNavigateTo: () -> Unit)

@Suppress("detekt:UnusedPrivateMember")
@PreviewAll
@Composable
private fun ProfileScreenPreview() {
    ParkEasyTheme {
        ProfileScreen(
            profile = Profile(),
            profileRoutes =
                listOf(ProfileButtonRoute(R.string.cars) {}, ProfileButtonRoute(R.string.spots) {}),
            onLogoutClick = {},
            navBar = { NavBar() },
            snackbarHost = {},
        )
    }
}
