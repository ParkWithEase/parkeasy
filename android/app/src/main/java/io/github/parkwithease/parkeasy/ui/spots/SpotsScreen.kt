@file:Suppress("detekt:TooManyFunctions")

package io.github.parkwithease.parkeasy.ui.spots

import androidx.compose.foundation.Image
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.heightIn
import androidx.compose.foundation.layout.imePadding
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.widthIn
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Add
import androidx.compose.material3.Button
import androidx.compose.material3.Card
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.FilterChip
import androidx.compose.material3.FloatingActionButton
import androidx.compose.material3.Icon
import androidx.compose.material3.ModalBottomSheet
import androidx.compose.material3.Scaffold
import androidx.compose.material3.SheetValue
import androidx.compose.material3.SnackbarHost
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.material3.rememberModalBottomSheetState
import androidx.compose.material3.rememberStandardBottomSheetState
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.saveable.rememberSaveable
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.tooling.preview.PreviewParameter
import androidx.compose.ui.tooling.preview.PreviewParameterProvider
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import io.github.parkwithease.parkeasy.R
import io.github.parkwithease.parkeasy.model.EditMode
import io.github.parkwithease.parkeasy.model.FieldState
import io.github.parkwithease.parkeasy.model.Spot
import io.github.parkwithease.parkeasy.ui.common.ParkEasyTextField
import io.github.parkwithease.parkeasy.ui.common.PreviewAll
import io.github.parkwithease.parkeasy.ui.common.PullToRefreshBox
import io.github.parkwithease.parkeasy.ui.common.TimeGrid
import io.github.parkwithease.parkeasy.ui.theme.ParkEasyTheme

const val NumColumns = 6
const val NumRows = 24
const val NumSlots = NumColumns * NumRows

