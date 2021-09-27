--
-- Copyright 2021 Alessandro De Blasis <alex@deblasis.net>
--
-- Licensed under the Apache License, Version 2.0 (the "License");
-- you may not use this file except in compliance with the License.
-- You may obtain a copy of the License at
--
-- http://www.apache.org/licenses/LICENSE-2.0
--
-- Unless required by applicable law or agreed to in writing, software
-- distributed under the License is distributed on an "AS IS" BASIS,
-- WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
-- See the License for the specific language governing permissions and
-- limitations under the License.
--
DROP VIEW ships_view;
DROP VIEW stations_view;
DROP FUNCTIOn stations_available_for_ship;
DROP FUNCTION reservations_expired;
DROP FUNCTION get_next_available_docking_station_for_ship;
DROP VIEW docks_view;

DROP TABLE ships;
DROP FUNCTION ships_have_left;
DROP TABLE docked_ships;
DROP TABLE docks;
DROP TABLE stations;

