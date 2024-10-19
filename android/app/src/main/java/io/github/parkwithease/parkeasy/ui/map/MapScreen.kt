package io.github.parkwithease.parkeasy.ui.map

import androidx.compose.material3.Surface
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.hilt.navigation.compose.hiltViewModel

@Composable
fun MapScreen(
    modifier: Modifier = Modifier,
    @Suppress("unused") viewModel: MapViewModel = hiltViewModel<MapViewModel>(),
) {
    Surface(modifier) {}
}
