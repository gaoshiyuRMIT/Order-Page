create table if not exists orders (
    id integer not null primary key,
    created_at timestamp not null,
    order_name varchar(256) not null,
    customer_id varchar(256) not null
);

create table if not exists order_items (
    id integer not null primary key,
    order_id integer not null references orders(id),
    price_per_unit decimal(20, 6) constraint price_greater_zero check (price_per_unit > 0),
    quantity integer constraint quantity_greater_zero check (quantity > 0),
    product varchar(256) not null
);

create table if not exists deliveries (
    id integer not null primary key,
    order_item_id integer not null references order_items(id),
    delivered_quantity integer not null constraint dlv_quantity_greater_zero check (delivered_quantity >= 0)
);