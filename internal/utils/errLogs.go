package utils

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"runtime"
	"strconv"
	"strings"
)

type ErrLogs struct {
	gorm.Model
	Error string `json:"error" db:"error"`
	Place string `json:"place" db:"place"`
	Count int    `json:"count" db:"count"`
}

func UppendErrorWithPath(err error, bd *gorm.DB) {
	badway := ""
	for stepCount := 1; stepCount <= 10; stepCount++ {
		// получаем через рантайм указатель на положение stepCount в цепочке шаг наверх, а так же строку вызова внутри функции.
		pc, _, line, _ := runtime.Caller(stepCount)
		// получаем через указатель полный путь к функции записанный на этом шагу
		fullFuncPath := runtime.FuncForPC(pc).Name()
		// дробим строку через точку, что бы получить массив строк, в котором последним элементом останется имя функции.
		splitedFuncPath := strings.Split(fullFuncPath, ".")
		// выбираем имя функции из массива в новую переменную, для лёгкого чтения дальнейших строк.
		funcName := splitedFuncPath[len(splitedFuncPath)-1]
		// если за stepCount шагов мы дошли до роута, то заканчиваем сбор имён в цепочке функций.
		if funcName == "ServeHTTP" {
			break
		} else {
			// исключаем пустые имена функций, что периодически могут появляться если шаги ушли достаточно далеко.
			if funcName != "" {
				// сохраняем/дописываем в crudLog.FuncPath имя и строку вызова функции.
				badway += funcName + "(" + strconv.Itoa(line) + ") | "
			}
		}
	}
	if err != nil {
		UppendError(badway, err.Error(), bd)
	} else {
		UppendError(badway, errors.New("передали пустую ошибку").Error(), bd)
	}
}

func UppendError(place string, newerr string, bd *gorm.DB) {
	var errLog ErrLogs
	log.Error("func uppendError( ", place, newerr)

	if err := bd.Where("place = ? and error = ?", place, newerr).First(&errLog).Error; err != nil {
		if err != nil {
			log.Warn(err)

		}
		errLog.Place = place
		errLog.Error = newerr
		errLog.Count = 1
		bd.Create(&errLog)
	} else {
		errLog.Count = errLog.Count + 1
		err = bd.Model(&errLog).Update("count", errLog.Count).Error
		if err != nil {

			log.Warn(err)
		}

	}
}
