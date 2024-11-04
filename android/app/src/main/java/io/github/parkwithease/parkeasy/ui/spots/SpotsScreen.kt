package io.github.parkwithease.parkeasy.ui.spots

import androidx.compose.foundation.Image
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.heightIn
import androidx.compose.foundation.layout.imePadding
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.widthIn
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Add
import androidx.compose.material3.Button
import androidx.compose.material3.Card
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.FloatingActionButton
import androidx.compose.material3.Icon
import androidx.compose.material3.ModalBottomSheet
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.material3.rememberModalBottomSheetState
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.saveable.rememberSaveable
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import io.github.parkwithease.parkeasy.R
import io.github.parkwithease.parkeasy.common.PullToRefreshBox
import io.github.parkwithease.parkeasy.model.EditMode
import io.github.parkwithease.parkeasy.model.Spot

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun SpotsScreen(modifier: Modifier = Modifier, viewModel: SpotsViewModel = hiltViewModel()) {
    val spots by viewModel.spots.collectAsState()
    val isRefreshing by viewModel.isRefreshing.collectAsState()
    var editMode by rememberSaveable { mutableStateOf(EditMode.ADD) }

    var openBottomSheet by rememberSaveable { mutableStateOf(false) }
    val skipPartiallyExpanded by rememberSaveable { mutableStateOf(false) }
    val bottomSheetState =
        rememberModalBottomSheetState(skipPartiallyExpanded = skipPartiallyExpanded)

    LaunchedEffect(Unit) { viewModel.onRefresh() }
    SpotsScreen(
        spots,
        { spot ->
            viewModel.onStreetAddressChange(spot.location.streetAddress)
            viewModel.onCityChange(spot.location.city)
            viewModel.onStateChange(spot.location.state)
            viewModel.onCountryCodeChange(spot.location.countryCode)
            editMode = EditMode.EDIT
            openBottomSheet = true
        },
        {
            viewModel.onStreetAddressChange("")
            viewModel.onCityChange("")
            viewModel.onStateChange("")
            viewModel.onCountryCodeChange("")
            editMode = EditMode.ADD
            openBottomSheet = true
        },
        isRefreshing,
        viewModel::onRefresh,
        modifier,
    )
    if (openBottomSheet) {
        ModalBottomSheet(
            onDismissRequest = { openBottomSheet = false },
            sheetState = bottomSheetState,
        ) {
            Column(
                horizontalAlignment = Alignment.CenterHorizontally,
                verticalArrangement = Arrangement.Center,
                modifier =
                    Modifier.padding(horizontal = 16.dp)
                        .fillMaxWidth()
                        .imePadding()
                        .verticalScroll(rememberScrollState(), reverseScrolling = true),
            ) {
                AddSpotScreen(
                    viewModel.formState,
                    viewModel::onStreetAddressChange,
                    viewModel::onCityChange,
                    viewModel::onStateChange,
                    viewModel::onCountryCodeChange,
                    viewModel::onPostalCodeChange,
                    viewModel::onAddSpotClick,
                )
            }
        }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun SpotsScreen(
    spots: List<Spot>,
    onSpotClick: (Spot) -> Unit,
    onShowAddSpotClick: () -> Unit,
    isRefreshing: Boolean,
    onRefresh: () -> Unit,
    modifier: Modifier = Modifier,
) {
    Scaffold(
        floatingActionButton = { AddSpotButton(onShowAddSpotClick = onShowAddSpotClick) },
        modifier = modifier,
    ) { innerPadding ->
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
                Text(spot.location.streetAddress)
                Text(spot.location.city + ' ' + spot.location.state)
                Text(spot.location.countryCode + ' ' + spot.location.postalCode)
            }
        }
    }
}

@Composable
fun AddSpotButton(onShowAddSpotClick: () -> Unit, modifier: Modifier = Modifier) {
    FloatingActionButton(onClick = onShowAddSpotClick, modifier) {
        Icon(imageVector = Icons.Filled.Add, contentDescription = stringResource(R.string.add_spot))
    }
}

@Composable
fun AddSpotScreen(
    state: AddSpotFormState,
    onStreetAddressChange: (String) -> Unit,
    onCityChange: (String) -> Unit,
    onStateChange: (String) -> Unit,
    onCountryCodeChange: (String) -> Unit,
    onPostalCodeChange: (String) -> Unit,
    onAddSpotClick: () -> Unit,
    modifier: Modifier = Modifier,
) {
    Column(
        verticalArrangement = Arrangement.spacedBy(2.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
        modifier = modifier.widthIn(max = 320.dp),
    ) {
        OutlinedTextField(
            value = state.streetAddress.value,
            onValueChange = onStreetAddressChange,
            label = { Text(stringResource(R.string.street_address)) },
            modifier = Modifier.fillMaxWidth(),
        )
        OutlinedTextField(
            value = state.city.value,
            onValueChange = onCityChange,
            label = { Text(stringResource(R.string.city)) },
            modifier = Modifier.fillMaxWidth(),
        )
        OutlinedTextField(
            value = state.state.value,
            onValueChange = onStateChange,
            label = { Text(stringResource(R.string.state)) },
            modifier = Modifier.fillMaxWidth(),
        )
        OutlinedTextField(
            value = state.countryCode.value,
            onValueChange = onCountryCodeChange,
            label = { Text(stringResource(R.string.country)) },
            modifier = Modifier.fillMaxWidth(),
        )
        OutlinedTextField(
            value = state.postalCode.value,
            onValueChange = onPostalCodeChange,
            label = { Text(stringResource(R.string.postal_code)) },
            modifier = Modifier.fillMaxWidth(),
        )
        Button(onClick = onAddSpotClick, modifier = Modifier.fillMaxWidth()) {
            Text(stringResource(R.string.add_spot))
        }
    }
}
