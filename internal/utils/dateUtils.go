package utils

import (
	"fmt"
	"time"
)

func getSpanishMonth(month int) string {
	months := map[int]string{
		1:  "Enero",
		2:  "Febrero",
		3:  "Marzo",
		4:  "Abril",
		5:  "Mayo",
		6:  "Junio",
		7:  "Julio",
		8:  "Agosto",
		9:  "Septiembre",
		10: "Octubre",
		11: "Noviembre",
		12: "Diciembre",
	}

	if name, ok := months[month]; ok {
		return name
	}
	return "Unknown month"
}

func ParseBirthday(date time.Time) string {
	return fmt.Sprintf("%d de %s", date.Day(), getSpanishMonth(int(date.Month())))

}
