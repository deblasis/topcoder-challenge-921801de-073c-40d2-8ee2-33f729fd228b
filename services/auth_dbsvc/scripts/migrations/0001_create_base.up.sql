CREATE TABLE if not exists roles(
   role VARCHAR(255) NOT NULL,
   PRIMARY KEY(role)
);

CREATE TABLE if not exists users(
   id integer generated always as identity,
   kong_id VARCHAR(255) NULL UNIQUE,
   role VARCHAR(255),
   username VARCHAR(255) NOT NULL UNIQUE,
   password TEXT NOT NULL,
 
   PRIMARY KEY(id),
   CONSTRAINT fk_role
      FOREIGN KEY(role) 
	  REFERENCES roles(role)
);


INSERT INTO roles (role) VALUES ('Ship'),('Station'),('Command');