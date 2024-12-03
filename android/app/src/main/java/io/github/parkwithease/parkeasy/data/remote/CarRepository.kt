package io.github.parkwithease.parkeasy.data.remote

import io.github.parkwithease.parkeasy.model.Car
import io.github.parkwithease.parkeasy.model.CarDetails

interface CarRepository {
    /**
     * Gets logged in user's cars.
     *
     * @return a Result List of Cars on success, failing Result otherwise.
     */
    suspend fun getCars(): Result<List<Car>>

    /**
     * Deletes a car in the database by id.
     *
     * @return whether deleting the car was successful or not.
     */
    suspend fun deleteCar(id: String): Result<Boolean>

    /**
     * Gets a car's details by id.
     *
     * @return Result with the Car if exists and authorized, failing Result otherwise.
     */
    suspend fun getCarInfo(id: String): Result<Car>

    /**
     * Updates a car's details.
     *
     * @param car the Car and its updated info.
     * @return whether the update was successful or not.
     */
    suspend fun updateCar(car: Car): Result<Unit>

    /**
     * Creates a car in the repository.
     *
     * @param car the Car to add to the repository.
     * @return whether the creation was successful or not.
     */
    suspend fun createCar(car: CarDetails): Result<Unit>
}
