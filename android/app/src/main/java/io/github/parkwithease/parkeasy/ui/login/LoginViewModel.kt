package io.github.parkwithease.parkeasy.ui.login

import androidx.lifecycle.ViewModel
import dagger.hilt.android.lifecycle.HiltViewModel
import io.github.parkwithease.parkeasy.data.ParkEasyRepository
import io.github.parkwithease.parkeasy.model.Credentials
import javax.inject.Inject

@HiltViewModel
class LoginViewModel @Inject constructor(private val repo: ParkEasyRepository) : ViewModel() {

    suspend fun login(credentials: Credentials) {
        return repo.login(credentials)
    }
}
