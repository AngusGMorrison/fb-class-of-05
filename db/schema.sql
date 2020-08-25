--
-- PostgreSQL database dump
--

-- Dumped from database version 12.3
-- Dumped by pg_dump version 12.3

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

--
-- Name: friend_request_status; Type: TYPE; Schema: public; Owner: fb05_dev
--

CREATE TYPE public.friend_request_status AS ENUM (
    'accepted',
    'declined',
    'pending'
);


ALTER TYPE public.friend_request_status OWNER TO fb05_dev;

--
-- Name: interest_category; Type: TYPE; Schema: public; Owner: fb05_dev
--

CREATE TYPE public.interest_category AS ENUM (
    'general',
    'clubs_and_jobs',
    'movies',
    'music',
    'books'
);


ALTER TYPE public.interest_category OWNER TO fb05_dev;

--
-- Name: interested_in; Type: TYPE; Schema: public; Owner: fb05_dev
--

CREATE TYPE public.interested_in AS ENUM (
    'men',
    'women',
    'both'
);


ALTER TYPE public.interested_in OWNER TO fb05_dev;

--
-- Name: looking_for; Type: TYPE; Schema: public; Owner: fb05_dev
--

CREATE TYPE public.looking_for AS ENUM (
    'friendship',
    'dating',
    'relationship',
    'random_play',
    'whatever'
);


ALTER TYPE public.looking_for OWNER TO fb05_dev;

--
-- Name: political_views; Type: TYPE; Schema: public; Owner: fb05_dev
--

CREATE TYPE public.political_views AS ENUM (
    'very_conservative',
    'conservative',
    'moderate',
    'liberal',
    'very_liberal'
);


ALTER TYPE public.political_views OWNER TO fb05_dev;

--
-- Name: relationship_status; Type: TYPE; Schema: public; Owner: fb05_dev
--

CREATE TYPE public.relationship_status AS ENUM (
    'single',
    'relationship',
    'engaged',
    'married',
    'complicated',
    'open'
);


ALTER TYPE public.relationship_status OWNER TO fb05_dev;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: fb05_dev
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO fb05_dev;

--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: fb05_dev
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- PostgreSQL database dump complete
--

