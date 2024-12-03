package io.github.parkwithease.parkeasy.ui.map

import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Surface
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.rememberUpdatedState
import androidx.compose.ui.Modifier
import androidx.hilt.navigation.compose.hiltViewModel
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

@Composable
fun MapScreen(
    onNavigateToLogin: () -> Unit,
    navBar: @Composable () -> Unit,
    modifier: Modifier = Modifier,
    viewModel: MapViewModel = hiltViewModel<MapViewModel>(),
) {
    val loggedIn by viewModel.loggedIn.collectAsState(true)
    val latestOnNavigateToLogin by rememberUpdatedState(onNavigateToLogin)

    if (!loggedIn) {
        LaunchedEffect(Unit) { latestOnNavigateToLogin() }
    } else {
        MapScreen(map = { MapLibreMap(modifier = it) }, navBar = navBar, modifier = modifier)
    }
}

@Composable
fun MapScreen(
    map: @Composable ((modifier: Modifier) -> Unit),
    navBar: @Composable (() -> Unit),
    modifier: Modifier = Modifier,
) {
    Scaffold(modifier = modifier, bottomBar = navBar) { innerPadding ->
        Surface(modifier = Modifier.padding(innerPadding)) { map(Modifier) }
    }
}

@Composable
fun MapLibreMap(modifier: Modifier = Modifier) {
    val cameraState = CameraState.Centered(latitude = DefaultLatitude, longitude = DefaultLongitude)
    val mapViewCamera = rememberSaveableMapViewCamera(MapViewCamera(state = cameraState))

    MapView(
        modifier = modifier,
        styleUrl = LocalMapLibreStyleProvider.current.getStyleUrl(),
        camera = mapViewCamera,
    )
}

@Suppress("detekt:UnusedPrivateMember")
@PreviewAll
@Composable
private fun MapScreenPreview() {
    ParkEasyTheme { MapScreen(map = {}, navBar = { NavBar() }) }
}
