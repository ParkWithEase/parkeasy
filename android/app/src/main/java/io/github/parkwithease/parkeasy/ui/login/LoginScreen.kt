package io.github.parkwithease.parkeasy.ui.login

import androidx.activity.compose.BackHandler
import androidx.compose.animation.AnimatedVisibility
import androidx.compose.foundation.Image
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
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
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import io.github.parkwithease.parkeasy.R
import io.github.parkwithease.parkeasy.model.LoginMode
import io.github.parkwithease.parkeasy.ui.common.ClickableText
import io.github.parkwithease.parkeasy.ui.common.ParkEasyTextField
import io.github.parkwithease.parkeasy.ui.theme.ParkEasyTheme

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
    val handler = rememberLoginFormHandler(viewModel)
    var mode by rememberSaveable { mutableStateOf(LoginMode.LOGIN) }

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
            state = viewModel.state,
            handler = handler,
            mode = mode,
            onSwitchReset = {
                mode = if (mode == LoginMode.LOGIN) LoginMode.FORGOT else LoginMode.LOGIN
            },
            onSwitchRegister = {
                mode = if (mode == LoginMode.LOGIN) LoginMode.REGISTER else LoginMode.LOGIN
            },
            modifier = Modifier.padding(innerPadding),
            enabled = formEnabled,
        )
    }
}

@Composable
private fun LoginScreen(
    state: LoginFormState,
    handler: LoginFormHandler,
    mode: LoginMode,
    onSwitchReset: () -> Unit,
    onSwitchRegister: () -> Unit,
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
                mode = mode,
                onSwitchReset = onSwitchReset,
                onSwitchRegister = onSwitchRegister,
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
    mode: LoginMode,
    onSwitchReset: () -> Unit,
    onSwitchRegister: () -> Unit,
    modifier: Modifier = Modifier,
    enabled: Boolean = true,
) {
    Column(horizontalAlignment = Alignment.CenterHorizontally, modifier = modifier) {
        LoginFields(
            state = state,
            handler = handler,
            mode = mode,
            modifier = Modifier.widthIn(max = 320.dp),
            enabled = enabled,
        )
        LoginButtons(
            handler = handler,
            mode = mode,
            onSwitchReset = onSwitchReset,
            onSwitchRegister = onSwitchRegister,
            modifier = Modifier.widthIn(max = 320.dp),
            enabled = enabled,
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
                enabled = enabled,
                labelId = R.string.name,
                leadingIconImage = Icons.Filled.Person,
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
    onSwitchReset: () -> Unit,
    onSwitchRegister: () -> Unit,
    modifier: Modifier = Modifier,
    enabled: Boolean = true,
) {
    Column(verticalArrangement = Arrangement.spacedBy(2.dp), modifier = modifier) {
        AnimatedVisibility(mode != LoginMode.REGISTER, modifier = Modifier.align(Alignment.End)) {
            ClickableText(text = stringResource(R.string.forgot_password), onClick = onSwitchReset)
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
        AnimatedVisibility(
            mode != LoginMode.FORGOT,
            modifier = Modifier.align(Alignment.CenterHorizontally),
        ) {
            ClickableText(
                supportingText =
                    stringResource(
                        if (mode == LoginMode.REGISTER) R.string.already_registered
                        else R.string.not_registered
                    ),
                text =
                    stringResource(
                        if (mode == LoginMode.REGISTER) R.string.login_instead
                        else R.string.register_instead
                    ),
                onClick = onSwitchRegister,
            )
        }
    }
}

@Preview
@Composable
private fun PreviewLoginScreen() {
    ParkEasyTheme {
        LoginScreen(
            state = LoginFormState(),
            handler = LoginFormHandler(),
            mode = LoginMode.LOGIN,
            onSwitchReset = {},
            onSwitchRegister = {},
        )
    }
}
