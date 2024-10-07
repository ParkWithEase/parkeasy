package io.github.parkwithease.parkeasy.data

import android.util.Log
import io.github.parkwithease.parkeasy.model.LoginCredentials
import io.github.parkwithease.parkeasy.model.RegistrationCredentials
import io.ktor.client.HttpClient
import io.ktor.client.request.post
import io.ktor.client.request.setBody
import io.ktor.http.ContentType
import io.ktor.http.Cookie
import io.ktor.http.contentType
import io.ktor.http.setCookie

class UserRepositoryImpl(private val client: HttpClient, private val authStore: AuthStore) :
    UserRepository {
    override suspend fun login(credentials: LoginCredentials) {
        var sessionCookie: Cookie? = null
        val response =
            client.post("/auth") {
                contentType(ContentType.Application.Json)
                setBody(credentials)
            }
        if (response.setCookie().size == 1) {
            sessionCookie = response.setCookie()[0]
            authStore.set(sessionCookie)
        }
        Log.d("HTTP", sessionCookie.toString())
    }

    override suspend fun register(credentials: RegistrationCredentials) {
        var sessionCookie: Cookie? = null
        val response =
            client.post("/user") {
                contentType(ContentType.Application.Json)
                setBody(credentials)
            }
        if (response.setCookie().size == 1) {
            sessionCookie = response.setCookie()[0]
            authStore.set(sessionCookie)
        }
        Log.d("HTTP", sessionCookie.toString())
    }
}
