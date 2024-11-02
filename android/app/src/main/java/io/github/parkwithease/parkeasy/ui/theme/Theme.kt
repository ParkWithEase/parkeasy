package io.github.parkwithease.parkeasy.ui.theme

import android.os.Build
import androidx.compose.foundation.isSystemInDarkTheme
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.darkColorScheme
import androidx.compose.material3.dynamicDarkColorScheme
import androidx.compose.material3.dynamicLightColorScheme
import androidx.compose.material3.lightColorScheme
import androidx.compose.runtime.Composable
import androidx.compose.runtime.Immutable
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.res.stringResource
import com.maplibre.compose.MapLibreStyleProviding
import com.maplibre.compose.MapLibreSystemThemeStyleProvider
import io.github.parkwithease.parkeasy.BuildConfig
import io.github.parkwithease.parkeasy.R

private val LightScheme =
    lightColorScheme(
        primary = PrimaryLight,
        onPrimary = OnPrimaryLight,
        primaryContainer = PrimaryContainerLight,
        onPrimaryContainer = OnPrimaryContainerLight,
        secondary = SecondaryLight,
        onSecondary = OnSecondaryLight,
        secondaryContainer = SecondaryContainerLight,
        onSecondaryContainer = OnSecondaryContainerLight,
        tertiary = TertiaryLight,
        onTertiary = OnTertiaryLight,
        tertiaryContainer = TertiaryContainerLight,
        onTertiaryContainer = OnTertiaryContainerLight,
        error = ErrorLight,
        onError = OnErrorLight,
        errorContainer = ErrorContainerLight,
        onErrorContainer = OnErrorContainerLight,
        background = BackgroundLight,
        onBackground = OnBackgroundLight,
        surface = SurfaceLight,
        onSurface = OnSurfaceLight,
        surfaceVariant = SurfaceVariantLight,
        onSurfaceVariant = OnSurfaceVariantLight,
        outline = OutlineLight,
        outlineVariant = OutlineVariantLight,
        scrim = ScrimLight,
        inverseSurface = InverseSurfaceLight,
        inverseOnSurface = InverseOnSurfaceLight,
        inversePrimary = InversePrimaryLight,
        surfaceDim = SurfaceDimLight,
        surfaceBright = SurfaceBrightLight,
        surfaceContainerLowest = SurfaceContainerLowestLight,
        surfaceContainerLow = SurfaceContainerLowLight,
        surfaceContainer = SurfaceContainerLight,
        surfaceContainerHigh = SurfaceContainerHighLight,
        surfaceContainerHighest = SurfaceContainerHighestLight,
    )

private val DarkScheme =
    darkColorScheme(
        primary = PrimaryDark,
        onPrimary = OnPrimaryDark,
        primaryContainer = PrimaryContainerDark,
        onPrimaryContainer = OnPrimaryContainerDark,
        secondary = SecondaryDark,
        onSecondary = OnSecondaryDark,
        secondaryContainer = SecondaryContainerDark,
        onSecondaryContainer = OnSecondaryContainerDark,
        tertiary = TertiaryDark,
        onTertiary = OnTertiaryDark,
        tertiaryContainer = TertiaryContainerDark,
        onTertiaryContainer = OnTertiaryContainerDark,
        error = ErrorDark,
        onError = OnErrorDark,
        errorContainer = ErrorContainerDark,
        onErrorContainer = OnErrorContainerDark,
        background = BackgroundDark,
        onBackground = OnBackgroundDark,
        surface = SurfaceDark,
        onSurface = OnSurfaceDark,
        surfaceVariant = SurfaceVariantDark,
        onSurfaceVariant = OnSurfaceVariantDark,
        outline = OutlineDark,
        outlineVariant = OutlineVariantDark,
        scrim = ScrimDark,
        inverseSurface = InverseSurfaceDark,
        inverseOnSurface = InverseOnSurfaceDark,
        inversePrimary = InversePrimaryDark,
        surfaceDim = SurfaceDimDark,
        surfaceBright = SurfaceBrightDark,
        surfaceContainerLowest = SurfaceContainerLowestDark,
        surfaceContainerLow = SurfaceContainerLowDark,
        surfaceContainer = SurfaceContainerDark,
        surfaceContainerHigh = SurfaceContainerHighDark,
        surfaceContainerHighest = SurfaceContainerHighestDark,
    )

