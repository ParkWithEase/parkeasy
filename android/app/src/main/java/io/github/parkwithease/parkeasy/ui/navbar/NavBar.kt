package io.github.parkwithease.parkeasy.ui.navbar

import androidx.compose.animation.AnimatedVisibility
import androidx.compose.animation.slideInVertically
import androidx.compose.animation.slideOutVertically
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Menu
import androidx.compose.material.icons.filled.Person
import androidx.compose.material.icons.filled.Place
import androidx.compose.material.icons.outlined.Menu
import androidx.compose.material.icons.outlined.Person
import androidx.compose.material.icons.outlined.Place
import androidx.compose.material3.Icon
import androidx.compose.material3.NavigationBar
import androidx.compose.material3.NavigationBarItem
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Modifier
import androidx.hilt.navigation.compose.hiltViewModel

@Composable
fun NavBar(
    navigateToList: () -> Unit,
    navigateToMap: () -> Unit,
    navigateToProfile: () -> Unit,
    modifier: Modifier = Modifier,
    viewModel: NavBarViewModel = hiltViewModel<NavBarViewModel>(),
) {
    val selectedItem by viewModel.selectedItem.collectAsState()
    val items = listOf("List", "Map", "Profile")
    val selectedIcons = listOf(Icons.Filled.Menu, Icons.Filled.Place, Icons.Filled.Person)
    val unselectedIcons = listOf(Icons.Outlined.Menu, Icons.Outlined.Place, Icons.Outlined.Person)
    val loggedIn by viewModel.loggedIn.collectAsState(false)
    val navigateTo = listOf(navigateToList, navigateToMap, navigateToProfile)

    AnimatedVisibility(
        visible = loggedIn,
        enter = slideInVertically(initialOffsetY = { it / 2 }),
        exit = slideOutVertically(targetOffsetY = { it / 2 }),
    ) {
        NavigationBar(modifier) {
            items.forEachIndexed { index, item ->
                NavigationBarItem(
                    icon = {
                        Icon(
                            if (selectedItem == index) selectedIcons[index]
                            else unselectedIcons[index],
                            contentDescription = item,
                        )
                    },
                    label = { Text(item) },
                    selected = selectedItem == index,
                    onClick = {
                        viewModel.onClick(index)
                        navigateTo[index]()
                    },
                )
            }
        }
    }
}
