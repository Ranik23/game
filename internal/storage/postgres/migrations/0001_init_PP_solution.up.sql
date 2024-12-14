CREATE TABLE pf_solution {
    ID INT PRIMARY KEY,
    pf_id INT NOT NULL,
    date_interval DATE,
    parametr_id INT NOT NULL,
    FOREIGN KEY (parametr_id) REFERENCES pf_solution_parametr_value(ID),
};

CREATE TABLE pf_solution_parametr_value {
    ID INT PRIMARY KEY,
    value VARCHAR(255),
};