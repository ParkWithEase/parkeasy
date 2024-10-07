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
import kotlinx.coroutines.CoroutineDispatcher
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext

class UserRepositoryImpl(
    private val client: HttpClient,
    private val authRepo: AuthRepository,
    @IoDispatcher private val ioDispatcher: CoroutineDispatcher = Dispatchers.IO,
) : UserRepository {
    override suspend fun login(credentials: LoginCredentials): Boolean {
        return withContext(ioDispatcher) {
            var success = false
            val sessionCookie: Cookie?
            val response =
                client.post("/auth") {
                    contentType(ContentType.Application.Json)
                    setBody(credentials)
                }
            sessionCookie = response.setCookie().firstOrNull()
            if (sessionCookie != null) {
                authRepo.set(sessionCookie)
                success = true
            }
            Log.d("HTTP", sessionCookie.toString())
            return@withContext success
        }
    }

    override suspend fun register(credentials: RegistrationCredentials): Boolean {
        return withContext(ioDispatcher) {
            var success = false
            val sessionCookie: Cookie?
            val response =
                client.post("/user") {
                    contentType(ContentType.Application.Json)
                    setBody(credentials)
                }
            sessionCookie = response.setCookie().firstOrNull()
            if (sessionCookie != null) {
                authRepo.set(sessionCookie)
                success = true
            }
            Log.d("HTTP", sessionCookie.toString())
            return@withContext success
        }
    }

    override suspend fun logout() {
        withContext(ioDispatcher) {
            val authCookie = authRepo.getSession()
            if (authCookie != null) {
                val response =
                    client.delete("/auth") {
                        contentType(ContentType.Application.Json)
                        cookie(authCookie.name, authCookie.value)
                    }
                authRepo.reset()
                Log.d("HTTP", response.toString())
            }
        }
    }

    override suspend fun reset(credentials: ResetCredentials): Boolean {
        return withContext(ioDispatcher) {
            val response =
                client.post("/auth/password:forgot") {
                    contentType(ContentType.Application.Json)
                    setBody(credentials)
                }
            return@withContext response.status == HttpStatusCode.OK
        }
    }
}
