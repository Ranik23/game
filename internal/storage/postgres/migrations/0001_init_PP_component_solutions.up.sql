CREATE TABLE pf_component_solution {
    ID INT PRIMARY KEY,
    pf_ID INT, -- ссылка на таблицу pf
    date_interval DATE,
    component_id INT NOT NULL, -- ссылка на ПП, проекты, контролируемые результаты
    made_decision VARCHAR(255),
    FOREIGN KEY (pf_ID) REFERENCES pf(ID)
    FOREIGN KEY (component_id) REFERENCES project_suggestion(ID)
};