package io.github.parkwithease.parkeasy.data

import android.util.Log
import io.github.parkwithease.parkeasy.model.Credentials
import io.ktor.client.HttpClient
import io.ktor.client.request.post
import io.ktor.client.request.setBody
import io.ktor.http.ContentType
import io.ktor.http.Cookie
import io.ktor.http.contentType
import io.ktor.http.setCookie

class ParkEasyRepositoryImpl(private val client: HttpClient) : ParkEasyRepository {
    override suspend fun login(credentials: Credentials) {
        var sessionCookie: Cookie? = null
        val response =
            client.post("/auth") {
                contentType(ContentType.Application.Json)
                setBody(credentials)
            }
        if (response.setCookie().size == 1) {
            sessionCookie = response.setCookie()[0]
        }
        Log.d("HTTP", sessionCookie.toString())
    }
}
