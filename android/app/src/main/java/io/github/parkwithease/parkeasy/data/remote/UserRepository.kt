package io.github.parkwithease.parkeasy.data.remote

import io.github.parkwithease.parkeasy.model.LoginCredentials
import io.github.parkwithease.parkeasy.model.RegistrationCredentials
import io.github.parkwithease.parkeasy.model.ResetCredentials

interface UserRepository {
    /**
     * Attempts to log user in with given credentials.
     *
     * @param credentials the login credentials.
     * @return true if user login success, false otherwise.
     */
    suspend fun login(credentials: LoginCredentials): Boolean

    /**
     * Attempts to register a new user with the given credentials.
     *
     * @param credentials the registration credentials.
     * @return true the new user is registered, false otherwise.
     */
    suspend fun register(credentials: RegistrationCredentials): Boolean

    /**
     * Logs the user out.
     */
    suspend fun logout()

    /**
     * Requests for a password reset token to be sent.
     *
     * @param credentials for the account which the password reset token is for.
     */
    suspend fun reset(credentials: ResetCredentials): Boolean
}
