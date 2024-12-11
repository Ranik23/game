CREATE TABLE pf {
    ID INT PRIMARY KEY, 
    date_fixation DATE, 
    project_suggestion_id INT NOT NULL, 
    FOREIGN KEY(project_suggestion_id) REFERENCES project_suggestion(ID),
    
};