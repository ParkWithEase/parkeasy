package io.github.parkwithease.parkeasy.model

import kotlinx.serialization.Serializable

@Serializable
data class ErrorModel(
    val detail: String? = null,
    val status: Int,
    val type: String? = null,
    val title: String,
    val errors: List<ErrorDetail> = emptyList(),
) {
    companion object {
        const val TYPE_INVALID_CREDENTIALS =
            "tag:parkwithease.github.io,2024-10-13:problem:invalid-credentials"
        const val TYPE_PASSWORD_LENGTH =
            "tag:parkwithease.github.io,2024-10-13:problem:password-length"
    }
}
