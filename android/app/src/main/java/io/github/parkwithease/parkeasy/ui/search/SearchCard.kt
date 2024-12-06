package io.github.parkwithease.parkeasy.ui.search

import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Card
import androidx.compose.material3.MaterialTheme.typography
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import io.github.parkwithease.parkeasy.model.Spot
import io.github.parkwithease.parkeasy.ui.common.FeaturesRow
import io.github.parkwithease.parkeasy.ui.common.SpotLocationText

@Suppress("DefaultLocale", "detekt:ImplicitDefaultLocale")
@Composable
fun SearchCard(spot: Spot, onClick: (Spot) -> Unit, modifier: Modifier = Modifier) {
    Card(onClick = { onClick(spot) }, modifier = modifier.fillMaxWidth()) {
        Row(modifier = Modifier.padding(8.dp)) {
            Column(modifier = Modifier.weight(1f)) {
                Text(
                    text = String.format("$%.2f", spot.pricePerHour),
                    style = typography.headlineLarge,
                )
                FeaturesRow(spot.features)
                Text(
                    text = String.format("%.0f metres away", spot.distanceToLocation),
                    style = typography.titleSmall,
                )
            }
            SpotLocationText(
                spotLocation = spot.location,
                modifier = Modifier.weight(1f),
                horizontalAlignment = Alignment.End,
            )
        }
    }
}
