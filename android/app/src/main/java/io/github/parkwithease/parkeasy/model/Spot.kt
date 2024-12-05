package io.github.parkwithease.parkeasy.model

import java.time.format.DateTimeFormatter
import kotlinx.datetime.LocalDateTime
import kotlinx.datetime.toJavaLocalDateTime
import kotlinx.datetime.toKotlinLocalDateTime
import kotlinx.serialization.KSerializer
import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable
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
    @Serializable(with = LocalDateTimeToRFC3339::class)
    @SerialName("start_time")
    val startTime: LocalDateTime,
    @Serializable(with = LocalDateTimeToRFC3339::class)
    @SerialName("end_time")
    val endTime: LocalDateTime,
    val status: String = "",
) {
    companion object {
        const val BOOKED = "booked"
        const val AVAILABLE = "available"
    }
}

object LocalDateTimeToRFC3339 : KSerializer<LocalDateTime> {
    private val format = DateTimeFormatter.ofPattern("yyyy-MM-dd'T'HH:mm:ss.SSS'Z'")

    override val descriptor: SerialDescriptor
        get() = TODO("Not yet implemented") // Seems like overkill to implement

    override fun serialize(encoder: Encoder, value: LocalDateTime) =
        encoder.encodeString(value.toJavaLocalDateTime().format(format))

    override fun deserialize(decoder: Decoder): LocalDateTime =
        java.time.LocalDateTime.parse(decoder.decodeString(), format).toKotlinLocalDateTime()
}
