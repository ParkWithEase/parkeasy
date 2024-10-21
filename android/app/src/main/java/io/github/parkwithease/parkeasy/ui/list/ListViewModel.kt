package io.github.parkwithease.parkeasy.ui.list

import androidx.lifecycle.ViewModel
import dagger.hilt.android.lifecycle.HiltViewModel
import io.github.parkwithease.parkeasy.data.local.AuthRepository
import javax.inject.Inject

@HiltViewModel
class ListViewModel @Inject constructor(authRepo: AuthRepository) : ViewModel() {
    val loggedIn = authRepo.statusFlow
}
