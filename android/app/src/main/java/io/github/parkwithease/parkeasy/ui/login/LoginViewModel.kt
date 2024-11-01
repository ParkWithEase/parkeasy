package io.github.parkwithease.parkeasy.ui.login

import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.setValue
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.assisted.Assisted
import dagger.assisted.AssistedFactory
import dagger.assisted.AssistedInject
import dagger.hilt.android.lifecycle.HiltViewModel
import io.github.parkwithease.parkeasy.data.local.AuthRepository
import io.github.parkwithease.parkeasy.data.remote.APIException
import io.github.parkwithease.parkeasy.data.remote.UserRepository
import io.github.parkwithease.parkeasy.model.ErrorDetail
import io.github.parkwithease.parkeasy.model.ErrorModel
import io.github.parkwithease.parkeasy.model.LoginCredentials
import io.github.parkwithease.parkeasy.model.RegistrationCredentials
import io.github.parkwithease.parkeasy.model.ResetCredentials
import java.io.IOException
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch

@HiltViewModel(assistedFactory = LoginViewModel.Factory::class)
// XXX: Lots of handlers and transforms, consider consolidation or splitting
@Suppress("detekt:TooManyFunctions")
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

    var formState by mutableStateOf(LoginFormState())
        private set

    private val _formEnabled = MutableStateFlow(true)
    val formEnabled = _formEnabled.asStateFlow()

    fun onLoginPress() {
        viewModelScope.launch {
            _formEnabled.value = false
            userRepo
                .login(LoginCredentials(formState.email.value, formState.password.value))
                .also { clearFieldErrors() }
                .onSuccess {
                    viewModelScope.launch { showSnackbar("Logged in successfully", null) }
                }
                .recoverRequestErrors("Login failed")
            _formEnabled.value = true
        }
    }

    fun onRegisterPress() {
        if (formState.password.value != formState.confirmPassword.value) return

        viewModelScope.launch {
            _formEnabled.value = false
            userRepo
                .register(
                    RegistrationCredentials(
                        formState.name.value,
                        formState.email.value,
                        formState.password.value,
                    )
                )
                .also { clearFieldErrors() }
                .onSuccess {
                    viewModelScope.launch { showSnackbar("Registered successfully", null) }
                }
                .recoverRequestErrors("Error registering")
            _formEnabled.value = true
        }
    }

    fun onRequestResetPress() {
        viewModelScope.launch {
            _formEnabled.value = false
            userRepo
                .requestReset(ResetCredentials(formState.email.value))
                .also { clearFieldErrors() }
                .onSuccess {
                    viewModelScope.launch {
                        showSnackbar("Reset email sent\nJk... we're working on it", null)
                    }
                }
                .recoverRequestErrors("Error resetting password")
            _formEnabled.value = true
        }
    }

    fun onNameChange(value: String) {
        formState = formState.run { copy(name = name.copy(value = value)) }
    }

    fun onEmailChange(value: String) {
        formState = formState.run { copy(email = email.copy(value = value)) }
    }

    fun onPasswordChange(value: String) {
        formState =
            formState.run {
                copy(
                    password = password.copy(value = value),
                    confirmPassword =
                        confirmPassword.copy(
                            error =
                                if (value == confirmPassword.value) null
                                else "Password does not match"
                        ),
                )
            }
    }

    fun onConfirmPasswordChange(value: String) {
        formState =
            formState.run {
                copy(
                    confirmPassword =
                        confirmPassword.copy(
                            value = value,
                            error = if (value == password.value) null else "Password does not match",
                        )
                )
            }
    }

    private fun Result<Unit>.recoverRequestErrors(operationFailMsg: String): Result<Unit> =
        recover {
            when (it) {
                is APIException -> {
                    errorToForm(it.error)
                    viewModelScope.launch { showSnackbar(operationFailMsg, null) }
                }
                is IOException -> {
                    viewModelScope.launch {
                        showSnackbar("Could not connect to server, are you online?", null)
                    }
                }
                else -> throw it
            }
        }

    private fun errorToForm(error: ErrorModel) {
        when (error.type) {
            ErrorModel.TYPE_INVALID_CREDENTIALS -> {
                formState =
                    formState.run {
                        copy(email = email.copy(error = ""), password = password.copy(error = ""))
                    }
            }

            ErrorModel.TYPE_PASSWORD_LENGTH ->
                formState =
                    formState.run {
                        copy(password = password.copy(error = "Password too long or too short"))
                    }

            else -> annotateErrorLocation(error.errors)
        }
    }

    private fun annotateErrorLocation(errors: List<ErrorDetail>) {
        for (err in errors) {
            when (err.location) {
                "body.email" ->
                    formState =
                        formState.run { copy(email = email.copy(error = "Invalid email address")) }

                "body.password" ->
                    formState =
                        formState.run { copy(password = password.copy(error = "Invalid password")) }
            }
        }
    }

    // Clear errors set via external services
    private fun clearFieldErrors() {
        formState =
            formState.run {
                copy(
                    name = name.copy(error = null),
                    email = email.copy(error = null),
                    password = password.copy(error = null),
                )
            }
    }
}

data class LoginFieldState(val value: String = "", val error: String? = null)

data class LoginFormState(
    val name: LoginFieldState = LoginFieldState(),
    val email: LoginFieldState = LoginFieldState(),
    val password: LoginFieldState = LoginFieldState(),
    val confirmPassword: LoginFieldState = LoginFieldState(),
)
