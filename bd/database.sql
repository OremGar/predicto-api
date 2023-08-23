-- SE CREA LA TABLA DE USUARIOS
CREATE TABLE public.usuarios
(
    id serial,
    nombre character varying NOT NULL,
    apellidos character varying NOT NULL,
    usuario character varying NOT NULL,
    correo character varying NOT NULL,
    contrasena character varying NOT NULL,
    telefono character varying NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT usuarios_usuario_unique UNIQUE (usuario),
    CONSTRAINT usuarios_correo_unique UNIQUE (correo)
)
WITH (
    OIDS = FALSE
);

ALTER TABLE IF EXISTS public.usuarios
    OWNER to oaorlsjq;

-- SE CREA LA TABLA DE RELACIÓN ENTRE USUARIOS Y SUS JWT
CREATE TABLE public.usuarios_jwt
(
    id_usuario integer NOT NULL,
    token character varying NOT NULL,
    fecha_inicio timestamp without time zone NOT NULL,
    CONSTRAINT fk_idusuario FOREIGN KEY (id_usuario)
        REFERENCES public.usuarios (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
)
WITH (
    OIDS = FALSE
);

ALTER TABLE IF EXISTS public.usuarios_jwt
    OWNER to oaorlsjq;