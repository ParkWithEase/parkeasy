package io.github.parkwithease.parkeasy.ui.nav

import androidx.lifecycle.ViewModel
import dagger.hilt.android.lifecycle.HiltViewModel
import javax.inject.Inject
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow

@HiltViewModel
class NavBarViewModel @Inject constructor() : ViewModel() {
    private val _selectedItem = MutableStateFlow(0)
    val selectedItem = _selectedItem.asStateFlow()

    fun onClick(index: Int) {
        _selectedItem.value = index
    }
}
