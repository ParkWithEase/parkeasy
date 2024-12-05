package io.github.parkwithease.parkeasy.data.remote

import io.github.parkwithease.parkeasy.model.Spot

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
}
