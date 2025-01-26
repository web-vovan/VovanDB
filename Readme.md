**VovanDB** - файловая база данных, написанная на go.

Поддерживается:
* CREATE TABLE
* DROP TABLE
* INSERT
* SELECT
* UPDATE

Примеры запросов:

```sql
-- CREATE TABLE
CREATE TABLE users (
    id int AUTO_INCREMENT,
    name text NULL,
    age int,
    is_admin bool,
    date date
);

-- DROP TABLE
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

-- UPDATE
UPDATE
    users
SET 
    is_admin = true
WHERE
    is_admin = false
```

Нет поддержки:
* JOIN
* Индексы
* ALTER TABLE
* GROUP BY
* HAVING
* Агрегатные функции

