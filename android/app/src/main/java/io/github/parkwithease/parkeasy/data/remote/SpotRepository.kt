package io.github.parkwithease.parkeasy.data.remote

import io.github.parkwithease.parkeasy.model.Spot

interface SpotRepository {
    /**
     * Gets logged in user's spots.
     *
     * @return a Result of List of Spots on success, failing Result otherwise.
     */
    suspend fun getSpots(): Result<List<Spot>>

    /**
     * Creates a spot in the repository.
     *
     * @param spot the Spot to add to the repository.
     * @return whether the creation was successful or not.
     */
    suspend fun createSpot(spot: Spot): Result<Unit>
}
