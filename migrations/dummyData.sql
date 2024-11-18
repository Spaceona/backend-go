insert into client(name,key) values ('Lafayette','fjdslka');
insert into client(name,key) values ('WPI','gfhdsaj');
insert into client(name,key,salt) values ('lafTest','5J*P)S6ufUcknP1IllEn%(tu7',0x4EA6D6C7AA008ADCC190177975AD34D3AC61BA9E30448EEC8B);

insert into building(name, client_name) values ('Watson Hall','Lafayette');
insert into building(name, client_name) values ('Reeder House','Lafayette');

insert into building(name, client_name) values ('Faraday','WPI');
insert into building(name, client_name) values ('East Hall','WPI');


insert into building(name, client_name) values ('Reeder','lafTest');
insert into building(name, client_name) values ('Watson','lafTest');

insert into board(mac_address, valid,client_name,heart_beat_interval) VALUES ('24:EC:4A:D7:0F:70',1,'lafTest',15);
insert into board(mac_address, valid,client_name,heart_beat_interval) VALUES ('E4:B0:63:07:9C:54',1,'lafTest',15);
insert into board(mac_address, valid,client_name,heart_beat_interval) VALUES ('24:EC:4A:D7:0F:80',1,'lafTest',15);
insert into board(mac_address, valid,client_name,heart_beat_interval) VALUES ('E4:B0:63:07:9C:58',1,'lafTest',15);
insert into board(mac_address, valid,client_name,heart_beat_interval) VALUES ('E4:B0:63:07:9C:44',1,'lafTest',15);
insert into board(mac_address, valid,client_name,heart_beat_interval) VALUES ('24:EC:4A:D7:0F:60',1,'lafTest',15);
insert into board(mac_address, valid,client_name,heart_beat_interval) VALUES ('24:EC:4A:D7:0F:78',1,'lafTest',15);
insert into board(mac_address, valid,client_name,heart_beat_interval) VALUES ('24:EC:4A:D7:0F:7C',1,'lafTest',15);
insert into board(mac_address, valid,client_name,heart_beat_interval) VALUES ('24:EC:4A:D7:0F:74',1,'lafTest',15);

insert into machine(number,mac_address,type,building_name,client_name,status) values (1,'24:EC:4A:D7:0F:70','Washer','Watson','lafTest',0);
insert into machine(number,mac_address,type,building_name,client_name,status) values (2,'E4:B0:63:07:9C:54','Washer','Watson','lafTest',0);
insert into machine(number,mac_address,type,building_name,client_name,status) values (3,'24:EC:4A:D7:0F:80','Washer','Watson','lafTest',0);
insert into machine(number,mac_address,type,building_name,client_name,status) values (4,'E4:B0:63:07:9C:58','Washer','Watson','lafTest',0);
insert into machine(number,mac_address,type,building_name,client_name,status) values (1,'E4:B0:63:07:9C:44','Dryer','Watson','lafTest',0);
insert into machine(number,mac_address,type,building_name,client_name,status) values (2,'24:EC:4A:D7:0F:60','Dryer','Watson','lafTest',0);
insert into machine(number,mac_address,type,building_name,client_name,status) values (3,'24:EC:4A:D7:0F:78','Dryer','Watson','lafTest',0);
insert into machine(number,mac_address,type,building_name,client_name,status) values (4,'24:EC:4A:D7:0F:7C','Dryer','Watson','lafTest',0);


insert into machine(number,mac_address,type,building_name,client_name,status) values (1,'24:EC:4A:D7:0F:74','Washer','Reeder','lafTest',0);
insert into machine(number,mac_address,type,building_name,client_name,status) values (1,'E4:B0:63:07:9C:6C','Dryer','Reeder','lafTest',0); -- robbie test board






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