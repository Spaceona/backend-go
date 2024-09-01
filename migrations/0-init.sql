DROP TABLE IF EXISTS StatusChange;
DROP TABLE IF EXISTS machine;
DROP TABLE IF EXISTS board;
DROP TABLE IF EXISTS building;
DROP TABLE IF EXISTS client;
CREATE TABLE StatusChange(
    id integer primary key, -- this is auto assigned
    mac_address text,
    current_version text,
    status integer, -- can only be 1 or 0
    confidence integer,
    date text,
    FOREIGN KEY(mac_address) REFERENCES board (mac_address)
);
CREATE TABLE board(
    mac_address text unique primary key ,
    valid integer,
    client_name text,
    foreign key (client_name) references client(name)
);
CREATE TABLE client(
    name text primary key,
    key text unique,
    salt blob
);
CREATE TABLE building(
    name text primary key,
    client_name text,
    FOREIGN KEY (client_name) REFERENCES client(name)
);
CREATE TABLE machine(
    id integer primary key, -- this must be set by the client
    number integer,
    status integer,
    mac_address text unique ,
    type text,
    building_name text,
    client_name text,
    last_write text,
    foreign key (client_name) REFERENCES client(name),
    FOREIGN KEY (building_name) REFERENCES building(name),
    FOREIGN KEY (mac_address) REFERENCES  board(mac_address)
);

CREATE INDEX buildings_machines on machine (building_name)