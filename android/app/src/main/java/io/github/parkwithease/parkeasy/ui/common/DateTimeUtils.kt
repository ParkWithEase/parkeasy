package io.github.parkwithease.parkeasy.ui.common

import kotlinx.datetime.Clock
import kotlinx.datetime.DateTimeUnit
import kotlinx.datetime.LocalDate
import kotlinx.datetime.LocalDateTime
import kotlinx.datetime.TimeZone
import kotlinx.datetime.atStartOfDayIn
import kotlinx.datetime.format
import kotlinx.datetime.plus
import kotlinx.datetime.toLocalDateTime

const val MinutesPerSlot = 30

val format =
    LocalDate.Format {
        dayOfMonth()
        chars("/")
        monthNumber()
    }

fun timezone() = TimeZone.currentSystemDefault()

fun startOfNextAvailableDay() =
    Clock.System.now()
        .plus(MinutesPerSlot, DateTimeUnit.MINUTE)
        .toLocalDateTime(timezone())
        .date
        .atStartOfDayIn(timezone())
        .toLocalDateTime(timezone())

fun LocalDateTime.isoDay(isoDayNumber: Int) = this.date.plus(isoDayNumber - 1, DateTimeUnit.DAY)

fun LocalDate.toShortDate() = this.format(format)
