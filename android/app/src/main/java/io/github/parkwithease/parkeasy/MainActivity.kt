package io.github.parkwithease.parkeasy

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.foundation.Image
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
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
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
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
import io.github.parkwithease.parkeasy.HttpService.login
import io.github.parkwithease.parkeasy.ui.theme.ParkEasyTheme
import kotlinx.coroutines.runBlocking

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        setContent {
            ParkEasyTheme {
                LoginPage()
            }
        }
    }
}

@Preview(showSystemUi = true, device = "id:pixel_8")
@Composable
fun LoginPagePreview() {
    ParkEasyTheme {
        LoginPage()
    }
}

@Composable
fun LoginPage() {
    Surface {
        Column {
            Row(
                verticalAlignment = Alignment.Bottom,
                horizontalArrangement = Arrangement.Center,
                modifier = Modifier
                    .weight(1f)
                    .fillMaxSize()
            ) {
                Image(
                    painter = painterResource(R.drawable.outlined_stacked),
                    contentDescription = stringResource(R.string.logo),
                    modifier = Modifier.size(280.dp)
                )
            }
            Row(
                modifier = Modifier
                    .weight(1f)
                    .fillMaxSize()
            ) {
                LoginForm()
            }
        }
    }
}

@Composable
fun LoginForm() {
    var email by remember { mutableStateOf("") }
    var password by remember { mutableStateOf("") }
    val onEmailChange = { input: String ->
        email = input
    }
    val onPasswordChange = { input: String ->
        password = input
    }
    val onLoginClick = { }
    Column(
        horizontalAlignment = Alignment.CenterHorizontally,
        modifier = Modifier
            .fillMaxSize()
    ) {
        EmailField(
            email,
            onEmailChange
        )
        Spacer(Modifier.size(4.dp))
        PasswordField(
            password,
            onPasswordChange
        )
        Spacer(Modifier.size(4.dp))
        Row(
            modifier = Modifier.width(280.dp)
        ) {
            Spacer(Modifier.weight(1f)) // Space for future register button
            Spacer(Modifier.size(4.dp))

            LoginButton(
                modifier = Modifier.weight(1f),
                onLoginClick
            )
        }
    }
}

@Composable
fun EmailField(
    text: String,
    onValueChange: (String) -> Unit,
) {
    OutlinedTextField(
        value = text,
        onValueChange = {
            onValueChange(it)
        },
        label = {
            Text(stringResource(R.string.email))
        },
        leadingIcon = {
            Icon(
                imageVector = ImageVector.vectorResource(R.drawable.email),
                contentDescription = stringResource(R.string.emailicon))
        },
        keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Email),
        singleLine = true
    )
}

@Composable
fun PasswordField(
    text: String,
    onValueChange: (String) -> Unit
) {
    OutlinedTextField(
        value = text,
        onValueChange = {
            onValueChange(it)
        },
        label = {
            Text(stringResource(R.string.password))
        },
        leadingIcon = {
            Icon(
                imageVector = ImageVector.vectorResource(R.drawable.password),
                contentDescription = stringResource(R.string.passwordicon))
        },
        visualTransformation = PasswordVisualTransformation(),
        keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Password),
        singleLine = true
    )
}

@Composable
fun LoginButton(
    modifier: Modifier = Modifier,
    onClick: () -> Unit
    ) {
    Button(
        onClick = { onClick() },
        modifier = modifier
    ) {
        Text(stringResource(R.string.login))
    }
}
