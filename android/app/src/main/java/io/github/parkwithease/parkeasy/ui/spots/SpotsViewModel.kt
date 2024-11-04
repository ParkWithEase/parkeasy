package io.github.parkwithease.parkeasy.ui.spots

import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.setValue
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import io.github.parkwithease.parkeasy.data.remote.SpotRepository
import io.github.parkwithease.parkeasy.model.FieldState
import io.github.parkwithease.parkeasy.model.Spot
import io.github.parkwithease.parkeasy.model.SpotFeatures
import io.github.parkwithease.parkeasy.model.SpotLocation
import io.github.parkwithease.parkeasy.model.TimeSlot
import javax.inject.Inject
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch

@HiltViewModel
class SpotsViewModel @Inject constructor(private val spotRepo: SpotRepository) : ViewModel() {
    private val _spots = MutableStateFlow(emptyList<Spot>())
    val spots = _spots.asStateFlow()

    private val _isRefreshing = MutableStateFlow(false)
    val isRefreshing = _isRefreshing.asStateFlow()

    var formState by mutableStateOf(AddSpotFormState())
        private set

    fun onRefresh() {
        viewModelScope.launch {
            _isRefreshing.value = true
            _spots.value = spotRepo.getSpots()
            _isRefreshing.value = false
        }
    }

    fun onAddSpotClick() {
        viewModelScope.launch {
            spotRepo
                .createSpot(
                    Spot(
                        availability =
                            listOf(
                                TimeSlot(
                                    startTime = "2024-11-04T01:30:00-06:00",
                                    endTime = "2024-11-04T02:00:00-06:00",
                                )
                            ),
                        features =
                            SpotFeatures(chargingStation = false, plugIn = false, shelter = false),
                        location =
                            SpotLocation(
                                streetAddress = formState.streetAddress.value,
                                city = formState.city.value,
                                state = formState.state.value,
                                countryCode = formState.countryCode.value,
                                postalCode = formState.postalCode.value,
                            ),
                        pricePerHour = 2.0,
                    )
                )
                .also { clearFieldErrors() }
                ?.onSuccess { onRefresh() }
        }
    }

    fun onStreetAddressChange(value: String) {
        formState = formState.run { copy(streetAddress = streetAddress.copy(value = value)) }
    }

    fun onCityChange(value: String) {
        formState = formState.run { copy(city = city.copy(value = value)) }
    }

    fun onStateChange(value: String) {
        formState = formState.run { copy(state = state.copy(value = value)) }
    }

    fun onCountryCodeChange(value: String) {
        formState = formState.run { copy(countryCode = countryCode.copy(value = value)) }
    }

    fun onPostalCodeChange(value: String) {
        formState = formState.run { copy(postalCode = postalCode.copy(value = value)) }
    }

    private fun clearFieldErrors() {
        formState =
            formState.run {
                copy(
                    streetAddress = streetAddress.copy(error = null),
                    city = city.copy(error = null),
                    state = state.copy(error = null),
                    countryCode = countryCode.copy(error = null),
                )
            }
    }
}

data class AddSpotFormState(
    val streetAddress: FieldState = FieldState(),
    val city: FieldState = FieldState(),
    val state: FieldState = FieldState(),
    val countryCode: FieldState = FieldState(),
    val postalCode: FieldState = FieldState(),
)
