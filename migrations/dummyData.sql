insert into client(name,key) values ('Lafayette','fjdslka');
insert into client(name,key) values ('WPI','gfhdsaj');

insert into building(name, client_name) values ('Watson Hall','Lafayette');
insert into building(name, client_name) values ('Reeder House','Lafayette');

insert into building(name, client_name) values ('Faraday','WPI');
insert into building(name, client_name) values ('East Hall','WPI');

-- watson
insert into board(mac_address, valid,client_name,heart_beat_interval) VALUES ('64:E8:33:86:DB:A4',1,'Lafayette',15);
insert into board(mac_address, valid,client_name,heart_beat_interval) VALUES ('testing2',1,'Lafayette',15);

-- reeder
insert into board(mac_address, valid,client_name,heart_beat_interval) VALUES ('testing3',1,'Lafayette',15);
insert into board(mac_address, valid,client_name,heart_beat_interval) VALUES ('testing4',1,'Lafayette',15);

--faraday
insert into board(mac_address, valid,client_name,heart_beat_interval) VALUES ('testing5',1,'WPI',15);
insert into board(mac_address, valid,client_name,heart_beat_interval) VALUES ('testing6',1,'WPI',15);

--east
insert into board(mac_address, valid,client_name,heart_beat_interval) VALUES ('testing7',1,'WPI',15);
insert into board(mac_address, valid,client_name,heart_beat_interval) VALUES ('testing8',1,'WPI',15);

--un assigned
insert into board(mac_address, valid,heart_beat_interval) VALUES ('testing9',1,15);
insert into board(mac_address, valid,heart_beat_interval) VALUES ('testing10',1,15);


insert into machine(number,mac_address,type,building_name,client_name,status) values (1,'64:E8:33:86:DB:A4','Washer','Watson Hall','Lafayette',0);
insert into machine(number,mac_address,type,building_name,client_name,status) values (2,'testing2','Dryer','Watson Hall','Lafayette',0);

insert into machine(number,mac_address,type,building_name,client_name,status) values (1,'testing3','Washer','Reeder House','Lafayette',0);
insert into machine(number,mac_address,type,building_name,client_name,status) values (2,'testing4','Dryer','Reeder House','Lafayette',0);

insert into machine(number,mac_address,type,building_name,client_name,status) values (1,'testing5','Washer','Faraday','WPI',0);
insert into machine(number,mac_address,type,building_name,client_name,status) values (2,'testing6','Dryer','Faraday','WPI',0);

insert into machine(number,mac_address,type,building_name,client_name,status) values (1,'testing7','Washer','East Hall','WPI',0);
insert into machine(number,mac_address,type,building_name,client_name,status) values (2,'testing8','Dryer','East Hall','WPI',0);