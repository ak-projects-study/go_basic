package main

import (
	"fmt"
)

// Трёхкомпонентный цикл
// Классическая форма цикла состоит из трёх компонентов:
// i := 1 — инициализация (pre-действие): выполняется единожды при входе в scope цикла;
// i < 10 — основное условие: пока условие true, итерации будут продолжаться;
// i++ — post-действие: выполняется по завершении каждой итерации цикла.

// func main() {
// 	// создаём переменную
// 	v := 0
// 	for i := 1; i < 10; i++ {
// 		// наращиваем переменную
// 		v++
// 		fmt.Println(v)
// 	}
// 	// выводим результат на экран
// 	fmt.Println(v)
// }

// Заполнять каждую компоненту необязательно — можно опускать.
// Есть ещё два варианта бесконечного цикла, но в форме трёх компонент:
// for ;; {}
// for ; true; {}
// Компоненты цикла могут принимать более комплексный вид:
// for a, b := 5, 10; a < 10 && b < 20; a, b = a + 1, b + 2 {
//     // do stuff
// }

// Цикл while

func cycle_while() {
	// создаём переменную
	i := 0
	// описываем предусловие
	for i < 5 {
		// наращиваем переменную
		i++
	}
	// выводим результат на экран
	fmt.Println(i)
}

// Цикл while похож на трёхкомпонентный, но здесь оставлено только основное условие.

// Цикл range
// func main() {

// 	// создаём массив
// 	array := [3]int{1, 2, 3}
// 	// итерируемся
// 	for arrayIndex, arrayValue := range array {
// 		fmt.Printf("array[%d]: %d\n", arrayIndex, arrayValue)
// 	}
// }

// Цикл range используется для комплексных типов — слайса и мапы (map).
// Подробнее об этом цикле расскажем в следующей теме, посвящённой композитным типам.

// func main() {
// 	sum, limit := 0, 100
// 	for i := 0; true; i++ {
// 		if i%2 != 0 {
// 			continue // переход к следующему числу, так как i — нечётное
// 		}

// 		if sum+i > limit {
// 			break // выход из цикла, так как сумма превысит заданный предел
// 		}

// 		sum += i
// 	}
// 	fmt.Println(sum)
// }

// package main

// import "fmt"

// func main() {
// 	for i := 1; i <= 100; i++ {
// 		found := false

// 		if i%3 == 0 {
// 			fmt.Printf("Fizz")
// 			found = true
// 		}
// 		if i%5 == 0 {
// 			fmt.Printf("Buzz")
// 			found = true
// 		}

// 		if !found {
// 			fmt.Println(i)
// 			continue
// 		}

// 		fmt.Println()
// 	}
// }

func GetPersonWithLastVisisted(p Person) Person {
    return Person{
        Name:        p.Name,
        Age:         p.Age,
        lastVisited: time.Now(), // time.Now() возвращает текущее время
    }
}

p := P{
Name: "Alex",
Age: 25,
lastVisited: time.Time{} // пустое значение времени — пользователь ещё не посещал наш сервис
}

p = GetPersonWithLastVisisted(p)

