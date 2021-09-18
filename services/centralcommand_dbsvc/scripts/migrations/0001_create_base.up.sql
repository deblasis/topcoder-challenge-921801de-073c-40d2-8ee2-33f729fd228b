CREATE TABLE if not exists stations(
   id UUID NOT NULL,
   capacity FLOAT NOT NULL,
  
   PRIMARY KEY(id)
);

CREATE TABLE if not exists docks(
   id UUID NOT NULL,
   station_id UUID NOT NULL,
   num_docking_ports INTEGER NOT NULL,
   
   weight FLOAT NOT NULL DEFAULT(0),
 
   PRIMARY KEY(id),
   CONSTRAINT fk_station
      FOREIGN KEY(station_id) 
	  REFERENCES stations(id)
);



CREATE TABLE if not exists docked_ships(
   dock_id UUID NOT NULL,
   ship_id UUID NOT NULL UNIQUE,
   docked_since TIMESTAMP,
   dock_duration INT,

   PRIMARY KEY(dock_id, ship_id)
);


CREATE OR REPLACE FUNCTION ships_have_left() RETURNS INTEGER
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
   id UUID NOT NULL,
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
      LEFT JOIN ships s on (s.id = ds.ship_id)
      GROUP BY (d.id, ds.ship_id)
   )
	SELECT d.id as id, d.station_id as station_id, 
   d.num_docking_ports as num_docking_ports,
   d.occupied as occupied,
   COALESCE(d.weight,0) as weight   
	from docked as d;   
	

CREATE OR REPLACE VIEW stations_view AS
	SELECT s.id as id, 
   s.capacity as capacity, 
   coalesce((select sum(weight) from docks_view where station_id=s.id group by station_id),0) as used_capacity
	from stations as s;


CREATE OR REPLACE FUNCTION get_next_available_docking_station_for_ship(ship_id UUID) RETURNS 
   TABLE (dock_id UUID, station_id UUID, available_capacity FLOAT, available_docks_at_station BIGINT)
    AS $$
BEGIN
  return query
      with ship as (
         select id, weight from ships where id = ship_id
      ), stations_with_capacity as (
         select st.id, st.capacity-st.used_capacity as available_capacity,
         d.num_docking_ports-d.occupied as available_docks_at_station,
         d.id as dock_id,
         st.id as station_id 
         from stations_view st 
         inner join docks_view d on (d.station_id = st.id and d.num_docking_ports-d.occupied>0)
         where capacity-used_capacity>(select weight from ship)
      ) 
      select swc.dock_id, swc.station_id, swc.available_capacity, swc.available_docks_at_station from stations_with_capacity swc order by available_capacity desc, available_docks_at_station desc limit 1;
END;
$$ LANGUAGE plpgsql;


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
CREATE INDEX IF NOT EXISTS num_docking_ports_idx on docks(num_docking_ports);
