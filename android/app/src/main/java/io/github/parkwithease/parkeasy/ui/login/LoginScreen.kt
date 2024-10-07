package io.github.parkwithease.parkeasy.ui.login

import androidx.compose.foundation.Image
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.material3.Button
import androidx.compose.material3.Icon
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.res.vectorResource
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import io.github.parkwithease.parkeasy.R

private data class LoginState(
    val name: String,
    val email: String,
    val password: String,
    val confirmPassword: String,
    val registering: Boolean,
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
}

// @Preview
// @Composable
// private fun PreviewLoginScreenInner() {
//    ParkEasyTheme { LoginScreenInner("", "", {}, {}, {}) }
// }

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
            Row(modifier = Modifier.weight(1f).fillMaxSize()) { LoginForm(state, events) }
        }
    }
}

@Composable
private fun LoginForm(state: LoginState, events: LoginEvents, modifier: Modifier = Modifier) {
    Column(
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.spacedBy(4.dp),
        modifier = modifier.fillMaxSize(),
    ) {
        EmailField(state.email, events.onEmailChange)
        PasswordField(state.password, events.onPasswordChange)
        Row(modifier = Modifier.width(280.dp), horizontalArrangement = Arrangement.spacedBy(4.dp)) {
            RegisterButton(modifier = Modifier.weight(1f), events.onRegisterPress)
            LoginButton(modifier = Modifier.weight(1f), events.onLoginPress)
        }
    }
}

@Composable
private fun EmailField(
    text: String,
    onValueChange: (String) -> Unit,
    modifier: Modifier = Modifier,
) {
    OutlinedTextField(
        value = text,
        onValueChange = { onValueChange(it) },
        label = { Text(stringResource(R.string.email)) },
        leadingIcon = {
            Icon(
                imageVector = ImageVector.vectorResource(R.drawable.email),
                contentDescription = stringResource(R.string.email_icon),
            )
        },
        keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Email),
        singleLine = true,
        modifier = modifier,
    )
}

@Composable
private fun PasswordField(
    text: String,
    onValueChange: (String) -> Unit,
    modifier: Modifier = Modifier,
) {
    OutlinedTextField(
        value = text,
        onValueChange = { onValueChange(it) },
        label = { Text(stringResource(R.string.password)) },
        leadingIcon = {
            Icon(
                imageVector = ImageVector.vectorResource(R.drawable.password),
                contentDescription = stringResource(R.string.password_icon),
            )
        },
        visualTransformation = PasswordVisualTransformation(),
        keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Password),
        singleLine = true,
        modifier = modifier,
    )
}

@Composable
private fun LoginButton(modifier: Modifier = Modifier, onClick: () -> Unit) {
    Button(onClick = { onClick() }, modifier = modifier) { Text(stringResource(R.string.login)) }
}

@Composable
private fun RegisterButton(modifier: Modifier = Modifier, onClick: () -> Unit) {
    Button(onClick = { onClick() }, modifier = modifier) { Text(stringResource(R.string.register)) }
}
