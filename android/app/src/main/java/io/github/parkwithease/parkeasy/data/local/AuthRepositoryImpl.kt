package io.github.parkwithease.parkeasy.data.local

import androidx.datastore.core.DataStore
import androidx.datastore.preferences.core.Preferences
import androidx.datastore.preferences.core.booleanPreferencesKey
import androidx.datastore.preferences.core.edit
import androidx.datastore.preferences.core.stringPreferencesKey
import io.ktor.http.Cookie
import io.ktor.http.parseServerSetCookieHeader
import io.ktor.http.renderSetCookieHeader
import javax.inject.Inject
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map

class AuthRepositoryImpl @Inject constructor(private val dataStore: DataStore<Preferences>) :
    AuthRepository {
    private val session = stringPreferencesKey("session")
    private val status = booleanPreferencesKey("status")

    override val sessionFlow: Flow<Cookie> =
        dataStore.data.map { preferences ->
            val session = preferences[session]
            if (session != null) parseServerSetCookieHeader(session) else Cookie("session", "")
        }

    override val statusFlow: Flow<Boolean> =
        dataStore.data.map { preferences -> preferences[status] ?: false }

    override suspend fun set(cookie: Cookie) {
        dataStore.edit { settings ->
            settings[session] = renderSetCookieHeader(cookie)
            settings[status] = true
        }
    }

    override suspend fun reset() {
        dataStore.edit { settings ->
            settings.remove(session)
            settings[status] = false
        }
    }
}
