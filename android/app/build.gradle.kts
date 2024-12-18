import io.github.reactivecircus.appversioning.toSemVer
import java.io.ByteArrayInputStream
import java.util.Properties
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

val secretProperties =
    rootProject.layout.projectDirectory
        .file("secrets.properties")
        .let { providers.fileContents(it) }
        .asBytes
        .map { Properties().apply { load(ByteArrayInputStream(it)) } }

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

        val apiHost = System.getenv("PARKEASY_ANDROID_API_HOST") ?: "http://10.0.2.2:8080"
        resValue("string", "api_host", apiHost)

        val protomapsApiKey =
            System.getenv("PARKEASY_ANDROID_PROTOMAPS_API_KEY")
                ?: secretProperties.map { it.getProperty("protomaps.apiKey") }.orNull
                ?: ""
        buildConfigField("String", "PROTOMAPS_API_KEY", "\"${protomapsApiKey}\"")
    }

    signingConfigs {
        register("staging") {
            keyAlias = System.getenv("PARKEASY_ANDROID_STAGING_KEYID")
            keyPassword = System.getenv("PARKEASY_ANDROID_STAGING_KEYPWD")
            storeFile = System.getenv("PARKEASY_ANDROID_STAGING_STORE")?.let { file(it) }
            storePassword = System.getenv("PARKEASY_ANDROID_STAGING_STOREPWD")
        }
    }

    buildTypes {
        val release by getting {
            isMinifyEnabled = true
            isShrinkResources = true
            proguardFiles(
                getDefaultProguardFile("proguard-android-optimize.txt"),
                "proguard-rules.pro",
            )
        }

        register("staging") {
            initWith(release)
            applicationIdSuffix = ".staging"
            signingConfig = signingConfigs.getByName("staging")
        }
    }
    compileOptions {
        isCoreLibraryDesugaringEnabled = true
        sourceCompatibility = JavaVersion.VERSION_11
        targetCompatibility = JavaVersion.VERSION_11
    }
    kotlinOptions { jvmTarget = "11" }
    buildFeatures {
        compose = true
        buildConfig = true
    }
    composeOptions { kotlinCompilerExtensionVersion = "1.5.15" }
    packaging { resources { excludes += "/META-INF/{AL2.0,LGPL2.1}" } }
    lint {
        sarifReport = true
        abortOnError = true

        baseline = file("lint-baseline.xml")
    }
}

appVersioning {
    enabled = providers.environmentVariable("DISABLE_APP_VERSIONING").map { false }.orElse(true)
    gitRootDirectory = rootProject.file("../")

    overrideVersionCode { tag, _, variant ->
        val semVer = tag.toSemVer()
        val baseVer = semVer.major * 1000000 + semVer.minor * 1000 + semVer.patch
        // Add commit number to debug builds
        when (variant.buildType) {
            "debug",
            "staging" -> baseVer * 1000 + min(999, tag.commitsSinceLatestTag)
            else -> baseVer
        }
    }

    overrideVersionName { tag, _, variant ->
        val suffix =
            when (variant.buildType) {
                "debug",
                "staging" -> " (${variant.variantName}, ${tag.commitHash})"
                else -> ""
            }
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

    // Support for JDK 11+
    coreLibraryDesugaring(libs.android.desugar)

    // Lifecycle
    implementation(libs.androidx.lifecycle.runtime.ktx)
    implementation(libs.androidx.lifecycle.runtime.compose)
    implementation(libs.androidx.lifecycle.viewmodel.compose)

    // Compose
    implementation(libs.androidx.activity.compose)
    implementation(platform(libs.androidx.compose.bom))
    implementation(libs.androidx.navigation.compose)
    implementation(libs.androidx.ui)
    implementation(libs.androidx.ui.graphics)
    implementation(libs.androidx.ui.tooling.preview)
    implementation(libs.androidx.material3)
    implementation(libs.material3)
    implementation(libs.accompanist.permissions)

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

    // DateTime
    implementation(libs.kotlinx.datetime)

    // Play Services
    implementation(libs.gms.location)

    // Testing
    testImplementation(libs.junit)
    androidTestImplementation(libs.androidx.junit)
    androidTestImplementation(libs.androidx.espresso.core)
    androidTestImplementation(platform(libs.androidx.compose.bom))
    androidTestImplementation(libs.androidx.ui.test.junit4)
    debugImplementation(libs.androidx.ui.tooling)
    debugImplementation(libs.androidx.ui.test.manifest)
}
