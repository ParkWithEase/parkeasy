package io.github.parkwithease.parkeasy.di

import dagger.Binds
import dagger.Module
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent
import io.github.parkwithease.parkeasy.data.local.AuthRepository
import io.github.parkwithease.parkeasy.data.local.AuthRepositoryImpl

@Module
@InstallIn(SingletonComponent::class)
interface AuthModule {
    @Binds fun bindAuthRepository(authRepositoryImpl: AuthRepositoryImpl): AuthRepository
}
