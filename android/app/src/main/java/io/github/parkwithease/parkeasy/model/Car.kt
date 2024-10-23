package io.github.parkwithease.parkeasy.model

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class Car(
    val details: CarDetails,
    val id: String,
    @SerialName("\$schema") val schema: String = "", // ContentNegotiation needs this.
)

@Serializable
data class CarDetails(
    val color: String = "",
    @SerialName("license_plate") val licensePlate: String = "",
    val make: String = "",
    val model: String = "",
)
