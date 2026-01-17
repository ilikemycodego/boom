package my.robi.boom.ui.magnet

import androidx.compose.foundation.layout.*

import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.*
import androidx.compose.runtime.saveable.rememberSaveable

import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import my.robi.boom.ui.BoomButton
import my.robi.boom.ui.BoomTextField
import androidx.compose.material3.HorizontalDivider

@Composable
fun Magnet() {
    var income by rememberSaveable { mutableStateOf("") }
    var magnet by rememberSaveable { mutableStateOf("") }

    Column(
        modifier = Modifier
            .fillMaxWidth()
            .padding(16.dp),
        verticalArrangement = Arrangement.spacedBy(16.dp)
    ) {
        Text("Магнит", style = MaterialTheme.typography.headlineMedium)

        // Поле "Доход" + кнопка
        BoomTextField(
            value = income,
            onValueChange = { income = it },
            label = "Доход",
            modifier = Modifier.fillMaxWidth()
        )
        BoomButton(
            text = "Сохранить",
            onClick = { /* сохранить income */ },
            enabled = income.isNotBlank(),

        )

        HorizontalDivider()

        // Поле "Магнит" + кнопка
        BoomTextField(
            value = magnet,
            onValueChange = { magnet = it },
            label = "Магнит",
            modifier = Modifier.fillMaxWidth()
        )
        BoomButton(
            text = "Сохранить",
            onClick = { /* сохранить magnet */ },
            enabled = magnet.isNotBlank(),

        )
    }
}
