package io.github.parkwithease.parkeasy.ui.common

import androidx.compose.foundation.Image
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.graphics.vector.ImageVector

fun String?.textOrNull(): @Composable (() -> Unit)? {
    return if (this != null) {
        { Text(this) }
    } else {
        null
    }
}

fun ImageVector?.imageOrNull(): @Composable (() -> Unit)? {
    return if (this != null) {
        { Image(this, contentDescription = null) }
    } else {
        null
    }
}
