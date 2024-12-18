package io.github.parkwithease.parkeasy.ui.cars

import androidx.activity.compose.BackHandler
import androidx.compose.foundation.Image
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.heightIn
import androidx.compose.foundation.layout.imePadding
import androidx.compose.foundation.layout.padding
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
import androidx.compose.ui.tooling.preview.PreviewParameter
import androidx.compose.ui.tooling.preview.PreviewParameterProvider
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import io.github.parkwithease.parkeasy.R
import io.github.parkwithease.parkeasy.model.Car
import io.github.parkwithease.parkeasy.model.CarDetails
import io.github.parkwithease.parkeasy.model.EditMode
import io.github.parkwithease.parkeasy.model.FieldState
import io.github.parkwithease.parkeasy.ui.common.CarDetailsText
import io.github.parkwithease.parkeasy.ui.common.ParkEasyTextField
import io.github.parkwithease.parkeasy.ui.common.PreviewAll
import io.github.parkwithease.parkeasy.ui.common.PullToRefreshBox
import io.github.parkwithease.parkeasy.ui.theme.ParkEasyTheme

@Suppress("detekt:LongMethod")
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun CarsScreen(modifier: Modifier = Modifier, viewModel: CarsViewModel = hiltViewModel()) {
    val handler = rememberAddCarFormHandler(viewModel)
    val cars by viewModel.cars.collectAsState()
    val formEnabled by viewModel.formEnabled.collectAsState()
    val isRefreshing by viewModel.isRefreshing.collectAsState()
    val showForm by viewModel.showForm.collectAsState()
    var editMode by rememberSaveable { mutableStateOf(EditMode.ADD) }

    val skipPartiallyExpanded by rememberSaveable { mutableStateOf(false) }
    val bottomSheetState =
        rememberModalBottomSheetState(skipPartiallyExpanded = skipPartiallyExpanded)

    BackHandler(enabled = showForm) { viewModel.onHideForm() }

    LaunchedEffect(Unit) { viewModel.onRefresh() }

    CarsScreen(
        cars = cars,
        onCarClick = { car ->
            handler.resetForm()
            handler.onColorChange(car.details.color)
            handler.onLicensePlateChange(car.details.licensePlate)
            handler.onMakeChange(car.details.make)
            handler.onModelChange(car.details.model)
            viewModel.currentlyEditingId = car.id
            editMode = EditMode.EDIT
            viewModel.onShowForm()
        },
        onShowAddCarClick = {
            handler.resetForm()
            viewModel.currentlyEditingId = ""
            editMode = EditMode.ADD
            viewModel.onShowForm()
        },
        isRefreshing = isRefreshing,
        onRefresh = viewModel::onRefresh,
        snackbarHost = { SnackbarHost(hostState = viewModel.snackbarState) },
        modifier = modifier,
    )
    if (showForm) {
        ModalBottomSheet(
            onDismissRequest = { viewModel.onHideForm() },
            modifier = Modifier.fillMaxWidth(),
            sheetState = bottomSheetState,
        ) {
            AddCarScreen(
                state = viewModel.formState,
                handler = handler,
                formEnabled = formEnabled,
                editMode = editMode,
            )
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
        modifier = modifier,
        snackbarHost = snackbarHost,
        floatingActionButton = { AddCarButton(onShowAddCarClick = onShowAddCarClick) },
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
    Card(onClick = { onClick(car) }, modifier = modifier.fillMaxWidth()) {
        Row(modifier = Modifier.padding(8.dp)) {
            Column(modifier = Modifier.weight(1f), verticalArrangement = Arrangement.Center) {
                Image(
                    painter = painterResource(R.drawable.wordmark),
                    contentDescription = null,
                    modifier = Modifier.heightIn(max = 64.dp),
                )
            }
            CarDetailsText(
                carDetails = car.details,
                modifier = Modifier.weight(1f),
                horizontalAlignment = Alignment.End,
            )
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
    state: AddCarFormState,
    handler: AddCarFormHandler,
    formEnabled: Boolean,
    editMode: EditMode,
    modifier: Modifier = Modifier,
) {
    Column(
        modifier =
            modifier
                .imePadding()
                .padding(horizontal = 32.dp)
                .verticalScroll(rememberScrollState(), reverseScrolling = true),
        verticalArrangement = Arrangement.spacedBy(2.dp, Alignment.CenterVertically),
        horizontalAlignment = Alignment.CenterHorizontally,
    ) {
        ParkEasyTextField(
            state = state.licensePlate,
            onValueChange = handler.onLicensePlateChange,
            modifier = Modifier.fillMaxWidth(),
            enabled = formEnabled,
            labelId = R.string.license_plate,
        )
        ParkEasyTextField(
            state = state.color,
            onValueChange = handler.onColorChange,
            modifier = Modifier.fillMaxWidth(),
            enabled = formEnabled,
            labelId = R.string.color,
        )
        ParkEasyTextField(
            state = state.make,
            onValueChange = handler.onMakeChange,
            modifier = Modifier.fillMaxWidth(),
            enabled = formEnabled,
            labelId = R.string.make,
        )
        ParkEasyTextField(
            state = state.model,
            onValueChange = handler.onModelChange,
            modifier = Modifier.fillMaxWidth(),
            enabled = formEnabled,
            labelId = R.string.model,
        )
        Button(
            onClick =
                if (editMode == EditMode.ADD) handler.onAddCarClick else handler.onEditCarClick,
            modifier = Modifier.fillMaxWidth(),
            enabled = formEnabled,
        ) {
            Text(
                stringResource(
                    if (editMode == EditMode.ADD) R.string.add_car else R.string.edit_car
                )
            )
        }
    }
}

@Suppress("detekt:UnusedPrivateMember")
@PreviewAll
@Composable
private fun CarsScreenPreview() {
    val cars =
        listOf(
                CarDetails("Red", "MV1", "Red Bull", "RB20"),
                CarDetails("Red", "SP11", "Red Bull", "RB20"),
                CarDetails("Silver", "LH44", "Mercedes", "W15"),
                CarDetails("Silver", "GR63", "Mercedes", "W15"),
                CarDetails("Scarlet", "CL16", "Ferrari", "SF-24"),
                CarDetails("Scarlet", "CS55", "Ferrari", "SF-24"),
                CarDetails("Papaya", "LN4", "McLaren", "MCL38"),
                CarDetails("Papaya", "OP81", "McLaren", "MCL38"),
            )
            .map { Car(it, "") }
    ParkEasyTheme {
        CarsScreen(
            cars = cars,
            onCarClick = {},
            onShowAddCarClick = {},
            isRefreshing = false,
            onRefresh = {},
            snackbarHost = {},
        )
    }
}

private class AddCarFormStateProvider : PreviewParameterProvider<AddCarFormState> {
    override val values =
        sequenceOf(
            AddCarFormState(),
            AddCarFormState(
                color = FieldState("Red"),
                licensePlate = FieldState("MV1"),
                make = FieldState("Red Bull"),
                model = FieldState("RB20"),
            ),
            AddCarFormState(
                color = FieldState("", "Color cannot be empty"),
                licensePlate = FieldState("", "License plate cannot be empty"),
                make = FieldState("", "Make cannot be empty"),
                model = FieldState("", "Model cannot be empty"),
            ),
        )
}

@Suppress("detekt:UnusedPrivateMember")
@OptIn(ExperimentalMaterial3Api::class)
@PreviewAll
@Composable
private fun AddCarScreenPreview(
    @PreviewParameter(AddCarFormStateProvider::class) state: AddCarFormState
) {
    ParkEasyTheme {
        ModalBottomSheet(
            onDismissRequest = {},
            sheetState = rememberStandardBottomSheetState(initialValue = SheetValue.Expanded),
        ) {
            AddCarScreen(
                state = state,
                handler = AddCarFormHandler(),
                formEnabled = true,
                editMode = EditMode.ADD,
            )
        }
    }
}
