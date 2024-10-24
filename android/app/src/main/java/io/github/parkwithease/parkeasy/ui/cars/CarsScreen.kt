package io.github.parkwithease.parkeasy.ui.cars

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.outlined.Add
import androidx.compose.material.icons.outlined.Refresh
import androidx.compose.material3.Button
import androidx.compose.material3.Card
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.FloatingActionButton
import androidx.compose.material3.Icon
import androidx.compose.material3.ListItem
import androidx.compose.material3.ModalBottomSheet
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.material3.pulltorefresh.PullToRefreshBox
import androidx.compose.material3.pulltorefresh.PullToRefreshState
import androidx.compose.material3.pulltorefresh.rememberPullToRefreshState
import androidx.compose.material3.rememberModalBottomSheetState
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.rememberCoroutineScope
import androidx.compose.runtime.saveable.rememberSaveable
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import io.github.parkwithease.parkeasy.R
import io.github.parkwithease.parkeasy.model.Car
import io.github.parkwithease.parkeasy.model.CarDetails
import kotlinx.coroutines.launch

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun CarsScreen(
    showSnackbar: suspend (String, String?) -> Boolean,
    onSelectCar: (Car) -> Unit,
    modifier: Modifier = Modifier,
    viewModel: CarsViewModel =
        hiltViewModel<CarsViewModel, CarsViewModel.Factory>(
            creationCallback = { factory -> factory.create(showSnackbar = showSnackbar) }
        ),
) {
    val cars by viewModel.cars.collectAsState()
    val isRefreshing by viewModel.isRefreshing.collectAsState()
    val state = rememberPullToRefreshState()
    CarsScreen(cars, isRefreshing, state, viewModel::onRefresh, onSelectCar, modifier)
    LaunchedEffect(Unit) { viewModel.onRefresh() }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun CarsScreen(
    cars: List<Car>,
    isRefreshing: Boolean,
    state: PullToRefreshState,
    onRefresh: () -> Unit,
    onCarClick: (Car) -> Unit,
    modifier: Modifier = Modifier,
) {
    Surface(modifier) {
        var addingCar by rememberSaveable { mutableStateOf(false) }
        val skipPartiallyExpanded by rememberSaveable { mutableStateOf(false) }
        val scope = rememberCoroutineScope()
        val bottomSheetState =
            rememberModalBottomSheetState(skipPartiallyExpanded = skipPartiallyExpanded)
        Box {
            PullToRefreshBox(
                isRefreshing = isRefreshing,
                onRefresh = onRefresh,
                state = state,
                modifier = Modifier.padding(4.dp),
            ) {
                LazyColumn(
                    Modifier.fillMaxSize(),
                    horizontalAlignment = Alignment.CenterHorizontally,
                ) {
                    items(cars.count()) { index ->
                        ListItem({ CarCard(cars[index], onCarClick, Modifier.fillMaxWidth()) })
                    }
                }
            }
            ButtonBar(
                onRefresh,
                { addingCar = true },
                Modifier.align(Alignment.BottomEnd).padding(16.dp),
            )
            if (addingCar) {
                ModalBottomSheet(
                    onDismissRequest = { addingCar = false },
                    sheetState = bottomSheetState,
                ) {
                    AddCarForm(
                        {
                            scope
                                .launch { bottomSheetState.hide() }
                                .invokeOnCompletion {
                                    if (!bottomSheetState.isVisible) {
                                        addingCar = false
                                    }
                                }
                        },
                        { true },
                    )
                }
            }
        }
    }
}

@Composable
private fun AddCarForm(
    onClose: () -> Unit,
    onAddCar: (CarDetails) -> Boolean,
    modifier: Modifier = Modifier,
) {
    var carDetails by remember { mutableStateOf(CarDetails()) }
    Column(modifier) {
        OutlinedTextField(
            value = carDetails.licensePlate,
            onValueChange = { carDetails = carDetails.copy(licensePlate = it) },
            label = { Text(stringResource(R.string.license_plate)) },
            modifier = Modifier.fillMaxWidth().padding(16.dp, 4.dp),
        )
        OutlinedTextField(
            value = carDetails.color,
            onValueChange = { carDetails = carDetails.copy(color = it) },
            label = { Text(stringResource(R.string.color)) },
            modifier = Modifier.fillMaxWidth().padding(16.dp, 4.dp),
        )
        OutlinedTextField(
            value = carDetails.model,
            onValueChange = { carDetails = carDetails.copy(model = it) },
            label = { Text(stringResource(R.string.model)) },
            modifier = Modifier.fillMaxWidth().padding(16.dp, 4.dp),
        )
        OutlinedTextField(
            value = carDetails.make,
            onValueChange = { carDetails = carDetails.copy(make = it) },
            label = { Text(stringResource(R.string.make)) },
            modifier = Modifier.fillMaxWidth().padding(16.dp, 4.dp),
        )
        Button(
            content = { Text(stringResource(R.string.add_car)) },
            onClick = {
                onClose()
                onAddCar(carDetails)
            },
            modifier = Modifier.fillMaxWidth().padding(16.dp, 4.dp),
        )
    }
}

@Composable
private fun ButtonBar(onRefresh: () -> Unit, onAdd: () -> Unit, modifier: Modifier = Modifier) {
    Row(modifier, Arrangement.spacedBy(8.dp)) {
        RefreshButton(onRefresh)
        AddButton(onAdd)
    }
}

@Composable
private fun RefreshButton(onClick: () -> Unit, modifier: Modifier = Modifier) {
    FloatingActionButton(
        onClick = onClick,
        modifier = modifier,
        content = {
            Icon(
                imageVector = Icons.Outlined.Refresh,
                contentDescription = stringResource(R.string.refresh),
            )
        },
    )
}

@Composable
private fun AddButton(onClick: () -> Unit, modifier: Modifier = Modifier) {
    FloatingActionButton(
        onClick = onClick,
        modifier = modifier,
        content = {
            Icon(
                imageVector = Icons.Outlined.Add,
                contentDescription = stringResource(R.string.add_car),
            )
        },
    )
}

@Composable
private fun CarCard(car: Car, onCarClick: (Car) -> Unit, modifier: Modifier = Modifier) {
    Card(modifier = modifier.clickable { onCarClick(car) }) {
        Text(car.details.color, Modifier.padding(8.dp, 8.dp, 0.dp, 0.dp))
        Text(car.details.model, Modifier.padding(8.dp, 8.dp, 0.dp, 0.dp))
        Text(car.details.make, Modifier.padding(8.dp, 8.dp, 0.dp, 0.dp))
        Text(car.details.licensePlate, Modifier.padding(8.dp, 8.dp, 0.dp, 8.dp))
    }
}
