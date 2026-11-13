package main

import (
	"errors"
	"fmt"
	"strings"
)

type Student struct {
	ID     int
	Name   string
	Group  string
	Grades []int
}

// Пять порогов используются при определении уровня по среднему баллу.
var gradeThresholds = [5]int{0, 50, 70, 85, 95}
var gradeLevels = [5]string{
	"Неудовлетворительно", "Удовлетворительно", "Хорошо", "Отлично", "Превосходно",
}

type Manager struct {
	students  []Student
	positions map[int]int
	lastID    int
}

func (m *Manager) AddStudent(name, group string) (int, error) {
	name, group = strings.TrimSpace(name), strings.TrimSpace(group)
	if name == "" || group == "" {
		return 0, errors.New("имя и группа не могут быть пустыми")
	}
	if m.positions == nil {
		m.positions = make(map[int]int)
	}
	m.lastID++
	m.students = append(m.students, Student{ID: m.lastID, Name: name, Group: group})
	m.positions[m.lastID] = len(m.students) - 1
	return m.lastID, nil
}

// Карта хранит индекс, поэтому поиск по ID не требует перебора всего среза.
// Полученный указатель используем сразу, до следующего добавления студента.
func (m *Manager) FindStudent(id int) (*Student, error) {
	position, ok := m.positions[id]
	if !ok {
		return nil, fmt.Errorf("студент с ID %d не найден", id)
	}
	return &m.students[position], nil
}

func CalculateAverage(grades []int) float64 {
	if len(grades) == 0 {
		return 0
	}
	var sum float64
	for _, grade := range grades {
		sum += float64(grade)
	}
	return sum / float64(len(grades))
}

func PerformanceLevel(average float64) string {
	level := gradeLevels[0]
	for i, threshold := range gradeThresholds {
		if average >= float64(threshold) {
			level = gradeLevels[i]
		}
	}
	return level
}

func AddGrade(student *Student, grade int) error {
	if student == nil {
		return errors.New("студент не найден")
	}
	if grade < 0 || grade > 100 {
		return errors.New("оценка должна быть целым числом от 0 до 100")
	}
	student.Grades = append(student.Grades, grade)
	return nil
}

// index — индекс с нуля; в меню пользователь вводит номер оценки с единицы.
func UpdateGrade(student *Student, index int, newGrade int) error {
	if student == nil {
		return errors.New("студент не найден")
	}
	if index < 0 || index >= len(student.Grades) {
		return errors.New("оценки с таким порядковым номером нет")
	}
	if newGrade < 0 || newGrade > 100 {
		return errors.New("оценка должна быть целым числом от 0 до 100")
	}
	student.Grades[index] = newGrade
	return nil
}

func (m *Manager) CountByGroup() map[string]int {
	counts := make(map[string]int)
	for _, student := range m.students {
		counts[student.Group]++
	}
	return counts
}

// Общий средний — среднее индивидуальных средних, а не всех оценок подряд.
// Студенты без оценок учитываются в списке и группах, но не в успеваемости.
func (m *Manager) Statistics() (average float64, best *Student, lowCount, assessed int) {
	var bestAverage float64
	for i := range m.students {
		student := &m.students[i]
		if len(student.Grades) == 0 {
			continue
		}
		current := CalculateAverage(student.Grades)
		average += current
		assessed++
		if current < float64(gradeThresholds[1]) {
			lowCount++
		}
		if best == nil || current > bestAverage {
			best, bestAverage = student, current
		}
	}
	if assessed > 0 {
		average /= float64(assessed)
	}
	return
}

func DemoManager() Manager {
	return Manager{
		students: []Student{
			{ID: 1, Name: "Иван Петров", Group: "CS-301", Grades: []int{85, 90, 78}},
			{ID: 2, Name: "Анна Сидорова", Group: "CS-302", Grades: []int{95, 92, 98}},
		},
		positions: map[int]int{1: 0, 2: 1},
		lastID:    2,
	}
}
