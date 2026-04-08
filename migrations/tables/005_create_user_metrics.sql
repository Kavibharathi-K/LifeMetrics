create table user_metrics (
    id serial primary key,
    age int not null,
    gender text not null,
    height_cm decimal(5,2) not null,
    weight_kg decimal(5,2) not null,
    activity_level text not null,
    maintenance_calories int not null,
    protein_goal int not null,
    carb_goal int not null,
    fat_goal int not null,
    created_at timestamp default now()
);