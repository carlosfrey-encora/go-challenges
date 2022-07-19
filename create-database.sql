DROP TABLE IF EXISTS task;
CREATE TABLE task (
  id         INT AUTO_INCREMENT NOT NULL,
  Name       VARCHAR(128) NOT NULL,
  Completed  BOOLEAN NOT NULL,
  
  PRIMARY KEY (`id`)
);

INSERT INTO task
  (Name, Completed)
VALUES
  ('Wash my dishes', false),
  ('Go out for jogging', false);


  