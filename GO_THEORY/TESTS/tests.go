Тестирование
В этой теме вы познакомитесь с базовыми возможностями для тестирования в экосистеме Go, а также поработаете с самой популярной библиотекой для тестов — testify.
Вы узнаете:
что такое юнит-тестирование;
как создавать и запускать тесты для Go-пакетов;
как определять процент покрытия кода тестами;
какие существуют паттерны тестирования;
зачем нужны и как использовать моки (заглушки).
Юнит-тесты и покрытие кода
В этом уроке расскажем о том, как тестировать Go-код стандартными утилитами.
Существует множество видов функциональных и нефункциональных тестов:
Юнит-тесты — тестируют минимальную часть функциональности (функцию или методы) в полной изоляции от внешних зависимостей. По сути, тестируются отдельные небольшие кусочки кода.
Интеграционные тесты — тестируют взаимодействие нескольких крупных частей приложения, например систем оформления заказов и оплаты.
End-to-end-тесты — тестируют работоспособность всей системы.
Мутационные тесты — тестируют код на устойчивость к случайным изменениям.
Нагрузочные тесты — используются для определения максимальной нагрузки, которую система способна выдержать с допустимым уровнем деградации.
Напомним, что функциональные тесты проверяют работоспособность разработанного кода, то есть соответствие функциональным требованиям. Нефункциональные (например, нагрузочные) определяют соответствие кода требованиям надёжности, качества, сопровождения и т. п.
Сконцентрируемся на юнит-тестах, потому что:
Это самый простой для анализа вид тестов.
При правильном подходе к тестированию таких тестов в кодовой базе будет большинство.
Техники и приёмы, которые будут продемонстрированы на их примере, можно использовать и в других типах функциональных тестов.
Инструментарий языка Go, поставляемый вместе с компилятором, включает в себя обширные средства тестирования. Будем рассматривать именно их.
Где размещаются юнит-тесты в Go?
В Go все тесты должны располагаться в файлах с суффиксом _test.go: например, user_test.go. Юнит-тесты принято располагать рядом с тестируемым кодом. Файлы user.go и user_test.go обычно лежат в одной и той же директории.
Файлы *_test.go не участвуют в компиляции финальной сборки проекта, поэтому можно не бояться импортировать в них большие библиотеки вроде stretchr/testify.
Тем не менее при компиляции тестов, как и основного кода, запрещены циклические импорты. Однако при написании тестов они могут возникать часто, поскольку в тестах может понадобиться код, который зависит от тестируемого кода. В таком случае тесты будут зависеть от тестируемого кода и, следовательно, импортировать сами себя. Поэтому для тестов сделали единственное исключение из правила «одна директория — один пакет». Тестовые файлы могут располагаться в пакете с суффиксом _test — и этой возможностью лучше пользоваться.
Также в пакете _test стоит располагать код, нужный исключительно для тестов. Допустим, надо сделать приватный тип публичным или добавить класс вспомогательных методов. Вспомогательный файл (часто его называют harness_test.go или common_test) может выглядеть так:
package user

import "context"

type UserDAO = userDAO

func (p *userProvider) ResetAllCaches(ctx context.Context) error {
    // сбрасываем кеши провайдера
} 
Тесты в Go
Теперь, когда вы знаете, где и как хранить тестовый код, поговорим о том, как именно писать тесты на Go.
В Go все тесты — это функции вида:
func TestXxx(t *testing.T) 
Префикс Test обязателен. В качестве Xxx обычно указывают название тестируемой функции. У каждой тестируемой функции может быть несколько тестов, и тогда нужно указать дополнительную информацию по конкретному тесту.
Для примера протестируем функцию Add, которая должна сложить два числа при условии, что они положительные. Если одно или оба числа равны нулю, функция должна вернуть ошибку.
Файл add.go:
package math

import "errors"

func Add(a, b int) (int, error) {
    if a == 0 || b ==  0 {
        return 0, errors.New("arg is zero")
    }
    
    if a < 0 || b < 0 {
        return 0, errors.New("arg is negative")
    } 
    return a + b, nil
} 
Файл add_test.go:
package math

import "testing"

func TestAddPositive(t *testing.T) {
    sum, err := Add(1, 2)
    if err != nil {
        t.Error("unexpected error")
    }
    if sum != 3 {
        t.Errorf("sum expected to be 3; got %d", sum)
    }
}


