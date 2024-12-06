package io.github.parkwithease.parkeasy.ui.search.list

import androidx.activity.compose.BackHandler
import androidx.compose.animation.AnimatedVisibility
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.Scaffold
import androidx.compose.material3.SnackbarHost
import androidx.compose.material3.Surface
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.rememberUpdatedState
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import com.google.accompanist.permissions.ExperimentalPermissionsApi
import com.google.accompanist.permissions.rememberMultiplePermissionsState
import io.github.parkwithease.parkeasy.model.Spot
import io.github.parkwithease.parkeasy.model.SpotLocation
import io.github.parkwithease.parkeasy.ui.common.PreviewAll
import io.github.parkwithease.parkeasy.ui.common.PullToRefreshBox
import io.github.parkwithease.parkeasy.ui.common.enterAnimation
import io.github.parkwithease.parkeasy.ui.common.exitAnimation
import io.github.parkwithease.parkeasy.ui.navbar.NavBar
import io.github.parkwithease.parkeasy.ui.search.CreateBookingScreen
import io.github.parkwithease.parkeasy.ui.search.SearchCard
import io.github.parkwithease.parkeasy.ui.search.SearchViewModel
import io.github.parkwithease.parkeasy.ui.search.rememberCreateHandler
import io.github.parkwithease.parkeasy.ui.theme.ParkEasyTheme

@OptIn(ExperimentalMaterial3Api::class, ExperimentalPermissionsApi::class)
@Suppress("detekt:LongMethod")
@Composable
fun ListScreen(
    onNavigateToLogin: () -> Unit,
    navBar: @Composable () -> Unit,
    modifier: Modifier = Modifier,
    viewModel: SearchViewModel = hiltViewModel<SearchViewModel>(),
) {
    val loggedIn by viewModel.loggedIn.collectAsState(true)
    val latestOnNavigateToLogin by rememberUpdatedState(onNavigateToLogin)
    var showRationale by remember { mutableStateOf(false) }

    if (!loggedIn) {
        LaunchedEffect(Unit) { latestOnNavigateToLogin() }
    } else {
        val createHandler = rememberCreateHandler(viewModel)
        val cars by viewModel.cars.collectAsState()
        val spots by viewModel.spots.collectAsState()
        val formEnabled by viewModel.formEnabled.collectAsState()
        val isRefreshing by viewModel.isRefreshing.collectAsState()
        val showForm by viewModel.showForm.collectAsState()

        val locationPermissions =
            rememberMultiplePermissionsState(
                listOf(
                    android.Manifest.permission.ACCESS_COARSE_LOCATION,
                    android.Manifest.permission.ACCESS_FINE_LOCATION,
                )
            ) {
                if (it.all { entry -> entry.value }) {
                    showRationale = false
                    viewModel.startLocationFlow()
                    viewModel.onRefresh()
                }
            }

        BackHandler(enabled = showForm) { viewModel.onHideForm() }

        LaunchedEffect(Unit) {
            viewModel.snackbarState.currentSnackbarData?.dismiss()
            viewModel.onRefresh()
        }

        AnimatedVisibility(visible = showForm, enter = enterAnimation(), exit = exitAnimation()) {
            CreateBookingScreen(
                cars = cars,
                state = viewModel.createState,
                handler = createHandler,
                formEnabled = formEnabled,
                getSelectedIds = { viewModel.createState.selectedIds.value },
                disabledIds = viewModel.createState.disabledIds.value,
            )
        }
        AnimatedVisibility(visible = !showForm, enter = enterAnimation(), exit = exitAnimation()) {
            ListScreen(
                spots = spots,
                onSpotClick = {
                    createHandler.reset()
                    createHandler.onSpotChange(it)
                    viewModel.onShowForm()
                },
                isRefreshing = isRefreshing,
                onRefresh = {
                    if (locationPermissions.allPermissionsGranted) {
                        viewModel.startLocationFlow()
                    } else if (locationPermissions.shouldShowRationale) {
                        showRationale = true
                    } else {
                        locationPermissions.launchMultiplePermissionRequest()
                    }
                    viewModel.onRefresh()
                },
                navBar = navBar,
                snackbarHost = { SnackbarHost(hostState = viewModel.snackbarState) },
                modifier = modifier,
            )
        }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun ListScreen(
    spots: List<Spot>,
    onSpotClick: (Spot) -> Unit,
    isRefreshing: Boolean,
    onRefresh: () -> Unit,
    navBar: @Composable () -> Unit,
    snackbarHost: @Composable (() -> Unit),
    modifier: Modifier = Modifier,
) {
    Scaffold(modifier = modifier, bottomBar = navBar, snackbarHost = snackbarHost) { innerPadding ->
        Surface(modifier = Modifier.padding(innerPadding)) {
            Column(
                modifier = Modifier.fillMaxSize(),
                verticalArrangement = Arrangement.Center,
                horizontalAlignment = Alignment.CenterHorizontally,
            ) {
                PullToRefreshBox(
                    items = spots,
                    onClick = onSpotClick,
                    isRefreshing = isRefreshing,
                    onRefresh = onRefresh,
                    modifier = Modifier.padding(4.dp),
                ) { spot, onClick ->
                    SearchCard(spot, onClick)
                }
            }
        }
    }
}

@Suppress("detekt:UnusedPrivateMember")
@PreviewAll
@Composable
private fun ListScreenPreview() {
    val spots =
        listOf(
            Spot(
                location =
                    SpotLocation(
                        streetAddress = "66 Chancellors Cir",
                        city = "Winnipeg",
                        state = "MB",
                        countryCode = "CA",
                        postalCode = "R3T2N2",
                    )
            )
        )
    ParkEasyTheme {
        ListScreen(
            spots = spots,
            onSpotClick = {},
            isRefreshing = false,
            onRefresh = {},
            navBar = { NavBar() },
            snackbarHost = {},
        )
    }
}
