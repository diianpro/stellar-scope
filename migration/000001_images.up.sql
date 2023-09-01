CREATE TABLE IF NOT EXISTS image
(
    "date"          date PRIMARY KEY,
    title           varchar(255),
    explanation     text,
    image_extension varchar(16),
    copyright       varchar(255)
);

