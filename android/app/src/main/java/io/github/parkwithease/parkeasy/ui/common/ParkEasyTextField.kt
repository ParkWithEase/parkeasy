package io.github.parkwithease.parkeasy.ui.common

import androidx.annotation.DrawableRes
import androidx.annotation.StringRes
import androidx.compose.foundation.interaction.MutableInteractionSource
import androidx.compose.foundation.text.KeyboardActions
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.material3.LocalTextStyle
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.OutlinedTextFieldDefaults
import androidx.compose.material3.TextFieldColors
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Shape
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.res.vectorResource
import androidx.compose.ui.text.TextStyle
import androidx.compose.ui.text.input.VisualTransformation
import io.github.parkwithease.parkeasy.model.FieldState

@Composable
fun ParkEasyTextField(
    state: FieldState<String>,
    onValueChange: (String) -> Unit,
    modifier: Modifier = Modifier,
    enabled: Boolean = true,
    visuallyEnabled: Boolean = enabled,
    readOnly: Boolean = false,
    textStyle: TextStyle = LocalTextStyle.current,
    @StringRes labelId: Int? = null,
    labelString: String? = labelId?.let { stringResource(it) },
    label: @Composable (() -> Unit)? = labelString.textOrNull(),
    placeholder: @Composable (() -> Unit)? = null,
    @DrawableRes leadingIconId: Int? = null,
    leadingIconImage: ImageVector? = leadingIconId?.let { ImageVector.vectorResource(it) },
    leadingIcon: @Composable (() -> Unit)? = leadingIconImage.imageOrNull(),
    trailingIcon: @Composable (() -> Unit)? = null,
    prefix: @Composable (() -> Unit)? = null,
    suffix: @Composable (() -> Unit)? = null,
    supportingText: @Composable (() -> Unit)? = state.error.textOrNull(),
    isError: Boolean = state.error != null,
    visualTransformation: VisualTransformation = VisualTransformation.None,
    keyboardOptions: KeyboardOptions = KeyboardOptions.Default,
    keyboardActions: KeyboardActions = KeyboardActions.Default,
    singleLine: Boolean = true,
    maxLines: Int = if (singleLine) 1 else Int.MAX_VALUE,
    minLines: Int = 1,
    interactionSource: MutableInteractionSource? = null,
    shape: Shape = OutlinedTextFieldDefaults.shape,
    colors: TextFieldColors =
        if (visuallyEnabled)
            OutlinedTextFieldDefaults.colors().run {
                copy(
                    disabledTextColor = unfocusedTextColor,
                    disabledLabelColor = unfocusedLabelColor,
                    disabledContainerColor = unfocusedContainerColor,
                    disabledPrefixColor = unfocusedPrefixColor,
                    disabledSuffixColor = unfocusedSuffixColor,
                    disabledIndicatorColor = unfocusedIndicatorColor,
                    disabledPlaceholderColor = unfocusedPlaceholderColor,
                    disabledLeadingIconColor = unfocusedLeadingIconColor,
                    disabledTrailingIconColor = unfocusedTrailingIconColor,
                    disabledSupportingTextColor = unfocusedSupportingTextColor,
                )
            }
        else OutlinedTextFieldDefaults.colors(),
) {
    OutlinedTextField(
        value = state.value,
        onValueChange = onValueChange,
        modifier = modifier,
        enabled = enabled,
        readOnly = readOnly,
        textStyle = textStyle,
        label = label,
        placeholder = placeholder,
        leadingIcon = leadingIcon,
        trailingIcon = trailingIcon,
        prefix = prefix,
        suffix = suffix,
        supportingText = supportingText,
        isError = isError,
        visualTransformation = visualTransformation,
        keyboardOptions = keyboardOptions,
        keyboardActions = keyboardActions,
        maxLines = maxLines,
        minLines = minLines,
        interactionSource = interactionSource,
        shape = shape,
        colors = colors,
    )
}
