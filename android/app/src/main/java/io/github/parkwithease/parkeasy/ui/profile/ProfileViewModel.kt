package io.github.parkwithease.parkeasy.ui.profile

import androidx.compose.material3.SnackbarHostState
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import io.github.parkwithease.parkeasy.data.local.AuthRepository
import io.github.parkwithease.parkeasy.data.remote.UserRepository
import io.github.parkwithease.parkeasy.model.Profile
import javax.inject.Inject
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch

@HiltViewModel
class ProfileViewModel
@Inject
constructor(authRepo: AuthRepository, private val userRepo: UserRepository) : ViewModel() {
    val loggedIn = authRepo.statusFlow
    val snackbarState = SnackbarHostState()

    private val _profile = MutableStateFlow(Profile("", ""))
    val profile = _profile.asStateFlow()

    fun refresh() =
        viewModelScope.launch {
            userRepo
                .getUser()
                .onSuccess { _profile.value = it }
                .onFailure {
                    viewModelScope.launch { snackbarState.showSnackbar("Error retrieving profile") }
                }
        }

    fun onLogoutClick() {
        viewModelScope.launch {
            snackbarState.showSnackbar(
                if (userRepo.logout()) "Logged out successfully" else "Error logging out"
            )
        }
    }
}
