package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",") // Разделила строку на части
	if len(parts) != 2 {              // Проверила, что длина слайса равна 2
		return 0, 0, fmt.Errorf("неверный формат данных")
	}
	// Преобразовала шаги из строки в int и обработала ошибку
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}
	// Проверила, что количество шагов больше нуля
	if steps <= 0 {
		return 0, 0, fmt.Errorf("количество шагов должно быть больше нуля")
	}
	// Преобразовала продолжительность
	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, err
	}
	if duration <= 0 {
		return 0, 0, fmt.Errorf("продолжительность должна быть больше нуля")
	}
	// Вернула результат
	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	// Вызвала функцию для получения шагов, продолжительность и ошибку.
	//
	//	Если ошибка есть, печатем ее и выходим из функции с пустой строкой
	parsedSteps, parsedDuration, err := parsePackage(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	// Проверяем,  что количество шагов было больше нуля
	if parsedSteps <= 0 {
		return ""
	}
	// Считаем дистанцию (сколько км прошли, переводим метры в километры разделив на 1000)
	distanceM := float64(parsedSteps) * stepLength
	distanceKm := distanceM / mInKm
	// Функция из другого пакета, она считает калории, пердает ей шаг, вес, рост и продолжительность
	// Снова проверяем ошибку
	calories, err := spentcalories.WalkingSpentCalories(parsedSteps, weight, height, parsedDuration)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	// Собрала строку с результатом
	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", parsedSteps, distanceKm, calories)

}
