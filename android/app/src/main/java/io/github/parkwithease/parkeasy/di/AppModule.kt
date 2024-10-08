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
import javax.inject.Qualifier
import javax.inject.Singleton
import kotlinx.coroutines.CoroutineDispatcher
import kotlinx.coroutines.Dispatchers

@Retention(AnnotationRetention.BINARY) @Qualifier annotation class IoDispatcher

@Module
@InstallIn(SingletonComponent::class)
object AppModule {
    private val Context.dataStore: DataStore<Preferences> by preferencesDataStore(name = "auth")

    @IoDispatcher
    @Provides
    fun providesIoDispatcher(
        dispatcher: CoroutineDispatcher = Dispatchers.IO
    ): CoroutineDispatcher = dispatcher

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
        return AuthRepositoryImpl(app.applicationContext.dataStore)
    }

    @Provides
    @Singleton
    fun provideUserRepository(client: HttpClient, authStore: AuthRepository): UserRepository {
        return UserRepositoryImpl(client, authStore)
    }
}
