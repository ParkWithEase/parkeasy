package io.github.parkwithease.parkeasy.ui.map

import androidx.compose.foundation.layout.fillMaxSize
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
import androidx.navigation.NavHostController
import com.maplibre.compose.LocalMapLibreStyleProvider
import com.maplibre.compose.MapView
import com.maplibre.compose.camera.CameraState
import com.maplibre.compose.camera.MapViewCamera
import com.maplibre.compose.rememberSaveableMapViewCamera
import io.github.parkwithease.parkeasy.ui.navbar.NavBar

@Composable
fun MapScreen(
    onNavigateToLogin: () -> Unit,
    navController: NavHostController,
    modifier: Modifier = Modifier,
    viewModel: MapViewModel = hiltViewModel<MapViewModel>(),
) {
    val latestOnNavigateToLogin by rememberUpdatedState(onNavigateToLogin)
    val loggedIn by viewModel.loggedIn.collectAsState(true)

    if (!loggedIn) {
        LaunchedEffect(Unit) { latestOnNavigateToLogin() }
    } else {
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
}
