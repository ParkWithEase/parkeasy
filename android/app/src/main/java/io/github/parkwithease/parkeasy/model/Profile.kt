package io.github.parkwithease.parkeasy.model

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable data class Profile(@SerialName("full_name") val name: String, val email: String)
