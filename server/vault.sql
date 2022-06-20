create table salts (
    salt_id varchar(255) NOT NULL,
    salt varchar(255) NOT NULL,
    PRIMARY KEY (salt_id)
);

create table objects (
    object_id varchar(255) NOT NULL,
    object_data varchar(65536) NOT NULL,
    PRIMARY KEY (object_id)
);

create table authenticated_objects (
    object_id varchar(255) NOT NULL,
    object_data varchar(65536) NOT NULL,
    PRIMARY KEY (object_id)
)