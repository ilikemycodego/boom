package my.robi.boom

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Surface
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import my.robi.boom.ui.BoomTheme
import my.robi.boom.ui.carousel.CarouselScreen

class Main : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContent {
            // Флаг темы (false = светлая, true = тёмная)
            var darkTheme by remember { mutableStateOf(false) }

            // Оборачиваем всё приложение в нашу тему
            BoomTheme(darkTheme = darkTheme) {
                // Фон на весь экран, чтобы не было “пятен” от системного фона
                Surface(
                    modifier = Modifier.fillMaxSize(),
                    color = MaterialTheme.colorScheme.background
                ) {
                    CarouselScreen(
                        darkTheme = darkTheme,
                        onToggleTheme = { darkTheme = it }
                    )
                }
            }
        }
    }
}
