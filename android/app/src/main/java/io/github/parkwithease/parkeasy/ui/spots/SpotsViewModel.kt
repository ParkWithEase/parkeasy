package io.github.parkwithease.parkeasy.ui.spots

import androidx.compose.material3.SnackbarHostState
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import io.github.parkwithease.parkeasy.data.remote.SpotRepository
import io.github.parkwithease.parkeasy.model.ErrorDetail
import io.github.parkwithease.parkeasy.model.ErrorModel
import io.github.parkwithease.parkeasy.model.FieldState
import io.github.parkwithease.parkeasy.model.Spot
import io.github.parkwithease.parkeasy.model.SpotFeatures
import io.github.parkwithease.parkeasy.model.SpotLocation
import io.github.parkwithease.parkeasy.model.TimeUnit
import io.github.parkwithease.parkeasy.ui.common.recoverRequestErrors
import io.github.parkwithease.parkeasy.ui.common.startOfNextAvailableDayInstant
import io.github.parkwithease.parkeasy.ui.common.toIndex
import io.github.parkwithease.parkeasy.ui.common.toLocalDateTime
import javax.inject.Inject
import kotlin.String
import kotlin.collections.plus
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import kotlinx.datetime.Clock.System.now

@Suppress("detekt:TooManyFunctions")
@HiltViewModel
class SpotsViewModel @Inject constructor(private val spotRepo: SpotRepository) : ViewModel() {
    val snackbarState = SnackbarHostState()

    private val _spots = MutableStateFlow(emptyList<Spot>())
    val spots = _spots.asStateFlow()

    private val _isRefreshing = MutableStateFlow(false)
    val isRefreshing = _isRefreshing.asStateFlow()

    private val _showForm = MutableStateFlow(false)
    val showForm = _showForm.asStateFlow()

    private val _formEnabled = MutableStateFlow(true)
    val formEnabled = _formEnabled.asStateFlow()

