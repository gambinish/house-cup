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