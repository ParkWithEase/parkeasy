package io.github.parkwithease.parkeasy.di

import dagger.Binds
import dagger.Module
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent
import io.github.parkwithease.parkeasy.data.local.AuthRepository
import io.github.parkwithease.parkeasy.data.local.AuthRepositoryImpl
import io.github.parkwithease.parkeasy.data.remote.UserRepository
import io.github.parkwithease.parkeasy.data.remote.UserRepositoryImpl

@Module
@InstallIn(SingletonComponent::class)
@Suppress("Unused")
interface RepoModule {
    @Binds fun bindAuthRepository(authRepositoryImpl: AuthRepositoryImpl): AuthRepository
    @Binds fun bindUserRepository(userRepositoryImpl: UserRepositoryImpl): UserRepository
}
