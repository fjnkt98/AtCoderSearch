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
-- Name: batch_histories; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.batch_histories (
    id bigint NOT NULL,
    name text NOT NULL,
    started_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    finished_at timestamp with time zone,
    status text DEFAULT 'working'::text NOT NULL,
    options json
);


--
-- Name: batch_histories_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.batch_histories_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: batch_histories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.batch_histories_id_seq OWNED BY public.batch_histories.id;


--
-- Name: contests; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.contests (
    contest_id text NOT NULL,
    start_epoch_second bigint NOT NULL,
    duration_second bigint NOT NULL,
    title text NOT NULL,
    rate_change text NOT NULL,
    category text NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: difficulties; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.difficulties (
    problem_id text NOT NULL,
    slope double precision,
    intercept double precision,
    variance double precision,
    difficulty bigint,
    discrimination double precision,
    irt_loglikelihood double precision,
    irt_users double precision,
    is_experimental boolean,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: languages; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.languages (
    language text NOT NULL,
    "group" text
);


--
-- Name: problems; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.problems (
    problem_id text NOT NULL,
    contest_id text NOT NULL,
    problem_index text NOT NULL,
    name text NOT NULL,
    title text NOT NULL,
    url text NOT NULL,
    html text NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version character varying(128) NOT NULL
);


--
-- Name: submission_crawl_histories; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submission_crawl_histories (
    id bigint NOT NULL,
    contest_id text NOT NULL,
    started_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    status text DEFAULT 'working'::text NOT NULL,
    finished_at timestamp with time zone
);


--
-- Name: submission_crawl_histories_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.submission_crawl_histories_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: submission_crawl_histories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.submission_crawl_histories_id_seq OWNED BY public.submission_crawl_histories.id;


--
-- Name: submissions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions (
    id bigint NOT NULL,
    epoch_second bigint NOT NULL,
    problem_id text NOT NULL,
    contest_id text,
    user_id text,
    language text,
    point double precision,
    length integer,
    result text,
    execution_time integer,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
)
PARTITION BY RANGE (epoch_second);


--
-- Name: submissions_2010; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_2010 (
    id bigint NOT NULL,
    epoch_second bigint NOT NULL,
    problem_id text NOT NULL,
    contest_id text,
    user_id text,
    language text,
    point double precision,
    length integer,
    result text,
    execution_time integer,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: submissions_2011; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_2011 (
    id bigint NOT NULL,
    epoch_second bigint NOT NULL,
    problem_id text NOT NULL,
    contest_id text,
    user_id text,
    language text,
    point double precision,
    length integer,
    result text,
    execution_time integer,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: submissions_2012; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_2012 (
    id bigint NOT NULL,
    epoch_second bigint NOT NULL,
    problem_id text NOT NULL,
    contest_id text,
    user_id text,
    language text,
    point double precision,
    length integer,
    result text,
    execution_time integer,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: submissions_2013; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_2013 (
    id bigint NOT NULL,
    epoch_second bigint NOT NULL,
    problem_id text NOT NULL,
    contest_id text,
    user_id text,
    language text,
    point double precision,
    length integer,
    result text,
    execution_time integer,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: submissions_2014; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_2014 (
    id bigint NOT NULL,
    epoch_second bigint NOT NULL,
    problem_id text NOT NULL,
    contest_id text,
    user_id text,
    language text,
    point double precision,
    length integer,
    result text,
    execution_time integer,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: submissions_2015; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_2015 (
    id bigint NOT NULL,
    epoch_second bigint NOT NULL,
    problem_id text NOT NULL,
    contest_id text,
    user_id text,
    language text,
    point double precision,
    length integer,
    result text,
    execution_time integer,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: submissions_2016; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_2016 (
    id bigint NOT NULL,
    epoch_second bigint NOT NULL,
    problem_id text NOT NULL,
    contest_id text,
    user_id text,
    language text,
    point double precision,
    length integer,
    result text,
    execution_time integer,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: submissions_2017; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_2017 (
    id bigint NOT NULL,
    epoch_second bigint NOT NULL,
    problem_id text NOT NULL,
    contest_id text,
    user_id text,
    language text,
    point double precision,
    length integer,
    result text,
    execution_time integer,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: submissions_2018; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_2018 (
    id bigint NOT NULL,
    epoch_second bigint NOT NULL,
    problem_id text NOT NULL,
    contest_id text,
    user_id text,
    language text,
    point double precision,
    length integer,
    result text,
    execution_time integer,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: submissions_2019; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_2019 (
    id bigint NOT NULL,
    epoch_second bigint NOT NULL,
    problem_id text NOT NULL,
    contest_id text,
    user_id text,
    language text,
    point double precision,
    length integer,
    result text,
    execution_time integer,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: submissions_2020; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_2020 (
    id bigint NOT NULL,
    epoch_second bigint NOT NULL,
    problem_id text NOT NULL,
    contest_id text,
    user_id text,
    language text,
    point double precision,
    length integer,
    result text,
    execution_time integer,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: submissions_2021; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_2021 (
    id bigint NOT NULL,
    epoch_second bigint NOT NULL,
    problem_id text NOT NULL,
    contest_id text,
    user_id text,
    language text,
    point double precision,
    length integer,
    result text,
    execution_time integer,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: submissions_2022; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_2022 (
    id bigint NOT NULL,
    epoch_second bigint NOT NULL,
    problem_id text NOT NULL,
    contest_id text,
    user_id text,
    language text,
    point double precision,
    length integer,
    result text,
    execution_time integer,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: submissions_2023; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_2023 (
    id bigint NOT NULL,
    epoch_second bigint NOT NULL,
    problem_id text NOT NULL,
    contest_id text,
    user_id text,
    language text,
    point double precision,
    length integer,
    result text,
    execution_time integer,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: submissions_2024; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_2024 (
    id bigint NOT NULL,
    epoch_second bigint NOT NULL,
    problem_id text NOT NULL,
    contest_id text,
    user_id text,
    language text,
    point double precision,
    length integer,
    result text,
    execution_time integer,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: submissions_2025; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_2025 (
    id bigint NOT NULL,
    epoch_second bigint NOT NULL,
    problem_id text NOT NULL,
    contest_id text,
    user_id text,
    language text,
    point double precision,
    length integer,
    result text,
    execution_time integer,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: submissions_2026; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_2026 (
    id bigint NOT NULL,
    epoch_second bigint NOT NULL,
    problem_id text NOT NULL,
    contest_id text,
    user_id text,
    language text,
    point double precision,
    length integer,
    result text,
    execution_time integer,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: submissions_2027; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_2027 (
    id bigint NOT NULL,
    epoch_second bigint NOT NULL,
    problem_id text NOT NULL,
    contest_id text,
    user_id text,
    language text,
    point double precision,
    length integer,
    result text,
    execution_time integer,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: submissions_2028; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_2028 (
    id bigint NOT NULL,
    epoch_second bigint NOT NULL,
    problem_id text NOT NULL,
    contest_id text,
    user_id text,
    language text,
    point double precision,
    length integer,
    result text,
    execution_time integer,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: submissions_2029; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_2029 (
    id bigint NOT NULL,
    epoch_second bigint NOT NULL,
    problem_id text NOT NULL,
    contest_id text,
    user_id text,
    language text,
    point double precision,
    length integer,
    result text,
    execution_time integer,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    user_id text NOT NULL,
    rating integer NOT NULL,
    highest_rating integer NOT NULL,
    affiliation text,
    birth_year integer,
    country text,
    crown text,
    join_count integer NOT NULL,
    rank integer NOT NULL,
    active_rank integer,
    wins integer NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: submissions_2010; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_2010 FOR VALUES FROM ('0') TO ('1293807600');


--
-- Name: submissions_2011; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_2011 FOR VALUES FROM ('1293807600') TO ('1325343600');


--
-- Name: submissions_2012; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_2012 FOR VALUES FROM ('1325343600') TO ('1356966000');


--
-- Name: submissions_2013; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_2013 FOR VALUES FROM ('1356966000') TO ('1388502000');


--
-- Name: submissions_2014; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_2014 FOR VALUES FROM ('1388502000') TO ('1420038000');


--
-- Name: submissions_2015; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_2015 FOR VALUES FROM ('1420038000') TO ('1451574000');


--
-- Name: submissions_2016; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_2016 FOR VALUES FROM ('1451574000') TO ('1483196400');


--
-- Name: submissions_2017; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_2017 FOR VALUES FROM ('1483196400') TO ('1514732400');


--
-- Name: submissions_2018; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_2018 FOR VALUES FROM ('1514732400') TO ('1546268400');


--
-- Name: submissions_2019; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_2019 FOR VALUES FROM ('1546268400') TO ('1577804400');


--
-- Name: submissions_2020; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_2020 FOR VALUES FROM ('1577804400') TO ('1609426800');


--
-- Name: submissions_2021; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_2021 FOR VALUES FROM ('1609426800') TO ('1640962800');


--
-- Name: submissions_2022; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_2022 FOR VALUES FROM ('1640962800') TO ('1672498800');


--
-- Name: submissions_2023; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_2023 FOR VALUES FROM ('1672498800') TO ('1704034800');


--
-- Name: submissions_2024; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_2024 FOR VALUES FROM ('1704034800') TO ('1735657200');


--
-- Name: submissions_2025; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_2025 FOR VALUES FROM ('1735657200') TO ('1767193200');


--
-- Name: submissions_2026; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_2026 FOR VALUES FROM ('1767193200') TO ('1798729200');


--
-- Name: submissions_2027; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_2027 FOR VALUES FROM ('1798729200') TO ('1830265200');


--
-- Name: submissions_2028; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_2028 FOR VALUES FROM ('1830265200') TO ('1861887600');


--
-- Name: submissions_2029; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_2029 FOR VALUES FROM ('1861887600') TO ('1893423600');


--
-- Name: batch_histories id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.batch_histories ALTER COLUMN id SET DEFAULT nextval('public.batch_histories_id_seq'::regclass);


--
-- Name: submission_crawl_histories id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submission_crawl_histories ALTER COLUMN id SET DEFAULT nextval('public.submission_crawl_histories_id_seq'::regclass);


--
-- Name: batch_histories batch_histories_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.batch_histories
    ADD CONSTRAINT batch_histories_pkey PRIMARY KEY (id);


--
-- Name: contests contests_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.contests
    ADD CONSTRAINT contests_pkey PRIMARY KEY (contest_id);


--
-- Name: difficulties difficulties_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.difficulties
    ADD CONSTRAINT difficulties_pkey PRIMARY KEY (problem_id);


--
-- Name: languages languages_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.languages
    ADD CONSTRAINT languages_pkey PRIMARY KEY (language);


--
-- Name: problems problems_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.problems
    ADD CONSTRAINT problems_pkey PRIMARY KEY (problem_id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: submission_crawl_histories submission_crawl_histories_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submission_crawl_histories
    ADD CONSTRAINT submission_crawl_histories_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (user_id);


--
-- Name: languages_group_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX languages_group_index ON public.languages USING btree ("group");


--
-- Name: submission_crawl_histories_contest_id_start_at_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submission_crawl_histories_contest_id_start_at_index ON public.submission_crawl_histories USING btree (contest_id, started_at);


--
-- Name: submissions_contest_id_epoch_second_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_contest_id_epoch_second_index ON ONLY public.submissions USING btree (contest_id, epoch_second);


--
-- Name: submissions_2010_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2010_contest_id_epoch_second_idx ON public.submissions_2010 USING btree (contest_id, epoch_second);


--
-- Name: submissions_contest_id_execution_time_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_contest_id_execution_time_index ON ONLY public.submissions USING btree (contest_id, execution_time);


--
-- Name: submissions_2010_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2010_contest_id_execution_time_idx ON public.submissions_2010 USING btree (contest_id, execution_time);


--
-- Name: submissions_contest_id_length_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_contest_id_length_index ON ONLY public.submissions USING btree (contest_id, length);


--
-- Name: submissions_2010_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2010_contest_id_length_idx ON public.submissions_2010 USING btree (contest_id, length);


--
-- Name: submissions_contest_id_point_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_contest_id_point_index ON ONLY public.submissions USING btree (contest_id, point);


--
-- Name: submissions_2010_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2010_contest_id_point_idx ON public.submissions_2010 USING btree (contest_id, point);


--
-- Name: submissions_unique; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_unique ON ONLY public.submissions USING btree (epoch_second, id);


--
-- Name: submissions_2010_epoch_second_id_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_2010_epoch_second_id_idx ON public.submissions_2010 USING btree (epoch_second, id);


--
-- Name: submissions_language_epoch_second_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_language_epoch_second_index ON ONLY public.submissions USING btree (language, epoch_second);


--
-- Name: submissions_2010_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2010_language_epoch_second_idx ON public.submissions_2010 USING btree (language, epoch_second);


--
-- Name: submissions_language_execution_time_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_language_execution_time_index ON ONLY public.submissions USING btree (language, execution_time);


--
-- Name: submissions_2010_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2010_language_execution_time_idx ON public.submissions_2010 USING btree (language, execution_time);


--
-- Name: submissions_language_length_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_language_length_index ON ONLY public.submissions USING btree (language, length);


--
-- Name: submissions_2010_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2010_language_length_idx ON public.submissions_2010 USING btree (language, length);


--
-- Name: submissions_language_point_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_language_point_index ON ONLY public.submissions USING btree (language, point);


--
-- Name: submissions_2010_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2010_language_point_idx ON public.submissions_2010 USING btree (language, point);


--
-- Name: submissions_problem_id_epoch_second_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_problem_id_epoch_second_index ON ONLY public.submissions USING btree (problem_id, epoch_second);


--
-- Name: submissions_2010_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2010_problem_id_epoch_second_idx ON public.submissions_2010 USING btree (problem_id, epoch_second);


--
-- Name: submissions_problem_id_execution_time_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_problem_id_execution_time_index ON ONLY public.submissions USING btree (problem_id, execution_time);


--
-- Name: submissions_2010_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2010_problem_id_execution_time_idx ON public.submissions_2010 USING btree (problem_id, execution_time);


--
-- Name: submissions_problem_id_length_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_problem_id_length_index ON ONLY public.submissions USING btree (problem_id, length);


--
-- Name: submissions_2010_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2010_problem_id_length_idx ON public.submissions_2010 USING btree (problem_id, length);


--
-- Name: submissions_problem_id_point_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_problem_id_point_index ON ONLY public.submissions USING btree (problem_id, point);


--
-- Name: submissions_2010_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2010_problem_id_point_idx ON public.submissions_2010 USING btree (problem_id, point);


--
-- Name: submissions_result_epoch_second_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_result_epoch_second_index ON ONLY public.submissions USING btree (result, epoch_second);


--
-- Name: submissions_2010_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2010_result_epoch_second_idx ON public.submissions_2010 USING btree (result, epoch_second);


--
-- Name: submissions_result_execution_time_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_result_execution_time_index ON ONLY public.submissions USING btree (result, execution_time);


--
-- Name: submissions_2010_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2010_result_execution_time_idx ON public.submissions_2010 USING btree (result, execution_time);


--
-- Name: submissions_result_length_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_result_length_index ON ONLY public.submissions USING btree (result, length);


--
-- Name: submissions_2010_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2010_result_length_idx ON public.submissions_2010 USING btree (result, length);


--
-- Name: submissions_result_point_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_result_point_index ON ONLY public.submissions USING btree (result, point);


--
-- Name: submissions_2010_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2010_result_point_idx ON public.submissions_2010 USING btree (result, point);


--
-- Name: submissions_user_id_epoch_second_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_user_id_epoch_second_index ON ONLY public.submissions USING btree (user_id, epoch_second);


--
-- Name: submissions_2010_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2010_user_id_epoch_second_idx ON public.submissions_2010 USING btree (user_id, epoch_second);


--
-- Name: submissions_user_id_execution_time_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_user_id_execution_time_index ON ONLY public.submissions USING btree (user_id, execution_time);


--
-- Name: submissions_2010_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2010_user_id_execution_time_idx ON public.submissions_2010 USING btree (user_id, execution_time);


--
-- Name: submissions_user_id_length_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_user_id_length_index ON ONLY public.submissions USING btree (user_id, length);


--
-- Name: submissions_2010_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2010_user_id_length_idx ON public.submissions_2010 USING btree (user_id, length);


--
-- Name: submissions_user_id_point_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_user_id_point_index ON ONLY public.submissions USING btree (user_id, point);


--
-- Name: submissions_2010_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2010_user_id_point_idx ON public.submissions_2010 USING btree (user_id, point);


--
-- Name: submissions_user_id_result_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_user_id_result_index ON ONLY public.submissions USING btree (user_id, result);


--
-- Name: submissions_2010_user_id_result_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2010_user_id_result_idx ON public.submissions_2010 USING btree (user_id, result);


--
-- Name: submissions_2011_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2011_contest_id_epoch_second_idx ON public.submissions_2011 USING btree (contest_id, epoch_second);


--
-- Name: submissions_2011_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2011_contest_id_execution_time_idx ON public.submissions_2011 USING btree (contest_id, execution_time);


--
-- Name: submissions_2011_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2011_contest_id_length_idx ON public.submissions_2011 USING btree (contest_id, length);


--
-- Name: submissions_2011_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2011_contest_id_point_idx ON public.submissions_2011 USING btree (contest_id, point);


--
-- Name: submissions_2011_epoch_second_id_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_2011_epoch_second_id_idx ON public.submissions_2011 USING btree (epoch_second, id);


--
-- Name: submissions_2011_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2011_language_epoch_second_idx ON public.submissions_2011 USING btree (language, epoch_second);


--
-- Name: submissions_2011_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2011_language_execution_time_idx ON public.submissions_2011 USING btree (language, execution_time);


--
-- Name: submissions_2011_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2011_language_length_idx ON public.submissions_2011 USING btree (language, length);


--
-- Name: submissions_2011_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2011_language_point_idx ON public.submissions_2011 USING btree (language, point);


--
-- Name: submissions_2011_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2011_problem_id_epoch_second_idx ON public.submissions_2011 USING btree (problem_id, epoch_second);


--
-- Name: submissions_2011_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2011_problem_id_execution_time_idx ON public.submissions_2011 USING btree (problem_id, execution_time);


--
-- Name: submissions_2011_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2011_problem_id_length_idx ON public.submissions_2011 USING btree (problem_id, length);


--
-- Name: submissions_2011_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2011_problem_id_point_idx ON public.submissions_2011 USING btree (problem_id, point);


--
-- Name: submissions_2011_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2011_result_epoch_second_idx ON public.submissions_2011 USING btree (result, epoch_second);


--
-- Name: submissions_2011_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2011_result_execution_time_idx ON public.submissions_2011 USING btree (result, execution_time);


--
-- Name: submissions_2011_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2011_result_length_idx ON public.submissions_2011 USING btree (result, length);


--
-- Name: submissions_2011_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2011_result_point_idx ON public.submissions_2011 USING btree (result, point);


--
-- Name: submissions_2011_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2011_user_id_epoch_second_idx ON public.submissions_2011 USING btree (user_id, epoch_second);


--
-- Name: submissions_2011_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2011_user_id_execution_time_idx ON public.submissions_2011 USING btree (user_id, execution_time);


--
-- Name: submissions_2011_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2011_user_id_length_idx ON public.submissions_2011 USING btree (user_id, length);


--
-- Name: submissions_2011_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2011_user_id_point_idx ON public.submissions_2011 USING btree (user_id, point);


--
-- Name: submissions_2011_user_id_result_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2011_user_id_result_idx ON public.submissions_2011 USING btree (user_id, result);


--
-- Name: submissions_2012_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2012_contest_id_epoch_second_idx ON public.submissions_2012 USING btree (contest_id, epoch_second);


--
-- Name: submissions_2012_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2012_contest_id_execution_time_idx ON public.submissions_2012 USING btree (contest_id, execution_time);


--
-- Name: submissions_2012_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2012_contest_id_length_idx ON public.submissions_2012 USING btree (contest_id, length);


--
-- Name: submissions_2012_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2012_contest_id_point_idx ON public.submissions_2012 USING btree (contest_id, point);


--
-- Name: submissions_2012_epoch_second_id_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_2012_epoch_second_id_idx ON public.submissions_2012 USING btree (epoch_second, id);


--
-- Name: submissions_2012_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2012_language_epoch_second_idx ON public.submissions_2012 USING btree (language, epoch_second);


--
-- Name: submissions_2012_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2012_language_execution_time_idx ON public.submissions_2012 USING btree (language, execution_time);


--
-- Name: submissions_2012_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2012_language_length_idx ON public.submissions_2012 USING btree (language, length);


--
-- Name: submissions_2012_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2012_language_point_idx ON public.submissions_2012 USING btree (language, point);


--
-- Name: submissions_2012_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2012_problem_id_epoch_second_idx ON public.submissions_2012 USING btree (problem_id, epoch_second);


--
-- Name: submissions_2012_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2012_problem_id_execution_time_idx ON public.submissions_2012 USING btree (problem_id, execution_time);


--
-- Name: submissions_2012_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2012_problem_id_length_idx ON public.submissions_2012 USING btree (problem_id, length);


--
-- Name: submissions_2012_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2012_problem_id_point_idx ON public.submissions_2012 USING btree (problem_id, point);


--
-- Name: submissions_2012_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2012_result_epoch_second_idx ON public.submissions_2012 USING btree (result, epoch_second);


--
-- Name: submissions_2012_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2012_result_execution_time_idx ON public.submissions_2012 USING btree (result, execution_time);


--
-- Name: submissions_2012_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2012_result_length_idx ON public.submissions_2012 USING btree (result, length);


--
-- Name: submissions_2012_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2012_result_point_idx ON public.submissions_2012 USING btree (result, point);


--
-- Name: submissions_2012_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2012_user_id_epoch_second_idx ON public.submissions_2012 USING btree (user_id, epoch_second);


--
-- Name: submissions_2012_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2012_user_id_execution_time_idx ON public.submissions_2012 USING btree (user_id, execution_time);


--
-- Name: submissions_2012_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2012_user_id_length_idx ON public.submissions_2012 USING btree (user_id, length);


--
-- Name: submissions_2012_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2012_user_id_point_idx ON public.submissions_2012 USING btree (user_id, point);


--
-- Name: submissions_2012_user_id_result_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2012_user_id_result_idx ON public.submissions_2012 USING btree (user_id, result);


--
-- Name: submissions_2013_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2013_contest_id_epoch_second_idx ON public.submissions_2013 USING btree (contest_id, epoch_second);


--
-- Name: submissions_2013_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2013_contest_id_execution_time_idx ON public.submissions_2013 USING btree (contest_id, execution_time);


--
-- Name: submissions_2013_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2013_contest_id_length_idx ON public.submissions_2013 USING btree (contest_id, length);


--
-- Name: submissions_2013_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2013_contest_id_point_idx ON public.submissions_2013 USING btree (contest_id, point);


--
-- Name: submissions_2013_epoch_second_id_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_2013_epoch_second_id_idx ON public.submissions_2013 USING btree (epoch_second, id);


--
-- Name: submissions_2013_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2013_language_epoch_second_idx ON public.submissions_2013 USING btree (language, epoch_second);


--
-- Name: submissions_2013_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2013_language_execution_time_idx ON public.submissions_2013 USING btree (language, execution_time);


--
-- Name: submissions_2013_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2013_language_length_idx ON public.submissions_2013 USING btree (language, length);


--
-- Name: submissions_2013_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2013_language_point_idx ON public.submissions_2013 USING btree (language, point);


--
-- Name: submissions_2013_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2013_problem_id_epoch_second_idx ON public.submissions_2013 USING btree (problem_id, epoch_second);


--
-- Name: submissions_2013_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2013_problem_id_execution_time_idx ON public.submissions_2013 USING btree (problem_id, execution_time);


--
-- Name: submissions_2013_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2013_problem_id_length_idx ON public.submissions_2013 USING btree (problem_id, length);


--
-- Name: submissions_2013_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2013_problem_id_point_idx ON public.submissions_2013 USING btree (problem_id, point);


--
-- Name: submissions_2013_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2013_result_epoch_second_idx ON public.submissions_2013 USING btree (result, epoch_second);


--
-- Name: submissions_2013_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2013_result_execution_time_idx ON public.submissions_2013 USING btree (result, execution_time);


--
-- Name: submissions_2013_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2013_result_length_idx ON public.submissions_2013 USING btree (result, length);


--
-- Name: submissions_2013_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2013_result_point_idx ON public.submissions_2013 USING btree (result, point);


--
-- Name: submissions_2013_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2013_user_id_epoch_second_idx ON public.submissions_2013 USING btree (user_id, epoch_second);


--
-- Name: submissions_2013_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2013_user_id_execution_time_idx ON public.submissions_2013 USING btree (user_id, execution_time);


--
-- Name: submissions_2013_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2013_user_id_length_idx ON public.submissions_2013 USING btree (user_id, length);


--
-- Name: submissions_2013_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2013_user_id_point_idx ON public.submissions_2013 USING btree (user_id, point);


--
-- Name: submissions_2013_user_id_result_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2013_user_id_result_idx ON public.submissions_2013 USING btree (user_id, result);


--
-- Name: submissions_2014_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2014_contest_id_epoch_second_idx ON public.submissions_2014 USING btree (contest_id, epoch_second);


--
-- Name: submissions_2014_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2014_contest_id_execution_time_idx ON public.submissions_2014 USING btree (contest_id, execution_time);


--
-- Name: submissions_2014_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2014_contest_id_length_idx ON public.submissions_2014 USING btree (contest_id, length);


--
-- Name: submissions_2014_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2014_contest_id_point_idx ON public.submissions_2014 USING btree (contest_id, point);


--
-- Name: submissions_2014_epoch_second_id_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_2014_epoch_second_id_idx ON public.submissions_2014 USING btree (epoch_second, id);


--
-- Name: submissions_2014_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2014_language_epoch_second_idx ON public.submissions_2014 USING btree (language, epoch_second);


--
-- Name: submissions_2014_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2014_language_execution_time_idx ON public.submissions_2014 USING btree (language, execution_time);


--
-- Name: submissions_2014_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2014_language_length_idx ON public.submissions_2014 USING btree (language, length);


--
-- Name: submissions_2014_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2014_language_point_idx ON public.submissions_2014 USING btree (language, point);


--
-- Name: submissions_2014_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2014_problem_id_epoch_second_idx ON public.submissions_2014 USING btree (problem_id, epoch_second);


--
-- Name: submissions_2014_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2014_problem_id_execution_time_idx ON public.submissions_2014 USING btree (problem_id, execution_time);


--
-- Name: submissions_2014_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2014_problem_id_length_idx ON public.submissions_2014 USING btree (problem_id, length);


--
-- Name: submissions_2014_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2014_problem_id_point_idx ON public.submissions_2014 USING btree (problem_id, point);


--
-- Name: submissions_2014_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2014_result_epoch_second_idx ON public.submissions_2014 USING btree (result, epoch_second);


--
-- Name: submissions_2014_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2014_result_execution_time_idx ON public.submissions_2014 USING btree (result, execution_time);


--
-- Name: submissions_2014_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2014_result_length_idx ON public.submissions_2014 USING btree (result, length);


--
-- Name: submissions_2014_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2014_result_point_idx ON public.submissions_2014 USING btree (result, point);


--
-- Name: submissions_2014_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2014_user_id_epoch_second_idx ON public.submissions_2014 USING btree (user_id, epoch_second);


--
-- Name: submissions_2014_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2014_user_id_execution_time_idx ON public.submissions_2014 USING btree (user_id, execution_time);


--
-- Name: submissions_2014_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2014_user_id_length_idx ON public.submissions_2014 USING btree (user_id, length);


--
-- Name: submissions_2014_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2014_user_id_point_idx ON public.submissions_2014 USING btree (user_id, point);


--
-- Name: submissions_2014_user_id_result_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2014_user_id_result_idx ON public.submissions_2014 USING btree (user_id, result);


--
-- Name: submissions_2015_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2015_contest_id_epoch_second_idx ON public.submissions_2015 USING btree (contest_id, epoch_second);


--
-- Name: submissions_2015_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2015_contest_id_execution_time_idx ON public.submissions_2015 USING btree (contest_id, execution_time);


--
-- Name: submissions_2015_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2015_contest_id_length_idx ON public.submissions_2015 USING btree (contest_id, length);


--
-- Name: submissions_2015_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2015_contest_id_point_idx ON public.submissions_2015 USING btree (contest_id, point);


--
-- Name: submissions_2015_epoch_second_id_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_2015_epoch_second_id_idx ON public.submissions_2015 USING btree (epoch_second, id);


--
-- Name: submissions_2015_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2015_language_epoch_second_idx ON public.submissions_2015 USING btree (language, epoch_second);


--
-- Name: submissions_2015_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2015_language_execution_time_idx ON public.submissions_2015 USING btree (language, execution_time);


--
-- Name: submissions_2015_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2015_language_length_idx ON public.submissions_2015 USING btree (language, length);


--
-- Name: submissions_2015_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2015_language_point_idx ON public.submissions_2015 USING btree (language, point);


--
-- Name: submissions_2015_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2015_problem_id_epoch_second_idx ON public.submissions_2015 USING btree (problem_id, epoch_second);


--
-- Name: submissions_2015_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2015_problem_id_execution_time_idx ON public.submissions_2015 USING btree (problem_id, execution_time);


--
-- Name: submissions_2015_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2015_problem_id_length_idx ON public.submissions_2015 USING btree (problem_id, length);


--
-- Name: submissions_2015_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2015_problem_id_point_idx ON public.submissions_2015 USING btree (problem_id, point);


--
-- Name: submissions_2015_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2015_result_epoch_second_idx ON public.submissions_2015 USING btree (result, epoch_second);


--
-- Name: submissions_2015_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2015_result_execution_time_idx ON public.submissions_2015 USING btree (result, execution_time);


--
-- Name: submissions_2015_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2015_result_length_idx ON public.submissions_2015 USING btree (result, length);


--
-- Name: submissions_2015_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2015_result_point_idx ON public.submissions_2015 USING btree (result, point);


--
-- Name: submissions_2015_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2015_user_id_epoch_second_idx ON public.submissions_2015 USING btree (user_id, epoch_second);


--
-- Name: submissions_2015_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2015_user_id_execution_time_idx ON public.submissions_2015 USING btree (user_id, execution_time);


--
-- Name: submissions_2015_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2015_user_id_length_idx ON public.submissions_2015 USING btree (user_id, length);


--
-- Name: submissions_2015_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2015_user_id_point_idx ON public.submissions_2015 USING btree (user_id, point);


--
-- Name: submissions_2015_user_id_result_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2015_user_id_result_idx ON public.submissions_2015 USING btree (user_id, result);


--
-- Name: submissions_2016_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2016_contest_id_epoch_second_idx ON public.submissions_2016 USING btree (contest_id, epoch_second);


--
-- Name: submissions_2016_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2016_contest_id_execution_time_idx ON public.submissions_2016 USING btree (contest_id, execution_time);


--
-- Name: submissions_2016_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2016_contest_id_length_idx ON public.submissions_2016 USING btree (contest_id, length);


--
-- Name: submissions_2016_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2016_contest_id_point_idx ON public.submissions_2016 USING btree (contest_id, point);


--
-- Name: submissions_2016_epoch_second_id_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_2016_epoch_second_id_idx ON public.submissions_2016 USING btree (epoch_second, id);


--
-- Name: submissions_2016_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2016_language_epoch_second_idx ON public.submissions_2016 USING btree (language, epoch_second);


--
-- Name: submissions_2016_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2016_language_execution_time_idx ON public.submissions_2016 USING btree (language, execution_time);


--
-- Name: submissions_2016_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2016_language_length_idx ON public.submissions_2016 USING btree (language, length);


--
-- Name: submissions_2016_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2016_language_point_idx ON public.submissions_2016 USING btree (language, point);


--
-- Name: submissions_2016_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2016_problem_id_epoch_second_idx ON public.submissions_2016 USING btree (problem_id, epoch_second);


--
-- Name: submissions_2016_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2016_problem_id_execution_time_idx ON public.submissions_2016 USING btree (problem_id, execution_time);


--
-- Name: submissions_2016_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2016_problem_id_length_idx ON public.submissions_2016 USING btree (problem_id, length);


--
-- Name: submissions_2016_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2016_problem_id_point_idx ON public.submissions_2016 USING btree (problem_id, point);


--
-- Name: submissions_2016_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2016_result_epoch_second_idx ON public.submissions_2016 USING btree (result, epoch_second);


--
-- Name: submissions_2016_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2016_result_execution_time_idx ON public.submissions_2016 USING btree (result, execution_time);


--
-- Name: submissions_2016_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2016_result_length_idx ON public.submissions_2016 USING btree (result, length);


--
-- Name: submissions_2016_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2016_result_point_idx ON public.submissions_2016 USING btree (result, point);


--
-- Name: submissions_2016_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2016_user_id_epoch_second_idx ON public.submissions_2016 USING btree (user_id, epoch_second);


--
-- Name: submissions_2016_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2016_user_id_execution_time_idx ON public.submissions_2016 USING btree (user_id, execution_time);


--
-- Name: submissions_2016_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2016_user_id_length_idx ON public.submissions_2016 USING btree (user_id, length);


--
-- Name: submissions_2016_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2016_user_id_point_idx ON public.submissions_2016 USING btree (user_id, point);


--
-- Name: submissions_2016_user_id_result_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2016_user_id_result_idx ON public.submissions_2016 USING btree (user_id, result);


--
-- Name: submissions_2017_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2017_contest_id_epoch_second_idx ON public.submissions_2017 USING btree (contest_id, epoch_second);


--
-- Name: submissions_2017_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2017_contest_id_execution_time_idx ON public.submissions_2017 USING btree (contest_id, execution_time);


--
-- Name: submissions_2017_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2017_contest_id_length_idx ON public.submissions_2017 USING btree (contest_id, length);


--
-- Name: submissions_2017_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2017_contest_id_point_idx ON public.submissions_2017 USING btree (contest_id, point);


--
-- Name: submissions_2017_epoch_second_id_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_2017_epoch_second_id_idx ON public.submissions_2017 USING btree (epoch_second, id);


--
-- Name: submissions_2017_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2017_language_epoch_second_idx ON public.submissions_2017 USING btree (language, epoch_second);


--
-- Name: submissions_2017_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2017_language_execution_time_idx ON public.submissions_2017 USING btree (language, execution_time);


--
-- Name: submissions_2017_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2017_language_length_idx ON public.submissions_2017 USING btree (language, length);


--
-- Name: submissions_2017_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2017_language_point_idx ON public.submissions_2017 USING btree (language, point);


--
-- Name: submissions_2017_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2017_problem_id_epoch_second_idx ON public.submissions_2017 USING btree (problem_id, epoch_second);


--
-- Name: submissions_2017_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2017_problem_id_execution_time_idx ON public.submissions_2017 USING btree (problem_id, execution_time);


--
-- Name: submissions_2017_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2017_problem_id_length_idx ON public.submissions_2017 USING btree (problem_id, length);


--
-- Name: submissions_2017_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2017_problem_id_point_idx ON public.submissions_2017 USING btree (problem_id, point);


--
-- Name: submissions_2017_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2017_result_epoch_second_idx ON public.submissions_2017 USING btree (result, epoch_second);


--
-- Name: submissions_2017_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2017_result_execution_time_idx ON public.submissions_2017 USING btree (result, execution_time);


--
-- Name: submissions_2017_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2017_result_length_idx ON public.submissions_2017 USING btree (result, length);


--
-- Name: submissions_2017_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2017_result_point_idx ON public.submissions_2017 USING btree (result, point);


--
-- Name: submissions_2017_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2017_user_id_epoch_second_idx ON public.submissions_2017 USING btree (user_id, epoch_second);


--
-- Name: submissions_2017_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2017_user_id_execution_time_idx ON public.submissions_2017 USING btree (user_id, execution_time);


--
-- Name: submissions_2017_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2017_user_id_length_idx ON public.submissions_2017 USING btree (user_id, length);


--
-- Name: submissions_2017_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2017_user_id_point_idx ON public.submissions_2017 USING btree (user_id, point);


--
-- Name: submissions_2017_user_id_result_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2017_user_id_result_idx ON public.submissions_2017 USING btree (user_id, result);


--
-- Name: submissions_2018_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2018_contest_id_epoch_second_idx ON public.submissions_2018 USING btree (contest_id, epoch_second);


--
-- Name: submissions_2018_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2018_contest_id_execution_time_idx ON public.submissions_2018 USING btree (contest_id, execution_time);


--
-- Name: submissions_2018_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2018_contest_id_length_idx ON public.submissions_2018 USING btree (contest_id, length);


--
-- Name: submissions_2018_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2018_contest_id_point_idx ON public.submissions_2018 USING btree (contest_id, point);


--
-- Name: submissions_2018_epoch_second_id_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_2018_epoch_second_id_idx ON public.submissions_2018 USING btree (epoch_second, id);


--
-- Name: submissions_2018_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2018_language_epoch_second_idx ON public.submissions_2018 USING btree (language, epoch_second);


--
-- Name: submissions_2018_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2018_language_execution_time_idx ON public.submissions_2018 USING btree (language, execution_time);


--
-- Name: submissions_2018_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2018_language_length_idx ON public.submissions_2018 USING btree (language, length);


--
-- Name: submissions_2018_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2018_language_point_idx ON public.submissions_2018 USING btree (language, point);


--
-- Name: submissions_2018_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2018_problem_id_epoch_second_idx ON public.submissions_2018 USING btree (problem_id, epoch_second);


--
-- Name: submissions_2018_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2018_problem_id_execution_time_idx ON public.submissions_2018 USING btree (problem_id, execution_time);


--
-- Name: submissions_2018_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2018_problem_id_length_idx ON public.submissions_2018 USING btree (problem_id, length);


--
-- Name: submissions_2018_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2018_problem_id_point_idx ON public.submissions_2018 USING btree (problem_id, point);


--
-- Name: submissions_2018_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2018_result_epoch_second_idx ON public.submissions_2018 USING btree (result, epoch_second);


--
-- Name: submissions_2018_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2018_result_execution_time_idx ON public.submissions_2018 USING btree (result, execution_time);


--
-- Name: submissions_2018_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2018_result_length_idx ON public.submissions_2018 USING btree (result, length);


--
-- Name: submissions_2018_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2018_result_point_idx ON public.submissions_2018 USING btree (result, point);


--
-- Name: submissions_2018_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2018_user_id_epoch_second_idx ON public.submissions_2018 USING btree (user_id, epoch_second);


--
-- Name: submissions_2018_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2018_user_id_execution_time_idx ON public.submissions_2018 USING btree (user_id, execution_time);


--
-- Name: submissions_2018_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2018_user_id_length_idx ON public.submissions_2018 USING btree (user_id, length);


--
-- Name: submissions_2018_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2018_user_id_point_idx ON public.submissions_2018 USING btree (user_id, point);


--
-- Name: submissions_2018_user_id_result_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2018_user_id_result_idx ON public.submissions_2018 USING btree (user_id, result);


--
-- Name: submissions_2019_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2019_contest_id_epoch_second_idx ON public.submissions_2019 USING btree (contest_id, epoch_second);


--
-- Name: submissions_2019_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2019_contest_id_execution_time_idx ON public.submissions_2019 USING btree (contest_id, execution_time);


--
-- Name: submissions_2019_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2019_contest_id_length_idx ON public.submissions_2019 USING btree (contest_id, length);


--
-- Name: submissions_2019_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2019_contest_id_point_idx ON public.submissions_2019 USING btree (contest_id, point);


--
-- Name: submissions_2019_epoch_second_id_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_2019_epoch_second_id_idx ON public.submissions_2019 USING btree (epoch_second, id);


--
-- Name: submissions_2019_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2019_language_epoch_second_idx ON public.submissions_2019 USING btree (language, epoch_second);


--
-- Name: submissions_2019_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2019_language_execution_time_idx ON public.submissions_2019 USING btree (language, execution_time);


--
-- Name: submissions_2019_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2019_language_length_idx ON public.submissions_2019 USING btree (language, length);


--
-- Name: submissions_2019_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2019_language_point_idx ON public.submissions_2019 USING btree (language, point);


--
-- Name: submissions_2019_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2019_problem_id_epoch_second_idx ON public.submissions_2019 USING btree (problem_id, epoch_second);


--
-- Name: submissions_2019_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2019_problem_id_execution_time_idx ON public.submissions_2019 USING btree (problem_id, execution_time);


--
-- Name: submissions_2019_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2019_problem_id_length_idx ON public.submissions_2019 USING btree (problem_id, length);


--
-- Name: submissions_2019_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2019_problem_id_point_idx ON public.submissions_2019 USING btree (problem_id, point);


--
-- Name: submissions_2019_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2019_result_epoch_second_idx ON public.submissions_2019 USING btree (result, epoch_second);


--
-- Name: submissions_2019_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2019_result_execution_time_idx ON public.submissions_2019 USING btree (result, execution_time);


--
-- Name: submissions_2019_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2019_result_length_idx ON public.submissions_2019 USING btree (result, length);


--
-- Name: submissions_2019_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2019_result_point_idx ON public.submissions_2019 USING btree (result, point);


--
-- Name: submissions_2019_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2019_user_id_epoch_second_idx ON public.submissions_2019 USING btree (user_id, epoch_second);


--
-- Name: submissions_2019_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2019_user_id_execution_time_idx ON public.submissions_2019 USING btree (user_id, execution_time);


--
-- Name: submissions_2019_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2019_user_id_length_idx ON public.submissions_2019 USING btree (user_id, length);


--
-- Name: submissions_2019_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2019_user_id_point_idx ON public.submissions_2019 USING btree (user_id, point);


--
-- Name: submissions_2019_user_id_result_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2019_user_id_result_idx ON public.submissions_2019 USING btree (user_id, result);


--
-- Name: submissions_2020_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2020_contest_id_epoch_second_idx ON public.submissions_2020 USING btree (contest_id, epoch_second);


--
-- Name: submissions_2020_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2020_contest_id_execution_time_idx ON public.submissions_2020 USING btree (contest_id, execution_time);


--
-- Name: submissions_2020_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2020_contest_id_length_idx ON public.submissions_2020 USING btree (contest_id, length);


--
-- Name: submissions_2020_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2020_contest_id_point_idx ON public.submissions_2020 USING btree (contest_id, point);


--
-- Name: submissions_2020_epoch_second_id_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_2020_epoch_second_id_idx ON public.submissions_2020 USING btree (epoch_second, id);


--
-- Name: submissions_2020_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2020_language_epoch_second_idx ON public.submissions_2020 USING btree (language, epoch_second);


--
-- Name: submissions_2020_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2020_language_execution_time_idx ON public.submissions_2020 USING btree (language, execution_time);


--
-- Name: submissions_2020_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2020_language_length_idx ON public.submissions_2020 USING btree (language, length);


--
-- Name: submissions_2020_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2020_language_point_idx ON public.submissions_2020 USING btree (language, point);


--
-- Name: submissions_2020_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2020_problem_id_epoch_second_idx ON public.submissions_2020 USING btree (problem_id, epoch_second);


--
-- Name: submissions_2020_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2020_problem_id_execution_time_idx ON public.submissions_2020 USING btree (problem_id, execution_time);


--
-- Name: submissions_2020_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2020_problem_id_length_idx ON public.submissions_2020 USING btree (problem_id, length);


--
-- Name: submissions_2020_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2020_problem_id_point_idx ON public.submissions_2020 USING btree (problem_id, point);


--
-- Name: submissions_2020_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2020_result_epoch_second_idx ON public.submissions_2020 USING btree (result, epoch_second);


--
-- Name: submissions_2020_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2020_result_execution_time_idx ON public.submissions_2020 USING btree (result, execution_time);


--
-- Name: submissions_2020_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2020_result_length_idx ON public.submissions_2020 USING btree (result, length);


--
-- Name: submissions_2020_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2020_result_point_idx ON public.submissions_2020 USING btree (result, point);


--
-- Name: submissions_2020_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2020_user_id_epoch_second_idx ON public.submissions_2020 USING btree (user_id, epoch_second);


--
-- Name: submissions_2020_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2020_user_id_execution_time_idx ON public.submissions_2020 USING btree (user_id, execution_time);


--
-- Name: submissions_2020_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2020_user_id_length_idx ON public.submissions_2020 USING btree (user_id, length);


--
-- Name: submissions_2020_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2020_user_id_point_idx ON public.submissions_2020 USING btree (user_id, point);


--
-- Name: submissions_2020_user_id_result_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2020_user_id_result_idx ON public.submissions_2020 USING btree (user_id, result);


--
-- Name: submissions_2021_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2021_contest_id_epoch_second_idx ON public.submissions_2021 USING btree (contest_id, epoch_second);


--
-- Name: submissions_2021_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2021_contest_id_execution_time_idx ON public.submissions_2021 USING btree (contest_id, execution_time);


--
-- Name: submissions_2021_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2021_contest_id_length_idx ON public.submissions_2021 USING btree (contest_id, length);


--
-- Name: submissions_2021_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2021_contest_id_point_idx ON public.submissions_2021 USING btree (contest_id, point);


--
-- Name: submissions_2021_epoch_second_id_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_2021_epoch_second_id_idx ON public.submissions_2021 USING btree (epoch_second, id);


--
-- Name: submissions_2021_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2021_language_epoch_second_idx ON public.submissions_2021 USING btree (language, epoch_second);


--
-- Name: submissions_2021_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2021_language_execution_time_idx ON public.submissions_2021 USING btree (language, execution_time);


--
-- Name: submissions_2021_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2021_language_length_idx ON public.submissions_2021 USING btree (language, length);


--
-- Name: submissions_2021_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2021_language_point_idx ON public.submissions_2021 USING btree (language, point);


--
-- Name: submissions_2021_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2021_problem_id_epoch_second_idx ON public.submissions_2021 USING btree (problem_id, epoch_second);


--
-- Name: submissions_2021_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2021_problem_id_execution_time_idx ON public.submissions_2021 USING btree (problem_id, execution_time);


--
-- Name: submissions_2021_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2021_problem_id_length_idx ON public.submissions_2021 USING btree (problem_id, length);


--
-- Name: submissions_2021_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2021_problem_id_point_idx ON public.submissions_2021 USING btree (problem_id, point);


--
-- Name: submissions_2021_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2021_result_epoch_second_idx ON public.submissions_2021 USING btree (result, epoch_second);


--
-- Name: submissions_2021_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2021_result_execution_time_idx ON public.submissions_2021 USING btree (result, execution_time);


--
-- Name: submissions_2021_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2021_result_length_idx ON public.submissions_2021 USING btree (result, length);


--
-- Name: submissions_2021_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2021_result_point_idx ON public.submissions_2021 USING btree (result, point);


--
-- Name: submissions_2021_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2021_user_id_epoch_second_idx ON public.submissions_2021 USING btree (user_id, epoch_second);


--
-- Name: submissions_2021_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2021_user_id_execution_time_idx ON public.submissions_2021 USING btree (user_id, execution_time);


--
-- Name: submissions_2021_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2021_user_id_length_idx ON public.submissions_2021 USING btree (user_id, length);


--
-- Name: submissions_2021_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2021_user_id_point_idx ON public.submissions_2021 USING btree (user_id, point);


--
-- Name: submissions_2021_user_id_result_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2021_user_id_result_idx ON public.submissions_2021 USING btree (user_id, result);


--
-- Name: submissions_2022_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2022_contest_id_epoch_second_idx ON public.submissions_2022 USING btree (contest_id, epoch_second);


--
-- Name: submissions_2022_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2022_contest_id_execution_time_idx ON public.submissions_2022 USING btree (contest_id, execution_time);


--
-- Name: submissions_2022_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2022_contest_id_length_idx ON public.submissions_2022 USING btree (contest_id, length);


--
-- Name: submissions_2022_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2022_contest_id_point_idx ON public.submissions_2022 USING btree (contest_id, point);


--
-- Name: submissions_2022_epoch_second_id_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_2022_epoch_second_id_idx ON public.submissions_2022 USING btree (epoch_second, id);


--
-- Name: submissions_2022_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2022_language_epoch_second_idx ON public.submissions_2022 USING btree (language, epoch_second);


--
-- Name: submissions_2022_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2022_language_execution_time_idx ON public.submissions_2022 USING btree (language, execution_time);


--
-- Name: submissions_2022_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2022_language_length_idx ON public.submissions_2022 USING btree (language, length);


--
-- Name: submissions_2022_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2022_language_point_idx ON public.submissions_2022 USING btree (language, point);


--
-- Name: submissions_2022_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2022_problem_id_epoch_second_idx ON public.submissions_2022 USING btree (problem_id, epoch_second);


--
-- Name: submissions_2022_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2022_problem_id_execution_time_idx ON public.submissions_2022 USING btree (problem_id, execution_time);


--
-- Name: submissions_2022_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2022_problem_id_length_idx ON public.submissions_2022 USING btree (problem_id, length);


--
-- Name: submissions_2022_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2022_problem_id_point_idx ON public.submissions_2022 USING btree (problem_id, point);


--
-- Name: submissions_2022_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2022_result_epoch_second_idx ON public.submissions_2022 USING btree (result, epoch_second);


--
-- Name: submissions_2022_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2022_result_execution_time_idx ON public.submissions_2022 USING btree (result, execution_time);


--
-- Name: submissions_2022_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2022_result_length_idx ON public.submissions_2022 USING btree (result, length);


--
-- Name: submissions_2022_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2022_result_point_idx ON public.submissions_2022 USING btree (result, point);


--
-- Name: submissions_2022_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2022_user_id_epoch_second_idx ON public.submissions_2022 USING btree (user_id, epoch_second);


--
-- Name: submissions_2022_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2022_user_id_execution_time_idx ON public.submissions_2022 USING btree (user_id, execution_time);


--
-- Name: submissions_2022_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2022_user_id_length_idx ON public.submissions_2022 USING btree (user_id, length);


--
-- Name: submissions_2022_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2022_user_id_point_idx ON public.submissions_2022 USING btree (user_id, point);


--
-- Name: submissions_2022_user_id_result_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2022_user_id_result_idx ON public.submissions_2022 USING btree (user_id, result);


--
-- Name: submissions_2023_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2023_contest_id_epoch_second_idx ON public.submissions_2023 USING btree (contest_id, epoch_second);


--
-- Name: submissions_2023_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2023_contest_id_execution_time_idx ON public.submissions_2023 USING btree (contest_id, execution_time);


--
-- Name: submissions_2023_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2023_contest_id_length_idx ON public.submissions_2023 USING btree (contest_id, length);


--
-- Name: submissions_2023_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2023_contest_id_point_idx ON public.submissions_2023 USING btree (contest_id, point);


--
-- Name: submissions_2023_epoch_second_id_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_2023_epoch_second_id_idx ON public.submissions_2023 USING btree (epoch_second, id);


--
-- Name: submissions_2023_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2023_language_epoch_second_idx ON public.submissions_2023 USING btree (language, epoch_second);


--
-- Name: submissions_2023_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2023_language_execution_time_idx ON public.submissions_2023 USING btree (language, execution_time);


--
-- Name: submissions_2023_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2023_language_length_idx ON public.submissions_2023 USING btree (language, length);


--
-- Name: submissions_2023_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2023_language_point_idx ON public.submissions_2023 USING btree (language, point);


--
-- Name: submissions_2023_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2023_problem_id_epoch_second_idx ON public.submissions_2023 USING btree (problem_id, epoch_second);


--
-- Name: submissions_2023_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2023_problem_id_execution_time_idx ON public.submissions_2023 USING btree (problem_id, execution_time);


--
-- Name: submissions_2023_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2023_problem_id_length_idx ON public.submissions_2023 USING btree (problem_id, length);


--
-- Name: submissions_2023_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2023_problem_id_point_idx ON public.submissions_2023 USING btree (problem_id, point);


--
-- Name: submissions_2023_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2023_result_epoch_second_idx ON public.submissions_2023 USING btree (result, epoch_second);


--
-- Name: submissions_2023_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2023_result_execution_time_idx ON public.submissions_2023 USING btree (result, execution_time);


--
-- Name: submissions_2023_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2023_result_length_idx ON public.submissions_2023 USING btree (result, length);


--
-- Name: submissions_2023_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2023_result_point_idx ON public.submissions_2023 USING btree (result, point);


--
-- Name: submissions_2023_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2023_user_id_epoch_second_idx ON public.submissions_2023 USING btree (user_id, epoch_second);


--
-- Name: submissions_2023_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2023_user_id_execution_time_idx ON public.submissions_2023 USING btree (user_id, execution_time);


--
-- Name: submissions_2023_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2023_user_id_length_idx ON public.submissions_2023 USING btree (user_id, length);


--
-- Name: submissions_2023_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2023_user_id_point_idx ON public.submissions_2023 USING btree (user_id, point);


--
-- Name: submissions_2023_user_id_result_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2023_user_id_result_idx ON public.submissions_2023 USING btree (user_id, result);


--
-- Name: submissions_2024_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2024_contest_id_epoch_second_idx ON public.submissions_2024 USING btree (contest_id, epoch_second);


--
-- Name: submissions_2024_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2024_contest_id_execution_time_idx ON public.submissions_2024 USING btree (contest_id, execution_time);


--
-- Name: submissions_2024_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2024_contest_id_length_idx ON public.submissions_2024 USING btree (contest_id, length);


--
-- Name: submissions_2024_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2024_contest_id_point_idx ON public.submissions_2024 USING btree (contest_id, point);


--
-- Name: submissions_2024_epoch_second_id_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_2024_epoch_second_id_idx ON public.submissions_2024 USING btree (epoch_second, id);


--
-- Name: submissions_2024_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2024_language_epoch_second_idx ON public.submissions_2024 USING btree (language, epoch_second);


--
-- Name: submissions_2024_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2024_language_execution_time_idx ON public.submissions_2024 USING btree (language, execution_time);


--
-- Name: submissions_2024_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2024_language_length_idx ON public.submissions_2024 USING btree (language, length);


--
-- Name: submissions_2024_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2024_language_point_idx ON public.submissions_2024 USING btree (language, point);


--
-- Name: submissions_2024_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2024_problem_id_epoch_second_idx ON public.submissions_2024 USING btree (problem_id, epoch_second);


--
-- Name: submissions_2024_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2024_problem_id_execution_time_idx ON public.submissions_2024 USING btree (problem_id, execution_time);


--
-- Name: submissions_2024_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2024_problem_id_length_idx ON public.submissions_2024 USING btree (problem_id, length);


--
-- Name: submissions_2024_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2024_problem_id_point_idx ON public.submissions_2024 USING btree (problem_id, point);


--
-- Name: submissions_2024_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2024_result_epoch_second_idx ON public.submissions_2024 USING btree (result, epoch_second);


--
-- Name: submissions_2024_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2024_result_execution_time_idx ON public.submissions_2024 USING btree (result, execution_time);


--
-- Name: submissions_2024_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2024_result_length_idx ON public.submissions_2024 USING btree (result, length);


--
-- Name: submissions_2024_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2024_result_point_idx ON public.submissions_2024 USING btree (result, point);


--
-- Name: submissions_2024_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2024_user_id_epoch_second_idx ON public.submissions_2024 USING btree (user_id, epoch_second);


--
-- Name: submissions_2024_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2024_user_id_execution_time_idx ON public.submissions_2024 USING btree (user_id, execution_time);


--
-- Name: submissions_2024_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2024_user_id_length_idx ON public.submissions_2024 USING btree (user_id, length);


--
-- Name: submissions_2024_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2024_user_id_point_idx ON public.submissions_2024 USING btree (user_id, point);


--
-- Name: submissions_2024_user_id_result_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2024_user_id_result_idx ON public.submissions_2024 USING btree (user_id, result);


--
-- Name: submissions_2025_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2025_contest_id_epoch_second_idx ON public.submissions_2025 USING btree (contest_id, epoch_second);


--
-- Name: submissions_2025_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2025_contest_id_execution_time_idx ON public.submissions_2025 USING btree (contest_id, execution_time);


--
-- Name: submissions_2025_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2025_contest_id_length_idx ON public.submissions_2025 USING btree (contest_id, length);


--
-- Name: submissions_2025_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2025_contest_id_point_idx ON public.submissions_2025 USING btree (contest_id, point);


--
-- Name: submissions_2025_epoch_second_id_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_2025_epoch_second_id_idx ON public.submissions_2025 USING btree (epoch_second, id);


--
-- Name: submissions_2025_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2025_language_epoch_second_idx ON public.submissions_2025 USING btree (language, epoch_second);


--
-- Name: submissions_2025_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2025_language_execution_time_idx ON public.submissions_2025 USING btree (language, execution_time);


--
-- Name: submissions_2025_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2025_language_length_idx ON public.submissions_2025 USING btree (language, length);


--
-- Name: submissions_2025_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2025_language_point_idx ON public.submissions_2025 USING btree (language, point);


--
-- Name: submissions_2025_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2025_problem_id_epoch_second_idx ON public.submissions_2025 USING btree (problem_id, epoch_second);


--
-- Name: submissions_2025_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2025_problem_id_execution_time_idx ON public.submissions_2025 USING btree (problem_id, execution_time);


--
-- Name: submissions_2025_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2025_problem_id_length_idx ON public.submissions_2025 USING btree (problem_id, length);


--
-- Name: submissions_2025_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2025_problem_id_point_idx ON public.submissions_2025 USING btree (problem_id, point);


--
-- Name: submissions_2025_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2025_result_epoch_second_idx ON public.submissions_2025 USING btree (result, epoch_second);


--
-- Name: submissions_2025_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2025_result_execution_time_idx ON public.submissions_2025 USING btree (result, execution_time);


--
-- Name: submissions_2025_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2025_result_length_idx ON public.submissions_2025 USING btree (result, length);


--
-- Name: submissions_2025_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2025_result_point_idx ON public.submissions_2025 USING btree (result, point);


--
-- Name: submissions_2025_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2025_user_id_epoch_second_idx ON public.submissions_2025 USING btree (user_id, epoch_second);


--
-- Name: submissions_2025_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2025_user_id_execution_time_idx ON public.submissions_2025 USING btree (user_id, execution_time);


--
-- Name: submissions_2025_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2025_user_id_length_idx ON public.submissions_2025 USING btree (user_id, length);


--
-- Name: submissions_2025_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2025_user_id_point_idx ON public.submissions_2025 USING btree (user_id, point);


--
-- Name: submissions_2025_user_id_result_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2025_user_id_result_idx ON public.submissions_2025 USING btree (user_id, result);


--
-- Name: submissions_2026_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2026_contest_id_epoch_second_idx ON public.submissions_2026 USING btree (contest_id, epoch_second);


--
-- Name: submissions_2026_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2026_contest_id_execution_time_idx ON public.submissions_2026 USING btree (contest_id, execution_time);


--
-- Name: submissions_2026_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2026_contest_id_length_idx ON public.submissions_2026 USING btree (contest_id, length);


--
-- Name: submissions_2026_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2026_contest_id_point_idx ON public.submissions_2026 USING btree (contest_id, point);


--
-- Name: submissions_2026_epoch_second_id_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_2026_epoch_second_id_idx ON public.submissions_2026 USING btree (epoch_second, id);


--
-- Name: submissions_2026_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2026_language_epoch_second_idx ON public.submissions_2026 USING btree (language, epoch_second);


--
-- Name: submissions_2026_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2026_language_execution_time_idx ON public.submissions_2026 USING btree (language, execution_time);


--
-- Name: submissions_2026_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2026_language_length_idx ON public.submissions_2026 USING btree (language, length);


--
-- Name: submissions_2026_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2026_language_point_idx ON public.submissions_2026 USING btree (language, point);


--
-- Name: submissions_2026_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2026_problem_id_epoch_second_idx ON public.submissions_2026 USING btree (problem_id, epoch_second);


--
-- Name: submissions_2026_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2026_problem_id_execution_time_idx ON public.submissions_2026 USING btree (problem_id, execution_time);


--
-- Name: submissions_2026_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2026_problem_id_length_idx ON public.submissions_2026 USING btree (problem_id, length);


--
-- Name: submissions_2026_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2026_problem_id_point_idx ON public.submissions_2026 USING btree (problem_id, point);


--
-- Name: submissions_2026_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2026_result_epoch_second_idx ON public.submissions_2026 USING btree (result, epoch_second);


--
-- Name: submissions_2026_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2026_result_execution_time_idx ON public.submissions_2026 USING btree (result, execution_time);


--
-- Name: submissions_2026_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2026_result_length_idx ON public.submissions_2026 USING btree (result, length);


--
-- Name: submissions_2026_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2026_result_point_idx ON public.submissions_2026 USING btree (result, point);


--
-- Name: submissions_2026_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2026_user_id_epoch_second_idx ON public.submissions_2026 USING btree (user_id, epoch_second);


--
-- Name: submissions_2026_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2026_user_id_execution_time_idx ON public.submissions_2026 USING btree (user_id, execution_time);


--
-- Name: submissions_2026_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2026_user_id_length_idx ON public.submissions_2026 USING btree (user_id, length);


--
-- Name: submissions_2026_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2026_user_id_point_idx ON public.submissions_2026 USING btree (user_id, point);


--
-- Name: submissions_2026_user_id_result_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2026_user_id_result_idx ON public.submissions_2026 USING btree (user_id, result);


--
-- Name: submissions_2027_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2027_contest_id_epoch_second_idx ON public.submissions_2027 USING btree (contest_id, epoch_second);


--
-- Name: submissions_2027_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2027_contest_id_execution_time_idx ON public.submissions_2027 USING btree (contest_id, execution_time);


--
-- Name: submissions_2027_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2027_contest_id_length_idx ON public.submissions_2027 USING btree (contest_id, length);


--
-- Name: submissions_2027_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2027_contest_id_point_idx ON public.submissions_2027 USING btree (contest_id, point);


--
-- Name: submissions_2027_epoch_second_id_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_2027_epoch_second_id_idx ON public.submissions_2027 USING btree (epoch_second, id);


--
-- Name: submissions_2027_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2027_language_epoch_second_idx ON public.submissions_2027 USING btree (language, epoch_second);


--
-- Name: submissions_2027_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2027_language_execution_time_idx ON public.submissions_2027 USING btree (language, execution_time);


--
-- Name: submissions_2027_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2027_language_length_idx ON public.submissions_2027 USING btree (language, length);


--
-- Name: submissions_2027_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2027_language_point_idx ON public.submissions_2027 USING btree (language, point);


--
-- Name: submissions_2027_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2027_problem_id_epoch_second_idx ON public.submissions_2027 USING btree (problem_id, epoch_second);


--
-- Name: submissions_2027_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2027_problem_id_execution_time_idx ON public.submissions_2027 USING btree (problem_id, execution_time);


--
-- Name: submissions_2027_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2027_problem_id_length_idx ON public.submissions_2027 USING btree (problem_id, length);


--
-- Name: submissions_2027_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2027_problem_id_point_idx ON public.submissions_2027 USING btree (problem_id, point);


--
-- Name: submissions_2027_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2027_result_epoch_second_idx ON public.submissions_2027 USING btree (result, epoch_second);


--
-- Name: submissions_2027_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2027_result_execution_time_idx ON public.submissions_2027 USING btree (result, execution_time);


--
-- Name: submissions_2027_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2027_result_length_idx ON public.submissions_2027 USING btree (result, length);


--
-- Name: submissions_2027_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2027_result_point_idx ON public.submissions_2027 USING btree (result, point);


--
-- Name: submissions_2027_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2027_user_id_epoch_second_idx ON public.submissions_2027 USING btree (user_id, epoch_second);


--
-- Name: submissions_2027_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2027_user_id_execution_time_idx ON public.submissions_2027 USING btree (user_id, execution_time);


--
-- Name: submissions_2027_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2027_user_id_length_idx ON public.submissions_2027 USING btree (user_id, length);


--
-- Name: submissions_2027_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2027_user_id_point_idx ON public.submissions_2027 USING btree (user_id, point);


--
-- Name: submissions_2027_user_id_result_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2027_user_id_result_idx ON public.submissions_2027 USING btree (user_id, result);


--
-- Name: submissions_2028_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2028_contest_id_epoch_second_idx ON public.submissions_2028 USING btree (contest_id, epoch_second);


--
-- Name: submissions_2028_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2028_contest_id_execution_time_idx ON public.submissions_2028 USING btree (contest_id, execution_time);


--
-- Name: submissions_2028_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2028_contest_id_length_idx ON public.submissions_2028 USING btree (contest_id, length);


--
-- Name: submissions_2028_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2028_contest_id_point_idx ON public.submissions_2028 USING btree (contest_id, point);


--
-- Name: submissions_2028_epoch_second_id_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_2028_epoch_second_id_idx ON public.submissions_2028 USING btree (epoch_second, id);


--
-- Name: submissions_2028_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2028_language_epoch_second_idx ON public.submissions_2028 USING btree (language, epoch_second);


--
-- Name: submissions_2028_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2028_language_execution_time_idx ON public.submissions_2028 USING btree (language, execution_time);


--
-- Name: submissions_2028_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2028_language_length_idx ON public.submissions_2028 USING btree (language, length);


--
-- Name: submissions_2028_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2028_language_point_idx ON public.submissions_2028 USING btree (language, point);


--
-- Name: submissions_2028_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2028_problem_id_epoch_second_idx ON public.submissions_2028 USING btree (problem_id, epoch_second);


--
-- Name: submissions_2028_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2028_problem_id_execution_time_idx ON public.submissions_2028 USING btree (problem_id, execution_time);


--
-- Name: submissions_2028_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2028_problem_id_length_idx ON public.submissions_2028 USING btree (problem_id, length);


--
-- Name: submissions_2028_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2028_problem_id_point_idx ON public.submissions_2028 USING btree (problem_id, point);


--
-- Name: submissions_2028_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2028_result_epoch_second_idx ON public.submissions_2028 USING btree (result, epoch_second);


--
-- Name: submissions_2028_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2028_result_execution_time_idx ON public.submissions_2028 USING btree (result, execution_time);


--
-- Name: submissions_2028_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2028_result_length_idx ON public.submissions_2028 USING btree (result, length);


--
-- Name: submissions_2028_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2028_result_point_idx ON public.submissions_2028 USING btree (result, point);


--
-- Name: submissions_2028_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2028_user_id_epoch_second_idx ON public.submissions_2028 USING btree (user_id, epoch_second);


--
-- Name: submissions_2028_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2028_user_id_execution_time_idx ON public.submissions_2028 USING btree (user_id, execution_time);


--
-- Name: submissions_2028_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2028_user_id_length_idx ON public.submissions_2028 USING btree (user_id, length);


--
-- Name: submissions_2028_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2028_user_id_point_idx ON public.submissions_2028 USING btree (user_id, point);


--
-- Name: submissions_2028_user_id_result_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2028_user_id_result_idx ON public.submissions_2028 USING btree (user_id, result);


--
-- Name: submissions_2029_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2029_contest_id_epoch_second_idx ON public.submissions_2029 USING btree (contest_id, epoch_second);


--
-- Name: submissions_2029_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2029_contest_id_execution_time_idx ON public.submissions_2029 USING btree (contest_id, execution_time);


--
-- Name: submissions_2029_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2029_contest_id_length_idx ON public.submissions_2029 USING btree (contest_id, length);


--
-- Name: submissions_2029_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2029_contest_id_point_idx ON public.submissions_2029 USING btree (contest_id, point);


--
-- Name: submissions_2029_epoch_second_id_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_2029_epoch_second_id_idx ON public.submissions_2029 USING btree (epoch_second, id);


--
-- Name: submissions_2029_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2029_language_epoch_second_idx ON public.submissions_2029 USING btree (language, epoch_second);


--
-- Name: submissions_2029_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2029_language_execution_time_idx ON public.submissions_2029 USING btree (language, execution_time);


--
-- Name: submissions_2029_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2029_language_length_idx ON public.submissions_2029 USING btree (language, length);


--
-- Name: submissions_2029_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2029_language_point_idx ON public.submissions_2029 USING btree (language, point);


--
-- Name: submissions_2029_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2029_problem_id_epoch_second_idx ON public.submissions_2029 USING btree (problem_id, epoch_second);


--
-- Name: submissions_2029_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2029_problem_id_execution_time_idx ON public.submissions_2029 USING btree (problem_id, execution_time);


--
-- Name: submissions_2029_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2029_problem_id_length_idx ON public.submissions_2029 USING btree (problem_id, length);


--
-- Name: submissions_2029_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2029_problem_id_point_idx ON public.submissions_2029 USING btree (problem_id, point);


--
-- Name: submissions_2029_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2029_result_epoch_second_idx ON public.submissions_2029 USING btree (result, epoch_second);


--
-- Name: submissions_2029_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2029_result_execution_time_idx ON public.submissions_2029 USING btree (result, execution_time);


--
-- Name: submissions_2029_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2029_result_length_idx ON public.submissions_2029 USING btree (result, length);


--
-- Name: submissions_2029_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2029_result_point_idx ON public.submissions_2029 USING btree (result, point);


--
-- Name: submissions_2029_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2029_user_id_epoch_second_idx ON public.submissions_2029 USING btree (user_id, epoch_second);


--
-- Name: submissions_2029_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2029_user_id_execution_time_idx ON public.submissions_2029 USING btree (user_id, execution_time);


--
-- Name: submissions_2029_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2029_user_id_length_idx ON public.submissions_2029 USING btree (user_id, length);


--
-- Name: submissions_2029_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2029_user_id_point_idx ON public.submissions_2029 USING btree (user_id, point);


--
-- Name: submissions_2029_user_id_result_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_2029_user_id_result_idx ON public.submissions_2029 USING btree (user_id, result);


--
-- Name: submissions_2010_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_2010_contest_id_epoch_second_idx;


--
-- Name: submissions_2010_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_2010_contest_id_execution_time_idx;


--
-- Name: submissions_2010_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_2010_contest_id_length_idx;


--
-- Name: submissions_2010_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_2010_contest_id_point_idx;


--
-- Name: submissions_2010_epoch_second_id_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_unique ATTACH PARTITION public.submissions_2010_epoch_second_id_idx;


--
-- Name: submissions_2010_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_2010_language_epoch_second_idx;


--
-- Name: submissions_2010_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_2010_language_execution_time_idx;


--
-- Name: submissions_2010_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_2010_language_length_idx;


--
-- Name: submissions_2010_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_2010_language_point_idx;


--
-- Name: submissions_2010_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_2010_problem_id_epoch_second_idx;


--
-- Name: submissions_2010_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_2010_problem_id_execution_time_idx;


--
-- Name: submissions_2010_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_2010_problem_id_length_idx;


--
-- Name: submissions_2010_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_2010_problem_id_point_idx;


--
-- Name: submissions_2010_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_2010_result_epoch_second_idx;


--
-- Name: submissions_2010_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_2010_result_execution_time_idx;


--
-- Name: submissions_2010_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_2010_result_length_idx;


--
-- Name: submissions_2010_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_2010_result_point_idx;


--
-- Name: submissions_2010_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_2010_user_id_epoch_second_idx;


--
-- Name: submissions_2010_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_2010_user_id_execution_time_idx;


--
-- Name: submissions_2010_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_2010_user_id_length_idx;


--
-- Name: submissions_2010_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_2010_user_id_point_idx;


--
-- Name: submissions_2010_user_id_result_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_result_index ATTACH PARTITION public.submissions_2010_user_id_result_idx;


--
-- Name: submissions_2011_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_2011_contest_id_epoch_second_idx;


--
-- Name: submissions_2011_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_2011_contest_id_execution_time_idx;


--
-- Name: submissions_2011_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_2011_contest_id_length_idx;


--
-- Name: submissions_2011_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_2011_contest_id_point_idx;


--
-- Name: submissions_2011_epoch_second_id_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_unique ATTACH PARTITION public.submissions_2011_epoch_second_id_idx;


--
-- Name: submissions_2011_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_2011_language_epoch_second_idx;


--
-- Name: submissions_2011_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_2011_language_execution_time_idx;


--
-- Name: submissions_2011_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_2011_language_length_idx;


--
-- Name: submissions_2011_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_2011_language_point_idx;


--
-- Name: submissions_2011_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_2011_problem_id_epoch_second_idx;


--
-- Name: submissions_2011_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_2011_problem_id_execution_time_idx;


--
-- Name: submissions_2011_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_2011_problem_id_length_idx;


--
-- Name: submissions_2011_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_2011_problem_id_point_idx;


--
-- Name: submissions_2011_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_2011_result_epoch_second_idx;


--
-- Name: submissions_2011_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_2011_result_execution_time_idx;


--
-- Name: submissions_2011_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_2011_result_length_idx;


--
-- Name: submissions_2011_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_2011_result_point_idx;


--
-- Name: submissions_2011_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_2011_user_id_epoch_second_idx;


--
-- Name: submissions_2011_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_2011_user_id_execution_time_idx;


--
-- Name: submissions_2011_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_2011_user_id_length_idx;


--
-- Name: submissions_2011_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_2011_user_id_point_idx;


--
-- Name: submissions_2011_user_id_result_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_result_index ATTACH PARTITION public.submissions_2011_user_id_result_idx;


--
-- Name: submissions_2012_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_2012_contest_id_epoch_second_idx;


--
-- Name: submissions_2012_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_2012_contest_id_execution_time_idx;


--
-- Name: submissions_2012_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_2012_contest_id_length_idx;


--
-- Name: submissions_2012_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_2012_contest_id_point_idx;


--
-- Name: submissions_2012_epoch_second_id_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_unique ATTACH PARTITION public.submissions_2012_epoch_second_id_idx;


--
-- Name: submissions_2012_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_2012_language_epoch_second_idx;


--
-- Name: submissions_2012_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_2012_language_execution_time_idx;


--
-- Name: submissions_2012_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_2012_language_length_idx;


--
-- Name: submissions_2012_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_2012_language_point_idx;


--
-- Name: submissions_2012_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_2012_problem_id_epoch_second_idx;


--
-- Name: submissions_2012_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_2012_problem_id_execution_time_idx;


--
-- Name: submissions_2012_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_2012_problem_id_length_idx;


--
-- Name: submissions_2012_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_2012_problem_id_point_idx;


--
-- Name: submissions_2012_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_2012_result_epoch_second_idx;


--
-- Name: submissions_2012_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_2012_result_execution_time_idx;


--
-- Name: submissions_2012_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_2012_result_length_idx;


--
-- Name: submissions_2012_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_2012_result_point_idx;


--
-- Name: submissions_2012_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_2012_user_id_epoch_second_idx;


--
-- Name: submissions_2012_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_2012_user_id_execution_time_idx;


--
-- Name: submissions_2012_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_2012_user_id_length_idx;


--
-- Name: submissions_2012_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_2012_user_id_point_idx;


--
-- Name: submissions_2012_user_id_result_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_result_index ATTACH PARTITION public.submissions_2012_user_id_result_idx;


--
-- Name: submissions_2013_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_2013_contest_id_epoch_second_idx;


--
-- Name: submissions_2013_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_2013_contest_id_execution_time_idx;


--
-- Name: submissions_2013_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_2013_contest_id_length_idx;


--
-- Name: submissions_2013_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_2013_contest_id_point_idx;


--
-- Name: submissions_2013_epoch_second_id_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_unique ATTACH PARTITION public.submissions_2013_epoch_second_id_idx;


--
-- Name: submissions_2013_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_2013_language_epoch_second_idx;


--
-- Name: submissions_2013_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_2013_language_execution_time_idx;


--
-- Name: submissions_2013_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_2013_language_length_idx;


--
-- Name: submissions_2013_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_2013_language_point_idx;


--
-- Name: submissions_2013_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_2013_problem_id_epoch_second_idx;


--
-- Name: submissions_2013_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_2013_problem_id_execution_time_idx;


--
-- Name: submissions_2013_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_2013_problem_id_length_idx;


--
-- Name: submissions_2013_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_2013_problem_id_point_idx;


--
-- Name: submissions_2013_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_2013_result_epoch_second_idx;


--
-- Name: submissions_2013_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_2013_result_execution_time_idx;


--
-- Name: submissions_2013_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_2013_result_length_idx;


--
-- Name: submissions_2013_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_2013_result_point_idx;


--
-- Name: submissions_2013_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_2013_user_id_epoch_second_idx;


--
-- Name: submissions_2013_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_2013_user_id_execution_time_idx;


--
-- Name: submissions_2013_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_2013_user_id_length_idx;


--
-- Name: submissions_2013_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_2013_user_id_point_idx;


--
-- Name: submissions_2013_user_id_result_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_result_index ATTACH PARTITION public.submissions_2013_user_id_result_idx;


--
-- Name: submissions_2014_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_2014_contest_id_epoch_second_idx;


--
-- Name: submissions_2014_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_2014_contest_id_execution_time_idx;


--
-- Name: submissions_2014_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_2014_contest_id_length_idx;


--
-- Name: submissions_2014_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_2014_contest_id_point_idx;


--
-- Name: submissions_2014_epoch_second_id_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_unique ATTACH PARTITION public.submissions_2014_epoch_second_id_idx;


--
-- Name: submissions_2014_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_2014_language_epoch_second_idx;


--
-- Name: submissions_2014_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_2014_language_execution_time_idx;


--
-- Name: submissions_2014_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_2014_language_length_idx;


--
-- Name: submissions_2014_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_2014_language_point_idx;


--
-- Name: submissions_2014_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_2014_problem_id_epoch_second_idx;


--
-- Name: submissions_2014_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_2014_problem_id_execution_time_idx;


--
-- Name: submissions_2014_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_2014_problem_id_length_idx;


--
-- Name: submissions_2014_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_2014_problem_id_point_idx;


--
-- Name: submissions_2014_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_2014_result_epoch_second_idx;


--
-- Name: submissions_2014_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_2014_result_execution_time_idx;


--
-- Name: submissions_2014_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_2014_result_length_idx;


--
-- Name: submissions_2014_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_2014_result_point_idx;


--
-- Name: submissions_2014_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_2014_user_id_epoch_second_idx;


--
-- Name: submissions_2014_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_2014_user_id_execution_time_idx;


--
-- Name: submissions_2014_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_2014_user_id_length_idx;


--
-- Name: submissions_2014_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_2014_user_id_point_idx;


--
-- Name: submissions_2014_user_id_result_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_result_index ATTACH PARTITION public.submissions_2014_user_id_result_idx;


--
-- Name: submissions_2015_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_2015_contest_id_epoch_second_idx;


--
-- Name: submissions_2015_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_2015_contest_id_execution_time_idx;


--
-- Name: submissions_2015_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_2015_contest_id_length_idx;


--
-- Name: submissions_2015_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_2015_contest_id_point_idx;


--
-- Name: submissions_2015_epoch_second_id_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_unique ATTACH PARTITION public.submissions_2015_epoch_second_id_idx;


--
-- Name: submissions_2015_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_2015_language_epoch_second_idx;


--
-- Name: submissions_2015_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_2015_language_execution_time_idx;


--
-- Name: submissions_2015_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_2015_language_length_idx;


--
-- Name: submissions_2015_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_2015_language_point_idx;


--
-- Name: submissions_2015_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_2015_problem_id_epoch_second_idx;


--
-- Name: submissions_2015_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_2015_problem_id_execution_time_idx;


--
-- Name: submissions_2015_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_2015_problem_id_length_idx;


--
-- Name: submissions_2015_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_2015_problem_id_point_idx;


--
-- Name: submissions_2015_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_2015_result_epoch_second_idx;


--
-- Name: submissions_2015_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_2015_result_execution_time_idx;


--
-- Name: submissions_2015_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_2015_result_length_idx;


--
-- Name: submissions_2015_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_2015_result_point_idx;


--
-- Name: submissions_2015_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_2015_user_id_epoch_second_idx;


--
-- Name: submissions_2015_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_2015_user_id_execution_time_idx;


--
-- Name: submissions_2015_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_2015_user_id_length_idx;


--
-- Name: submissions_2015_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_2015_user_id_point_idx;


--
-- Name: submissions_2015_user_id_result_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_result_index ATTACH PARTITION public.submissions_2015_user_id_result_idx;


--
-- Name: submissions_2016_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_2016_contest_id_epoch_second_idx;


--
-- Name: submissions_2016_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_2016_contest_id_execution_time_idx;


--
-- Name: submissions_2016_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_2016_contest_id_length_idx;


--
-- Name: submissions_2016_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_2016_contest_id_point_idx;


--
-- Name: submissions_2016_epoch_second_id_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_unique ATTACH PARTITION public.submissions_2016_epoch_second_id_idx;


--
-- Name: submissions_2016_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_2016_language_epoch_second_idx;


--
-- Name: submissions_2016_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_2016_language_execution_time_idx;


--
-- Name: submissions_2016_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_2016_language_length_idx;


--
-- Name: submissions_2016_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_2016_language_point_idx;


--
-- Name: submissions_2016_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_2016_problem_id_epoch_second_idx;


--
-- Name: submissions_2016_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_2016_problem_id_execution_time_idx;


--
-- Name: submissions_2016_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_2016_problem_id_length_idx;


--
-- Name: submissions_2016_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_2016_problem_id_point_idx;


--
-- Name: submissions_2016_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_2016_result_epoch_second_idx;


--
-- Name: submissions_2016_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_2016_result_execution_time_idx;


--
-- Name: submissions_2016_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_2016_result_length_idx;


--
-- Name: submissions_2016_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_2016_result_point_idx;


--
-- Name: submissions_2016_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_2016_user_id_epoch_second_idx;


--
-- Name: submissions_2016_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_2016_user_id_execution_time_idx;


--
-- Name: submissions_2016_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_2016_user_id_length_idx;


--
-- Name: submissions_2016_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_2016_user_id_point_idx;


--
-- Name: submissions_2016_user_id_result_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_result_index ATTACH PARTITION public.submissions_2016_user_id_result_idx;


--
-- Name: submissions_2017_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_2017_contest_id_epoch_second_idx;


--
-- Name: submissions_2017_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_2017_contest_id_execution_time_idx;


--
-- Name: submissions_2017_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_2017_contest_id_length_idx;


--
-- Name: submissions_2017_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_2017_contest_id_point_idx;


--
-- Name: submissions_2017_epoch_second_id_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_unique ATTACH PARTITION public.submissions_2017_epoch_second_id_idx;


--
-- Name: submissions_2017_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_2017_language_epoch_second_idx;


--
-- Name: submissions_2017_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_2017_language_execution_time_idx;


--
-- Name: submissions_2017_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_2017_language_length_idx;


--
-- Name: submissions_2017_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_2017_language_point_idx;


--
-- Name: submissions_2017_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_2017_problem_id_epoch_second_idx;


--
-- Name: submissions_2017_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_2017_problem_id_execution_time_idx;


--
-- Name: submissions_2017_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_2017_problem_id_length_idx;


--
-- Name: submissions_2017_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_2017_problem_id_point_idx;


--
-- Name: submissions_2017_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_2017_result_epoch_second_idx;


--
-- Name: submissions_2017_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_2017_result_execution_time_idx;


--
-- Name: submissions_2017_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_2017_result_length_idx;


--
-- Name: submissions_2017_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_2017_result_point_idx;


--
-- Name: submissions_2017_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_2017_user_id_epoch_second_idx;


--
-- Name: submissions_2017_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_2017_user_id_execution_time_idx;


--
-- Name: submissions_2017_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_2017_user_id_length_idx;


--
-- Name: submissions_2017_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_2017_user_id_point_idx;


--
-- Name: submissions_2017_user_id_result_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_result_index ATTACH PARTITION public.submissions_2017_user_id_result_idx;


--
-- Name: submissions_2018_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_2018_contest_id_epoch_second_idx;


--
-- Name: submissions_2018_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_2018_contest_id_execution_time_idx;


--
-- Name: submissions_2018_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_2018_contest_id_length_idx;


--
-- Name: submissions_2018_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_2018_contest_id_point_idx;


--
-- Name: submissions_2018_epoch_second_id_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_unique ATTACH PARTITION public.submissions_2018_epoch_second_id_idx;


--
-- Name: submissions_2018_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_2018_language_epoch_second_idx;


--
-- Name: submissions_2018_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_2018_language_execution_time_idx;


--
-- Name: submissions_2018_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_2018_language_length_idx;


--
-- Name: submissions_2018_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_2018_language_point_idx;


--
-- Name: submissions_2018_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_2018_problem_id_epoch_second_idx;


--
-- Name: submissions_2018_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_2018_problem_id_execution_time_idx;


--
-- Name: submissions_2018_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_2018_problem_id_length_idx;


--
-- Name: submissions_2018_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_2018_problem_id_point_idx;


--
-- Name: submissions_2018_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_2018_result_epoch_second_idx;


--
-- Name: submissions_2018_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_2018_result_execution_time_idx;


--
-- Name: submissions_2018_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_2018_result_length_idx;


--
-- Name: submissions_2018_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_2018_result_point_idx;


--
-- Name: submissions_2018_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_2018_user_id_epoch_second_idx;


--
-- Name: submissions_2018_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_2018_user_id_execution_time_idx;


--
-- Name: submissions_2018_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_2018_user_id_length_idx;


--
-- Name: submissions_2018_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_2018_user_id_point_idx;


--
-- Name: submissions_2018_user_id_result_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_result_index ATTACH PARTITION public.submissions_2018_user_id_result_idx;


--
-- Name: submissions_2019_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_2019_contest_id_epoch_second_idx;


--
-- Name: submissions_2019_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_2019_contest_id_execution_time_idx;


--
-- Name: submissions_2019_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_2019_contest_id_length_idx;


--
-- Name: submissions_2019_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_2019_contest_id_point_idx;


--
-- Name: submissions_2019_epoch_second_id_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_unique ATTACH PARTITION public.submissions_2019_epoch_second_id_idx;


--
-- Name: submissions_2019_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_2019_language_epoch_second_idx;


--
-- Name: submissions_2019_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_2019_language_execution_time_idx;


--
-- Name: submissions_2019_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_2019_language_length_idx;


--
-- Name: submissions_2019_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_2019_language_point_idx;


--
-- Name: submissions_2019_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_2019_problem_id_epoch_second_idx;


--
-- Name: submissions_2019_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_2019_problem_id_execution_time_idx;


--
-- Name: submissions_2019_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_2019_problem_id_length_idx;


--
-- Name: submissions_2019_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_2019_problem_id_point_idx;


--
-- Name: submissions_2019_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_2019_result_epoch_second_idx;


--
-- Name: submissions_2019_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_2019_result_execution_time_idx;


--
-- Name: submissions_2019_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_2019_result_length_idx;


--
-- Name: submissions_2019_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_2019_result_point_idx;


--
-- Name: submissions_2019_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_2019_user_id_epoch_second_idx;


--
-- Name: submissions_2019_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_2019_user_id_execution_time_idx;


--
-- Name: submissions_2019_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_2019_user_id_length_idx;


--
-- Name: submissions_2019_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_2019_user_id_point_idx;


--
-- Name: submissions_2019_user_id_result_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_result_index ATTACH PARTITION public.submissions_2019_user_id_result_idx;


--
-- Name: submissions_2020_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_2020_contest_id_epoch_second_idx;


--
-- Name: submissions_2020_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_2020_contest_id_execution_time_idx;


--
-- Name: submissions_2020_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_2020_contest_id_length_idx;


--
-- Name: submissions_2020_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_2020_contest_id_point_idx;


--
-- Name: submissions_2020_epoch_second_id_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_unique ATTACH PARTITION public.submissions_2020_epoch_second_id_idx;


--
-- Name: submissions_2020_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_2020_language_epoch_second_idx;


--
-- Name: submissions_2020_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_2020_language_execution_time_idx;


--
-- Name: submissions_2020_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_2020_language_length_idx;


--
-- Name: submissions_2020_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_2020_language_point_idx;


--
-- Name: submissions_2020_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_2020_problem_id_epoch_second_idx;


--
-- Name: submissions_2020_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_2020_problem_id_execution_time_idx;


--
-- Name: submissions_2020_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_2020_problem_id_length_idx;


--
-- Name: submissions_2020_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_2020_problem_id_point_idx;


--
-- Name: submissions_2020_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_2020_result_epoch_second_idx;


--
-- Name: submissions_2020_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_2020_result_execution_time_idx;


--
-- Name: submissions_2020_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_2020_result_length_idx;


--
-- Name: submissions_2020_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_2020_result_point_idx;


--
-- Name: submissions_2020_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_2020_user_id_epoch_second_idx;


--
-- Name: submissions_2020_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_2020_user_id_execution_time_idx;


--
-- Name: submissions_2020_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_2020_user_id_length_idx;


--
-- Name: submissions_2020_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_2020_user_id_point_idx;


--
-- Name: submissions_2020_user_id_result_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_result_index ATTACH PARTITION public.submissions_2020_user_id_result_idx;


--
-- Name: submissions_2021_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_2021_contest_id_epoch_second_idx;


--
-- Name: submissions_2021_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_2021_contest_id_execution_time_idx;


--
-- Name: submissions_2021_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_2021_contest_id_length_idx;


--
-- Name: submissions_2021_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_2021_contest_id_point_idx;


--
-- Name: submissions_2021_epoch_second_id_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_unique ATTACH PARTITION public.submissions_2021_epoch_second_id_idx;


--
-- Name: submissions_2021_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_2021_language_epoch_second_idx;


--
-- Name: submissions_2021_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_2021_language_execution_time_idx;


--
-- Name: submissions_2021_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_2021_language_length_idx;


--
-- Name: submissions_2021_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_2021_language_point_idx;


--
-- Name: submissions_2021_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_2021_problem_id_epoch_second_idx;


--
-- Name: submissions_2021_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_2021_problem_id_execution_time_idx;


--
-- Name: submissions_2021_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_2021_problem_id_length_idx;


--
-- Name: submissions_2021_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_2021_problem_id_point_idx;


--
-- Name: submissions_2021_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_2021_result_epoch_second_idx;


--
-- Name: submissions_2021_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_2021_result_execution_time_idx;


--
-- Name: submissions_2021_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_2021_result_length_idx;


--
-- Name: submissions_2021_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_2021_result_point_idx;


--
-- Name: submissions_2021_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_2021_user_id_epoch_second_idx;


--
-- Name: submissions_2021_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_2021_user_id_execution_time_idx;


--
-- Name: submissions_2021_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_2021_user_id_length_idx;


--
-- Name: submissions_2021_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_2021_user_id_point_idx;


--
-- Name: submissions_2021_user_id_result_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_result_index ATTACH PARTITION public.submissions_2021_user_id_result_idx;


--
-- Name: submissions_2022_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_2022_contest_id_epoch_second_idx;


--
-- Name: submissions_2022_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_2022_contest_id_execution_time_idx;


--
-- Name: submissions_2022_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_2022_contest_id_length_idx;


--
-- Name: submissions_2022_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_2022_contest_id_point_idx;


--
-- Name: submissions_2022_epoch_second_id_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_unique ATTACH PARTITION public.submissions_2022_epoch_second_id_idx;


--
-- Name: submissions_2022_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_2022_language_epoch_second_idx;


--
-- Name: submissions_2022_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_2022_language_execution_time_idx;


--
-- Name: submissions_2022_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_2022_language_length_idx;


--
-- Name: submissions_2022_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_2022_language_point_idx;


--
-- Name: submissions_2022_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_2022_problem_id_epoch_second_idx;


--
-- Name: submissions_2022_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_2022_problem_id_execution_time_idx;


--
-- Name: submissions_2022_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_2022_problem_id_length_idx;


--
-- Name: submissions_2022_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_2022_problem_id_point_idx;


--
-- Name: submissions_2022_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_2022_result_epoch_second_idx;


--
-- Name: submissions_2022_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_2022_result_execution_time_idx;


--
-- Name: submissions_2022_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_2022_result_length_idx;


--
-- Name: submissions_2022_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_2022_result_point_idx;


--
-- Name: submissions_2022_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_2022_user_id_epoch_second_idx;


--
-- Name: submissions_2022_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_2022_user_id_execution_time_idx;


--
-- Name: submissions_2022_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_2022_user_id_length_idx;


--
-- Name: submissions_2022_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_2022_user_id_point_idx;


--
-- Name: submissions_2022_user_id_result_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_result_index ATTACH PARTITION public.submissions_2022_user_id_result_idx;


--
-- Name: submissions_2023_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_2023_contest_id_epoch_second_idx;


--
-- Name: submissions_2023_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_2023_contest_id_execution_time_idx;


--
-- Name: submissions_2023_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_2023_contest_id_length_idx;


--
-- Name: submissions_2023_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_2023_contest_id_point_idx;


--
-- Name: submissions_2023_epoch_second_id_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_unique ATTACH PARTITION public.submissions_2023_epoch_second_id_idx;


--
-- Name: submissions_2023_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_2023_language_epoch_second_idx;


--
-- Name: submissions_2023_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_2023_language_execution_time_idx;


--
-- Name: submissions_2023_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_2023_language_length_idx;


--
-- Name: submissions_2023_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_2023_language_point_idx;


--
-- Name: submissions_2023_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_2023_problem_id_epoch_second_idx;


--
-- Name: submissions_2023_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_2023_problem_id_execution_time_idx;


--
-- Name: submissions_2023_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_2023_problem_id_length_idx;


--
-- Name: submissions_2023_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_2023_problem_id_point_idx;


--
-- Name: submissions_2023_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_2023_result_epoch_second_idx;


--
-- Name: submissions_2023_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_2023_result_execution_time_idx;


--
-- Name: submissions_2023_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_2023_result_length_idx;


--
-- Name: submissions_2023_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_2023_result_point_idx;


--
-- Name: submissions_2023_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_2023_user_id_epoch_second_idx;


--
-- Name: submissions_2023_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_2023_user_id_execution_time_idx;


--
-- Name: submissions_2023_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_2023_user_id_length_idx;


--
-- Name: submissions_2023_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_2023_user_id_point_idx;


--
-- Name: submissions_2023_user_id_result_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_result_index ATTACH PARTITION public.submissions_2023_user_id_result_idx;


--
-- Name: submissions_2024_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_2024_contest_id_epoch_second_idx;


--
-- Name: submissions_2024_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_2024_contest_id_execution_time_idx;


--
-- Name: submissions_2024_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_2024_contest_id_length_idx;


--
-- Name: submissions_2024_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_2024_contest_id_point_idx;


--
-- Name: submissions_2024_epoch_second_id_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_unique ATTACH PARTITION public.submissions_2024_epoch_second_id_idx;


--
-- Name: submissions_2024_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_2024_language_epoch_second_idx;


--
-- Name: submissions_2024_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_2024_language_execution_time_idx;


--
-- Name: submissions_2024_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_2024_language_length_idx;


--
-- Name: submissions_2024_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_2024_language_point_idx;


--
-- Name: submissions_2024_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_2024_problem_id_epoch_second_idx;


--
-- Name: submissions_2024_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_2024_problem_id_execution_time_idx;


--
-- Name: submissions_2024_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_2024_problem_id_length_idx;


--
-- Name: submissions_2024_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_2024_problem_id_point_idx;


--
-- Name: submissions_2024_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_2024_result_epoch_second_idx;


--
-- Name: submissions_2024_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_2024_result_execution_time_idx;


--
-- Name: submissions_2024_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_2024_result_length_idx;


--
-- Name: submissions_2024_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_2024_result_point_idx;


--
-- Name: submissions_2024_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_2024_user_id_epoch_second_idx;


--
-- Name: submissions_2024_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_2024_user_id_execution_time_idx;


--
-- Name: submissions_2024_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_2024_user_id_length_idx;


--
-- Name: submissions_2024_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_2024_user_id_point_idx;


--
-- Name: submissions_2024_user_id_result_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_result_index ATTACH PARTITION public.submissions_2024_user_id_result_idx;


--
-- Name: submissions_2025_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_2025_contest_id_epoch_second_idx;


--
-- Name: submissions_2025_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_2025_contest_id_execution_time_idx;


--
-- Name: submissions_2025_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_2025_contest_id_length_idx;


--
-- Name: submissions_2025_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_2025_contest_id_point_idx;


--
-- Name: submissions_2025_epoch_second_id_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_unique ATTACH PARTITION public.submissions_2025_epoch_second_id_idx;


--
-- Name: submissions_2025_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_2025_language_epoch_second_idx;


--
-- Name: submissions_2025_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_2025_language_execution_time_idx;


--
-- Name: submissions_2025_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_2025_language_length_idx;


--
-- Name: submissions_2025_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_2025_language_point_idx;


--
-- Name: submissions_2025_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_2025_problem_id_epoch_second_idx;


--
-- Name: submissions_2025_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_2025_problem_id_execution_time_idx;


--
-- Name: submissions_2025_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_2025_problem_id_length_idx;


--
-- Name: submissions_2025_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_2025_problem_id_point_idx;


--
-- Name: submissions_2025_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_2025_result_epoch_second_idx;


--
-- Name: submissions_2025_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_2025_result_execution_time_idx;


--
-- Name: submissions_2025_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_2025_result_length_idx;


--
-- Name: submissions_2025_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_2025_result_point_idx;


--
-- Name: submissions_2025_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_2025_user_id_epoch_second_idx;


--
-- Name: submissions_2025_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_2025_user_id_execution_time_idx;


--
-- Name: submissions_2025_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_2025_user_id_length_idx;


--
-- Name: submissions_2025_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_2025_user_id_point_idx;


--
-- Name: submissions_2025_user_id_result_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_result_index ATTACH PARTITION public.submissions_2025_user_id_result_idx;


--
-- Name: submissions_2026_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_2026_contest_id_epoch_second_idx;


--
-- Name: submissions_2026_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_2026_contest_id_execution_time_idx;


--
-- Name: submissions_2026_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_2026_contest_id_length_idx;


--
-- Name: submissions_2026_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_2026_contest_id_point_idx;


--
-- Name: submissions_2026_epoch_second_id_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_unique ATTACH PARTITION public.submissions_2026_epoch_second_id_idx;


--
-- Name: submissions_2026_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_2026_language_epoch_second_idx;


--
-- Name: submissions_2026_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_2026_language_execution_time_idx;


--
-- Name: submissions_2026_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_2026_language_length_idx;


--
-- Name: submissions_2026_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_2026_language_point_idx;


--
-- Name: submissions_2026_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_2026_problem_id_epoch_second_idx;


--
-- Name: submissions_2026_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_2026_problem_id_execution_time_idx;


--
-- Name: submissions_2026_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_2026_problem_id_length_idx;


--
-- Name: submissions_2026_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_2026_problem_id_point_idx;


--
-- Name: submissions_2026_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_2026_result_epoch_second_idx;


--
-- Name: submissions_2026_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_2026_result_execution_time_idx;


--
-- Name: submissions_2026_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_2026_result_length_idx;


--
-- Name: submissions_2026_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_2026_result_point_idx;


--
-- Name: submissions_2026_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_2026_user_id_epoch_second_idx;


--
-- Name: submissions_2026_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_2026_user_id_execution_time_idx;


--
-- Name: submissions_2026_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_2026_user_id_length_idx;


--
-- Name: submissions_2026_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_2026_user_id_point_idx;


--
-- Name: submissions_2026_user_id_result_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_result_index ATTACH PARTITION public.submissions_2026_user_id_result_idx;


--
-- Name: submissions_2027_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_2027_contest_id_epoch_second_idx;


--
-- Name: submissions_2027_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_2027_contest_id_execution_time_idx;


--
-- Name: submissions_2027_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_2027_contest_id_length_idx;


--
-- Name: submissions_2027_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_2027_contest_id_point_idx;


--
-- Name: submissions_2027_epoch_second_id_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_unique ATTACH PARTITION public.submissions_2027_epoch_second_id_idx;


--
-- Name: submissions_2027_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_2027_language_epoch_second_idx;


--
-- Name: submissions_2027_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_2027_language_execution_time_idx;


--
-- Name: submissions_2027_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_2027_language_length_idx;


--
-- Name: submissions_2027_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_2027_language_point_idx;


--
-- Name: submissions_2027_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_2027_problem_id_epoch_second_idx;


--
-- Name: submissions_2027_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_2027_problem_id_execution_time_idx;


--
-- Name: submissions_2027_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_2027_problem_id_length_idx;


--
-- Name: submissions_2027_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_2027_problem_id_point_idx;


--
-- Name: submissions_2027_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_2027_result_epoch_second_idx;


--
-- Name: submissions_2027_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_2027_result_execution_time_idx;


--
-- Name: submissions_2027_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_2027_result_length_idx;


--
-- Name: submissions_2027_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_2027_result_point_idx;


--
-- Name: submissions_2027_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_2027_user_id_epoch_second_idx;


--
-- Name: submissions_2027_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_2027_user_id_execution_time_idx;


--
-- Name: submissions_2027_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_2027_user_id_length_idx;


--
-- Name: submissions_2027_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_2027_user_id_point_idx;


--
-- Name: submissions_2027_user_id_result_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_result_index ATTACH PARTITION public.submissions_2027_user_id_result_idx;


--
-- Name: submissions_2028_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_2028_contest_id_epoch_second_idx;


--
-- Name: submissions_2028_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_2028_contest_id_execution_time_idx;


--
-- Name: submissions_2028_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_2028_contest_id_length_idx;


--
-- Name: submissions_2028_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_2028_contest_id_point_idx;


--
-- Name: submissions_2028_epoch_second_id_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_unique ATTACH PARTITION public.submissions_2028_epoch_second_id_idx;


--
-- Name: submissions_2028_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_2028_language_epoch_second_idx;


--
-- Name: submissions_2028_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_2028_language_execution_time_idx;


--
-- Name: submissions_2028_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_2028_language_length_idx;


--
-- Name: submissions_2028_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_2028_language_point_idx;


--
-- Name: submissions_2028_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_2028_problem_id_epoch_second_idx;


--
-- Name: submissions_2028_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_2028_problem_id_execution_time_idx;


--
-- Name: submissions_2028_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_2028_problem_id_length_idx;


--
-- Name: submissions_2028_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_2028_problem_id_point_idx;


--
-- Name: submissions_2028_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_2028_result_epoch_second_idx;


--
-- Name: submissions_2028_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_2028_result_execution_time_idx;


--
-- Name: submissions_2028_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_2028_result_length_idx;


--
-- Name: submissions_2028_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_2028_result_point_idx;


--
-- Name: submissions_2028_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_2028_user_id_epoch_second_idx;


--
-- Name: submissions_2028_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_2028_user_id_execution_time_idx;


--
-- Name: submissions_2028_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_2028_user_id_length_idx;


--
-- Name: submissions_2028_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_2028_user_id_point_idx;


--
-- Name: submissions_2028_user_id_result_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_result_index ATTACH PARTITION public.submissions_2028_user_id_result_idx;


--
-- Name: submissions_2029_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_2029_contest_id_epoch_second_idx;


--
-- Name: submissions_2029_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_2029_contest_id_execution_time_idx;


--
-- Name: submissions_2029_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_2029_contest_id_length_idx;


--
-- Name: submissions_2029_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_2029_contest_id_point_idx;


--
-- Name: submissions_2029_epoch_second_id_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_unique ATTACH PARTITION public.submissions_2029_epoch_second_id_idx;


--
-- Name: submissions_2029_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_2029_language_epoch_second_idx;


--
-- Name: submissions_2029_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_2029_language_execution_time_idx;


--
-- Name: submissions_2029_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_2029_language_length_idx;


--
-- Name: submissions_2029_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_2029_language_point_idx;


--
-- Name: submissions_2029_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_2029_problem_id_epoch_second_idx;


--
-- Name: submissions_2029_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_2029_problem_id_execution_time_idx;


--
-- Name: submissions_2029_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_2029_problem_id_length_idx;


--
-- Name: submissions_2029_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_2029_problem_id_point_idx;


--
-- Name: submissions_2029_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_2029_result_epoch_second_idx;


--
-- Name: submissions_2029_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_2029_result_execution_time_idx;


--
-- Name: submissions_2029_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_2029_result_length_idx;


--
-- Name: submissions_2029_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_2029_result_point_idx;


--
-- Name: submissions_2029_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_2029_user_id_epoch_second_idx;


--
-- Name: submissions_2029_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_2029_user_id_execution_time_idx;


--
-- Name: submissions_2029_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_2029_user_id_length_idx;


--
-- Name: submissions_2029_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_2029_user_id_point_idx;


--
-- Name: submissions_2029_user_id_result_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_result_index ATTACH PARTITION public.submissions_2029_user_id_result_idx;


--
-- PostgreSQL database dump complete
--


--
-- Dbmate schema migrations
--

INSERT INTO public.schema_migrations (version) VALUES
    ('20240824090834'),
    ('20240824091009'),
    ('20240824091057'),
    ('20240824091152'),
    ('20240824091232'),
    ('20240824091311'),
    ('20240824091447'),
    ('20240824091517'),
    ('20240830102505'),
    ('20240831033102'),
    ('20240831033424'),
    ('20241104132841');