@Suppress("unused")
private val MediumContrastLightColorScheme =
    lightColorScheme(
        primary = PrimaryLightMediumContrast,
        onPrimary = OnPrimaryLightMediumContrast,
        primaryContainer = PrimaryContainerLightMediumContrast,
        onPrimaryContainer = OnPrimaryContainerLightMediumContrast,
        secondary = SecondaryLightMediumContrast,
        onSecondary = OnSecondaryLightMediumContrast,
        secondaryContainer = SecondaryContainerLightMediumContrast,
        onSecondaryContainer = OnSecondaryContainerLightMediumContrast,
        tertiary = TertiaryLightMediumContrast,
        onTertiary = OnTertiaryLightMediumContrast,
        tertiaryContainer = TertiaryContainerLightMediumContrast,
        onTertiaryContainer = OnTertiaryContainerLightMediumContrast,
        error = ErrorLightMediumContrast,
        onError = OnErrorLightMediumContrast,
        errorContainer = ErrorContainerLightMediumContrast,
        onErrorContainer = OnErrorContainerLightMediumContrast,
        background = BackgroundLightMediumContrast,
        onBackground = OnBackgroundLightMediumContrast,
        surface = SurfaceLightMediumContrast,
        onSurface = OnSurfaceLightMediumContrast,
        surfaceVariant = SurfaceVariantLightMediumContrast,
        onSurfaceVariant = OnSurfaceVariantLightMediumContrast,
        outline = OutlineLightMediumContrast,
        outlineVariant = OutlineVariantLightMediumContrast,
        scrim = ScrimLightMediumContrast,
        inverseSurface = InverseSurfaceLightMediumContrast,
        inverseOnSurface = InverseOnSurfaceLightMediumContrast,
        inversePrimary = InversePrimaryLightMediumContrast,
        surfaceDim = SurfaceDimLightMediumContrast,
        surfaceBright = SurfaceBrightLightMediumContrast,
        surfaceContainerLowest = SurfaceContainerLowestLightMediumContrast,
        surfaceContainerLow = SurfaceContainerLowLightMediumContrast,
        surfaceContainer = SurfaceContainerLightMediumContrast,
        surfaceContainerHigh = SurfaceContainerHighLightMediumContrast,
        surfaceContainerHighest = SurfaceContainerHighestLightMediumContrast,
    )

@Suppress("unused")
private val HighContrastLightColorScheme =
    lightColorScheme(
        primary = PrimaryLightHighContrast,
        onPrimary = OnPrimaryLightHighContrast,
        primaryContainer = PrimaryContainerLightHighContrast,
        onPrimaryContainer = OnPrimaryContainerLightHighContrast,
        secondary = SecondaryLightHighContrast,
        onSecondary = OnSecondaryLightHighContrast,
        secondaryContainer = SecondaryContainerLightHighContrast,
        onSecondaryContainer = OnSecondaryContainerLightHighContrast,
        tertiary = TertiaryLightHighContrast,
        onTertiary = OnTertiaryLightHighContrast,
        tertiaryContainer = TertiaryContainerLightHighContrast,
        onTertiaryContainer = OnTertiaryContainerLightHighContrast,
        error = ErrorLightHighContrast,
        onError = OnErrorLightHighContrast,
        errorContainer = ErrorContainerLightHighContrast,
        onErrorContainer = OnErrorContainerLightHighContrast,
        background = BackgroundLightHighContrast,
        onBackground = OnBackgroundLightHighContrast,
        surface = SurfaceLightHighContrast,
        onSurface = OnSurfaceLightHighContrast,
        surfaceVariant = SurfaceVariantLightHighContrast,
        onSurfaceVariant = OnSurfaceVariantLightHighContrast,
        outline = OutlineLightHighContrast,
        outlineVariant = OutlineVariantLightHighContrast,
        scrim = ScrimLightHighContrast,
        inverseSurface = InverseSurfaceLightHighContrast,
        inverseOnSurface = InverseOnSurfaceLightHighContrast,
        inversePrimary = InversePrimaryLightHighContrast,
        surfaceDim = SurfaceDimLightHighContrast,
        surfaceBright = SurfaceBrightLightHighContrast,
        surfaceContainerLowest = SurfaceContainerLowestLightHighContrast,
        surfaceContainerLow = SurfaceContainerLowLightHighContrast,
        surfaceContainer = SurfaceContainerLightHighContrast,
        surfaceContainerHigh = SurfaceContainerHighLightHighContrast,
        surfaceContainerHighest = SurfaceContainerHighestLightHighContrast,
    )

