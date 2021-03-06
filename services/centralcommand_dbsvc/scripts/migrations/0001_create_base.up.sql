-- The MIT License (MIT)
--
-- Copyright (c) 2021 Alessandro De Blasis <alex@deblasis.net>  
--
-- Permission is hereby granted, free of charge, to any person obtaining a copy
-- of this software and associated documentation files (the "Software"), to deal
-- in the Software without restriction, including without limitation the rights
-- to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
-- copies of the Software, and to permit persons to whom the Software is
-- furnished to do so, subject to the following conditions:
--
-- The above copyright notice and this permission notice shall be included in all
-- copies or substantial portions of the Software.
--
-- THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
-- IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
-- FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
-- AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
-- LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
-- OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
-- SOFTWARE. 
--

CREATE TABLE if not exists stations(
   id UUID NOT NULL,
   capacity FLOAT NOT NULL,
  
   PRIMARY KEY(id)
);

CREATE TABLE if not exists docks(
   id UUID NOT NULL,
   station_id UUID NOT NULL,
   num_docking_ports INTEGER NOT NULL,
 
   PRIMARY KEY(id),
   CONSTRAINT fk_station
      FOREIGN KEY(station_id) 
	  REFERENCES stations(id)
);


CREATE TABLE if not exists docked_ships(
   dock_id UUID NOT NULL,
   ship_id UUID NOT NULL UNIQUE,
   docked_since TIMESTAMP NULL,
   dock_duration INT,
   waiting_for_ship_since TIMESTAMP,

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

CREATE OR REPLACE FUNCTION reservations_expired(dock_holding_period INT) RETURNS INTEGER
    AS $$
DECLARE ret INTEGER;    
BEGIN
  WITH deleted AS (
     DELETE FROM docked_ships WHERE docked_since is null AND waiting_for_ship_since + INTERVAL '1 second'*dock_holding_period < NOW() RETURNING *
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
      min((ds.docked_since+ INTERVAL '1 second'*ds.dock_duration)-NOW()) as seconds_until_next_available,
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
   COALESCE(d.weight,0) as weight,   
   d.seconds_until_next_available as seconds_until_next_available
	from docked as d;   
	

CREATE OR REPLACE VIEW stations_view AS
	SELECT s.id as id, 
   s.capacity as capacity, 
   coalesce((select sum(weight) from docks_view where station_id=s.id group by station_id),0) as used_capacity
	from stations as s;

CREATE OR REPLACE  FUNCTION stations_available_for_ship(ship_id UUID) RETURNS 
   TABLE (station_id UUID, capacity FLOAT, used_capacity FLOAT, dock_id UUID, num_docking_ports INTEGER, occupied BIGINT, weight FLOAT)
    AS $$
BEGIN
  return query
      with ship as (
         select s.id, s.weight from ships s where id = ship_id
      ), stations_with_capacity as (
         select st.id as station_id, 
         st.capacity as capacity,
         st.used_capacity as used_capacity,
         d.id as dock_id,
         d.num_docking_ports as num_docking_ports,
         d.occupied as occupied,
         d.weight as weight
         from stations_view st 
         inner join docks_view d on (d.station_id = st.id and d.num_docking_ports-d.occupied>0)
         inner join ship on (ship.id = ship_id)
         where st.capacity-st.used_capacity>=ship.weight
      ) 
      select swc.station_id, swc.capacity, swc.used_capacity, swc.dock_id, swc.num_docking_ports, swc.occupied, swc.weight
      from stations_with_capacity swc;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION get_next_available_docking_station_for_ship(_ship_id UUID) 
RETURNS TABLE (dock_id UUID, station_id UUID, ship_weight FLOAT, available_capacity FLOAT, available_docks_at_station BIGINT, seconds_until_next_available INT) AS 
$$
BEGIN
--if the dock is reserved, we return that as the next available, it will be kept on hold and result as "occupied" until released or used for landing
   IF EXISTS (SELECT FROM docked_ships ds where ds.ship_id = _ship_id) then
      return query with ship as (
         select id, weight from ships where id = _ship_id
      )
      select ds.dock_id, d.station_id, ship.weight as ship_weight, 
      ship.weight as available_capacity, 1::bigint as available_docks_at_station, 0 as seconds_until_next_available
      from docked_ships ds 
      INNER JOIN ship on (ship.id=_ship_id)
      INNER JOIN docks d on (d.id=ds.dock_id);
   ELSE
-- if the ship didn't reserve a dock already, we look for the first available, reserve it and return it   
      return query with ship as (
         select id, weight from ships where id = _ship_id
      ), stations_with_capacity as (
         select st.id, st.capacity-st.used_capacity as available_capacity,
         d.num_docking_ports-d.occupied as available_docks_at_station,
         d.seconds_until_next_available as seconds_until_next_available,
         d.id as dock_id,
         st.id as station_id 
         from stations_view st 
         inner join docks_view d on (d.station_id = st.id)
      ), next_available as ( 
      select swc.dock_id, swc.station_id, ship.weight as ship_weight, swc.available_capacity, swc.available_docks_at_station, 
      CASE 
         when swc.seconds_until_next_available is null then 0
         ELSE (select extract(epoch from swc.seconds_until_next_available))::int
      END
      as seconds_until_next_available
      from stations_with_capacity swc 
      inner join ship on (ship.id = _ship_id)
      order by available_capacity desc, available_docks_at_station desc, seconds_until_next_available asc limit 1
      ), reservation as (
         insert into docked_ships(dock_id, ship_id, waiting_for_ship_since)
         select n.dock_id, _ship_id, NOW() from next_available n
         where n.available_capacity >= n.ship_weight and n.available_docks_at_station >= 1
      )
      select * from next_available;
   END IF;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE VIEW ships_view AS
   WITH docked as (
      SELECT s.id as id, s.weight as weight,
      ds.dock_id as dock_id FROM ships s 
      LEFT JOIN docked_ships ds ON (s.id = ds.ship_id AND ds.docked_since is not null)
   )
   select id, weight, dock_id,
   CASE 
      WHEN dock_id is null then 'in-flight'
      WHEN dock_id is not null then 'docked'
   END status
   from docked;

CREATE INDEX IF NOT EXISTS station_id_idx on docks(station_id); 
CREATE INDEX IF NOT EXISTS num_docking_ports_idx on docks(num_docking_ports);
