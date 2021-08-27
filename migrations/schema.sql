--
-- PostgreSQL database dump
--

-- Dumped from database version 13.3
-- Dumped by pg_dump version 13.3

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: documents; Type: TABLE; Schema: public; Owner: ensetse_user
--

CREATE TABLE public.documents (
    id uuid NOT NULL,
    doc_name character varying(255) NOT NULL,
    student_id character varying(255) NOT NULL,
    is_done boolean DEFAULT false NOT NULL,
    status character varying(255) DEFAULT 'P'::character varying NOT NULL,
    doc_path character varying(255) NOT NULL,
    message character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.documents OWNER TO ensetse_user;

--
-- Name: schema_migration; Type: TABLE; Schema: public; Owner: ensetse_user
--

CREATE TABLE public.schema_migration (
    version character varying(14) NOT NULL
);


ALTER TABLE public.schema_migration OWNER TO ensetse_user;

--
-- Name: documents documents_pkey; Type: CONSTRAINT; Schema: public; Owner: ensetse_user
--

ALTER TABLE ONLY public.documents
    ADD CONSTRAINT documents_pkey PRIMARY KEY (id);


--
-- Name: schema_migration_version_idx; Type: INDEX; Schema: public; Owner: ensetse_user
--

CREATE UNIQUE INDEX schema_migration_version_idx ON public.schema_migration USING btree (version);


--
-- PostgreSQL database dump complete
--
