package io.github.parkwithease.parkeasy.ui.profile

import androidx.lifecycle.ViewModel
import dagger.hilt.android.lifecycle.HiltViewModel
import io.github.parkwithease.parkeasy.data.local.AuthRepository
import io.github.parkwithease.parkeasy.data.remote.UserRepository
import javax.inject.Inject
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import kotlinx.coroutines.runBlocking

@HiltViewModel
class ProfileViewModel
@Inject
constructor(private val authRepo: AuthRepository, private val userRepo: UserRepository) :
    ViewModel() {
    private val _loggedIn = MutableStateFlow(false)
    val loggedIn = _loggedIn.asStateFlow()

    init {
        runBlocking { launch { _loggedIn.value = authRepo.getStatus() } }
    }

    private fun logout() {
        runBlocking { launch { userRepo.logout() } }
        _loggedIn.value = false
    }

    fun onLogoutPress() {
        logout()
    }
}
