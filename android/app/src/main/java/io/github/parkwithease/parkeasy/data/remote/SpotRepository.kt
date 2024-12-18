package io.github.parkwithease.parkeasy.data.remote

import io.github.parkwithease.parkeasy.model.Booking
import io.github.parkwithease.parkeasy.model.BookingHistory
import io.github.parkwithease.parkeasy.model.Spot
import io.github.parkwithease.parkeasy.model.TimeUnit

interface SpotRepository {
    /**
     * Gets logged in user's [Spot]s.
     *
     * @return a [Result] containing [List] of [Spot]s on success, failing [Result] otherwise.
     */
    suspend fun getSpots(): Result<List<Spot>>

    /**
     * Creates a [Spot] in the repository.
     *
     * @param spot - the [Spot] to add to the repository.
     * @return whether the creation was successful or not.
     */
    suspend fun createSpot(spot: Spot): Result<Unit>

    /**
     * Gets [Spot]s around the location (latitude, longitude).
     *
     * @return a [Result] containing [List] of [Spot]s on success, failing [Result] otherwise.
     */
    suspend fun getSpotsAround(
        latitude: Double,
        longitude: Double,
        distance: Int,
    ): Result<List<Spot>>

    /**
     * Gets a [List] of [TimeUnit]s for a given [Spot], given its id.
     *
     * @param id - the id of the [Spot].
     * @return a [Result] containing [List] of [TimeUnit]s on success, failing [Result] otherwise.
     */
    suspend fun getSpotAvailability(id: String): Result<List<TimeUnit>>

    /**
     * Creates a booking for a parking spot.
     *
     * @return whether the creation was successful or not.
     */
    suspend fun bookSpot(spotId: String, booking: Booking): Result<Unit>

    /**
     * Gets logged in user's booking history as [BookingHistory]s.
     *
     * @return a [Result] containing [List] of [BookingHistory]s on success, failing [Result]
     *   otherwise.
     */
    suspend fun getBookings(): Result<List<BookingHistory>>

    /**
     * Gets logged in user's leasing history as [BookingHistory]s.
     *
     * @return a [Result] containing [List] of [BookingHistory]s on success, failing [Result]
     *   otherwise.
     */
    suspend fun getLeasings(): Result<List<BookingHistory>>
}
