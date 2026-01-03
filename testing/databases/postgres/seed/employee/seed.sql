-- Create schema and access control
CREATE SCHEMA employee;

CREATE ROLE employee_access;
GRANT USAGE ON SCHEMA employee TO employee_access;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA employee TO employee_access;
ALTER DEFAULT PRIVILEGES IN SCHEMA employee GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO employee_access;
ALTER DEFAULT PRIVILEGES IN SCHEMA employee GRANT USAGE, SELECT ON SEQUENCES TO employee_access;

-- Create user with access to the schema
CREATE USER "Jimbobby" WITH PASSWORD 'Jimbobby';
GRANT employee_access TO "Jimbobby";
ALTER USER "Jimbobby" SET search_path TO employee;

-- Revoke public access
REVOKE ALL ON SCHEMA employee FROM PUBLIC;

-- Set search path for this session
SET search_path TO employee;


CREATE TABLE employees (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    annual_salary NUMERIC(10, 2) NOT NULL,
    CONSTRAINT employees_email_check CHECK (email LIKE '%_@__%.__%'),
    CONSTRAINT employees_annual_salary_check CHECK (annual_salary >= 0)
);
CREATE INDEX idx_employees_name ON employees (name);


CREATE TABLE departments (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    location VARCHAR(100) NOT NULL
);
CREATE INDEX idx_departments_name ON departments (name);


CREATE TABLE positions (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    level INT NOT NULL,
    CONSTRAINT positions_level_check CHECK (level >= 1 AND level <= 10)
);
CREATE INDEX idx_positions_title ON positions (title);
CREATE INDEX idx_positions_level ON positions (level);


CREATE TABLE employee_department (
    employee_id INT REFERENCES employees(id),
    department_id INT REFERENCES departments(id),
    position INT REFERENCES positions(id),
    joined_date DATE,
    parent_department_id INT,
    PRIMARY KEY (employee_id, department_id)
);
CREATE INDEX idx_employee_department_position ON employee_department (position);
CREATE INDEX idx_employee_department_joined_date ON employee_department (joined_date);
CREATE INDEX idx_employee_department_parent_department_id ON employee_department (parent_department_id);


CREATE VIEW employee_overview AS
SELECT e.id AS employee_id,
       e.name AS employee_name,
       e.email AS employee_email,
       e.annual_salary AS employee_annual_salary,
       d.name AS department_name,
       p.title AS position_title,
       p.level AS position_level,
       ed.joined_date AS date_joined
FROM employees e
LEFT JOIN employee_department ed ON e.id = ed.employee_id
LEFT JOIN departments d ON ed.department_id = d.id
LEFT JOIN positions p ON ed.position = p.id;


CREATE MATERIALIZED VIEW total_salary_by_department AS
SELECT d.id AS department_id,
       d.name AS department_name,
       SUM(e.annual_salary) AS total_annual_salary
FROM departments d
LEFT JOIN employee_department ed ON d.id = ed.department_id
LEFT JOIN employees e ON ed.employee_id = e.id
GROUP BY d.id, d.name;


CREATE PROCEDURE refresh_total_salary_by_department()
LANGUAGE plpgsql
AS $$
BEGIN
    REFRESH MATERIALIZED VIEW employee.total_salary_by_department;
END;
$$;


CREATE FUNCTION refresh_total_salary_by_department_trigger_fn()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    CALL employee.refresh_total_salary_by_department();
    RETURN NULL;
END;
$$;

CREATE TRIGGER trg_refresh_total_salary_by_department
AFTER INSERT OR UPDATE OR DELETE ON employee_department
FOR EACH STATEMENT
EXECUTE FUNCTION refresh_total_salary_by_department_trigger_fn();


-- Grant execute on procedures/functions to the group
GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA employee TO employee_access;
GRANT EXECUTE ON ALL PROCEDURES IN SCHEMA employee TO employee_access;


-- Initial data population
INSERT INTO employees (name, email, annual_salary) VALUES
('Alice Johnson', 'alice.johnson@example.com', 75000.00),
('Bob Smith', 'bob.smith@example.com', 68000.00),
('Charlie Kirk', 'charlie.kirk@example.com', 72000.00),
('Diana Prince', 'diana.prince@example.com', 78000.00),
('Mike Hunt', 'mike.hunt@example.com', 82000.00);

INSERT INTO departments (name, location) VALUES
('Human Resources', 'Brisbane'),
('Engineering', 'Sydney');

INSERT INTO positions (title, level) VALUES
('Manager', 5),
('Software Engineer', 4);

INSERT INTO employee_department (employee_id, department_id, position, joined_date, parent_department_id) VALUES
(1, 1, 1, '2020-01-15', NULL),
(2, 2, 2, '2019-03-22', NULL),
(3, 2, 2, '2021-07-30', NULL),
(4, 1, 1, '2018-11-05', NULL),
(5, 2, 2, '2022-02-14', NULL);










