package io.github.parkwithease.parkeasy.data.common

import io.github.parkwithease.parkeasy.data.remote.APIException
import io.github.parkwithease.parkeasy.model.ErrorModel
import io.ktor.client.call.body
import io.ktor.client.statement.HttpResponse
import io.ktor.http.isSuccess

// Convert API error into a failing Result
suspend fun Result<HttpResponse>.mapAPIError(): Result<HttpResponse> = mapCatching {
    if (!it.status.isSuccess()) throw APIException(it.body<ErrorModel>())
    it
}
