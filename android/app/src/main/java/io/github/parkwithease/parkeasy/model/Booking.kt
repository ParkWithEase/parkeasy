package io.github.parkwithease.parkeasy.model

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class Booking(
    @SerialName("car_id") val carId: String = "",
    @SerialName("booked_times") val bookedTimes: List<TimeUnit> = emptyList(),
)
