package io.github.parkwithease.parkeasy.ui.spots

import androidx.compose.material3.SnackbarHostState
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
import io.github.parkwithease.parkeasy.ui.common.startOfWeek
import io.github.parkwithease.parkeasy.ui.common.timezone
import javax.inject.Inject
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import kotlinx.datetime.DateTimeUnit
import kotlinx.datetime.plus
import kotlinx.datetime.toInstant
import kotlinx.datetime.toLocalDateTime

private const val MinutesPerSlot = 30

@HiltViewModel
@Suppress("detekt:TooManyFunctions")
class SpotsViewModel @Inject constructor(private val spotRepo: SpotRepository) : ViewModel() {
    val snackbarState = SnackbarHostState()

    private val _spots = MutableStateFlow(emptyList<Spot>())
    val spots = _spots.asStateFlow()

    private val _isRefreshing = MutableStateFlow(false)
    val isRefreshing = _isRefreshing.asStateFlow()

    var formState by mutableStateOf(AddSpotFormState())
        private set

    fun onRefresh() {
        viewModelScope.launch {
            _isRefreshing.value = true
            spotRepo
                .getSpots()
                .onSuccess { _spots.value = it }
                .onFailure {
                    viewModelScope.launch { snackbarState.showSnackbar("Error retrieving spots") }
                }
            _isRefreshing.value = false
        }
    }

    @Suppress("detekt:LongMethod")
    fun onAddSpotClick() {
        val timezone = timezone()
        val startOfWeek = startOfWeek()
        viewModelScope.launch {
            spotRepo
                .createSpot(
                    Spot(
                        availability =
                            formState.times.value
                                .sortedBy { it }
                                .map {
                                    TimeSlot(
                                        startTime =
                                            startOfWeek
                                                .toInstant(timezone)
                                                .plus(
                                                    MinutesPerSlot * it,
                                                    DateTimeUnit.MINUTE,
                                                    timezone,
                                                )
                                                .toLocalDateTime(timezone),
                                        endTime =
                                            startOfWeek
                                                .toInstant(timezone)
                                                .plus(
                                                    MinutesPerSlot * (it + 1),
                                                    DateTimeUnit.MINUTE,
                                                    timezone,
                                                )
                                                .toLocalDateTime(timezone),
                                    )
                                },
                        features =
                            SpotFeatures(
                                chargingStation = formState.chargingStation.value,
                                plugIn = formState.plugIn.value,
                                shelter = formState.shelter.value,
                            ),
                        location =
                            SpotLocation(
                                streetAddress = formState.streetAddress.value,
                                city = formState.city.value,
                                state = formState.state.value,
                                countryCode = formState.countryCode.value,
                                postalCode = formState.postalCode.value,
                            ),
                        pricePerHour = formState.pricePerHour.value.toDoubleOrNull() ?: -1.0,
                    )
                )
                .also { clearFieldErrors() }
                .onSuccess { onRefresh() }
                .onFailure {
                    viewModelScope.launch { snackbarState.showSnackbar("Error creating spot") }
                }
        }
    }

    fun onStreetAddressChange(value: String) {
        formState =
            formState.run {
                copy(
                    streetAddress =
                        streetAddress.copy(
                            value = value,
                            error = if (value != "") null else "Address cannot be empty",
                        )
                )
            }
    }

    fun onCityChange(value: String) {
        formState =
            formState.run {
                copy(
                    city =
                        city.copy(
                            value = value,
                            error = if (value != "") null else "City cannot be empty",
                        )
                )
            }
    }

    fun onStateChange(value: String) {
        formState =
            formState.run {
                copy(
                    state =
                        state.copy(
                            value = value,
                            error = if (value != "") null else "State cannot be empty",
                        )
                )
            }
    }

    fun onCountryCodeChange(value: String) {
        formState =
            formState.run {
                copy(
                    countryCode =
                        countryCode.copy(
                            value = value,
                            error = if (value != "") null else "Country cannot be empty",
                        )
                )
            }
    }

    fun onPostalCodeChange(value: String) {
        formState =
            formState.run {
                copy(
                    postalCode =
                        postalCode.copy(
                            value = value,
                            error = if (value != "") null else "Postal code cannot be empty",
                        )
                )
            }
    }

    fun onChargingStationChange(value: Boolean) {
        formState = formState.run { copy(chargingStation = chargingStation.copy(value = value)) }
    }

    fun onPlugInChange(value: Boolean) {
        formState = formState.run { copy(plugIn = plugIn.copy(value = value)) }
    }

    fun onShelterChange(value: Boolean) {
        formState = formState.run { copy(shelter = shelter.copy(value = value)) }
    }

    fun onPricePerHourChange(value: String) {
        formState =
            formState.run {
                copy(
                    pricePerHour =
                        pricePerHour.copy(
                            value = value,
                            error = if (value != "") null else "Price cannot be empty",
                        )
                )
            }
    }

    fun onPlusTime(elements: Iterable<Int>) {
        formState =
            formState.run { copy(times = times.copy(value = formState.times.value.plus(elements))) }
    }

    fun onMinusTime(elements: Iterable<Int>) {
        formState =
            formState.run {
                copy(times = times.copy(value = formState.times.value.minus(elements)))
            }
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
    val streetAddress: FieldState<String> = FieldState(""),
    val city: FieldState<String> = FieldState(""),
    val state: FieldState<String> = FieldState(""),
    val countryCode: FieldState<String> = FieldState("CA"),
    val postalCode: FieldState<String> = FieldState(""),
    val chargingStation: FieldState<Boolean> = FieldState(false),
    val plugIn: FieldState<Boolean> = FieldState(false),
    val shelter: FieldState<Boolean> = FieldState(false),
    val pricePerHour: FieldState<String> = FieldState(""),
    val times: FieldState<Set<Int>> = FieldState<Set<Int>>(emptySet()),
)
