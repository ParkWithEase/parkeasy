package io.github.parkwithease.parkeasy.ui.search.list

import androidx.activity.compose.BackHandler
import androidx.compose.foundation.Image
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.imePadding
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.Button
import androidx.compose.material3.Card
import androidx.compose.material3.DropdownMenuItem
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.ExposedDropdownMenuBox
import androidx.compose.material3.ExposedDropdownMenuDefaults
import androidx.compose.material3.FilterChip
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.MenuAnchorType
import androidx.compose.material3.Scaffold
import androidx.compose.material3.SnackbarHost
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.rememberUpdatedState
import androidx.compose.runtime.saveable.rememberSaveable
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import io.github.parkwithease.parkeasy.R
import io.github.parkwithease.parkeasy.model.Car
import io.github.parkwithease.parkeasy.model.FieldState
import io.github.parkwithease.parkeasy.model.Spot
import io.github.parkwithease.parkeasy.model.SpotLocation
import io.github.parkwithease.parkeasy.ui.common.ParkEasyTextField
import io.github.parkwithease.parkeasy.ui.common.PreviewAll
import io.github.parkwithease.parkeasy.ui.common.PullToRefreshBox
import io.github.parkwithease.parkeasy.ui.common.TimeGrid
import io.github.parkwithease.parkeasy.ui.navbar.NavBar
import io.github.parkwithease.parkeasy.ui.search.CreateHandler
import io.github.parkwithease.parkeasy.ui.search.CreateState
import io.github.parkwithease.parkeasy.ui.search.SearchViewModel
import io.github.parkwithease.parkeasy.ui.search.rememberCreateHandler
import io.github.parkwithease.parkeasy.ui.search.rememberSearchHandler
import io.github.parkwithease.parkeasy.ui.theme.ParkEasyTheme

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun ListScreen(
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
        val isRefreshing by viewModel.isRefreshing.collectAsState()

        var showForm by rememberSaveable { mutableStateOf(false) }

        BackHandler(enabled = showForm) { showForm = false }

        LaunchedEffect(Unit) { viewModel.onRefresh() }

        if (showForm)
            CreateBookingScreen(
                cars = cars,
                state = viewModel.createState,
                handler = createHandler,
                getSelectedIds = { viewModel.createState.selectedIds.value },
                disabledIds = viewModel.createState.disabledIds.value,
            )
        else
            ListScreen(
                spots = spots,
                onSpotClick = {
                    createHandler.reset()
                    createHandler.onSpotChange(it)
                    showForm = true
                },
                isRefreshing = isRefreshing,
                onRefresh = viewModel::onRefresh,
                navBar = navBar,
                snackbarHost = { SnackbarHost(hostState = viewModel.snackbarState) },
                modifier = modifier,
            )
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
                    SpotCard(spot, onClick)
                }
            }
        }
    }
}

