package io.github.parkwithease.parkeasy.data.common

import io.github.parkwithease.parkeasy.data.remote.APIException
import io.github.parkwithease.parkeasy.model.ErrorModel
import io.ktor.client.call.body
import io.ktor.client.statement.HttpResponse
import io.ktor.http.isSuccess

/**
 * Map a [Result] with an unsuccessful [HttpResponse] to an unsuccessful [Result] with an
 * [APIException] or itself if the [HttpResponse] was successful.
 */
suspend fun Result<Any>.mapAPIError(): Result<Any> = mapCatching {
    if (it is HttpResponse && !it.status.isSuccess()) throw APIException(it.body<ErrorModel>())
    it
}
