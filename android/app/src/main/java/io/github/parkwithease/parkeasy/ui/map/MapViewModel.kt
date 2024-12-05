package io.github.parkwithease.parkeasy.ui.map

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.maplibre.compose.StaticLocationEngine
import dagger.hilt.android.lifecycle.HiltViewModel
import io.github.parkwithease.parkeasy.data.local.AuthRepository
import io.github.parkwithease.parkeasy.domain.GetLocationUseCase
import javax.inject.Inject
import kotlinx.coroutines.channels.BufferOverflow
import kotlinx.coroutines.flow.MutableSharedFlow
import kotlinx.coroutines.flow.emitAll
import kotlinx.coroutines.flow.transform
import kotlinx.coroutines.launch

@HiltViewModel
class MapViewModel @Inject constructor(authRepo: AuthRepository, getLocation: GetLocationUseCase) :
    ViewModel() {
    val loggedIn = authRepo.statusFlow
    private val locationTrigger =
        MutableSharedFlow<Unit>(extraBufferCapacity = 1, onBufferOverflow = BufferOverflow.DROP_OLDEST)
    private val locationFlow = getLocation()
    val currentLocation = locationTrigger.transform { emitAll(locationFlow) }
    val engine = StaticLocationEngine()

    init {
        viewModelScope.launch { currentLocation.collect { engine.lastLocation = it } }
    }

    fun startLocationFlow() {
        locationTrigger.tryEmit(Unit)
    }
}
