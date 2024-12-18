package io.github.parkwithease.parkeasy.data.remote

import io.github.parkwithease.parkeasy.data.common.mapAPIError
import io.github.parkwithease.parkeasy.data.local.AuthRepository
import io.github.parkwithease.parkeasy.di.IoDispatcher
import io.github.parkwithease.parkeasy.model.Booking
import io.github.parkwithease.parkeasy.model.BookingHistory
import io.github.parkwithease.parkeasy.model.LoggedOutException
import io.github.parkwithease.parkeasy.model.Spot
import io.github.parkwithease.parkeasy.model.TimeUnit
import io.ktor.client.HttpClient
import io.ktor.client.call.body
import io.ktor.client.request.cookie
import io.ktor.client.request.get
import io.ktor.client.request.parameter
import io.ktor.client.request.post
import io.ktor.client.request.setBody
import io.ktor.http.ContentType
import io.ktor.http.contentType
import javax.inject.Inject
import kotlinx.coroutines.CoroutineDispatcher
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.flow.firstOrNull
import kotlinx.coroutines.withContext

class SpotRepositoryImpl
@Inject
constructor(
    private val client: HttpClient,
    private val authRepo: AuthRepository,
    @IoDispatcher private val ioDispatcher: CoroutineDispatcher = Dispatchers.IO,
) : SpotRepository {
    override suspend fun getSpots(): Result<List<Spot>> =
        authRepo.sessionFlow
            .firstOrNull()
            .runCatching {
                if (this == null) throw LoggedOutException()
                withContext(ioDispatcher) { client.get("/user/spots") { cookie(name, value) } }
            }
            .mapAPIError()
            .let { result ->
                result.mapCatching { if (result.isSuccess) it.body<List<Spot>>() else emptyList() }
            }

    override suspend fun createSpot(spot: Spot): Result<Unit> =
        authRepo.sessionFlow
            .firstOrNull()
            .runCatching {
                if (this == null) throw LoggedOutException()
                withContext(ioDispatcher) {
                    client.post("/spots") {
                        contentType(ContentType.Application.Json)
                        setBody(spot)
                        cookie(name = name, value = value)
                    }
                }
            }
            .mapAPIError()
            .map {}

    override suspend fun getSpotsAround(
        latitude: Double,
        longitude: Double,
        distance: Int,
    ): Result<List<Spot>> =
        authRepo.sessionFlow
            .firstOrNull()
            .runCatching {
                if (this == null) throw LoggedOutException()
                withContext(ioDispatcher) {
                    client.get("/spots") {
                        cookie(name, value)
                        parameter("latitude", latitude)
                        parameter("longitude", longitude)
                        parameter("distance", distance)
                    }
                }
            }
            .mapAPIError()
            .let { result ->
                result.mapCatching { if (result.isSuccess) it.body<List<Spot>>() else emptyList() }
            }

    override suspend fun getSpotAvailability(id: String): Result<List<TimeUnit>> =
        authRepo.sessionFlow
            .firstOrNull()
            .runCatching {
                if (this == null) throw LoggedOutException()
                withContext(ioDispatcher) {
                    client.get("/spots/${id}/availability") { cookie(name, value) }
                }
            }
            .mapAPIError()
            .let { result ->
                result.mapCatching { if (result.isSuccess) it.body() else emptyList() }
            }

    override suspend fun bookSpot(spotId: String, booking: Booking): Result<Unit> =
        authRepo.sessionFlow
            .firstOrNull()
            .runCatching {
                if (this == null) throw LoggedOutException()
                withContext(ioDispatcher) {
                    client.post("/spots/${spotId}/bookings") {
                        contentType(ContentType.Application.Json)
                        setBody(booking)
                        cookie(name = name, value = value)
                    }
                }
            }
            .mapAPIError()
            .map {}

    override suspend fun getBookings(): Result<List<BookingHistory>> =
        authRepo.sessionFlow
            .firstOrNull()
            .runCatching {
                if (this == null) throw LoggedOutException()
                withContext(ioDispatcher) { client.get("/user/bookings") { cookie(name, value) } }
            }
            .mapAPIError()
            .let { result ->
                result.mapCatching {
                    if (result.isSuccess) it.body<List<BookingHistory>>() else emptyList()
                }
            }

    override suspend fun getLeasings(): Result<List<BookingHistory>> =
        authRepo.sessionFlow
            .firstOrNull()
            .runCatching {
                if (this == null) throw LoggedOutException()
                withContext(ioDispatcher) { client.get("/user/leasings") { cookie(name, value) } }
            }
            .mapAPIError()
            .let { result ->
                result.mapCatching {
                    if (result.isSuccess) it.body<List<BookingHistory>>() else emptyList()
                }
            }
}
