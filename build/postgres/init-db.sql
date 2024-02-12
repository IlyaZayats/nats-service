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

-- select json_build_object(
--                'order_uid', W.order_uid,
--                'track_number', W.track_number,
--                'entry', W.entry,
--                'locale', W.locale,
--                'internal_signature', W.internal_signature,
--                'customer_id', W.customer_id,
--                'delivery_service', W.delivery_service,
--                'shard_key', W.shard_key,
--                'sm_id', W.sm_id,
--                'date_created', W.date_created,
--                'oof_shard', W.oof_shard,
--                'delivery', (select row_to_json(x) from (select D.name, D.phone, D.zip, D.city, D.address, D.region, D.email from Delivery as D where D.order_uid = W.order_uid) as x),
--                'payment', (select row_to_json(x) from (select P.order_uid as transaction, P.request_id, P.currency, P.provider, P.amount, P.payment_dt, P.Bank, P.delivery_cost, P.goods_total, P.custom_fee from Payment as P where P.order_uid = W.order_uid) as x),
--                'items', array(select row_to_json(x) from (select I.track_number, I.rid, I.name, I.size, I.brand, I.chrt_id, I.price, I.sale, I.total_price, I.nm_id, i.status from Item as I where I.track_number = W.track_number) as x)
--        ) from WBOrder as W;
--
-- select row_to_json(x2) from (select * from Delivery) as x2;
--
-- select json_build_object('items', array(select row_to_json(x1) from (select Item.track_number, Item.name from Item limit 1) as x1));
--
-- select row_to_json(x) from (select P.order_uid as transaction, P.request_id, P.currency, P.provider, P.Bank, P.amount, P.payment_dt, P.delivery_cost, P.goods_total, P.custom_fee from Payment as P) as x;