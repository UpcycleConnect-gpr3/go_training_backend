CREATE TABLE IF NOT EXISTS TRAINING_TRAINING_CONTENT (
    training_id INT NOT NULL,
    training_content_id INT NOT NULL,
    PRIMARY KEY (training_id, training_content_id),
    FOREIGN KEY (training_id) REFERENCES TRAININGS(id) ON DELETE CASCADE,
    FOREIGN KEY (training_content_id) REFERENCES TRAINING_CONTENT(id) ON DELETE CASCADE
);
