package io.github.parkwithease.parkeasy.data.remote

import io.github.parkwithease.parkeasy.model.LoginCredentials
import io.github.parkwithease.parkeasy.model.Profile
import io.github.parkwithease.parkeasy.model.RegistrationCredentials
import io.github.parkwithease.parkeasy.model.ResetCredentials

interface UserRepository {
    /**
     * Attempts to log user in with given credentials.
     *
     * @param credentials - the login credentials.
     * @return whether the login was successful or not.
     */
    suspend fun login(credentials: LoginCredentials): Result<Unit>

    /**
     * Attempts to register a new user with the given credentials.
     *
     * @param credentials - the registration credentials.
     * @return whether the registration was successful or not.
     */
    suspend fun register(credentials: RegistrationCredentials): Result<Unit>

    /**
     * Logs the user out.
     *
     * @return whether the logout was successful or not.
     */
    suspend fun logout(): Result<Unit>

    /**
     * Requests for a password reset token to be sent.
     *
     * @param credentials - for the account which the password reset token is for.
     * @return whether the request was successful or not.
     */
    suspend fun requestReset(credentials: ResetCredentials): Result<Unit>

    /**
     * Gets the user details.
     *
     * @return Result with [Profile] of the user if valid user, failing [Result] otherwise.
     */
    suspend fun getUser(): Result<Profile>
}
