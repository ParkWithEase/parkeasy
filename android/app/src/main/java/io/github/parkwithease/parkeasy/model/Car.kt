package io.github.parkwithease.parkeasy.model

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class Car(
    val details: CarDetails = CarDetails(),
    val id: String = "",
    @SerialName("\$schema") val schema: String = "", // ContentNegotiation needs this.
)

@Serializable
data class CarDetails(
    var color: String = "",
    @SerialName("license_plate") var licensePlate: String = "",
    var make: String = "",
    var model: String = "",
)
