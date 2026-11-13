package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

const menu = `===== STUDENT GRADE MANAGER =====
1. Добавить студента
2. Добавить оценку
3. Показать всех студентов
4. Найти студента
5. Изменить оценку
6. Статистика успеваемости
0. Выход`

func readText(scanner *bufio.Scanner, prompt string) (string, error) {
	fmt.Print(prompt)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", err
		}
		return "", io.EOF
	}
	return strings.TrimSpace(scanner.Text()), nil
}

func readNumber(scanner *bufio.Scanner, prompt string) (int, error) {
	text, err := readText(scanner, prompt)
	if err != nil {
		return 0, err
	}
	number, err := strconv.Atoi(text)
	if err != nil {
		return 0, errors.New("введите целое число")
	}
	return number, nil
}

func addFromConsole(scanner *bufio.Scanner, manager *Manager) error {
	name, err := readText(scanner, "ФИО: ")
	if err != nil {
		return err
	}
	if name == "" {
		return errors.New("имя не может быть пустым")
	}
	group, err := readText(scanner, "Группа: ")
	if err != nil {
		return err
	}
	id, err := manager.AddStudent(name, group)
	if err == nil {
		fmt.Printf("Студент добавлен. ID: %d\n", id)
	}
	return err
}

func findFromConsole(scanner *bufio.Scanner, manager *Manager) (*Student, error) {
	id, err := readNumber(scanner, "ID студента: ")
	if err != nil {
		return nil, err
	}
	return manager.FindStudent(id)
}

func gradeFromConsole(scanner *bufio.Scanner, manager *Manager, update bool) error {
	student, err := findFromConsole(scanner, manager)
	if err != nil {
		return err
	}
	index := 0
	if update {
		if len(student.Grades) == 0 {
			return errors.New("у студента ещё нет оценок для изменения")
		}
		for i, grade := range student.Grades {
			fmt.Printf("%d. %d\n", i+1, grade)
		}
		number, err := readNumber(scanner, "Номер оценки (с 1): ")
		if err != nil {
			return err
		}
		if number < 1 || number > len(student.Grades) {
			return errors.New("оценки с таким порядковым номером нет")
		}
		index = number - 1
	}
	grade, err := readNumber(scanner, "Оценка (0–100): ")
	if err != nil {
		return err
	}
	if update {
		err = UpdateGrade(student, index, grade)
	} else {
		err = AddGrade(student, grade)
	}
	if err == nil {
		fmt.Println("Оценка сохранена.")
		printStudent(student)
	}
	return err
}

func printStudent(student *Student) {
	fmt.Printf("ID: %d\nИмя: %s\nГруппа: %s\nОценки: %v\n",
		student.ID, student.Name, student.Group, student.Grades)
	if len(student.Grades) == 0 {
		fmt.Println("Средний балл: нет оценок\nУровень: не определён")
	} else {
		average := CalculateAverage(student.Grades)
		fmt.Printf("Средний балл: %.2f\nУровень: %s\n", average, PerformanceLevel(average))
	}
}

func printStatistics(manager *Manager) {
	average, best, lowCount, assessed := manager.Statistics()
	fmt.Printf("Всего студентов: %d\nБез оценок: %d\n", len(manager.students), len(manager.students)-assessed)
	if best == nil {
		fmt.Println("Средний балл: нет оценок\nЛучший студент: нет данных")
	} else {
		fmt.Printf("Средний балл студентов с оценками: %.2f\n", average)
		fmt.Printf("Лучший студент: %s (ID: %d; %.2f)\n", best.Name, best.ID, CalculateAverage(best.Grades))
	}
	fmt.Printf("Студентов со средним баллом ниже 50: %d\n", lowCount)
	counts := manager.CountByGroup()
	var groups []string
	for group := range counts {
		groups = append(groups, group)
	}
	sort.Strings(groups)
	fmt.Println("Количество студентов по академическим группам:")
	for _, group := range groups {
		fmt.Printf("%s: %d\n", group, counts[group])
	}
}

func main() {
	manager := DemoManager()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println(menu)
		choice, err := readText(scanner, "Выберите действие: ")
		if err != nil {
			if !errors.Is(err, io.EOF) {
				fmt.Fprintln(os.Stderr, "Ошибка чтения:", err)
			}
			return
		}
		switch choice {
		case "1":
			err = addFromConsole(scanner, &manager)
		case "2":
			err = gradeFromConsole(scanner, &manager, false)
		case "3":
			if len(manager.students) == 0 {
				fmt.Println("Студентов пока нет.")
			}
			for i := range manager.students {
				printStudent(&manager.students[i])
				fmt.Println()
			}
		case "4":
			var student *Student
			student, err = findFromConsole(scanner, &manager)
			if err == nil {
				printStudent(student)
			}
		case "5":
			err = gradeFromConsole(scanner, &manager, true)
		case "6":
			printStatistics(&manager)
		case "0":
			return
		default:
			err = errors.New("выберите пункт от 0 до 6")
		}
		if errors.Is(err, io.EOF) {
			return
		}
		if err != nil {
			fmt.Println("Ошибка:", err)
		}
		fmt.Println()
	}
}