@Suppress("unused")
private val MediumContrastDarkColorScheme =
    darkColorScheme(
        primary = PrimaryDarkMediumContrast,
        onPrimary = OnPrimaryDarkMediumContrast,
        primaryContainer = PrimaryContainerDarkMediumContrast,
        onPrimaryContainer = OnPrimaryContainerDarkMediumContrast,
        secondary = SecondaryDarkMediumContrast,
        onSecondary = OnSecondaryDarkMediumContrast,
        secondaryContainer = SecondaryContainerDarkMediumContrast,
        onSecondaryContainer = OnSecondaryContainerDarkMediumContrast,
        tertiary = TertiaryDarkMediumContrast,
        onTertiary = OnTertiaryDarkMediumContrast,
        tertiaryContainer = TertiaryContainerDarkMediumContrast,
        onTertiaryContainer = OnTertiaryContainerDarkMediumContrast,
        error = ErrorDarkMediumContrast,
        onError = OnErrorDarkMediumContrast,
        errorContainer = ErrorContainerDarkMediumContrast,
        onErrorContainer = OnErrorContainerDarkMediumContrast,
        background = BackgroundDarkMediumContrast,
        onBackground = OnBackgroundDarkMediumContrast,
        surface = SurfaceDarkMediumContrast,
        onSurface = OnSurfaceDarkMediumContrast,
        surfaceVariant = SurfaceVariantDarkMediumContrast,
        onSurfaceVariant = OnSurfaceVariantDarkMediumContrast,
        outline = OutlineDarkMediumContrast,
        outlineVariant = OutlineVariantDarkMediumContrast,
        scrim = ScrimDarkMediumContrast,
        inverseSurface = InverseSurfaceDarkMediumContrast,
        inverseOnSurface = InverseOnSurfaceDarkMediumContrast,
        inversePrimary = InversePrimaryDarkMediumContrast,
        surfaceDim = SurfaceDimDarkMediumContrast,
        surfaceBright = SurfaceBrightDarkMediumContrast,
        surfaceContainerLowest = SurfaceContainerLowestDarkMediumContrast,
        surfaceContainerLow = SurfaceContainerLowDarkMediumContrast,
        surfaceContainer = SurfaceContainerDarkMediumContrast,
        surfaceContainerHigh = SurfaceContainerHighDarkMediumContrast,
        surfaceContainerHighest = SurfaceContainerHighestDarkMediumContrast,
    )

@Suppress("unused")
private val HighContrastDarkColorScheme =
    darkColorScheme(
        primary = PrimaryDarkHighContrast,
        onPrimary = OnPrimaryDarkHighContrast,
        primaryContainer = PrimaryContainerDarkHighContrast,
        onPrimaryContainer = OnPrimaryContainerDarkHighContrast,
        secondary = SecondaryDarkHighContrast,
        onSecondary = OnSecondaryDarkHighContrast,
        secondaryContainer = SecondaryContainerDarkHighContrast,
        onSecondaryContainer = OnSecondaryContainerDarkHighContrast,
        tertiary = TertiaryDarkHighContrast,
        onTertiary = OnTertiaryDarkHighContrast,
        tertiaryContainer = TertiaryContainerDarkHighContrast,
        onTertiaryContainer = OnTertiaryContainerDarkHighContrast,
        error = ErrorDarkHighContrast,
        onError = OnErrorDarkHighContrast,
        errorContainer = ErrorContainerDarkHighContrast,
        onErrorContainer = OnErrorContainerDarkHighContrast,
        background = BackgroundDarkHighContrast,
        onBackground = OnBackgroundDarkHighContrast,
        surface = SurfaceDarkHighContrast,
        onSurface = OnSurfaceDarkHighContrast,
        surfaceVariant = SurfaceVariantDarkHighContrast,
        onSurfaceVariant = OnSurfaceVariantDarkHighContrast,
        outline = OutlineDarkHighContrast,
        outlineVariant = OutlineVariantDarkHighContrast,
        scrim = ScrimDarkHighContrast,
        inverseSurface = InverseSurfaceDarkHighContrast,
        inverseOnSurface = InverseOnSurfaceDarkHighContrast,
        inversePrimary = InversePrimaryDarkHighContrast,
        surfaceDim = SurfaceDimDarkHighContrast,
        surfaceBright = SurfaceBrightDarkHighContrast,
        surfaceContainerLowest = SurfaceContainerLowestDarkHighContrast,
        surfaceContainerLow = SurfaceContainerLowDarkHighContrast,
        surfaceContainer = SurfaceContainerDarkHighContrast,
        surfaceContainerHigh = SurfaceContainerHighDarkHighContrast,
        surfaceContainerHighest = SurfaceContainerHighestDarkHighContrast,
    )

@Immutable
data class ColorFamily(
    val color: Color,
    val onColor: Color,
    val colorContainer: Color,
    val onColorContainer: Color,
)

@Suppress("unused")
val unspecified_scheme =
    ColorFamily(Color.Unspecified, Color.Unspecified, Color.Unspecified, Color.Unspecified)

@Composable
fun ParkEasyTheme(
    darkTheme: Boolean = isSystemInDarkTheme(),
    // Dynamic color is available on Android 12+
    dynamicColor: Boolean = false,
    content: @Composable () -> Unit,
) {
    val colorScheme =
        when {
            dynamicColor && Build.VERSION.SDK_INT >= Build.VERSION_CODES.S -> {
                val context = LocalContext.current
                if (darkTheme) dynamicDarkColorScheme(context) else dynamicLightColorScheme(context)
            }

            darkTheme -> DarkScheme
            else -> LightScheme
        }

    MaterialTheme(colorScheme = colorScheme, typography = Typography) {
        MapLibreStyleProviding(
            MapLibreSystemThemeStyleProvider(
                lightModeStyleUrl =
                    stringResource(R.string.map_style_light_format)
                        .format(BuildConfig.PROTOMAPS_API_KEY),
                darkModeStyleUrl =
                    stringResource(R.string.map_style_dark_format)
                        .format(BuildConfig.PROTOMAPS_API_KEY),
            ),
            content = content,
        )
    }
}
