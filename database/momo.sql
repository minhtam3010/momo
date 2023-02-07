create table if not exists momo_payment (
    id int not null auto_increment primary key,
    order_id text not null,
    request_id text not null,
    trans_id text not null,
    amount int not null,
    pay_type text not null,
    date_created timestamp not null default current_timestamp,
    date_updated timestamp
);