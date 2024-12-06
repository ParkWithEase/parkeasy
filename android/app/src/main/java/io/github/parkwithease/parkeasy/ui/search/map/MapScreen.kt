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
import androidx.compose.runtime.rememberUpdatedState
import androidx.compose.ui.Modifier
import androidx.hilt.navigation.compose.hiltViewModel
import com.mapbox.mapboxsdk.geometry.LatLng
import com.maplibre.compose.LocalMapLibreStyleProvider
import com.maplibre.compose.MapView
import com.maplibre.compose.camera.CameraState
import com.maplibre.compose.camera.MapViewCamera
import com.maplibre.compose.rememberSaveableMapViewCamera
import com.maplibre.compose.symbols.Symbol
import io.github.parkwithease.parkeasy.R
import io.github.parkwithease.parkeasy.model.Spot
import io.github.parkwithease.parkeasy.ui.common.PreviewAll
import io.github.parkwithease.parkeasy.ui.navbar.NavBar
import io.github.parkwithease.parkeasy.ui.search.CreateBookingScreen
import io.github.parkwithease.parkeasy.ui.search.DefaultLatitude
import io.github.parkwithease.parkeasy.ui.search.DefaultLongitude
import io.github.parkwithease.parkeasy.ui.search.SearchViewModel
import io.github.parkwithease.parkeasy.ui.search.rememberCreateHandler
import io.github.parkwithease.parkeasy.ui.search.rememberSearchHandler
import io.github.parkwithease.parkeasy.ui.theme.ParkEasyTheme

@Composable
fun MapScreen(
    onNavigateToLogin: () -> Unit,
    navBar: @Composable () -> Unit,
    modifier: Modifier = Modifier,
    viewModel: SearchViewModel = hiltViewModel<SearchViewModel>(),
) {
    val loggedIn by viewModel.loggedIn.collectAsState(true)
    val latestOnNavigateToLogin by rememberUpdatedState(onNavigateToLogin)

    if (!loggedIn) {
        LaunchedEffect(Unit) { latestOnNavigateToLogin() }
    } else {
        @Suppress("unused") val searchHandler = rememberSearchHandler(viewModel)
        val createHandler = rememberCreateHandler(viewModel)
        val cars by viewModel.cars.collectAsState()
        val spots by viewModel.spots.collectAsState()
        val formEnabled by viewModel.formEnabled.collectAsState()
        val isRefreshing by viewModel.isRefreshing.collectAsState()
        val showForm by viewModel.showForm.collectAsState()

        BackHandler(enabled = showForm) { viewModel.onHideForm() }

        LaunchedEffect(Unit) { viewModel.onRefresh() }

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
            MapScreen(
                map = { mapSpots, onSpotClick, mapModifier ->
                    MapLibreMap(mapSpots, onSpotClick, mapModifier)
                },
                spots = spots,
                onSpotClick = {
                    createHandler.reset()
                    createHandler.onSpotChange(it)
                    viewModel.onShowForm()
                },
                isRefreshing = isRefreshing,
                onRefresh = viewModel::onRefresh,
                navBar = navBar,
                snackbarHost = { SnackbarHost(hostState = viewModel.snackbarState) },
                modifier = modifier,
            )
    }
}

@Composable
fun MapScreen(
    map: @Composable ((spots: List<Spot>, onSpotClick: (Spot) -> Unit, modifier: Modifier) -> Unit),
    spots: List<Spot>,
    onSpotClick: (Spot) -> Unit,
    isRefreshing: Boolean,
    onRefresh: () -> Unit,
    navBar: @Composable (() -> Unit),
    snackbarHost: @Composable (() -> Unit),
    modifier: Modifier = Modifier,
) {
    Scaffold(
        modifier = modifier,
        bottomBar = navBar,
        snackbarHost = snackbarHost,
        floatingActionButton = { RefreshButton(isRefreshing = isRefreshing, onRefresh = onRefresh) },
    ) { innerPadding ->
        Surface(modifier = Modifier.padding(innerPadding)) { map(spots, onSpotClick, modifier) }
    }
}

@Composable
fun MapLibreMap(spots: List<Spot>, onSpotClick: (Spot) -> Unit, modifier: Modifier = Modifier) {
    val cameraState = CameraState.Centered(latitude = DefaultLatitude, longitude = DefaultLongitude)
    val mapViewCamera = rememberSaveableMapViewCamera(MapViewCamera(state = cameraState))

    MapView(
        modifier = modifier,
        styleUrl = LocalMapLibreStyleProvider.current.getStyleUrl(),
        camera = mapViewCamera,
    ) {
        spots.forEach {
            Symbol(
                center = LatLng(it.location.latitude, it.location.longitude),
                imageId = R.drawable.location,
                onTap = { onSpotClick(it) },
            )
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

@Suppress("detekt:UnusedPrivateMember")
@PreviewAll
@Composable
private fun MapScreenPreview() {
    ParkEasyTheme {
        MapScreen(
            map = { _, _, _ -> },
            spots = emptyList<Spot>(),
            onSpotClick = {},
            isRefreshing = false,
            onRefresh = {},
            navBar = { NavBar() },
            snackbarHost = {},
        )
    }
}
