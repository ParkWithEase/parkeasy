package io.github.parkwithease.parkeasy.ui.navbar

import android.annotation.SuppressLint
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
import androidx.compose.runtime.getValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.res.stringResource
import androidx.navigation.NavController
import androidx.navigation.NavDestination.Companion.hierarchy
import androidx.navigation.NavGraph.Companion.findStartDestination
import androidx.navigation.compose.currentBackStackEntryAsState
import androidx.navigation.compose.rememberNavController
import io.github.parkwithease.parkeasy.R
import io.github.parkwithease.parkeasy.ui.profile.ProfileRoute
import io.github.parkwithease.parkeasy.ui.search.list.ListRoute
import io.github.parkwithease.parkeasy.ui.search.map.MapRoute

@SuppressLint("RestrictedApi")
@Composable
fun NavBar(modifier: Modifier = Modifier, navController: NavController = rememberNavController()) {
    val topLevelRoutes =
        listOf(
            TopLevelRoute(
                stringResource(R.string.list),
                ListRoute,
                Icons.Filled.Menu,
                Icons.Outlined.Menu,
            ),
            TopLevelRoute(
                stringResource(R.string.map),
                MapRoute,
                Icons.Filled.Place,
                Icons.Outlined.Place,
            ),
            TopLevelRoute(
                stringResource(R.string.profile),
                ProfileRoute,
                Icons.Filled.Person,
                Icons.Outlined.Person,
            ),
        )
    val navBackStackEntry by navController.currentBackStackEntryAsState()
    val currentDestination = navBackStackEntry?.destination

    NavigationBar(modifier) {
        topLevelRoutes.forEach { topLevelRoute ->
            val selected =
                currentDestination?.hierarchy?.any { it.hasRoute(topLevelRoute.route, null) } ==
                    true
            NavigationBarItem(
                icon = {
                    Icon(
                        if (selected) topLevelRoute.selectedIcon else topLevelRoute.unselectedIcon,
                        contentDescription = null,
                    )
                },
                label = { Text(topLevelRoute.name) },
                selected = selected,
                onClick = {
                    navController.navigate(topLevelRoute.route) {
                        popUpTo(navController.graph.findStartDestination().id) { saveState = true }
                        launchSingleTop = true
                        restoreState = true
                    }
                },
            )
        }
    }
}

data class TopLevelRoute<T : Any>(
    val name: String,
    val route: T,
    val selectedIcon: ImageVector,
    val unselectedIcon: ImageVector,
)
