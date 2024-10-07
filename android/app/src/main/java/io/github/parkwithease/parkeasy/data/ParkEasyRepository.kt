package io.github.parkwithease.parkeasy.data

import io.github.parkwithease.parkeasy.model.Credentials

interface ParkEasyRepository {
    suspend fun login(credentials: Credentials)
}
