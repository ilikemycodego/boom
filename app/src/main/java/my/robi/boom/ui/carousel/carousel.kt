package my.robi.boom.ui.carousel

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.pager.HorizontalPager
import androidx.compose.foundation.pager.rememberPagerState
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Surface
import androidx.compose.runtime.Composable
import androidx.compose.runtime.rememberCoroutineScope
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import kotlinx.coroutines.launch
import my.robi.boom.ui.aim.Aim
import my.robi.boom.ui.feed.Feed
import my.robi.boom.ui.logo.Logo
import my.robi.boom.ui.magnet.Magnet
import my.robi.boom.ui.sport.Sport

@Composable
fun CarouselScreen(
    darkTheme: Boolean,
    onToggleTheme: (Boolean) -> Unit
) {
    // Список экранов (каждый экран — composable)
    val pages = listOf<@Composable () -> Unit>(
        { Aim() },
        { Sport() },
        { Feed() },
        { Magnet() }
    )

    val pageCount = pages.size
    val pagerState = rememberPagerState(pageCount = { pageCount })
    val scope = rememberCoroutineScope()

    Scaffold(
        // Верхний бар (логотип + переключатель темы)
        topBar = {
            Logo(
                darkTheme = darkTheme,
                onToggleTheme = { onToggleTheme(!darkTheme) }
            )
        }
    ) { padding ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(padding)
        ) {
            // 1) Свайпающаяся карусель
            HorizontalPager(
                state = pagerState,
                modifier = Modifier
                    .weight(1f)
                    .fillMaxWidth()
            ) { page ->
                Box(
                    modifier = Modifier.fillMaxSize(),
                    contentAlignment = Alignment.Center
                ) {
                    pages[page].invoke()
                }
            }

            // 2) Точки навигации (вторая функция ниже)
            DotsIndicator(
                total = pageCount,
                selected = pagerState.currentPage,
                onSelect = { index ->
                    scope.launch { pagerState.animateScrollToPage(index) }
                },
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(16.dp)
            )
        }
    }
}

@Composable
fun DotsIndicator(
    total: Int,                 // сколько всего точек
    selected: Int,              // какая точка активна (0..total-1)
    onSelect: (Int) -> Unit,    // что делать при тапе на точку
    modifier: Modifier = Modifier
) {
    Row(
        modifier = modifier,
        horizontalArrangement = Arrangement.Center
    ) {
        repeat(total) { i ->
            // Для активной точки делаем ярче/больше
            val alpha = if (i == selected) 1f else 0.35f
            val size = if (i == selected) 10.dp else 8.dp

            Box(
                modifier = Modifier
                    .padding(horizontal = 6.dp)
                    .size(size)
                    .clickable { onSelect(i) } // тап по точке
            ) {
                Surface(
                    modifier = Modifier.fillMaxSize(),
                    color = MaterialTheme.colorScheme.primary.copy(alpha = alpha),
                    shape = MaterialTheme.shapes.small
                ) {}
            }
        }
    }
}
