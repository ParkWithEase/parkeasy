package io.github.parkwithease.parkeasy.model

import kotlinx.datetime.LocalDateTime
import kotlinx.datetime.format
import kotlinx.datetime.format.char
import kotlinx.serialization.KSerializer
import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable
import kotlinx.serialization.descriptors.PrimitiveKind
import kotlinx.serialization.descriptors.PrimitiveSerialDescriptor
import kotlinx.serialization.descriptors.SerialDescriptor
import kotlinx.serialization.encoding.Decoder
import kotlinx.serialization.encoding.Encoder

@Serializable
data class Spot(
    val availability: List<TimeUnit> = emptyList(),
    val features: SpotFeatures = SpotFeatures(),
    val id: String = "",
    val location: SpotLocation = SpotLocation(),
    @SerialName("price_per_hour") val pricePerHour: Double = 0.0,
    @SerialName("distance_to_location") val distanceToLocation: Double = 0.0,
)

@Serializable
data class SpotFeatures(
    @SerialName("charging_station") val chargingStation: Boolean = false,
    @SerialName("plug_in") val plugIn: Boolean = false,
    val shelter: Boolean = false,
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

@Serializable
data class TimeUnit(
    @Serializable(with = LocalDateTimeRFC3339Serializer::class)
    @SerialName("start_time")
    val startTime: LocalDateTime,
    @Serializable(with = LocalDateTimeRFC3339Serializer::class)
    @SerialName("end_time")
    val endTime: LocalDateTime,
    val status: String = "",
) {
    companion object {
        const val BOOKED = "booked"
        const val AVAILABLE = "available"
    }
}

object LocalDateTimeRFC3339Serializer : KSerializer<LocalDateTime> {
    private val format =
        LocalDateTime.Format {
            year()
            char('-')
            monthNumber()
            char('-')
            dayOfMonth()
            char('T')
            hour()
            char(':')
            minute()
            char(':')
            second()
            char('Z')
        }

    override val descriptor: SerialDescriptor =
        PrimitiveSerialDescriptor("kotlinx.datetime.LocalDateTime", PrimitiveKind.STRING)

    override fun serialize(encoder: Encoder, value: LocalDateTime) {
        encoder.encodeString(value.format(format))
    }

    override fun deserialize(decoder: Decoder): LocalDateTime =
        LocalDateTime.parse(decoder.decodeString(), format)
}
