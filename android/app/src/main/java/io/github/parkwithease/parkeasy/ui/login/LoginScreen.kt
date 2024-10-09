package io.github.parkwithease.parkeasy.ui.login

import android.widget.Toast
import androidx.compose.animation.AnimatedVisibility
import androidx.compose.foundation.Image
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Person
import androidx.compose.material3.Button
import androidx.compose.material3.Icon
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.rememberUpdatedState
import androidx.compose.runtime.saveable.rememberSaveable
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.res.vectorResource
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.text.input.VisualTransformation
import androidx.compose.ui.unit.Dp
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import io.github.parkwithease.parkeasy.R
import io.github.parkwithease.parkeasy.model.LoginMode

private data class LoginUiState(
    val name: String,
    val email: String,
    val password: String,
    val confirmPassword: String,
    val loginMode: LoginMode,
)

private data class LoginUiEvents(
    val onNameChange: (String) -> Unit,
    val onEmailChange: (String) -> Unit,
    val onPasswordChange: (String) -> Unit,
    val onConfirmPasswordChange: (String) -> Unit,
    val onSwitchMode: (LoginMode) -> Unit,
)

private data class LoginEvents(
    val onLoginPress: (String, String) -> Unit,
    val onRegisterPress: (String, String, String, String) -> Unit,
    val onRequestResetPress: (String) -> Unit,
)

@Composable
fun LoginScreen(
    onLogin: () -> Unit,
    modifier: Modifier = Modifier,
    viewModel: LoginViewModel = hiltViewModel<LoginViewModel>(),
) {
    val context = LocalContext.current
    val message by viewModel.message.collectAsState()
    val loggedIn by viewModel.loggedIn.collectAsState(false)
    val events =
        LoginEvents(
            viewModel::onLoginPress,
            viewModel::onRegisterPress,
            viewModel::onRequestResetPress,
        )
    val latestOnLogin by rememberUpdatedState(onLogin)
    LaunchedEffect(loggedIn) {
        if (loggedIn) {
            latestOnLogin()
        }
    }
    LoginScreenInner(events, modifier)
    LaunchedEffect(message) {
        message.getContentIfNotHandled()?.let {
            Toast.makeText(context, it, Toast.LENGTH_SHORT).show()
        }
    }
}

@Composable
private fun LoginScreenInner(events: LoginEvents, modifier: Modifier = Modifier) {
    Surface(modifier) {
        Column {
            Row(
                verticalAlignment = Alignment.Bottom,
                horizontalArrangement = Arrangement.Center,
                modifier = Modifier.weight(1f).fillMaxSize(),
            ) {
                Image(
                    painter = painterResource(R.drawable.outlined_stacked),
                    contentDescription = stringResource(R.string.logo),
                    modifier = Modifier.size(280.dp),
                )
            }
            Row(modifier = Modifier.weight(1f).fillMaxSize()) { LoginForm(events, 280.dp) }
        }
    }
}

@Composable
private fun LoginForm(events: LoginEvents, width: Dp, modifier: Modifier = Modifier) {
    var name by rememberSaveable { mutableStateOf("") }
    var email by rememberSaveable { mutableStateOf("") }
    var password by rememberSaveable { mutableStateOf("") }
    var confirmPassword by rememberSaveable { mutableStateOf("") }
    var loginMode by rememberSaveable { mutableStateOf(LoginMode.LOGIN) }
    val uiState = LoginUiState(name, email, password, confirmPassword, loginMode)
    val uiEvents =
        LoginUiEvents(
            { name = it },
            { email = it },
            { password = it },
            { confirmPassword = it },
            { loginMode = it },
        )
    Column(
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.spacedBy(4.dp),
        modifier = modifier.fillMaxSize(),
    ) {
        LoginFields(uiState, uiEvents, width)
        Button(
            onClick = {
                when (loginMode) {
                    LoginMode.LOGIN -> events.onLoginPress(email, password)
                    LoginMode.REGISTER ->
                        events.onRegisterPress(name, email, password, confirmPassword)
                    LoginMode.FORGOT -> events.onRequestResetPress(email)
                }
            },
            modifier = Modifier.width(width),
        ) {
            Text(
                stringResource(
                    when (loginMode) {
                        LoginMode.LOGIN -> R.string.login
                        LoginMode.REGISTER -> R.string.register
                        LoginMode.FORGOT -> R.string.forgot_password
                    }
                )
            )
        }
        AnimatedVisibility(loginMode != LoginMode.FORGOT) {
            SwitchRegisterText(
                stringResource(
                    if (loginMode == LoginMode.REGISTER) R.string.already_registered
                    else R.string.not_registered
                ),
                stringResource(
                    if (loginMode == LoginMode.REGISTER) R.string.login_instead
                    else R.string.register_instead
                ),
                if (loginMode == LoginMode.REGISTER) LoginMode.LOGIN else LoginMode.REGISTER,
                { loginMode = it },
                Modifier.width(width),
            )
        }
    }
}

