package io.github.parkwithease.parkeasy.ui.cars

import androidx.compose.material3.SnackbarHostState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.setValue
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import io.github.parkwithease.parkeasy.data.remote.APIException
import io.github.parkwithease.parkeasy.data.remote.CarRepository
import io.github.parkwithease.parkeasy.model.Car
import io.github.parkwithease.parkeasy.model.CarDetails
import io.github.parkwithease.parkeasy.model.ErrorDetail
import io.github.parkwithease.parkeasy.model.ErrorModel
import io.github.parkwithease.parkeasy.model.FieldState
import java.io.IOException
import javax.inject.Inject
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch

@HiltViewModel
@Suppress("detekt:TooManyFunctions")
class CarsViewModel @Inject constructor(private val carRepo: CarRepository) : ViewModel() {
    val snackbarState = SnackbarHostState()

    private val _cars = MutableStateFlow(emptyList<Car>())
    val cars = _cars.asStateFlow()

    private val _isRefreshing = MutableStateFlow(false)
    val isRefreshing = _isRefreshing.asStateFlow()

    var currentlyEditingId = ""

    var formState by mutableStateOf(AddCarFormState())
        private set

    fun onRefresh() {
        viewModelScope.launch {
            _isRefreshing.value = true
            carRepo
                .getCars()
                .onSuccess { _cars.value = it }
                .onFailure {
                    viewModelScope.launch { snackbarState.showSnackbar("Error retrieving cars") }
                }
            _isRefreshing.value = false
        }
    }

    fun onAddCarClick() {
        viewModelScope.launch {
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
                .onSuccess { onRefresh() }
                .recoverRequestErrors("Error adding car")
        }
    }

    fun onEditCarClick() {
        viewModelScope.launch {
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
                .onSuccess { onRefresh() }
                .recoverRequestErrors("Error adding car")
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

    private fun Result<Unit>.recoverRequestErrors(operationFailMsg: String): Result<Unit> =
        recover {
            when (it) {
                is APIException -> {
                    errorToForm(it.error)
                    viewModelScope.launch { snackbarState.showSnackbar(operationFailMsg) }
                }
                is IOException -> {
                    viewModelScope.launch {
                        snackbarState.showSnackbar("Could not connect to server, are you online?")
                    }
                }
                else -> throw it
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

    // Clear errors set via external services
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
}

data class AddCarFormState(
    val color: FieldState<String> = FieldState(""),
    val licensePlate: FieldState<String> = FieldState(""),
    val make: FieldState<String> = FieldState(""),
    val model: FieldState<String> = FieldState(""),
)
