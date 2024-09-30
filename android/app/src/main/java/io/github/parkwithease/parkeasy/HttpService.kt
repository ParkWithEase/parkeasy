package io.github.parkwithease.parkeasy

import android.util.Log
import io.ktor.client.HttpClient
import io.ktor.client.engine.android.Android
import io.ktor.client.plugins.contentnegotiation.ContentNegotiation
import io.ktor.client.plugins.cookies.HttpCookies
import io.ktor.client.plugins.logging.LogLevel
import io.ktor.client.plugins.logging.Logger
import io.ktor.client.plugins.logging.Logging
import io.ktor.client.request.post
import io.ktor.client.request.setBody
import io.ktor.http.ContentType
import io.ktor.http.Cookie
import io.ktor.http.contentType
import io.ktor.http.setCookie
import io.ktor.serialization.kotlinx.json.json

object HttpService {
    const val API_HOST = "http://10.0.2.2:8080"
    var SESSION_COOKIE: Cookie? = null

    private val client: HttpClient by lazy {
        HttpClient(Android) {
            install(ContentNegotiation) { json() }
            install(HttpCookies) {}

            install(Logging) {
                logger =
                    object : Logger {
                        override fun log(message: String) {
                            Log.d("HTTP call", message)
                        }
                    }
                level = LogLevel.ALL
            }
        }
    }

    suspend fun login(credentials: Credentials) {
        val response =
            client.post("$API_HOST/auth") {
                contentType(ContentType.Application.Json)
                setBody(credentials)
            }
        if (response.setCookie().size == 1) {
            SESSION_COOKIE = response.setCookie()[0]
        }
        Log.d("HTTP call", SESSION_COOKIE.toString())
    }
}
