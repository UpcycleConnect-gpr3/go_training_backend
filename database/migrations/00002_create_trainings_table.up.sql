CREATE TABLE IF NOT EXISTS TRAININGS (
    id INT AUTO_INCREMENT PRIMARY KEY,
    type VARCHAR(255),
    name VARCHAR(255) NOT NULL,
    mode_of_delivery VARCHAR(255),
    duration VARCHAR(255),
    target_audience VARCHAR(255),
    minimum_number_of_participants INT DEFAULT 0,
    maximum_number_of_participants INT DEFAULT 0,
    location VARCHAR(255),
    trainer_profile TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)
