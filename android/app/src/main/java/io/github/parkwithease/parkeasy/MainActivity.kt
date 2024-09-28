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
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.material3.Icon
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
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
import io.github.parkwithease.parkeasy.ui.theme.ParkEasyTheme

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        setContent {
            ParkEasyTheme {
                LoginForm()
            }
        }
    }
}

@Preview(showSystemUi = true, device = "id:pixel_8")
@Composable
fun LoginFormPreview() {
    ParkEasyTheme {
        LoginForm()
    }
}

@Composable
fun LoginForm() {
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
                Column(
                    horizontalAlignment = Alignment.CenterHorizontally,
                    modifier = Modifier
                        .fillMaxSize()
                ) {
                    EmailField()
                    Spacer(Modifier.size(4.dp))
                    PasswordField()
                    Spacer(Modifier.size(4.dp))
                }
            }
        }
    }
}

@Composable
fun EmailField() {
    val email = remember { mutableStateOf("") }

    OutlinedTextField(
        value = email.value,
        onValueChange = {email.value = it},
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
fun PasswordField() {
    val password = remember { mutableStateOf("") }

    OutlinedTextField(
        value = password.value,
        onValueChange = {password.value = it},
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