func TestAddNegative(t *testing.T) {
    _, err := Add(-1, 2)
    if err == nil {
        t.Error("first arg negative - expected error not be nil" )
    }
    _, err = Add(1, -2)
    if err == nil {
        t.Error("second arg negative - expected error not be nil" )
    }
    _, err = Add(-1, -2)
    if err == nil {
        t.Error("all arg negative - expected error not be nil" )
    }
}
 
Объект *testing.T предоставляет доступ к нескольким базовым методам:
Error, Errorf — записывает сообщение в error-лог и помечает тест как непройденный. Исполнение теста продолжается.
Fatal, Fatalf — делает то же самое, но исполнение теста немедленно завершается. Этот метод часто используется в рабочих проектах при обработке ошибок. Очень удобен при отладке, когда тестируется какой-то конкретный участок кода.
Skip, Skipf — позволяет пропустить тест с сообщением. Используется, когда окружение для теста не задано. Типичный сценарий — прогон интеграционных тестов с внешним сервисом только на CI, где к нему есть доступы.
Log, Logf — позволяет выводить лог-сообщения внутри теста. Преимущество перед методами пакета fmt в том, что из лога сразу видно, к какому тесту относится сообщение.
Run(name string, testf func(t *testing.T) ) — запускает функцию в качестве теста, что удобно при выполнении нескольких запусков теста, например, с разными именами.
go test
Теперь разберём запуск написанных тестов. Для этого в экосистеме Go используется стандартная утилита go test. Она позволяет запускать тесты следующими способами:
Тесты на основании положения кода в директории. Чтобы запустить все тесты в директории, достаточно перейти в неё и выполнить команду go test или go test -v. Флаг -v перенаправляет на stdout всё, что тесты логируют в stdout и stderr.
Все тесты пакета. Чтобы запустить все тесты в пакете, утилите go test надо передать пути импорта этих пакетов, разделённые пробелами. Например: go test math github.com/username/packagename github.com/username/packagename2.
Тесты, подходящие под регулярное выражение. Также есть возможность протестировать некоторое подмножество тестов пакета. Для этого используется флаг -run утилиты go test.
Например, если нужно протестировать все тест-кейсы с префиксом TestFunc в пакете github.com/ytuser/ytpackage, вызов команды go test будет выглядеть так:
go test github.com/ytuser/ytpackage -run ^TestFunc 
В качестве аргумента передаётся регулярное выражение, под которое должны подходить названия тестов. Вот только регулярные выражения выходят за пределы темы урока, и если вы не знакомы с ними, то лучше сперва прочитать эту статью.
Кеширование тестов
Если повторно запустим команду go test <PACKAGE_NAME>, то увидим, что вывод команды изменился. В случае с пакетом math получим:
ok      math    (cached) 
Дело в том, что в режиме тестирования пакета go test кеширует результат прогона тестов и, если код и тесты не были изменены, использует закешированный результат.
Отключить кеширование можно двумя способами:
Передать флаг -count 1, который определяет, сколько раз нужно запустить каждый тест (по умолчанию — один). Соответственно, -count 1 не изменяет количество запусков — если сравнивать со значением по умолчанию, — но выключает кеширование.
Запустить команду go test clear, очищающую кеш.
Дополнительные настройки тестирования
-cpu 1,2,4 — позволяет прогнать все тесты несколько раз с использованием разного количества потоков. Пригодится, если нужно протестировать параллельный код и убедиться, что на машинах с разным количеством ядер он будет работать корректно.
-list regexp — вместо того чтобы запускать тесты, go test выведет в консоль имена тестов, подходящих под переданное регулярное выражение.
-parallel n — позволяет параллельно выполнять тесты, которые в теле вызывают t.Parallel.
-run regexp — позволяет запускать конкретные тесты.
-short — если этот флаг передан, то t.Short() == true. В этом случае можно либо пропустить длительные тесты, либо урезать их функционал.
-v — подробное логирование. Даже в случае успешного прохождения тестов весь их лог будет выведен в консоль.

