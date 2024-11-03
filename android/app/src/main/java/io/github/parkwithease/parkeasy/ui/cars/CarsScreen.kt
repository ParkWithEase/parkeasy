package io.github.parkwithease.parkeasy.ui.cars

import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Add
import androidx.compose.material3.Card
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.FloatingActionButton
import androidx.compose.material3.Icon
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import io.github.parkwithease.parkeasy.R
import io.github.parkwithease.parkeasy.common.PullToRefreshBox
import io.github.parkwithease.parkeasy.model.Car

@Composable
fun CarsScreen(
    onCarClick: (Car) -> Unit,
    modifier: Modifier = Modifier,
    viewModel: CarsViewModel = hiltViewModel(),
) {
    val cars by viewModel.cars.collectAsState()
    val isRefreshing by viewModel.isRefreshing.collectAsState()

    LaunchedEffect(Unit) { viewModel.onRefresh() }

    CarsScreen(cars, onCarClick, {}, isRefreshing, viewModel::onRefresh, modifier)
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun CarsScreen(
    cars: List<Car>,
    onCarClick: (Car) -> Unit,
    onAddCar: () -> Unit,
    isRefreshing: Boolean,
    onRefresh: () -> Unit,
    modifier: Modifier = Modifier,
) {
    Scaffold(floatingActionButton = { AddCarButton(onAddCar = onAddCar) }, modifier = modifier) {
        innerPadding ->
        Surface(Modifier.padding(innerPadding)) {
            PullToRefreshBox(
                items = cars,
                onClick = onCarClick,
                isRefreshing = isRefreshing,
                onRefresh = onRefresh,
                modifier = Modifier.padding(4.dp),
                card = { car, onClick -> CarCard(car, onClick) },
            )
        }
    }
}

@Composable
fun CarCard(car: Car, onClick: (Car) -> Unit, modifier: Modifier = Modifier) {
    Card(onClick = { onClick(car) }, modifier = modifier.fillMaxWidth().padding(8.dp, 0.dp)) {
        Text(car.details.color)
        Text(car.details.model)
        Text(car.details.make)
        Text(car.details.licensePlate)
    }
}

@Composable
fun AddCarButton(onAddCar: () -> Unit, modifier: Modifier = Modifier) {
    FloatingActionButton(onClick = onAddCar, modifier) {
        Icon(imageVector = Icons.Filled.Add, contentDescription = stringResource(R.string.add_car))
    }
}
