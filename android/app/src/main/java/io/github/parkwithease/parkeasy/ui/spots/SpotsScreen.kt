package io.github.parkwithease.parkeasy.ui.spots

import androidx.compose.foundation.Image
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.heightIn
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Card
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import io.github.parkwithease.parkeasy.R
import io.github.parkwithease.parkeasy.common.PullToRefreshBox
import io.github.parkwithease.parkeasy.model.Spot
import io.github.parkwithease.parkeasy.ui.theme.Typography

@Composable
fun SpotsScreen(
    onSpotClick: (Spot) -> Unit,
    modifier: Modifier = Modifier,
    viewModel: SpotsViewModel = hiltViewModel(),
) {
    val spots by viewModel.spots.collectAsState()
    val isRefreshing by viewModel.isRefreshing.collectAsState()

    LaunchedEffect(Unit) { viewModel.onRefresh() }
    SpotsScreen(spots, onSpotClick, isRefreshing, viewModel::onRefresh, modifier)
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun SpotsScreen(
    spots: List<Spot>,
    onSpotClick: (Spot) -> Unit,
    isRefreshing: Boolean,
    onRefresh: () -> Unit,
    modifier: Modifier = Modifier,
) {
    Scaffold(modifier = modifier) { innerPadding ->
        Surface(Modifier.padding(innerPadding)) {
            PullToRefreshBox(
                items = spots,
                onClick = onSpotClick,
                isRefreshing = isRefreshing,
                onRefresh = onRefresh,
                modifier = Modifier.padding(4.dp),
            ) { spot, onClick ->
                SpotCard(spot, onClick)
            }
        }
    }
}

@Composable
fun SpotCard(spot: Spot, onClick: (Spot) -> Unit, modifier: Modifier = Modifier) {
    Card(onClick = { onClick(spot) }, modifier = modifier.fillMaxWidth().padding(4.dp, 0.dp)) {
        Row(modifier = Modifier.padding(8.dp)) {
            Column(modifier = Modifier.weight(1f)) {
                Image(
                    painter = painterResource(R.drawable.wordmark_outlined),
                    contentDescription = null,
                    modifier = Modifier.heightIn(max = 64.dp),
                )
            }
            Column(horizontalAlignment = Alignment.End, modifier = Modifier.weight(1f)) {
                Text(text = spot.location.streetAddress, style = Typography.headlineMedium)
                Text(spot.location.city + ' ' + spot.location.state)
                Text(spot.location.countryCode + ' ' + spot.location.postalCode)
            }
        }
    }
}