Покрытие кода
Одна из важнейших метрик качества кода — степень покрытия тестами (test coverage). В Go для вычисления этой метрики используется флаг -cover утилиты go test, подробнее о котором можно почитать в официальном блоге Go.
Например, если вы хотите узнать степень покрытия тестами пакета math из стандартной библиотеки, надо вызвать команду go test math -cover:
ok      math    0.003s    coverage: 86.8% of statements 
Определим покрытие тестами функции Add из примера выше:
% go test -cover
PASS
coverage: 80.0% of statements
ok      tests_06        0.233s 
Но знания одной только метрики зачастую не хватает, и нужно выяснить, какие именно строки кода не были задействованы при прогоне тестов. Этот функционал тоже идёт «из коробки» в виде флага -coverprofile утилиты go test и специальной утилиты go cover для анализа профиля покрытия тестами.
Итак, если вы хотите узнать, какой именно код пакета math не был покрыт тестами, надо сделать следующее:
Запустить на нём утилиту go test и сохранить файл профиля покрытия тестами. Путь к файлу с профилем — это значение флага -coverprofile. В данном случае сохраним его в файл coverage.out текущей директории:
go test . -coverprofile=coverage.out 
Проанализировать полученный файл утилитой cover. Например, по собранному профилю можно получить HTML-представление исходного кода с дополнительной разметкой, связанной с покрытием тестами:
go tool cover -html=coverage.out 
После выполнения автоматически запустится браузер, где будет отображена информация по покрытию.
image
Задание 3
Допишите тесты в math_test, чтобы добиться стопроцентного покрытия кода. Как будете готовы проверить себя, нажмите на кнопку ниже:


Правильный ответ
Кнопка ниже
func TestAddZero(t *testing.T) {
    _, err := Add(0, 2)
    if err == nil {
        t.Error("first arg is zero  - expected error not be nil")
    }
    _, err = Add(1, 0)
    if err == nil {
        t.Error("second arg is zero  - expected error not be nil")
    }
    _, err = Add(0, 0)
    if err == nil {
        t.Error("all arg negative - expected error not be nil")
    }
} 
Удобное тестирование
Можно тестировать код, используя только *testing.T, но он не предоставляет доступ к функциям вроде проверки на равенство, проверки на возврат ошибки или проверки на наличие паники при вызове переданного колбэка.
Этот пробел заполняют сторонние библиотеки, среди которых самая популярная — testify. Следует отметить, что они не являются заменой стандартного testing, а расширяют и дополняют его.
testify — это швейцарский нож в мире тестирования Go-кода. Поэтому поговорим про часто используемые пакеты из этого репозитория отдельно.
assert{target="_blank"} — пакет, в котором собрано множество удобных функций с говорящими названиями: вроде assert.Equal или assert.Nil. Предназначен для использования внутри обычных Go-тестов вида func TestXxx(t *testing.T) (про необычные поговорим, когда доберёмся до suite).
require — то же самое, что assert, но в случае, если проверки из этого пакета падают, выполнение теста останавливается. То есть внутри используется Fatal вместо Error.
suite — пакет, который вводит концепцию тест-сьюта (test suite). Если вы работали с тестами на Java или Python, то, скорее всего, она вам знакома.
Сьют — это объект, содержащий сами тесты в виде методов, а также некоторый набор переменных, доступный всем тестам. Кроме того, сьют даёт возможность определить код, который будет вызван перед и/или после исполнения теста. Эта функция полезна, например, для очистки базы данных между тестами (в случае интеграционных тестов).
Типичный пример работы со suite из официальной документации выглядит так:
import (
    "testing"
    "github.com/stretchr/testify/suite"
)

// ExampleTestSuite — это тестовый сьют, который создан путём эмбеддинга suite.Suite.
type ExampleTestSuite struct {
    suite.Suite
    VariableThatShouldStartAtFive int
}

// SetupTest заполняет переменную VariableThatShouldStartAtFive перед началом теста.
func (suite *ExampleTestSuite) SetupTest() {
    suite.VariableThatShouldStartAtFive = 5
}

func (suite *ExampleTestSuite) TestExample() { // все тесты должны начинаться со слова Test
    suite.Equal(5, suite.VariableThatShouldStartAtFive)
}

