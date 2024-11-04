package io.github.parkwithease.parkeasy.ui.spots

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import io.github.parkwithease.parkeasy.data.remote.SpotRepository
import io.github.parkwithease.parkeasy.model.Spot
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

    fun onRefresh() {
        viewModelScope.launch {
            _isRefreshing.value = true
            _spots.value = spotRepo.getSpots()
            _isRefreshing.value = false
        }
    }
}
