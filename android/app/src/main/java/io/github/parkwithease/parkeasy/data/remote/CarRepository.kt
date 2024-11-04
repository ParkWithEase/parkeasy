package io.github.parkwithease.parkeasy.data.remote

import io.github.parkwithease.parkeasy.model.Car
import io.github.parkwithease.parkeasy.model.CarDetails

interface CarRepository {
    /**
     * Gets logged in user's cars.
     *
     * @return a List of Cars.
     */
    suspend fun getCars(): List<Car>

    /**
     * Deletes a car in the database by id.
     *
     * @return true on success, false on fail.
     */
    suspend fun deleteCar(id: String): Boolean

    /**
     * Gets a car's details by id.
     *
     * @return the Car if exists and authorized, null otherwise.
     */
    suspend fun getCarInfo(id: String): Car

    /**
     * Updates a car's details.
     *
     * @param car the Car and its updated info.
     * @return Result
     */
    suspend fun updateCar(car: Car): Result<Unit>?

    /**
     * Creates a car in the repository.
     *
     * @param car the Car to add to the repository.
     * @return Result
     */
    suspend fun createCar(car: CarDetails): Result<Unit>?
}
