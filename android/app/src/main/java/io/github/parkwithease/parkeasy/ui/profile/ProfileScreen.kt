package io.github.parkwithease.parkeasy.ui.profile

import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Surface
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import io.github.parkwithease.parkeasy.ui.nav.NavBar

@Composable
fun ProfileScreen(modifier: Modifier = Modifier) {
    Scaffold(bottomBar = { NavBar() }, modifier = modifier) { innerPadding ->
        Surface(Modifier.padding(innerPadding)) {}
    }
}
