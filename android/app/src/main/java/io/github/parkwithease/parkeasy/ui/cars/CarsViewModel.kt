package io.github.parkwithease.parkeasy.ui.cars

import androidx.compose.material3.SnackbarHostState
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import io.github.parkwithease.parkeasy.data.remote.CarRepository
import io.github.parkwithease.parkeasy.model.Car
import io.github.parkwithease.parkeasy.model.CarDetails
import io.github.parkwithease.parkeasy.model.ErrorDetail
import io.github.parkwithease.parkeasy.model.ErrorModel
import io.github.parkwithease.parkeasy.model.FieldState
import io.github.parkwithease.parkeasy.ui.common.recoverRequestErrors
import javax.inject.Inject
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch

@Suppress("detekt:TooManyFunctions")
@HiltViewModel
class CarsViewModel @Inject constructor(private val carRepo: CarRepository) : ViewModel() {
    val snackbarState = SnackbarHostState()

    private val _cars = MutableStateFlow(emptyList<Car>())
    val cars = _cars.asStateFlow()

    private val _isRefreshing = MutableStateFlow(false)
    val isRefreshing = _isRefreshing.asStateFlow()

    private val _showForm = MutableStateFlow(false)
    val showForm = _showForm.asStateFlow()

    private val _formEnabled = MutableStateFlow(true)
    val formEnabled = _formEnabled.asStateFlow()

    var currentlyEditingId = ""

    var formState by mutableStateOf(AddCarFormState())
        private set

    fun onShowForm() {
        _formEnabled.value = true
        viewModelScope.launch { _showForm.value = true }
    }

    fun onHideForm() {
        viewModelScope.launch { _showForm.value = false }
    }

    fun onRefresh() {
        viewModelScope.launch {
            _isRefreshing.value = true
            carRepo
                .getCars()
                .onSuccess { _cars.value = it }
                .recoverRequestErrors(
                    "Error retrieving cars",
                    { errorToForm(it) },
                    snackbarState,
                    viewModelScope,
                )
            _isRefreshing.value = false
        }
    }

    fun onAddCarClick() {
        viewModelScope.launch {
            _formEnabled.value = false
            carRepo
                .createCar(
                    CarDetails(
                        formState.color.value,
                        formState.licensePlate.value,
                        formState.make.value,
                        formState.model.value,
                    )
                )
                .also { clearFieldErrors() }
                .onSuccess {
                    onRefresh()
                    onHideForm()
                }
                .recoverRequestErrors(
                    "Error adding car",
                    { errorToForm(it) },
                    snackbarState,
                    viewModelScope,
                )
            _formEnabled.value = true
        }
    }

    fun onEditCarClick() {
        viewModelScope.launch {
            _formEnabled.value = false
            carRepo
                .updateCar(
                    Car(
                        details =
                            CarDetails(
                                formState.color.value,
                                formState.licensePlate.value,
                                formState.make.value,
                                formState.model.value,
                            ),
                        id = currentlyEditingId,
                    )
                )
                .also { clearFieldErrors() }
                .onSuccess {
                    onRefresh()
                    onHideForm()
                }
                .recoverRequestErrors(
                    "Error adding car",
                    { errorToForm(it) },
                    snackbarState,
                    viewModelScope,
                )
            _formEnabled.value = true
        }
    }

    fun onLicensePlateChange(value: String) {
        formState =
            formState.run {
                copy(
                    licensePlate =
                        licensePlate.copy(
                            value = value,
                            error = if (value != "") null else "License plate cannot be empty",
                        )
                )
            }
    }

    fun onColorChange(value: String) {
        formState =
            formState.run {
                copy(
                    color =
                        color.copy(
                            value = value,
                            error = if (value != "") null else "Colour cannot be empty",
                        )
                )
            }
    }

    fun onMakeChange(value: String) {
        formState =
            formState.run {
                copy(
                    make =
                        make.copy(
                            value = value,
                            error = if (value != "") null else "Make cannot be empty",
                        )
                )
            }
    }

    fun onModelChange(value: String) {
        formState =
            formState.run {
                copy(
                    model =
                        model.copy(
                            value = value,
                            error = if (value != "") null else "Model cannot be empty",
                        )
                )
            }
    }

    fun resetForm() {
        formState =
            formState.run {
                copy(
                    color = FieldState(""),
                    licensePlate = FieldState(""),
                    make = FieldState(""),
                    model = FieldState(""),
                )
            }
    }

    private fun errorToForm(error: ErrorModel) {
        annotateErrorLocation(error.errors)
    }

    private fun annotateErrorLocation(errors: List<ErrorDetail>) {
        for (err in errors) {
            when (err.location) {
                "body" ->
                    formState =
                        formState.run {
                            copy(
                                licensePlate = licensePlate.copy(error = "Invalid license plate"),
                                color = color.copy(error = "Invalid color"),
                                make = make.copy(error = "Invalid make"),
                                model = model.copy(error = "Invalid model"),
                            )
                        }

                "body.license_plate" ->
                    formState =
                        formState.run {
                            copy(licensePlate = licensePlate.copy(error = "Invalid license plate"))
                        }

                "body.color" ->
                    formState = formState.run { copy(color = color.copy(error = "Invalid color")) }

                "body.make" ->
                    formState = formState.run { copy(make = make.copy(error = "Invalid make")) }

                "body.model" ->
                    formState = formState.run { copy(model = model.copy(error = "Invalid model")) }
            }
        }
    }

    private fun clearFieldErrors() {
        formState =
            formState.run {
                copy(
                    licensePlate = licensePlate.copy(error = null),
                    color = color.copy(error = null),
                    make = make.copy(error = null),
                    model = model.copy(error = null),
                )
            }
    }

    fun createHandler() =
        AddCarFormHandler(
            onColorChange = this::onColorChange,
            onLicensePlateChange = this::onLicensePlateChange,
            onMakeChange = this::onMakeChange,
            onModelChange = this::onModelChange,
            onAddCarClick = this::onAddCarClick,
            onEditCarClick = this::onEditCarClick,
            resetForm = this::resetForm,
        )
}

@Composable
fun rememberAddCarFormHandler(viewModel: CarsViewModel) =
    remember(viewModel) { viewModel.createHandler() }

data class AddCarFormState(
    val color: FieldState<String> = FieldState(""),
    val licensePlate: FieldState<String> = FieldState(""),
    val make: FieldState<String> = FieldState(""),
    val model: FieldState<String> = FieldState(""),
)

data class AddCarFormHandler(
    val onColorChange: (String) -> Unit = {},
    val onLicensePlateChange: (String) -> Unit = {},
    val onMakeChange: (String) -> Unit = {},
    val onModelChange: (String) -> Unit = {},
    val onAddCarClick: () -> Unit = {},
    val onEditCarClick: () -> Unit = {},
    val resetForm: () -> Unit = {},
)
