package io.github.parkwithease.parkeasy.model

import io.github.parkwithease.parkeasy.ui.common.LocalDateTimeRFC3339Serializer
import kotlinx.datetime.LocalDateTime
import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class BookingHistory(
    @Serializable(with = LocalDateTimeRFC3339Serializer::class)
    @SerialName("booking_time")
    val bookingTime: LocalDateTime = LocalDateTime(0, 0, 0, 0, 0),
    @SerialName("car_details") val carDetails: CarDetails = CarDetails(),
    @SerialName("car_id") val carId: String = "",
    val id: String = "",
    @SerialName("paid_amount") val paidAmount: Double = 0.0,
    @SerialName("parkingspot_id") val parkingSpotId: String = "",
    @SerialName("parkingspot_location") val parkingSpotLocation: SpotLocation = SpotLocation(),
)
