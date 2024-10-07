package io.github.parkwithease.parkeasy.ui.login

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
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.vector.ImageVector
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
    val registering: Boolean,
    val loggedIn: Boolean,
)

private data class LoginEvents(
    val onNameChange: (String) -> Unit,
    val onEmailChange: (String) -> Unit,
    val onPasswordChange: (String) -> Unit,
    val onConfirmPasswordChange: (String) -> Unit,
    val onLoginPress: () -> Unit,
    val onRegisterPress: () -> Unit,
    val onSwitchPress: () -> Unit,
)

@Composable
fun LoginScreen(
    onLogin: () -> Unit,
    modifier: Modifier = Modifier,
    viewModel: LoginViewModel = hiltViewModel<LoginViewModel>(),
) {
    val state =
        LoginState(
            viewModel.name.collectAsState().value,
            viewModel.email.collectAsState().value,
            viewModel.password.collectAsState().value,
            viewModel.confirmPassword.collectAsState().value,
            viewModel.registering.collectAsState().value,
            viewModel.loggedIn.collectAsState().value,
        )
    val events =
        LoginEvents(
            viewModel::onNameChange,
            viewModel::onEmailChange,
            viewModel::onPasswordChange,
            viewModel::onConfirmPasswordChange,
            viewModel::onLoginPress,
            viewModel::onRegisterPress,
            viewModel::onSwitchPress,
        )
    LoginScreenInner(state, events, modifier)
    val onLoginEvent: () -> Unit = onLogin
    LaunchedEffect(state.loggedIn) {
        if (state.loggedIn) {
            onLoginEvent()
        }
    }
}

@Preview
@Composable
private fun PreviewLoginScreenInner() {
    val state = LoginState("", "", "", "", registering = false, loggedIn = false)
    val events = LoginEvents({}, {}, {}, {}, {}, {}, {})
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
        AnimatedVisibility(state.registering) {
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
        LoginField(
            state.password,
            events.onPasswordChange,
            stringResource(R.string.password),
            ImageVector.vectorResource(R.drawable.password),
            KeyboardOptions(keyboardType = KeyboardType.Password),
            visualTransformation = PasswordVisualTransformation(),
        )
        AnimatedVisibility(state.registering) {
            LoginField(
                state.confirmPassword,
                events.onConfirmPasswordChange,
                stringResource(R.string.confirm_password),
                ImageVector.vectorResource(R.drawable.password),
                KeyboardOptions(keyboardType = KeyboardType.Password),
                visualTransformation = PasswordVisualTransformation(),
            )
        }
        LoginButton(
            stringResource(if (state.registering) R.string.register else R.string.login),
            if (state.registering) events.onRegisterPress else events.onLoginPress,
            Modifier.width(width),
        )
        SwitchText(
            stringResource(
                if (state.registering) R.string.already_registered else R.string.not_registered
            ),
            stringResource(
                if (state.registering) R.string.login_instead else R.string.register_instead
            ),
            events.onSwitchPress,
        )
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
    visualTransformation: VisualTransformation = VisualTransformation.None,
) {
    OutlinedTextField(
        value = value,
        onValueChange = { onValueChange(it) },
        label = { Text(label) },
        leadingIcon = { Icon(imageVector = imageVector, contentDescription = label) },
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
private fun SwitchText(
    text: String,
    clickableText: String,
    onClick: () -> Unit,
    modifier: Modifier = Modifier,
) {
    Row(modifier) {
        Text(text, color = MaterialTheme.colorScheme.onSurface)
        Text(" ")
        Text(
            clickableText,
            Modifier.clickable { onClick() },
            color = MaterialTheme.colorScheme.primary,
        )
    }
}
