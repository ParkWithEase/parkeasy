package io.github.parkwithease.parkeasy.ui

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Scaffold
import androidx.compose.ui.Modifier
import dagger.hilt.android.AndroidEntryPoint
import io.github.parkwithease.parkeasy.ui.navbar.NavBar
import io.github.parkwithease.parkeasy.ui.theme.ParkEasyTheme

@AndroidEntryPoint
class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        setContent { ParkEasyTheme { MainNavGraph() } }
        setContent {
            ParkEasyTheme {
                Scaffold(
                    bottomBar = { NavBar() },
                ) { innerPadding ->
                    MainNavGraph(modifier = Modifier.padding(innerPadding))
                }
            }
        }
    }
}
