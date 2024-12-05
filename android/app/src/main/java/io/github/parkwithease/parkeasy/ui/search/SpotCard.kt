package io.github.parkwithease.parkeasy.ui.search

import androidx.compose.foundation.Image
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.width
import androidx.compose.material3.Card
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.unit.dp
import io.github.parkwithease.parkeasy.R
import io.github.parkwithease.parkeasy.model.Spot

@Suppress("DefaultLocale", "detekt:ImplicitDefaultLocale")
@Composable
fun SpotCard(spot: Spot, onClick: (Spot) -> Unit, modifier: Modifier = Modifier) {
    Card(onClick = { onClick(spot) }, modifier = modifier.fillMaxWidth().padding(4.dp, 0.dp)) {
        Row(modifier = Modifier.padding(8.dp)) {
            Column(modifier = Modifier.weight(1f)) {
                Text(
                    text = String.format("$ %.2f", spot.pricePerHour),
                    style = MaterialTheme.typography.titleLarge,
                )
                Row(Modifier.height(24.dp)) {
                    if (spot.features.chargingStation)
                        Image(
                            painter = painterResource(R.drawable.charging_station),
                            contentDescription = null,
                            modifier = Modifier.width(24.dp),
                        )
                    if (spot.features.plugIn)
                        Image(
                            painter = painterResource(R.drawable.plug_in),
                            contentDescription = null,
                            modifier = Modifier.width(24.dp),
                        )
                    if (spot.features.shelter)
                        Image(
                            painter = painterResource(R.drawable.shelter),
                            contentDescription = null,
                            modifier = Modifier.width(24.dp),
                        )
                }
                Text(
                    text = String.format("%.0f meters away", spot.distanceToLocation),
                    style = MaterialTheme.typography.titleSmall,
                )
            }
            Column(modifier = Modifier.weight(1f), horizontalAlignment = Alignment.End) {
                Text(spot.location.streetAddress)
                Text(spot.location.city + ' ' + spot.location.state)
                Text(spot.location.countryCode + ' ' + spot.location.postalCode)
            }
        }
    }
}
