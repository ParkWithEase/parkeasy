package io.github.parkwithease.parkeasy.data.local

import androidx.datastore.core.DataStore
import androidx.datastore.preferences.core.Preferences
import androidx.datastore.preferences.core.booleanPreferencesKey
import androidx.datastore.preferences.core.edit
import androidx.datastore.preferences.core.stringPreferencesKey
import io.ktor.http.Cookie
import io.ktor.http.parseServerSetCookieHeader
import io.ktor.http.renderSetCookieHeader
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.flow.map

class AuthRepositoryImpl(private val dataStore: DataStore<Preferences>) : AuthRepository {
    private val session = stringPreferencesKey("session")
    private val status = booleanPreferencesKey("status")

    override suspend fun getSession(): Cookie {
        val sessionFlow: Flow<String> =
            dataStore.data.map { preferences -> preferences[session] ?: "" }
        return if (sessionFlow.first() != "") parseServerSetCookieHeader(sessionFlow.first())
        else Cookie("session", "")
    }

    override suspend fun getStatus(): Boolean {
        val sessionFlow: Flow<Boolean> =
            dataStore.data.map { preferences -> preferences[status] ?: false }
        return sessionFlow.first()
    }

    override suspend fun set(cookie: Cookie) {
        dataStore.edit { settings ->
            settings[session] = renderSetCookieHeader(cookie)
            settings[status] = true
        }
    }

    override suspend fun reset() {
        dataStore.edit { settings ->
            settings[session] = ""
            settings[status] = false
        }
    }
}
