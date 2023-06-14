CREATE TABLE chat_logs (

    id  INT NOT NULL AUTO_INCREMENT,
    chat VARCHAR(101) NOT NULL,
    message VARCHAR(500) NOT NULL,
	sender VARCHAR(50) NOT NULL,
    date_time INT NOT NULL,
    PRIMARY KEY (id)

)