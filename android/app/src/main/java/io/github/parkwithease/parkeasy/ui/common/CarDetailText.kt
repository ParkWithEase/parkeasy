package io.github.parkwithease.parkeasy.ui.common

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import io.github.parkwithease.parkeasy.model.CarDetails

@Composable
fun CarDetailsText(
    carDetails: CarDetails,
    modifier: Modifier = Modifier,
    verticalArrangement: Arrangement.Vertical = Arrangement.Top,
    horizontalAlignment: Alignment.Horizontal = Alignment.Start,
) {
    Column(
        modifier = modifier,
        verticalArrangement = verticalArrangement,
        horizontalAlignment = horizontalAlignment,
    ) {
        Text(text = carDetails.licensePlate, style = MaterialTheme.typography.headlineLarge)
        Text(
            text = carDetails.color + ' ' + carDetails.make + ' ' + carDetails.model,
            style = MaterialTheme.typography.titleSmall,
        )
    }
}
