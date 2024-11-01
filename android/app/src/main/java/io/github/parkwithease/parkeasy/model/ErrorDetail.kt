package io.github.parkwithease.parkeasy.model

import kotlinx.serialization.Serializable
import kotlinx.serialization.json.JsonElement

@Serializable
data class ErrorDetail(val location: String, val message: String? = null, val value: JsonElement)
