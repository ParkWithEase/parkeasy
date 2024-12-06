package io.github.parkwithease.parkeasy.ui.history.bookings

import android.content.Intent
import android.net.Uri
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
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import io.github.parkwithease.parkeasy.model.BookingHistory
import io.github.parkwithease.parkeasy.ui.common.PullToRefreshBox
import io.github.parkwithease.parkeasy.ui.history.HistoryCard

@Suppress("detekt:LongMethod")
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun BookingsScreen(modifier: Modifier = Modifier, viewModel: BookingsViewModel = hiltViewModel()) {
    val bookings by viewModel.bookings.collectAsState()
    val isRefreshing by viewModel.isRefreshing.collectAsState()
    val context = LocalContext.current

    LaunchedEffect(Unit) { viewModel.onRefresh() }

    BookingsScreen(
        bookings = bookings,
        onBookingClick = {
            with(it.parkingSpotLocation) {
                val gmmIntentUri =
                    Uri.parse("google.navigation:q=$streetAddress, $city, $state, $countryCode")
                val mapIntent = Intent(Intent.ACTION_VIEW, gmmIntentUri)
                mapIntent.setPackage("com.google.android.apps.maps")
                context.startActivity(mapIntent)
            }
        },
        isRefreshing = isRefreshing,
        onRefresh = viewModel::onRefresh,
        snackbarHost = { SnackbarHost(hostState = viewModel.snackbarState) },
        modifier = modifier,
    )
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun BookingsScreen(
    bookings: List<BookingHistory>,
    onBookingClick: (BookingHistory) -> Unit,
    isRefreshing: Boolean,
    onRefresh: () -> Unit,
    snackbarHost: @Composable () -> Unit,
    modifier: Modifier = Modifier,
) {
    Scaffold(modifier = modifier, snackbarHost = snackbarHost) { innerPadding ->
        Surface(Modifier.padding(innerPadding)) {
            PullToRefreshBox(
                items = bookings,
                onClick = onBookingClick,
                isRefreshing = isRefreshing,
                onRefresh = onRefresh,
                modifier = Modifier.padding(4.dp),
            ) { booking, onClick ->
                HistoryCard(booking, onClick)
            }
        }
    }
}
