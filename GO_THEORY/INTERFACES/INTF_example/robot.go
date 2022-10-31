package robot

import "fmt"

// Robot — тип робота
type Robot struct {
	model       string
	serialId    int
	workCounter int
}

func (r Robot) String() string {
	return fmt.Sprintf("Robot %s serialID %d", r.model, r.serialId)
}

// Work — робот выполняет работы и запоминает количество выполненных задач. Поэтому получатель метода — по указателю
func (r *Robot) Work(tasks []string) string {
	res := fmt.Sprintf("%s work:", r)
	for _, task := range tasks {
		res += "\n I do " + task
	}
	r.workCounter += len(tasks)
	return res
}

// Важно
// С точки зрения Go типы Robot и *Robot (указатель) — разные. В примере метод Work привязан именно к *Robot. Так как формально тип Robot не реализует интерфейс Worker, такой код не скомпилируется: go
// robo := Robot{};
// comp.Hire(robo);
// Поэтому будем использовать указатель на робота. Действительно, в этом есть логика. Раз работа в компании изменяет внутреннее состояние робота, то нужно передать указатель именно на неё.
//  robo := &Robot{};
//  comp.Hire(robo);
