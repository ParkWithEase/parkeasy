package io.github.parkwithease.parkeasy.ui.login

import androidx.lifecycle.ViewModel
import dagger.hilt.android.lifecycle.HiltViewModel
import io.github.parkwithease.parkeasy.data.ParkEasyRepository
import io.github.parkwithease.parkeasy.model.Credentials
import javax.inject.Inject
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.runBlocking

@HiltViewModel
class LoginViewModel @Inject constructor(private val repo: ParkEasyRepository) : ViewModel() {
    private val _email = MutableStateFlow("")
    val email = _email.asStateFlow()

    private val _password = MutableStateFlow("")
    val password = _password.asStateFlow()

    suspend fun login(credentials: Credentials) {
        return repo.login(credentials)
    }

    fun onEmailChange(input: String) {
        _email.value = input
    }

    fun onPasswordChange(input: String) {
        _password.value = input
    }

    fun onLoginPress() {
        runBlocking { login(Credentials(_email.value, _password.value)) }
    }
}
