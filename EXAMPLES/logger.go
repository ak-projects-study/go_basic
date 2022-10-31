// Реализуйте для вашей структуры вспомогательный метод, который будет печатать сообщение в зависимости от уровня.
// func (l *LogExtended) println(srcLogLvl LogLevel, prefix, msg string) {
//     if l.logLevel < srcLogLvl {
//         return
//     }

//     l.Logger.Println(prefix + msg)
// }
// Задание 5
// Каждый уважающий себя разработчик должен написать свой логгер 😁.
// В стандартной библиотеке Go уже есть логгер — в пакете log. У этой реализации отсутствует параметр уровня лога (log level), то есть для вывода есть только метод log.Println и log.Printf. Предлагаем расширить этот объект следующими методами:
// SetLogLevel(logLvl LogLevel);
// Infoln(msg string);
// Warnln(msg string);
// Errorln(msg string).
// LogLevel — перечисление с такими возможными значениями:
// Info;
// Warning;
// Error.
// Логгер можно писать долго, прикручивая множество нужных фич, но стоит остановиться на простом расширении типа. Функции Infoln, Warnln и Errorln должны оборачивать метод log.Println, добавляя соответствующий префикс.
// Пример работы:
// func main() {
//     logger := NewLogExtended()
//     logger.SetLogLevel(LogLevelWarning)
//     logger.Infoln("Не должно напечататься")
//     logger.Warnln("Hello")
//     logger.Errorln("World")
//     logger.Println("Debug")
// }
// 2021/05/19 23:04:14 WARN Hello
// 2021/05/19 23:04:14 ERR World
// 2021/05/19 23:04:14 Debug
// Примечание
// Просто вызов log.Println() обращается к глобальной переменной std пакета log (загляните в реализацию). Но нужно создать свой объект типа LogExtended, то есть расширить тип log.Logger.
// Обратите внимание на конструктор типа log.Logger:
// func New(out io.Writer, prefix string, flag int) *Logger
// Конструктор возвращает указатель, поэтому тип должен выглядеть так:
// type LogExtended struct {
//     *log.Logger
//     logLevel LogLevel // LogLevel это enum
// }

package main

import (
	"log"
	"os"
)

type LogLevel int

func (l LogLevel) IsValid() bool {
	switch l {
	case LogLevelInfo, LogLevelWarning, LogLevelError:
		return true
	default:
		return false
	}
}

const (
	LogLevelError LogLevel = iota
	LogLevelWarning
	LogLevelInfo
)

type LogExtended struct {
	*log.Logger
	logLevel LogLevel
}

func (l *LogExtended) SetLogLevel(logLvl LogLevel) {
	if !logLvl.IsValid() {
		return
	}
	l.logLevel = logLvl
}

func (l *LogExtended) Infoln(msg string) {
	l.println(LogLevelInfo, "INFO ", msg)
}

func (l *LogExtended) Warnln(msg string) {
	l.println(LogLevelWarning, "WARN ", msg)
}

func (l *LogExtended) Errorln(msg string) {
	l.println(LogLevelError, "ERR ", msg)
}

func (l *LogExtended) println(srcLogLvl LogLevel, prefix, msg string) {
	if l.logLevel < srcLogLvl {
		return
	}

	l.Logger.Println(prefix + msg)
}

func NewLogExtended() *LogExtended {
	return &LogExtended{
		Logger:   log.New(os.Stderr, "", log.LstdFlags),
		logLevel: LogLevelError,
	}
}

func main() {
	logger := NewLogExtended()
	logger.SetLogLevel(LogLevelWarning)
	logger.Infoln("Не должно напечататься")
	logger.Warnln("Hello")
	logger.Errorln("World")
}
