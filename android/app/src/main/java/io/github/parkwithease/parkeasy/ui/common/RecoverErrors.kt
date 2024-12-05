package io.github.parkwithease.parkeasy.ui.common

import androidx.compose.material3.SnackbarHostState
import io.github.parkwithease.parkeasy.data.remote.APIException
import io.github.parkwithease.parkeasy.model.ErrorModel
import java.io.IOException
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.launch

/** Recovers errors, handling [APIException]s and [IOException]s. */
fun Result<Any>.recoverRequestErrors(
    operationFailMsg: String,
    handleError: (ErrorModel) -> Unit,
    snackbarState: SnackbarHostState,
    scope: CoroutineScope,
): Result<Any> = recover {
    when (it) {
        is APIException -> {
            handleError(it.error)
            scope.launch { snackbarState.showSnackbar(operationFailMsg) }
        }
        is IOException -> {
            scope.launch {
                snackbarState.showSnackbar("Could not connect to server, are you online?")
            }
        }
        else -> throw it
    }
}
