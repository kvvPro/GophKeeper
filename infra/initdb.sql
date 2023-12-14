create table if not exists users
(
    login    varchar not null,
    password varchar not null,
    constraint users_pk
        primary key (login)
);

create table if not exists data
(
    id        integer generated always as identity (minvalue 0),
    key       varchar not null,
    owner     varchar not null,
    data_type varchar not null,
    constraint data_pk
        primary key (id),
    constraint data_pk2
        unique (key, owner),
    constraint data_users_login_fk
        foreign key (owner) references users ON DELETE CASCADE
);

create table if not exists auth_info
(
    id       integer not null,
    login    varchar not null,
    password varchar not null,
    constraint auth_info_data_id_fk
        foreign key (id) references data ON DELETE CASCADE
);

create table if not exists metainfo
(
    id    integer generated always as identity (minvalue 0),
    owner integer not null,
    key   varchar not null,
    value varchar not null,
    constraint metainfo_pk
        primary key (id),
    constraint metainfo_data_id_fk
        foreign key (owner) references data ON DELETE CASCADE
);

create table if not exists textinfo
(
    id   integer not null,
    info text    not null,
    constraint textinfo_data_id_fk
        foreign key (id) references data ON DELETE CASCADE
);

create table if not exists bindata
(
    id   integer,
    data bytea not null,
    constraint bindata_data_id_fk
        foreign key (id) references data ON DELETE CASCADE
);

create table if not exists cards
(
    id          integer,
    card_number varchar not null,
    pin         varchar not null,
    cvc         varchar not null,
    constraint cards_data_id_fk
        foreign key (id) references data ON DELETE CASCADE
);

create table if not exists data_change_history
(
    id         integer   not null,
    updated_at timestamp not null,
    constraint data_change_history_data_id_fk
        foreign key (id) references data ON DELETE CASCADE
);

