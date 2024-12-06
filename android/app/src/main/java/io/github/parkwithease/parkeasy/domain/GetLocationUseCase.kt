package io.github.parkwithease.parkeasy.domain

import android.location.Location
import kotlinx.coroutines.flow.Flow

interface GetLocationUseCase {
    operator fun invoke(): Flow<Location>
}
