package io.github.parkwithease.parkeasy.ui.spots

import androidx.compose.foundation.Image
import androidx.compose.foundation.gestures.detectDragGestures
import androidx.compose.foundation.interaction.MutableInteractionSource
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.heightIn
import androidx.compose.foundation.layout.imePadding
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.widthIn
import androidx.compose.foundation.lazy.grid.GridCells
import androidx.compose.foundation.lazy.grid.LazyGridState
import androidx.compose.foundation.lazy.grid.LazyVerticalGrid
import androidx.compose.foundation.lazy.grid.items
import androidx.compose.foundation.lazy.grid.rememberLazyGridState
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.selection.toggleable
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Add
import androidx.compose.material3.Button
import androidx.compose.material3.Card
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.FloatingActionButton
import androidx.compose.material3.Icon
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.ModalBottomSheet
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Scaffold
import androidx.compose.material3.SnackbarHost
import androidx.compose.material3.SnackbarHostState
import androidx.compose.material3.Surface
import androidx.compose.material3.Switch
import androidx.compose.material3.Text
import androidx.compose.material3.rememberModalBottomSheetState
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.saveable.rememberSaveable
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.geometry.Offset
import androidx.compose.ui.input.pointer.pointerInput
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.round
import androidx.compose.ui.unit.toIntRect
import androidx.hilt.navigation.compose.hiltViewModel
import io.github.parkwithease.parkeasy.R
import io.github.parkwithease.parkeasy.model.EditMode
import io.github.parkwithease.parkeasy.model.Spot
import io.github.parkwithease.parkeasy.ui.common.PullToRefreshBox

private const val NumColumns = 7
private const val NumRows = 48

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
            viewModel.onCountryCodeChange("")
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
                    viewModel::onPlusTime,
                    viewModel::onMinusTime,
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
    plus: (elements: Iterable<Int>) -> Unit,
    minus: (elements: Iterable<Int>) -> Unit,
    modifier: Modifier = Modifier,
) {
    Column(
        verticalArrangement = Arrangement.spacedBy(2.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
        modifier = modifier.widthIn(max = 320.dp),
    ) {
        OutlinedTextField(
            value = state.streetAddress.value,
            onValueChange = onStreetAddressChange,
            label = { Text(stringResource(R.string.street_address)) },
            modifier = Modifier.fillMaxWidth(),
            singleLine = true,
            isError = state.streetAddress.error != null,
            supportingText = { state.streetAddress.error?.also { Text(it) } },
        )
        OutlinedTextField(
            value = state.city.value,
            onValueChange = onCityChange,
            label = { Text(stringResource(R.string.city)) },
            modifier = Modifier.fillMaxWidth(),
            singleLine = true,
            isError = state.city.error != null,
            supportingText = { state.city.error?.also { Text(it) } },
        )
        OutlinedTextField(
            value = state.state.value,
            onValueChange = onStateChange,
            label = { Text(stringResource(R.string.state)) },
            modifier = Modifier.fillMaxWidth(),
            singleLine = true,
            isError = state.state.error != null,
            supportingText = { state.state.error?.also { Text(it) } },
        )
        OutlinedTextField(
            value = state.countryCode.value,
            onValueChange = onCountryCodeChange,
            label = { Text(stringResource(R.string.country)) },
            modifier = Modifier.fillMaxWidth(),
            singleLine = true,
            isError = state.countryCode.error != null,
            supportingText = { state.countryCode.error?.also { Text(it) } },
        )
        OutlinedTextField(
            value = state.postalCode.value,
            onValueChange = onPostalCodeChange,
            label = { Text(stringResource(R.string.postal_code)) },
            modifier = Modifier.fillMaxWidth(),
            singleLine = true,
            isError = state.postalCode.error != null,
            supportingText = { state.postalCode.error?.also { Text(it) } },
        )
        OutlinedTextField(
            value = state.pricePerHour.value,
            onValueChange = onPricePerHourChange,
            label = { Text(stringResource(R.string.price_per_hour)) },
            modifier = Modifier.fillMaxWidth(),
            singleLine = true,
            isError = state.pricePerHour.error != null,
            supportingText = { state.pricePerHour.error?.also { Text(it) } },
            keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Decimal),
        )
        Row {
            Column(modifier = Modifier.weight(1f)) {
                Text(stringResource(R.string.charging_station))
            }
            Column(horizontalAlignment = Alignment.End, modifier = Modifier.weight(1f)) {
                Switch(
                    checked = state.chargingStation.value,
                    onCheckedChange = onChargingStationChange,
                )
            }
        }
        Row {
            Column(modifier = Modifier.weight(1f)) { Text(stringResource(R.string.plug_in)) }
            Column(horizontalAlignment = Alignment.End, modifier = Modifier.weight(1f)) {
                Switch(checked = state.plugIn.value, onCheckedChange = onPlugInChange)
            }
        }
        Row {
            Column(modifier = Modifier.weight(1f)) { Text(stringResource(R.string.shelter)) }
            Column(horizontalAlignment = Alignment.End, modifier = Modifier.weight(1f)) {
                Switch(checked = state.shelter.value, onCheckedChange = onShelterChange)
            }
        }
        TimeGrid(state.times.value, plus, minus, Modifier.height(576.dp))
        Button(onClick = onAddSpotClick, modifier = Modifier.fillMaxWidth()) {
            Text(stringResource(R.string.add_spot))
        }
    }
}

