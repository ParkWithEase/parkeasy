package io.github.parkwithease.parkeasy.data.local

import io.ktor.http.Cookie

interface AuthRepository {
    /**
     * Gets the session Cookie.
     *
     * @return the session Cookie.
     */
    suspend fun getSession(): Cookie?

    /**
     * Gets the session status.
     *
     * @return true if user logged in, false otherwise.
     */
    suspend fun getStatus(): Boolean

    /**
     * Sets the session Cookie and status to logged in.
     *
     * @param cookie the session Cookie to set.
     */
    suspend fun set(cookie: Cookie)

    /** Clears the session Cookie and sets session status to logged out. */
    suspend fun reset()
}
