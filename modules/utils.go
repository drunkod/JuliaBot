package modules

import (
	"strconv"
	"time"
	"fmt"
	"os"	
)

func parseBirthday(dat, month, year int32) string {
	months := []string{
		"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December",
	}
	result := strconv.Itoa(int(dat)) + ", " + months[month-1]
	if year != 0 {
		result += ", " + string(year)
	}

	return result + "; is in " + tillDate(dat, month)
}

func tillDate(dat, month int32) string {
	currYear := time.Now().Year()

	timeBday := time.Date(currYear, time.Month(month), int(dat), 0, 0, 0, 0, time.UTC)
	currTime := time.Now()

	if timeBday.Before(currTime) {
		timeBday = time.Date(currYear+1, time.Month(month), int(dat), 0, 0, 0, 0, time.UTC)
	}

	// convert to days only
	days := timeBday.Sub(currTime).Hours() / 24

	return strconv.Itoa(int(days)) + " days"
}


func AskInputOrEnv[T any](key string) T {
	var val T
	var envValue string

	if v, ok := os.LookupEnv(key); ok {
		envValue = v
	}

	switch any(val).(type) {
	case string:
		if envValue != "" {
			return any(envValue).(T)
		}
	case int:
		if envValue != "" {
			i, err := strconv.Atoi(envValue)
			if err == nil {
				return any(i).(T)
			}
		}
	case int32:
		if envValue != "" {
			i, err := strconv.Atoi(envValue)
			if err == nil {
				return any(int32(i)).(T)
			}
		}
	}

	fmt.Printf("Enter %s: ", key)
	var input string
	fmt.Scanln(&input)

	switch any(val).(type) {
	case string:
		return any(input).(T)
	case int:
		i, err := strconv.Atoi(input)
		if err == nil {
			return any(i).(T)
		}
	}

	return val
}