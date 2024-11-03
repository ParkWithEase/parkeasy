package io.github.parkwithease.parkeasy.data.remote

import io.github.parkwithease.parkeasy.model.Spot

interface SpotRepository {
    /**
     * Gets logged in user's spots.
     *
     * @return a List of Spots.
     */
    suspend fun getSpots(): List<Spot>
}
