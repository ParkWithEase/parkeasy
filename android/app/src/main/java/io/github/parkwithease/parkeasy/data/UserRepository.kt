package io.github.parkwithease.parkeasy.data

import io.github.parkwithease.parkeasy.model.LoginCredentials
import io.github.parkwithease.parkeasy.model.RegistrationCredentials

interface UserRepository {
    suspend fun login(credentials: LoginCredentials)

    suspend fun register(credentials: RegistrationCredentials)
}
