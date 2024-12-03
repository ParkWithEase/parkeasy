package io.github.parkwithease.parkeasy.ui.spots

import android.util.Log
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
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.layout.widthIn
import androidx.compose.foundation.lazy.grid.GridCells
import androidx.compose.foundation.lazy.grid.GridItemSpan
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
import androidx.compose.material3.FilterChip
import androidx.compose.material3.FloatingActionButton
import androidx.compose.material3.Icon
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.ModalBottomSheet
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
import io.github.parkwithease.parkeasy.ui.common.MinutesPerSlot
import io.github.parkwithease.parkeasy.ui.common.ParkEasyTextField
import io.github.parkwithease.parkeasy.ui.common.PullToRefreshBox
import io.github.parkwithease.parkeasy.ui.common.isoDay
import io.github.parkwithease.parkeasy.ui.common.startOfNextAvailableDay
import io.github.parkwithease.parkeasy.ui.common.timezone
import io.github.parkwithease.parkeasy.ui.common.toShortDate
import kotlin.time.DurationUnit
import kotlinx.datetime.Clock
import kotlinx.datetime.LocalTime
import kotlinx.datetime.toInstant

private const val NumColumns = 6
private const val NumRows = 24

@OptIn(ExperimentalMaterial3Api::class)
@Suppress("detekt:LongMethod")
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
        spots,
        {
            handler.resetForm()
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
                    handler,
                    viewModel::onAddSpotClick,
                    { viewModel.formState.times.value },
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
    handler: AddSpotFormHandler,
    onAddSpotClick: () -> Unit,
    getSelectedIds: () -> Set<Int>,
    modifier: Modifier = Modifier,
) {
    Column(
        verticalArrangement = Arrangement.spacedBy(2.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
        modifier = modifier.widthIn(max = 360.dp),
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
        Row(Modifier.height(336.dp)) {
            ColumnHeader(Modifier.width(48.dp))
            TimeGrid(getSelectedIds, handler.onPlusTime, handler.onMinusTime)
        }
        Button(onClick = onAddSpotClick, modifier = Modifier.fillMaxWidth()) {
            Text(stringResource(R.string.add_spot))
        }
    }
}

@Composable
private fun ColumnHeader(
    modifier: Modifier = Modifier,
    state: LazyGridState = rememberLazyGridState(),
) {
    LazyVerticalGrid(state = state, columns = GridCells.Fixed(1), modifier = modifier) {
        items(2) {
            Text(text = "", Modifier.height(24.dp), style = MaterialTheme.typography.labelLarge)
        }
        items(NumRows) { num ->
            if (num % 2 == 0) {
                Text(
                    text = LocalTime((num / 2), 0).toString(),
                    Modifier.height(24.dp),
                    style = MaterialTheme.typography.labelLarge,
                )
            }
        }
    }
}

@Suppress("detekt:LongMethod")
@Composable
private fun TimeGrid(
    getSelectedIds: () -> Set<Int>,
    plus: (elements: Iterable<Int>) -> Unit,
    minus: (elements: Iterable<Int>) -> Unit,
    modifier: Modifier = Modifier,
    state: LazyGridState = rememberLazyGridState(),
    slots: List<Int> = List(NumColumns * NumRows) { it % NumColumns * NumRows + it / NumColumns },
) {
    val disabledIds: List<Int> =
        if (Clock.System.now() > startOfNextAvailableDay().toInstant(timezone()))
            (0..Clock.System.now()
                        .minus(startOfNextAvailableDay().toInstant(timezone()))
                        .toInt(DurationUnit.MINUTES) / MinutesPerSlot)
                .asSequence()
                .toList()
        else emptyList()
    val selectedIds: Set<Int> = getSelectedIds()

    LazyVerticalGrid(
        state = state,
        columns = GridCells.Fixed(NumColumns),
        horizontalArrangement = Arrangement.spacedBy(3.dp),
        modifier =
            modifier.timeGridDragHandler(
                lazyGridState = state,
                disabledIds = disabledIds,
                getSelectedIds = getSelectedIds,
                plus = plus,
                minus = minus,
            ),
    ) {
        items(count = NumColumns / 2, span = { GridItemSpan(2) }) {
            Text(
                text = startOfNextAvailableDay().isoDay(it + 1).toShortDate(),
                Modifier.height(24.dp),
                style = MaterialTheme.typography.labelLarge,
            )
        }
        items(NumColumns) {
            Text(
                text = if (it % 2 == 0) "AM" else "PM",
                Modifier.height(24.dp),
                style = MaterialTheme.typography.labelLarge,
            )
        }
        items(slots, key = { it }) { id ->
            val disabled = disabledIds.contains(id)
            val selected = selectedIds.contains(id)

            Surface(
                tonalElevation = 3.dp,
                color =
                    if (disabled) {
                        MaterialTheme.colorScheme.onSurface
                    } else {
                        if (selected) MaterialTheme.colorScheme.primary
                        else MaterialTheme.colorScheme.primaryContainer
                    },
                modifier =
                    Modifier.height(12.dp)
                        .padding(top = if (id % 24 > 0 && id % 24 % 2 == 0) 3.dp else 0.dp)
                        .toggleable(
                            value = selected,
                            interactionSource = remember { MutableInteractionSource() },
                            indication = null, // do not show a ripple
                            onValueChange = {
                                Log.e("", selectedIds.toString())
                                if (!disabled) {
                                    if (it) {
                                        plus(id..id)
                                    } else {
                                        minus(id..id)
                                    }
                                }
                            },
                        ),
            ) {}
        }
    }
}

@Suppress("detekt:UnsafeCallOnNullableType") // code provided by a Google engineer -> probably fine
private fun Modifier.timeGridDragHandler(
    lazyGridState: LazyGridState,
    disabledIds: List<Int>,
    getSelectedIds: () -> Set<Int>,
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
                    val selectedIds = getSelectedIds()
                    if (!disabledIds.contains(key)) {
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
                            minus(disabledIds.min()..disabledIds.max())
                            currentKey = key
                        }
                    }
                }
            },
        )
    }
