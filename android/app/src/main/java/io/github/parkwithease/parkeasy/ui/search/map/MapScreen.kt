package io.github.parkwithease.parkeasy.ui.search.map

import androidx.activity.compose.BackHandler
import androidx.compose.foundation.layout.padding
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Refresh
import androidx.compose.material3.FloatingActionButton
import androidx.compose.material3.Icon
import androidx.compose.material3.MaterialTheme
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
import androidx.compose.ui.Modifier
import androidx.hilt.navigation.compose.hiltViewModel
import com.google.accompanist.permissions.ExperimentalPermissionsApi
import com.google.accompanist.permissions.rememberMultiplePermissionsState
import com.mapbox.mapboxsdk.geometry.LatLng
import com.maplibre.compose.LocalMapLibreStyleProvider
import com.maplibre.compose.MapView
import com.maplibre.compose.camera.CameraState
import com.maplibre.compose.camera.MapViewCamera
import com.maplibre.compose.rememberSaveableMapViewCamera
import com.maplibre.compose.symbols.Symbol
import io.github.parkwithease.parkeasy.R
import io.github.parkwithease.parkeasy.ui.search.CreateBookingScreen
import io.github.parkwithease.parkeasy.ui.search.DefaultLatitude
import io.github.parkwithease.parkeasy.ui.search.DefaultLongitude
import io.github.parkwithease.parkeasy.ui.search.SearchViewModel
import io.github.parkwithease.parkeasy.ui.search.rememberCreateHandler

@OptIn(ExperimentalPermissionsApi::class)
@Suppress("detekt:LongMethod")
@Composable
fun MapScreen(
    onNavigateToLogin: () -> Unit,
    navBar: @Composable () -> Unit,
    modifier: Modifier = Modifier,
    viewModel: SearchViewModel = hiltViewModel<SearchViewModel>(),
) {
    val loggedIn by viewModel.loggedIn.collectAsState(false)
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

        val cameraState =
            CameraState.Centered(latitude = DefaultLatitude, longitude = DefaultLongitude)
        val mapViewCamera = rememberSaveableMapViewCamera(MapViewCamera(state = cameraState))
        val engine = remember { viewModel.engine }

        var launchedPermissions = false

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
                    mapViewCamera.value = MapViewCamera.TrackingUserLocation()
                }
            }

        BackHandler(enabled = showForm) { viewModel.onHideForm() }

        LaunchedEffect(Unit) {
            if (locationPermissions.allPermissionsGranted) {
                viewModel.startLocationFlow()
                mapViewCamera.value = MapViewCamera.TrackingUserLocation()
            }
            viewModel.snackbarState.currentSnackbarData?.dismiss()
            viewModel.onRefresh()
        }

        if (showForm)
            CreateBookingScreen(
                cars = cars,
                state = viewModel.createState,
                handler = createHandler,
                formEnabled = formEnabled,
                getSelectedIds = { viewModel.createState.selectedIds.value },
                disabledIds = viewModel.createState.disabledIds.value,
            )
        else
            Scaffold(
                modifier = modifier,
                bottomBar = navBar,
                snackbarHost = { SnackbarHost(hostState = viewModel.snackbarState) },
                floatingActionButton = {
                    RefreshButton(
                        isRefreshing = isRefreshing,
                        onRefresh = {
                            if (!launchedPermissions) {
                                launchedPermissions = true
                                if (locationPermissions.allPermissionsGranted) {
                                    viewModel.startLocationFlow()
                                    mapViewCamera.value = MapViewCamera.TrackingUserLocation()
                                } else if (locationPermissions.shouldShowRationale) {
                                    showRationale = true
                                } else {
                                    locationPermissions.launchMultiplePermissionRequest()
                                }
                            }
                            viewModel.onRefresh()
                        },
                    )
                },
            ) { innerPadding ->
                Surface(modifier = Modifier.padding(innerPadding)) {
                    MapView(
                        modifier = Modifier,
                        styleUrl = LocalMapLibreStyleProvider.current.getStyleUrl(),
                        camera = mapViewCamera,
                        locationEngine = engine,
                    ) {
                        spots.forEach {
                            Symbol(
                                center = LatLng(it.location.latitude, it.location.longitude),
                                imageId = R.drawable.location,
                                onTap = {
                                    createHandler.reset()
                                    createHandler.onSpotChange(it)
                                    viewModel.onShowForm()
                                },
                            )
                        }
                    }
                }
            }
    }
}

@Composable
fun RefreshButton(isRefreshing: Boolean, onRefresh: () -> Unit, modifier: Modifier = Modifier) {
    FloatingActionButton(
        onClick = onRefresh,
        modifier = modifier,
        containerColor =
            if (isRefreshing) MaterialTheme.colorScheme.primary
            else MaterialTheme.colorScheme.primaryContainer,
    ) {
        Icon(imageVector = Icons.Filled.Refresh, contentDescription = null)
    }
}
