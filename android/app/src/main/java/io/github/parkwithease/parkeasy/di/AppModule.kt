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
import dagger.hilt.android.qualifiers.ApplicationContext
import dagger.hilt.components.SingletonComponent
import io.github.parkwithease.parkeasy.R
import io.github.parkwithease.parkeasy.data.local.AuthRepository
import io.github.parkwithease.parkeasy.data.local.AuthRepositoryImpl
import io.github.parkwithease.parkeasy.data.remote.UserRepository
import io.github.parkwithease.parkeasy.data.remote.UserRepositoryImpl
import io.ktor.client.HttpClient
import io.ktor.client.engine.okhttp.OkHttp
import io.ktor.client.plugins.contentnegotiation.ContentNegotiation
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
    fun provideAuthStore(@ApplicationContext context: Context): AuthRepository {
        return AuthRepositoryImpl(context.dataStore)
    }

    @Provides
    @Singleton
    fun provideUserRepository(client: HttpClient, authStore: AuthRepository): UserRepository {
        return UserRepositoryImpl(client, authStore)
    }
}
