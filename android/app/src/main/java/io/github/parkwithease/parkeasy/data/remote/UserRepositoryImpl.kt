package io.github.parkwithease.parkeasy.data.remote

import io.github.parkwithease.parkeasy.data.local.AuthRepository
import io.github.parkwithease.parkeasy.di.IoDispatcher
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
import io.ktor.http.ContentType
import io.ktor.http.Cookie
import io.ktor.http.contentType
import io.ktor.http.isSuccess
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
    override suspend fun login(credentials: LoginCredentials): Boolean {
        val sessionCookie: Cookie?
        var success = false
        val response =
            withContext(ioDispatcher) {
                client.post("/auth") {
                    contentType(ContentType.Application.Json)
                    setBody(credentials)
                }
            }
        sessionCookie = response.setCookie().firstOrNull()
        if (sessionCookie != null) {
            authRepo.set(sessionCookie)
            success = true
        }
        return success
    }

    override suspend fun register(credentials: RegistrationCredentials): Boolean {
        val sessionCookie: Cookie?
        var success = false
        val response =
            withContext(ioDispatcher) {
                client.post("/user") {
                    contentType(ContentType.Application.Json)
                    setBody(credentials)
                }
            }
        sessionCookie = response.setCookie().firstOrNull()
        if (sessionCookie != null) {
            authRepo.set(sessionCookie)
            success = true
        }
        return success
    }

    override suspend fun logout(): Boolean {
        val authCookie = authRepo.sessionFlow.firstOrNull()
        var success = false
        if (authCookie != null) {
            val response =
                withContext(ioDispatcher) {
                    client.delete("/auth") { cookie(authCookie.name, authCookie.value) }
                }
            if (response.status == HttpStatusCode.NoContent) {
                authRepo.reset()
                success = true
            }
        }
        return success
    }

    override suspend fun requestReset(credentials: ResetCredentials): Boolean {
        val response =
            withContext(ioDispatcher) {
                client.post("/auth/password:forgot") { setBody(credentials) }
            }
        return response.status.isSuccess()
    }

    override suspend fun getUser(): Profile? {
        val authCookie = authRepo.sessionFlow.firstOrNull()
        var profile: Profile? = null
        if (authCookie != null) {
            val response =
                withContext(ioDispatcher) {
                    client.get("/user") { cookie(authCookie.name, authCookie.value) }
                }
            if (response.status == HttpStatusCode.OK) {
                profile = response.body()
            }
        }
        return profile
    }
}
