package io.github.parkwithease.parkeasy.ui.login

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.assisted.Assisted
import dagger.assisted.AssistedFactory
import dagger.assisted.AssistedInject
import dagger.hilt.android.lifecycle.HiltViewModel
import io.github.parkwithease.parkeasy.data.local.AuthRepository
import io.github.parkwithease.parkeasy.data.remote.UserRepository
import io.github.parkwithease.parkeasy.model.LoginCredentials
import io.github.parkwithease.parkeasy.model.RegistrationCredentials
import io.github.parkwithease.parkeasy.model.ResetCredentials
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch

@HiltViewModel(assistedFactory = LoginViewModel.Factory::class)
class LoginViewModel
@AssistedInject
constructor(
    authRepo: AuthRepository,
    private val userRepo: UserRepository,
    @Assisted val showSnackbar: suspend (String, String?) -> Boolean,
) : ViewModel() {
    @AssistedFactory
    interface Factory {
        fun create(showSnackbar: suspend (String, String?) -> Boolean): LoginViewModel
    }

    val loggedIn = authRepo.statusFlow

    private val _formEnabled = MutableStateFlow(true)
    val formEnabled = _formEnabled.asStateFlow()

    fun onLoginPress(email: String, password: String) {
        viewModelScope.launch {
            _formEnabled.value = false
            val result: Boolean = userRepo.login(LoginCredentials(email, password))
            _formEnabled.value = true
            showSnackbar(if (result) "Logged in successfully" else "Error logging in", null)
        }
    }

    fun onRegisterPress(name: String, email: String, password: String, confirmPassword: String) {
        if (password == confirmPassword) {
            viewModelScope.launch {
                _formEnabled.value = false
                val result: Boolean =
                    userRepo.register(RegistrationCredentials(name, email, password))
                _formEnabled.value = true
                showSnackbar(if (result) "Registered successfully" else "Error registering", null)
            }
        } else {
            viewModelScope.launch { showSnackbar("Passwords don't match", null) }
        }
    }

    fun onRequestResetPress(email: String) {
        viewModelScope.launch {
            val result = userRepo.requestReset(ResetCredentials(email))
            showSnackbar(
                if (result) "Reset email sent\nJk... we're working on it"
                else "Error resetting password",
                null,
            )
        }
    }
}
