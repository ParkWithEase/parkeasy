package io.github.parkwithease.parkeasy.ui.cars

import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.setValue
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import io.github.parkwithease.parkeasy.data.remote.CarRepository
import io.github.parkwithease.parkeasy.model.Car
import io.github.parkwithease.parkeasy.model.CarDetails
import javax.inject.Inject
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch

@HiltViewModel
class CarsViewModel @Inject constructor(private val carRepo: CarRepository) : ViewModel() {
    private val _cars = MutableStateFlow(emptyList<Car>())
    val cars = _cars.asStateFlow()

    private val _isRefreshing = MutableStateFlow(false)
    val isRefreshing = _isRefreshing.asStateFlow()

    var formState by mutableStateOf(AddCarFormState())
        private set

    fun onRefresh() {
        viewModelScope.launch {
            _isRefreshing.value = true
            _cars.value = carRepo.getCars()
            _isRefreshing.value = false
        }
    }

    fun onAddCarClick() {
        viewModelScope.launch {
            carRepo.createCar(
                CarDetails(
                    formState.color.value,
                    formState.licensePlate.value,
                    formState.make.value,
                    formState.model.value,
                )
            )
            onRefresh()
        }
    }

    fun onColorChange(value: String) {
        formState = formState.run { copy(color = color.copy(value = value)) }
    }

    fun onLicensePlateChange(value: String) {
        formState = formState.run { copy(licensePlate = licensePlate.copy(value = value)) }
    }

    fun onMakeChange(value: String) {
        formState = formState.run { copy(make = make.copy(value = value)) }
    }

    fun onModelChange(value: String) {
        formState = formState.run { copy(model = model.copy(value = value)) }
    }
}

data class AddCarFieldState(val value: String = "", val error: String? = null)

data class AddCarFormState(
    val color: AddCarFieldState = AddCarFieldState(),
    val licensePlate: AddCarFieldState = AddCarFieldState(),
    val make: AddCarFieldState = AddCarFieldState(),
    val model: AddCarFieldState = AddCarFieldState(),
)
