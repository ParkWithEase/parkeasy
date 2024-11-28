package io.github.parkwithease.parkeasy.data.remote

import io.github.parkwithease.parkeasy.data.common.mapAPIError
import io.github.parkwithease.parkeasy.data.local.AuthRepository
import io.github.parkwithease.parkeasy.di.IoDispatcher
import io.github.parkwithease.parkeasy.model.LoggedOutException
import io.github.parkwithease.parkeasy.model.LoginCredentials
import io.github.parkwithease.parkeasy.model.Profile
import io.github.parkwithease.parkeasy.model.RegistrationCredentials
import io.github.parkwithease.parkeasy.model.ResetCredentials
import io.ktor.client.HttpClient
import io.ktor.client.call.body
import io.ktor.client.request.cookie
import io.ktor.client.request.delete
import io.ktor.client.request.get
import io.ktor.client.request.post
import io.ktor.client.request.setBody
import io.ktor.client.statement.HttpResponse
import io.ktor.http.ContentType
import io.ktor.http.contentType
import io.ktor.http.setCookie
import javax.inject.Inject
import kotlinx.coroutines.CoroutineDispatcher
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.flow.firstOrNull
import kotlinx.coroutines.withContext

class UserRepositoryImpl
@Inject
constructor(
    private val client: HttpClient,
    private val authRepo: AuthRepository,
    @IoDispatcher private val ioDispatcher: CoroutineDispatcher = Dispatchers.IO,
) : UserRepository {
    override suspend fun login(credentials: LoginCredentials): Result<Unit> =
        runCatching {
                withContext(ioDispatcher) {
                    client.post("/auth") {
                        contentType(ContentType.Application.Json)
                        setBody(credentials)
                    }
                }
            }
            .mapAPIError()
            .updateAuthCookie()
            .map {}

    override suspend fun register(credentials: RegistrationCredentials): Result<Unit> =
        runCatching {
                withContext(ioDispatcher) {
                    client.post("/user") {
                        contentType(ContentType.Application.Json)
                        setBody(credentials)
                    }
                }
            }
            .mapAPIError()
            .updateAuthCookie()
            .map {}

    override suspend fun logout(): Boolean =
        authRepo.sessionFlow
            .firstOrNull()
            ?.runCatching {
                withContext(ioDispatcher) {
                    client.delete("/auth") { cookie(name = name, value = value) }
                }
            }
            ?.mapAPIError()
            ?.updateAuthCookie()
            .let { it == null || it.onFailure { e -> throw e }.isSuccess }

    override suspend fun requestReset(credentials: ResetCredentials): Result<Unit> =
        runCatching {
                withContext(ioDispatcher) {
                    client.post("/auth/password:forgot") {
                        contentType(ContentType.Application.Json)
                        setBody(credentials)
                    }
                }
            }
            .mapAPIError()
            .map {}

    override suspend fun getUser(): Result<Profile> =
        authRepo.sessionFlow
            .firstOrNull()
            .runCatching {
                if (this == null) throw LoggedOutException()
                withContext(ioDispatcher) {
                    client.get("/user") { cookie(name = name, value = value) }
                }
            }
            .mapAPIError()
            .let { result ->
                if (result.isSuccess) result.mapCatching { it.body<Profile>() }
                else result.map { Profile() }
            }

    // Update authentication status based on the response assuming that the request alters
    // authentication status
    private suspend fun Result<HttpResponse>.updateAuthCookie(): Result<HttpResponse> =
        onSuccess { response ->
                response
                    .setCookie()
                    .find { it.name == "session" && it.value != "" }
                    .let { if (it != null) authRepo.set(it) else authRepo.reset() }
            }
            .onFailure { if (it is APIException) authRepo.reset() }
}
