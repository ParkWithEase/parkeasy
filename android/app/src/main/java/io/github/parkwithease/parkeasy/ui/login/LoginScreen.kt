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
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.res.vectorResource
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import io.github.parkwithease.parkeasy.R
import io.github.parkwithease.parkeasy.ui.theme.ParkEasyTheme

@Composable
fun LoginScreen(
    modifier: Modifier = Modifier,
    viewModel: LoginViewModel = hiltViewModel<LoginViewModel>(),
) {
    val email by viewModel.email.collectAsState()
    val password by viewModel.password.collectAsState()
    val onEmailChange = viewModel::onEmailChange
    val onPasswordChange = viewModel::onPasswordChange
    val onLoginPress = viewModel::onLoginPress
    LoginScreenInner(email, password, onEmailChange, onPasswordChange, onLoginPress, modifier)
}

@Preview
@Composable
private fun PreviewLoginScreenInner() {
    ParkEasyTheme { LoginScreenInner("", "", {}, {}, {}) }
}

@Composable
fun LoginScreenInner(
    email: String,
    password: String,
    onEmailChange: (String) -> Unit,
    onPasswordChange: (String) -> Unit,
    onLoginPress: () -> Unit,
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
            Row(modifier = Modifier.weight(1f).fillMaxSize()) {
                LoginForm(email, password, onEmailChange, onPasswordChange, onLoginPress)
            }
        }
    }
}

@Composable
fun LoginForm(
    email: String,
    password: String,
    onEmailChange: (String) -> Unit,
    onPasswordChange: (String) -> Unit,
    onLoginPress: () -> Unit,
    modifier: Modifier = Modifier,
) {
    Column(
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.spacedBy(4.dp),
        modifier = modifier.fillMaxSize(),
    ) {
        EmailField(email, onEmailChange)
        PasswordField(password, onPasswordChange)
        Row(modifier = Modifier.width(280.dp), horizontalArrangement = Arrangement.spacedBy(4.dp)) {
            RegisterButton(modifier = Modifier.weight(1f)) {}
            LoginButton(modifier = Modifier.weight(1f), onLoginPress)
        }
    }
}

@Composable
fun EmailField(text: String, onValueChange: (String) -> Unit, modifier: Modifier = Modifier) {
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
fun PasswordField(text: String, onValueChange: (String) -> Unit, modifier: Modifier = Modifier) {
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
fun LoginButton(modifier: Modifier = Modifier, onClick: () -> Unit) {
    Button(onClick = { onClick() }, modifier = modifier) { Text(stringResource(R.string.login)) }
}

@Composable
internal fun RegisterButton(modifier: Modifier = Modifier, onClick: () -> Unit) {
    Button(onClick = { onClick() }, modifier = modifier, enabled = false) {
        Text(stringResource(R.string.register))
    }
}
