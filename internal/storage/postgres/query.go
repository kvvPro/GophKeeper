package postgres

func getInitQuery() string {
	return `
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
			foreign key (owner) references users
	);

	create table if not exists auth_info
	(
		id       integer not null,
		login    varchar not null,
		password varchar not null,
		constraint auth_info_data_id_fk
			foreign key (id) references data
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
			foreign key (owner) references data
	);

	create table if not exists textinfo
	(
		id   integer not null,
		info text    not null,
		constraint textinfo_data_id_fk
			foreign key (id) references data
	);

	create table if not exists bindata
	(
		id   integer,
		data bytea not null,
		constraint bindata_data_id_fk
			foreign key (id) references data
	);

	create table if not exists cards
	(
		id          integer,
		card_number varchar not null,
		pin         varchar not null,
		cvc         varchar not null,
		constraint cards_data_id_fk
			foreign key (id) references data
	);
	`
}

func addUserQuery() string {
	return `
	insert into public.users (login, password)
	values ($1, $2);
	`
}

func getUserQuery() string {
	return `
	select login, password
	from public.users
	where
		login = $1
	`
}

func getUserDataQuery() string {
	return `
	select auth_info.id,
			auth_info.login, 
			auth_info.password
	from auth_info 
		inner join data data_info 
		on data_info.id = auth_info.id
	where 
		data_info.key = $1 and data_info.owner = $2
	`
}

func updateUserDataQuery() string {
	return `
	update public.auth_info
	set login=$1, password=$2
	where 
		id in (select 
					data_info.id 
				from 
					data data_info 
				where 
					data_info.key=$3 
					and data_info.owner=$4)
	`
}

func addUserDataQuery() string {
	return `
	DO $$ 
	DECLARE
		dataid INTEGER;
	BEGIN
		insert into public.data (owner, key, data_type)
		values ($1, $2, $3)
		returning id into dataid;

		insert into public.auth_info (owner, key, value)
		values (dataid, $4, $5);

		RETURNING id INTO myid;
	END $$
	`
}

func deleteDataQuery() string {
	return `
	delete
	from public.data
	where public.data.key=$1 and public.data.owner=$2
	`
}

func getMetadataQuery() string {
	return `
	select key, value
	from metainfo
	where owner = $1
	`
}

func updateMetadataQuery() string {
	return `
	update metainfo
	set key=$1, value=$2 
	where
		owner=(select data.id 
			   from data
			   where data.owner=$3
			   		 and data.key=$4);
	`
}

func addMetadataQuery() string {
	return `
	insert into public.metainfo (owner, key, value)
	select data.id, $1, $2
	from data
	where data.owner=$3 and data.key=$4;
	`
}
