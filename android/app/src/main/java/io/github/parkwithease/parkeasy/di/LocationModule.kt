package io.github.parkwithease.parkeasy.di

import dagger.Binds
import dagger.Module
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent
import io.github.parkwithease.parkeasy.domain.GetLocationUseCase
import io.github.parkwithease.parkeasy.domain.GetLocationUseCaseImpl
import javax.inject.Singleton

@Module
@InstallIn(SingletonComponent::class)
interface LocationModule {
    @Binds
    @Singleton
    fun bindGetLocationUseCase(getLocationUseCaseImpl: GetLocationUseCaseImpl): GetLocationUseCase
}
