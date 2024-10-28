package io.github.parkwithease.parkeasy.ui.cars

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material3.Card
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.ListItem
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.material3.pulltorefresh.PullToRefreshBox
import androidx.compose.material3.pulltorefresh.PullToRefreshState
import androidx.compose.material3.pulltorefresh.rememberPullToRefreshState
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import io.github.parkwithease.parkeasy.model.Car

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun CarsScreen(
    showSnackbar: suspend (String, String?) -> Boolean,
    onSelectCar: (Car) -> Unit,
    modifier: Modifier = Modifier,
    viewModel: CarsViewModel =
        hiltViewModel<CarsViewModel, CarsViewModel.Factory>(
            creationCallback = { factory -> factory.create(showSnackbar = showSnackbar) }
        ),
) {
    val cars by viewModel.cars.collectAsState()
    val isRefreshing by viewModel.isRefreshing.collectAsState()
    val state = rememberPullToRefreshState()
    CarsScreen(cars, isRefreshing, state, viewModel::onRefresh, onSelectCar, modifier)
    LaunchedEffect(Unit) { viewModel.onRefresh() }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun CarsScreen(
    cars: List<Car>,
    isRefreshing: Boolean,
    state: PullToRefreshState,
    onRefresh: () -> Unit,
    onSelectCar: (Car) -> Unit,
    modifier: Modifier = Modifier,
) {
    Surface(modifier) {
        Box {
            PullToRefreshBox(
                isRefreshing = isRefreshing,
                onRefresh = onRefresh,
                state = state,
                modifier = Modifier.padding(4.dp),
            ) {
                LazyColumn(
                    Modifier.fillMaxSize(),
                    horizontalAlignment = Alignment.CenterHorizontally,
                ) {
                    items(cars) { car ->
                        ListItem({
                            Card(
                                modifier = Modifier.fillMaxWidth().clickable { onSelectCar(car) }
                            ) {
                                Text(car.details.color, Modifier.padding(8.dp, 8.dp, 0.dp, 0.dp))
                                Text(car.details.model, Modifier.padding(8.dp, 8.dp, 0.dp, 0.dp))
                                Text(car.details.make, Modifier.padding(8.dp, 8.dp, 0.dp, 0.dp))
                                Text(
                                    car.details.licensePlate,
                                    Modifier.padding(8.dp, 8.dp, 0.dp, 8.dp),
                                )
                            }
                        })
                    }
                }
            }
        }
    }
}
