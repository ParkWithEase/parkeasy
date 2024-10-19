package io.github.parkwithease.parkeasy.ui.login

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import io.github.parkwithease.parkeasy.data.local.AuthRepository
import io.github.parkwithease.parkeasy.data.remote.UserRepository
import io.github.parkwithease.parkeasy.model.LoginCredentials
import io.github.parkwithease.parkeasy.model.RegistrationCredentials
import io.github.parkwithease.parkeasy.model.ResetCredentials
import io.github.parkwithease.parkeasy.ui.SnackbarController
import io.github.parkwithease.parkeasy.ui.SnackbarEvent
import javax.inject.Inject
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch

@HiltViewModel
class LoginViewModel
@Inject
constructor(authRepo: AuthRepository, private val userRepo: UserRepository) : ViewModel() {
    val loggedIn = authRepo.statusFlow

    private val _formEnabled = MutableStateFlow(true)
    val formEnabled = _formEnabled.asStateFlow()

    fun onLoginPress(email: String, password: String) {
        viewModelScope.launch {
            _formEnabled.value = false
            if (userRepo.login(LoginCredentials(email, password))) {
                SnackbarController.sendEvent(
                    event = SnackbarEvent(message = "Logged in successfully")
                )
            } else {
                SnackbarController.sendEvent(event = SnackbarEvent(message = "Error logging in"))
            }
            _formEnabled.value = true
        }
    }

    fun onRegisterPress(name: String, email: String, password: String, confirmPassword: String) {
        if (password == confirmPassword) {
            viewModelScope.launch {
                _formEnabled.value = false
                if (userRepo.register(RegistrationCredentials(name, email, password))) {
                    SnackbarController.sendEvent(
                        event = SnackbarEvent(message = "Registered successfully")
                    )
                } else {
                    SnackbarController.sendEvent(
                        event = SnackbarEvent(message = "Error registering")
                    )
                }
                _formEnabled.value = true
            }
        } else {
            viewModelScope.launch {
                SnackbarController.sendEvent(
                    event = SnackbarEvent(message = "Passwords don't match")
                )
            }
        }
    }

    fun onRequestResetPress(email: String) {
        viewModelScope.launch {
            if (userRepo.requestReset(ResetCredentials(email))) {
                SnackbarController.sendEvent(
                    event = SnackbarEvent(message = "Reset email sent\nJk... we're working on it")
                )
            } else {
                SnackbarController.sendEvent(
                    event = SnackbarEvent(message = "Error resetting password")
                )
            }
        }
    }
}
