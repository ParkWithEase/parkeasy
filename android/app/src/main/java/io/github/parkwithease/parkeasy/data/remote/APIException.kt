package io.github.parkwithease.parkeasy.data.remote

import io.github.parkwithease.parkeasy.model.ErrorModel

data class APIException(val error: ErrorModel) : Exception(error.detail)
