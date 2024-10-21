package io.github.parkwithease.parkeasy.ui.list

import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.hilt.navigation.compose.hiltViewModel

@Composable
fun ListScreen(
    modifier: Modifier = Modifier,
    @Suppress("unused") viewModel: ListViewModel = hiltViewModel<ListViewModel>(),
) {
    Surface(modifier) { Text("List") }
}