func TestExampleTestSuite(t *testing.T) {
    // чтобы go test смог запустить сьют, нужно создать обычную тестовую функцию
    // и вызвать в ней suite.Run
    suite.Run(t, new(ExampleTestSuite))
} 
Паттерны тестирования
Есть несколько рекомендаций, как стоит и не стоит писать тесты. Они не специфичны для Go, поэтому их можно применить для тестирования на любом языке.
Каждый тест должен тестировать что-то одно
Например, если вы тестируете функцию func Divide(a, b int) (int, error), не стоит писать такой код:
func TestDivision(t *testing.T) {
    result, err := Divide(0, 1)
    require.NoError(t, err)
    assert.Equal(t, 0, result)

    result, err = Divide(4, 2)
    require.NoError(t, err)
    assert.Equal(t, 2, result)

    _, err = Divide(1, 0)
    require.Error(err)
} 
Он будет падать при любой проблеме в тестируемой функции. Лучше разбить его на три теста (или три подтеста внутри этого теста), каждый из которых тестирует свой специфический сценарий:
func TestDivision(t *testing.T) {
    t.Run("ZeroNumerator", func(t *testing.T) {
        result, err := Divide(0, 1)
        require.NoError(t, err)
        assert.Equal(t, 0, result)
    })

    t.Run("BothNonZero", func(t *testing.T) {
        result, err = Divide(4, 2)
        require.NoError(t, err)
        assert.Equal(t, 2, result)
    })

    t.Run("ZeroDenominator", func(t *testing.T) {
        _, err = Divide(1, 0)
        require.Error(err)
    })
} 
Тесты не должны зависеть друг от друга
Если один тест опирается на глобальное состояние, устанавливаемое другими тестами, — это проблема. Тестовая архитектура такого типа приводит к неприятным ошибкам, когда локально тесты проходят, а на CI/CD иногда падают, потому что время от времени порядок выполнения тестов меняется.
Эта же рекомендация относится к тестам внутри одного сьюта. При написании каждого теста нужно исходить из того, что ему на вход передаётся состояние после вызова подготовительного метода SetupTest.
Результат работы теста — это не лог
При написании теста стоит считать, что при успешном прохождении теста никто на его логи смотреть не будет. Все контракты, которые фиксирует тест, должны быть прописаны в виде проверок. В этом случае можно безопасно гонять тесты на CI/CD без участия человека.
Table-driven-тесты
В примерах выше мы видели, что тесты содержат много повторяющегося кода, и это неспроста. Существует паттерн тестирования, который называется table-driven, или, говоря по-русски, «табличное тестирование». Он реализуется не только в Go, но и в других языках программирования.
Суть паттерна в отделении тестовых данных от выполнения самих тестов. Для коротких тестов паттерн может показаться избыточным, однако он позволяет эффективно организовать и расширять тесты, если тестовых наборов требуется несколько.
В Go есть пакет генерации шаблонов тестов именно в стиле table-driven, поэтому для рассмотрения примера будем использовать генерацию кода.
Если вы используете GoLand или VS Code с плагинами, то инструменты генерации уже встроены в IDE. Если нет, то потребуется немного дополнительных действий.
Установите пакет github.com/cweill/gotests:
$ go get -u github.com/cweill/gotests/... 
Вернёмся к нашей функции Add. Сгенерируем для неё заготовку теста:
gotests -only Add . 
В файл math_test.go добавился следующий код. Мы добавили к нему комментарии, чтобы было понятнее, что он делает.
func TestAdd(t *testing.T) {
    // args — описывает аргументы тестируемой функции
    type args struct {
        a int
        b int
    }
    // описывает структуру тестовых данных и сами тесты
    tests := []struct {
        name    string // название теста
        args    args // аргументы 
        want    int // ожидаемое значение
        wantErr bool // должна ли функция вернуть ошибку 
    }{
        // TODO: Add test cases.
        // сюда добавляем тестовые случаи (testcases)
    }
    // вызываем тестируемую функцию для каждого тестового случая  
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Add(tt.args.a, tt.args.b)
            if (err != nil) != tt.wantErr {
                t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("Add() = %v, want %v", got, tt.want)
            }
        })
    }
}
 
