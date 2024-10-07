package io.github.parkwithease.parkeasy.ui.login

import androidx.lifecycle.ViewModel
import dagger.hilt.android.lifecycle.HiltViewModel
import io.github.parkwithease.parkeasy.data.local.AuthRepository
import io.github.parkwithease.parkeasy.data.remote.UserRepository
import io.github.parkwithease.parkeasy.model.LoginCredentials
import io.github.parkwithease.parkeasy.model.RegistrationCredentials
import javax.inject.Inject
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import kotlinx.coroutines.runBlocking

@HiltViewModel
class LoginViewModel
@Inject
constructor(private val authRepo: AuthRepository, private val userRepo: UserRepository) :
    ViewModel() {
    private val _name = MutableStateFlow("")
    val name = _name.asStateFlow()

    private val _email = MutableStateFlow("")
    val email = _email.asStateFlow()

    private val _password = MutableStateFlow("")
    val password = _password.asStateFlow()

    private val _confirmPassword = MutableStateFlow("")
    val confirmPassword = _confirmPassword.asStateFlow()

    private val _registering = MutableStateFlow(false)
    val registering = _registering.asStateFlow()

    private val _loggedIn = MutableStateFlow(false)
    val loggedIn = _loggedIn.asStateFlow()

    init {
        runBlocking { launch { _loggedIn.value = authRepo.getStatus() } }
    }

    private suspend fun login(credentials: LoginCredentials) {
        _loggedIn.value = userRepo.login(credentials)
    }

    private suspend fun register(credentials: RegistrationCredentials) {
        _loggedIn.value = userRepo.register(credentials)
    }

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
    }

    fun onLoginPress() {
        runBlocking { launch { login(LoginCredentials(_email.value, _password.value)) } }
    }

    fun onRegisterPress() {
        runBlocking {
            launch { register(RegistrationCredentials(_name.value, _email.value, _password.value)) }
        }
    }

    fun onSwitchPress() {
        _registering.value = !_registering.value
    }
}
