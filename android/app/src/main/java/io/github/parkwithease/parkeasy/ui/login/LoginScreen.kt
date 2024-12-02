package io.github.parkwithease.parkeasy.ui.login

import androidx.activity.compose.BackHandler
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
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.SnackbarHost
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
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import io.github.parkwithease.parkeasy.R
import io.github.parkwithease.parkeasy.model.LoginMode
import io.github.parkwithease.parkeasy.ui.common.ParkEasyTextField

@Composable
fun LoginScreen(
    onExitApp: () -> Unit,
    onNavigateFromLogin: () -> Unit,
    modifier: Modifier = Modifier,
    viewModel: LoginViewModel = hiltViewModel(),
) {
    val loggedIn by viewModel.loggedIn.collectAsState(false)
    val latestOnNavigateFromLogin by rememberUpdatedState(onNavigateFromLogin)

    val formEnabled by viewModel.formEnabled.collectAsState()

    val state = rememberLoginFormState(viewModel)
    val handler = rememberLoginFormHandler(viewModel)

    LaunchedEffect(loggedIn) {
        if (loggedIn) {
            latestOnNavigateFromLogin()
        }
    }

    BackHandler(enabled = true) { onExitApp() }

    Scaffold(
        snackbarHost = { SnackbarHost(hostState = viewModel.snackbarState) },
        modifier = modifier,
    ) { innerPadding ->
        LoginScreen(
            state = state,
            handler = handler,
            enabled = formEnabled,
            modifier = Modifier.padding(innerPadding),
        )
    }
}

@Composable
private fun LoginScreen(
    state: LoginFormState,
    handler: LoginFormHandler,
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
                    contentDescription = null,
                    modifier = Modifier.aspectRatio(1f),
                )
            }
            LoginForm(
                state = state,
                handler = handler,
                enabled = enabled,
                modifier = Modifier.fillMaxWidth(),
            )
            Spacer(modifier = Modifier.windowInsetsBottomHeight(WindowInsets.navigationBars))
        }
    }
}

@Composable
private fun LoginForm(
    state: LoginFormState,
    handler: LoginFormHandler,
    modifier: Modifier = Modifier,
    enabled: Boolean = true,
) {
    var loginMode by rememberSaveable { mutableStateOf(LoginMode.LOGIN) }
    Column(horizontalAlignment = Alignment.CenterHorizontally, modifier = modifier) {
        LoginFields(
            state,
            handler,
            loginMode,
            enabled = enabled,
            modifier = Modifier.widthIn(max = 320.dp),
        )
        LoginButtons(
            loginMode,
            handler,
            onSwitchModeClick = { loginMode = it },
            enabled = enabled,
            modifier = Modifier.widthIn(max = 320.dp),
        )
    }
}

@Composable
private fun LoginFields(
    state: LoginFormState,
    handler: LoginFormHandler,
    mode: LoginMode,
    modifier: Modifier = Modifier,
    enabled: Boolean = true,
) {
    Column(modifier = modifier) {
        AnimatedVisibility(mode == LoginMode.REGISTER) {
            ParkEasyTextField(
                state = state.name,
                onValueChange = handler.onNameChange,
                modifier = Modifier.fillMaxWidth(),
                leadingIconImage = Icons.Filled.Person,
                enabled = enabled,
                labelId = R.string.name,
                keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Text),
            )
        }
        ParkEasyTextField(
            state = state.email,
            onValueChange = handler.onEmailChange,
            modifier = Modifier.fillMaxWidth(),
            enabled = enabled,
            labelId = R.string.email,
            leadingIconImage = ImageVector.vectorResource(R.drawable.email),
            keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Email),
        )
        AnimatedVisibility(mode != LoginMode.FORGOT) {
            ParkEasyTextField(
                state = state.password,
                onValueChange = handler.onPasswordChange,
                modifier = Modifier.fillMaxWidth(),
                enabled = enabled,
                labelId = R.string.password,
                leadingIconImage = ImageVector.vectorResource(R.drawable.password),
                visualTransformation = PasswordVisualTransformation(),
                keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Password),
            )
        }
        AnimatedVisibility(mode == LoginMode.REGISTER) {
            ParkEasyTextField(
                state = state.confirmPassword,
                onValueChange = handler.onConfirmPasswordChange,
                modifier = Modifier.fillMaxWidth(),
                enabled = enabled,
                labelId = R.string.confirm_password,
                leadingIconImage = ImageVector.vectorResource(R.drawable.password),
                visualTransformation = PasswordVisualTransformation(),
                keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Password),
            )
        }
    }
}

@Composable
private fun LoginButtons(
    mode: LoginMode,
    handler: LoginFormHandler,
    onSwitchModeClick: (LoginMode) -> Unit,
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
                onSwitchModeClick,
            )
        }
        Button(
            onClick =
                when (mode) {
                    LoginMode.LOGIN -> handler.onLoginClick
                    LoginMode.REGISTER -> handler.onRegisterClick
                    LoginMode.FORGOT -> handler.onRequestResetClick
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
                onSwitchModeClick,
                modifier = Modifier.fillMaxWidth(),
            )
        }
    }
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

// @Composable
// @Preview
// private fun PreviewLoginScreen() {
//    LoginScreen(LoginFormState(), {}, {}, {}, {}, {}, {}, {})
// }
//
// @Composable
// @Preview
// private fun PreviewLoginScreenError() {
//    LoginScreen(
//        LoginFormState(
//            email = FieldState(value = "not-an-email", error = "Invalid email"),
//            confirmPassword = FieldState("", error = "Password does not match"),
//        ),
//        {},
//        {},
//        {},
//        {},
//        {},
//        {},
//        {},
//    )
// }
