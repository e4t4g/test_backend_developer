package main

const (
	FIND_BY_SQUARE = `
		SELECT  
			id, 
			name, 
			website, 
			ST_AsText(coordinates), 
			description, 
			rating
		FROM 
			"MY_TABLE"
		WHERE 
			ST_DWithin(coordinates::geography, ST_MakeEnvelope(%f, %f, %f, %f, 4326), 0)
	ORDER BY 				CASE
					WHEN ST_DWithin(coordinates::geography, 'SRID=4326;POINT(%f %f)', 50)
					THEN -rating
					ELSE ST_Distance(coordinates::geography, 'SRID=4326;POINT(%f %f)')
				END;
`
)

const (
	FIND_BY_CIRCLE = `
		SELECT  
			id, 
			name, 
			website, 
			ST_AsText(coordinates), 
			description, 
			rating
		FROM 
			"MY_TABLE"
		WHERE 
			ST_DWithin(coordinates::geography, ST_GeographyFromText('SRID=4326;POINT(%f %f)'), %f)
	ORDER BY 				CASE
					WHEN ST_DWithin(coordinates::geography, 'SRID=4326;POINT(%f %f)', 50)
					THEN -rating
					ELSE ST_Distance(coordinates::geography, 'SRID=4326;POINT(%f %f)')
				END;
`
)
