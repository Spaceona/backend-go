-- select all lattest
with status as (
    select
        StatusChange.id,
        StatusChange.mac_address,
        StatusChange.current_version,
        StatusChange.status,
        StatusChange.confidence,
        StatusChange.date,
        row_number() over (partition by StatusChange.mac_address order by StatusChange.date desc) as RN
    FROM StatusChange
) select * from status where RN = 1 order by status.mac_address;

-- select all latest with client name
with status as (
    select
        StatusChange.id,
        StatusChange.mac_address,
        StatusChange.current_version,
        StatusChange.status,
        StatusChange.confidence,
        StatusChange.date,
        row_number() over (partition by StatusChange.mac_address order by StatusChange.date desc) as RN
    FROM StatusChange
) select status.id,
         status.mac_address,
         status.current_version,
         status.status,
         status.confidence,
         status.date,
         machine.id as machine_id,
         machine.number as machine_num
    from status
                              inner join main.board board on board.mac_address = status.mac_address
                              inner join main.machine machine on board.mac_address = machine.mac_address
where RN = 1
  and client_name = 'Lafayette'
order by status.mac_address;

-- client and building
with status as (
    select
        StatusChange.id,
        StatusChange.mac_address,
        StatusChange.current_version,
        StatusChange.status,
        StatusChange.confidence,
        StatusChange.date,
        row_number() over (partition by StatusChange.mac_address order by StatusChange.date desc) as RN
    FROM StatusChange
) select status.id,
         status.mac_address,
         status.current_version,
         status.status,
         status.confidence,
         status.date,
         machine.id as machine_id,
         machine.number as machine_num
         status.date
from status
                                inner join main.board board on board.mac_address = status.mac_address
                                inner join main.machine machine on board.mac_address = machine.mac_address
where RN = 1
  and client_name = 'Lafayette' and building_name = 'Watson'
order by status.mac_address;

-- client and building and type
with status as (
    select
        StatusChange.id,
        StatusChange.mac_address,
        StatusChange.current_version,
        StatusChange.status,
        StatusChange.confidence,
        StatusChange.date,
        row_number() over (partition by StatusChange.mac_address order by StatusChange.date desc) as RN
    FROM StatusChange
) select status.id,
         status.mac_address,
         status.current_version,
         status.status,
         status.confidence,
         status.date,
         machine.id as machine_id,
         machine.number as machine_num
        from status
                                inner join main.board board on board.mac_address = status.mac_address
                                inner join main.machine machine on board.mac_address = machine.mac_address
where RN = 1
  and client_name = 'Lafayette'
  and building_name = 'Watson'
  and type = 'Washer'
order by status.mac_address;


-- by id
with status as (
select
    StatusChange.id,
    StatusChange.mac_address,
    StatusChange.current_version,
    StatusChange.status,
    StatusChange.confidence,
    StatusChange.date,
    row_number() over (partition by StatusChange.mac_address order by StatusChange.date desc) as RN
FROM StatusChange
) select status.id,
         status.mac_address,
         status.current_version,
         status.status,
         status.confidence,
         status.date,
         machine.id as machine_id,
         machine.number as machine_num,
         machine.building_name as building_name,
         machine.client_name as client_name
  from status
                          inner join main.board board on board.mac_address = status.mac_address
                          inner join main.machine machine on board.mac_address = machine.mac_address
where RN = 1
and machine.id = 1
order by status.mac_address;