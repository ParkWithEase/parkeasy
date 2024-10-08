package io.github.parkwithease.parkeasy.ui.login

import androidx.lifecycle.ViewModel
import dagger.hilt.android.lifecycle.HiltViewModel
import io.github.parkwithease.parkeasy.data.local.AuthRepository
import io.github.parkwithease.parkeasy.data.remote.UserRepository
import io.github.parkwithease.parkeasy.model.Event
import io.github.parkwithease.parkeasy.model.LoginCredentials
import io.github.parkwithease.parkeasy.model.RegistrationCredentials
import io.github.parkwithease.parkeasy.model.ResetCredentials
import javax.inject.Inject
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import kotlinx.coroutines.runBlocking

@HiltViewModel
class LoginViewModel
@Inject
constructor(authRepo: AuthRepository, private val userRepo: UserRepository) : ViewModel() {
    private val _name = MutableStateFlow("")
    val name = _name.asStateFlow()

    private val _email = MutableStateFlow("")
    val email = _email.asStateFlow()

    private val _password = MutableStateFlow("")
    val password = _password.asStateFlow()

    private val _confirmPassword = MutableStateFlow("")
    val confirmPassword = _confirmPassword.asStateFlow()

    private val _matchingPasswords = MutableStateFlow(true)
    val matchingPasswords = _matchingPasswords.asStateFlow()

    private val _registering = MutableStateFlow(false)
    val registering = _registering.asStateFlow()

    private val _requestingReset = MutableStateFlow(false)
    val requestingReset = _requestingReset.asStateFlow()

    val loggedIn = authRepo.statusFlow

    private val _message = MutableStateFlow(Event.initial(""))
    val message = _message.asStateFlow()

    fun onNameChange(input: String) {
        _name.value = input
    }

    fun onEmailChange(input: String) {
        _email.value = input
    }

    fun onPasswordChange(input: String) {
        _password.value = input
    }

    fun onConfirmPasswordChange(input: String) {
        _confirmPassword.value = input
        _matchingPasswords.value = _password.value == _confirmPassword.value
    }

    fun onLoginPress() {
        runBlocking {
            launch {
                if (userRepo.login(LoginCredentials(_email.value, _password.value))) {
                    _message.value = Event("Logged in successfully")
                } else {
                    _message.value = Event("Error logging in")
                }
            }
        }
    }

    fun onRegisterPress() {
        if (_password.value == _confirmPassword.value) {
            runBlocking {
                launch {
                    if (
                        userRepo.register(
                            RegistrationCredentials(_name.value, _email.value, _password.value)
                        )
                    ) {
                        _message.value = Event("Registered successfully")
                    } else {
                        _message.value = Event("Error registering")
                    }
                }
            }
        } else {
            _message.value = Event("Passwords don't match")
        }
    }

    fun onRequestResetPress() {
        runBlocking {
            launch {
                if (userRepo.requestReset(ResetCredentials(_email.value))) {
                    _message.value = Event("Reset email sent\nJk... we're working on it")
                } else {
                    _message.value = Event("Error resetting password")
                }
            }
        }
    }

    fun onSwitchRegisterPress() {
        _registering.value = !_registering.value
    }

    fun onSwitchRequestResetPress() {
        _requestingReset.value = !_requestingReset.value
    }
}
