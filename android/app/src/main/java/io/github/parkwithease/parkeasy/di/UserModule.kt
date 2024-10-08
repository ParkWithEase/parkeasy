package io.github.parkwithease.parkeasy.di

import dagger.Binds
import dagger.Module
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent
import io.github.parkwithease.parkeasy.data.remote.UserRepository
import io.github.parkwithease.parkeasy.data.remote.UserRepositoryImpl

@Module
@InstallIn(SingletonComponent::class)
interface UserModule {
    @Binds fun bindUserRepository(userRepositoryImpl: UserRepositoryImpl): UserRepository
}
