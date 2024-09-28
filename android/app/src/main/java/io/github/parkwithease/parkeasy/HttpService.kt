package io.github.parkwithease.parkeasy

import io.ktor.client.HttpClient
import io.ktor.client.engine.android.Android
import io.ktor.client.request.post
import io.ktor.client.request.setBody
import io.ktor.http.ContentType
import io.ktor.http.contentType

object HttpService {
    private val client: HttpClient by lazy {
        HttpClient(Android)
    }

    suspend fun login() {
        val response = client.post("") {
            contentType(ContentType.Application.Json)
            setBody("")
        }
    }
}
