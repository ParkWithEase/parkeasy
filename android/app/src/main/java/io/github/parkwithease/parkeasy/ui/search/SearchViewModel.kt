package io.github.parkwithease.parkeasy.ui.search

import android.location.Location
import android.util.Log
import android.util.SparseArray
import androidx.compose.material3.SnackbarHostState
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.core.util.containsKey
import androidx.core.util.forEach
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.maplibre.compose.StaticLocationEngine
import dagger.hilt.android.lifecycle.HiltViewModel
import io.github.parkwithease.parkeasy.data.local.AuthRepository
import io.github.parkwithease.parkeasy.data.remote.CarRepository
import io.github.parkwithease.parkeasy.data.remote.SpotRepository
import io.github.parkwithease.parkeasy.domain.GetLocationUseCase
import io.github.parkwithease.parkeasy.model.Booking
import io.github.parkwithease.parkeasy.model.Car
import io.github.parkwithease.parkeasy.model.FieldState
import io.github.parkwithease.parkeasy.model.Spot
import io.github.parkwithease.parkeasy.model.TimeUnit
import io.github.parkwithease.parkeasy.ui.common.recoverRequestErrors
import io.github.parkwithease.parkeasy.ui.common.startOfNextAvailableDayInstant
import io.github.parkwithease.parkeasy.ui.common.timezone
import io.github.parkwithease.parkeasy.ui.common.toIndex
import io.github.parkwithease.parkeasy.ui.spots.NumSlots
import javax.inject.Inject
import kotlin.collections.plus
import kotlinx.coroutines.channels.BufferOverflow
import kotlinx.coroutines.flow.MutableSharedFlow
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.emitAll
import kotlinx.coroutines.flow.transform
import kotlinx.coroutines.launch
import kotlinx.datetime.Clock.System.now
import kotlinx.datetime.toInstant

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
    getLocation: GetLocationUseCase,
) : ViewModel() {
    val loggedIn = authRepo.statusFlow
    val snackbarState = SnackbarHostState()

    private var times = SparseArray<TimeUnit>()

    private val _spots = MutableStateFlow(emptyList<Spot>())
    val spots = _spots.asStateFlow()

    private val _cars = MutableStateFlow(emptyList<Car>())
    val cars = _cars.asStateFlow()

    private val _isRefreshing = MutableStateFlow(false)
    val isRefreshing = _isRefreshing.asStateFlow()

    private val _showForm = MutableStateFlow(false)
    val showForm = _showForm.asStateFlow()

    private val _formEnabled = MutableStateFlow(true)
    val formEnabled = _formEnabled.asStateFlow()

    private val locationTrigger =
        MutableSharedFlow<Unit>(
            extraBufferCapacity = 1,
            onBufferOverflow = BufferOverflow.DROP_OLDEST,
        )
    private val locationFlow = getLocation()
    val currentLocation = locationTrigger.transform { emitAll(locationFlow) }
    val engine =
        StaticLocationEngine().apply {
            lastLocation =
                Location("static").apply {
                    latitude = DefaultLatitude
                    longitude = DefaultLongitude
                }
        }

    init {
        viewModelScope.launch { currentLocation.collect { engine.lastLocation = it } }
    }

    fun startLocationFlow() {
        locationTrigger.tryEmit(Unit)
    }

    var createState by mutableStateOf(CreateState())
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
                .getSpotsAround(
                    engine.lastLocation?.latitude ?: DefaultLatitude,
                    engine.lastLocation?.longitude ?: DefaultLongitude,
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

    fun refreshAvailability() {
        viewModelScope.launch {
            spotRepo
                .getSpotAvailability(createState.selectedSpot.value.id)
                .onSuccess {
                    times = SparseArray<TimeUnit>()
                    it.forEach { timeUnit ->
                        val index = timeUnit.startTime.toInstant(timezone()).toIndex()
                        if (index > 0 && index < NumSlots && timeUnit.status != TimeUnit.BOOKED) {
                            times.put(index, timeUnit)
                        }
                    }
                    updateDisabled()
                }
                .recoverRequestErrors(
                    "Error retrieving availabilities",
                    {},
                    snackbarState,
                    viewModelScope,
                )
        }
    }

    fun updateDisabled() {
        createState =
            createState.run {
                copy(
                    disabledIds =
                        FieldState(
                            (if (now() > startOfNextAvailableDayInstant())
                                    (0..now().toIndex()).toSet()
                                else emptySet())
                                .plus((0..NumSlots - 1).filter { !times.containsKey(it) })
                        )
                )
            }
    }

    fun onCreateBookingClick() {
        viewModelScope.launch {
            var bookedTimes = emptyList<TimeUnit>()
            times.forEach { k, v ->
                if (k in createState.selectedIds.value) bookedTimes = bookedTimes.plus(v)
            }
            Log.e("", bookedTimes.toString())
            spotRepo
                .bookSpot(
                    spotId = createState.selectedSpot.value.id,
                    booking =
                        Booking(carId = createState.selectedCar.value.id, bookedTimes = bookedTimes),
                )
                .onSuccess {
                    onRefresh()
                    onHideForm()
                }
                .onFailure { onHideForm() }
                .recoverRequestErrors(
                    "Error booking parking spot",
                    {},
                    snackbarState,
                    viewModelScope,
                )
        }
    }

    fun onCarChange(value: Car) {
        createState = createState.run { copy(selectedCar = selectedCar.copy(value = value)) }
    }

    fun onSpotChange(value: Spot) {
        createState = createState.run { copy(selectedSpot = selectedSpot.copy(value = value)) }
        refreshAvailability()
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
        updateDisabled()
    }

    fun createCreateHandler() =
        CreateHandler(
            onCarChange = this::onCarChange,
            onSpotChange = this::onSpotChange,
            onAddTime = this::onAddTime,
            onRemoveTime = this::onRemoveTime,
            onCreateBookingClick = this::onCreateBookingClick,
            reset = this::resetCreate,
        )
}

@Composable
fun rememberCreateHandler(viewModel: SearchViewModel) =
    remember(viewModel) { viewModel.createCreateHandler() }

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
