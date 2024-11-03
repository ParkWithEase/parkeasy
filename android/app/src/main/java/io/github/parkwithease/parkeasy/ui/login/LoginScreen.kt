package io.github.parkwithease.parkeasy.ui.login

import androidx.compose.animation.AnimatedVisibility
import androidx.compose.foundation.Image
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.WindowInsets
import androidx.compose.foundation.layout.aspectRatio
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.imePadding
import androidx.compose.foundation.layout.navigationBars
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.widthIn
import androidx.compose.foundation.layout.windowInsetsBottomHeight
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.foundation.verticalScroll
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
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.res.vectorResource
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.text.input.VisualTransformation
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import io.github.parkwithease.parkeasy.R
import io.github.parkwithease.parkeasy.model.LoginMode

@Composable
fun LoginScreen(onLogin: () -> Unit, viewModel: LoginViewModel, modifier: Modifier = Modifier) {
    val loggedIn by viewModel.loggedIn.collectAsState(false)
    val formEnabled by viewModel.formEnabled.collectAsState()
    val latestOnLogin by rememberUpdatedState(onLogin)
    LaunchedEffect(loggedIn) {
        if (loggedIn) {
            latestOnLogin()
        }
    }

    LoginScreenInner(
        viewModel.formState,
        viewModel::onNameChange,
        viewModel::onEmailChange,
        viewModel::onPasswordChange,
        viewModel::onConfirmPasswordChange,
        viewModel::onLoginPress,
        viewModel::onRegisterPress,
        viewModel::onRequestResetPress,
        enabled = formEnabled,
        modifier = modifier,
    )
}

@Composable
private fun LoginScreenInner(
    formState: LoginFormState,
    onNameChange: (String) -> Unit,
    onEmailChange: (String) -> Unit,
    onPasswordChange: (String) -> Unit,
    onConfirmPasswordChange: (String) -> Unit,
    onLogin: () -> Unit,
    onRegister: () -> Unit,
    onRequestReset: () -> Unit,
    modifier: Modifier = Modifier,
    enabled: Boolean = true,
) {
    Surface(modifier = modifier) {
        Column(
            horizontalAlignment = Alignment.CenterHorizontally,
            verticalArrangement = Arrangement.Center,
            modifier =
                Modifier.padding(horizontal = 16.dp)
                    .fillMaxSize()
                    .imePadding()
                    .verticalScroll(rememberScrollState(), reverseScrolling = true),
        ) {
            Box(modifier = Modifier.widthIn(min = 80.dp, max = 240.dp)) {
                Image(
                    painter = painterResource(R.drawable.logo_stacked_outlined),
                    contentDescription = stringResource(R.string.logo),
                    modifier = Modifier.aspectRatio(1f),
                )
            }
            LoginForm(
                formState = formState,
                onNameChange = onNameChange,
                onEmailChange = onEmailChange,
                onPasswordChange = onPasswordChange,
                onConfirmPasswordChange = onConfirmPasswordChange,
                onLogin = onLogin,
                onRegister = onRegister,
                onRequestReset = onRequestReset,
                enabled = enabled,
                modifier = Modifier.fillMaxWidth(),
            )
            Spacer(modifier = Modifier.windowInsetsBottomHeight(WindowInsets.navigationBars))
        }
    }
}

@Composable
private fun LoginForm(
    formState: LoginFormState,
    onNameChange: (String) -> Unit,
    onEmailChange: (String) -> Unit,
    onPasswordChange: (String) -> Unit,
    onConfirmPasswordChange: (String) -> Unit,
    onLogin: () -> Unit,
    onRegister: () -> Unit,
    onRequestReset: () -> Unit,
    modifier: Modifier = Modifier,
    enabled: Boolean = true,
) {
    var loginMode by rememberSaveable { mutableStateOf(LoginMode.LOGIN) }
    Column(horizontalAlignment = Alignment.CenterHorizontally, modifier = modifier) {
        LoginFields(
            formState,
            loginMode,
            onNameChange = onNameChange,
            onEmailChange = onEmailChange,
            onPasswordChange = onPasswordChange,
            onConfirmPasswordChange = onConfirmPasswordChange,
            enabled = enabled,
            modifier = Modifier.widthIn(max = 320.dp),
        )
        LoginButtons(
            loginMode,
            onLogin = onLogin,
            onRegister = onRegister,
            onRequestReset = onRequestReset,
            onSwitchMode = { loginMode = it },
            enabled = enabled,
            modifier = Modifier.widthIn(max = 320.dp),
        )
    }
}

