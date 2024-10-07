package io.github.parkwithease.parkeasy.di

import android.app.Application
import android.content.Context
import android.util.Log
import androidx.datastore.core.DataStore
import androidx.datastore.preferences.core.Preferences
import androidx.datastore.preferences.preferencesDataStore
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent
import io.github.parkwithease.parkeasy.R
import io.github.parkwithease.parkeasy.data.AuthStore
import io.github.parkwithease.parkeasy.data.AuthStoreImpl
import io.github.parkwithease.parkeasy.data.UserRepository
import io.github.parkwithease.parkeasy.data.UserRepositoryImpl
import io.ktor.client.HttpClient
import io.ktor.client.engine.okhttp.OkHttp
import io.ktor.client.plugins.contentnegotiation.ContentNegotiation
import io.ktor.client.plugins.cookies.HttpCookies
import io.ktor.client.plugins.defaultRequest
import io.ktor.client.plugins.logging.LogLevel
import io.ktor.client.plugins.logging.Logger
import io.ktor.client.plugins.logging.Logging
import io.ktor.serialization.kotlinx.json.json
import javax.inject.Singleton

@Module
@InstallIn(SingletonComponent::class)
object AppModule {
    private val Context.dataStore: DataStore<Preferences> by preferencesDataStore(name = "auth")

    @Provides
    @Singleton
    fun provideHttpClient(app: Application): HttpClient =
        HttpClient(OkHttp) {
            defaultRequest { url(app.getString(R.string.api_host)) }
            install(ContentNegotiation) { json() }
            install(HttpCookies) {}
            install(Logging) {
                logger =
                    object : Logger {
                        override fun log(message: String) {
                            Log.d("HTTP", message)
                        }
                    }
                level = LogLevel.ALL
            }
        }

    @Provides
    @Singleton
    fun provideAuthStore(app: Application): AuthStore {
        return AuthStoreImpl(app.applicationContext.dataStore)
    }

    @Provides
    @Singleton
    fun provideUserRepository(client: HttpClient, authStore: AuthStore): UserRepository {
        return UserRepositoryImpl(client, authStore)
    }
}
