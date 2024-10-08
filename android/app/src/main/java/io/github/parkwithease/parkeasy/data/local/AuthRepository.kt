package io.github.parkwithease.parkeasy.data.local

import io.ktor.http.Cookie
import kotlinx.coroutines.flow.Flow

interface AuthRepository {
    /** Flow of the session Cookie. */
    val sessionFlow: Flow<Cookie>

    /** Flow of the session status (true if logged in, false otherwise). */
    val statusFlow: Flow<Boolean>

    /**
     * Sets the session Cookie and status to logged in.
     *
     * @param cookie the session Cookie to set.
     */
    suspend fun set(cookie: Cookie)

    /** Clears the session Cookie and sets session status to logged out. */
    suspend fun reset()
}
