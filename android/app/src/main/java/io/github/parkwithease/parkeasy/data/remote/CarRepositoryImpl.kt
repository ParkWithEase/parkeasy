package io.github.parkwithease.parkeasy.data.remote

import io.github.parkwithease.parkeasy.data.local.AuthRepository
import io.github.parkwithease.parkeasy.di.IoDispatcher
import io.github.parkwithease.parkeasy.model.Car
import io.github.parkwithease.parkeasy.model.CarDetails
import io.ktor.client.HttpClient
import io.ktor.client.call.body
import io.ktor.client.request.cookie
import io.ktor.client.request.get
import io.ktor.client.request.post
import io.ktor.client.request.setBody
import io.ktor.http.ContentType
import io.ktor.http.HttpStatusCode
import io.ktor.http.contentType
import javax.inject.Inject
import kotlinx.coroutines.CoroutineDispatcher
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.flow.firstOrNull
import kotlinx.coroutines.withContext

class CarRepositoryImpl
@Inject
constructor(
    private val client: HttpClient,
    private val authRepo: AuthRepository,
    @IoDispatcher private val ioDispatcher: CoroutineDispatcher = Dispatchers.IO,
) : CarRepository {
    override suspend fun getCars(): List<Car> {
        val authCookie = authRepo.sessionFlow.firstOrNull()
        var cars = emptyList<Car>()
        if (authCookie != null) {
            val response =
                withContext(ioDispatcher) {
                    client.get("/cars") { cookie(authCookie.name, authCookie.value) }
                }
            if (response.status == HttpStatusCode.OK) {
                cars = response.body()
            }
        }
        return cars
    }

    override suspend fun deleteCar(id: String): Boolean {
        TODO("Not yet implemented")
    }

    override suspend fun getCarInfo(id: String): Car {
        TODO("Not yet implemented")
    }

    override suspend fun updateCar(car: Car): Car {
        TODO("Not yet implemented")
    }

    override suspend fun createCar(car: CarDetails): Result<Unit> = runCatching {
        val authCookie = authRepo.sessionFlow.firstOrNull()
        if (authCookie != null) {
            withContext(ioDispatcher) {
                client.post("/cars") {
                    contentType(ContentType.Application.Json)
                    setBody(car)
                    cookie(authCookie.name, authCookie.value)
                }
            }
        }
    }
}
