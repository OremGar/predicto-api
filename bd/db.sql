CREATE TABLE public.usuarios (
    id integer NOT NULL,
    nombre character varying NOT NULL,
    apellidos character varying NOT NULL,
    usuario character varying NOT NULL,
    correo character varying NOT NULL,
    contrasena character varying NOT NULL,
    telefono character varying NOT NULL
);

CREATE SEQUENCE public.usuarios_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.usuarios_id_seq OWNED BY public.usuarios.id;
ALTER TABLE ONLY public.usuarios ALTER COLUMN id SET DEFAULT nextval('public.usuarios_id_seq'::regclass);



CREATE TABLE public.usuarios_jwt (
    id_usuario integer NOT NULL,
    token character varying NOT NULL,
    fecha_inicio timestamp without time zone NOT NULL
);


SELECT pg_catalog.setval('public.usuarios_id_seq', 10, true);


ALTER TABLE ONLY public.usuarios
    ADD CONSTRAINT usuarios_correo_unique UNIQUE (correo);


ALTER TABLE ONLY public.usuarios
    ADD CONSTRAINT usuarios_pkey PRIMARY KEY (id);


ALTER TABLE ONLY public.usuarios
    ADD CONSTRAINT usuarios_usuario_unique UNIQUE (usuario);

ALTER TABLE ONLY public.usuarios_jwt
    ADD CONSTRAINT fk_idusuario FOREIGN KEY (id_usuario) REFERENCES public.usuarios(id) ON DELETE CASCADE;

