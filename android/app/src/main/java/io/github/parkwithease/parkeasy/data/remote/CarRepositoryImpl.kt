package io.github.parkwithease.parkeasy.data.remote

import io.github.parkwithease.parkeasy.data.local.AuthRepository
import io.github.parkwithease.parkeasy.di.IoDispatcher
import io.github.parkwithease.parkeasy.model.Car
import io.github.parkwithease.parkeasy.model.CarDetails
import io.github.parkwithease.parkeasy.model.ErrorModel
import io.ktor.client.HttpClient
import io.ktor.client.call.body
import io.ktor.client.request.cookie
import io.ktor.client.request.get
import io.ktor.client.request.post
import io.ktor.client.request.put
import io.ktor.client.request.setBody
import io.ktor.client.statement.HttpResponse
import io.ktor.http.ContentType
import io.ktor.http.HttpStatusCode
import io.ktor.http.contentType
import io.ktor.http.isSuccess
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

    override suspend fun updateCar(car: Car): Result<Unit>? =
        authRepo.sessionFlow
            .firstOrNull()
            ?.runCatching {
                withContext(ioDispatcher) {
                    client.put("/cars/" + car.id) {
                        contentType(ContentType.Application.Json)
                        setBody(car.details)
                        cookie(name = name, value = value)
                    }
                }
            }
            ?.mapAPIError()
            ?.map {}

    override suspend fun createCar(car: CarDetails): Result<Unit>? =
        authRepo.sessionFlow
            .firstOrNull()
            ?.runCatching {
                withContext(ioDispatcher) {
                    client.post("/cars") {
                        contentType(ContentType.Application.Json)
                        setBody(car)
                        cookie(name = name, value = value)
                    }
                }
            }
            ?.mapAPIError()
            ?.map {}

    // Convert API error into a failing Result
    private suspend fun Result<HttpResponse>.mapAPIError(): Result<HttpResponse> = mapCatching {
        if (!it.status.isSuccess()) throw APIException(it.body<ErrorModel>())
        it
    }
}
