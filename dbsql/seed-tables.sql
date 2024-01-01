DROP TABLE IF EXISTS Points;
DROP TABLE IF EXISTS Students;
DROP TABLE IF EXISTS Houses;
DROP TABLE IF EXISTS Tournaments;

-- Create Tournament table
CREATE TABLE Tournaments (
    ID INT AUTO_INCREMENT PRIMARY KEY,
    Tournament_Name VARCHAR(255) NOT NULL,
    Created_At TIMESTAMP NOT NULL,
    Ended_At TIMESTAMP
);

-- Create House table
CREATE TABLE Houses (
    ID INT AUTO_INCREMENT PRIMARY KEY,
    House_Name VARCHAR(255) NOT NULL,
    House_Points INT NOT NULL,
    Tournament_ID INT,
    FOREIGN KEY (Tournament_ID) REFERENCES Tournaments(ID)
);

-- Create Student table
CREATE TABLE Students (
    ID INT AUTO_INCREMENT PRIMARY KEY,
    Student_Name VARCHAR(255) NOT NULL,
    Points INT NOT NULL,
    House_ID INT,
    FOREIGN KEY (House_ID) REFERENCES Houses(ID)
);

-- Create Point table
CREATE TABLE Points (
    ID INT AUTO_INCREMENT PRIMARY KEY,
    Points INT NOT NULL,
    Notes VARCHAR(255),
    Student_ID INT,
    FOREIGN KEY (Student_ID) REFERENCES Students(ID),
    House_ID INT NOT NULL,
    FOREIGN KEY (House_ID) REFERENCES Houses(ID)
);
-- Seed Data

-- Insert seed data for Tournament table
INSERT INTO Tournaments
  (Tournament_Name, Created_At, Ended_At)
VALUES
  ('Spring Tournament', '2023-01-01 00:00:00', NULL),
  ('Summer Tournament', '2023-05-15 12:30:00', NULL);

-- Insert seed data for House table
INSERT INTO Houses
  (House_Name, House_Points, Tournament_ID)
VALUES
  ('Gryffindor', 23, 1),
  ('Slytherin', 7, 1),
  ('Ravenclaw', 10, 1),
  ('Hufflepuff', 11, 1),
  ('Durmstrang', 8, 2),
  ('Beauxbatons', 7, 2);

-- Insert seed data for Student table
INSERT INTO Students
  (Student_Name, Points, House_ID)
VALUES
  ('Harry Potter', 10, 1),
  ('Hermione Granger', 5, 1),
  ('Ron Weasley', 8, 1),
  ('Draco Malfoy', 7, 2),
  ('Luna Lovegood', 10, 3),
  ('Cedric Diggory', 11, 4),
  ('Fleur Delacour', 7, 6),
  ('Viktor Krum', 8, 5),
  ('Cho Chang', 0, 3);

-- Insert seed data for Points table
INSERT INTO Points
  (Points, Notes, Student_ID, House_ID)
VALUES
  (10, 'Quidditch match victory', 1, 1),
  (5, 'Excellent potion brewing', 2, 1),
  (8, 'Prefect duties', 3, 1),
  (7, 'Slytherin common room points', 4, 2),
  (10, 'Outstanding in Charms class', 5, 3),
  (5, 'Herbology achievement', 6, 4),
  (8, 'Durmstrang team victory', 8, 5),
  (7, 'Beauxbatons team victory', 7, 6),
  (6, 'Participation in Triwizard Tournament', 6, 4);