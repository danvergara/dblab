-- This script is intended to be run inside a Pluggable Database (PDB), such as FREEPDB1.
-- Ensure you are connected to the PDB before running this script.
-- Example connection: sqlplus sys/password@//localhost:1521/FREEPDB1 as sysdba

-- Drop users if they exist, to make the script re-runnable

-- How to use:

-- Step 1: Make sure the Oracle container is running:

-- docker compose up -d oracle

-- Step 2: Run the script via docker exec:

-- docker exec -i dblab-oracle-1 sqlplus sys/password@//localhost:1521/FREEPDB1 as sysdba < scripts/oracle_test_setup.sql

-- If the container name is different, check it with:
-- docker compose ps

BEGIN
   EXECUTE IMMEDIATE 'DROP USER user1 CASCADE';
EXCEPTION
   WHEN OTHERS THEN
      IF SQLCODE != -1918 THEN
         RAISE;
      END IF;
END;
/

BEGIN
   EXECUTE IMMEDIATE 'DROP USER user2 CASCADE';
EXCEPTION
   WHEN OTHERS THEN
      IF SQLCODE != -1918 THEN
         RAISE;
      END IF;
END;
/

-- Create user 1 and grant privileges
CREATE USER user1 IDENTIFIED BY password;
GRANT CREATE SESSION, RESOURCE TO user1;
ALTER USER user1 QUOTA UNLIMITED ON USERS;

-- Create table in user1's schema
CREATE TABLE user1.test_table (
    id NUMBER PRIMARY KEY,
    name VARCHAR2(50)
);

INSERT INTO user1.test_table (id, name) VALUES (1, 'test_data_1');
INSERT INTO user1.test_table (id, name) VALUES (2, 'test_data_2');

-- Create user 2 and grant privileges
CREATE USER user2 IDENTIFIED BY password;
GRANT CREATE SESSION TO user2;

-- Grant user2 access to user1's table
GRANT SELECT ON user1.test_table TO user2;

-- Grant user2 access to view all tables (necessary for the application's ShowTables function)
GRANT SELECT ANY DICTIONARY TO user2;
