package io.github.parkwithease.parkeasy.ui.login

import androidx.lifecycle.ViewModel
import dagger.hilt.android.lifecycle.HiltViewModel
import io.github.parkwithease.parkeasy.data.local.AuthRepository
import io.github.parkwithease.parkeasy.data.remote.UserRepository
import io.github.parkwithease.parkeasy.di.IoDispatcher
import io.github.parkwithease.parkeasy.model.Event
import io.github.parkwithease.parkeasy.model.LoginCredentials
import io.github.parkwithease.parkeasy.model.RegistrationCredentials
import io.github.parkwithease.parkeasy.model.ResetCredentials
import javax.inject.Inject
import kotlinx.coroutines.CoroutineDispatcher
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch

@HiltViewModel
class LoginViewModel
@Inject
constructor(
    authRepo: AuthRepository,
    private val userRepo: UserRepository,
    @IoDispatcher private val ioDispatcher: CoroutineDispatcher = Dispatchers.IO,
) : ViewModel() {
    val loggedIn = authRepo.statusFlow

    private val _formEnabled = MutableStateFlow(true)
    val formEnabled = _formEnabled.asStateFlow()

    private val _message = MutableStateFlow(Event.initial(""))
    val message = _message.asStateFlow()

    fun onLoginPress(email: String, password: String) {
        CoroutineScope(ioDispatcher).launch {
            _formEnabled.value = false
            if (userRepo.login(LoginCredentials(email, password))) {
                _message.value = Event("Logged in successfully")
            } else {
                _message.value = Event("Error logging in")
            }
            _formEnabled.value = true
        }
    }

    fun onRegisterPress(name: String, email: String, password: String, confirmPassword: String) {
        if (password == confirmPassword) {
            CoroutineScope(ioDispatcher).launch {
                _formEnabled.value = false
                if (userRepo.register(RegistrationCredentials(name, email, password))) {
                    _message.value = Event("Registered successfully")
                } else {
                    _message.value = Event("Error registering")
                }
                _formEnabled.value = true
            }
        } else {
            _message.value = Event("Passwords don't match")
        }
    }

    fun onRequestResetPress(email: String) {
        CoroutineScope(ioDispatcher).launch {
            if (userRepo.requestReset(ResetCredentials(email))) {
                _message.value = Event("Reset email sent\nJk... we're working on it")
            } else {
                _message.value = Event("Error resetting password")
            }
        }
    }
}
