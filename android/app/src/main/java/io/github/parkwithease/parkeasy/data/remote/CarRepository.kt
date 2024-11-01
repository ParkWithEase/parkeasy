package io.github.parkwithease.parkeasy.data.remote

import io.github.parkwithease.parkeasy.model.Car

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
     * @return the updated Car on success, null otherwise.
     */
    suspend fun updateCar(car: Car): Car

    /**
     * Creates a car in the repository.
     *
     * @param car the Car to add to the repository.
     * @return the created Car on success, null otherwise.
     */
    suspend fun createCar(car: Car): Car
}
