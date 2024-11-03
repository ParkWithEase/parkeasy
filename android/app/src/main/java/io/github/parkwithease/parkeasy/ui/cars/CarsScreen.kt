package io.github.parkwithease.parkeasy.ui.cars

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.imePadding
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.widthIn
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Add
import androidx.compose.material3.Button
import androidx.compose.material3.Card
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.FloatingActionButton
import androidx.compose.material3.Icon
import androidx.compose.material3.ModalBottomSheet
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Scaffold
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
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import io.github.parkwithease.parkeasy.R
import io.github.parkwithease.parkeasy.common.PullToRefreshBox
import io.github.parkwithease.parkeasy.model.Car

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun CarsScreen(
    onCarClick: (Car) -> Unit,
    modifier: Modifier = Modifier,
    viewModel: CarsViewModel = hiltViewModel(),
) {
    val cars by viewModel.cars.collectAsState()
    val isRefreshing by viewModel.isRefreshing.collectAsState()

    var openBottomSheet by rememberSaveable { mutableStateOf(false) }
    val skipPartiallyExpanded by rememberSaveable { mutableStateOf(false) }
    val bottomSheetState =
        rememberModalBottomSheetState(skipPartiallyExpanded = skipPartiallyExpanded)

    LaunchedEffect(Unit) { viewModel.onRefresh() }
    CarsScreen(
        cars,
        onCarClick,
        { openBottomSheet = true },
        isRefreshing,
        viewModel::onRefresh,
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
                AddCarScreen(
                    viewModel.formState,
                    viewModel::onColorChange,
                    viewModel::onLicensePlateChange,
                    viewModel::onMakeChange,
                    viewModel::onModelChange,
                    viewModel::onAddCarClick,
                )
            }
        }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun CarsScreen(
    cars: List<Car>,
    onCarClick: (Car) -> Unit,
    onShowAddCarClick: () -> Unit,
    isRefreshing: Boolean,
    onRefresh: () -> Unit,
    modifier: Modifier = Modifier,
) {
    Scaffold(
        floatingActionButton = { AddCarButton(onShowAddCarClick = onShowAddCarClick) },
        modifier = modifier,
    ) { innerPadding ->
        Surface(Modifier.padding(innerPadding)) {
            PullToRefreshBox(
                items = cars,
                onClick = onCarClick,
                isRefreshing = isRefreshing,
                onRefresh = onRefresh,
                modifier = Modifier.padding(4.dp),
                card = { car, onClick -> CarCard(car, onClick) },
            )
        }
    }
}

@Composable
fun CarCard(car: Car, onClick: (Car) -> Unit, modifier: Modifier = Modifier) {
    Card(onClick = { onClick(car) }, modifier = modifier.fillMaxWidth().padding(8.dp, 0.dp)) {
        Text(car.details.color)
        Text(car.details.model)
        Text(car.details.make)
        Text(car.details.licensePlate)
    }
}

@Composable
fun AddCarButton(onShowAddCarClick: () -> Unit, modifier: Modifier = Modifier) {
    FloatingActionButton(onClick = onShowAddCarClick, modifier) {
        Icon(imageVector = Icons.Filled.Add, contentDescription = stringResource(R.string.add_car))
    }
}

@Composable
fun AddCarScreen(
    state: AddCarFormState,
    onColorChange: (String) -> Unit,
    onLicensePlateChange: (String) -> Unit,
    onMakeChange: (String) -> Unit,
    onModelChange: (String) -> Unit,
    onAddCarClick: () -> Unit,
    modifier: Modifier = Modifier,
) {
    Column(
        verticalArrangement = Arrangement.spacedBy(2.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
        modifier = modifier.widthIn(max = 320.dp),
    ) {
        OutlinedTextField(
            value = state.color.value,
            onValueChange = onColorChange,
            label = { Text(stringResource(R.string.color)) },
            modifier = Modifier.fillMaxWidth(),
        )
        OutlinedTextField(
            value = state.licensePlate.value,
            onValueChange = onLicensePlateChange,
            label = { Text(stringResource(R.string.license_plate)) },
            modifier = Modifier.fillMaxWidth(),
        )
        OutlinedTextField(
            value = state.make.value,
            onValueChange = onMakeChange,
            label = { Text(stringResource(R.string.make)) },
            modifier = Modifier.fillMaxWidth(),
        )
        OutlinedTextField(
            value = state.model.value,
            onValueChange = onModelChange,
            label = { Text(stringResource(R.string.model)) },
            modifier = Modifier.fillMaxWidth(),
        )
        Button(onClick = onAddCarClick, modifier = Modifier.fillMaxWidth()) {
            Text(stringResource(R.string.add_car))
        }
    }
}
