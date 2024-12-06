package io.github.parkwithease.parkeasy.ui.profile

import androidx.compose.material3.SnackbarHostState
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import io.github.parkwithease.parkeasy.data.local.AuthRepository
import io.github.parkwithease.parkeasy.data.remote.UserRepository
import io.github.parkwithease.parkeasy.model.Profile
import javax.inject.Inject
import kotlinx.coroutines.channels.BufferOverflow
import kotlinx.coroutines.flow.MutableSharedFlow
import kotlinx.coroutines.flow.SharingStarted
import kotlinx.coroutines.flow.map
import kotlinx.coroutines.flow.onStart
import kotlinx.coroutines.flow.stateIn
import kotlinx.coroutines.launch

@HiltViewModel
class ProfileViewModel
@Inject
constructor(authRepo: AuthRepository, private val userRepo: UserRepository) : ViewModel() {
    val loggedIn = authRepo.statusFlow
    val snackbarState = SnackbarHostState()

    private val refreshTrigger =
        MutableSharedFlow<Unit>(
            extraBufferCapacity = 1,
            onBufferOverflow = BufferOverflow.DROP_LATEST,
        )
    val profile =
        refreshTrigger
            .onStart { emit(Unit) }
            .map { userRepo.getUser().getOrDefault(Profile()) }
            .stateIn(viewModelScope, SharingStarted.Lazily, Profile("", ""))

    fun refresh() {
        refreshTrigger.tryEmit(Unit)
    }

    fun onLogoutClick() =
        viewModelScope.launch {
            userRepo.logout().onSuccess { snackbarState.showSnackbar("Logged out successfully") }
        }
}
