package io.github.parkwithease.parkeasy.ui.history

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Card
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import io.github.parkwithease.parkeasy.model.BookingHistory
import io.github.parkwithease.parkeasy.ui.common.CarDetailsText
import io.github.parkwithease.parkeasy.ui.common.SpotLocationText
import io.github.parkwithease.parkeasy.ui.common.toReadable

@Suppress("DefaultLocale", "detekt:ImplicitDefaultLocale")
@Composable
fun HistoryCard(
    booking: BookingHistory,
    onClick: (BookingHistory) -> Unit,
    modifier: Modifier = Modifier,
) {
    Card(onClick = { onClick(booking) }, modifier = modifier.fillMaxWidth()) {
        Column {
            Row(modifier = Modifier.padding(8.dp)) {
                SpotLocationText(
                    spotLocation = booking.parkingSpotLocation,
                    modifier = Modifier.weight(1f),
                )
                CarDetailsText(
                    carDetails = booking.carDetails,
                    modifier = Modifier.weight(1f),
                    horizontalAlignment = Alignment.End,
                )
            }
            Row(modifier = Modifier.padding(8.dp)) {
                Column(modifier = Modifier.weight(1f)) { Text(booking.bookingTime.toReadable()) }
                Column(
                    modifier = Modifier.weight(1f),
                    verticalArrangement = Arrangement.Bottom,
                    horizontalAlignment = Alignment.End,
                ) {
                    Text(
                        text = String.format("$%.2f", booking.paidAmount),
                        style = MaterialTheme.typography.headlineLarge,
                    )
                }
            }
        }
    }
}
