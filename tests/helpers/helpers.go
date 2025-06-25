package helpers

import (
	"fmt"
	"os"
	"vovanDB/internal/helpers"
	"vovanDB/internal/database"
)

func CreateTestDataBase() error {
	sql := `
        CREATE TABLE users (
            id int AUTO_INCREMENT,
            name text NULL,
            age int,
            is_admin bool,
            date date
        );`

	executeResult := database.Execute(sql)

	if executeResult.Success != true {
		return fmt.Errorf("не удалось создать тестовую таблицу: %s", executeResult.Error)
	}

	return nil
}

func SeedTestDatabase() error {
	sql := `
        INSERT INTO users (id, name, age, is_admin, date)
        VALUES
        (1, 'vova', 38, true, '2025-01-28'),
        (2, 'katya', 33, false, '2025-01-28'),
        (3, 'sacha', 38, false, '2025-01-28');
    `

	executeResult := database.Execute(sql)

	if executeResult.Success != true {
		return fmt.Errorf("не удалось наполнить таблицу тестовыми данными: %s", executeResult.Error)
	}

	return nil
}

func ClearAllTestDatabase() {
	os.RemoveAll(helpers.GetDataBaseDir())
}
