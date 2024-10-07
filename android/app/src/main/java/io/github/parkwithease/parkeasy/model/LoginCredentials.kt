package io.github.parkwithease.parkeasy.model

import kotlinx.serialization.Serializable

@Serializable
data class LoginCredentials(val email: String, val password: String, val persist: Boolean = true)
