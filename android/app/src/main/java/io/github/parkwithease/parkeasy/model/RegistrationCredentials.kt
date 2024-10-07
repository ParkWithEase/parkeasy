package io.github.parkwithease.parkeasy.model

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class RegistrationCredentials(
    @SerialName("full_name") val name: String,
    val email: String,
    val password: String,
    val persist: Boolean = true,
)
