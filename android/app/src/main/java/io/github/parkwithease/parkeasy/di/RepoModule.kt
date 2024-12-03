package io.github.parkwithease.parkeasy.di

import dagger.Binds
import dagger.Module
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent
import io.github.parkwithease.parkeasy.data.local.AuthRepository
import io.github.parkwithease.parkeasy.data.local.AuthRepositoryImpl
import io.github.parkwithease.parkeasy.data.remote.CarRepository
import io.github.parkwithease.parkeasy.data.remote.CarRepositoryImpl
import io.github.parkwithease.parkeasy.data.remote.SpotRepository
import io.github.parkwithease.parkeasy.data.remote.SpotRepositoryImpl
import io.github.parkwithease.parkeasy.data.remote.UserRepository
import io.github.parkwithease.parkeasy.data.remote.UserRepositoryImpl

@Suppress("unused") // IntelliJ does not recognize that these are used
@Module
@InstallIn(SingletonComponent::class)
interface RepoModule {
    @Binds fun bindAuthRepository(authRepositoryImpl: AuthRepositoryImpl): AuthRepository

    @Binds fun bindUserRepository(userRepositoryImpl: UserRepositoryImpl): UserRepository

    @Binds fun bindCarRepository(carRepositoryImpl: CarRepositoryImpl): CarRepository

    @Binds fun bindSpotRepository(spotRepositoryImpl: SpotRepositoryImpl): SpotRepository
}