@Composable
private fun LoginFields(
    state: LoginUiState,
    events: LoginUiEvents,
    width: Dp,
    modifier: Modifier = Modifier,
) {
    Column(modifier) {
        AnimatedVisibility(state.loginMode == LoginMode.REGISTER) {
            LoginField(
                state.name,
                events.onNameChange,
                stringResource(R.string.name),
                Icons.Filled.Person,
                KeyboardOptions(keyboardType = KeyboardType.Text),
            )
        }
        LoginField(
            state.email,
            events.onEmailChange,
            stringResource(R.string.email),
            ImageVector.vectorResource(R.drawable.email),
            KeyboardOptions(keyboardType = KeyboardType.Email),
        )
        AnimatedVisibility(state.loginMode != LoginMode.FORGOT) {
            LoginField(
                state.password,
                events.onPasswordChange,
                stringResource(R.string.password),
                ImageVector.vectorResource(R.drawable.password),
                KeyboardOptions(keyboardType = KeyboardType.Password),
                visualTransformation = PasswordVisualTransformation(),
            )
        }
        AnimatedVisibility(state.loginMode == LoginMode.REGISTER) {
            LoginField(
                state.confirmPassword,
                events.onConfirmPasswordChange,
                stringResource(R.string.confirm_password),
                ImageVector.vectorResource(R.drawable.password),
                KeyboardOptions(keyboardType = KeyboardType.Password),
                isError = state.password != state.confirmPassword,
                visualTransformation = PasswordVisualTransformation(),
            )
        }
        AnimatedVisibility(state.loginMode != LoginMode.REGISTER) {
            SwitchRequestResetText(
                stringResource(
                    if (state.loginMode == LoginMode.FORGOT) R.string.return_login
                    else R.string.forgot_password
                ),
                if (state.loginMode == LoginMode.FORGOT) LoginMode.LOGIN else LoginMode.FORGOT,
                events.onSwitchMode,
                Modifier.width(width),
            )
        }
    }
}

@Composable
private fun LoginField(
    value: String,
    onValueChange: (String) -> Unit,
    label: String,
    imageVector: ImageVector,
    keyboardOptions: KeyboardOptions,
    modifier: Modifier = Modifier,
    isError: Boolean = false,
    visualTransformation: VisualTransformation = VisualTransformation.None,
) {
    OutlinedTextField(
        value = value,
        onValueChange = { onValueChange(it) },
        label = { Text(label) },
        leadingIcon = { Icon(imageVector = imageVector, contentDescription = label) },
        isError = isError,
        visualTransformation = visualTransformation,
        keyboardOptions = keyboardOptions,
        singleLine = true,
        modifier = modifier,
    )
}

@Composable
private fun SwitchRequestResetText(
    text: String,
    revertMode: LoginMode,
    onClick: (LoginMode) -> Unit,
    modifier: Modifier = Modifier,
) {
    Row(modifier, Arrangement.End) {
        Text(
            text,
            Modifier.clickable { onClick(revertMode) },
            color = MaterialTheme.colorScheme.primary,
        )
    }
}

@Composable
private fun SwitchRegisterText(
    supportingText: String,
    clickableText: String,
    revertMode: LoginMode,
    onClick: (LoginMode) -> Unit,
    modifier: Modifier = Modifier,
) {
    Row(modifier, Arrangement.Center) {
        Text(supportingText, color = MaterialTheme.colorScheme.onSurface)
        Text(" ")
        Text(
            clickableText,
            Modifier.clickable { onClick(revertMode) },
            color = MaterialTheme.colorScheme.primary,
        )
    }
}
