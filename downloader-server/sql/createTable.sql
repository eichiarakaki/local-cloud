-- Do not change the name of these property.
CREATE TABLE IF NOT EXISTS localcloud ( -- TABLE NAME MUST MATCH WITH THE config.json's 'mysql-table-name'
    id INT AUTO_INCREMENT,
    filepath VARCHAR(250) NOT NULL,
    filename VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
)
-- CREATE TABLE IF NOT EXISTS testdb (
--     id INT AUTO_INCREMENT,
--     filepath VARCHAR(250) NOT NULL,
--     filename VARCHAR(100) NOT NULL,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     PRIMARY KEY (id)
-- )
