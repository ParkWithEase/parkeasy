package io.github.parkwithease.parkeasy.model

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class Spot(
    val features: SpotFeatures = SpotFeatures(),
    val id: String = "",
    val location: SpotLocation = SpotLocation(),
    @SerialName("price_per_hour") val pricePerHour: Double = 1.0,
)

@Serializable
data class SpotFeatures(
    @SerialName("charging_station") val chargingStation: Boolean = true,
    @SerialName("plug_in") val plugIn: Boolean = true,
    val shelter: Boolean = true,
)

@Serializable
data class SpotLocation(
    val city: String = "",
    @SerialName("country_code") val countryCode: String = "",
    val latitude: Double = 1.0,
    val longitude: Double = 1.0,
    @SerialName("postal_code") val postalCode: String = "",
    val state: String = "",
    @SerialName("street_address") val streetAddress: String = "",
)
