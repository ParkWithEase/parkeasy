package io.github.parkwithease.parkeasy.ui

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Scaffold
import androidx.compose.material3.SnackbarDuration
import androidx.compose.material3.SnackbarHost
import androidx.compose.material3.SnackbarHostState
import androidx.compose.material3.SnackbarResult
import androidx.compose.runtime.remember
import androidx.compose.runtime.rememberCoroutineScope
import androidx.compose.ui.Modifier
import dagger.hilt.android.AndroidEntryPoint
import io.github.parkwithease.parkeasy.ui.navbar.NavBar
import io.github.parkwithease.parkeasy.ui.theme.ParkEasyTheme
import kotlinx.coroutines.launch

@AndroidEntryPoint
class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        setContent {
            ParkEasyTheme {
                val snackbarHostState = remember { SnackbarHostState() }
                val scope = rememberCoroutineScope()
                ObserveAsEvents(flow = SnackbarController.events, snackbarHostState) { event ->
                    scope.launch {
                        snackbarHostState.currentSnackbarData?.dismiss()
                        val result =
                            snackbarHostState.showSnackbar(
                                message = event.message,
                                actionLabel = event.action?.name,
                                duration = SnackbarDuration.Short,
                            )
                        if (result == SnackbarResult.ActionPerformed) {
                            event.action?.action?.invoke()
                        }
                    }
                }
                Scaffold(
                    bottomBar = { NavBar() },
                    snackbarHost = { SnackbarHost(hostState = snackbarHostState) },
                ) { innerPadding ->
                    MainNavGraph(modifier = Modifier.padding(innerPadding))
                }
            }
        }
    }
}
