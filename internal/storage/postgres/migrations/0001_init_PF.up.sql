CREATE TABLE pf (
    ID                          INT PRIMARY KEY, 
    date_fixation               DATE, 
    pp_id                       INT NOT NULL, 
    status_id                   INT NOT NULL,
    range_results_id            INT NOT NULL,
    duration_fact_realization   VARCHAR(255),
    fact_realization_param_id   INT NOT NULL,
    param_po_otchetam_realiz_id INT NOT NULL,
    cost_on_realization_id      INT NOT NULL,
    FOREIGN KEY(fact_realization_param_id)      REFERENCES fact_realization_param(ID),
    FOREIGN KEY(param_po_otchetam_realiz_id)    REFERENCES param_po_otchetam_realiz(ID),
    FOREIGN KEY(cost_on_realization_id)         REFERENCES cost_on_realization(ID),
    FOREIGN KEY(status_id)                      REFERENCES status(ID),
    FOREIGN KEY(pp_id)                          REFERENCES project_suggestion(ID),
    FOREIGN KEY(range_results_id)               REFERENCES range_results(ID)
);

CREATE TABLE status (
    ID                      INT PRIMARY KEY,
    pp_review_status        VARCHAR(255),
    pp_realization_status   VARCHAR(255),
    pp_usage_result         VARCHAR(255)
);

CREATE TABLE range_results (
    ID          INT PRIMARY KEY,
    alghoritm_1 VARCHAR(255),
    alghoritm_2 VARCHAR(255)
);

CREATE TABLE fact_realization_param (
    ID                      INT PRIMARY KEY,
    narast_itog             VARCHAR(255),
    otklon_nakopl_znach_BL  VARCHAR(255),
    prognoz_do_zaversh      VARCHAR(255),
    prognoz_po_zaversh      VARCHAR(255)
);

CREATE TABLE param_po_otchetam_realiz (
    ID                      INT PRIMARY KEY,
    narast_itog             VARCHAR(255),
    otklon_nakopl_znach_BL  VARCHAR(255),
    prognoz_do_zaversh      VARCHAR(255),
    prognoz_po_zaversh      VARCHAR(255)
);

CREATE TABLE cost_on_realization (
    ID                      INT PRIMARY KEY,
    narast_itog             VARCHAR(255),
    otklon_nakopl_znach_BL  VARCHAR(255),
    prognoz_do_zaversh      VARCHAR(255),
    prognoz_po_zaversh      VARCHAR(255)
);