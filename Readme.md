**VovanDB** - файловая база данных, написанная на go.

Типы данных:

- text
- int
- bool
- datetime
- date

Подходит для работы с небольшими проектами, в которых нужны самые базовые операции по работе с данными.

Схема хранится в формате json, данные в csv.

Поддерживается CRUD-запросы (create, select, update, delete).

Не поддерживаются join, индексы, вложенные запросы, агрегатные функции, группировка, обновление схемы (alter).

Примеры запросов:

```sql
-- CREATE
CREATE TABLE users (
    id int AUTO_INCREMENT,
    name text NULL,
    age int,
    is_admin bool,
    date date
);

-- DROP
DROP TABLE users;

-- INSERT
INSERT INTO users 
    (id, name, age, is_admin, date)
VALUES
    (1, 'vova', 38, true, '2025-01-28'),
    (2, 'katay', 33, false, '2025-01-28'),
    (3, 'sacha', 38, false, '2025-01-28');

-- SELECT
SELECT 
    id, name
FROM
    users
WHERE
    is_admin = false
ORDER BY
    id DESC

-- UPDATE
UPDATE
    users
SET 
    is_admin = true,
    date = '2025-02-28'
WHERE
    is_admin = false

-- DELETE
DELETE FROM
    users
WHERE
    is_admin = false
```

Запуск тестов

```bach
go test ./tests -v
```

