package io.github.parkwithease.parkeasy.data

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

class AuthStoreImpl(private val dataStore: DataStore<Preferences>) : AuthStore {
    private val loggedIn = booleanPreferencesKey("loggedIn")
    private val session = stringPreferencesKey("session")

    override suspend fun get(): Cookie {
        val sessionFlow: Flow<String> =
            dataStore.data.map { preferences -> preferences[session] ?: "" }
        return parseServerSetCookieHeader(sessionFlow.first())
    }

    override suspend fun set(cookie: Cookie) {
        dataStore.edit { settings ->
            settings[loggedIn] = true
            settings[session] = renderSetCookieHeader(cookie)
        }
    }

    override suspend fun reset() {
        dataStore.edit { settings ->
            settings[loggedIn] = false
            settings[session] = ""
        }
    }
}
