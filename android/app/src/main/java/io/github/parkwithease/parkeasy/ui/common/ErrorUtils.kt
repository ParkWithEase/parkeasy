package io.github.parkwithease.parkeasy.ui.common

import androidx.compose.material3.Text
import androidx.compose.runtime.Composable

fun String?.textOrNull(): @Composable (() -> Unit)? {
    if (this != null) {
        return { Text(this) }
    } else {
        return null
    }
}
