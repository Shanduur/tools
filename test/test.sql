CREATE TABLE TEST (
     A1 CHAR(1),
     A2 CHAR(2),
     A3 INTEGER,
     -- A4 TIMESTAMP,
     A4 DATE
);
INSERT INTO Customers(
          CustomerName,
          ContactName,
          Address,
          City,
          PostalCode,
          Country
     )
VALUES (
          'Cardinal',
          'Tom B. Erichsen',
          'Skagen 21',
          'Stavanger',
          '4006',
          'Norway'
     );
-- COMMENT TEST
INSERT INTO Customers(CustomerName, City, Country)
VALUES ('Cardinal', 'Stavanger', 'Norway');
/* COMMENT TEST */
INSERT INTO Customers(CustomerName, City, Country)
VALUES ('Cardinal', '--Stavanger', 'Norway');
INSERT INTO Customers(CustomerName, City, CountryIsSuper)
VALUES ('Cardinal', 'Stavanger', 'Norway');
SELECT *
FROM TEST;
FORMAT ME TEST;
CREATE TABLE TEST (
     A1 CHAR(1),
     A2 CHAR(2),
     A3 INTEGER,
     -- A4 TIMESTAMP,
     A4 DATE
);
INSERT INTO Customers(CustomerName, City, CountryIsSuper)
VALUES ('Cardinal', 'Sta, vanger', 'No,rway');