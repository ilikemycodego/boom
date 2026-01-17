package my.robi.boom.ui.logo

import androidx.compose.material3.*
import androidx.compose.runtime.Composable

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun Logo(
    darkTheme: Boolean,
    onToggleTheme: () -> Unit
) {
    TopAppBar(
        title = { Text("BOOM") },
        colors = TopAppBarDefaults.topAppBarColors(
            containerColor = MaterialTheme.colorScheme.surface,
            titleContentColor = MaterialTheme.colorScheme.onSurface,
            actionIconContentColor = MaterialTheme.colorScheme.onSurface
        ),
        actions = {
            IconButton(onClick = onToggleTheme) {
                Text(if (darkTheme) "🌙" else "☀️")
            }
        }
    )
}
