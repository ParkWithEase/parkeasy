package io.github.parkwithease.parkeasy.data.remote

import android.util.Log
import io.github.parkwithease.parkeasy.data.local.AuthRepository
import io.github.parkwithease.parkeasy.model.LoginCredentials
import io.github.parkwithease.parkeasy.model.RegistrationCredentials
import io.ktor.client.HttpClient
import io.ktor.client.request.cookie
import io.ktor.client.request.delete
import io.ktor.client.request.post
import io.ktor.client.request.setBody
import io.ktor.http.ContentType
import io.ktor.http.Cookie
import io.ktor.http.contentType
import io.ktor.http.setCookie
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext

class UserRepositoryImpl(private val client: HttpClient, private val authStore: AuthRepository) :
    UserRepository {
    override suspend fun login(credentials: LoginCredentials): Boolean {
        return withContext(Dispatchers.IO) {
            var success = false
            var sessionCookie: Cookie? = null
            val response =
                client.post("/auth") {
                    contentType(ContentType.Application.Json)
                    setBody(credentials)
                }
            if (response.setCookie().size == 1) {
                sessionCookie = response.setCookie()[0]
                authStore.set(sessionCookie)
                success = true
            }
            Log.d("HTTP", sessionCookie.toString())
            return@withContext success
        }
    }

    override suspend fun register(credentials: RegistrationCredentials): Boolean {
        return withContext(Dispatchers.IO) {
            var success = false
            var sessionCookie: Cookie? = null
            val response =
                client.post("/user") {
                    contentType(ContentType.Application.Json)
                    setBody(credentials)
                }
            if (response.setCookie().size == 1) {
                sessionCookie = response.setCookie()[0]
                authStore.set(sessionCookie)
                success = true
            }
            Log.d("HTTP", sessionCookie.toString())
            return@withContext success
        }
    }

    override suspend fun logout() {
        withContext(Dispatchers.IO) {
            val authCookie = authStore.getSession()
            if (authCookie != null) {
                val response =
                    client.delete("/auth") {
                        contentType(ContentType.Application.Json)
                        cookie(authCookie.name, authCookie.value)
                    }
                authStore.reset()
                Log.d("HTTP", response.toString())
            }
        }
    }
}
