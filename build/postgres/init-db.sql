create table WBOrder(
    order_uid varchar primary key,
    track_number varchar unique,
    entry varchar,
    locale varchar,
    internal_signature varchar,
    customer_id varchar,
    delivery_service varchar,
    shard_key varchar,
    sm_id int,
    date_created varchar,
    oof_shard varchar
);

create table Delivery(
    id serial not null primary key,
    order_uid varchar not null references WBOrder(order_uid) on delete cascade,
    name varchar,
    phone varchar,
    zip varchar,
    city varchar,
    address varchar,
    region varchar,
    email varchar
);

create table Payment(
    id serial not null primary key,
    order_uid varchar not null references WBOrder(order_uid) on delete cascade,
    request_id varchar,
    currency varchar,
    provider varchar,
    bank varchar,
    amount int,
    payment_dt int,
    delivery_cost int,
    goods_total int,
    custom_fee int
);

create table Item(
    id serial not null primary key,
    track_number varchar references WBOrder(track_number) on delete cascade,
    rid varchar,
    name varchar,
    size varchar,
    brand varchar,
    chrt_id int,
    price int,
    sale int,
    total_price int,
    nm_id int,
    status int
);

insert into WBOrder (order_uid, track_number, entry) VALUES ('test1', 'wbtest1', 'lmao1'), ('test2', 'wbtest2', 'lmao2');
insert into Delivery (order_uid, name) VALUES ('test1', 'meme1'), ('test2', 'meme2');
insert into Payment (order_uid, amount) VALUES ('test1', 1), ('test2', 2);
insert into Item (track_number, name) VALUES ('wbtest1', 'nice1_1'), ('wbtest1', 'nice1_2'), ('wbtest1', 'nice1_3'), ('wbtest2', 'nice2_1');