Добавим тестовые случаи в место, указанное в //TODO:
        {
            name: "Test Positive",
            args: args{
                a: 1,
                b: 2,
            },
            want:    3,
            wantErr: false,
        },
        {
            name: "Test Negative 1",
            args: args{
                a: -1,
                b: 2,
            },
            want:    0,
            wantErr: true,
        },
        {
            name: "Test Negative 2",
            args: args{
                a: 1,
                b: -2,
            },
            want:    0,
            wantErr: true,
        },

        {
            name: "Test Negative all",
            args: args{
                a: -1,
                b: -2,
            },
            want:    0,
            wantErr: true,
        }, 
И запустим тесты.
Удобство такого подхода в том, что не требуется добавлять вызов и проверки, чтобы добавить ещё несколько случаев. Достаточно описать аргументы и желаемое поведение.
Задание
Напишите набор тестов для функции EstimateValue(value int) string:
func EstimateValue(value int) string {
    switch {
    case value < 10:
        return "small"

    case value < 100:
        return "medium"

    default:
        return "big"
    }
} 
Постарайтесь добиться 100%-го покрытия функции тестами. Рекомендуем использовать библиотеку testify.
Готовы проверить себя?


Правильный ответ
Да
package main

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestEstimateValue(t *testing.T) {
    t.Run("Small", func(t *testing.T) {
        assert.Equal(t, "small", EstimateValue(9))
    })

    t.Run("Medium", func(t *testing.T) {
        assert.Equal(t, "medium", EstimateValue(99))
    })

    t.Run("Big", func(t *testing.T) {
        assert.Equal(t, "big", EstimateValue(100))
    })
} 
Бонусное задание
Напишите тест для функции EstimateValue в стиле table-driven.
Готовы проверить себя?


Правильный ответ
Готово
package main

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestEstimateValueTableDriven(t *testing.T) {
    testCases := []struct {
        Name          string
        InputValue    int
        ExpectedValue string
    }{
        {
            Name:          "Small",
            InputValue:    9,
            ExpectedValue: "small",
        },
        {
            Name:          "Medium",
            InputValue:    99,
            ExpectedValue: "medium",
        },
        {
            Name:          "Big",
            InputValue:    100,
            ExpectedValue: "big",
        },
    }

    for _, tc := range testCases {
        t.Run(tc.Name, func(t *testing.T) {
            assert.EqualValues(t, tc.ExpectedValue, EstimateValue(tc.InputValue))
        })
    }
} 
Ключевые мысли
В Go инструменты тестирования обеспечивают всё необходимое для проведения юнит-тестирования.

Код тестов определяется компилятором и не компилируется в конечную сборку.

Необходимые инструменты можно получить из сторонних библиотек. Особенно полезны testify, suite, gotests.

Table-driven test — удобный способ организации тестов, который часто встречается на практике.


Интерфейсы в тестировании
В предыдущих уроках вы изучили различные подходы и методики в тестировании программного продукта. В этом вы узнаете, как работает тестирование с использованием моков (mock), или заглушек, и зачем оно нужно.
При работе с моками создаётся объект, который доступен извне так же, как настоящий, но разработчик полностью контролирует его поведение. Именно этот объект и называется mock.
В языке Go мок-тестирование особенно удобно благодаря концепции интерфейсов. По сути, всё, что нужно для создания объекта-заглушки, — это удовлетворить интерфейсу реального объекта. В других ООП-языках — вроде Python или Java — мок-тестирование также существует, однако может быть сложнее в случае сложной иерархии наследования.
Это бывает полезно, когда:
нужно протестировать только работу бизнес-логики;
процессы занимают много времени и его можно сэкономить при тестировании;
нельзя или нежелательно при тестировании выполнять какую-то операцию, например отправку email или уведомлений;
невозможно развернуть копию БД или она представляет собой чёрный ящик;
сложно тестировать необходимые состояния во внешних источниках данных и проще установить нужные граничные условия на моках.
image
Проще всего понять принцип работы с моками на примере.
Предположим, есть БД, которую нельзя использовать для тестирования, но надо проверить, правильно ли работает написанный код для работы с ней. Возьмём простейший случай: пакет для работы с БД имеет тип DB и метод для проверки существования пользователя по его email. Метод UserExists возвращает true, если пользователь с указанным адресом существует, и false, если нет.
func (db *DB) UserExists(email string) bool 
В прошлых темах рассматривалось понятие интерфейсного типа, в котором описывается только поведение (методы) какого-то объекта. При этом структура и его внутренняя реализация не имеют значения — можно описать набор методов для работы с БД в виде интерфейса.
В продакшене будем использовать тип, который подключается к базе данных и отправляет запросы, а для тестирования создадим тип с такими же методами, при вызове которых сможем сравнить результаты с эталонными значениями.
type DBStorage interface {
    UserExists(email string) bool
}


