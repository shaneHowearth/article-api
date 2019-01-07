CREATE TABLE IF NOT EXISTS article (
 id		SERIAL PRIMARY KEY,
 title	TEXT,
 pub_date	TIMESTAMP,
 body	TEXT,
 tags	TEXT[]
)
