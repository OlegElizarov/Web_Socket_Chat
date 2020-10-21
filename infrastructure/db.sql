CREATE TABLE message (
    id       bigserial primary key,
    sender integer,
    reciever integer,
    message_line text,
    shown bool
);