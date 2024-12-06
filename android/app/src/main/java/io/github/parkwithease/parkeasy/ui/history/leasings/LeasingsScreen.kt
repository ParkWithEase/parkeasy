package io.github.parkwithease.parkeasy.ui.history.leasings

import androidx.compose.foundation.layout.padding
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.Scaffold
import androidx.compose.material3.SnackbarHost
import androidx.compose.material3.Surface
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import io.github.parkwithease.parkeasy.model.BookingHistory
import io.github.parkwithease.parkeasy.ui.common.PullToRefreshBox
import io.github.parkwithease.parkeasy.ui.history.HistoryCard

@Suppress("detekt:LongMethod")
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun LeasingsScreen(modifier: Modifier = Modifier, viewModel: LeasingsViewModel = hiltViewModel()) {
    val bookings by viewModel.bookings.collectAsState()
    val isRefreshing by viewModel.isRefreshing.collectAsState()

    LaunchedEffect(Unit) { viewModel.onRefresh() }

    LeasingsScreen(
        bookings = bookings,
        isRefreshing = isRefreshing,
        onRefresh = viewModel::onRefresh,
        snackbarHost = { SnackbarHost(hostState = viewModel.snackbarState) },
        modifier = modifier,
    )
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun LeasingsScreen(
    bookings: List<BookingHistory>,
    isRefreshing: Boolean,
    onRefresh: () -> Unit,
    snackbarHost: @Composable () -> Unit,
    modifier: Modifier = Modifier,
) {
    Scaffold(modifier = modifier, snackbarHost = snackbarHost) { innerPadding ->
        Surface(Modifier.padding(innerPadding)) {
            PullToRefreshBox(
                items = bookings,
                onClick = {},
                isRefreshing = isRefreshing,
                onRefresh = onRefresh,
                modifier = Modifier.padding(4.dp),
            ) { booking, onClick ->
                HistoryCard(booking, onClick)
            }
        }
    }
}
