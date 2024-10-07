package io.github.parkwithease.parkeasy.data.local

import io.ktor.http.Cookie

interface AuthRepository {
    suspend fun getSession(): Cookie?

    suspend fun getStatus(): Boolean

    suspend fun set(cookie: Cookie)

    suspend fun reset()
}
