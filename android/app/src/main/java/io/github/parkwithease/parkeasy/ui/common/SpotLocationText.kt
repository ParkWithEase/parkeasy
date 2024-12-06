package io.github.parkwithease.parkeasy.ui.common

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.material3.MaterialTheme.typography
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import io.github.parkwithease.parkeasy.model.SpotLocation

@Composable
fun SpotLocationText(
    spotLocation: SpotLocation,
    modifier: Modifier = Modifier,
    verticalArrangement: Arrangement.Vertical = Arrangement.Top,
    horizontalAlignment: Alignment.Horizontal = Alignment.Start,
) {
    Column(
        modifier = modifier,
        verticalArrangement = verticalArrangement,
        horizontalAlignment = horizontalAlignment,
    ) {
        Text(text = spotLocation.streetAddress, style = typography.titleMedium)
        Text(text = spotLocation.city + ' ' + spotLocation.state, style = typography.titleSmall)
        Text(
            text = spotLocation.countryCode + ' ' + spotLocation.postalCode,
            style = typography.titleSmall,
        )
    }
}
