package io.github.parkwithease.parkeasy.ui.bookings

import androidx.compose.material3.SnackbarHostState
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import io.github.parkwithease.parkeasy.data.remote.SpotRepository
import io.github.parkwithease.parkeasy.model.BookingHistory
import io.github.parkwithease.parkeasy.ui.common.recoverRequestErrors
import javax.inject.Inject
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch

@Suppress("detekt:TooManyFunctions")
@HiltViewModel
class BookingsViewModel @Inject constructor(private val spotRepo: SpotRepository) : ViewModel() {
    val snackbarState = SnackbarHostState()

    private val _bookings = MutableStateFlow(emptyList<BookingHistory>())
    val bookings = _bookings.asStateFlow()

    private val _isRefreshing = MutableStateFlow(false)
    val isRefreshing = _isRefreshing.asStateFlow()

    fun onRefresh() {
        viewModelScope.launch {
            _isRefreshing.value = true
            spotRepo
                .getBookings()
                .onSuccess { _bookings.value = it }
                .recoverRequestErrors("Error booking history", {}, snackbarState, viewModelScope)
            _isRefreshing.value = false
        }
    }
}
