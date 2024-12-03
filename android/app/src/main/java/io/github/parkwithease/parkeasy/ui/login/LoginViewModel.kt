package io.github.parkwithease.parkeasy.ui.login

import androidx.compose.material3.SnackbarHostState
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import io.github.parkwithease.parkeasy.data.local.AuthRepository
import io.github.parkwithease.parkeasy.data.remote.APIException
import io.github.parkwithease.parkeasy.data.remote.UserRepository
import io.github.parkwithease.parkeasy.model.ErrorDetail
import io.github.parkwithease.parkeasy.model.ErrorModel
import io.github.parkwithease.parkeasy.model.FieldState
import io.github.parkwithease.parkeasy.model.LoginCredentials
import io.github.parkwithease.parkeasy.model.RegistrationCredentials
import io.github.parkwithease.parkeasy.model.ResetCredentials
import java.io.IOException
import javax.inject.Inject
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch

@Suppress("detekt:TooManyFunctions")
@HiltViewModel
class LoginViewModel
@Inject
constructor(authRepo: AuthRepository, private val userRepo: UserRepository) : ViewModel() {
    val loggedIn = authRepo.statusFlow
    val snackbarState = SnackbarHostState()

    private val _formEnabled = MutableStateFlow(true)
    val formEnabled = _formEnabled.asStateFlow()

    var state by mutableStateOf(LoginFormState())
        private set

    fun onLoginClick() {
        viewModelScope.launch {
            _formEnabled.value = false
            userRepo
                .login(LoginCredentials(state.email.value, state.password.value))
                .also { clearFieldErrors() }
                .onSuccess {
                    viewModelScope.launch { snackbarState.showSnackbar("Logged in successfully") }
                }
                .recoverRequestErrors("Login failed")
            _formEnabled.value = true
        }
    }

    fun onRegisterClick() {
        if (state.password.value != state.confirmPassword.value) return

        viewModelScope.launch {
            _formEnabled.value = false
            userRepo
                .register(
                    RegistrationCredentials(
                        state.name.value,
                        state.email.value,
                        state.password.value,
                    )
                )
                .also { clearFieldErrors() }
                .onSuccess {
                    viewModelScope.launch { snackbarState.showSnackbar("Registered successfully") }
                }
                .recoverRequestErrors("Error registering")
            _formEnabled.value = true
        }
    }

    fun onRequestResetClick() {
        viewModelScope.launch {
            _formEnabled.value = false
            userRepo
                .requestReset(ResetCredentials(state.email.value))
                .also { clearFieldErrors() }
                .onSuccess {
                    viewModelScope.launch {
                        snackbarState.showSnackbar("Reset email sent\nJk... we're working on it")
                    }
                }
                .recoverRequestErrors("Error resetting password")
            _formEnabled.value = true
        }
    }

    fun onNameChange(value: String) {
        state =
            state.run {
                copy(
                    name =
                        name.copy(
                            value = value,
                            error = if (value != "") null else "Name cannot be empty",
                        )
                )
            }
    }

    fun onEmailChange(value: String) {
        state =
            state.run {
                copy(
                    email =
                        email.copy(
                            value = value,
                            error = if (value != "") null else "Email cannot be empty",
                        )
                )
            }
    }

    fun onPasswordChange(value: String) {
        state =
            state.run {
                copy(
                    password =
                        password.copy(
                            value = value,
                            error = if (value != "") null else "Password cannot be empty",
                        ),
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
        state =
            state.run {
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
                    viewModelScope.launch { snackbarState.showSnackbar(operationFailMsg) }
                }
                is IOException -> {
                    viewModelScope.launch {
                        snackbarState.showSnackbar("Could not connect to server, are you online?")
                    }
                }
                else -> throw it
            }
        }

    private fun errorToForm(error: ErrorModel) {
        when (error.type) {
            ErrorModel.TYPE_INVALID_CREDENTIALS -> {
                state =
                    state.run {
                        copy(email = email.copy(error = ""), password = password.copy(error = ""))
                    }
            }

            ErrorModel.TYPE_PASSWORD_LENGTH ->
                state =
                    state.run {
                        copy(password = password.copy(error = "Password too long or too short"))
                    }

            else -> annotateErrorLocation(error.errors)
        }
    }

    private fun annotateErrorLocation(errors: List<ErrorDetail>) {
        for (err in errors) {
            when (err.location) {
                "body.email" ->
                    state = state.run { copy(email = email.copy(error = "Invalid email address")) }

                "body.password" ->
                    state = state.run { copy(password = password.copy(error = "Invalid password")) }
            }
        }
    }

    // Clear errors set via external services
    private fun clearFieldErrors() {
        state =
            state.run {
                copy(
                    name = name.copy(error = null),
                    email = email.copy(error = null),
                    password = password.copy(error = null),
                )
            }
    }

    fun createHandler() =
        LoginFormHandler(
            onNameChange = this::onNameChange,
            onEmailChange = this::onEmailChange,
            onPasswordChange = this::onPasswordChange,
            onConfirmPasswordChange = this::onConfirmPasswordChange,
            onLoginClick = this::onLoginClick,
            onRegisterClick = this::onRegisterClick,
            onRequestResetClick = this::onRequestResetClick,
        )
}

@Composable
fun rememberLoginFormHandler(viewModel: LoginViewModel) =
    remember(viewModel) { viewModel.createHandler() }

data class LoginFormState(
    val name: FieldState<String> = FieldState(""),
    val email: FieldState<String> = FieldState(""),
    val password: FieldState<String> = FieldState(""),
    val confirmPassword: FieldState<String> = FieldState(""),
)

data class LoginFormHandler(
    val onNameChange: (String) -> Unit = {},
    val onEmailChange: (String) -> Unit = {},
    val onPasswordChange: (String) -> Unit = {},
    val onConfirmPasswordChange: (String) -> Unit = {},
    val onLoginClick: () -> Unit = {},
    val onRegisterClick: () -> Unit = {},
    val onRequestResetClick: () -> Unit = {},
)
