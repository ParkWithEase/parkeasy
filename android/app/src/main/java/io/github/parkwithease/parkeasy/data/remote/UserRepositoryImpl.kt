package io.github.parkwithease.parkeasy.data.remote

import android.util.Log
import io.github.parkwithease.parkeasy.data.local.AuthRepository
import io.github.parkwithease.parkeasy.di.IoDispatcher
import io.github.parkwithease.parkeasy.model.LoginCredentials
import io.github.parkwithease.parkeasy.model.RegistrationCredentials
import io.github.parkwithease.parkeasy.model.ResetCredentials
import io.ktor.client.HttpClient
import io.ktor.client.request.cookie
import io.ktor.client.request.delete
import io.ktor.client.request.post
import io.ktor.client.request.setBody
import io.ktor.http.ContentType
import io.ktor.http.Cookie
import io.ktor.http.HttpStatusCode
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
    override suspend fun login(credentials: LoginCredentials): Boolean {
        var success = false
        val sessionCookie: Cookie?
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
        Log.d("HTTP", sessionCookie.toString())
        return success
    }

    override suspend fun register(credentials: RegistrationCredentials): Boolean {
        var success = false
        val sessionCookie: Cookie?
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
        Log.d("HTTP", sessionCookie.toString())
        return success
    }

    override suspend fun logout() {
        val authCookie = authRepo.sessionFlow.firstOrNull()
        if (authCookie != null) {
            val response =
                withContext(ioDispatcher) {
                    client.delete("/auth") {
                        contentType(ContentType.Application.Json)
                        cookie(authCookie.name, authCookie.value)
                    }
                }
            authRepo.reset()
            Log.d("HTTP", response.toString())
        }
    }

    override suspend fun requestReset(credentials: ResetCredentials): Boolean {
        val response =
            withContext(ioDispatcher) {
                client.post("/auth/password:forgot") {
                    contentType(ContentType.Application.Json)
                    setBody(credentials)
                }
            }
        return response.status == HttpStatusCode.Accepted
    }
}
