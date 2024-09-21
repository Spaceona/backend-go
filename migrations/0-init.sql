DROP TABLE IF EXISTS StatusChange;
DROP TABLE IF EXISTS machine;
DROP TABLE IF EXISTS board;
DROP TABLE IF EXISTS building;
DROP TABLE IF EXISTS client;
PRAGMA journal_mode=WAL;
PRAGMA busy_timeout = 3000;
CREATE TABLE StatusChange(
    id integer primary key, -- this is auto assigned
    machine_id integer,
    firmware_version text,
    status integer, -- can only be 1 or 0
    confidence integer,
    date text,
    month text, -- this is the name of the month used for historical data
    day_of_the_week text, --used for histograms
    hour integer,
    year integer,
    FOREIGN KEY(machine_id) REFERENCES machine (id)
);
CREATE TABLE board(
    mac_address text unique primary key ,
    valid integer,
    client_name text,
    heart_beat_interval integer,
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
    id integer primary key autoincrement,
    number integer, -- this must be set by the client
    status integer,
    mac_address text,
    type text,
    building_name text,
    client_name text,
    last_write text,
    estimated_duration integer, -- we get this by averaging run durations
    number_of_runs integer, -- how many runs we have compleated
    foreign key (client_name) REFERENCES client(name),
    FOREIGN KEY (building_name) REFERENCES building(name),
    FOREIGN KEY (mac_address) REFERENCES  board(mac_address)
);

CREATE INDEX buildings_machines on machine (building_name)