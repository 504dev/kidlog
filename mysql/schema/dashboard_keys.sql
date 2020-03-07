CREATE TABLE IF NOT EXISTS dashboard_keys (
    id INT NOT NULL AUTO_INCREMENT,
    dash_id INT,
    public_key VARCHAR(1024),
    private_key VARCHAR(1024),
    PRIMARY KEY (id),
    UNIQUE (public_key),
    FOREIGN KEY (dash_id) REFERENCES dashboards(id),
);