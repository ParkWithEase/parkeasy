package io.github.parkwithease.parkeasy.ui.cars

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
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Add
import androidx.compose.material3.Button
import androidx.compose.material3.Card
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.FloatingActionButton
import androidx.compose.material3.Icon
import androidx.compose.material3.ModalBottomSheet
import androidx.compose.material3.Scaffold
import androidx.compose.material3.SnackbarHost
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
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import io.github.parkwithease.parkeasy.R
import io.github.parkwithease.parkeasy.model.Car
import io.github.parkwithease.parkeasy.model.EditMode
import io.github.parkwithease.parkeasy.ui.common.ParkEasyTextField
import io.github.parkwithease.parkeasy.ui.common.PullToRefreshBox
import io.github.parkwithease.parkeasy.ui.theme.Typography

@OptIn(ExperimentalMaterial3Api::class)
@Composable
@Suppress("detekt:LongMethod")
fun CarsScreen(modifier: Modifier = Modifier, viewModel: CarsViewModel = hiltViewModel()) {
    val cars by viewModel.cars.collectAsState()
    val isRefreshing by viewModel.isRefreshing.collectAsState()
    var editMode by rememberSaveable { mutableStateOf(EditMode.ADD) }

    var openBottomSheet by rememberSaveable { mutableStateOf(false) }
    val skipPartiallyExpanded by rememberSaveable { mutableStateOf(false) }
    val bottomSheetState =
        rememberModalBottomSheetState(skipPartiallyExpanded = skipPartiallyExpanded)

    LaunchedEffect(Unit) { viewModel.onRefresh() }
    CarsScreen(
        cars,
        { car ->
            viewModel.currentlyEditingId = car.id
            viewModel.onLicensePlateChange(car.details.licensePlate)
            viewModel.onColorChange(car.details.color)
            viewModel.onMakeChange(car.details.make)
            viewModel.onModelChange(car.details.model)
            editMode = EditMode.EDIT
            openBottomSheet = true
        },
        {
            viewModel.currentlyEditingId = ""
            viewModel.onLicensePlateChange("")
            viewModel.onColorChange("")
            viewModel.onMakeChange("")
            viewModel.onModelChange("")
            editMode = EditMode.ADD
            openBottomSheet = true
        },
        isRefreshing,
        viewModel::onRefresh,
        { SnackbarHost(hostState = viewModel.snackbarState) },
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
                    editMode,
                    viewModel.formState,
                    viewModel::onLicensePlateChange,
                    viewModel::onColorChange,
                    viewModel::onMakeChange,
                    viewModel::onModelChange,
                    viewModel::onAddCarClick,
                    viewModel::onEditCarClick,
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
    snackbarHost: @Composable () -> Unit,
    modifier: Modifier = Modifier,
) {
    Scaffold(
        floatingActionButton = { AddCarButton(onShowAddCarClick = onShowAddCarClick) },
        snackbarHost = snackbarHost,
        modifier = modifier,
    ) { innerPadding ->
        Surface(Modifier.padding(innerPadding)) {
            PullToRefreshBox(
                items = cars,
                onClick = onCarClick,
                isRefreshing = isRefreshing,
                onRefresh = onRefresh,
                modifier = Modifier.padding(4.dp),
            ) { car, onClick ->
                CarCard(car, onClick)
            }
        }
    }
}

@Composable
fun CarCard(car: Car, onClick: (Car) -> Unit, modifier: Modifier = Modifier) {
    Card(onClick = { onClick(car) }, modifier = modifier.fillMaxWidth().padding(4.dp, 0.dp)) {
        Row(modifier = Modifier.padding(8.dp)) {
            Column(modifier = Modifier.weight(1f)) {
                Image(
                    painter = painterResource(R.drawable.wordmark_outlined),
                    contentDescription = null,
                    modifier = Modifier.heightIn(max = 64.dp),
                )
            }
            Column(horizontalAlignment = Alignment.End, modifier = Modifier.weight(1f)) {
                Text(text = car.details.licensePlate, style = Typography.headlineLarge)
                Text(car.details.color + ' ' + car.details.make + " " + car.details.model)
            }
        }
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
    editMode: EditMode,
    state: AddCarFormState,
    onLicensePlateChange: (String) -> Unit,
    onColorChange: (String) -> Unit,
    onMakeChange: (String) -> Unit,
    onModelChange: (String) -> Unit,
    onAddCarClick: () -> Unit,
    onEditCarClick: () -> Unit,
    modifier: Modifier = Modifier,
) {
    Column(
        verticalArrangement = Arrangement.spacedBy(2.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
        modifier = modifier.widthIn(max = 320.dp),
    ) {
        ParkEasyTextField(
            state = state.licensePlate,
            onValueChange = onLicensePlateChange,
            modifier = Modifier.fillMaxWidth(),
            labelId = R.string.license_plate,
        )
        ParkEasyTextField(
            state = state.color,
            onValueChange = onColorChange,
            modifier = Modifier.fillMaxWidth(),
            labelId = R.string.color,
        )
        ParkEasyTextField(
            state = state.make,
            onValueChange = onMakeChange,
            modifier = Modifier.fillMaxWidth(),
            labelId = R.string.make,
        )
        ParkEasyTextField(
            state = state.model,
            onValueChange = onModelChange,
            modifier = Modifier.fillMaxWidth(),
            labelId = R.string.model,
        )
        Button(
            onClick = if (editMode == EditMode.ADD) onAddCarClick else onEditCarClick,
            modifier = Modifier.fillMaxWidth(),
        ) {
            Text(
                stringResource(
                    if (editMode == EditMode.ADD) R.string.add_car else R.string.edit_car
                )
            )
        }
    }
}
