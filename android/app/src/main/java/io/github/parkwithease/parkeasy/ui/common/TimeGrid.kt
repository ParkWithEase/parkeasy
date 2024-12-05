package io.github.parkwithease.parkeasy.ui.common

import androidx.compose.foundation.gestures.detectDragGestures
import androidx.compose.foundation.interaction.MutableInteractionSource
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.lazy.grid.GridCells
import androidx.compose.foundation.lazy.grid.GridItemSpan
import androidx.compose.foundation.lazy.grid.LazyGridState
import androidx.compose.foundation.lazy.grid.LazyVerticalGrid
import androidx.compose.foundation.lazy.grid.items
import androidx.compose.foundation.lazy.grid.rememberLazyGridState
import androidx.compose.foundation.selection.toggleable
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.remember
import androidx.compose.ui.Modifier
import androidx.compose.ui.geometry.Offset
import androidx.compose.ui.input.pointer.pointerInput
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.round
import androidx.compose.ui.unit.toIntRect
import io.github.parkwithease.parkeasy.ui.spots.NumColumns
import io.github.parkwithease.parkeasy.ui.spots.NumRows
import io.github.parkwithease.parkeasy.ui.spots.NumSlots
import kotlinx.datetime.LocalTime

@Suppress("detekt:LongMethod")
@Composable
fun TimeGrid(
    getSelectedIds: () -> Set<Int>,
    disabledIds: Set<Int>,
    onAddTime: (elements: Iterable<Int>) -> Unit,
    onRemoveTime: (elements: Iterable<Int>) -> Unit,
    modifier: Modifier = Modifier,
    state: LazyGridState = rememberLazyGridState(),
    slots: List<Int> = List(NumSlots) { it % NumColumns * NumRows + it / NumColumns },
) {
    val selectedIds: Set<Int> = getSelectedIds()

    Row(modifier.height(336.dp)) {
        ColumnHeader(Modifier.width(48.dp))
        LazyVerticalGrid(
            columns = GridCells.Fixed(NumColumns),
            modifier =
                Modifier.timeGridDragHandler(
                    lazyGridState = state,
                    disabledIds = disabledIds,
                    getSelectedIds = getSelectedIds,
                    onAddTime = onAddTime,
                    onRemoveTime = onRemoveTime,
                ),
            state = state,
            horizontalArrangement = Arrangement.spacedBy(3.dp),
        ) {
            items(count = NumColumns / 2, span = { GridItemSpan(2) }) {
                Text(
                    text = startOfNextAvailableDay().isoDay(it + 1).toShortDate(),
                    Modifier.height(24.dp),
                    style = MaterialTheme.typography.labelLarge,
                )
            }
            items(NumColumns) {
                Text(
                    text = if (it % 2 == 0) "AM" else "PM",
                    Modifier.height(24.dp),
                    style = MaterialTheme.typography.labelLarge,
                )
            }
            items(slots, key = { it }) { id ->
                val disabled = disabledIds.contains(id)
                val selected = selectedIds.contains(id)

                Surface(
                    tonalElevation = 3.dp,
                    color =
                        if (disabled) {
                            MaterialTheme.colorScheme.onSurface
                        } else {
                            if (selected) MaterialTheme.colorScheme.primary
                            else MaterialTheme.colorScheme.primaryContainer
                        },
                    modifier =
                        Modifier.height(12.dp)
                            .padding(top = if (id % 24 > 0 && id % 24 % 2 == 0) 3.dp else 0.dp)
                            .toggleable(
                                value = selected,
                                interactionSource = remember { MutableInteractionSource() },
                                indication = null, // do not show a ripple
                                onValueChange = {
                                    if (!disabled) {
                                        if (it) {
                                            onAddTime(id..id)
                                        } else {
                                            onRemoveTime(id..id)
                                        }
                                    }
                                },
                            ),
                    content = {},
                )
            }
        }
    }
}

@Composable
fun ColumnHeader(modifier: Modifier = Modifier, state: LazyGridState = rememberLazyGridState()) {
    LazyVerticalGrid(state = state, columns = GridCells.Fixed(1), modifier = modifier) {
        items(2) {
            Text(text = "", Modifier.height(24.dp), style = MaterialTheme.typography.labelLarge)
        }
        items(NumRows) { num ->
            if (num % 2 == 0) {
                Text(
                    text = LocalTime((num / 2), 0).toString(),
                    Modifier.height(24.dp),
                    style = MaterialTheme.typography.labelLarge,
                )
            }
        }
    }
}

@Suppress("detekt:UnsafeCallOnNullableType") // code provided by a Google engineer -> probably fine
private fun Modifier.timeGridDragHandler(
    lazyGridState: LazyGridState,
    disabledIds: Set<Int>,
    getSelectedIds: () -> Set<Int>,
    onAddTime: (elements: Iterable<Int>) -> Unit,
    onRemoveTime: (elements: Iterable<Int>) -> Unit,
) =
    pointerInput(Unit) {
        fun LazyGridState.gridItemKeyAtPosition(hitPoint: Offset): Int? =
            layoutInfo.visibleItemsInfo
                .find { itemInfo ->
                    itemInfo.size.toIntRect().contains(hitPoint.round() - itemInfo.offset)
                }
                ?.key as? Int
        var initialKey: Int? = null
        var currentKey: Int? = null
        var adding = false
        detectDragGestures(
            onDragStart = { offset ->
                lazyGridState.gridItemKeyAtPosition(offset)?.let { key ->
                    val selectedIds = getSelectedIds()
                    if (!disabledIds.contains(key)) {
                        initialKey = key
                        currentKey = key
                        if (!selectedIds.contains(key)) {
                            onAddTime(key..key)
                            adding = true
                        } else {
                            onRemoveTime(key..key)
                            adding = false
                        }
                    }
                }
            },
            onDragCancel = { initialKey = null },
            onDragEnd = { initialKey = null },
            onDrag = { change, _ ->
                if (initialKey != null) {
                    lazyGridState.gridItemKeyAtPosition(change.position)?.let { key ->
                        if (currentKey != key) {
                            if (adding) {
                                onRemoveTime(initialKey!!..currentKey!!)
                                onRemoveTime(currentKey!!..initialKey!!)
                                onAddTime(initialKey!!..key)
                                onAddTime(key..initialKey!!)
                            } else {
                                onAddTime(initialKey!!..currentKey!!)
                                onAddTime(currentKey!!..initialKey!!)
                                onRemoveTime(initialKey!!..key)
                                onRemoveTime(key..initialKey!!)
                            }
                            disabledIds.forEach { onRemoveTime(it..it) }
                            currentKey = key
                        }
                    }
                }
            },
        )
    }
