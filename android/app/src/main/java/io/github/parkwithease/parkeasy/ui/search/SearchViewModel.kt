package io.github.parkwithease.parkeasy.ui.search

import androidx.compose.material3.SnackbarHostState
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import io.github.parkwithease.parkeasy.data.local.AuthRepository
import io.github.parkwithease.parkeasy.data.remote.CarRepository
import io.github.parkwithease.parkeasy.data.remote.SpotRepository
import io.github.parkwithease.parkeasy.model.Car
import io.github.parkwithease.parkeasy.model.FieldState
import io.github.parkwithease.parkeasy.model.Spot
import io.github.parkwithease.parkeasy.ui.common.recoverRequestErrors
import io.github.parkwithease.parkeasy.ui.common.startOfNextAvailableDayInstant
import io.github.parkwithease.parkeasy.ui.common.toIndex
import javax.inject.Inject
import kotlin.collections.plus
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import kotlinx.datetime.Clock.System.now

// Somewhere above Winnipeg
const val DefaultLatitude = 49.895077
const val DefaultLongitude = -97.138451
private const val DefaultDistance = 25000

@Suppress("detekt:TooManyFunctions")
@HiltViewModel
class SearchViewModel
@Inject
constructor(
    authRepo: AuthRepository,
    private val spotRepo: SpotRepository,
    private val carRepo: CarRepository,
) : ViewModel() {
    val loggedIn = authRepo.statusFlow
    val snackbarState = SnackbarHostState()

    private val _spots = MutableStateFlow(emptyList<Spot>())
    val spots = _spots.asStateFlow()

    private val _cars = MutableStateFlow(emptyList<Car>())
    val cars = _cars.asStateFlow()

    private val _isRefreshing = MutableStateFlow(false)
    val isRefreshing = _isRefreshing.asStateFlow()

    var searchState by mutableStateOf(SearchState())
        private set

    var createState by mutableStateOf(CreateState())
        private set

    fun onRefresh() {
        viewModelScope.launch {
            _isRefreshing.value = true
            spotRepo
                .getSpotsAround(
                    searchState.latitude.value,
                    searchState.longitude.value,
                    DefaultDistance,
                )
                .onSuccess { _spots.value = it }
                .recoverRequestErrors(
                    "Error retrieving parking spots",
                    {},
                    snackbarState,
                    viewModelScope,
                )
            carRepo
                .getCars()
                .onSuccess { _cars.value = it }
                .recoverRequestErrors("Error retrieving cars", {}, snackbarState, viewModelScope)
            _isRefreshing.value = false
        }
    }

    fun onCarChange(value: Car) {
        createState = createState.run { copy(selectedCar = selectedCar.copy(value = value)) }
    }

    fun onSpotChange(value: Spot) {
        createState = createState.run { copy(selectedSpot = selectedSpot.copy(value = value)) }
    }

    fun onLatitudeChange(value: Double) {
        searchState = searchState.run { copy(latitude = latitude.copy(value = value)) }
    }

    fun onLongitudeChange(value: Double) {
        searchState = searchState.run { copy(longitude = longitude.copy(value = value)) }
    }

    fun onAddTime(elements: Iterable<Int>) {
        createState =
            createState.run {
                copy(
                    selectedIds =
                        selectedIds.copy(value = createState.selectedIds.value.plus(elements))
                )
            }
    }

    fun onRemoveTime(elements: Iterable<Int>) {
        createState =
            createState.run {
                copy(
                    selectedIds =
                        selectedIds.copy(value = createState.selectedIds.value.minus(elements))
                )
            }
    }

    fun resetSearch() {
        searchState =
            searchState.run {
                copy(
                    latitude = FieldState(DefaultLatitude),
                    longitude = FieldState(DefaultLongitude),
                )
            }
    }

    fun resetCreate() {
        createState =
            createState.run {
                copy(
                    selectedCar = FieldState(Car()),
                    selectedSpot = FieldState(Spot()),
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

    fun createSearchHandler() =
        SearchHandler(
            onLatitudeChange = this::onLatitudeChange,
            onLongitudeChange = this::onLongitudeChange,
            reset = this::resetSearch,
        )

    fun createCreateHandler() =
        CreateHandler(
            onCarChange = this::onCarChange,
            onSpotChange = this::onSpotChange,
            onAddTime = this::onAddTime,
            onRemoveTime = this::onRemoveTime,
            onCreateBookingClick = {},
            reset = this::resetCreate,
        )
}

@Composable
fun rememberSearchHandler(viewModel: SearchViewModel) =
    remember(viewModel) { viewModel.createSearchHandler() }

@Composable
fun rememberCreateHandler(viewModel: SearchViewModel) =
    remember(viewModel) { viewModel.createCreateHandler() }

data class SearchState(
    val latitude: FieldState<Double> = FieldState(DefaultLatitude),
    val longitude: FieldState<Double> = FieldState(DefaultLongitude),
)

data class SearchHandler(
    val onLatitudeChange: (Double) -> Unit,
    val onLongitudeChange: (Double) -> Unit,
    val reset: () -> Unit,
)

data class CreateState(
    val selectedCar: FieldState<Car> = FieldState(Car()),
    val selectedSpot: FieldState<Spot> = FieldState(Spot()),
    val selectedIds: FieldState<Set<Int>> = FieldState(emptySet()),
    val disabledIds: FieldState<Set<Int>> = FieldState(emptySet()),
)

data class CreateHandler(
    val onCarChange: (Car) -> Unit,
    val onSpotChange: (Spot) -> Unit,
    val onAddTime: (Iterable<Int>) -> Unit,
    val onRemoveTime: (Iterable<Int>) -> Unit,
    val onCreateBookingClick: () -> Unit,
    val reset: () -> Unit,
)
