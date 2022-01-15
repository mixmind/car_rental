package db

var (
	createCarTable = `CREATE TABLE cars(car_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
					car_comp_name text,
					doors INTEGER,
					big_lag INTEGER,
					small_lag INTEGER,
					adult_place INTEGER,
					condition boolean,
					min_age INTEGER,
					locations TEXT,
					car_group INTEGER,
					description TEXT,
					price INTEGER);`
	createRentTable = `CREATE TABLE rents(rent_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
					car_id INTEGER,
					from_time TIMESTAMP,
					to_time TIMESTAMP,
					location TEXT,
					extras TEXT,
					discounts TEXT,
					rent_detail text,
					FOREIGN KEY(car_id) REFERENCES cars(car_id) ON DELETE RESTRICT
					);`
	InsertIntoCarTable = `INSERT INTO cars(car_comp_name , doors,
												big_lag, small_lag,
												adult_place, condition,
												min_age, locations,
												car_group, description,price) VALUES (?,?,?,?,?,?,?,?,?,?,?)`
	InsertIntoRentTable = `INSERT INTO rents(car_id,
											from_time,
											to_time,
											location,
											extras,
											discounts,
											rent_detail) VALUES (?,?,?,?,?,?,?)`
	SelectCars = `SELECT car_id,
					car_comp_name ,
					doors,
					big_lag,
					small_lag,
					adult_place,
					condition,
					min_age,
					locations,
					car_group,
					description,
					price FROM cars`
	UpdateCar = `UPDATE cars 
				SET car_comp_name = ? ,
				doors = ? ,
				big_lag = ? ,
				small_lag = ? ,
				adult_place = ? ,
				condition = ? ,
				min_age = ? ,
				locations = ? ,
				car_group = ? ,
				description = ?,
				price = ?
				WHERE car_id = ?`
	RemoveCar = `DELETE FROM cars 
				WHERE car_id = ?`
	SelectRents = `SELECT rent_id,
						car_id,
						from_time,
						to_time,
						location,
						extras,
						discounts,
						rent_detail
						FROM rents`
	RemoveRent = `DELETE FROM rents 
							  WHERE rent_id = ?`
	SelectCarsRents = `SELECT car_id,
								car_comp_name ,
								doors,
								big_lag,
								small_lag,
								adult_place,
								condition,
								min_age,
								locations,
								car_group,
								description,
								price,
								rent_id,
								from_time,
								to_time,
								location,
								extras,
								discounts,
								rent_detail
								FROM cars
								LEFT JOIN rents using (car_id)`
)
