package io.github.parkwithease.parkeasy.ui.history.bookings

import androidx.navigation.NavController
import androidx.navigation.NavGraphBuilder
import androidx.navigation.compose.composable
import io.github.parkwithease.parkeasy.ui.common.enterAnimation
import io.github.parkwithease.parkeasy.ui.common.exitAnimation

private const val BookingsRoute = "bookings"

fun NavGraphBuilder.bookingsScreen() {
    composable(
        route = BookingsRoute,
        enterTransition = { enterAnimation() },
        exitTransition = { exitAnimation() },
    ) {
        BookingsScreen()
    }
}

fun NavController.navigateToBookings() = this.navigate(BookingsRoute)
