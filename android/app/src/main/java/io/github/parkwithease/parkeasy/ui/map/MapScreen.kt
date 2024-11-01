package io.github.parkwithease.parkeasy.ui.map

import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.Surface
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.hilt.navigation.compose.hiltViewModel
import com.maplibre.compose.LocalMapLibreStyleProvider
import com.maplibre.compose.MapView
import com.maplibre.compose.camera.CameraState
import com.maplibre.compose.camera.MapViewCamera
import com.maplibre.compose.rememberSaveableMapViewCamera

@Composable
fun MapScreen(
    modifier: Modifier = Modifier,
    @Suppress("unused") viewModel: MapViewModel = hiltViewModel<MapViewModel>(),
) {
    val mapViewCamera =
        rememberSaveableMapViewCamera(
            // Somewhere above Winnipeg
            MapViewCamera(state = CameraState.Centered(latitude = 49.9, longitude = -97.1))
        )
    Surface(modifier = modifier) {
        MapView(
            modifier = Modifier.fillMaxSize(),
            styleUrl = LocalMapLibreStyleProvider.current.getStyleUrl(),
            camera = mapViewCamera,
        )
    }
}
