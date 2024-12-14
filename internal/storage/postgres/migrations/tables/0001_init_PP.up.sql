CREATE TABLE project_suggestion {
    ID                                                 INT PRIMARY KEY,
    appearance_interval                                TEXT,
    appearance_interval2                               DATA,
    status_PP                                          VARCHAR(255),
    general_classification_ID                          INT NOT NULL,
    project_life_cycle_classification_ID               INT NOT NULL,
    product_life_cycle_classification_ID               INT NO NULL,
    estm_PP_ID                                         INT NOT NULL,
    range_PP_ID                                        INT NOT NULL,
    FOREIGN KEY (general_classification_ID)            REFERENCES general_classification(ID),
    FOREIGN KEY (project_life_cycle_classification_ID) REFERENCES project_life_cycle_classification(ID),
    FOREIGN KEY (product_life_cycle_classification_ID) REFERENCES product_life_cycle_classification(ID),
    FOREIGN KEY (estm_PP_ID)                           REFERENCES estm_PP(ID),
    FOREIGN KEY (range_PP_ID)                          REFERENCES range_PP(ID),
};


CREATE TABLE general_classification {
    ID          INT PRIMARY KEY,
    ptrbnst     VARCHAR(255),
    prd_usl     VARCHAR(255),
    napr_invest VARCHAR(255),
};


CREATE TABLE project_life_cycle_classification {
    ID          INT PRIMARY KEY,
    rlstn       VARCHAR(255),
    mngmnt      VARCHAR(255),
    technology  VARCHAR(255),
    duration    VARCHAR(255),
    cost        VARCHAR(255),
    risk        VARCHAR(255),
};

CREATE TABLE product_life_cycle_classification {
    ID                  INT PRIMARY KEY,
    duration            VARCHAR(255),
    profability         VARCHAR(255),
    risk                VARCHAR(255),
    non_impl            VARCHAR(255),
    strst_algn          VARCHAR(255),
    strt_conflict       VARCHAR(255),
};


CREATE TABLE estm_PP {
    ID                                  INT PRIMARY KEY,
    prdct_dur_estm                      VARCHAR(255),
    prj_dur_estm                        VARCHAR(255),
    prj_fin_res_estm_ID                 INT NOT NULL,
    prdct_P_infl_estm_ID                INT NOT NULL,
    FOREIGN KEY (prj_fin_res_estm_ID)   REFERENCES prj_fin_res_estm(ID),
    FOREIGN KEY (prdct_P_infl_estm_ID)  REFERENCES prdct_P_infl_estm(ID)
};

CREATE TABLE prj_fin_res_estm {
    ID                  INT PRIMARY KEY,
    prj_cost_estm       VARCHAR(255),
    prj_rsrs_reserv     VARCHAR(255),
};

CREATE TABLE prdct_P_infl_estm {
    ID                      INT PRIMARY KEY,
    p1...4_rls_infl_estm    VARCHAR(255),
    p1...4_ntrls_infl_estm  VARCHAR(255),
};

CREATE TABLE range_pp {
    ID                           INT PRIMARY KEY,
    abslt_prior_ID               INT NOT NULL,
    strat_prior_ID               INT NOT NULL,
    rang_reserv                  VARCHAR(255),
    FOREIGN KEY (strat_prior_ID) REFERENCES strat_prior(ID),
    FOREIGN KEY (abslt_prior_ID) REFERENCES abslt_prior(ID),
};

CREATE TABLE abslt_prior {
    ID                  INT PRIMARY KEY,
    algr_1              VARCHAR(255),
    algr_2              VARCHAR(255),
    sum_algr_1_2        VARCHAR(255),
};

CREATE TABLE strat_prior {
    ID              INT PRIMARY KEY,
    algr_3          VARCHAR(255),
    algr_4          VARCHAR(255),
    sum_algr_3_4    VARCHAR(255),
};