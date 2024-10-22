package io.github.parkwithease.parkeasy.ui.profile

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
import kotlinx.coroutines.runBlocking

@HiltViewModel
class ProfileViewModel
@Inject
constructor(authRepo: AuthRepository, private val userRepo: UserRepository) : ViewModel() {
    private val _profile = MutableStateFlow(Profile("", "", ""))
    val profile = _profile.asStateFlow()
        get() {
            refresh()
            return field
        }

    val loggedIn = authRepo.statusFlow

    private fun refresh() {
        viewModelScope.launch {
            val profile = userRepo.getUser()
            if (profile != null) {
                _profile.value = profile
            }
        }
    }

    fun onLogoutClick() {
        runBlocking { userRepo.logout() }
    }
}