    var formState by mutableStateOf(AddSpotFormState())
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
            spotRepo
                .getSpots()
                .onSuccess { _spots.value = it }
                .recoverRequestErrors(
                    "Error retrieving parking spots",
                    { errorToForm(it) },
                    snackbarState,
                    viewModelScope,
                )
            _isRefreshing.value = false
        }
    }

    @Suppress("detekt:LongMethod")
    fun onAddSpotClick() {
        viewModelScope.launch {
            _formEnabled.value = false
            spotRepo
                .createSpot(
                    Spot(
                        availability =
                            formState.selectedIds.value
                                .sortedBy { it }
                                .map {
                                    TimeUnit(
                                        startTime = it.toLocalDateTime(),
                                        endTime = (it + 1).toLocalDateTime(),
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
                        pricePerHour = formState.pricePerHour.value.toDoubleOrNull() ?: -0.0,
                    )
                )
                .also { clearFieldErrors() }
                .onSuccess {
                    onRefresh()
                    onHideForm()
                }
                .recoverRequestErrors(
                    "Error adding parking spot",
                    { errorToForm(it) },
                    snackbarState,
                    viewModelScope,
                )
            _formEnabled.value = true
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
                            error =
                                if (value == "") "Price cannot be empty"
                                else if (value.toDoubleOrNull() == null) "Price must be a number"
                                else null,
                        )
                )
            }
    }

    fun onAddTime(elements: Iterable<Int>) {
        formState =
            formState.run {
                copy(
                    selectedIds =
                        selectedIds.copy(value = formState.selectedIds.value.plus(elements))
                )
            }
    }

    fun onRemoveTime(elements: Iterable<Int>) {
        formState =
            formState.run {
                copy(
                    selectedIds =
                        selectedIds.copy(value = formState.selectedIds.value.minus(elements))
                )
            }
    }

    fun resetForm() {
        formState =
            formState.run {
                copy(
                    streetAddress = FieldState(""),
                    city = FieldState(""),
                    state = FieldState(""),
                    countryCode = FieldState("CA"),
                    postalCode = FieldState(""),
                    chargingStation = FieldState(false),
                    plugIn = FieldState(false),
                    shelter = FieldState(false),
                    pricePerHour = FieldState(""),
                    selectedIds = FieldState(emptySet()),
                    disabledIds =
                        FieldState(
                            if (now() > startOfNextAvailableDayInstant())
                                (0..now().toIndex()).toSet()
                            else emptySet()
                        ),
                )
            }
    }

    private fun errorToForm(error: ErrorModel) {
        when (error.type) {
            else -> annotateErrorLocation(error.errors)
        }
    }

    @Suppress("detekt:CyclomaticComplexMethod")
    private fun annotateErrorLocation(errors: List<ErrorDetail>) {
        for (err in errors) {
            when (err.location) {
                "body",
                "body.location" ->
                    formState =
                        formState.run {
                            copy(
                                streetAddress =
                                    streetAddress.copy(error = "Invalid street address"),
                                city = city.copy(error = "Invalid city"),
                                state = state.copy(error = "Invalid state"),
                                countryCode = countryCode.copy(error = "Invalid country code"),
                                postalCode = postalCode.copy(error = "Invalid postal code"),
                                pricePerHour = pricePerHour.copy(error = "Invalid price"),
                            )
                        }

                "body.location.street_address" ->
                    formState =
                        formState.run {
                            copy(
                                streetAddress = streetAddress.copy(error = "Invalid street address")
                            )
                        }
                "body.location.city" ->
                    formState = formState.run { copy(city = city.copy(error = "Invalid city")) }
                "body.location.state" ->
                    formState = formState.run { copy(state = state.copy(error = "Invalid state")) }

                "body.location.country" ->
                    formState =
                        formState.run {
                            copy(countryCode = countryCode.copy(error = "Invalid country"))
                        }

                "body.location.postal_code" ->
                    formState =
                        formState.run {
                            copy(postalCode = postalCode.copy(error = "Invalid postal code"))
                        }

                "body.price_per_hour" ->
                    formState =
                        formState.run {
                            copy(pricePerHour = pricePerHour.copy(error = "Invalid price"))
                        }
            }
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
                    postalCode = postalCode.copy(error = null),
                    pricePerHour = pricePerHour.copy(error = null),
                )
            }
    }

    fun createHandler() =
        AddSpotFormHandler(
            onStreetAddressChange = this::onStreetAddressChange,
            onCityChange = this::onCityChange,
            onStateChange = this::onStateChange,
            onCountryCodeChange = this::onCountryCodeChange,
            onPostalCodeChange = this::onPostalCodeChange,
            onChargingStationChange = this::onChargingStationChange,
            onPlugInChange = this::onPlugInChange,
            onShelterChange = this::onShelterChange,
            onPricePerHourChange = this::onPricePerHourChange,
            onAddTime = this::onAddTime,
            onRemoveTime = this::onRemoveTime,
            onAddSpotClick = this::onAddSpotClick,
            resetForm = this::resetForm,
        )
}

@Composable
fun rememberAddSpotFormHandler(viewModel: SpotsViewModel) =
    remember(viewModel) { viewModel.createHandler() }

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
    val selectedIds: FieldState<Set<Int>> = FieldState(emptySet()),
    val disabledIds: FieldState<Set<Int>> = FieldState(emptySet()),
)

data class AddSpotFormHandler(
    val onStreetAddressChange: (String) -> Unit = {},
    val onCityChange: (String) -> Unit = {},
    val onStateChange: (String) -> Unit = {},
    val onCountryCodeChange: (String) -> Unit = {},
    val onPostalCodeChange: (String) -> Unit = {},
    val onChargingStationChange: (Boolean) -> Unit = {},
    val onPlugInChange: (Boolean) -> Unit = {},
    val onShelterChange: (Boolean) -> Unit = {},
    val onPricePerHourChange: (String) -> Unit = {},
    val onAddTime: (Iterable<Int>) -> Unit = {},
    val onRemoveTime: (Iterable<Int>) -> Unit = {},
    val onAddSpotClick: () -> Unit = {},
    val resetForm: () -> Unit = {},
)
