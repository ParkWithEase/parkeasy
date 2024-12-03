package io.github.parkwithease.parkeasy.ui.common

import kotlinx.datetime.Clock
import kotlinx.datetime.DateTimeUnit
import kotlinx.datetime.LocalDate
import kotlinx.datetime.LocalDateTime
import kotlinx.datetime.TimeZone
import kotlinx.datetime.atStartOfDayIn
import kotlinx.datetime.format
import kotlinx.datetime.isoDayNumber
import kotlinx.datetime.minus
import kotlinx.datetime.plus
import kotlinx.datetime.toLocalDateTime

val format =
    LocalDate.Format {
        dayOfMonth()
        chars("/")
        monthNumber()
    }

fun timezone() = TimeZone.currentSystemDefault()

fun startOfWeek() =
    Clock.System.now()
        .toLocalDateTime(timezone())
        .date
        .apply {
            minus(
                dayOfWeek.isoDayNumber - 1, // Monday is 1, Sunday is 7
                DateTimeUnit.DAY,
            )
        }
        .atStartOfDayIn(timezone())
        .toLocalDateTime(timezone())

fun LocalDateTime.isoDay(isoDayNumber: Int) = this.date.plus(isoDayNumber - 1, DateTimeUnit.DAY)

fun LocalDate.toShortDate() = this.format(format)