@Composable
private fun TimeGrid(
    selectedIds: Set<Int>,
    plus: (elements: Iterable<Int>) -> Unit,
    minus: (elements: Iterable<Int>) -> Unit,
    modifier: Modifier = Modifier,
    slots: List<Int> = List(NumColumns * NumRows) { it % NumColumns * NumRows + it / NumColumns },
) {
    val state = rememberLazyGridState()

    LazyVerticalGrid(
        state = state,
        columns = GridCells.Fixed(NumColumns),
        horizontalArrangement = Arrangement.spacedBy(3.dp),
        modifier =
            modifier.timeGridDragHandler(
                lazyGridState = state,
                selectedIds = selectedIds,
                plus = plus,
                minus = minus,
            ),
    ) {
        items(slots, key = { it }) { id ->
            val selected = selectedIds.contains(id)

            Surface(
                tonalElevation = 3.dp,
                color =
                    if (selected) MaterialTheme.colorScheme.onPrimary
                    else MaterialTheme.colorScheme.onSurface,
                modifier =
                    Modifier.height(12.dp)
                        .padding(top = if (id % 48 > 0 && id % 48 % 2 == 0) 3.dp else 0.dp)
                        .toggleable(
                            value = selected,
                            interactionSource = remember { MutableInteractionSource() },
                            indication = null, // do not show a ripple
                            onValueChange = {
                                if (it) {
                                    plus(id..id)
                                } else {
                                    minus(id..id)
                                }
                            },
                        ),
            ) {}
        }
    }
}

@Suppress("detekt:UnsafeCallOnNullableType") // code provided by a Google engineer -> probably fine
fun Modifier.timeGridDragHandler(
    lazyGridState: LazyGridState,
    selectedIds: Set<Int>,
    plus: (elements: Iterable<Int>) -> Unit,
    minus: (elements: Iterable<Int>) -> Unit,
) =
    pointerInput(Unit) {
        fun LazyGridState.gridItemKeyAtPosition(hitPoint: Offset): Int? =
            layoutInfo.visibleItemsInfo
                .find { itemInfo ->
                    itemInfo.size.toIntRect().contains(hitPoint.round() - itemInfo.offset)
                }
                ?.key as? Int

        var initialKey: Int? = null
        var currentKey: Int? = null
        var adding = false
        detectDragGestures(
            onDragStart = { offset ->
                lazyGridState.gridItemKeyAtPosition(offset)?.let { key ->
                    initialKey = key
                    currentKey = key
                    if (!selectedIds.contains(key)) {
                        plus(key..key)
                        adding = true
                    } else {
                        minus(key..key)
                        adding = false
                    }
                }
            },
            onDragCancel = { initialKey = null },
            onDragEnd = { initialKey = null },
            onDrag = { change, _ ->
                if (initialKey != null) {
                    lazyGridState.gridItemKeyAtPosition(change.position)?.let { key ->
                        if (currentKey != key) {
                            if (adding) {
                                minus(initialKey!!..currentKey!!)
                                minus(currentKey!!..initialKey!!)
                                plus(initialKey!!..key)
                                plus(key..initialKey!!)
                            } else {
                                plus(initialKey!!..currentKey!!)
                                plus(currentKey!!..initialKey!!)
                                minus(initialKey!!..key)
                                minus(key..initialKey!!)
                            }
                            currentKey = key
                        }
                    }
                }
            },
        )
    }
