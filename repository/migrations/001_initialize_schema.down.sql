-- use in case of message that migration version is dirty:

DROP TABLE IF EXISTS schema_migrations;

-- order of removing tables (in case it is needed) as they rely on foreign keys
DROP TABLE IF EXISTS spec_execution;
DROP TABLE IF EXISTS session_execution;
DROP TABLE IF EXISTS spec;
DROP TABLE IF EXISTS userProject;
DROP TABLE IF EXISTS project;
DROP TABLE IF EXISTS apiKey;
DROP TABLE IF EXISTS users;