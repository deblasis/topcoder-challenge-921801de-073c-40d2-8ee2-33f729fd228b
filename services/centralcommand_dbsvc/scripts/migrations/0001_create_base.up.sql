CREATE TABLE if not exists stations(
   id VARCHAR(255),
   capacity FLOAT NOT NULL,
  
   PRIMARY KEY(id)
);

CREATE TABLE if not exists docks(
   id VARCHAR(255) ,
   station_id VARCHAR(255) NOT NULL,
   num_docking_ports INTEGER NOT NULL,
   weight FLOAT NOT NULL DEFAULT(0),
 
   PRIMARY KEY(id),
   CONSTRAINT fk_station
      FOREIGN KEY(station_id) 
	  REFERENCES stations(id)
);



CREATE TABLE if not exists docked_ships(
   dock_id VARCHAR(255),
   ship_id VARCHAR(255) UNIQUE,
   docked_since TIMESTAMP,
   dock_duration INT,

   --SELECT date_part('epoch',CURRENT_TIMESTAMP)::int --UNIXTIMESTAMP in seconds
   
   PRIMARY KEY(dock_id, ship_id)
);


CREATE FUNCTION ships_have_left() RETURNS INTEGER
    AS $$
DECLARE ret INTEGER;    
BEGIN
  WITH deleted AS (
     DELETE FROM docked_ships WHERE docked_since + INTERVAL '1 second'*dock_duration < NOW() RETURNING *
  ) select count(*) into ret from deleted;
  RETURN ret;
END;
$$ LANGUAGE plpgsql;


CREATE TABLE if not exists ships(
   id VARCHAR(255),
   weight FLOAT NOT NULL,

   PRIMARY KEY(id)
);

CREATE OR REPLACE VIEW docks_view AS
   WITH docked as (
      SELECT 
      d.id as id,
      d.station_id as station_id,
      d.num_docking_ports as num_docking_ports,
      ds.ship_id as ship_id, 
      count(ship_id) as occupied,
      sum(s.weight) as weight
      FROM docks d 
      LEFT JOIN docked_ships ds ON (d.id = ds.dock_id)
      INNER JOIN ships s on (s.id = ds.ship_id)
      GROUP BY (d.id, ds.ship_id)
   )
	SELECT d.id as id, d.station_id as station_id, 
   d.num_docking_ports as num_docking_ports,
   d.occupied as occupied,
   d.weight as weight   
	from docked as d;   
	

CREATE OR REPLACE VIEW stations_view AS
	SELECT s.id as id, 
   s.capacity as capacity, 
   (select sum(weight) from docks_view where station_id=s.id group by station_id) as used_capacity
	from stations as s;


CREATE OR REPLACE VIEW ships_view AS
   WITH docked as (
      SELECT s.id as id, s.weight as weight,
      ds.dock_id as dock_id FROM ships s 
      LEFT JOIN docked_ships ds ON (s.id = ds.ship_id)
   )
   select id, weight, dock_id,
   CASE 
      WHEN dock_id is null then 'in-flight'
      WHEN dock_id is not null then 'docked'
   END status
   from docked;

CREATE INDEX IF NOT EXISTS station_id_idx on docks(station_id); 	
