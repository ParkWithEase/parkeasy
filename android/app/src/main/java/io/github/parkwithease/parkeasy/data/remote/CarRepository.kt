package io.github.parkwithease.parkeasy.data.remote

import io.github.parkwithease.parkeasy.model.Car
import io.github.parkwithease.parkeasy.model.CarDetails

interface CarRepository {
    /**
     * Gets logged in user's [Car]s.
     *
     * @return a [Result] containing [List] of [Car]s on success, failing [Result] otherwise.
     */
    suspend fun getCars(): Result<List<Car>>

    /**
     * Deletes a [Car] in the database by id.
     *
     * @param id - the [String] id of the car to delete.
     * @return whether deleting the car was successful or not.
     */
    suspend fun deleteCar(id: String): Result<Unit>

    /**
     * Gets a [Car]'s details by id.
     *
     * @param id - the [String] id of the car to get.
     * @return [Result] with the [Car] if exists and authorized, failing [Result] otherwise.
     */
    suspend fun getCarInfo(id: String): Result<Car>

    /**
     * Updates a [Car]'s details.
     *
     * @param car - the [Car] and its updated info.
     * @return whether the update was successful or not.
     */
    suspend fun updateCar(car: Car): Result<Unit>

    /**
     * Creates a [Car] in the repository.
     *
     * @param car - the [Car] to add to the repository.
     * @return whether the creation was successful or not.
     */
    suspend fun createCar(car: CarDetails): Result<Unit>
}
