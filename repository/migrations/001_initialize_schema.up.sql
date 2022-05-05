CREATE TABLE users (
   id uuid NOT NULL, 
   email varchar(100) NOT NULL,
   password varchar(100) NOT NULL,
   CONSTRAINT userPK PRIMARY KEY (id)
);

CREATE TABLE userProject (
   id uuid NOT NULL, 
   userId uuid NOT NULL, 
   projectId uuid NOT NULL,
   CONSTRAINT userProjectPK PRIMARY KEY (id)
);

CREATE TABLE apiKey (
   id uuid NOT NULL,
   userId uuid NOT NULL, 
   name varchar(100) NOT NULL,
   expireAt bigint NOT NULL,
   CONSTRAINT apiKeyPK PRIMARY KEY (id)
);

CREATE TABLE project (
   id uuid NOT NULL,
   name varchar(100) NOT NULL,
   CONSTRAINT projectPK PRIMARY KEY (id)
);

CREATE TABLE spec (
   id uuid NOT NULL,
   filePath varchar(255) NOT NULL,
   projectId uuid NOT NULL,
   CONSTRAINT specPK PRIMARY KEY (id)
);

CREATE TABLE session_execution (
   id uuid NOT NULL,
   projectId uuid NOT NULL,
   startedAt bigint DEFAULT 0 NOT NULL, 
	finishedAt bigint DEFAULT 0 NOT NULL, 
	createdAt bigint NOT NULL,
   CONSTRAINT sessionPK PRIMARY KEY (id)
);

CREATE TABLE spec_execution (
   id uuid NOT NULL,
   specId uuid NOT NULL,
   specName varchar(100) NOT NULL,
   sessionId uuid NOT NULL,
   machineId varchar(50) DEFAULT 'default' NOT NULL,
   startedAt bigint DEFAULT 0 NOT NULL,
	finishedAt bigint DEFAULT 0 NOT NULL,
   estimatedDuration bigint DEFAULT 0  NOT NULL,
   status varchar(10) DEFAULT 'unknown' NOT NULL,
   CONSTRAINT executionPK PRIMARY KEY (id)
);

ALTER TABLE users ADD CONSTRAINT user_email UNIQUE(email);
ALTER TABLE userProject ADD CONSTRAINT userProject_userFK FOREIGN KEY (userId) REFERENCES users (id);
ALTER TABLE userProject ADD CONSTRAINT userProject_projectFK FOREIGN KEY (projectId) REFERENCES project (id);
ALTER TABLE apiKey ADD CONSTRAINT apiKey_userFK FOREIGN KEY (userId) REFERENCES users (id);
ALTER TABLE spec ADD CONSTRAINT spec_projectFK FOREIGN KEY (projectId) REFERENCES project (id);
ALTER TABLE session_execution ADD CONSTRAINT session_projectFK FOREIGN KEY (projectId) REFERENCES project (id);
ALTER TABLE spec_execution ADD CONSTRAINT execution_specFK FOREIGN KEY (specId) REFERENCES spec (id);
ALTER TABLE spec_execution ADD CONSTRAINT execution_sessionFK FOREIGN KEY (sessionId) REFERENCES session_execution (id);