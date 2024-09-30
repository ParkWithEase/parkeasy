package io.github.parkwithease.parkeasy

import kotlinx.serialization.Serializable

@Serializable
data class Credentials(
    val email: String,
    val password: String,
    val persist: Boolean = true
)
