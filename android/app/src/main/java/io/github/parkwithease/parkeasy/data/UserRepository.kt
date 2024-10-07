package io.github.parkwithease.parkeasy.data

import io.github.parkwithease.parkeasy.model.LoginCredentials
import io.github.parkwithease.parkeasy.model.RegistrationCredentials

interface UserRepository {
    suspend fun login(credentials: LoginCredentials): Boolean

    suspend fun register(credentials: RegistrationCredentials): Boolean
}
