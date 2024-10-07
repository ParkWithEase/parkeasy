// Top-level build file where you can add configuration options common to all sub-projects/modules.
plugins {
    alias(libs.plugins.android.application) apply false
    alias(libs.plugins.kotlin.android) apply false
    alias(libs.plugins.ktfmt)
    alias(libs.plugins.detekt)
    alias(libs.plugins.hilt) apply false
    alias(libs.plugins.kotlin.ksp) apply false
    alias(libs.plugins.compose.compiler) apply false
}

val mergeDetektReports by
    tasks.registering(io.gitlab.arturbosch.detekt.report.ReportMergeTask::class) {
        output.set(rootProject.layout.buildDirectory.file("reports/detekt/merged.sarif"))
    }

val mergeLintReports by
    tasks.registering(io.gitlab.arturbosch.detekt.report.ReportMergeTask::class) {
        output.set(rootProject.layout.buildDirectory.file("reports/lint/merged.sarif"))
    }

val detektCompose = libs.detekt.compose

subprojects {
    apply {
        plugin("com.ncorti.ktfmt.gradle")
        plugin("io.gitlab.arturbosch.detekt")
    }

    detekt {
        config.from(rootProject.file("config/detekt/detekt.yml"))
        basePath = projectDir.toString()
        buildUponDefaultConfig = true
    }

    tasks.withType<io.gitlab.arturbosch.detekt.Detekt>().configureEach {
        reports.sarif.required = true
        finalizedBy(mergeDetektReports)
    }

    tasks.withType<com.android.build.gradle.internal.lint.AndroidLintTask>().configureEach {
        if (name == "lintReportDebug") {
            sarifReportEnabled = true
            finalizedBy(mergeLintReports)
        }
    }

    mergeDetektReports {
        input.from(tasks.withType<io.gitlab.arturbosch.detekt.Detekt>().map { it.sarifReportFile })
    }

    mergeLintReports {
        input.from(
            tasks.named("lintReportDebug", com.android.build.gradle.internal.lint.AndroidLintTask::class).map {
                it.sarifReportOutputFile
            }
        )
    }

    ktfmt { kotlinLangStyle() }

    dependencies { detektPlugins(detektCompose) }
}
