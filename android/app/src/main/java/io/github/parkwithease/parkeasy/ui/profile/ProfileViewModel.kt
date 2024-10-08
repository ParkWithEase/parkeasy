package io.github.parkwithease.parkeasy.ui.profile

import androidx.lifecycle.ViewModel
import dagger.hilt.android.lifecycle.HiltViewModel
import io.github.parkwithease.parkeasy.data.local.AuthRepository
import io.github.parkwithease.parkeasy.data.remote.UserRepository
import javax.inject.Inject
import kotlinx.coroutines.runBlocking

@HiltViewModel
class ProfileViewModel
@Inject
constructor(authRepo: AuthRepository, private val userRepo: UserRepository) : ViewModel() {
    val loggedIn = authRepo.statusFlow

    private fun logout() {
        runBlocking { userRepo.logout() }
    }

    fun onLogoutPress() {
        logout()
    }
}
