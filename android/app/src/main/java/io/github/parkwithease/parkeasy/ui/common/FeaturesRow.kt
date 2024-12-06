package io.github.parkwithease.parkeasy.ui.common

import androidx.compose.foundation.Image
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.width
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.unit.dp
import io.github.parkwithease.parkeasy.R
import io.github.parkwithease.parkeasy.model.SpotFeatures

private val FeatureImageSize = 24.dp

@Composable
fun FeaturesRow(
    features: SpotFeatures,
    modifier: Modifier = Modifier,
    horizontalArrangement: Arrangement.Horizontal = Arrangement.Start,
    verticalAlignment: Alignment.Vertical = Alignment.Top,
) {
    Row(
        modifier = modifier.height(FeatureImageSize),
        horizontalArrangement = horizontalArrangement,
        verticalAlignment = verticalAlignment,
    ) {
        if (features.chargingStation)
            Image(
                painter = painterResource(R.drawable.charging_station),
                contentDescription = null,
                modifier = Modifier.width(FeatureImageSize),
            )
        if (features.plugIn)
            Image(
                painter = painterResource(R.drawable.plug_in),
                contentDescription = null,
                modifier = Modifier.width(FeatureImageSize),
            )
        if (features.shelter)
            Image(
                painter = painterResource(R.drawable.shelter),
                contentDescription = null,
                modifier = Modifier.width(FeatureImageSize),
            )
    }
}
