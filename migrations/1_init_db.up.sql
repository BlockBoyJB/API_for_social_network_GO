create table if not exists public.user
(
    id         serial primary key,
    username   varchar unique not null,
    first_name varchar        not null,
    last_name  varchar        not null,
    email      varchar unique not null,
    password   varchar        not null
);
create table if not exists public.post
(
    id       serial PRIMARY KEY,
    username varchar references public.user (username) on delete cascade on update cascade,
    post_id  varchar unique not null,
    title    varchar        not null,
    text     text           not null
);
create table if not exists public.reaction
(
    id          serial primary key,
    post_id     varchar        not null references public.post (post_id) on delete cascade,
    reaction_id varchar unique not null,
    reaction    varchar        not null
);
create table if not exists public.comment
(
    id         serial primary key,
    username   varchar        not null references public.user (username) on delete cascade on update cascade,
    post_id    varchar        not null references public.post (post_id) on delete cascade,
    comment_id varchar unique not null,
    comment    varchar        not null
)
