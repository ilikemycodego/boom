package my.robi.boom.ui

import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.interaction.MutableInteractionSource
import androidx.compose.foundation.interaction.collectIsPressedAsState
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.wrapContentWidth
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.luminance
import androidx.compose.ui.unit.dp
import androidx.compose.foundation.shape.RoundedCornerShape



// --------------------
// 1) Палитра (строго ч/б + серые нейтральные для тёмной темы)
// --------------------
private val ControlShape = RoundedCornerShape(6.dp) // меньше радиус = более “квадратно”

private val LightScheme = lightColorScheme(
    // Фон/поверхности
    background = Color(0xFFFFFFFF),
    surface = Color(0xFFFFFFFF),

    // Текст
    onBackground = Color(0xFF000000),
    onSurface = Color(0xFF000000),

    // Важно: primary тоже ч/б, чтобы не было “синих” акцентов
    primary = Color(0xFF000000),
    onPrimary = Color(0xFFFFFFFF),

    // Обводки (в светлой теме можно прям чёрный)
    outline = Color(0xFF000000)
)

private val DarkScheme = darkColorScheme(
    // Фон/поверхности (как “безопасный” тёмный UI)
    background = Color(0xFF0B0F14),
    surface = Color(0xFF0F141A),

    // Текст
    onBackground = Color(0xFFE6E6E6),
    onSurface = Color(0xFFE6E6E6),

    // primary тоже нейтральный (без цвета)
    primary = Color(0xFFE6E6E6),
    onPrimary = Color(0xFF0B0F14),

    // Обводки мягкие серые
    outline = Color(0xFFE6E6E6)
)

/**
 * Главная тема приложения: все цвета задаются здесь.
 */
@Composable
fun BoomTheme(
    darkTheme: Boolean,
    content: @Composable () -> Unit
) {
    MaterialTheme(
        colorScheme = if (darkTheme) DarkScheme else LightScheme,
        content = content
    )
}

// --------------------
// 2) Единый стиль контролов (аналог “CSS для кнопок и инпутов”)
// --------------------

/**
 * Единая кнопка для всего приложения:
 * - прозрачная
 * - с обводкой
 * - ширина как у поля ввода
 * - высота как у стандартного TextField (56dp)
 */

@Composable
fun BoomButton(
    text: String,
    onClick: () -> Unit,
    modifier: Modifier = Modifier,
    enabled: Boolean = true
) {
    val isDark = MaterialTheme.colorScheme.background.luminance() < 0.5f

    val interaction = remember { MutableInteractionSource() }
    val pressed by interaction.collectIsPressedAsState()

    // “Сильная” обводка как у выделенного TextField
    val strongOutline = MaterialTheme.colorScheme.onSurface.copy(alpha = if (isDark) 0.9f else 1f)

    // Небольшая подсветка при нажатии (чтобы клик был заметнее)
    val pressedBg = if (isDark) Color(0x26FFFFFF) else Color(0x14000000)

    OutlinedButton(
        onClick = onClick,
        enabled = enabled,
        interactionSource = interaction,
        modifier = modifier
            .wrapContentWidth() // ширина по тексту
            .height(56.dp),
        shape = ControlShape, // такая же форма как у поля
        border = BorderStroke(
            width = if (isDark) 1.dp else 1.dp, // усилить в тёмной теме
            color = strongOutline               // ВСЕГДА сильная обводка
        ),
        colors = ButtonDefaults.outlinedButtonColors(
            containerColor = if (pressed) pressedBg else Color.Transparent,
            contentColor = MaterialTheme.colorScheme.onSurface
        )
    ) {
        Text(text)
    }
}

/**
 * Единое поле ввода для всего приложения:
 * - прозрачное
 * - обводка всегда “ч/б” (без синего при фокусе)
 */
@Composable
fun BoomTextField(
    value: String,
    onValueChange: (String) -> Unit,
    label: String,
    modifier: Modifier = Modifier
) {
    val isDark = MaterialTheme.colorScheme.background.luminance() < 0.5f

    val outline = MaterialTheme.colorScheme.outline
    val strongOutline = MaterialTheme.colorScheme.onSurface.copy(alpha = if (isDark) 0.9f else 1f)

    // Подсветка фокуса (чтобы выделение было заметнее)
    val focusedFill = if (isDark) Color(0x0FFFFFFF) else Color(0x0A000000)

    OutlinedTextField(
        value = value,
        onValueChange = onValueChange,
        modifier = modifier,
        label = { Text(label) },
        singleLine = true,
        shape = ControlShape, // чтобы форма совпадала с кнопкой
        colors = OutlinedTextFieldDefaults.colors(
            focusedContainerColor = focusedFill,
            unfocusedContainerColor = Color.Transparent,

            focusedTextColor = MaterialTheme.colorScheme.onSurface,
            unfocusedTextColor = MaterialTheme.colorScheme.onSurface,

            // в фокусе — сильная рамка (как ты хочешь)
            focusedBorderColor = strongOutline,
            unfocusedBorderColor = outline,

            cursorColor = MaterialTheme.colorScheme.onSurface
        )
    )
}