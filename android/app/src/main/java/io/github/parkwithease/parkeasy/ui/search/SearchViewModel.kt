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
import io.github.parkwithease.parkeasy.data.remote.SpotRepository
import io.github.parkwithease.parkeasy.model.Car
import io.github.parkwithease.parkeasy.model.FieldState
import io.github.parkwithease.parkeasy.model.Spot
import io.github.parkwithease.parkeasy.ui.common.recoverRequestErrors
import javax.inject.Inject
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch

// Somewhere above Winnipeg
const val DefaultLatitude = 49.895077
const val DefaultLongitude = -97.138451
private const val DefaultDistance = 25000

@Suppress("detekt:TooManyFunctions")
@HiltViewModel
class SearchViewModel
@Inject
constructor(authRepo: AuthRepository, private val spotRepo: SpotRepository) : ViewModel() {
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
            _isRefreshing.value = false
        }
    }

    fun onLatitudeChange(value: Double) {
        searchState = searchState.run { copy(latitude = latitude.copy(value = value)) }
    }

    fun onLongitudeChange(value: Double) {
        searchState = searchState.run { copy(longitude = longitude.copy(value = value)) }
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

    fun createSearchHandler() =
        SearchHandler(
            onLatitudeChange = this::onLatitudeChange,
            onLongitudeChange = this::onLongitudeChange,
            reset = this::resetSearch,
        )
}

@Composable
fun rememberSearchHandler(viewModel: SearchViewModel) =
    remember(viewModel) { viewModel.createSearchHandler() }

data class SearchState(
    val latitude: FieldState<Double> = FieldState(DefaultLatitude),
    val longitude: FieldState<Double> = FieldState(DefaultLongitude),
)

data class SearchHandler(
    val onLatitudeChange: (Double) -> Unit,
    val onLongitudeChange: (Double) -> Unit,
    val reset: () -> Unit,
)
