package io.github.parkwithease.parkeasy.ui.spots

import androidx.compose.foundation.Image
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
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
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Scaffold
import androidx.compose.material3.SnackbarHost
import androidx.compose.material3.SnackbarHostState
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.material3.rememberModalBottomSheetState
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
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import io.github.parkwithease.parkeasy.R
import io.github.parkwithease.parkeasy.model.EditMode
import io.github.parkwithease.parkeasy.model.Spot
import io.github.parkwithease.parkeasy.ui.common.PullToRefreshBox
import io.github.parkwithease.parkeasy.ui.common.textOrNull

@OptIn(ExperimentalMaterial3Api::class)
@Suppress("detekt:LongMethod")
@Composable
fun SpotsScreen(modifier: Modifier = Modifier, viewModel: SpotsViewModel = hiltViewModel()) {
    val spots by viewModel.spots.collectAsState()
    val isRefreshing by viewModel.isRefreshing.collectAsState()
    var editMode by rememberSaveable { mutableStateOf(EditMode.ADD) }

    var openBottomSheet by rememberSaveable { mutableStateOf(false) }
    val skipPartiallyExpanded by rememberSaveable { mutableStateOf(true) }
    val bottomSheetState =
        rememberModalBottomSheetState(skipPartiallyExpanded = skipPartiallyExpanded)

    LaunchedEffect(Unit) { viewModel.onRefresh() }
    SpotsScreen(
        spots,
        {
            viewModel.onStreetAddressChange("")
            viewModel.onCityChange("")
            viewModel.onStateChange("")
            viewModel.onCountryCodeChange("CA")
            viewModel.onPostalCodeChange("")
            viewModel.onChargingStationChange(false)
            viewModel.onPlugInChange(false)
            viewModel.onShelterChange(false)
            viewModel.onPricePerHourChange("")
            editMode = EditMode.ADD
            openBottomSheet = true
        },
        isRefreshing,
        viewModel::onRefresh,
        viewModel.snackbarState,
        modifier,
    )
    if (openBottomSheet) {
        ModalBottomSheet(
            onDismissRequest = { openBottomSheet = false },
            sheetState = bottomSheetState,
        ) {
            Column(
                horizontalAlignment = Alignment.CenterHorizontally,
                verticalArrangement = Arrangement.Center,
                modifier =
                    Modifier.padding(horizontal = 16.dp)
                        .fillMaxWidth()
                        .imePadding()
                        .verticalScroll(rememberScrollState(), reverseScrolling = true),
            ) {
                AddSpotScreen(
                    viewModel.formState,
                    viewModel::onStreetAddressChange,
                    viewModel::onCityChange,
                    viewModel::onStateChange,
                    viewModel::onCountryCodeChange,
                    viewModel::onPostalCodeChange,
                    viewModel::onChargingStationChange,
                    viewModel::onPlugInChange,
                    viewModel::onShelterChange,
                    viewModel::onPricePerHourChange,
                    viewModel::onAddSpotClick,
                )
            }
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
    snackbarState: SnackbarHostState,
    modifier: Modifier = Modifier,
) {
    Scaffold(
        floatingActionButton = { AddSpotButton(onShowAddSpotClick = onShowAddSpotClick) },
        modifier = modifier,
        snackbarHost = { SnackbarHost(hostState = snackbarState) },
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
                    painter = painterResource(R.drawable.wordmark_outlined),
                    contentDescription = null,
                    modifier = Modifier.heightIn(max = 64.dp),
                )
            }
            Column(horizontalAlignment = Alignment.End, modifier = Modifier.weight(1f)) {
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

@Composable
@Suppress("detekt:LongMethod")
fun AddSpotScreen(
    state: AddSpotFormState,
    onStreetAddressChange: (String) -> Unit,
    onCityChange: (String) -> Unit,
    onStateChange: (String) -> Unit,
    onCountryCodeChange: (String) -> Unit,
    onPostalCodeChange: (String) -> Unit,
    onChargingStationChange: (Boolean) -> Unit,
    onPlugInChange: (Boolean) -> Unit,
    onShelterChange: (Boolean) -> Unit,
    onPricePerHourChange: (String) -> Unit,
    onAddSpotClick: () -> Unit,
    modifier: Modifier = Modifier,
) {
    Column(
        verticalArrangement = Arrangement.spacedBy(2.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
        modifier = modifier.widthIn(max = 360.dp),
    ) {
        OutlinedTextField(
            value = state.streetAddress.value,
            onValueChange = onStreetAddressChange,
            label = { Text(stringResource(R.string.street_address)) },
            modifier = Modifier.fillMaxWidth(),
            singleLine = true,
            isError = state.streetAddress.error != null,
            supportingText = state.streetAddress.error.textOrNull(),
        )
        Row(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
            OutlinedTextField(
                value = state.city.value,
                onValueChange = onCityChange,
                label = { Text(stringResource(R.string.city)) },
                modifier = Modifier.weight(1f),
                singleLine = true,
                isError = state.city.error != null,
                supportingText = state.city.error.textOrNull(),
            )
            OutlinedTextField(
                value = state.state.value,
                onValueChange = onStateChange,
                label = { Text(stringResource(R.string.state)) },
                modifier = Modifier.weight(1f),
                singleLine = true,
                isError = state.state.error != null,
                supportingText = state.state.error.textOrNull(),
            )
        }
        Row(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
            OutlinedTextField(
                value = state.countryCode.value,
                onValueChange = onCountryCodeChange,
                label = { Text(stringResource(R.string.country)) },
                modifier = Modifier.weight(1f),
                singleLine = true,
                isError = state.countryCode.error != null,
                supportingText = state.countryCode.error.textOrNull(),
            )
            OutlinedTextField(
                value = state.postalCode.value,
                onValueChange = onPostalCodeChange,
                label = { Text(stringResource(R.string.postal_code)) },
                modifier = Modifier.weight(1f),
                singleLine = true,
                isError = state.postalCode.error != null,
                supportingText = state.postalCode.error.textOrNull(),
            )
        }
        OutlinedTextField(
            value = state.pricePerHour.value,
            onValueChange = onPricePerHourChange,
            label = { Text(stringResource(R.string.price_per_hour)) },
            modifier = Modifier.fillMaxWidth(),
            singleLine = true,
            isError = state.pricePerHour.error != null,
            supportingText = state.pricePerHour.error.textOrNull(),
            keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Decimal),
        )
        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.spacedBy(4.dp),
        ) {
            FilterChip(
                selected = state.chargingStation.value,
                onClick = { onChargingStationChange(!state.chargingStation.value) },
                label = { Text(stringResource(R.string.charging_station)) },
            )
            FilterChip(
                selected = state.plugIn.value,
                onClick = { onPlugInChange(!state.plugIn.value) },
                label = { Text(stringResource(R.string.plug_in)) },
            )
            FilterChip(
                selected = state.shelter.value,
                onClick = { onShelterChange(!state.shelter.value) },
                label = { Text(stringResource(R.string.shelter)) },
            )
        }
        Button(onClick = onAddSpotClick, modifier = Modifier.fillMaxWidth()) {
            Text(stringResource(R.string.add_spot))
        }
    }
}