@Suppress("detekt:LongMethod")
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun SpotsScreen(modifier: Modifier = Modifier, viewModel: SpotsViewModel = hiltViewModel()) {
    val handler = rememberAddSpotFormHandler(viewModel)
    val spots by viewModel.spots.collectAsState()
    val isRefreshing by viewModel.isRefreshing.collectAsState()
    var editMode by rememberSaveable { mutableStateOf(EditMode.ADD) }

    var openBottomSheet by rememberSaveable { mutableStateOf(false) }
    val skipPartiallyExpanded by rememberSaveable { mutableStateOf(true) }
    val bottomSheetState =
        rememberModalBottomSheetState(skipPartiallyExpanded = skipPartiallyExpanded)

    LaunchedEffect(Unit) { viewModel.onRefresh() }
    SpotsScreen(
        spots = spots,
        onShowAddSpotClick = {
            handler.resetForm()
            editMode = EditMode.ADD
            openBottomSheet = true
        },
        isRefreshing = isRefreshing,
        onRefresh = viewModel::onRefresh,
        snackbarHost = { SnackbarHost(hostState = viewModel.snackbarState) },
        modifier = modifier,
    )
    if (openBottomSheet) {
        ModalBottomSheet(
            onDismissRequest = { openBottomSheet = false },
            sheetState = bottomSheetState,
        ) {
            AddSpotScreen(
                state = viewModel.formState,
                handler = handler,
                getSelectedIds = { viewModel.formState.selectedIds.value },
                disabledIds = viewModel.formState.disabledIds.value,
            )
        }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun SpotsScreen(
    spots: List<Spot>,
    onShowAddSpotClick: () -> Unit,
    isRefreshing: Boolean,
    onRefresh: () -> Unit,
    snackbarHost: @Composable (() -> Unit),
    modifier: Modifier = Modifier,
) {
    Scaffold(
        modifier = modifier,
        snackbarHost = snackbarHost,
        floatingActionButton = { AddSpotButton(onShowAddSpotClick = onShowAddSpotClick) },
    ) { innerPadding ->
        Surface(Modifier.padding(innerPadding)) {
            PullToRefreshBox(
                items = spots,
                onClick = {},
                isRefreshing = isRefreshing,
                onRefresh = onRefresh,
                modifier = Modifier.padding(4.dp),
            ) { spot, onClick ->
                SpotCard(spot, onClick)
            }
        }
    }
}

@Composable
fun SpotCard(spot: Spot, onClick: (Spot) -> Unit, modifier: Modifier = Modifier) {
    Card(onClick = { onClick(spot) }, modifier = modifier.fillMaxWidth().padding(4.dp, 0.dp)) {
        Row(modifier = Modifier.padding(8.dp)) {
            Column(modifier = Modifier.weight(1f)) {
                Image(
                    painter = painterResource(R.drawable.wordmark),
                    contentDescription = null,
                    modifier = Modifier.heightIn(max = 64.dp),
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

@Composable
fun AddSpotButton(onShowAddSpotClick: () -> Unit, modifier: Modifier = Modifier) {
    FloatingActionButton(onClick = onShowAddSpotClick, modifier) {
        Icon(imageVector = Icons.Filled.Add, contentDescription = stringResource(R.string.add_spot))
    }
}

@Suppress("detekt:LongMethod")
@Composable
fun AddSpotScreen(
    state: AddSpotFormState,
    handler: AddSpotFormHandler,
    getSelectedIds: () -> Set<Int>,
    disabledIds: Set<Int>,
    modifier: Modifier = Modifier,
) {
    Column(
        modifier =
            modifier
                .fillMaxWidth()
                .imePadding()
                .padding(horizontal = 16.dp)
                .verticalScroll(rememberScrollState(), reverseScrolling = true)
                .widthIn(max = 360.dp),
        verticalArrangement = Arrangement.spacedBy(2.dp, Alignment.CenterVertically),
        horizontalAlignment = Alignment.CenterHorizontally,
    ) {
        ParkEasyTextField(
            state = state.streetAddress,
            onValueChange = handler.onStreetAddressChange,
            modifier = Modifier.fillMaxWidth(),
            labelId = R.string.street_address,
        )
        Row(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
            ParkEasyTextField(
                state = state.city,
                onValueChange = handler.onCityChange,
                modifier = Modifier.weight(1f),
                labelId = R.string.city,
            )
            ParkEasyTextField(
                state = state.state,
                onValueChange = handler.onStateChange,
                modifier = Modifier.weight(1f),
                labelId = R.string.state,
            )
        }
        Row(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
            ParkEasyTextField(
                state = state.countryCode,
                onValueChange = handler.onCountryCodeChange,
                modifier = Modifier.weight(1f),
                labelId = R.string.country,
            )
            ParkEasyTextField(
                state = state.postalCode,
                onValueChange = handler.onPostalCodeChange,
                modifier = Modifier.weight(1f),
                labelId = R.string.postal_code,
            )
        }
        ParkEasyTextField(
            state = state.pricePerHour,
            onValueChange = handler.onPricePerHourChange,
            modifier = Modifier.fillMaxWidth(),
            labelId = R.string.price_per_hour,
            keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Decimal),
        )
        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.spacedBy(4.dp),
        ) {
            FilterChip(
                selected = state.chargingStation.value,
                onClick = { handler.onChargingStationChange(!state.chargingStation.value) },
                label = { Text(stringResource(R.string.charging_station)) },
            )
            FilterChip(
                selected = state.plugIn.value,
                onClick = { handler.onPlugInChange(!state.plugIn.value) },
                label = { Text(stringResource(R.string.plug_in)) },
            )
            FilterChip(
                selected = state.shelter.value,
                onClick = { handler.onShelterChange(!state.shelter.value) },
                label = { Text(stringResource(R.string.shelter)) },
            )
        }
        TimeGrid(getSelectedIds, disabledIds, handler.onAddTime, handler.onRemoveTime)
        Button(onClick = handler.onAddSpotClick, modifier = Modifier.fillMaxWidth()) {
            Text(stringResource(R.string.add_spot))
        }
    }
}

@Suppress("detekt:UnusedPrivateMember")
@OptIn(ExperimentalMaterial3Api::class)
@PreviewAll
@Composable
private fun SpotsScreenPreview() {
    val spots = listOf(Spot())
    ParkEasyTheme {
        SpotsScreen(
            spots = spots,
            onShowAddSpotClick = {},
            isRefreshing = false,
            onRefresh = {},
            snackbarHost = {},
        )
    }
}

private class AddSpotFormStateProvider : PreviewParameterProvider<AddSpotFormState> {
    override val values =
        sequenceOf(
            AddSpotFormState(),
            AddSpotFormState(
                streetAddress = FieldState("66 Chancellors Cir"),
                city = FieldState("Winnipeg"),
                state = FieldState("MB"),
                countryCode = FieldState("CA"),
                postalCode = FieldState("R3T2N2"),
                chargingStation = FieldState(true),
                plugIn = FieldState(true),
                shelter = FieldState(false),
                pricePerHour = FieldState("4.49"),
            ),
            AddSpotFormState(
                streetAddress = FieldState("", "Address cannot be empty"),
                city = FieldState("", "City cannot be empty"),
                state = FieldState("", "State cannot be empty"),
                countryCode = FieldState("", "Country cannot be empty"),
                postalCode = FieldState("", "Postal code cannot be empty"),
                chargingStation = FieldState(false),
                plugIn = FieldState(false),
                shelter = FieldState(false),
                pricePerHour = FieldState("", "Price cannot be empty"),
            ),
        )
}

@Suppress("detekt:UnusedPrivateMember")
@OptIn(ExperimentalMaterial3Api::class)
@PreviewAll
@Composable
private fun AddCarScreenPreview(
    @PreviewParameter(AddSpotFormStateProvider::class) state: AddSpotFormState
) {
    ParkEasyTheme {
        ModalBottomSheet(
            onDismissRequest = {},
            sheetState = rememberStandardBottomSheetState(initialValue = SheetValue.Expanded),
        ) {
            AddSpotScreen(
                state = state,
                handler = AddSpotFormHandler(),
                getSelectedIds = { emptySet() },
                disabledIds = emptySet(),
            )
        }
    }
}

@Suppress("detekt:UnusedPrivateMember")
@OptIn(ExperimentalMaterial3Api::class)
@PreviewAll
@Composable
private fun ExpandedAddCarScreenPreview(
    @PreviewParameter(AddSpotFormStateProvider::class) state: AddSpotFormState
) {
    ParkEasyTheme {
        Surface(
            Modifier.fillMaxSize()
        ) { // Can't actually preview fully expanded BottomSheets apparently
            AddSpotScreen(
                state = state,
                handler = AddSpotFormHandler(),
                getSelectedIds = { emptySet() },
                disabledIds = emptySet(),
            )
        }
    }
}

private class AvailabilityProvider : PreviewParameterProvider<Set<Int>> {
    override val values =
        sequenceOf(
            (0..72).toSet(),
            (72..144).toSet(),
            (49..94).toSet(),
            (0..72).map { it * 2 }.toSet(),
            (0..72).map { it * 2 + 1 }.toSet(),
        )
}

@Suppress("detekt:UnusedPrivateMember")
@OptIn(ExperimentalMaterial3Api::class)
@PreviewAll
@Composable
private fun AddCarScreenTimeGridPreview(
    @PreviewParameter(AvailabilityProvider::class) selectedIds: Set<Int>
) {
    ParkEasyTheme {
        Surface(Modifier.fillMaxSize()) {
            AddSpotScreen(
                state = AddSpotFormState(),
                handler = AddSpotFormHandler(),
                getSelectedIds = { selectedIds },
                disabledIds = emptySet(),
            )
        }
    }
}
