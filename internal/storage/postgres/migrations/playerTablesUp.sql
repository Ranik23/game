-- Таблица Игроков
CREATE TABLE Players (
    PlayerID SERIAL PRIMARY KEY, -- Уникальный идентификатор игрока
    PlayerName VARCHAR(255), -- Имя игрока
    PlayerEmail VARCHAR(255) -- Электронная почта игрока
);

-- Реестр счетов Игрока
CREATE TABLE PlayerAccounts (
    ID SERIAL PRIMARY KEY, -- Уникальный идентификатор записи счёта
    PlayerID INTEGER REFERENCES Players(PlayerID), -- Ссылка на игрока
    AccountArticle VARCHAR(255), -- Статья счёта
    AccountSum DECIMAL(10, 2), -- Сумма счёта
    Comments TEXT, -- Комментарии
    Reserve VARCHAR(255) -- Резервное поле
);

-- Журнал учета затрат/поступлений Игрока
CREATE TABLE PlayerExpensesIncomes (
    ID SERIAL PRIMARY KEY, -- Уникальный идентификатор записи
    PlayerID INTEGER REFERENCES Players(PlayerID), -- Ссылка на игрока
    Classifier VARCHAR(255), -- Классификатор
    Amount DECIMAL(10, 2), -- Сумма (положительная или отрицательная)
    Comments TEXT, -- Комментарии
    Reserve VARCHAR(255), -- Резервное поле
    FixationDate DATE -- Дата фиксации
);

-- Реестр ресурсов Игрока
CREATE TABLE PlayerResources (
    ResourceID SERIAL PRIMARY KEY, -- Уникальный идентификатор ресурса
    PlayerID INTEGER REFERENCES Players(PlayerID), -- Ссылка на игрока
    ResourceName VARCHAR(255), -- Название ресурса
    UnitOfMeasure VARCHAR(50), -- Единица измерения ресурса
    InitialQuantity DECIMAL(10, 2), -- Начальное количество ресурса
    Comments TEXT, -- Комментарии
    Reserve VARCHAR(255) -- Резервное поле
);

-- Журнал учета использования Ресурсов Игроком
CREATE TABLE PlayerResourceUsage (
    ID SERIAL PRIMARY KEY, -- Уникальный идентификатор записи использования ресурса
    PlayerID INTEGER REFERENCES Players(PlayerID), -- Ссылка на игрока
    ResourceID INTEGER REFERENCES PlayerResources(ResourceID), -- Ссылка на ресурс
    Events VARCHAR(255), -- Событие, связанное с использованием ресурса
    Comments TEXT, -- Комментарии
    Reserve VARCHAR(255), -- Резервное поле
    EntryNumber INTEGER, -- Учетный номер записи
    FixationDate DATE -- Дата фиксации события
);

-- Реестр кредитов Игрока
CREATE TABLE PlayerCredits (
    ID SERIAL PRIMARY KEY, -- Уникальный идентификатор записи кредита
    FixationDate DATE, -- Дата фиксации
    PlayerID INTEGER REFERENCES Players(PlayerID), -- Ссылка на игрока
    Comments TEXT, -- Комментарии
    Reserve VARCHAR(255), -- Резервное поле
    CreditDirection VARCHAR(255), -- Направление кредита (получен/выдан)
    CreditDuration1 VARCHAR(50), -- Срок кредита (краткосрочный/долгосрочный)
    CreditDuration2 VARCHAR(50), -- Срок кредита (месяцы)
    Percent INTEGER,
    PeriodOfPayment VARCHAR(50) -- периодичность платежа (мес / кв)
);
-- почему нет суммы, а только процент?