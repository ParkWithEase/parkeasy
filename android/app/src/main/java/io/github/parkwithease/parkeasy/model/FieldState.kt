package io.github.parkwithease.parkeasy.model

data class FieldState<T>(val value: T, val error: String? = null)