@Composable
private fun LoginFields(
    state: LoginFormState,
    mode: LoginMode,
    onNameChange: (String) -> Unit,
    onEmailChange: (String) -> Unit,
    onPasswordChange: (String) -> Unit,
    onConfirmPasswordChange: (String) -> Unit,
    modifier: Modifier = Modifier,
    enabled: Boolean = true,
) {
    Column(modifier = modifier) {
        AnimatedVisibility(mode == LoginMode.REGISTER) {
            LoginField(
                state.name.value,
                onNameChange,
                stringResource(R.string.name),
                Icons.Filled.Person,
                KeyboardOptions(keyboardType = KeyboardType.Text),
                isError = state.name.error != null,
                supportingText = { state.name.error?.also { Text(it) } },
                enabled = enabled,
                modifier = Modifier.fillMaxWidth(),
            )
        }
        LoginField(
            state.email.value,
            onEmailChange,
            stringResource(R.string.email),
            ImageVector.vectorResource(R.drawable.email),
            KeyboardOptions(keyboardType = KeyboardType.Email),
            isError = state.email.error != null,
            supportingText = { state.email.error?.also { Text(it) } },
            enabled = enabled,
            modifier = Modifier.fillMaxWidth(),
        )
        AnimatedVisibility(mode != LoginMode.FORGOT) {
            LoginField(
                state.password.value,
                onPasswordChange,
                stringResource(R.string.password),
                ImageVector.vectorResource(R.drawable.password),
                KeyboardOptions(keyboardType = KeyboardType.Password),
                visualTransformation = PasswordVisualTransformation(),
                isError = state.password.error != null,
                supportingText = { state.password.error?.also { Text(it) } },
                enabled = enabled,
                modifier = Modifier.fillMaxWidth(),
            )
        }
        AnimatedVisibility(mode == LoginMode.REGISTER) {
            LoginField(
                state.confirmPassword.value,
                onConfirmPasswordChange,
                stringResource(R.string.confirm_password),
                ImageVector.vectorResource(R.drawable.password),
                KeyboardOptions(keyboardType = KeyboardType.Password),
                isError = state.confirmPassword.error != null,
                supportingText = { state.confirmPassword.error?.also { Text(it) } },
                visualTransformation = PasswordVisualTransformation(),
                enabled = enabled,
                modifier = Modifier.fillMaxWidth(),
            )
        }
    }
}

@Composable
private fun LoginButtons(
    mode: LoginMode,
    onLogin: () -> Unit,
    onRegister: () -> Unit,
    onRequestReset: () -> Unit,
    onSwitchMode: (LoginMode) -> Unit,
    modifier: Modifier = Modifier,
    enabled: Boolean = true,
) {
    Column(verticalArrangement = Arrangement.spacedBy(2.dp), modifier = modifier) {
        AnimatedVisibility(mode != LoginMode.REGISTER, modifier = Modifier.align(Alignment.End)) {
            SwitchRequestResetText(
                stringResource(
                    if (mode == LoginMode.FORGOT) R.string.return_login
                    else R.string.forgot_password
                ),
                if (mode == LoginMode.FORGOT) LoginMode.LOGIN else LoginMode.FORGOT,
                onSwitchMode,
            )
        }
        Button(
            onClick =
                when (mode) {
                    LoginMode.LOGIN -> onLogin
                    LoginMode.REGISTER -> onRegister
                    LoginMode.FORGOT -> onRequestReset
                },
            enabled = enabled,
            modifier = Modifier.fillMaxWidth(),
        ) {
            Text(
                stringResource(
                    when (mode) {
                        LoginMode.LOGIN -> R.string.login
                        LoginMode.REGISTER -> R.string.register
                        LoginMode.FORGOT -> R.string.forgot_password
                    }
                )
            )
        }
        AnimatedVisibility(mode != LoginMode.FORGOT) {
            SwitchRegisterText(
                stringResource(
                    if (mode == LoginMode.REGISTER) R.string.already_registered
                    else R.string.not_registered
                ),
                stringResource(
                    if (mode == LoginMode.REGISTER) R.string.login_instead
                    else R.string.register_instead
                ),
                if (mode == LoginMode.REGISTER) LoginMode.LOGIN else LoginMode.REGISTER,
                onSwitchMode,
                modifier = Modifier.fillMaxWidth(),
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
    supportingText: @Composable (() -> Unit)? = null,
    visualTransformation: VisualTransformation = VisualTransformation.None,
    enabled: Boolean = true,
) {
    OutlinedTextField(
        value = value,
        onValueChange = { onValueChange(it) },
        label = { Text(label) },
        leadingIcon = { Icon(imageVector = imageVector, contentDescription = label) },
        isError = isError,
        visualTransformation = visualTransformation,
        keyboardOptions = keyboardOptions,
        enabled = enabled,
        singleLine = true,
        modifier = modifier,
        supportingText = supportingText,
    )
}

@Composable
private fun SwitchRequestResetText(
    text: String,
    revertMode: LoginMode,
    onClick: (LoginMode) -> Unit,
    modifier: Modifier = Modifier,
) {
    Text(
        text,
        modifier = modifier.clickable { onClick(revertMode) }.padding(vertical = 16.dp),
        color = MaterialTheme.colorScheme.primary,
    )
}

@Composable
private fun SwitchRegisterText(
    supportingText: String,
    clickableText: String,
    revertMode: LoginMode,
    onClick: (LoginMode) -> Unit,
    modifier: Modifier = Modifier,
) {
    Row(
        modifier = modifier,
        horizontalArrangement = Arrangement.spacedBy(2.dp, Alignment.CenterHorizontally),
        verticalAlignment = Alignment.CenterVertically,
    ) {
        Text(supportingText, color = MaterialTheme.colorScheme.onSurface)
        Text(
            clickableText,
            Modifier.clickable { onClick(revertMode) }.padding(vertical = 16.dp),
            color = MaterialTheme.colorScheme.primary,
        )
    }
}

@Composable
@Preview
private fun PreviewLoginScreen() {
    LoginScreenInner(LoginFormState(), {}, {}, {}, {}, {}, {}, {})
}

@Composable
@Preview
private fun PreviewLoginScreenError() {
    LoginScreenInner(
        LoginFormState(
            email = LoginFieldState(value = "not-an-email", error = "Invalid email"),
            confirmPassword = LoginFieldState(error = "Password does not match"),
        ),
        {},
        {},
        {},
        {},
        {},
        {},
        {},
    )
}
