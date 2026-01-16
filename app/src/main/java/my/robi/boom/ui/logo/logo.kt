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
        actions = {
            IconButton(onClick = onToggleTheme) {
                Text(if (darkTheme) "🌙" else "☀️")
            }
        }
    )
}
