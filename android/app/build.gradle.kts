import io.github.reactivecircus.appversioning.toSemVer
import kotlin.math.min

plugins {
    alias(libs.plugins.android.application)
    alias(libs.plugins.kotlin.android)
    alias(libs.plugins.kotlin.serialization)
    alias(libs.plugins.compose.compiler)
    alias(libs.plugins.kotlin.ksp)
    alias(libs.plugins.hilt)
    alias(libs.plugins.kotlinx.kover)
    alias(libs.plugins.app.versioning)
}

android {
    namespace = "io.github.parkwithease.parkeasy"
    compileSdk = 35

    defaultConfig {
        applicationId = "io.github.parkwithease.parkeasy"
        minSdk = 25
        targetSdk = 35
        versionCode = 1
        versionName = "1.0"

        testInstrumentationRunner = "androidx.test.runner.AndroidJUnitRunner"
        vectorDrawables { useSupportLibrary = true }
    }

    buildTypes {
        release {
            isMinifyEnabled = false
            proguardFiles(
                getDefaultProguardFile("proguard-android-optimize.txt"),
                "proguard-rules.pro",
            )
        }
    }
    compileOptions {
        sourceCompatibility = JavaVersion.VERSION_1_8
        targetCompatibility = JavaVersion.VERSION_1_8
    }
    kotlinOptions { jvmTarget = "1.8" }
    buildFeatures { compose = true }
    composeOptions { kotlinCompilerExtensionVersion = "1.5.1" }
    packaging { resources { excludes += "/META-INF/{AL2.0,LGPL2.1}" } }
    lint {
        sarifReport = true
        abortOnError = true

        baseline = file("lint-baseline.xml")
    }
}

appVersioning {
    gitRootDirectory = rootProject.file("../")

    overrideVersionCode { tag, _, variant ->
        val semVer = tag.toSemVer()
        val baseVer = semVer.major * 1000000 + semVer.minor * 1000 + semVer.patch
        // Add commit number to debug builds
        if (variant.isDebugBuild) baseVer * 1000 + min(999, tag.commitsSinceLatestTag) else baseVer
    }

    overrideVersionName { tag, _, variant ->
        val suffix =
            if (variant.isDebugBuild) " (${variant.variantName}, ${tag.commitHash})" else ""
        tag.toString().removePrefix("v") + suffix
    }
}

kover {
    reports {
        filters {
            excludes {
                annotatedBy("*Generated*")
                classes("*\$BindsModule", "*\$KeyModule", "*\$InstanceHolder")
                packages("hilt_aggregated_deps", "dagger.*")
            }
        }
    }
}

dependencies {
    implementation(libs.androidx.core.ktx)
    implementation(libs.androidx.datastore.preferences)

    // Lifecycle
    implementation(libs.androidx.lifecycle.runtime.ktx)
    implementation(libs.androidx.lifecycle.runtime.compose)
    implementation(libs.androidx.lifecycle.viewmodel.compose)

    // Compose
    implementation(libs.androidx.activity.compose)
    implementation(platform(libs.androidx.compose.bom))
    implementation(libs.androidx.ui)
    implementation(libs.androidx.ui.graphics)
    implementation(libs.androidx.ui.tooling.preview)
    implementation(libs.androidx.material3)
    implementation(libs.material3)

    // MapLibre
    implementation(libs.maplibre.compose)

    // Ktor
    implementation(libs.ktor.client.core)
    implementation(libs.ktor.client.okhttp)
    implementation(libs.ktor.client.logging)
    implementation(libs.ktor.client.content.negotiation)
    implementation(libs.ktor.serialization.kotlinx.json)

    // Hilt
    implementation(libs.hilt.android)
    ksp(libs.hilt.compiler)
    implementation(libs.androidx.hilt.navigation.compose)

    // Testing
    testImplementation(libs.junit)
    androidTestImplementation(libs.androidx.junit)
    androidTestImplementation(libs.androidx.espresso.core)
    androidTestImplementation(platform(libs.androidx.compose.bom))
    androidTestImplementation(libs.androidx.ui.test.junit4)
    debugImplementation(libs.androidx.ui.tooling)
    debugImplementation(libs.androidx.ui.test.manifest)
}