// обратите внимание, что DBStorage передаётся в функцию в качестве параметра, таким образом мы можем при тестировании подменить реальную БД тестовой заглушкой.
func NewUser(db DBStorage, email string) error {
    if db.UserExists(email) {
        return fmt.Errorf(`user with '%s' email already exists`, email)
    }
    // добавляем запись
    return nil
} 
Здесь определены интерфейсный тип DBStorage и функция NewUser, в которой происходит проверка на существование пользователя с таким же почтовым ящиком. В продакшене эту функцию будем вызывать для переменной типа DB, а сейчас напишем для неё тест.
Если есть много вариантов для проверки, то для тестирования лучше использовать таблицы (table-driven tests) с входящими данными и ожидаемыми результатами:
import (
    "github.com/stretchr/testify/require"
)

// тип объекта-заглушки
type DBMock struct {
    emails map[string]bool
}

// для удовлетворения интерфейсу DBStorage реализуем  
func (db *DBMock) UserExists(email string) bool {
    return db.emails[email]
}
// вспомогательный метод, для подсовывания тестовых данных
func (db *DBMock) addUser(email string) {
    db.emails[email] = true
}

func TestNewUser(t *testing.T) {
    errPattern := `user with '%s' email already exists`
    tbl := []struct {
        name    string
        email   string
        preset  bool
        wanterr bool
    }{
        {`want success`, `gregorysmith@myexampledomain.com`, false, false},
        {`want error`, `johndoe@myexampledomain.com`, true, true},
    }
    for _, item := range tbl {
        t.Run(item.name, func(t *testing.T) {
            // создаём объект-заглушку 
            dbMock := &DBMock{emails: make(map[string]bool)}
            if item.preset {
                dbMock.addUser(item.email)
            }
             // выполняем наш код, передавая объект-заглушку
            err := NewUser(dbMock, item.email)
            if !item.wanterr {
                require.NoError(err)
            } else {
                require.EqualError(t, err, fmt.Sprintf(errPattern, err.email))
            }
        })
    }
} 
В тестовой функции TestNewUser проверяем, какой результат возвращает функция NewUser в зависимости от того, существует пользователь с таким email или нет. Таким образом можно протестировать поведение своей функции без обращения к реальной базе данных. В примере на каждую итерацию создаём мок и при необходимости добавляем туда email, чтобы проверить различные варианты.
Моки можно использовать не только для подмены операций, но и для получения детальной информации, например количества вызовов функции, контроля параметров и т. д. Для этого достаточно добавить нужные поля в структуру заглушки и изменять их при каждом вызове интерфейсного метода:
type DBMock struct {
    emails  map[string]bool
    counter int
}

func (db *DBMock) UserExists(email string) bool {
    db.counter++
    return db.emails[email]
} 
Итак, мы разобрались с общими принципами работы моков, которые дают дополнительные возможности для тестирования ПО. В простых случаях можно самостоятельно создавать подобные тесты, а на практике лучше использовать готовые библиотеки для тестирования с использованием моков:
testify/mock;
golang/mock;
vektra/mockery.
Дополнительные материалы
GoMock vs. Testify: Mocking frameworks for Go
Mock Solutions for Golang Unit Test
Задание
Определена API со следующим интерфейсом:
type APIClient interface {
    GetData(query string) (Response, error)
}

type Response struct {
    Text       string
    StatusCode int
} 
Реализуйте тип Mock с интерфейсом MockAPIClient:
type MockAPIClient interface {
    APIClient
    SetResponse(resp Response, err error)
} 


package main

type APIClient interface {
    GetData(query string) (Response, error)
}

type Response struct {
    Text       string
    StatusCode int
}

type MockAPIClient interface {
    APIClient
    SetResponse(resp Response, err error)
}

type Mock struct {
    Resp Response
    Err  error
}

func (m *Mock) GetData(query string) (Response, error) {
    return m.Resp, m.Err
}

func (m *Mock) SetResponse(resp Response, err error) {
    m.Resp = resp
    m.Err = err
}