package io.github.parkwithease.parkeasy.data.remote

import io.github.parkwithease.parkeasy.data.common.mapAPIError
import io.github.parkwithease.parkeasy.data.local.AuthRepository
import io.github.parkwithease.parkeasy.di.IoDispatcher
import io.github.parkwithease.parkeasy.model.Car
import io.github.parkwithease.parkeasy.model.CarDetails
import io.github.parkwithease.parkeasy.model.LoggedOutException
import io.ktor.client.HttpClient
import io.ktor.client.call.body
import io.ktor.client.request.cookie
import io.ktor.client.request.get
import io.ktor.client.request.post
import io.ktor.client.request.put
import io.ktor.client.request.setBody
import io.ktor.client.statement.HttpResponse
import io.ktor.http.ContentType
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
    override suspend fun getCars(): Result<List<Car>> =
        authRepo.sessionFlow
            .firstOrNull()
            .runCatching {
                if (this == null) throw LoggedOutException()
                withContext(ioDispatcher) { client.get("/cars") { cookie(name, value) } }
            }
            .mapAPIError()
            .let { result ->
                result.mapCatching {
                    if (it is HttpResponse && result.isSuccess) it.body<List<Car>>()
                    else emptyList()
                }
            }

    override suspend fun deleteCar(id: String): Result<Unit> {
        TODO("Not yet implemented")
    }

    override suspend fun getCarInfo(id: String): Result<Car> {
        TODO("Not yet implemented")
    }

    override suspend fun updateCar(car: Car): Result<Unit> =
        authRepo.sessionFlow
            .firstOrNull()
            .runCatching {
                if (this == null) throw LoggedOutException()
                withContext(ioDispatcher) {
                    client.put("/cars/" + car.id) {
                        contentType(ContentType.Application.Json)
                        setBody(car.details)
                        cookie(name = name, value = value)
                    }
                }
            }
            .mapAPIError()
            .map {}

    override suspend fun createCar(car: CarDetails): Result<Unit> =
        authRepo.sessionFlow
            .firstOrNull()
            .runCatching {
                if (this == null) throw LoggedOutException()
                withContext(ioDispatcher) {
                    client.post("/cars") {
                        contentType(ContentType.Application.Json)
                        setBody(car)
                        cookie(name = name, value = value)
                    }
                }
            }
            .mapAPIError()
            .map {}
}
