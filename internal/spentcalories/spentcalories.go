package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	// Разделила строку на слайс строк
	parts := strings.Split(data, ",")
	// Проверила, что длина слайса равна 3
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("некорректные данные")
	}
	// Преоьразовала количество шагов в int и обработала ошибку
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, err
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("некорректные данные")
	}
	// Преобразовала продолжительность в time.Duration и обработала ошибку
	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, err
	}
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("некорректные данные")
	}
	// Вернула кол-во шагов, вид активности, продолжительность и nil
	return steps, parts[1], duration, nil

}

func distance(steps int, height float64) float64 {
	// Считаю длину шага (рост умножаем на коэффициент)
	lenStep := height * stepLengthCoefficient
	// Считаю общую дистанцию в метрах (кол-во шагов умножаем на длину)
	distanceM := float64(steps) * lenStep
	// Перевела метры в километры и разделила на 1000
	distanceKm := distanceM / mInKm
	return distanceKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// Проверила что продолжительность больше 0, если нет, возвращаю 0
	if duration <= 0 {
		return 0
	}
	// Вызвала готовую distance
	distanceKm := distance(steps, height)
	// Вычислила и вернула среднюю скорость
	return distanceKm / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	// Получила значения и обработала ошибку
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}
	// Создала переменные для хранения дистанции, скорости и калорий
	var dist, speed, calories float64
	// В зависимости от типа тренировки, считаю дистанцию, скорость и калории
	switch activityType {
	case "Бег":
		dist = distance(steps, height)
		speed = meanSpeed(steps, height, duration)
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	case "Ходьба":
		dist = distance(steps, height)
		speed = meanSpeed(steps, height, duration)
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		// Если был передан неизвестный тип тренировки, вернуть ошибку с текстом
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}
	// Собрала строку с результатом
	result := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", activityType, duration.Hours(), dist, speed, calories)
	return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверила, что все входные данные больше нуля, если нет, возвращаю ошибку
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("некорректные данные")
	}
	// Вычислила среднюю скорость и сохранила е в speed
	speed := meanSpeed(steps, height, duration)
	// Перевела продолжительность в минуты
	durationInMinutes := duration.Minutes()
	// Расчет количества калорий по формуле
	calories := (weight * speed * durationInMinutes) / minInH
	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверила параметры и в случае ошибки вернула 0 и ошибку
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("некорректные данные")
	}
	// Рассчитала среднюю скорость
	speed := meanSpeed(steps, height, duration)
	// Перевела продолжительность в минуты
	durationInMinutes := duration.Minutes()
	// Рассчитала количество калорий по формуле из задания
	calories := (weight * speed * durationInMinutes * walkingCaloriesCoefficient) / minInH
	return calories, nil
}
