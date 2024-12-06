package io.github.parkwithease.parkeasy.ui.common

import androidx.compose.animation.core.AnimationConstants.DefaultDurationMillis
import androidx.compose.animation.core.tween
import androidx.compose.animation.fadeIn
import androidx.compose.animation.fadeOut

fun enterAnimation() = fadeIn(animationSpec = tween())

fun exitAnimation() = fadeOut(animationSpec = tween(delayMillis = DefaultDurationMillis))
