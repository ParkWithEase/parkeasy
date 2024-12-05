package io.github.parkwithease.parkeasy.domain

import android.content.Context
import android.content.pm.PackageManager
import android.location.Location
import android.os.Looper
import androidx.core.content.ContextCompat
import com.google.android.gms.location.LocationCallback
import com.google.android.gms.location.LocationRequest
import com.google.android.gms.location.LocationResult
import com.google.android.gms.location.LocationServices
import com.google.android.gms.location.Priority
import dagger.hilt.android.qualifiers.ApplicationContext
import io.github.parkwithease.parkeasy.di.IoDispatcher
import javax.inject.Inject
import kotlin.time.Duration.Companion.seconds
import kotlinx.coroutines.CoroutineDispatcher
import kotlinx.coroutines.channels.awaitClose
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.callbackFlow
import kotlinx.coroutines.flow.flowOn
import kotlinx.coroutines.launch

private const val MIN_UPDATE_DISTANCE = 100F

class GetLocationUseCaseImpl
@Inject
constructor(
    @ApplicationContext private val context: Context,
    @IoDispatcher private val ioDispatcher: CoroutineDispatcher,
) : GetLocationUseCase {
    private val client by lazy { LocationServices.getFusedLocationProviderClient(context) }
    private val request =
        LocationRequest.Builder(Priority.PRIORITY_HIGH_ACCURACY, 10.seconds.inWholeMilliseconds)
            .setMinUpdateDistanceMeters(MIN_UPDATE_DISTANCE)
            .build()

    override fun invoke(): Flow<Location> =
        callbackFlow {
                check(
                    ContextCompat.checkSelfPermission(
                        context,
                        android.Manifest.permission.ACCESS_FINE_LOCATION,
                    ) == PackageManager.PERMISSION_GRANTED
                )

                val callback =
                    object : LocationCallback() {
                        override fun onLocationResult(result: LocationResult) {
                            super.onLocationResult(result)
                            result.lastLocation?.let { launch { send(it) } }
                        }
                    }

                client.requestLocationUpdates(request, callback, Looper.getMainLooper())
                awaitClose { client.removeLocationUpdates(callback) }
            }
            .flowOn(ioDispatcher)
}
