package io.github.parkwithease.parkeasy.data

import kotlinx.serialization.Serializable

@Serializable
data class Credentials(val email: String, val password: String, val persist: Boolean = true)
