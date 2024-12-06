package io.github.parkwithease.parkeasy.ui.common

import java.time.format.DateTimeFormatter
import kotlin.time.DurationUnit
import kotlinx.datetime.Clock
import kotlinx.datetime.DateTimeUnit
import kotlinx.datetime.Instant
import kotlinx.datetime.LocalDate
import kotlinx.datetime.LocalDateTime
import kotlinx.datetime.TimeZone
import kotlinx.datetime.atStartOfDayIn
import kotlinx.datetime.format
import kotlinx.datetime.format.char
import kotlinx.datetime.plus
import kotlinx.datetime.toInstant
import kotlinx.datetime.toJavaLocalDateTime
import kotlinx.datetime.toJavaZoneId
import kotlinx.datetime.toLocalDateTime

const val MinutesPerSlot = 30

private val ShortDate by lazy {
    LocalDate.Format {
        dayOfMonth()
        char('/')
        monthNumber()
    }
}

private val FullDate by lazy { DateTimeFormatter.ofPattern("EEE, MMM dd, uuuu HH:mm:ss") }

fun timezone() = TimeZone.currentSystemDefault()

fun startOfNextAvailableDay() =
    Clock.System.now()
        .plus(MinutesPerSlot, DateTimeUnit.MINUTE)
        .toLocalDateTime(timezone())
        .date
        .atStartOfDayIn(timezone())
        .toLocalDateTime(timezone())

fun startOfNextAvailableDayInstant() = startOfNextAvailableDay().toInstant(timezone())

fun LocalDateTime.isoDay(isoDayNumber: Int) = this.date.plus(isoDayNumber - 1, DateTimeUnit.DAY)

fun LocalDate.toShortDate() = this.format(ShortDate)

fun LocalDateTime.toReadable(): String =
    with(this.toJavaLocalDateTime().atZone(timezone().toJavaZoneId())) { this.format(FullDate) }

fun Instant.toIndex() =
    this.minus(startOfNextAvailableDay().toInstant(timezone())).toInt(DurationUnit.MINUTES) /
        MinutesPerSlot

fun Int.toLocalDateTime() =
    startOfNextAvailableDayInstant()
        .plus(MinutesPerSlot * this, DateTimeUnit.MINUTE, timezone())
        .toLocalDateTime(timezone())
