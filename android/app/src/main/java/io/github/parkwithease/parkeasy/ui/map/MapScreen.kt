package io.github.parkwithease.parkeasy.ui.map

import android.location.Location
import android.util.Log
import androidx.compose.foundation.layout.padding
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Place
import androidx.compose.material3.FloatingActionButton
import androidx.compose.material3.Icon
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Surface
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.MutableState
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.setValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.rememberUpdatedState
import androidx.compose.ui.Modifier
import androidx.hilt.navigation.compose.hiltViewModel
import com.google.accompanist.permissions.ExperimentalPermissionsApi
import com.google.accompanist.permissions.rememberMultiplePermissionsState
import com.mapbox.mapboxsdk.location.engine.LocationEngine
import com.maplibre.compose.LocalMapLibreStyleProvider
import com.maplibre.compose.MapView
import com.maplibre.compose.camera.CameraState
import com.maplibre.compose.camera.MapViewCamera
import com.maplibre.compose.rememberSaveableMapViewCamera
import io.github.parkwithease.parkeasy.ui.common.PreviewAll
import io.github.parkwithease.parkeasy.ui.navbar.NavBar
import io.github.parkwithease.parkeasy.ui.theme.ParkEasyTheme

// Somewhere above Winnipeg
private const val DefaultLatitude = 49.9
private const val DefaultLongitude = -97.1

@OptIn(ExperimentalPermissionsApi::class)
@Composable
fun MapScreen(
    onNavigateToLogin: () -> Unit,
    navBar: @Composable () -> Unit,
    modifier: Modifier = Modifier,
    viewModel: MapViewModel = hiltViewModel<MapViewModel>(),
) {
    val loggedIn by viewModel.loggedIn.collectAsState(true)
    val latestOnNavigateToLogin by rememberUpdatedState(onNavigateToLogin)
    var showRationale by remember { mutableStateOf(false) }


    val cameraState = CameraState.Centered(latitude = DefaultLatitude, longitude = DefaultLongitude)
    val mapViewCamera = rememberSaveableMapViewCamera(MapViewCamera(state = cameraState))

    val locationPermissions = rememberMultiplePermissionsState(
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

    if (!loggedIn) {
        LaunchedEffect(Unit) { latestOnNavigateToLogin() }
    } else {
        MapScreen(
            begForPermission = showRationale,
            map = { MapLibreMap(mapViewCamera, viewModel.engine, modifier = it) },
            navBar = navBar,
            onUseUserLocation = {
                if (locationPermissions.allPermissionsGranted) {
                    viewModel.startLocationFlow()
                    mapViewCamera.value = MapViewCamera.TrackingUserLocation()
                } else if (locationPermissions.shouldShowRationale) {
                    showRationale = true
                } else {
                    locationPermissions.launchMultiplePermissionRequest()
                }
            },
            onRequestLocationPermission = {},
            modifier = modifier,
        )
    }
}

@Composable
fun MapScreen(
    begForPermission: Boolean,
    map: @Composable ((modifier: Modifier) -> Unit),
    navBar: @Composable (() -> Unit),
    onUseUserLocation: () -> Unit,
    onRequestLocationPermission: () -> Unit,
    modifier: Modifier = Modifier,
) {
    Scaffold(
        modifier = modifier,
        bottomBar = navBar,
        floatingActionButton = {
            FloatingActionButton(onClick = onUseUserLocation) { Icon(Icons.Filled.Place, "Get current location") }
        }
    ) { innerPadding ->
        Surface(modifier = Modifier.padding(innerPadding)) {
            if (begForPermission) {
                // Show a dialog here
            }
            map(Modifier)
        }
    }
}

@Composable
fun MapLibreMap(camera: MutableState<MapViewCamera>, engine: LocationEngine, modifier: Modifier = Modifier) {
    MapView(
        modifier = modifier,
        styleUrl = LocalMapLibreStyleProvider.current.getStyleUrl(),
        locationEngine = engine,
        camera = camera,
    )
}

@Suppress("detekt:UnusedPrivateMember")
@PreviewAll
@Composable
private fun MapScreenPreview() {
    ParkEasyTheme { MapScreen(
        begForPermission = false,
        map = {},
        navBar = { NavBar() },
        onRequestLocationPermission = {},
        onUseUserLocation = {}
    ) }
}
