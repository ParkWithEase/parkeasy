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
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.Dp
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import io.github.parkwithease.parkeasy.R
import io.github.parkwithease.parkeasy.ui.theme.ParkEasyTheme

private data class LoginState(
    val name: String,
    val email: String,
    val password: String,
    val confirmPassword: String,
    val matchingPasswords: Boolean,
    val loginMode: LoginMode,
    val loggedIn: Boolean,
)

private data class LoginEvents(
    val onNameChange: (String) -> Unit,
    val onEmailChange: (String) -> Unit,
    val onPasswordChange: (String) -> Unit,
    val onConfirmPasswordChange: (String) -> Unit,
    val onLoginPress: () -> Unit,
    val onRegisterPress: () -> Unit,
    val onRequestResetPress: () -> Unit,
    val onSwitchModePress: (LoginMode) -> Unit,
)

@Composable
fun LoginScreen(
    onLogin: () -> Unit,
    modifier: Modifier = Modifier,
    viewModel: LoginViewModel = hiltViewModel<LoginViewModel>(),
) {
    val context = LocalContext.current
    val message by viewModel.message.collectAsState()
    val state =
        LoginState(
            viewModel.name.collectAsState().value,
            viewModel.email.collectAsState().value,
            viewModel.password.collectAsState().value,
            viewModel.confirmPassword.collectAsState().value,
            viewModel.matchingPasswords.collectAsState().value,
            viewModel.loginMode.collectAsState().value,
            viewModel.loggedIn.collectAsState(false).value,
        )
    val events =
        LoginEvents(
            viewModel::onNameChange,
            viewModel::onEmailChange,
            viewModel::onPasswordChange,
            viewModel::onConfirmPasswordChange,
            viewModel::onLoginPress,
            viewModel::onRegisterPress,
            viewModel::onRequestResetPress,
            viewModel::onSwitchModePress,
        )
    val onLoginEvent: () -> Unit = onLogin
    LaunchedEffect(state.loggedIn) {
        if (state.loggedIn) {
            onLoginEvent()
        }
    }
    LoginScreenInner(state, events, modifier)
    LaunchedEffect(message) {
        message.getContentIfNotHandled()?.let {
            Toast.makeText(context, it, Toast.LENGTH_SHORT).show()
        }
    }
}

@Preview
@Composable
private fun PreviewLoginScreenInner() {
    val state =
        LoginState(
            "",
            "",
            "",
            "",
            matchingPasswords = false,
            loginMode = LoginMode.REGISTER,
            loggedIn = false,
        )
    val events = LoginEvents({}, {}, {}, {}, {}, {}, {}, {})
    ParkEasyTheme { LoginScreenInner(state, events) }
}

@Composable
private fun LoginScreenInner(
    state: LoginState,
    events: LoginEvents,
    modifier: Modifier = Modifier,
) {
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
            Row(modifier = Modifier.weight(1f).fillMaxSize()) { LoginForm(state, events, 280.dp) }
        }
    }
}

@Composable
private fun LoginForm(
    state: LoginState,
    events: LoginEvents,
    width: Dp,
    modifier: Modifier = Modifier,
) {
    Column(
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.spacedBy(4.dp),
        modifier = modifier.fillMaxSize(),
    ) {
        LoginFields(state, events, Modifier.width(width))
        AnimatedVisibility(state.loginMode != LoginMode.REGISTER) {
            SwitchRequestResetText(
                stringResource(
                    if (state.loginMode == LoginMode.FORGOT) R.string.return_login
                    else R.string.forgot_password
                ),
                if (state.loginMode == LoginMode.FORGOT) LoginMode.LOGIN else LoginMode.FORGOT,
                events.onSwitchModePress,
                Modifier.width(width),
            )
        }
        LoginButton(
            stringResource(
                when (state.loginMode) {
                    LoginMode.LOGIN -> R.string.login
                    LoginMode.REGISTER -> R.string.register
                    LoginMode.FORGOT -> R.string.forgot_password
                }
            ),
            when (state.loginMode) {
                LoginMode.LOGIN -> events.onLoginPress
                LoginMode.REGISTER -> events.onRegisterPress
                LoginMode.FORGOT -> events.onRequestResetPress
            },
            Modifier.width(width),
        )
        AnimatedVisibility(state.loginMode != LoginMode.FORGOT) {
            SwitchRegisterText(
                stringResource(
                    if (state.loginMode == LoginMode.REGISTER) R.string.already_registered
                    else R.string.not_registered
                ),
                stringResource(
                    if (state.loginMode == LoginMode.REGISTER) R.string.login_instead
                    else R.string.register_instead
                ),
                if (state.loginMode == LoginMode.REGISTER) LoginMode.LOGIN else LoginMode.REGISTER,
                events.onSwitchModePress,
                Modifier.width(width),
            )
        }
    }
}

@Composable
private fun LoginFields(state: LoginState, events: LoginEvents, modifier: Modifier = Modifier) {
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
                isError = !state.matchingPasswords,
                visualTransformation = PasswordVisualTransformation(),
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
private fun LoginButton(label: String, onClick: () -> Unit, modifier: Modifier = Modifier) {
    Button(onClick = { onClick() }, modifier = modifier) { Text(label) }
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
