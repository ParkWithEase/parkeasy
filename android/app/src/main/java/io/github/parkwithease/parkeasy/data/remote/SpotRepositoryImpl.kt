package io.github.parkwithease.parkeasy.data.remote

import io.github.parkwithease.parkeasy.data.local.AuthRepository
import io.github.parkwithease.parkeasy.di.IoDispatcher
import io.github.parkwithease.parkeasy.model.Spot
import io.ktor.client.HttpClient
import io.ktor.client.call.body
import io.ktor.client.request.cookie
import io.ktor.client.request.get
import io.ktor.http.HttpStatusCode
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
    override suspend fun getSpots(): List<Spot> {
        val authCookie = authRepo.sessionFlow.firstOrNull()
        var spots = emptyList<Spot>()
        if (authCookie != null) {
            val response =
                withContext(ioDispatcher) {
                    client.get("/user/spots") { cookie(authCookie.name, authCookie.value) }
                }
            if (response.status == HttpStatusCode.OK) {
                spots = response.body()
            }
        }
        return spots
    }
}