@Suppress("DefaultLocale", "detekt:ImplicitDefaultLocale")
@Composable
fun SpotCard(spot: Spot, onClick: (Spot) -> Unit, modifier: Modifier = Modifier) {
    Card(onClick = { onClick(spot) }, modifier = modifier.fillMaxWidth().padding(4.dp, 0.dp)) {
        Row(modifier = Modifier.padding(8.dp)) {
            Column(modifier = Modifier.weight(1f)) {
                Text(
                    text = String.format("$ %.2f", spot.pricePerHour),
                    style = MaterialTheme.typography.titleLarge,
                )
                Row(Modifier.height(24.dp)) {
                    if (spot.features.chargingStation)
                        Image(
                            painter = painterResource(R.drawable.charging_station),
                            contentDescription = null,
                            modifier = Modifier.width(24.dp),
                        )
                    if (spot.features.plugIn)
                        Image(
                            painter = painterResource(R.drawable.plug_in),
                            contentDescription = null,
                            modifier = Modifier.width(24.dp),
                        )
                    if (spot.features.shelter)
                        Image(
                            painter = painterResource(R.drawable.shelter),
                            contentDescription = null,
                            modifier = Modifier.width(24.dp),
                        )
                }
                Text(
                    text = String.format("%.0f meters away", spot.distanceToLocation),
                    style = MaterialTheme.typography.titleSmall,
                )
            }
            Column(modifier = Modifier.weight(1f), horizontalAlignment = Alignment.End) {
                Text(spot.location.streetAddress)
                Text(spot.location.city + ' ' + spot.location.state)
                Text(spot.location.countryCode + ' ' + spot.location.postalCode)
            }
        }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Suppress("detekt:LongMethod", "DefaultLocale", "detekt:ImplicitDefaultLocale")
@Composable
fun CreateBookingScreen(
    cars: List<Car>,
    state: CreateState,
    handler: CreateHandler,
    getSelectedIds: () -> Set<Int>,
    disabledIds: Set<Int>,
    modifier: Modifier = Modifier,
) {
    var expanded by remember { mutableStateOf(false) }

    Surface(modifier.fillMaxSize()) {
        Column(
            modifier =
                Modifier.imePadding()
                    .padding(32.dp)
                    .verticalScroll(rememberScrollState(), reverseScrolling = true),
            verticalArrangement = Arrangement.spacedBy(2.dp, Alignment.CenterVertically),
            horizontalAlignment = Alignment.CenterHorizontally,
        ) {
            ExposedDropdownMenuBox(expanded = expanded, onExpandedChange = { expanded = it }) {
                ParkEasyTextField(
                    state = FieldState(state.selectedCar.value.details.licensePlate),
                    onValueChange = {},
                    modifier =
                        Modifier.fillMaxWidth()
                            .menuAnchor(MenuAnchorType.PrimaryNotEditable)
                            .clickable { expanded = true },
                    enabled = false,
                    visuallyEnabled = true,
                    readOnly = true,
                    labelId = R.string.select_car,
                    trailingIcon = { ExposedDropdownMenuDefaults.TrailingIcon(expanded = expanded) },
                )
                ExposedDropdownMenu(expanded = expanded, onDismissRequest = { expanded = false }) {
                    cars.forEach {
                        DropdownMenuItem(
                            text = { Text(it.details.licensePlate) },
                            onClick = {
                                handler.onCarChange(it)
                                expanded = false
                            },
                        )
                    }
                }
            }
            ParkEasyTextField(
                state = FieldState(state.selectedSpot.value.location.streetAddress),
                onValueChange = {},
                modifier = Modifier.fillMaxWidth(),
                enabled = false,
                visuallyEnabled = true,
                readOnly = true,
                labelId = R.string.street_address,
            )
            Row(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                ParkEasyTextField(
                    state = FieldState(state.selectedSpot.value.location.city),
                    onValueChange = {},
                    modifier = Modifier.weight(1f),
                    enabled = false,
                    visuallyEnabled = true,
                    readOnly = true,
                    labelId = R.string.city,
                )
                ParkEasyTextField(
                    state = FieldState(state.selectedSpot.value.location.state),
                    onValueChange = {},
                    modifier = Modifier.weight(1f),
                    enabled = false,
                    visuallyEnabled = true,
                    readOnly = true,
                    labelId = R.string.state,
                )
            }
            Row(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                ParkEasyTextField(
                    state = FieldState(state.selectedSpot.value.location.countryCode),
                    onValueChange = {},
                    modifier = Modifier.weight(1f),
                    enabled = false,
                    visuallyEnabled = true,
                    readOnly = true,
                    labelId = R.string.country,
                )
                ParkEasyTextField(
                    state = FieldState(state.selectedSpot.value.location.postalCode),
                    onValueChange = {},
                    modifier = Modifier.weight(1f),
                    enabled = false,
                    visuallyEnabled = true,
                    readOnly = true,
                    labelId = R.string.postal_code,
                )
            }
            ParkEasyTextField(
                state = FieldState(String.format("$ %.2f", state.selectedSpot.value.pricePerHour)),
                onValueChange = {},
                modifier = Modifier.fillMaxWidth(),
                enabled = false,
                visuallyEnabled = true,
                readOnly = true,
                labelId = R.string.price_per_hour,
            )
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.spacedBy(4.dp),
            ) {
                FilterChip(
                    selected = state.selectedSpot.value.features.chargingStation,
                    onClick = {},
                    label = { Text(stringResource(R.string.charging_station)) },
                )
                FilterChip(
                    selected = state.selectedSpot.value.features.plugIn,
                    onClick = {},
                    label = { Text(stringResource(R.string.plug_in)) },
                )
                FilterChip(
                    selected = state.selectedSpot.value.features.plugIn,
                    onClick = {},
                    label = { Text(stringResource(R.string.shelter)) },
                )
            }
            TimeGrid(getSelectedIds, disabledIds, handler.onAddTime, handler.onRemoveTime)
            ParkEasyTextField(
                state =
                    FieldState(
                        String.format(
                            "$ %.2f",
                            state.selectedSpot.value.pricePerHour * getSelectedIds().size / 2,
                        )
                    ),
                onValueChange = {},
                modifier = Modifier.fillMaxWidth(),
                enabled = false,
                visuallyEnabled = true,
                readOnly = true,
                labelId = R.string.total_price,
            )
            Button(onClick = handler.onCreateBookingClick, modifier = Modifier.fillMaxWidth()) {
                Text(stringResource(R.string.create_booking))
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
