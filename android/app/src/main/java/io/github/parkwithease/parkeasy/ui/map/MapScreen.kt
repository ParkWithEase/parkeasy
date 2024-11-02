package io.github.parkwithease.parkeasy.ui.map

import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Surface
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.navigation.NavHostController
import com.maplibre.compose.LocalMapLibreStyleProvider
import com.maplibre.compose.MapView
import com.maplibre.compose.camera.CameraState
import com.maplibre.compose.camera.MapViewCamera
import com.maplibre.compose.rememberSaveableMapViewCamera
import io.github.parkwithease.parkeasy.ui.navbar.NavBar

@Composable
fun MapScreen(
    navController: NavHostController,
    modifier: Modifier = Modifier,
    @Suppress("unused") viewModel: MapViewModel = hiltViewModel<MapViewModel>(),
) {
    val mapViewCamera =
        rememberSaveableMapViewCamera(
            // Somewhere above Winnipeg
            MapViewCamera(state = CameraState.Centered(latitude = 49.9, longitude = -97.1))
        )
    Scaffold(modifier = modifier, bottomBar = { NavBar(navController = navController) }) {
        innerPadding ->
        Surface(modifier = Modifier.padding(innerPadding)) {
            MapView(
                modifier = Modifier.fillMaxSize(),
                styleUrl = LocalMapLibreStyleProvider.current.getStyleUrl(),
                camera = mapViewCamera,
            )
        }
    }
}
