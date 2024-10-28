package io.github.parkwithease.parkeasy.ui.cars

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.assisted.Assisted
import dagger.assisted.AssistedFactory
import dagger.assisted.AssistedInject
import dagger.hilt.android.lifecycle.HiltViewModel
import io.github.parkwithease.parkeasy.data.remote.CarRepository
import io.github.parkwithease.parkeasy.model.Car
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch

@HiltViewModel(assistedFactory = CarsViewModel.Factory::class)
class CarsViewModel
@AssistedInject
constructor(
    private val carRepo: CarRepository,
    @Assisted val showSnackbar: suspend (String, String?) -> Boolean,
) : ViewModel() {
    @AssistedFactory
    interface Factory {
        fun create(showSnackbar: suspend (String, String?) -> Boolean): CarsViewModel
    }

    private val _cars = MutableStateFlow(emptyList<Car>())
    val cars = _cars.asStateFlow()

    private val _isRefreshing = MutableStateFlow(false)
    val isRefreshing = _isRefreshing.asStateFlow()

    fun onRefresh() {
        viewModelScope.launch {
            _isRefreshing.value = true
            _cars.value = carRepo.getCars()
            _isRefreshing.value = false
        }
    }
}
