
CREATE TABLE students (
    id CHAR(36) NOT NULL,         
    name VARCHAR(255) NOT NULL,   
    age INT NOT NULL,            
    created_by CHAR(36) NOT NULL, 
    created_on DATETIME NOT NULL, 
    updated_by CHAR(36),          
    updated_on DATETIME,          
    PRIMARY KEY (id),            
    FOREIGN KEY (created_by) REFERENCES users(id),  
    FOREIGN KEY (updated_by) REFERENCES users(id)  
);

CREATE TABLE users (
    id CHAR(36) NOT NULL,         
    email VARCHAR(255) NOT NULL,  
    password VARCHAR(255) NOT NULL, 
    PRIMARY KEY (id),             
    UNIQUE (email)               
);
