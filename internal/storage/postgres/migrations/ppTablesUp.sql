-- Обобщенная классификация
CREATE TABLE CommonClasification (
  ID SERIAL PRIMARY KEY,
  DemandInRealization INTEGER, -- потребность в реализации
  ProductOrService    INTEGER, -- 0 если не используется
  InvestDirection     TEXT -- направление инвестирования
);

-- Классификация ЖЦ проекта
CREATE TABLE ProjectLifeCycleClassification (
    ID SERIAL PRIMARY KEY, -- Уникальный идентификатор классификации
    RealizationForm VARCHAR(255), -- Форма реализации проекта
    TechnologyType VARCHAR(255), -- Тип технологии, применяемой в проекте
    DurationCategory VARCHAR(50), -- Категория продолжительности (долгосрочный, краткосрочный и т. д.)
    Costs DECIMAL(10, 2), -- Затраты
    Risks VARCHAR(255) -- Уровень риска
);

-- Классификация ЖЦ продукта/услуги
CREATE TABLE ProductLifeCycleClassification (
    ID SERIAL PRIMARY KEY, -- Уникальный идентификатор классификации
    Profitability VARCHAR(255), -- Доходность
    Risks VARCHAR(255), -- Риски
    ImpactConsequences TEXT, -- Последствия реализации или не реализации
    StrategicAgreement VARCHAR(255), -- Стратегическое согласование
    StrategicEfficiency VARCHAR(255) -- Стратегическая эффективность
);

-- Журнал Проектных Предложений Игрока
CREATE TABLE ProjectProposals (
    ProposalID SERIAL PRIMARY KEY, -- Уникальный идентификатор проектного предложения
    PlayerID INTEGER REFERENCES Players(PlayerID), -- Ссылка на игрока
    AppearanceInterval VARCHAR(50), -- Интервал появления проектного предложения (например, YYYY-MM)
    Status VARCHAR(255), -- Статус рассмотрения проектного предложения
    LifeCycleClassificationID INTEGER REFERENCES ProjectLifeCycleClassification(ID), -- Cсылка 
    ProductLifeCycleClassificationID INTEGER REFERENCES ProductLifeCycleClassification(ID), -- Ссылка
    CommonClasificationID INTEGER REFERENCES CommonClasification(ID)
);

-- Оценки фин. и ресурсных потребностей
CREATE TABLE FinancialResourceNeeds (
    ID SERIAL PRIMARY KEY, -- Уникальный идентификатор оценки
    ProposalID INTEGER REFERENCES ProjectProposals(ProposalID), -- Ссылка на проектное предложение
    RealizationCosts DECIMAL(10, 2), -- Оценка затрат на реализацию
    ResourceNeeds TEXT -- Оценка ресурсов, необходимых для реализации
);


CREATE TABLE ResultImpactOnKPC (
  ID SERIAL PRIMARY KEY,
  RealizedImpact INTEGER, -- оценка влияния реализованных результатов на КПС
  UnrealizedImpact INTEGER -- оценка влияния нереализованных результатов на КПС
);
-- Оценки ПП
CREATE TABLE ProjectRate (
  ID SERIAL PRIMARY KEY,
  FinancialResourceNeedsID INTEGER REFERENCES FinancialResourceNeeds(ID),
  ResultImpactOnKPCID INTEGER REFERENCES ResultImpactOnKPC(ID),
  DurationOfProjectRealization INTEGER,
  DurationOfProductUsing   INTEGER
);

-- Ранжирование проектных предложений
CREATE TABLE ProposalRanking (
    ID SERIAL PRIMARY KEY, -- Уникальный идентификатор ранжирования
    ProposalID INTEGER REFERENCES ProjectProposals(ProposalID), -- Ссылка на проектное предложение
    OtherPrinciplesOfRanking TEXT, -- Другие принципы ранжирования
    RealizationImportance DECIMAL(10, 2), -- Важность реализации (из абсолютной приоритетности)
    UnrealizedCriticality DECIMAL(10, 2), -- Критичность нереализации (из абсолютной приоритетности)
    AbsoluteRank DECIMAL(10, 2), -- Абсолютный ранг (из абсолютной приоритетности)
    StrategicRealizationImportance DECIMAL(10, 2), -- Важность реализации (из стратегической приоритетности)
    StrategicUnrealizedCriticality DECIMAL(10, 2), -- Критичность нереализации (из стратегической приоритетности)
    StrategicRank DECIMAL(10, 2) -- Стратегический ранг (из стратегической приоритетности)
);
