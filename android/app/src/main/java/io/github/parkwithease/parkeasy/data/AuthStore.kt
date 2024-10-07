package io.github.parkwithease.parkeasy.data

import io.ktor.http.Cookie

interface AuthStore {
    suspend fun get(): Cookie?

    suspend fun set(cookie: Cookie)

    suspend fun reset()
}
