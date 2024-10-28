package io.github.parkwithease.parkeasy.ui.profile

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.assisted.Assisted
import dagger.assisted.AssistedFactory
import dagger.assisted.AssistedInject
import dagger.hilt.android.lifecycle.HiltViewModel
import io.github.parkwithease.parkeasy.data.local.AuthRepository
import io.github.parkwithease.parkeasy.data.remote.UserRepository
import io.github.parkwithease.parkeasy.model.Profile
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch

@HiltViewModel(assistedFactory = ProfileViewModel.Factory::class)
class ProfileViewModel
@AssistedInject
constructor(
    authRepo: AuthRepository,
    private val userRepo: UserRepository,
    @Assisted val showSnackbar: suspend (String, String?) -> Boolean,
) : ViewModel() {
    @AssistedFactory
    interface Factory {
        fun create(showSnackbar: suspend (String, String?) -> Boolean): ProfileViewModel
    }

    private val _profile = MutableStateFlow(Profile("", ""))
    val profile = _profile.asStateFlow()

    val loggedIn = authRepo.statusFlow

    fun refresh() {
        viewModelScope.launch {
            val profile = userRepo.getUser()
            if (profile != null) {
                _profile.value = profile
            }
        }
    }

    fun onLogoutClick() {
        viewModelScope.launch {
            showSnackbar(
                if (userRepo.logout()) "Logged out successfully" else "Error logging out",
                null,
            )
        }
    }
}
