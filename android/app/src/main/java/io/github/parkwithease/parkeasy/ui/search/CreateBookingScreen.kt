package io.github.parkwithease.parkeasy.ui.search

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.imePadding
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.Button
import androidx.compose.material3.DropdownMenuItem
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.ExposedDropdownMenuBox
import androidx.compose.material3.ExposedDropdownMenuDefaults
import androidx.compose.material3.FilterChip
import androidx.compose.material3.MenuAnchorType
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.unit.dp
import io.github.parkwithease.parkeasy.R
import io.github.parkwithease.parkeasy.model.Car
import io.github.parkwithease.parkeasy.model.FieldState
import io.github.parkwithease.parkeasy.ui.common.ParkEasyTextField
import io.github.parkwithease.parkeasy.ui.common.TimeGrid
import kotlin.collections.forEach

@OptIn(ExperimentalMaterial3Api::class)
@Suppress("detekt:LongMethod", "DefaultLocale", "detekt:ImplicitDefaultLocale")
@Composable
fun CreateBookingScreen(
    cars: List<Car>,
    state: CreateState,
    handler: CreateHandler,
    getSelectedIds: () -> Set<Int>,
    disabledIds: Set<Int>,
    modifier: Modifier = Modifier,
) {
    var expanded by remember { mutableStateOf(false) }

    Surface(modifier.fillMaxSize()) {
        Column(
            modifier =
                Modifier.imePadding()
                    .padding(32.dp)
                    .verticalScroll(rememberScrollState(), reverseScrolling = true),
            verticalArrangement = Arrangement.spacedBy(2.dp, Alignment.CenterVertically),
            horizontalAlignment = Alignment.CenterHorizontally,
        ) {
            ExposedDropdownMenuBox(expanded = expanded, onExpandedChange = { expanded = it }) {
                ParkEasyTextField(
                    state = FieldState(state.selectedCar.value.details.licensePlate),
                    onValueChange = {},
                    modifier =
                        Modifier.fillMaxWidth()
                            .menuAnchor(MenuAnchorType.PrimaryNotEditable)
                            .clickable { expanded = true },
                    enabled = false,
                    visuallyEnabled = true,
                    readOnly = true,
                    labelId = R.string.select_car,
                    trailingIcon = { ExposedDropdownMenuDefaults.TrailingIcon(expanded = expanded) },
                )
                ExposedDropdownMenu(expanded = expanded, onDismissRequest = { expanded = false }) {
                    cars.forEach {
                        DropdownMenuItem(
                            text = { Text(it.details.licensePlate) },
                            onClick = {
                                handler.onCarChange(it)
                                expanded = false
                            },
                        )
                    }
                }
            }
            ParkEasyTextField(
                state = FieldState(state.selectedSpot.value.location.streetAddress),
                onValueChange = {},
                modifier = Modifier.fillMaxWidth(),
                enabled = false,
                visuallyEnabled = true,
                readOnly = true,
                labelId = R.string.street_address,
            )
            Row(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                ParkEasyTextField(
                    state = FieldState(state.selectedSpot.value.location.city),
                    onValueChange = {},
                    modifier = Modifier.weight(1f),
                    enabled = false,
                    visuallyEnabled = true,
                    readOnly = true,
                    labelId = R.string.city,
                )
                ParkEasyTextField(
                    state = FieldState(state.selectedSpot.value.location.state),
                    onValueChange = {},
                    modifier = Modifier.weight(1f),
                    enabled = false,
                    visuallyEnabled = true,
                    readOnly = true,
                    labelId = R.string.state,
                )
            }
            Row(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                ParkEasyTextField(
                    state = FieldState(state.selectedSpot.value.location.countryCode),
                    onValueChange = {},
                    modifier = Modifier.weight(1f),
                    enabled = false,
                    visuallyEnabled = true,
                    readOnly = true,
                    labelId = R.string.country,
                )
                ParkEasyTextField(
                    state = FieldState(state.selectedSpot.value.location.postalCode),
                    onValueChange = {},
                    modifier = Modifier.weight(1f),
                    enabled = false,
                    visuallyEnabled = true,
                    readOnly = true,
                    labelId = R.string.postal_code,
                )
            }
            ParkEasyTextField(
                state = FieldState(String.format("$ %.2f", state.selectedSpot.value.pricePerHour)),
                onValueChange = {},
                modifier = Modifier.fillMaxWidth(),
                enabled = false,
                visuallyEnabled = true,
                readOnly = true,
                labelId = R.string.price_per_hour,
            )
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.spacedBy(4.dp),
            ) {
                FilterChip(
                    selected = state.selectedSpot.value.features.chargingStation,
                    onClick = {},
                    label = { Text(stringResource(R.string.charging_station)) },
                )
                FilterChip(
                    selected = state.selectedSpot.value.features.plugIn,
                    onClick = {},
                    label = { Text(stringResource(R.string.plug_in)) },
                )
                FilterChip(
                    selected = state.selectedSpot.value.features.plugIn,
                    onClick = {},
                    label = { Text(stringResource(R.string.shelter)) },
                )
            }
            TimeGrid(getSelectedIds, disabledIds, handler.onAddTime, handler.onRemoveTime)
            ParkEasyTextField(
                state =
                    FieldState(
                        String.format(
                            "$ %.2f",
                            state.selectedSpot.value.pricePerHour * getSelectedIds().size / 2,
                        )
                    ),
                onValueChange = {},
                modifier = Modifier.fillMaxWidth(),
                enabled = false,
                visuallyEnabled = true,
                readOnly = true,
                labelId = R.string.total_price,
            )
            Button(onClick = handler.onCreateBookingClick, modifier = Modifier.fillMaxWidth()) {
                Text(stringResource(R.string.create_booking))
            }
        }
    }
}
