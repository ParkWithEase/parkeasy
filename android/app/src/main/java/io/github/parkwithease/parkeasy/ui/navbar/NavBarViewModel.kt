package io.github.parkwithease.parkeasy.ui.navbar

import androidx.lifecycle.ViewModel
import dagger.hilt.android.lifecycle.HiltViewModel
import io.github.parkwithease.parkeasy.data.local.AuthRepository
import javax.inject.Inject
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow

@HiltViewModel
class NavBarViewModel @Inject constructor(authRepo: AuthRepository) : ViewModel() {
    private val _selectedItem = MutableStateFlow(2)
    val selectedItem = _selectedItem.asStateFlow()
    val loggedIn = authRepo.statusFlow

    fun onClick(index: Int) {
        _selectedItem.value = index
    }
}
