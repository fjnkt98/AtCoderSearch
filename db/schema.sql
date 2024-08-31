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
    started_at bigint NOT NULL,
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
-- Name: submissions_0_1304175600; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_0_1304175600 (
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
-- Name: submissions_1304175600_1321455600; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_1304175600_1321455600 (
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
-- Name: submissions_1321455600_1338735600; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_1321455600_1338735600 (
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
-- Name: submissions_1338735600_1356015600; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_1338735600_1356015600 (
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
-- Name: submissions_1356015600_1373295600; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_1356015600_1373295600 (
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
-- Name: submissions_1373295600_1390575600; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_1373295600_1390575600 (
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
-- Name: submissions_1390575600_1407855600; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_1390575600_1407855600 (
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
-- Name: submissions_1407855600_1425135600; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_1407855600_1425135600 (
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
-- Name: submissions_1425135600_1442415600; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_1425135600_1442415600 (
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
-- Name: submissions_1442415600_1459695600; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_1442415600_1459695600 (
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
-- Name: submissions_1459695600_1476975600; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_1459695600_1476975600 (
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
-- Name: submissions_1476975600_1494255600; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_1476975600_1494255600 (
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
-- Name: submissions_1494255600_1511535600; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_1494255600_1511535600 (
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
-- Name: submissions_1511535600_1528815600; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_1511535600_1528815600 (
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
-- Name: submissions_1528815600_1546095600; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_1528815600_1546095600 (
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
-- Name: submissions_1546095600_1563375600; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_1546095600_1563375600 (
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
-- Name: submissions_1563375600_1580655600; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_1563375600_1580655600 (
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
-- Name: submissions_1580655600_1597935600; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_1580655600_1597935600 (
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
-- Name: submissions_1597935600_1615215600; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_1597935600_1615215600 (
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
-- Name: submissions_1615215600_1632495600; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_1615215600_1632495600 (
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
-- Name: submissions_1632495600_1649775600; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_1632495600_1649775600 (
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
-- Name: submissions_1649775600_1667055600; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_1649775600_1667055600 (
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
-- Name: submissions_1667055600_1684335600; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_1667055600_1684335600 (
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
-- Name: submissions_1684335600_1701615600; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_1684335600_1701615600 (
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
-- Name: submissions_1701615600_1718895600; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_1701615600_1718895600 (
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
-- Name: submissions_1718895600_1736175600; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_1718895600_1736175600 (
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
-- Name: submissions_1736175600_1753455600; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_1736175600_1753455600 (
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
-- Name: submissions_1753455600_1770735600; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.submissions_1753455600_1770735600 (
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
-- Name: submissions_0_1304175600; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_0_1304175600 FOR VALUES FROM ('0') TO ('1304175600');


--
-- Name: submissions_1304175600_1321455600; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_1304175600_1321455600 FOR VALUES FROM ('1304175600') TO ('1321455600');


--
-- Name: submissions_1321455600_1338735600; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_1321455600_1338735600 FOR VALUES FROM ('1321455600') TO ('1338735600');


--
-- Name: submissions_1338735600_1356015600; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_1338735600_1356015600 FOR VALUES FROM ('1338735600') TO ('1356015600');


--
-- Name: submissions_1356015600_1373295600; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_1356015600_1373295600 FOR VALUES FROM ('1356015600') TO ('1373295600');


--
-- Name: submissions_1373295600_1390575600; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_1373295600_1390575600 FOR VALUES FROM ('1373295600') TO ('1390575600');


--
-- Name: submissions_1390575600_1407855600; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_1390575600_1407855600 FOR VALUES FROM ('1390575600') TO ('1407855600');


--
-- Name: submissions_1407855600_1425135600; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_1407855600_1425135600 FOR VALUES FROM ('1407855600') TO ('1425135600');


--
-- Name: submissions_1425135600_1442415600; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_1425135600_1442415600 FOR VALUES FROM ('1425135600') TO ('1442415600');


--
-- Name: submissions_1442415600_1459695600; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_1442415600_1459695600 FOR VALUES FROM ('1442415600') TO ('1459695600');


--
-- Name: submissions_1459695600_1476975600; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_1459695600_1476975600 FOR VALUES FROM ('1459695600') TO ('1476975600');


--
-- Name: submissions_1476975600_1494255600; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_1476975600_1494255600 FOR VALUES FROM ('1476975600') TO ('1494255600');


--
-- Name: submissions_1494255600_1511535600; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_1494255600_1511535600 FOR VALUES FROM ('1494255600') TO ('1511535600');


--
-- Name: submissions_1511535600_1528815600; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_1511535600_1528815600 FOR VALUES FROM ('1511535600') TO ('1528815600');


--
-- Name: submissions_1528815600_1546095600; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_1528815600_1546095600 FOR VALUES FROM ('1528815600') TO ('1546095600');


--
-- Name: submissions_1546095600_1563375600; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_1546095600_1563375600 FOR VALUES FROM ('1546095600') TO ('1563375600');


--
-- Name: submissions_1563375600_1580655600; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_1563375600_1580655600 FOR VALUES FROM ('1563375600') TO ('1580655600');


--
-- Name: submissions_1580655600_1597935600; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_1580655600_1597935600 FOR VALUES FROM ('1580655600') TO ('1597935600');


--
-- Name: submissions_1597935600_1615215600; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_1597935600_1615215600 FOR VALUES FROM ('1597935600') TO ('1615215600');


--
-- Name: submissions_1615215600_1632495600; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_1615215600_1632495600 FOR VALUES FROM ('1615215600') TO ('1632495600');


--
-- Name: submissions_1632495600_1649775600; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_1632495600_1649775600 FOR VALUES FROM ('1632495600') TO ('1649775600');


--
-- Name: submissions_1649775600_1667055600; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_1649775600_1667055600 FOR VALUES FROM ('1649775600') TO ('1667055600');


--
-- Name: submissions_1667055600_1684335600; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_1667055600_1684335600 FOR VALUES FROM ('1667055600') TO ('1684335600');


--
-- Name: submissions_1684335600_1701615600; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_1684335600_1701615600 FOR VALUES FROM ('1684335600') TO ('1701615600');


--
-- Name: submissions_1701615600_1718895600; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_1701615600_1718895600 FOR VALUES FROM ('1701615600') TO ('1718895600');


--
-- Name: submissions_1718895600_1736175600; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_1718895600_1736175600 FOR VALUES FROM ('1718895600') TO ('1736175600');


--
-- Name: submissions_1736175600_1753455600; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_1736175600_1753455600 FOR VALUES FROM ('1736175600') TO ('1753455600');


--
-- Name: submissions_1753455600_1770735600; Type: TABLE ATTACH; Schema: public; Owner: -
--

ALTER TABLE ONLY public.submissions ATTACH PARTITION public.submissions_1753455600_1770735600 FOR VALUES FROM ('1753455600') TO ('1770735600');


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
-- Name: submissions_0_1304175600_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_0_1304175600_contest_id_epoch_second_idx ON public.submissions_0_1304175600 USING btree (contest_id, epoch_second);


--
-- Name: submissions_contest_id_execution_time_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_contest_id_execution_time_index ON ONLY public.submissions USING btree (contest_id, execution_time);


--
-- Name: submissions_0_1304175600_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_0_1304175600_contest_id_execution_time_idx ON public.submissions_0_1304175600 USING btree (contest_id, execution_time);


--
-- Name: submissions_contest_id_length_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_contest_id_length_index ON ONLY public.submissions USING btree (contest_id, length);


--
-- Name: submissions_0_1304175600_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_0_1304175600_contest_id_length_idx ON public.submissions_0_1304175600 USING btree (contest_id, length);


--
-- Name: submissions_contest_id_point_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_contest_id_point_index ON ONLY public.submissions USING btree (contest_id, point);


--
-- Name: submissions_0_1304175600_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_0_1304175600_contest_id_point_idx ON public.submissions_0_1304175600 USING btree (contest_id, point);


--
-- Name: submissions_updated_at_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_updated_at_index ON ONLY public.submissions USING btree (epoch_second, updated_at);


--
-- Name: submissions_0_1304175600_epoch_second_updated_at_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_0_1304175600_epoch_second_updated_at_idx ON public.submissions_0_1304175600 USING btree (epoch_second, updated_at);


--
-- Name: submissions_execution_time_epoch_second_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_execution_time_epoch_second_index ON ONLY public.submissions USING btree (execution_time, epoch_second);


--
-- Name: submissions_0_1304175600_execution_time_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_0_1304175600_execution_time_epoch_second_idx ON public.submissions_0_1304175600 USING btree (execution_time, epoch_second);


--
-- Name: submissions_execution_time_length_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_execution_time_length_index ON ONLY public.submissions USING btree (execution_time, length);


--
-- Name: submissions_0_1304175600_execution_time_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_0_1304175600_execution_time_length_idx ON public.submissions_0_1304175600 USING btree (execution_time, length);


--
-- Name: submissions_execution_time_point_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_execution_time_point_index ON ONLY public.submissions USING btree (execution_time, point);


--
-- Name: submissions_0_1304175600_execution_time_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_0_1304175600_execution_time_point_idx ON public.submissions_0_1304175600 USING btree (execution_time, point);


--
-- Name: submissions_id_epoch_second_unique; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_id_epoch_second_unique ON ONLY public.submissions USING btree (id, epoch_second);


--
-- Name: submissions_0_1304175600_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_0_1304175600_id_epoch_second_idx ON public.submissions_0_1304175600 USING btree (id, epoch_second);


--
-- Name: submissions_language_epoch_second_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_language_epoch_second_index ON ONLY public.submissions USING btree (language, epoch_second);


--
-- Name: submissions_0_1304175600_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_0_1304175600_language_epoch_second_idx ON public.submissions_0_1304175600 USING btree (language, epoch_second);


--
-- Name: submissions_language_execution_time_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_language_execution_time_index ON ONLY public.submissions USING btree (language, execution_time);


--
-- Name: submissions_0_1304175600_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_0_1304175600_language_execution_time_idx ON public.submissions_0_1304175600 USING btree (language, execution_time);


--
-- Name: submissions_language_length_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_language_length_index ON ONLY public.submissions USING btree (language, length);


--
-- Name: submissions_0_1304175600_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_0_1304175600_language_length_idx ON public.submissions_0_1304175600 USING btree (language, length);


--
-- Name: submissions_language_point_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_language_point_index ON ONLY public.submissions USING btree (language, point);


--
-- Name: submissions_0_1304175600_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_0_1304175600_language_point_idx ON public.submissions_0_1304175600 USING btree (language, point);


--
-- Name: submissions_length_epoch_second_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_length_epoch_second_index ON ONLY public.submissions USING btree (length, epoch_second);


--
-- Name: submissions_0_1304175600_length_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_0_1304175600_length_epoch_second_idx ON public.submissions_0_1304175600 USING btree (length, epoch_second);


--
-- Name: submissions_length_execution_time_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_length_execution_time_index ON ONLY public.submissions USING btree (length, execution_time);


--
-- Name: submissions_0_1304175600_length_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_0_1304175600_length_execution_time_idx ON public.submissions_0_1304175600 USING btree (length, execution_time);


--
-- Name: submissions_length_point_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_length_point_index ON ONLY public.submissions USING btree (length, point);


--
-- Name: submissions_0_1304175600_length_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_0_1304175600_length_point_idx ON public.submissions_0_1304175600 USING btree (length, point);


--
-- Name: submissions_point_epoch_second_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_point_epoch_second_index ON ONLY public.submissions USING btree (point, epoch_second);


--
-- Name: submissions_0_1304175600_point_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_0_1304175600_point_epoch_second_idx ON public.submissions_0_1304175600 USING btree (point, epoch_second);


--
-- Name: submissions_point_execution_time_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_point_execution_time_index ON ONLY public.submissions USING btree (point, execution_time);


--
-- Name: submissions_0_1304175600_point_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_0_1304175600_point_execution_time_idx ON public.submissions_0_1304175600 USING btree (point, execution_time);


--
-- Name: submissions_point_length_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_point_length_index ON ONLY public.submissions USING btree (point, length);


--
-- Name: submissions_0_1304175600_point_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_0_1304175600_point_length_idx ON public.submissions_0_1304175600 USING btree (point, length);


--
-- Name: submissions_problem_id_epoch_second_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_problem_id_epoch_second_index ON ONLY public.submissions USING btree (problem_id, epoch_second);


--
-- Name: submissions_0_1304175600_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_0_1304175600_problem_id_epoch_second_idx ON public.submissions_0_1304175600 USING btree (problem_id, epoch_second);


--
-- Name: submissions_problem_id_execution_time_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_problem_id_execution_time_index ON ONLY public.submissions USING btree (problem_id, execution_time);


--
-- Name: submissions_0_1304175600_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_0_1304175600_problem_id_execution_time_idx ON public.submissions_0_1304175600 USING btree (problem_id, execution_time);


--
-- Name: submissions_problem_id_length_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_problem_id_length_index ON ONLY public.submissions USING btree (problem_id, length);


--
-- Name: submissions_0_1304175600_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_0_1304175600_problem_id_length_idx ON public.submissions_0_1304175600 USING btree (problem_id, length);


--
-- Name: submissions_problem_id_point_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_problem_id_point_index ON ONLY public.submissions USING btree (problem_id, point);


--
-- Name: submissions_0_1304175600_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_0_1304175600_problem_id_point_idx ON public.submissions_0_1304175600 USING btree (problem_id, point);


--
-- Name: submissions_result_epoch_second_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_result_epoch_second_index ON ONLY public.submissions USING btree (result, epoch_second);


--
-- Name: submissions_0_1304175600_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_0_1304175600_result_epoch_second_idx ON public.submissions_0_1304175600 USING btree (result, epoch_second);


--
-- Name: submissions_result_execution_time_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_result_execution_time_index ON ONLY public.submissions USING btree (result, execution_time);


--
-- Name: submissions_0_1304175600_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_0_1304175600_result_execution_time_idx ON public.submissions_0_1304175600 USING btree (result, execution_time);


--
-- Name: submissions_result_length_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_result_length_index ON ONLY public.submissions USING btree (result, length);


--
-- Name: submissions_0_1304175600_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_0_1304175600_result_length_idx ON public.submissions_0_1304175600 USING btree (result, length);


--
-- Name: submissions_result_point_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_result_point_index ON ONLY public.submissions USING btree (result, point);


--
-- Name: submissions_0_1304175600_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_0_1304175600_result_point_idx ON public.submissions_0_1304175600 USING btree (result, point);


--
-- Name: submissions_user_id_epoch_second_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_user_id_epoch_second_index ON ONLY public.submissions USING btree (user_id, epoch_second);


--
-- Name: submissions_0_1304175600_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_0_1304175600_user_id_epoch_second_idx ON public.submissions_0_1304175600 USING btree (user_id, epoch_second);


--
-- Name: submissions_user_id_execution_time_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_user_id_execution_time_index ON ONLY public.submissions USING btree (user_id, execution_time);


--
-- Name: submissions_0_1304175600_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_0_1304175600_user_id_execution_time_idx ON public.submissions_0_1304175600 USING btree (user_id, execution_time);


--
-- Name: submissions_user_id_length_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_user_id_length_index ON ONLY public.submissions USING btree (user_id, length);


--
-- Name: submissions_0_1304175600_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_0_1304175600_user_id_length_idx ON public.submissions_0_1304175600 USING btree (user_id, length);


--
-- Name: submissions_user_id_point_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_user_id_point_index ON ONLY public.submissions USING btree (user_id, point);


--
-- Name: submissions_0_1304175600_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_0_1304175600_user_id_point_idx ON public.submissions_0_1304175600 USING btree (user_id, point);


--
-- Name: submissions_1304175600_1321455600_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1304175600_1321455600_contest_id_epoch_second_idx ON public.submissions_1304175600_1321455600 USING btree (contest_id, epoch_second);


--
-- Name: submissions_1304175600_1321455600_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1304175600_1321455600_contest_id_execution_time_idx ON public.submissions_1304175600_1321455600 USING btree (contest_id, execution_time);


--
-- Name: submissions_1304175600_1321455600_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1304175600_1321455600_contest_id_length_idx ON public.submissions_1304175600_1321455600 USING btree (contest_id, length);


--
-- Name: submissions_1304175600_1321455600_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1304175600_1321455600_contest_id_point_idx ON public.submissions_1304175600_1321455600 USING btree (contest_id, point);


--
-- Name: submissions_1304175600_1321455600_epoch_second_updated_at_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1304175600_1321455600_epoch_second_updated_at_idx ON public.submissions_1304175600_1321455600 USING btree (epoch_second, updated_at);


--
-- Name: submissions_1304175600_1321455600_execution_time_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1304175600_1321455600_execution_time_length_idx ON public.submissions_1304175600_1321455600 USING btree (execution_time, length);


--
-- Name: submissions_1304175600_1321455600_execution_time_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1304175600_1321455600_execution_time_point_idx ON public.submissions_1304175600_1321455600 USING btree (execution_time, point);


--
-- Name: submissions_1304175600_1321455600_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_1304175600_1321455600_id_epoch_second_idx ON public.submissions_1304175600_1321455600 USING btree (id, epoch_second);


--
-- Name: submissions_1304175600_1321455600_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1304175600_1321455600_language_epoch_second_idx ON public.submissions_1304175600_1321455600 USING btree (language, epoch_second);


--
-- Name: submissions_1304175600_1321455600_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1304175600_1321455600_language_execution_time_idx ON public.submissions_1304175600_1321455600 USING btree (language, execution_time);


--
-- Name: submissions_1304175600_1321455600_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1304175600_1321455600_language_length_idx ON public.submissions_1304175600_1321455600 USING btree (language, length);


--
-- Name: submissions_1304175600_1321455600_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1304175600_1321455600_language_point_idx ON public.submissions_1304175600_1321455600 USING btree (language, point);


--
-- Name: submissions_1304175600_1321455600_length_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1304175600_1321455600_length_epoch_second_idx ON public.submissions_1304175600_1321455600 USING btree (length, epoch_second);


--
-- Name: submissions_1304175600_1321455600_length_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1304175600_1321455600_length_execution_time_idx ON public.submissions_1304175600_1321455600 USING btree (length, execution_time);


--
-- Name: submissions_1304175600_1321455600_length_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1304175600_1321455600_length_point_idx ON public.submissions_1304175600_1321455600 USING btree (length, point);


--
-- Name: submissions_1304175600_1321455600_point_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1304175600_1321455600_point_epoch_second_idx ON public.submissions_1304175600_1321455600 USING btree (point, epoch_second);


--
-- Name: submissions_1304175600_1321455600_point_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1304175600_1321455600_point_execution_time_idx ON public.submissions_1304175600_1321455600 USING btree (point, execution_time);


--
-- Name: submissions_1304175600_1321455600_point_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1304175600_1321455600_point_length_idx ON public.submissions_1304175600_1321455600 USING btree (point, length);


--
-- Name: submissions_1304175600_1321455600_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1304175600_1321455600_problem_id_epoch_second_idx ON public.submissions_1304175600_1321455600 USING btree (problem_id, epoch_second);


--
-- Name: submissions_1304175600_1321455600_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1304175600_1321455600_problem_id_execution_time_idx ON public.submissions_1304175600_1321455600 USING btree (problem_id, execution_time);


--
-- Name: submissions_1304175600_1321455600_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1304175600_1321455600_problem_id_length_idx ON public.submissions_1304175600_1321455600 USING btree (problem_id, length);


--
-- Name: submissions_1304175600_1321455600_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1304175600_1321455600_problem_id_point_idx ON public.submissions_1304175600_1321455600 USING btree (problem_id, point);


--
-- Name: submissions_1304175600_1321455600_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1304175600_1321455600_result_epoch_second_idx ON public.submissions_1304175600_1321455600 USING btree (result, epoch_second);


--
-- Name: submissions_1304175600_1321455600_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1304175600_1321455600_result_execution_time_idx ON public.submissions_1304175600_1321455600 USING btree (result, execution_time);


--
-- Name: submissions_1304175600_1321455600_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1304175600_1321455600_result_length_idx ON public.submissions_1304175600_1321455600 USING btree (result, length);


--
-- Name: submissions_1304175600_1321455600_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1304175600_1321455600_result_point_idx ON public.submissions_1304175600_1321455600 USING btree (result, point);


--
-- Name: submissions_1304175600_1321455600_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1304175600_1321455600_user_id_epoch_second_idx ON public.submissions_1304175600_1321455600 USING btree (user_id, epoch_second);


--
-- Name: submissions_1304175600_1321455600_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1304175600_1321455600_user_id_execution_time_idx ON public.submissions_1304175600_1321455600 USING btree (user_id, execution_time);


--
-- Name: submissions_1304175600_1321455600_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1304175600_1321455600_user_id_length_idx ON public.submissions_1304175600_1321455600 USING btree (user_id, length);


--
-- Name: submissions_1304175600_1321455600_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1304175600_1321455600_user_id_point_idx ON public.submissions_1304175600_1321455600 USING btree (user_id, point);


--
-- Name: submissions_1304175600_13214556_execution_time_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1304175600_13214556_execution_time_epoch_second_idx ON public.submissions_1304175600_1321455600 USING btree (execution_time, epoch_second);


--
-- Name: submissions_1321455600_1338735600_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1321455600_1338735600_contest_id_epoch_second_idx ON public.submissions_1321455600_1338735600 USING btree (contest_id, epoch_second);


--
-- Name: submissions_1321455600_1338735600_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1321455600_1338735600_contest_id_execution_time_idx ON public.submissions_1321455600_1338735600 USING btree (contest_id, execution_time);


--
-- Name: submissions_1321455600_1338735600_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1321455600_1338735600_contest_id_length_idx ON public.submissions_1321455600_1338735600 USING btree (contest_id, length);


--
-- Name: submissions_1321455600_1338735600_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1321455600_1338735600_contest_id_point_idx ON public.submissions_1321455600_1338735600 USING btree (contest_id, point);


--
-- Name: submissions_1321455600_1338735600_epoch_second_updated_at_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1321455600_1338735600_epoch_second_updated_at_idx ON public.submissions_1321455600_1338735600 USING btree (epoch_second, updated_at);


--
-- Name: submissions_1321455600_1338735600_execution_time_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1321455600_1338735600_execution_time_length_idx ON public.submissions_1321455600_1338735600 USING btree (execution_time, length);


--
-- Name: submissions_1321455600_1338735600_execution_time_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1321455600_1338735600_execution_time_point_idx ON public.submissions_1321455600_1338735600 USING btree (execution_time, point);


--
-- Name: submissions_1321455600_1338735600_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_1321455600_1338735600_id_epoch_second_idx ON public.submissions_1321455600_1338735600 USING btree (id, epoch_second);


--
-- Name: submissions_1321455600_1338735600_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1321455600_1338735600_language_epoch_second_idx ON public.submissions_1321455600_1338735600 USING btree (language, epoch_second);


--
-- Name: submissions_1321455600_1338735600_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1321455600_1338735600_language_execution_time_idx ON public.submissions_1321455600_1338735600 USING btree (language, execution_time);


--
-- Name: submissions_1321455600_1338735600_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1321455600_1338735600_language_length_idx ON public.submissions_1321455600_1338735600 USING btree (language, length);


--
-- Name: submissions_1321455600_1338735600_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1321455600_1338735600_language_point_idx ON public.submissions_1321455600_1338735600 USING btree (language, point);


--
-- Name: submissions_1321455600_1338735600_length_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1321455600_1338735600_length_epoch_second_idx ON public.submissions_1321455600_1338735600 USING btree (length, epoch_second);


--
-- Name: submissions_1321455600_1338735600_length_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1321455600_1338735600_length_execution_time_idx ON public.submissions_1321455600_1338735600 USING btree (length, execution_time);


--
-- Name: submissions_1321455600_1338735600_length_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1321455600_1338735600_length_point_idx ON public.submissions_1321455600_1338735600 USING btree (length, point);


--
-- Name: submissions_1321455600_1338735600_point_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1321455600_1338735600_point_epoch_second_idx ON public.submissions_1321455600_1338735600 USING btree (point, epoch_second);


--
-- Name: submissions_1321455600_1338735600_point_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1321455600_1338735600_point_execution_time_idx ON public.submissions_1321455600_1338735600 USING btree (point, execution_time);


--
-- Name: submissions_1321455600_1338735600_point_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1321455600_1338735600_point_length_idx ON public.submissions_1321455600_1338735600 USING btree (point, length);


--
-- Name: submissions_1321455600_1338735600_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1321455600_1338735600_problem_id_epoch_second_idx ON public.submissions_1321455600_1338735600 USING btree (problem_id, epoch_second);


--
-- Name: submissions_1321455600_1338735600_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1321455600_1338735600_problem_id_execution_time_idx ON public.submissions_1321455600_1338735600 USING btree (problem_id, execution_time);


--
-- Name: submissions_1321455600_1338735600_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1321455600_1338735600_problem_id_length_idx ON public.submissions_1321455600_1338735600 USING btree (problem_id, length);


--
-- Name: submissions_1321455600_1338735600_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1321455600_1338735600_problem_id_point_idx ON public.submissions_1321455600_1338735600 USING btree (problem_id, point);


--
-- Name: submissions_1321455600_1338735600_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1321455600_1338735600_result_epoch_second_idx ON public.submissions_1321455600_1338735600 USING btree (result, epoch_second);


--
-- Name: submissions_1321455600_1338735600_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1321455600_1338735600_result_execution_time_idx ON public.submissions_1321455600_1338735600 USING btree (result, execution_time);


--
-- Name: submissions_1321455600_1338735600_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1321455600_1338735600_result_length_idx ON public.submissions_1321455600_1338735600 USING btree (result, length);


--
-- Name: submissions_1321455600_1338735600_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1321455600_1338735600_result_point_idx ON public.submissions_1321455600_1338735600 USING btree (result, point);


--
-- Name: submissions_1321455600_1338735600_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1321455600_1338735600_user_id_epoch_second_idx ON public.submissions_1321455600_1338735600 USING btree (user_id, epoch_second);


--
-- Name: submissions_1321455600_1338735600_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1321455600_1338735600_user_id_execution_time_idx ON public.submissions_1321455600_1338735600 USING btree (user_id, execution_time);


--
-- Name: submissions_1321455600_1338735600_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1321455600_1338735600_user_id_length_idx ON public.submissions_1321455600_1338735600 USING btree (user_id, length);


--
-- Name: submissions_1321455600_1338735600_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1321455600_1338735600_user_id_point_idx ON public.submissions_1321455600_1338735600 USING btree (user_id, point);


--
-- Name: submissions_1321455600_13387356_execution_time_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1321455600_13387356_execution_time_epoch_second_idx ON public.submissions_1321455600_1338735600 USING btree (execution_time, epoch_second);


--
-- Name: submissions_1338735600_1356015600_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1338735600_1356015600_contest_id_epoch_second_idx ON public.submissions_1338735600_1356015600 USING btree (contest_id, epoch_second);


--
-- Name: submissions_1338735600_1356015600_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1338735600_1356015600_contest_id_execution_time_idx ON public.submissions_1338735600_1356015600 USING btree (contest_id, execution_time);


--
-- Name: submissions_1338735600_1356015600_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1338735600_1356015600_contest_id_length_idx ON public.submissions_1338735600_1356015600 USING btree (contest_id, length);


--
-- Name: submissions_1338735600_1356015600_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1338735600_1356015600_contest_id_point_idx ON public.submissions_1338735600_1356015600 USING btree (contest_id, point);


--
-- Name: submissions_1338735600_1356015600_epoch_second_updated_at_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1338735600_1356015600_epoch_second_updated_at_idx ON public.submissions_1338735600_1356015600 USING btree (epoch_second, updated_at);


--
-- Name: submissions_1338735600_1356015600_execution_time_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1338735600_1356015600_execution_time_length_idx ON public.submissions_1338735600_1356015600 USING btree (execution_time, length);


--
-- Name: submissions_1338735600_1356015600_execution_time_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1338735600_1356015600_execution_time_point_idx ON public.submissions_1338735600_1356015600 USING btree (execution_time, point);


--
-- Name: submissions_1338735600_1356015600_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_1338735600_1356015600_id_epoch_second_idx ON public.submissions_1338735600_1356015600 USING btree (id, epoch_second);


--
-- Name: submissions_1338735600_1356015600_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1338735600_1356015600_language_epoch_second_idx ON public.submissions_1338735600_1356015600 USING btree (language, epoch_second);


--
-- Name: submissions_1338735600_1356015600_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1338735600_1356015600_language_execution_time_idx ON public.submissions_1338735600_1356015600 USING btree (language, execution_time);


--
-- Name: submissions_1338735600_1356015600_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1338735600_1356015600_language_length_idx ON public.submissions_1338735600_1356015600 USING btree (language, length);


--
-- Name: submissions_1338735600_1356015600_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1338735600_1356015600_language_point_idx ON public.submissions_1338735600_1356015600 USING btree (language, point);


--
-- Name: submissions_1338735600_1356015600_length_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1338735600_1356015600_length_epoch_second_idx ON public.submissions_1338735600_1356015600 USING btree (length, epoch_second);


--
-- Name: submissions_1338735600_1356015600_length_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1338735600_1356015600_length_execution_time_idx ON public.submissions_1338735600_1356015600 USING btree (length, execution_time);


--
-- Name: submissions_1338735600_1356015600_length_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1338735600_1356015600_length_point_idx ON public.submissions_1338735600_1356015600 USING btree (length, point);


--
-- Name: submissions_1338735600_1356015600_point_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1338735600_1356015600_point_epoch_second_idx ON public.submissions_1338735600_1356015600 USING btree (point, epoch_second);


--
-- Name: submissions_1338735600_1356015600_point_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1338735600_1356015600_point_execution_time_idx ON public.submissions_1338735600_1356015600 USING btree (point, execution_time);


--
-- Name: submissions_1338735600_1356015600_point_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1338735600_1356015600_point_length_idx ON public.submissions_1338735600_1356015600 USING btree (point, length);


--
-- Name: submissions_1338735600_1356015600_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1338735600_1356015600_problem_id_epoch_second_idx ON public.submissions_1338735600_1356015600 USING btree (problem_id, epoch_second);


--
-- Name: submissions_1338735600_1356015600_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1338735600_1356015600_problem_id_execution_time_idx ON public.submissions_1338735600_1356015600 USING btree (problem_id, execution_time);


--
-- Name: submissions_1338735600_1356015600_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1338735600_1356015600_problem_id_length_idx ON public.submissions_1338735600_1356015600 USING btree (problem_id, length);


--
-- Name: submissions_1338735600_1356015600_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1338735600_1356015600_problem_id_point_idx ON public.submissions_1338735600_1356015600 USING btree (problem_id, point);


--
-- Name: submissions_1338735600_1356015600_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1338735600_1356015600_result_epoch_second_idx ON public.submissions_1338735600_1356015600 USING btree (result, epoch_second);


--
-- Name: submissions_1338735600_1356015600_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1338735600_1356015600_result_execution_time_idx ON public.submissions_1338735600_1356015600 USING btree (result, execution_time);


--
-- Name: submissions_1338735600_1356015600_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1338735600_1356015600_result_length_idx ON public.submissions_1338735600_1356015600 USING btree (result, length);


--
-- Name: submissions_1338735600_1356015600_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1338735600_1356015600_result_point_idx ON public.submissions_1338735600_1356015600 USING btree (result, point);


--
-- Name: submissions_1338735600_1356015600_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1338735600_1356015600_user_id_epoch_second_idx ON public.submissions_1338735600_1356015600 USING btree (user_id, epoch_second);


--
-- Name: submissions_1338735600_1356015600_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1338735600_1356015600_user_id_execution_time_idx ON public.submissions_1338735600_1356015600 USING btree (user_id, execution_time);


--
-- Name: submissions_1338735600_1356015600_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1338735600_1356015600_user_id_length_idx ON public.submissions_1338735600_1356015600 USING btree (user_id, length);


--
-- Name: submissions_1338735600_1356015600_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1338735600_1356015600_user_id_point_idx ON public.submissions_1338735600_1356015600 USING btree (user_id, point);


--
-- Name: submissions_1338735600_13560156_execution_time_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1338735600_13560156_execution_time_epoch_second_idx ON public.submissions_1338735600_1356015600 USING btree (execution_time, epoch_second);


--
-- Name: submissions_1356015600_1373295600_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1356015600_1373295600_contest_id_epoch_second_idx ON public.submissions_1356015600_1373295600 USING btree (contest_id, epoch_second);


--
-- Name: submissions_1356015600_1373295600_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1356015600_1373295600_contest_id_execution_time_idx ON public.submissions_1356015600_1373295600 USING btree (contest_id, execution_time);


--
-- Name: submissions_1356015600_1373295600_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1356015600_1373295600_contest_id_length_idx ON public.submissions_1356015600_1373295600 USING btree (contest_id, length);


--
-- Name: submissions_1356015600_1373295600_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1356015600_1373295600_contest_id_point_idx ON public.submissions_1356015600_1373295600 USING btree (contest_id, point);


--
-- Name: submissions_1356015600_1373295600_epoch_second_updated_at_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1356015600_1373295600_epoch_second_updated_at_idx ON public.submissions_1356015600_1373295600 USING btree (epoch_second, updated_at);


--
-- Name: submissions_1356015600_1373295600_execution_time_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1356015600_1373295600_execution_time_length_idx ON public.submissions_1356015600_1373295600 USING btree (execution_time, length);


--
-- Name: submissions_1356015600_1373295600_execution_time_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1356015600_1373295600_execution_time_point_idx ON public.submissions_1356015600_1373295600 USING btree (execution_time, point);


--
-- Name: submissions_1356015600_1373295600_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_1356015600_1373295600_id_epoch_second_idx ON public.submissions_1356015600_1373295600 USING btree (id, epoch_second);


--
-- Name: submissions_1356015600_1373295600_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1356015600_1373295600_language_epoch_second_idx ON public.submissions_1356015600_1373295600 USING btree (language, epoch_second);


--
-- Name: submissions_1356015600_1373295600_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1356015600_1373295600_language_execution_time_idx ON public.submissions_1356015600_1373295600 USING btree (language, execution_time);


--
-- Name: submissions_1356015600_1373295600_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1356015600_1373295600_language_length_idx ON public.submissions_1356015600_1373295600 USING btree (language, length);


--
-- Name: submissions_1356015600_1373295600_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1356015600_1373295600_language_point_idx ON public.submissions_1356015600_1373295600 USING btree (language, point);


--
-- Name: submissions_1356015600_1373295600_length_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1356015600_1373295600_length_epoch_second_idx ON public.submissions_1356015600_1373295600 USING btree (length, epoch_second);


--
-- Name: submissions_1356015600_1373295600_length_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1356015600_1373295600_length_execution_time_idx ON public.submissions_1356015600_1373295600 USING btree (length, execution_time);


--
-- Name: submissions_1356015600_1373295600_length_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1356015600_1373295600_length_point_idx ON public.submissions_1356015600_1373295600 USING btree (length, point);


--
-- Name: submissions_1356015600_1373295600_point_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1356015600_1373295600_point_epoch_second_idx ON public.submissions_1356015600_1373295600 USING btree (point, epoch_second);


--
-- Name: submissions_1356015600_1373295600_point_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1356015600_1373295600_point_execution_time_idx ON public.submissions_1356015600_1373295600 USING btree (point, execution_time);


--
-- Name: submissions_1356015600_1373295600_point_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1356015600_1373295600_point_length_idx ON public.submissions_1356015600_1373295600 USING btree (point, length);


--
-- Name: submissions_1356015600_1373295600_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1356015600_1373295600_problem_id_epoch_second_idx ON public.submissions_1356015600_1373295600 USING btree (problem_id, epoch_second);


--
-- Name: submissions_1356015600_1373295600_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1356015600_1373295600_problem_id_execution_time_idx ON public.submissions_1356015600_1373295600 USING btree (problem_id, execution_time);


--
-- Name: submissions_1356015600_1373295600_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1356015600_1373295600_problem_id_length_idx ON public.submissions_1356015600_1373295600 USING btree (problem_id, length);


--
-- Name: submissions_1356015600_1373295600_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1356015600_1373295600_problem_id_point_idx ON public.submissions_1356015600_1373295600 USING btree (problem_id, point);


--
-- Name: submissions_1356015600_1373295600_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1356015600_1373295600_result_epoch_second_idx ON public.submissions_1356015600_1373295600 USING btree (result, epoch_second);


--
-- Name: submissions_1356015600_1373295600_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1356015600_1373295600_result_execution_time_idx ON public.submissions_1356015600_1373295600 USING btree (result, execution_time);


--
-- Name: submissions_1356015600_1373295600_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1356015600_1373295600_result_length_idx ON public.submissions_1356015600_1373295600 USING btree (result, length);


--
-- Name: submissions_1356015600_1373295600_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1356015600_1373295600_result_point_idx ON public.submissions_1356015600_1373295600 USING btree (result, point);


--
-- Name: submissions_1356015600_1373295600_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1356015600_1373295600_user_id_epoch_second_idx ON public.submissions_1356015600_1373295600 USING btree (user_id, epoch_second);


--
-- Name: submissions_1356015600_1373295600_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1356015600_1373295600_user_id_execution_time_idx ON public.submissions_1356015600_1373295600 USING btree (user_id, execution_time);


--
-- Name: submissions_1356015600_1373295600_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1356015600_1373295600_user_id_length_idx ON public.submissions_1356015600_1373295600 USING btree (user_id, length);


--
-- Name: submissions_1356015600_1373295600_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1356015600_1373295600_user_id_point_idx ON public.submissions_1356015600_1373295600 USING btree (user_id, point);


--
-- Name: submissions_1356015600_13732956_execution_time_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1356015600_13732956_execution_time_epoch_second_idx ON public.submissions_1356015600_1373295600 USING btree (execution_time, epoch_second);


--
-- Name: submissions_1373295600_1390575600_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1373295600_1390575600_contest_id_epoch_second_idx ON public.submissions_1373295600_1390575600 USING btree (contest_id, epoch_second);


--
-- Name: submissions_1373295600_1390575600_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1373295600_1390575600_contest_id_execution_time_idx ON public.submissions_1373295600_1390575600 USING btree (contest_id, execution_time);


--
-- Name: submissions_1373295600_1390575600_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1373295600_1390575600_contest_id_length_idx ON public.submissions_1373295600_1390575600 USING btree (contest_id, length);


--
-- Name: submissions_1373295600_1390575600_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1373295600_1390575600_contest_id_point_idx ON public.submissions_1373295600_1390575600 USING btree (contest_id, point);


--
-- Name: submissions_1373295600_1390575600_epoch_second_updated_at_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1373295600_1390575600_epoch_second_updated_at_idx ON public.submissions_1373295600_1390575600 USING btree (epoch_second, updated_at);


--
-- Name: submissions_1373295600_1390575600_execution_time_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1373295600_1390575600_execution_time_length_idx ON public.submissions_1373295600_1390575600 USING btree (execution_time, length);


--
-- Name: submissions_1373295600_1390575600_execution_time_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1373295600_1390575600_execution_time_point_idx ON public.submissions_1373295600_1390575600 USING btree (execution_time, point);


--
-- Name: submissions_1373295600_1390575600_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_1373295600_1390575600_id_epoch_second_idx ON public.submissions_1373295600_1390575600 USING btree (id, epoch_second);


--
-- Name: submissions_1373295600_1390575600_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1373295600_1390575600_language_epoch_second_idx ON public.submissions_1373295600_1390575600 USING btree (language, epoch_second);


--
-- Name: submissions_1373295600_1390575600_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1373295600_1390575600_language_execution_time_idx ON public.submissions_1373295600_1390575600 USING btree (language, execution_time);


--
-- Name: submissions_1373295600_1390575600_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1373295600_1390575600_language_length_idx ON public.submissions_1373295600_1390575600 USING btree (language, length);


--
-- Name: submissions_1373295600_1390575600_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1373295600_1390575600_language_point_idx ON public.submissions_1373295600_1390575600 USING btree (language, point);


--
-- Name: submissions_1373295600_1390575600_length_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1373295600_1390575600_length_epoch_second_idx ON public.submissions_1373295600_1390575600 USING btree (length, epoch_second);


--
-- Name: submissions_1373295600_1390575600_length_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1373295600_1390575600_length_execution_time_idx ON public.submissions_1373295600_1390575600 USING btree (length, execution_time);


--
-- Name: submissions_1373295600_1390575600_length_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1373295600_1390575600_length_point_idx ON public.submissions_1373295600_1390575600 USING btree (length, point);


--
-- Name: submissions_1373295600_1390575600_point_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1373295600_1390575600_point_epoch_second_idx ON public.submissions_1373295600_1390575600 USING btree (point, epoch_second);


--
-- Name: submissions_1373295600_1390575600_point_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1373295600_1390575600_point_execution_time_idx ON public.submissions_1373295600_1390575600 USING btree (point, execution_time);


--
-- Name: submissions_1373295600_1390575600_point_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1373295600_1390575600_point_length_idx ON public.submissions_1373295600_1390575600 USING btree (point, length);


--
-- Name: submissions_1373295600_1390575600_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1373295600_1390575600_problem_id_epoch_second_idx ON public.submissions_1373295600_1390575600 USING btree (problem_id, epoch_second);


--
-- Name: submissions_1373295600_1390575600_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1373295600_1390575600_problem_id_execution_time_idx ON public.submissions_1373295600_1390575600 USING btree (problem_id, execution_time);


--
-- Name: submissions_1373295600_1390575600_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1373295600_1390575600_problem_id_length_idx ON public.submissions_1373295600_1390575600 USING btree (problem_id, length);


--
-- Name: submissions_1373295600_1390575600_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1373295600_1390575600_problem_id_point_idx ON public.submissions_1373295600_1390575600 USING btree (problem_id, point);


--
-- Name: submissions_1373295600_1390575600_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1373295600_1390575600_result_epoch_second_idx ON public.submissions_1373295600_1390575600 USING btree (result, epoch_second);


--
-- Name: submissions_1373295600_1390575600_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1373295600_1390575600_result_execution_time_idx ON public.submissions_1373295600_1390575600 USING btree (result, execution_time);


--
-- Name: submissions_1373295600_1390575600_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1373295600_1390575600_result_length_idx ON public.submissions_1373295600_1390575600 USING btree (result, length);


--
-- Name: submissions_1373295600_1390575600_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1373295600_1390575600_result_point_idx ON public.submissions_1373295600_1390575600 USING btree (result, point);


--
-- Name: submissions_1373295600_1390575600_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1373295600_1390575600_user_id_epoch_second_idx ON public.submissions_1373295600_1390575600 USING btree (user_id, epoch_second);


--
-- Name: submissions_1373295600_1390575600_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1373295600_1390575600_user_id_execution_time_idx ON public.submissions_1373295600_1390575600 USING btree (user_id, execution_time);


--
-- Name: submissions_1373295600_1390575600_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1373295600_1390575600_user_id_length_idx ON public.submissions_1373295600_1390575600 USING btree (user_id, length);


--
-- Name: submissions_1373295600_1390575600_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1373295600_1390575600_user_id_point_idx ON public.submissions_1373295600_1390575600 USING btree (user_id, point);


--
-- Name: submissions_1373295600_13905756_execution_time_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1373295600_13905756_execution_time_epoch_second_idx ON public.submissions_1373295600_1390575600 USING btree (execution_time, epoch_second);


--
-- Name: submissions_1390575600_1407855600_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1390575600_1407855600_contest_id_epoch_second_idx ON public.submissions_1390575600_1407855600 USING btree (contest_id, epoch_second);


--
-- Name: submissions_1390575600_1407855600_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1390575600_1407855600_contest_id_execution_time_idx ON public.submissions_1390575600_1407855600 USING btree (contest_id, execution_time);


--
-- Name: submissions_1390575600_1407855600_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1390575600_1407855600_contest_id_length_idx ON public.submissions_1390575600_1407855600 USING btree (contest_id, length);


--
-- Name: submissions_1390575600_1407855600_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1390575600_1407855600_contest_id_point_idx ON public.submissions_1390575600_1407855600 USING btree (contest_id, point);


--
-- Name: submissions_1390575600_1407855600_epoch_second_updated_at_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1390575600_1407855600_epoch_second_updated_at_idx ON public.submissions_1390575600_1407855600 USING btree (epoch_second, updated_at);


--
-- Name: submissions_1390575600_1407855600_execution_time_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1390575600_1407855600_execution_time_length_idx ON public.submissions_1390575600_1407855600 USING btree (execution_time, length);


--
-- Name: submissions_1390575600_1407855600_execution_time_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1390575600_1407855600_execution_time_point_idx ON public.submissions_1390575600_1407855600 USING btree (execution_time, point);


--
-- Name: submissions_1390575600_1407855600_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_1390575600_1407855600_id_epoch_second_idx ON public.submissions_1390575600_1407855600 USING btree (id, epoch_second);


--
-- Name: submissions_1390575600_1407855600_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1390575600_1407855600_language_epoch_second_idx ON public.submissions_1390575600_1407855600 USING btree (language, epoch_second);


--
-- Name: submissions_1390575600_1407855600_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1390575600_1407855600_language_execution_time_idx ON public.submissions_1390575600_1407855600 USING btree (language, execution_time);


--
-- Name: submissions_1390575600_1407855600_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1390575600_1407855600_language_length_idx ON public.submissions_1390575600_1407855600 USING btree (language, length);


--
-- Name: submissions_1390575600_1407855600_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1390575600_1407855600_language_point_idx ON public.submissions_1390575600_1407855600 USING btree (language, point);


--
-- Name: submissions_1390575600_1407855600_length_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1390575600_1407855600_length_epoch_second_idx ON public.submissions_1390575600_1407855600 USING btree (length, epoch_second);


--
-- Name: submissions_1390575600_1407855600_length_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1390575600_1407855600_length_execution_time_idx ON public.submissions_1390575600_1407855600 USING btree (length, execution_time);


--
-- Name: submissions_1390575600_1407855600_length_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1390575600_1407855600_length_point_idx ON public.submissions_1390575600_1407855600 USING btree (length, point);


--
-- Name: submissions_1390575600_1407855600_point_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1390575600_1407855600_point_epoch_second_idx ON public.submissions_1390575600_1407855600 USING btree (point, epoch_second);


--
-- Name: submissions_1390575600_1407855600_point_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1390575600_1407855600_point_execution_time_idx ON public.submissions_1390575600_1407855600 USING btree (point, execution_time);


--
-- Name: submissions_1390575600_1407855600_point_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1390575600_1407855600_point_length_idx ON public.submissions_1390575600_1407855600 USING btree (point, length);


--
-- Name: submissions_1390575600_1407855600_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1390575600_1407855600_problem_id_epoch_second_idx ON public.submissions_1390575600_1407855600 USING btree (problem_id, epoch_second);


--
-- Name: submissions_1390575600_1407855600_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1390575600_1407855600_problem_id_execution_time_idx ON public.submissions_1390575600_1407855600 USING btree (problem_id, execution_time);


--
-- Name: submissions_1390575600_1407855600_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1390575600_1407855600_problem_id_length_idx ON public.submissions_1390575600_1407855600 USING btree (problem_id, length);


--
-- Name: submissions_1390575600_1407855600_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1390575600_1407855600_problem_id_point_idx ON public.submissions_1390575600_1407855600 USING btree (problem_id, point);


--
-- Name: submissions_1390575600_1407855600_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1390575600_1407855600_result_epoch_second_idx ON public.submissions_1390575600_1407855600 USING btree (result, epoch_second);


--
-- Name: submissions_1390575600_1407855600_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1390575600_1407855600_result_execution_time_idx ON public.submissions_1390575600_1407855600 USING btree (result, execution_time);


--
-- Name: submissions_1390575600_1407855600_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1390575600_1407855600_result_length_idx ON public.submissions_1390575600_1407855600 USING btree (result, length);


--
-- Name: submissions_1390575600_1407855600_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1390575600_1407855600_result_point_idx ON public.submissions_1390575600_1407855600 USING btree (result, point);


--
-- Name: submissions_1390575600_1407855600_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1390575600_1407855600_user_id_epoch_second_idx ON public.submissions_1390575600_1407855600 USING btree (user_id, epoch_second);


--
-- Name: submissions_1390575600_1407855600_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1390575600_1407855600_user_id_execution_time_idx ON public.submissions_1390575600_1407855600 USING btree (user_id, execution_time);


--
-- Name: submissions_1390575600_1407855600_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1390575600_1407855600_user_id_length_idx ON public.submissions_1390575600_1407855600 USING btree (user_id, length);


--
-- Name: submissions_1390575600_1407855600_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1390575600_1407855600_user_id_point_idx ON public.submissions_1390575600_1407855600 USING btree (user_id, point);


--
-- Name: submissions_1390575600_14078556_execution_time_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1390575600_14078556_execution_time_epoch_second_idx ON public.submissions_1390575600_1407855600 USING btree (execution_time, epoch_second);


--
-- Name: submissions_1407855600_1425135600_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1407855600_1425135600_contest_id_epoch_second_idx ON public.submissions_1407855600_1425135600 USING btree (contest_id, epoch_second);


--
-- Name: submissions_1407855600_1425135600_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1407855600_1425135600_contest_id_execution_time_idx ON public.submissions_1407855600_1425135600 USING btree (contest_id, execution_time);


--
-- Name: submissions_1407855600_1425135600_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1407855600_1425135600_contest_id_length_idx ON public.submissions_1407855600_1425135600 USING btree (contest_id, length);


--
-- Name: submissions_1407855600_1425135600_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1407855600_1425135600_contest_id_point_idx ON public.submissions_1407855600_1425135600 USING btree (contest_id, point);


--
-- Name: submissions_1407855600_1425135600_epoch_second_updated_at_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1407855600_1425135600_epoch_second_updated_at_idx ON public.submissions_1407855600_1425135600 USING btree (epoch_second, updated_at);


--
-- Name: submissions_1407855600_1425135600_execution_time_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1407855600_1425135600_execution_time_length_idx ON public.submissions_1407855600_1425135600 USING btree (execution_time, length);


--
-- Name: submissions_1407855600_1425135600_execution_time_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1407855600_1425135600_execution_time_point_idx ON public.submissions_1407855600_1425135600 USING btree (execution_time, point);


--
-- Name: submissions_1407855600_1425135600_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_1407855600_1425135600_id_epoch_second_idx ON public.submissions_1407855600_1425135600 USING btree (id, epoch_second);


--
-- Name: submissions_1407855600_1425135600_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1407855600_1425135600_language_epoch_second_idx ON public.submissions_1407855600_1425135600 USING btree (language, epoch_second);


--
-- Name: submissions_1407855600_1425135600_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1407855600_1425135600_language_execution_time_idx ON public.submissions_1407855600_1425135600 USING btree (language, execution_time);


--
-- Name: submissions_1407855600_1425135600_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1407855600_1425135600_language_length_idx ON public.submissions_1407855600_1425135600 USING btree (language, length);


--
-- Name: submissions_1407855600_1425135600_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1407855600_1425135600_language_point_idx ON public.submissions_1407855600_1425135600 USING btree (language, point);


--
-- Name: submissions_1407855600_1425135600_length_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1407855600_1425135600_length_epoch_second_idx ON public.submissions_1407855600_1425135600 USING btree (length, epoch_second);


--
-- Name: submissions_1407855600_1425135600_length_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1407855600_1425135600_length_execution_time_idx ON public.submissions_1407855600_1425135600 USING btree (length, execution_time);


--
-- Name: submissions_1407855600_1425135600_length_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1407855600_1425135600_length_point_idx ON public.submissions_1407855600_1425135600 USING btree (length, point);


--
-- Name: submissions_1407855600_1425135600_point_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1407855600_1425135600_point_epoch_second_idx ON public.submissions_1407855600_1425135600 USING btree (point, epoch_second);


--
-- Name: submissions_1407855600_1425135600_point_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1407855600_1425135600_point_execution_time_idx ON public.submissions_1407855600_1425135600 USING btree (point, execution_time);


--
-- Name: submissions_1407855600_1425135600_point_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1407855600_1425135600_point_length_idx ON public.submissions_1407855600_1425135600 USING btree (point, length);


--
-- Name: submissions_1407855600_1425135600_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1407855600_1425135600_problem_id_epoch_second_idx ON public.submissions_1407855600_1425135600 USING btree (problem_id, epoch_second);


--
-- Name: submissions_1407855600_1425135600_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1407855600_1425135600_problem_id_execution_time_idx ON public.submissions_1407855600_1425135600 USING btree (problem_id, execution_time);


--
-- Name: submissions_1407855600_1425135600_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1407855600_1425135600_problem_id_length_idx ON public.submissions_1407855600_1425135600 USING btree (problem_id, length);


--
-- Name: submissions_1407855600_1425135600_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1407855600_1425135600_problem_id_point_idx ON public.submissions_1407855600_1425135600 USING btree (problem_id, point);


--
-- Name: submissions_1407855600_1425135600_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1407855600_1425135600_result_epoch_second_idx ON public.submissions_1407855600_1425135600 USING btree (result, epoch_second);


--
-- Name: submissions_1407855600_1425135600_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1407855600_1425135600_result_execution_time_idx ON public.submissions_1407855600_1425135600 USING btree (result, execution_time);


--
-- Name: submissions_1407855600_1425135600_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1407855600_1425135600_result_length_idx ON public.submissions_1407855600_1425135600 USING btree (result, length);


--
-- Name: submissions_1407855600_1425135600_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1407855600_1425135600_result_point_idx ON public.submissions_1407855600_1425135600 USING btree (result, point);


--
-- Name: submissions_1407855600_1425135600_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1407855600_1425135600_user_id_epoch_second_idx ON public.submissions_1407855600_1425135600 USING btree (user_id, epoch_second);


--
-- Name: submissions_1407855600_1425135600_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1407855600_1425135600_user_id_execution_time_idx ON public.submissions_1407855600_1425135600 USING btree (user_id, execution_time);


--
-- Name: submissions_1407855600_1425135600_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1407855600_1425135600_user_id_length_idx ON public.submissions_1407855600_1425135600 USING btree (user_id, length);


--
-- Name: submissions_1407855600_1425135600_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1407855600_1425135600_user_id_point_idx ON public.submissions_1407855600_1425135600 USING btree (user_id, point);


--
-- Name: submissions_1407855600_14251356_execution_time_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1407855600_14251356_execution_time_epoch_second_idx ON public.submissions_1407855600_1425135600 USING btree (execution_time, epoch_second);


--
-- Name: submissions_1425135600_1442415600_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1425135600_1442415600_contest_id_epoch_second_idx ON public.submissions_1425135600_1442415600 USING btree (contest_id, epoch_second);


--
-- Name: submissions_1425135600_1442415600_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1425135600_1442415600_contest_id_execution_time_idx ON public.submissions_1425135600_1442415600 USING btree (contest_id, execution_time);


--
-- Name: submissions_1425135600_1442415600_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1425135600_1442415600_contest_id_length_idx ON public.submissions_1425135600_1442415600 USING btree (contest_id, length);


--
-- Name: submissions_1425135600_1442415600_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1425135600_1442415600_contest_id_point_idx ON public.submissions_1425135600_1442415600 USING btree (contest_id, point);


--
-- Name: submissions_1425135600_1442415600_epoch_second_updated_at_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1425135600_1442415600_epoch_second_updated_at_idx ON public.submissions_1425135600_1442415600 USING btree (epoch_second, updated_at);


--
-- Name: submissions_1425135600_1442415600_execution_time_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1425135600_1442415600_execution_time_length_idx ON public.submissions_1425135600_1442415600 USING btree (execution_time, length);


--
-- Name: submissions_1425135600_1442415600_execution_time_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1425135600_1442415600_execution_time_point_idx ON public.submissions_1425135600_1442415600 USING btree (execution_time, point);


--
-- Name: submissions_1425135600_1442415600_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_1425135600_1442415600_id_epoch_second_idx ON public.submissions_1425135600_1442415600 USING btree (id, epoch_second);


--
-- Name: submissions_1425135600_1442415600_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1425135600_1442415600_language_epoch_second_idx ON public.submissions_1425135600_1442415600 USING btree (language, epoch_second);


--
-- Name: submissions_1425135600_1442415600_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1425135600_1442415600_language_execution_time_idx ON public.submissions_1425135600_1442415600 USING btree (language, execution_time);


--
-- Name: submissions_1425135600_1442415600_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1425135600_1442415600_language_length_idx ON public.submissions_1425135600_1442415600 USING btree (language, length);


--
-- Name: submissions_1425135600_1442415600_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1425135600_1442415600_language_point_idx ON public.submissions_1425135600_1442415600 USING btree (language, point);


--
-- Name: submissions_1425135600_1442415600_length_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1425135600_1442415600_length_epoch_second_idx ON public.submissions_1425135600_1442415600 USING btree (length, epoch_second);


--
-- Name: submissions_1425135600_1442415600_length_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1425135600_1442415600_length_execution_time_idx ON public.submissions_1425135600_1442415600 USING btree (length, execution_time);


--
-- Name: submissions_1425135600_1442415600_length_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1425135600_1442415600_length_point_idx ON public.submissions_1425135600_1442415600 USING btree (length, point);


--
-- Name: submissions_1425135600_1442415600_point_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1425135600_1442415600_point_epoch_second_idx ON public.submissions_1425135600_1442415600 USING btree (point, epoch_second);


--
-- Name: submissions_1425135600_1442415600_point_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1425135600_1442415600_point_execution_time_idx ON public.submissions_1425135600_1442415600 USING btree (point, execution_time);


--
-- Name: submissions_1425135600_1442415600_point_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1425135600_1442415600_point_length_idx ON public.submissions_1425135600_1442415600 USING btree (point, length);


--
-- Name: submissions_1425135600_1442415600_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1425135600_1442415600_problem_id_epoch_second_idx ON public.submissions_1425135600_1442415600 USING btree (problem_id, epoch_second);


--
-- Name: submissions_1425135600_1442415600_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1425135600_1442415600_problem_id_execution_time_idx ON public.submissions_1425135600_1442415600 USING btree (problem_id, execution_time);


--
-- Name: submissions_1425135600_1442415600_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1425135600_1442415600_problem_id_length_idx ON public.submissions_1425135600_1442415600 USING btree (problem_id, length);


--
-- Name: submissions_1425135600_1442415600_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1425135600_1442415600_problem_id_point_idx ON public.submissions_1425135600_1442415600 USING btree (problem_id, point);


--
-- Name: submissions_1425135600_1442415600_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1425135600_1442415600_result_epoch_second_idx ON public.submissions_1425135600_1442415600 USING btree (result, epoch_second);


--
-- Name: submissions_1425135600_1442415600_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1425135600_1442415600_result_execution_time_idx ON public.submissions_1425135600_1442415600 USING btree (result, execution_time);


--
-- Name: submissions_1425135600_1442415600_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1425135600_1442415600_result_length_idx ON public.submissions_1425135600_1442415600 USING btree (result, length);


--
-- Name: submissions_1425135600_1442415600_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1425135600_1442415600_result_point_idx ON public.submissions_1425135600_1442415600 USING btree (result, point);


--
-- Name: submissions_1425135600_1442415600_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1425135600_1442415600_user_id_epoch_second_idx ON public.submissions_1425135600_1442415600 USING btree (user_id, epoch_second);


--
-- Name: submissions_1425135600_1442415600_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1425135600_1442415600_user_id_execution_time_idx ON public.submissions_1425135600_1442415600 USING btree (user_id, execution_time);


--
-- Name: submissions_1425135600_1442415600_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1425135600_1442415600_user_id_length_idx ON public.submissions_1425135600_1442415600 USING btree (user_id, length);


--
-- Name: submissions_1425135600_1442415600_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1425135600_1442415600_user_id_point_idx ON public.submissions_1425135600_1442415600 USING btree (user_id, point);


--
-- Name: submissions_1425135600_14424156_execution_time_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1425135600_14424156_execution_time_epoch_second_idx ON public.submissions_1425135600_1442415600 USING btree (execution_time, epoch_second);


--
-- Name: submissions_1442415600_1459695600_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1442415600_1459695600_contest_id_epoch_second_idx ON public.submissions_1442415600_1459695600 USING btree (contest_id, epoch_second);


--
-- Name: submissions_1442415600_1459695600_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1442415600_1459695600_contest_id_execution_time_idx ON public.submissions_1442415600_1459695600 USING btree (contest_id, execution_time);


--
-- Name: submissions_1442415600_1459695600_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1442415600_1459695600_contest_id_length_idx ON public.submissions_1442415600_1459695600 USING btree (contest_id, length);


--
-- Name: submissions_1442415600_1459695600_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1442415600_1459695600_contest_id_point_idx ON public.submissions_1442415600_1459695600 USING btree (contest_id, point);


--
-- Name: submissions_1442415600_1459695600_epoch_second_updated_at_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1442415600_1459695600_epoch_second_updated_at_idx ON public.submissions_1442415600_1459695600 USING btree (epoch_second, updated_at);


--
-- Name: submissions_1442415600_1459695600_execution_time_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1442415600_1459695600_execution_time_length_idx ON public.submissions_1442415600_1459695600 USING btree (execution_time, length);


--
-- Name: submissions_1442415600_1459695600_execution_time_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1442415600_1459695600_execution_time_point_idx ON public.submissions_1442415600_1459695600 USING btree (execution_time, point);


--
-- Name: submissions_1442415600_1459695600_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_1442415600_1459695600_id_epoch_second_idx ON public.submissions_1442415600_1459695600 USING btree (id, epoch_second);


--
-- Name: submissions_1442415600_1459695600_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1442415600_1459695600_language_epoch_second_idx ON public.submissions_1442415600_1459695600 USING btree (language, epoch_second);


--
-- Name: submissions_1442415600_1459695600_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1442415600_1459695600_language_execution_time_idx ON public.submissions_1442415600_1459695600 USING btree (language, execution_time);


--
-- Name: submissions_1442415600_1459695600_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1442415600_1459695600_language_length_idx ON public.submissions_1442415600_1459695600 USING btree (language, length);


--
-- Name: submissions_1442415600_1459695600_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1442415600_1459695600_language_point_idx ON public.submissions_1442415600_1459695600 USING btree (language, point);


--
-- Name: submissions_1442415600_1459695600_length_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1442415600_1459695600_length_epoch_second_idx ON public.submissions_1442415600_1459695600 USING btree (length, epoch_second);


--
-- Name: submissions_1442415600_1459695600_length_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1442415600_1459695600_length_execution_time_idx ON public.submissions_1442415600_1459695600 USING btree (length, execution_time);


--
-- Name: submissions_1442415600_1459695600_length_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1442415600_1459695600_length_point_idx ON public.submissions_1442415600_1459695600 USING btree (length, point);


--
-- Name: submissions_1442415600_1459695600_point_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1442415600_1459695600_point_epoch_second_idx ON public.submissions_1442415600_1459695600 USING btree (point, epoch_second);


--
-- Name: submissions_1442415600_1459695600_point_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1442415600_1459695600_point_execution_time_idx ON public.submissions_1442415600_1459695600 USING btree (point, execution_time);


--
-- Name: submissions_1442415600_1459695600_point_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1442415600_1459695600_point_length_idx ON public.submissions_1442415600_1459695600 USING btree (point, length);


--
-- Name: submissions_1442415600_1459695600_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1442415600_1459695600_problem_id_epoch_second_idx ON public.submissions_1442415600_1459695600 USING btree (problem_id, epoch_second);


--
-- Name: submissions_1442415600_1459695600_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1442415600_1459695600_problem_id_execution_time_idx ON public.submissions_1442415600_1459695600 USING btree (problem_id, execution_time);


--
-- Name: submissions_1442415600_1459695600_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1442415600_1459695600_problem_id_length_idx ON public.submissions_1442415600_1459695600 USING btree (problem_id, length);


--
-- Name: submissions_1442415600_1459695600_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1442415600_1459695600_problem_id_point_idx ON public.submissions_1442415600_1459695600 USING btree (problem_id, point);


--
-- Name: submissions_1442415600_1459695600_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1442415600_1459695600_result_epoch_second_idx ON public.submissions_1442415600_1459695600 USING btree (result, epoch_second);


--
-- Name: submissions_1442415600_1459695600_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1442415600_1459695600_result_execution_time_idx ON public.submissions_1442415600_1459695600 USING btree (result, execution_time);


--
-- Name: submissions_1442415600_1459695600_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1442415600_1459695600_result_length_idx ON public.submissions_1442415600_1459695600 USING btree (result, length);


--
-- Name: submissions_1442415600_1459695600_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1442415600_1459695600_result_point_idx ON public.submissions_1442415600_1459695600 USING btree (result, point);


--
-- Name: submissions_1442415600_1459695600_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1442415600_1459695600_user_id_epoch_second_idx ON public.submissions_1442415600_1459695600 USING btree (user_id, epoch_second);


--
-- Name: submissions_1442415600_1459695600_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1442415600_1459695600_user_id_execution_time_idx ON public.submissions_1442415600_1459695600 USING btree (user_id, execution_time);


--
-- Name: submissions_1442415600_1459695600_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1442415600_1459695600_user_id_length_idx ON public.submissions_1442415600_1459695600 USING btree (user_id, length);


--
-- Name: submissions_1442415600_1459695600_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1442415600_1459695600_user_id_point_idx ON public.submissions_1442415600_1459695600 USING btree (user_id, point);


--
-- Name: submissions_1442415600_14596956_execution_time_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1442415600_14596956_execution_time_epoch_second_idx ON public.submissions_1442415600_1459695600 USING btree (execution_time, epoch_second);


--
-- Name: submissions_1459695600_1476975600_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1459695600_1476975600_contest_id_epoch_second_idx ON public.submissions_1459695600_1476975600 USING btree (contest_id, epoch_second);


--
-- Name: submissions_1459695600_1476975600_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1459695600_1476975600_contest_id_execution_time_idx ON public.submissions_1459695600_1476975600 USING btree (contest_id, execution_time);


--
-- Name: submissions_1459695600_1476975600_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1459695600_1476975600_contest_id_length_idx ON public.submissions_1459695600_1476975600 USING btree (contest_id, length);


--
-- Name: submissions_1459695600_1476975600_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1459695600_1476975600_contest_id_point_idx ON public.submissions_1459695600_1476975600 USING btree (contest_id, point);


--
-- Name: submissions_1459695600_1476975600_epoch_second_updated_at_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1459695600_1476975600_epoch_second_updated_at_idx ON public.submissions_1459695600_1476975600 USING btree (epoch_second, updated_at);


--
-- Name: submissions_1459695600_1476975600_execution_time_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1459695600_1476975600_execution_time_length_idx ON public.submissions_1459695600_1476975600 USING btree (execution_time, length);


--
-- Name: submissions_1459695600_1476975600_execution_time_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1459695600_1476975600_execution_time_point_idx ON public.submissions_1459695600_1476975600 USING btree (execution_time, point);


--
-- Name: submissions_1459695600_1476975600_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_1459695600_1476975600_id_epoch_second_idx ON public.submissions_1459695600_1476975600 USING btree (id, epoch_second);


--
-- Name: submissions_1459695600_1476975600_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1459695600_1476975600_language_epoch_second_idx ON public.submissions_1459695600_1476975600 USING btree (language, epoch_second);


--
-- Name: submissions_1459695600_1476975600_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1459695600_1476975600_language_execution_time_idx ON public.submissions_1459695600_1476975600 USING btree (language, execution_time);


--
-- Name: submissions_1459695600_1476975600_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1459695600_1476975600_language_length_idx ON public.submissions_1459695600_1476975600 USING btree (language, length);


--
-- Name: submissions_1459695600_1476975600_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1459695600_1476975600_language_point_idx ON public.submissions_1459695600_1476975600 USING btree (language, point);


--
-- Name: submissions_1459695600_1476975600_length_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1459695600_1476975600_length_epoch_second_idx ON public.submissions_1459695600_1476975600 USING btree (length, epoch_second);


--
-- Name: submissions_1459695600_1476975600_length_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1459695600_1476975600_length_execution_time_idx ON public.submissions_1459695600_1476975600 USING btree (length, execution_time);


--
-- Name: submissions_1459695600_1476975600_length_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1459695600_1476975600_length_point_idx ON public.submissions_1459695600_1476975600 USING btree (length, point);


--
-- Name: submissions_1459695600_1476975600_point_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1459695600_1476975600_point_epoch_second_idx ON public.submissions_1459695600_1476975600 USING btree (point, epoch_second);


--
-- Name: submissions_1459695600_1476975600_point_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1459695600_1476975600_point_execution_time_idx ON public.submissions_1459695600_1476975600 USING btree (point, execution_time);


--
-- Name: submissions_1459695600_1476975600_point_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1459695600_1476975600_point_length_idx ON public.submissions_1459695600_1476975600 USING btree (point, length);


--
-- Name: submissions_1459695600_1476975600_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1459695600_1476975600_problem_id_epoch_second_idx ON public.submissions_1459695600_1476975600 USING btree (problem_id, epoch_second);


--
-- Name: submissions_1459695600_1476975600_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1459695600_1476975600_problem_id_execution_time_idx ON public.submissions_1459695600_1476975600 USING btree (problem_id, execution_time);


--
-- Name: submissions_1459695600_1476975600_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1459695600_1476975600_problem_id_length_idx ON public.submissions_1459695600_1476975600 USING btree (problem_id, length);


--
-- Name: submissions_1459695600_1476975600_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1459695600_1476975600_problem_id_point_idx ON public.submissions_1459695600_1476975600 USING btree (problem_id, point);


--
-- Name: submissions_1459695600_1476975600_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1459695600_1476975600_result_epoch_second_idx ON public.submissions_1459695600_1476975600 USING btree (result, epoch_second);


--
-- Name: submissions_1459695600_1476975600_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1459695600_1476975600_result_execution_time_idx ON public.submissions_1459695600_1476975600 USING btree (result, execution_time);


--
-- Name: submissions_1459695600_1476975600_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1459695600_1476975600_result_length_idx ON public.submissions_1459695600_1476975600 USING btree (result, length);


--
-- Name: submissions_1459695600_1476975600_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1459695600_1476975600_result_point_idx ON public.submissions_1459695600_1476975600 USING btree (result, point);


--
-- Name: submissions_1459695600_1476975600_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1459695600_1476975600_user_id_epoch_second_idx ON public.submissions_1459695600_1476975600 USING btree (user_id, epoch_second);


--
-- Name: submissions_1459695600_1476975600_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1459695600_1476975600_user_id_execution_time_idx ON public.submissions_1459695600_1476975600 USING btree (user_id, execution_time);


--
-- Name: submissions_1459695600_1476975600_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1459695600_1476975600_user_id_length_idx ON public.submissions_1459695600_1476975600 USING btree (user_id, length);


--
-- Name: submissions_1459695600_1476975600_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1459695600_1476975600_user_id_point_idx ON public.submissions_1459695600_1476975600 USING btree (user_id, point);


--
-- Name: submissions_1459695600_14769756_execution_time_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1459695600_14769756_execution_time_epoch_second_idx ON public.submissions_1459695600_1476975600 USING btree (execution_time, epoch_second);


--
-- Name: submissions_1476975600_1494255600_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1476975600_1494255600_contest_id_epoch_second_idx ON public.submissions_1476975600_1494255600 USING btree (contest_id, epoch_second);


--
-- Name: submissions_1476975600_1494255600_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1476975600_1494255600_contest_id_execution_time_idx ON public.submissions_1476975600_1494255600 USING btree (contest_id, execution_time);


--
-- Name: submissions_1476975600_1494255600_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1476975600_1494255600_contest_id_length_idx ON public.submissions_1476975600_1494255600 USING btree (contest_id, length);


--
-- Name: submissions_1476975600_1494255600_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1476975600_1494255600_contest_id_point_idx ON public.submissions_1476975600_1494255600 USING btree (contest_id, point);


--
-- Name: submissions_1476975600_1494255600_epoch_second_updated_at_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1476975600_1494255600_epoch_second_updated_at_idx ON public.submissions_1476975600_1494255600 USING btree (epoch_second, updated_at);


--
-- Name: submissions_1476975600_1494255600_execution_time_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1476975600_1494255600_execution_time_length_idx ON public.submissions_1476975600_1494255600 USING btree (execution_time, length);


--
-- Name: submissions_1476975600_1494255600_execution_time_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1476975600_1494255600_execution_time_point_idx ON public.submissions_1476975600_1494255600 USING btree (execution_time, point);


--
-- Name: submissions_1476975600_1494255600_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_1476975600_1494255600_id_epoch_second_idx ON public.submissions_1476975600_1494255600 USING btree (id, epoch_second);


--
-- Name: submissions_1476975600_1494255600_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1476975600_1494255600_language_epoch_second_idx ON public.submissions_1476975600_1494255600 USING btree (language, epoch_second);


--
-- Name: submissions_1476975600_1494255600_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1476975600_1494255600_language_execution_time_idx ON public.submissions_1476975600_1494255600 USING btree (language, execution_time);


--
-- Name: submissions_1476975600_1494255600_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1476975600_1494255600_language_length_idx ON public.submissions_1476975600_1494255600 USING btree (language, length);


--
-- Name: submissions_1476975600_1494255600_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1476975600_1494255600_language_point_idx ON public.submissions_1476975600_1494255600 USING btree (language, point);


--
-- Name: submissions_1476975600_1494255600_length_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1476975600_1494255600_length_epoch_second_idx ON public.submissions_1476975600_1494255600 USING btree (length, epoch_second);


--
-- Name: submissions_1476975600_1494255600_length_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1476975600_1494255600_length_execution_time_idx ON public.submissions_1476975600_1494255600 USING btree (length, execution_time);


--
-- Name: submissions_1476975600_1494255600_length_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1476975600_1494255600_length_point_idx ON public.submissions_1476975600_1494255600 USING btree (length, point);


--
-- Name: submissions_1476975600_1494255600_point_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1476975600_1494255600_point_epoch_second_idx ON public.submissions_1476975600_1494255600 USING btree (point, epoch_second);


--
-- Name: submissions_1476975600_1494255600_point_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1476975600_1494255600_point_execution_time_idx ON public.submissions_1476975600_1494255600 USING btree (point, execution_time);


--
-- Name: submissions_1476975600_1494255600_point_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1476975600_1494255600_point_length_idx ON public.submissions_1476975600_1494255600 USING btree (point, length);


--
-- Name: submissions_1476975600_1494255600_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1476975600_1494255600_problem_id_epoch_second_idx ON public.submissions_1476975600_1494255600 USING btree (problem_id, epoch_second);


--
-- Name: submissions_1476975600_1494255600_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1476975600_1494255600_problem_id_execution_time_idx ON public.submissions_1476975600_1494255600 USING btree (problem_id, execution_time);


--
-- Name: submissions_1476975600_1494255600_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1476975600_1494255600_problem_id_length_idx ON public.submissions_1476975600_1494255600 USING btree (problem_id, length);


--
-- Name: submissions_1476975600_1494255600_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1476975600_1494255600_problem_id_point_idx ON public.submissions_1476975600_1494255600 USING btree (problem_id, point);


--
-- Name: submissions_1476975600_1494255600_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1476975600_1494255600_result_epoch_second_idx ON public.submissions_1476975600_1494255600 USING btree (result, epoch_second);


--
-- Name: submissions_1476975600_1494255600_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1476975600_1494255600_result_execution_time_idx ON public.submissions_1476975600_1494255600 USING btree (result, execution_time);


--
-- Name: submissions_1476975600_1494255600_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1476975600_1494255600_result_length_idx ON public.submissions_1476975600_1494255600 USING btree (result, length);


--
-- Name: submissions_1476975600_1494255600_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1476975600_1494255600_result_point_idx ON public.submissions_1476975600_1494255600 USING btree (result, point);


--
-- Name: submissions_1476975600_1494255600_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1476975600_1494255600_user_id_epoch_second_idx ON public.submissions_1476975600_1494255600 USING btree (user_id, epoch_second);


--
-- Name: submissions_1476975600_1494255600_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1476975600_1494255600_user_id_execution_time_idx ON public.submissions_1476975600_1494255600 USING btree (user_id, execution_time);


--
-- Name: submissions_1476975600_1494255600_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1476975600_1494255600_user_id_length_idx ON public.submissions_1476975600_1494255600 USING btree (user_id, length);


--
-- Name: submissions_1476975600_1494255600_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1476975600_1494255600_user_id_point_idx ON public.submissions_1476975600_1494255600 USING btree (user_id, point);


--
-- Name: submissions_1476975600_14942556_execution_time_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1476975600_14942556_execution_time_epoch_second_idx ON public.submissions_1476975600_1494255600 USING btree (execution_time, epoch_second);


--
-- Name: submissions_1494255600_1511535600_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1494255600_1511535600_contest_id_epoch_second_idx ON public.submissions_1494255600_1511535600 USING btree (contest_id, epoch_second);


--
-- Name: submissions_1494255600_1511535600_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1494255600_1511535600_contest_id_execution_time_idx ON public.submissions_1494255600_1511535600 USING btree (contest_id, execution_time);


--
-- Name: submissions_1494255600_1511535600_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1494255600_1511535600_contest_id_length_idx ON public.submissions_1494255600_1511535600 USING btree (contest_id, length);


--
-- Name: submissions_1494255600_1511535600_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1494255600_1511535600_contest_id_point_idx ON public.submissions_1494255600_1511535600 USING btree (contest_id, point);


--
-- Name: submissions_1494255600_1511535600_epoch_second_updated_at_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1494255600_1511535600_epoch_second_updated_at_idx ON public.submissions_1494255600_1511535600 USING btree (epoch_second, updated_at);


--
-- Name: submissions_1494255600_1511535600_execution_time_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1494255600_1511535600_execution_time_length_idx ON public.submissions_1494255600_1511535600 USING btree (execution_time, length);


--
-- Name: submissions_1494255600_1511535600_execution_time_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1494255600_1511535600_execution_time_point_idx ON public.submissions_1494255600_1511535600 USING btree (execution_time, point);


--
-- Name: submissions_1494255600_1511535600_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_1494255600_1511535600_id_epoch_second_idx ON public.submissions_1494255600_1511535600 USING btree (id, epoch_second);


--
-- Name: submissions_1494255600_1511535600_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1494255600_1511535600_language_epoch_second_idx ON public.submissions_1494255600_1511535600 USING btree (language, epoch_second);


--
-- Name: submissions_1494255600_1511535600_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1494255600_1511535600_language_execution_time_idx ON public.submissions_1494255600_1511535600 USING btree (language, execution_time);


--
-- Name: submissions_1494255600_1511535600_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1494255600_1511535600_language_length_idx ON public.submissions_1494255600_1511535600 USING btree (language, length);


--
-- Name: submissions_1494255600_1511535600_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1494255600_1511535600_language_point_idx ON public.submissions_1494255600_1511535600 USING btree (language, point);


--
-- Name: submissions_1494255600_1511535600_length_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1494255600_1511535600_length_epoch_second_idx ON public.submissions_1494255600_1511535600 USING btree (length, epoch_second);


--
-- Name: submissions_1494255600_1511535600_length_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1494255600_1511535600_length_execution_time_idx ON public.submissions_1494255600_1511535600 USING btree (length, execution_time);


--
-- Name: submissions_1494255600_1511535600_length_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1494255600_1511535600_length_point_idx ON public.submissions_1494255600_1511535600 USING btree (length, point);


--
-- Name: submissions_1494255600_1511535600_point_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1494255600_1511535600_point_epoch_second_idx ON public.submissions_1494255600_1511535600 USING btree (point, epoch_second);


--
-- Name: submissions_1494255600_1511535600_point_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1494255600_1511535600_point_execution_time_idx ON public.submissions_1494255600_1511535600 USING btree (point, execution_time);


--
-- Name: submissions_1494255600_1511535600_point_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1494255600_1511535600_point_length_idx ON public.submissions_1494255600_1511535600 USING btree (point, length);


--
-- Name: submissions_1494255600_1511535600_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1494255600_1511535600_problem_id_epoch_second_idx ON public.submissions_1494255600_1511535600 USING btree (problem_id, epoch_second);


--
-- Name: submissions_1494255600_1511535600_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1494255600_1511535600_problem_id_execution_time_idx ON public.submissions_1494255600_1511535600 USING btree (problem_id, execution_time);


--
-- Name: submissions_1494255600_1511535600_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1494255600_1511535600_problem_id_length_idx ON public.submissions_1494255600_1511535600 USING btree (problem_id, length);


--
-- Name: submissions_1494255600_1511535600_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1494255600_1511535600_problem_id_point_idx ON public.submissions_1494255600_1511535600 USING btree (problem_id, point);


--
-- Name: submissions_1494255600_1511535600_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1494255600_1511535600_result_epoch_second_idx ON public.submissions_1494255600_1511535600 USING btree (result, epoch_second);


--
-- Name: submissions_1494255600_1511535600_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1494255600_1511535600_result_execution_time_idx ON public.submissions_1494255600_1511535600 USING btree (result, execution_time);


--
-- Name: submissions_1494255600_1511535600_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1494255600_1511535600_result_length_idx ON public.submissions_1494255600_1511535600 USING btree (result, length);


--
-- Name: submissions_1494255600_1511535600_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1494255600_1511535600_result_point_idx ON public.submissions_1494255600_1511535600 USING btree (result, point);


--
-- Name: submissions_1494255600_1511535600_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1494255600_1511535600_user_id_epoch_second_idx ON public.submissions_1494255600_1511535600 USING btree (user_id, epoch_second);


--
-- Name: submissions_1494255600_1511535600_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1494255600_1511535600_user_id_execution_time_idx ON public.submissions_1494255600_1511535600 USING btree (user_id, execution_time);


--
-- Name: submissions_1494255600_1511535600_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1494255600_1511535600_user_id_length_idx ON public.submissions_1494255600_1511535600 USING btree (user_id, length);


--
-- Name: submissions_1494255600_1511535600_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1494255600_1511535600_user_id_point_idx ON public.submissions_1494255600_1511535600 USING btree (user_id, point);


--
-- Name: submissions_1494255600_15115356_execution_time_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1494255600_15115356_execution_time_epoch_second_idx ON public.submissions_1494255600_1511535600 USING btree (execution_time, epoch_second);


--
-- Name: submissions_1511535600_1528815600_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1511535600_1528815600_contest_id_epoch_second_idx ON public.submissions_1511535600_1528815600 USING btree (contest_id, epoch_second);


--
-- Name: submissions_1511535600_1528815600_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1511535600_1528815600_contest_id_execution_time_idx ON public.submissions_1511535600_1528815600 USING btree (contest_id, execution_time);


--
-- Name: submissions_1511535600_1528815600_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1511535600_1528815600_contest_id_length_idx ON public.submissions_1511535600_1528815600 USING btree (contest_id, length);


--
-- Name: submissions_1511535600_1528815600_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1511535600_1528815600_contest_id_point_idx ON public.submissions_1511535600_1528815600 USING btree (contest_id, point);


--
-- Name: submissions_1511535600_1528815600_epoch_second_updated_at_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1511535600_1528815600_epoch_second_updated_at_idx ON public.submissions_1511535600_1528815600 USING btree (epoch_second, updated_at);


--
-- Name: submissions_1511535600_1528815600_execution_time_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1511535600_1528815600_execution_time_length_idx ON public.submissions_1511535600_1528815600 USING btree (execution_time, length);


--
-- Name: submissions_1511535600_1528815600_execution_time_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1511535600_1528815600_execution_time_point_idx ON public.submissions_1511535600_1528815600 USING btree (execution_time, point);


--
-- Name: submissions_1511535600_1528815600_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_1511535600_1528815600_id_epoch_second_idx ON public.submissions_1511535600_1528815600 USING btree (id, epoch_second);


--
-- Name: submissions_1511535600_1528815600_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1511535600_1528815600_language_epoch_second_idx ON public.submissions_1511535600_1528815600 USING btree (language, epoch_second);


--
-- Name: submissions_1511535600_1528815600_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1511535600_1528815600_language_execution_time_idx ON public.submissions_1511535600_1528815600 USING btree (language, execution_time);


--
-- Name: submissions_1511535600_1528815600_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1511535600_1528815600_language_length_idx ON public.submissions_1511535600_1528815600 USING btree (language, length);


--
-- Name: submissions_1511535600_1528815600_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1511535600_1528815600_language_point_idx ON public.submissions_1511535600_1528815600 USING btree (language, point);


--
-- Name: submissions_1511535600_1528815600_length_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1511535600_1528815600_length_epoch_second_idx ON public.submissions_1511535600_1528815600 USING btree (length, epoch_second);


--
-- Name: submissions_1511535600_1528815600_length_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1511535600_1528815600_length_execution_time_idx ON public.submissions_1511535600_1528815600 USING btree (length, execution_time);


--
-- Name: submissions_1511535600_1528815600_length_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1511535600_1528815600_length_point_idx ON public.submissions_1511535600_1528815600 USING btree (length, point);


--
-- Name: submissions_1511535600_1528815600_point_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1511535600_1528815600_point_epoch_second_idx ON public.submissions_1511535600_1528815600 USING btree (point, epoch_second);


--
-- Name: submissions_1511535600_1528815600_point_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1511535600_1528815600_point_execution_time_idx ON public.submissions_1511535600_1528815600 USING btree (point, execution_time);


--
-- Name: submissions_1511535600_1528815600_point_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1511535600_1528815600_point_length_idx ON public.submissions_1511535600_1528815600 USING btree (point, length);


--
-- Name: submissions_1511535600_1528815600_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1511535600_1528815600_problem_id_epoch_second_idx ON public.submissions_1511535600_1528815600 USING btree (problem_id, epoch_second);


--
-- Name: submissions_1511535600_1528815600_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1511535600_1528815600_problem_id_execution_time_idx ON public.submissions_1511535600_1528815600 USING btree (problem_id, execution_time);


--
-- Name: submissions_1511535600_1528815600_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1511535600_1528815600_problem_id_length_idx ON public.submissions_1511535600_1528815600 USING btree (problem_id, length);


--
-- Name: submissions_1511535600_1528815600_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1511535600_1528815600_problem_id_point_idx ON public.submissions_1511535600_1528815600 USING btree (problem_id, point);


--
-- Name: submissions_1511535600_1528815600_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1511535600_1528815600_result_epoch_second_idx ON public.submissions_1511535600_1528815600 USING btree (result, epoch_second);


--
-- Name: submissions_1511535600_1528815600_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1511535600_1528815600_result_execution_time_idx ON public.submissions_1511535600_1528815600 USING btree (result, execution_time);


--
-- Name: submissions_1511535600_1528815600_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1511535600_1528815600_result_length_idx ON public.submissions_1511535600_1528815600 USING btree (result, length);


--
-- Name: submissions_1511535600_1528815600_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1511535600_1528815600_result_point_idx ON public.submissions_1511535600_1528815600 USING btree (result, point);


--
-- Name: submissions_1511535600_1528815600_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1511535600_1528815600_user_id_epoch_second_idx ON public.submissions_1511535600_1528815600 USING btree (user_id, epoch_second);


--
-- Name: submissions_1511535600_1528815600_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1511535600_1528815600_user_id_execution_time_idx ON public.submissions_1511535600_1528815600 USING btree (user_id, execution_time);


--
-- Name: submissions_1511535600_1528815600_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1511535600_1528815600_user_id_length_idx ON public.submissions_1511535600_1528815600 USING btree (user_id, length);


--
-- Name: submissions_1511535600_1528815600_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1511535600_1528815600_user_id_point_idx ON public.submissions_1511535600_1528815600 USING btree (user_id, point);


--
-- Name: submissions_1511535600_15288156_execution_time_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1511535600_15288156_execution_time_epoch_second_idx ON public.submissions_1511535600_1528815600 USING btree (execution_time, epoch_second);


--
-- Name: submissions_1528815600_1546095600_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1528815600_1546095600_contest_id_epoch_second_idx ON public.submissions_1528815600_1546095600 USING btree (contest_id, epoch_second);


--
-- Name: submissions_1528815600_1546095600_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1528815600_1546095600_contest_id_execution_time_idx ON public.submissions_1528815600_1546095600 USING btree (contest_id, execution_time);


--
-- Name: submissions_1528815600_1546095600_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1528815600_1546095600_contest_id_length_idx ON public.submissions_1528815600_1546095600 USING btree (contest_id, length);


--
-- Name: submissions_1528815600_1546095600_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1528815600_1546095600_contest_id_point_idx ON public.submissions_1528815600_1546095600 USING btree (contest_id, point);


--
-- Name: submissions_1528815600_1546095600_epoch_second_updated_at_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1528815600_1546095600_epoch_second_updated_at_idx ON public.submissions_1528815600_1546095600 USING btree (epoch_second, updated_at);


--
-- Name: submissions_1528815600_1546095600_execution_time_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1528815600_1546095600_execution_time_length_idx ON public.submissions_1528815600_1546095600 USING btree (execution_time, length);


--
-- Name: submissions_1528815600_1546095600_execution_time_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1528815600_1546095600_execution_time_point_idx ON public.submissions_1528815600_1546095600 USING btree (execution_time, point);


--
-- Name: submissions_1528815600_1546095600_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_1528815600_1546095600_id_epoch_second_idx ON public.submissions_1528815600_1546095600 USING btree (id, epoch_second);


--
-- Name: submissions_1528815600_1546095600_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1528815600_1546095600_language_epoch_second_idx ON public.submissions_1528815600_1546095600 USING btree (language, epoch_second);


--
-- Name: submissions_1528815600_1546095600_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1528815600_1546095600_language_execution_time_idx ON public.submissions_1528815600_1546095600 USING btree (language, execution_time);


--
-- Name: submissions_1528815600_1546095600_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1528815600_1546095600_language_length_idx ON public.submissions_1528815600_1546095600 USING btree (language, length);


--
-- Name: submissions_1528815600_1546095600_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1528815600_1546095600_language_point_idx ON public.submissions_1528815600_1546095600 USING btree (language, point);


--
-- Name: submissions_1528815600_1546095600_length_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1528815600_1546095600_length_epoch_second_idx ON public.submissions_1528815600_1546095600 USING btree (length, epoch_second);


--
-- Name: submissions_1528815600_1546095600_length_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1528815600_1546095600_length_execution_time_idx ON public.submissions_1528815600_1546095600 USING btree (length, execution_time);


--
-- Name: submissions_1528815600_1546095600_length_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1528815600_1546095600_length_point_idx ON public.submissions_1528815600_1546095600 USING btree (length, point);


--
-- Name: submissions_1528815600_1546095600_point_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1528815600_1546095600_point_epoch_second_idx ON public.submissions_1528815600_1546095600 USING btree (point, epoch_second);


--
-- Name: submissions_1528815600_1546095600_point_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1528815600_1546095600_point_execution_time_idx ON public.submissions_1528815600_1546095600 USING btree (point, execution_time);


--
-- Name: submissions_1528815600_1546095600_point_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1528815600_1546095600_point_length_idx ON public.submissions_1528815600_1546095600 USING btree (point, length);


--
-- Name: submissions_1528815600_1546095600_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1528815600_1546095600_problem_id_epoch_second_idx ON public.submissions_1528815600_1546095600 USING btree (problem_id, epoch_second);


--
-- Name: submissions_1528815600_1546095600_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1528815600_1546095600_problem_id_execution_time_idx ON public.submissions_1528815600_1546095600 USING btree (problem_id, execution_time);


--
-- Name: submissions_1528815600_1546095600_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1528815600_1546095600_problem_id_length_idx ON public.submissions_1528815600_1546095600 USING btree (problem_id, length);


--
-- Name: submissions_1528815600_1546095600_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1528815600_1546095600_problem_id_point_idx ON public.submissions_1528815600_1546095600 USING btree (problem_id, point);


--
-- Name: submissions_1528815600_1546095600_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1528815600_1546095600_result_epoch_second_idx ON public.submissions_1528815600_1546095600 USING btree (result, epoch_second);


--
-- Name: submissions_1528815600_1546095600_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1528815600_1546095600_result_execution_time_idx ON public.submissions_1528815600_1546095600 USING btree (result, execution_time);


--
-- Name: submissions_1528815600_1546095600_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1528815600_1546095600_result_length_idx ON public.submissions_1528815600_1546095600 USING btree (result, length);


--
-- Name: submissions_1528815600_1546095600_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1528815600_1546095600_result_point_idx ON public.submissions_1528815600_1546095600 USING btree (result, point);


--
-- Name: submissions_1528815600_1546095600_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1528815600_1546095600_user_id_epoch_second_idx ON public.submissions_1528815600_1546095600 USING btree (user_id, epoch_second);


--
-- Name: submissions_1528815600_1546095600_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1528815600_1546095600_user_id_execution_time_idx ON public.submissions_1528815600_1546095600 USING btree (user_id, execution_time);


--
-- Name: submissions_1528815600_1546095600_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1528815600_1546095600_user_id_length_idx ON public.submissions_1528815600_1546095600 USING btree (user_id, length);


--
-- Name: submissions_1528815600_1546095600_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1528815600_1546095600_user_id_point_idx ON public.submissions_1528815600_1546095600 USING btree (user_id, point);


--
-- Name: submissions_1528815600_15460956_execution_time_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1528815600_15460956_execution_time_epoch_second_idx ON public.submissions_1528815600_1546095600 USING btree (execution_time, epoch_second);


--
-- Name: submissions_1546095600_1563375600_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1546095600_1563375600_contest_id_epoch_second_idx ON public.submissions_1546095600_1563375600 USING btree (contest_id, epoch_second);


--
-- Name: submissions_1546095600_1563375600_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1546095600_1563375600_contest_id_execution_time_idx ON public.submissions_1546095600_1563375600 USING btree (contest_id, execution_time);


--
-- Name: submissions_1546095600_1563375600_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1546095600_1563375600_contest_id_length_idx ON public.submissions_1546095600_1563375600 USING btree (contest_id, length);


--
-- Name: submissions_1546095600_1563375600_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1546095600_1563375600_contest_id_point_idx ON public.submissions_1546095600_1563375600 USING btree (contest_id, point);


--
-- Name: submissions_1546095600_1563375600_epoch_second_updated_at_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1546095600_1563375600_epoch_second_updated_at_idx ON public.submissions_1546095600_1563375600 USING btree (epoch_second, updated_at);


--
-- Name: submissions_1546095600_1563375600_execution_time_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1546095600_1563375600_execution_time_length_idx ON public.submissions_1546095600_1563375600 USING btree (execution_time, length);


--
-- Name: submissions_1546095600_1563375600_execution_time_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1546095600_1563375600_execution_time_point_idx ON public.submissions_1546095600_1563375600 USING btree (execution_time, point);


--
-- Name: submissions_1546095600_1563375600_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_1546095600_1563375600_id_epoch_second_idx ON public.submissions_1546095600_1563375600 USING btree (id, epoch_second);


--
-- Name: submissions_1546095600_1563375600_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1546095600_1563375600_language_epoch_second_idx ON public.submissions_1546095600_1563375600 USING btree (language, epoch_second);


--
-- Name: submissions_1546095600_1563375600_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1546095600_1563375600_language_execution_time_idx ON public.submissions_1546095600_1563375600 USING btree (language, execution_time);


--
-- Name: submissions_1546095600_1563375600_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1546095600_1563375600_language_length_idx ON public.submissions_1546095600_1563375600 USING btree (language, length);


--
-- Name: submissions_1546095600_1563375600_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1546095600_1563375600_language_point_idx ON public.submissions_1546095600_1563375600 USING btree (language, point);


--
-- Name: submissions_1546095600_1563375600_length_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1546095600_1563375600_length_epoch_second_idx ON public.submissions_1546095600_1563375600 USING btree (length, epoch_second);


--
-- Name: submissions_1546095600_1563375600_length_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1546095600_1563375600_length_execution_time_idx ON public.submissions_1546095600_1563375600 USING btree (length, execution_time);


--
-- Name: submissions_1546095600_1563375600_length_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1546095600_1563375600_length_point_idx ON public.submissions_1546095600_1563375600 USING btree (length, point);


--
-- Name: submissions_1546095600_1563375600_point_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1546095600_1563375600_point_epoch_second_idx ON public.submissions_1546095600_1563375600 USING btree (point, epoch_second);


--
-- Name: submissions_1546095600_1563375600_point_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1546095600_1563375600_point_execution_time_idx ON public.submissions_1546095600_1563375600 USING btree (point, execution_time);


--
-- Name: submissions_1546095600_1563375600_point_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1546095600_1563375600_point_length_idx ON public.submissions_1546095600_1563375600 USING btree (point, length);


--
-- Name: submissions_1546095600_1563375600_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1546095600_1563375600_problem_id_epoch_second_idx ON public.submissions_1546095600_1563375600 USING btree (problem_id, epoch_second);


--
-- Name: submissions_1546095600_1563375600_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1546095600_1563375600_problem_id_execution_time_idx ON public.submissions_1546095600_1563375600 USING btree (problem_id, execution_time);


--
-- Name: submissions_1546095600_1563375600_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1546095600_1563375600_problem_id_length_idx ON public.submissions_1546095600_1563375600 USING btree (problem_id, length);


--
-- Name: submissions_1546095600_1563375600_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1546095600_1563375600_problem_id_point_idx ON public.submissions_1546095600_1563375600 USING btree (problem_id, point);


--
-- Name: submissions_1546095600_1563375600_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1546095600_1563375600_result_epoch_second_idx ON public.submissions_1546095600_1563375600 USING btree (result, epoch_second);


--
-- Name: submissions_1546095600_1563375600_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1546095600_1563375600_result_execution_time_idx ON public.submissions_1546095600_1563375600 USING btree (result, execution_time);


--
-- Name: submissions_1546095600_1563375600_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1546095600_1563375600_result_length_idx ON public.submissions_1546095600_1563375600 USING btree (result, length);


--
-- Name: submissions_1546095600_1563375600_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1546095600_1563375600_result_point_idx ON public.submissions_1546095600_1563375600 USING btree (result, point);


--
-- Name: submissions_1546095600_1563375600_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1546095600_1563375600_user_id_epoch_second_idx ON public.submissions_1546095600_1563375600 USING btree (user_id, epoch_second);


--
-- Name: submissions_1546095600_1563375600_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1546095600_1563375600_user_id_execution_time_idx ON public.submissions_1546095600_1563375600 USING btree (user_id, execution_time);


--
-- Name: submissions_1546095600_1563375600_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1546095600_1563375600_user_id_length_idx ON public.submissions_1546095600_1563375600 USING btree (user_id, length);


--
-- Name: submissions_1546095600_1563375600_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1546095600_1563375600_user_id_point_idx ON public.submissions_1546095600_1563375600 USING btree (user_id, point);


--
-- Name: submissions_1546095600_15633756_execution_time_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1546095600_15633756_execution_time_epoch_second_idx ON public.submissions_1546095600_1563375600 USING btree (execution_time, epoch_second);


--
-- Name: submissions_1563375600_1580655600_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1563375600_1580655600_contest_id_epoch_second_idx ON public.submissions_1563375600_1580655600 USING btree (contest_id, epoch_second);


--
-- Name: submissions_1563375600_1580655600_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1563375600_1580655600_contest_id_execution_time_idx ON public.submissions_1563375600_1580655600 USING btree (contest_id, execution_time);


--
-- Name: submissions_1563375600_1580655600_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1563375600_1580655600_contest_id_length_idx ON public.submissions_1563375600_1580655600 USING btree (contest_id, length);


--
-- Name: submissions_1563375600_1580655600_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1563375600_1580655600_contest_id_point_idx ON public.submissions_1563375600_1580655600 USING btree (contest_id, point);


--
-- Name: submissions_1563375600_1580655600_epoch_second_updated_at_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1563375600_1580655600_epoch_second_updated_at_idx ON public.submissions_1563375600_1580655600 USING btree (epoch_second, updated_at);


--
-- Name: submissions_1563375600_1580655600_execution_time_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1563375600_1580655600_execution_time_length_idx ON public.submissions_1563375600_1580655600 USING btree (execution_time, length);


--
-- Name: submissions_1563375600_1580655600_execution_time_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1563375600_1580655600_execution_time_point_idx ON public.submissions_1563375600_1580655600 USING btree (execution_time, point);


--
-- Name: submissions_1563375600_1580655600_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_1563375600_1580655600_id_epoch_second_idx ON public.submissions_1563375600_1580655600 USING btree (id, epoch_second);


--
-- Name: submissions_1563375600_1580655600_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1563375600_1580655600_language_epoch_second_idx ON public.submissions_1563375600_1580655600 USING btree (language, epoch_second);


--
-- Name: submissions_1563375600_1580655600_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1563375600_1580655600_language_execution_time_idx ON public.submissions_1563375600_1580655600 USING btree (language, execution_time);


--
-- Name: submissions_1563375600_1580655600_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1563375600_1580655600_language_length_idx ON public.submissions_1563375600_1580655600 USING btree (language, length);


--
-- Name: submissions_1563375600_1580655600_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1563375600_1580655600_language_point_idx ON public.submissions_1563375600_1580655600 USING btree (language, point);


--
-- Name: submissions_1563375600_1580655600_length_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1563375600_1580655600_length_epoch_second_idx ON public.submissions_1563375600_1580655600 USING btree (length, epoch_second);


--
-- Name: submissions_1563375600_1580655600_length_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1563375600_1580655600_length_execution_time_idx ON public.submissions_1563375600_1580655600 USING btree (length, execution_time);


--
-- Name: submissions_1563375600_1580655600_length_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1563375600_1580655600_length_point_idx ON public.submissions_1563375600_1580655600 USING btree (length, point);


--
-- Name: submissions_1563375600_1580655600_point_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1563375600_1580655600_point_epoch_second_idx ON public.submissions_1563375600_1580655600 USING btree (point, epoch_second);


--
-- Name: submissions_1563375600_1580655600_point_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1563375600_1580655600_point_execution_time_idx ON public.submissions_1563375600_1580655600 USING btree (point, execution_time);


--
-- Name: submissions_1563375600_1580655600_point_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1563375600_1580655600_point_length_idx ON public.submissions_1563375600_1580655600 USING btree (point, length);


--
-- Name: submissions_1563375600_1580655600_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1563375600_1580655600_problem_id_epoch_second_idx ON public.submissions_1563375600_1580655600 USING btree (problem_id, epoch_second);


--
-- Name: submissions_1563375600_1580655600_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1563375600_1580655600_problem_id_execution_time_idx ON public.submissions_1563375600_1580655600 USING btree (problem_id, execution_time);


--
-- Name: submissions_1563375600_1580655600_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1563375600_1580655600_problem_id_length_idx ON public.submissions_1563375600_1580655600 USING btree (problem_id, length);


--
-- Name: submissions_1563375600_1580655600_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1563375600_1580655600_problem_id_point_idx ON public.submissions_1563375600_1580655600 USING btree (problem_id, point);


--
-- Name: submissions_1563375600_1580655600_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1563375600_1580655600_result_epoch_second_idx ON public.submissions_1563375600_1580655600 USING btree (result, epoch_second);


--
-- Name: submissions_1563375600_1580655600_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1563375600_1580655600_result_execution_time_idx ON public.submissions_1563375600_1580655600 USING btree (result, execution_time);


--
-- Name: submissions_1563375600_1580655600_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1563375600_1580655600_result_length_idx ON public.submissions_1563375600_1580655600 USING btree (result, length);


--
-- Name: submissions_1563375600_1580655600_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1563375600_1580655600_result_point_idx ON public.submissions_1563375600_1580655600 USING btree (result, point);


--
-- Name: submissions_1563375600_1580655600_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1563375600_1580655600_user_id_epoch_second_idx ON public.submissions_1563375600_1580655600 USING btree (user_id, epoch_second);


--
-- Name: submissions_1563375600_1580655600_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1563375600_1580655600_user_id_execution_time_idx ON public.submissions_1563375600_1580655600 USING btree (user_id, execution_time);


--
-- Name: submissions_1563375600_1580655600_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1563375600_1580655600_user_id_length_idx ON public.submissions_1563375600_1580655600 USING btree (user_id, length);


--
-- Name: submissions_1563375600_1580655600_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1563375600_1580655600_user_id_point_idx ON public.submissions_1563375600_1580655600 USING btree (user_id, point);


--
-- Name: submissions_1563375600_15806556_execution_time_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1563375600_15806556_execution_time_epoch_second_idx ON public.submissions_1563375600_1580655600 USING btree (execution_time, epoch_second);


--
-- Name: submissions_1580655600_1597935600_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1580655600_1597935600_contest_id_epoch_second_idx ON public.submissions_1580655600_1597935600 USING btree (contest_id, epoch_second);


--
-- Name: submissions_1580655600_1597935600_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1580655600_1597935600_contest_id_execution_time_idx ON public.submissions_1580655600_1597935600 USING btree (contest_id, execution_time);


--
-- Name: submissions_1580655600_1597935600_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1580655600_1597935600_contest_id_length_idx ON public.submissions_1580655600_1597935600 USING btree (contest_id, length);


--
-- Name: submissions_1580655600_1597935600_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1580655600_1597935600_contest_id_point_idx ON public.submissions_1580655600_1597935600 USING btree (contest_id, point);


--
-- Name: submissions_1580655600_1597935600_epoch_second_updated_at_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1580655600_1597935600_epoch_second_updated_at_idx ON public.submissions_1580655600_1597935600 USING btree (epoch_second, updated_at);


--
-- Name: submissions_1580655600_1597935600_execution_time_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1580655600_1597935600_execution_time_length_idx ON public.submissions_1580655600_1597935600 USING btree (execution_time, length);


--
-- Name: submissions_1580655600_1597935600_execution_time_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1580655600_1597935600_execution_time_point_idx ON public.submissions_1580655600_1597935600 USING btree (execution_time, point);


--
-- Name: submissions_1580655600_1597935600_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_1580655600_1597935600_id_epoch_second_idx ON public.submissions_1580655600_1597935600 USING btree (id, epoch_second);


--
-- Name: submissions_1580655600_1597935600_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1580655600_1597935600_language_epoch_second_idx ON public.submissions_1580655600_1597935600 USING btree (language, epoch_second);


--
-- Name: submissions_1580655600_1597935600_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1580655600_1597935600_language_execution_time_idx ON public.submissions_1580655600_1597935600 USING btree (language, execution_time);


--
-- Name: submissions_1580655600_1597935600_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1580655600_1597935600_language_length_idx ON public.submissions_1580655600_1597935600 USING btree (language, length);


--
-- Name: submissions_1580655600_1597935600_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1580655600_1597935600_language_point_idx ON public.submissions_1580655600_1597935600 USING btree (language, point);


--
-- Name: submissions_1580655600_1597935600_length_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1580655600_1597935600_length_epoch_second_idx ON public.submissions_1580655600_1597935600 USING btree (length, epoch_second);


--
-- Name: submissions_1580655600_1597935600_length_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1580655600_1597935600_length_execution_time_idx ON public.submissions_1580655600_1597935600 USING btree (length, execution_time);


--
-- Name: submissions_1580655600_1597935600_length_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1580655600_1597935600_length_point_idx ON public.submissions_1580655600_1597935600 USING btree (length, point);


--
-- Name: submissions_1580655600_1597935600_point_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1580655600_1597935600_point_epoch_second_idx ON public.submissions_1580655600_1597935600 USING btree (point, epoch_second);


--
-- Name: submissions_1580655600_1597935600_point_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1580655600_1597935600_point_execution_time_idx ON public.submissions_1580655600_1597935600 USING btree (point, execution_time);


--
-- Name: submissions_1580655600_1597935600_point_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1580655600_1597935600_point_length_idx ON public.submissions_1580655600_1597935600 USING btree (point, length);


--
-- Name: submissions_1580655600_1597935600_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1580655600_1597935600_problem_id_epoch_second_idx ON public.submissions_1580655600_1597935600 USING btree (problem_id, epoch_second);


--
-- Name: submissions_1580655600_1597935600_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1580655600_1597935600_problem_id_execution_time_idx ON public.submissions_1580655600_1597935600 USING btree (problem_id, execution_time);


--
-- Name: submissions_1580655600_1597935600_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1580655600_1597935600_problem_id_length_idx ON public.submissions_1580655600_1597935600 USING btree (problem_id, length);


--
-- Name: submissions_1580655600_1597935600_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1580655600_1597935600_problem_id_point_idx ON public.submissions_1580655600_1597935600 USING btree (problem_id, point);


--
-- Name: submissions_1580655600_1597935600_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1580655600_1597935600_result_epoch_second_idx ON public.submissions_1580655600_1597935600 USING btree (result, epoch_second);


--
-- Name: submissions_1580655600_1597935600_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1580655600_1597935600_result_execution_time_idx ON public.submissions_1580655600_1597935600 USING btree (result, execution_time);


--
-- Name: submissions_1580655600_1597935600_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1580655600_1597935600_result_length_idx ON public.submissions_1580655600_1597935600 USING btree (result, length);


--
-- Name: submissions_1580655600_1597935600_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1580655600_1597935600_result_point_idx ON public.submissions_1580655600_1597935600 USING btree (result, point);


--
-- Name: submissions_1580655600_1597935600_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1580655600_1597935600_user_id_epoch_second_idx ON public.submissions_1580655600_1597935600 USING btree (user_id, epoch_second);


--
-- Name: submissions_1580655600_1597935600_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1580655600_1597935600_user_id_execution_time_idx ON public.submissions_1580655600_1597935600 USING btree (user_id, execution_time);


--
-- Name: submissions_1580655600_1597935600_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1580655600_1597935600_user_id_length_idx ON public.submissions_1580655600_1597935600 USING btree (user_id, length);


--
-- Name: submissions_1580655600_1597935600_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1580655600_1597935600_user_id_point_idx ON public.submissions_1580655600_1597935600 USING btree (user_id, point);


--
-- Name: submissions_1580655600_15979356_execution_time_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1580655600_15979356_execution_time_epoch_second_idx ON public.submissions_1580655600_1597935600 USING btree (execution_time, epoch_second);


--
-- Name: submissions_1597935600_1615215600_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1597935600_1615215600_contest_id_epoch_second_idx ON public.submissions_1597935600_1615215600 USING btree (contest_id, epoch_second);


--
-- Name: submissions_1597935600_1615215600_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1597935600_1615215600_contest_id_execution_time_idx ON public.submissions_1597935600_1615215600 USING btree (contest_id, execution_time);


--
-- Name: submissions_1597935600_1615215600_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1597935600_1615215600_contest_id_length_idx ON public.submissions_1597935600_1615215600 USING btree (contest_id, length);


--
-- Name: submissions_1597935600_1615215600_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1597935600_1615215600_contest_id_point_idx ON public.submissions_1597935600_1615215600 USING btree (contest_id, point);


--
-- Name: submissions_1597935600_1615215600_epoch_second_updated_at_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1597935600_1615215600_epoch_second_updated_at_idx ON public.submissions_1597935600_1615215600 USING btree (epoch_second, updated_at);


--
-- Name: submissions_1597935600_1615215600_execution_time_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1597935600_1615215600_execution_time_length_idx ON public.submissions_1597935600_1615215600 USING btree (execution_time, length);


--
-- Name: submissions_1597935600_1615215600_execution_time_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1597935600_1615215600_execution_time_point_idx ON public.submissions_1597935600_1615215600 USING btree (execution_time, point);


--
-- Name: submissions_1597935600_1615215600_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_1597935600_1615215600_id_epoch_second_idx ON public.submissions_1597935600_1615215600 USING btree (id, epoch_second);


--
-- Name: submissions_1597935600_1615215600_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1597935600_1615215600_language_epoch_second_idx ON public.submissions_1597935600_1615215600 USING btree (language, epoch_second);


--
-- Name: submissions_1597935600_1615215600_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1597935600_1615215600_language_execution_time_idx ON public.submissions_1597935600_1615215600 USING btree (language, execution_time);


--
-- Name: submissions_1597935600_1615215600_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1597935600_1615215600_language_length_idx ON public.submissions_1597935600_1615215600 USING btree (language, length);


--
-- Name: submissions_1597935600_1615215600_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1597935600_1615215600_language_point_idx ON public.submissions_1597935600_1615215600 USING btree (language, point);


--
-- Name: submissions_1597935600_1615215600_length_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1597935600_1615215600_length_epoch_second_idx ON public.submissions_1597935600_1615215600 USING btree (length, epoch_second);


--
-- Name: submissions_1597935600_1615215600_length_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1597935600_1615215600_length_execution_time_idx ON public.submissions_1597935600_1615215600 USING btree (length, execution_time);


--
-- Name: submissions_1597935600_1615215600_length_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1597935600_1615215600_length_point_idx ON public.submissions_1597935600_1615215600 USING btree (length, point);


--
-- Name: submissions_1597935600_1615215600_point_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1597935600_1615215600_point_epoch_second_idx ON public.submissions_1597935600_1615215600 USING btree (point, epoch_second);


--
-- Name: submissions_1597935600_1615215600_point_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1597935600_1615215600_point_execution_time_idx ON public.submissions_1597935600_1615215600 USING btree (point, execution_time);


--
-- Name: submissions_1597935600_1615215600_point_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1597935600_1615215600_point_length_idx ON public.submissions_1597935600_1615215600 USING btree (point, length);


--
-- Name: submissions_1597935600_1615215600_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1597935600_1615215600_problem_id_epoch_second_idx ON public.submissions_1597935600_1615215600 USING btree (problem_id, epoch_second);


--
-- Name: submissions_1597935600_1615215600_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1597935600_1615215600_problem_id_execution_time_idx ON public.submissions_1597935600_1615215600 USING btree (problem_id, execution_time);


--
-- Name: submissions_1597935600_1615215600_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1597935600_1615215600_problem_id_length_idx ON public.submissions_1597935600_1615215600 USING btree (problem_id, length);


--
-- Name: submissions_1597935600_1615215600_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1597935600_1615215600_problem_id_point_idx ON public.submissions_1597935600_1615215600 USING btree (problem_id, point);


--
-- Name: submissions_1597935600_1615215600_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1597935600_1615215600_result_epoch_second_idx ON public.submissions_1597935600_1615215600 USING btree (result, epoch_second);


--
-- Name: submissions_1597935600_1615215600_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1597935600_1615215600_result_execution_time_idx ON public.submissions_1597935600_1615215600 USING btree (result, execution_time);


--
-- Name: submissions_1597935600_1615215600_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1597935600_1615215600_result_length_idx ON public.submissions_1597935600_1615215600 USING btree (result, length);


--
-- Name: submissions_1597935600_1615215600_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1597935600_1615215600_result_point_idx ON public.submissions_1597935600_1615215600 USING btree (result, point);


--
-- Name: submissions_1597935600_1615215600_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1597935600_1615215600_user_id_epoch_second_idx ON public.submissions_1597935600_1615215600 USING btree (user_id, epoch_second);


--
-- Name: submissions_1597935600_1615215600_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1597935600_1615215600_user_id_execution_time_idx ON public.submissions_1597935600_1615215600 USING btree (user_id, execution_time);


--
-- Name: submissions_1597935600_1615215600_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1597935600_1615215600_user_id_length_idx ON public.submissions_1597935600_1615215600 USING btree (user_id, length);


--
-- Name: submissions_1597935600_1615215600_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1597935600_1615215600_user_id_point_idx ON public.submissions_1597935600_1615215600 USING btree (user_id, point);


--
-- Name: submissions_1597935600_16152156_execution_time_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1597935600_16152156_execution_time_epoch_second_idx ON public.submissions_1597935600_1615215600 USING btree (execution_time, epoch_second);


--
-- Name: submissions_1615215600_1632495600_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1615215600_1632495600_contest_id_epoch_second_idx ON public.submissions_1615215600_1632495600 USING btree (contest_id, epoch_second);


--
-- Name: submissions_1615215600_1632495600_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1615215600_1632495600_contest_id_execution_time_idx ON public.submissions_1615215600_1632495600 USING btree (contest_id, execution_time);


--
-- Name: submissions_1615215600_1632495600_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1615215600_1632495600_contest_id_length_idx ON public.submissions_1615215600_1632495600 USING btree (contest_id, length);


--
-- Name: submissions_1615215600_1632495600_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1615215600_1632495600_contest_id_point_idx ON public.submissions_1615215600_1632495600 USING btree (contest_id, point);


--
-- Name: submissions_1615215600_1632495600_epoch_second_updated_at_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1615215600_1632495600_epoch_second_updated_at_idx ON public.submissions_1615215600_1632495600 USING btree (epoch_second, updated_at);


--
-- Name: submissions_1615215600_1632495600_execution_time_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1615215600_1632495600_execution_time_length_idx ON public.submissions_1615215600_1632495600 USING btree (execution_time, length);


--
-- Name: submissions_1615215600_1632495600_execution_time_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1615215600_1632495600_execution_time_point_idx ON public.submissions_1615215600_1632495600 USING btree (execution_time, point);


--
-- Name: submissions_1615215600_1632495600_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_1615215600_1632495600_id_epoch_second_idx ON public.submissions_1615215600_1632495600 USING btree (id, epoch_second);


--
-- Name: submissions_1615215600_1632495600_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1615215600_1632495600_language_epoch_second_idx ON public.submissions_1615215600_1632495600 USING btree (language, epoch_second);


--
-- Name: submissions_1615215600_1632495600_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1615215600_1632495600_language_execution_time_idx ON public.submissions_1615215600_1632495600 USING btree (language, execution_time);


--
-- Name: submissions_1615215600_1632495600_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1615215600_1632495600_language_length_idx ON public.submissions_1615215600_1632495600 USING btree (language, length);


--
-- Name: submissions_1615215600_1632495600_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1615215600_1632495600_language_point_idx ON public.submissions_1615215600_1632495600 USING btree (language, point);


--
-- Name: submissions_1615215600_1632495600_length_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1615215600_1632495600_length_epoch_second_idx ON public.submissions_1615215600_1632495600 USING btree (length, epoch_second);


--
-- Name: submissions_1615215600_1632495600_length_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1615215600_1632495600_length_execution_time_idx ON public.submissions_1615215600_1632495600 USING btree (length, execution_time);


--
-- Name: submissions_1615215600_1632495600_length_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1615215600_1632495600_length_point_idx ON public.submissions_1615215600_1632495600 USING btree (length, point);


--
-- Name: submissions_1615215600_1632495600_point_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1615215600_1632495600_point_epoch_second_idx ON public.submissions_1615215600_1632495600 USING btree (point, epoch_second);


--
-- Name: submissions_1615215600_1632495600_point_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1615215600_1632495600_point_execution_time_idx ON public.submissions_1615215600_1632495600 USING btree (point, execution_time);


--
-- Name: submissions_1615215600_1632495600_point_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1615215600_1632495600_point_length_idx ON public.submissions_1615215600_1632495600 USING btree (point, length);


--
-- Name: submissions_1615215600_1632495600_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1615215600_1632495600_problem_id_epoch_second_idx ON public.submissions_1615215600_1632495600 USING btree (problem_id, epoch_second);


--
-- Name: submissions_1615215600_1632495600_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1615215600_1632495600_problem_id_execution_time_idx ON public.submissions_1615215600_1632495600 USING btree (problem_id, execution_time);


--
-- Name: submissions_1615215600_1632495600_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1615215600_1632495600_problem_id_length_idx ON public.submissions_1615215600_1632495600 USING btree (problem_id, length);


--
-- Name: submissions_1615215600_1632495600_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1615215600_1632495600_problem_id_point_idx ON public.submissions_1615215600_1632495600 USING btree (problem_id, point);


--
-- Name: submissions_1615215600_1632495600_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1615215600_1632495600_result_epoch_second_idx ON public.submissions_1615215600_1632495600 USING btree (result, epoch_second);


--
-- Name: submissions_1615215600_1632495600_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1615215600_1632495600_result_execution_time_idx ON public.submissions_1615215600_1632495600 USING btree (result, execution_time);


--
-- Name: submissions_1615215600_1632495600_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1615215600_1632495600_result_length_idx ON public.submissions_1615215600_1632495600 USING btree (result, length);


--
-- Name: submissions_1615215600_1632495600_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1615215600_1632495600_result_point_idx ON public.submissions_1615215600_1632495600 USING btree (result, point);


--
-- Name: submissions_1615215600_1632495600_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1615215600_1632495600_user_id_epoch_second_idx ON public.submissions_1615215600_1632495600 USING btree (user_id, epoch_second);


--
-- Name: submissions_1615215600_1632495600_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1615215600_1632495600_user_id_execution_time_idx ON public.submissions_1615215600_1632495600 USING btree (user_id, execution_time);


--
-- Name: submissions_1615215600_1632495600_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1615215600_1632495600_user_id_length_idx ON public.submissions_1615215600_1632495600 USING btree (user_id, length);


--
-- Name: submissions_1615215600_1632495600_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1615215600_1632495600_user_id_point_idx ON public.submissions_1615215600_1632495600 USING btree (user_id, point);


--
-- Name: submissions_1615215600_16324956_execution_time_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1615215600_16324956_execution_time_epoch_second_idx ON public.submissions_1615215600_1632495600 USING btree (execution_time, epoch_second);


--
-- Name: submissions_1632495600_1649775600_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1632495600_1649775600_contest_id_epoch_second_idx ON public.submissions_1632495600_1649775600 USING btree (contest_id, epoch_second);


--
-- Name: submissions_1632495600_1649775600_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1632495600_1649775600_contest_id_execution_time_idx ON public.submissions_1632495600_1649775600 USING btree (contest_id, execution_time);


--
-- Name: submissions_1632495600_1649775600_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1632495600_1649775600_contest_id_length_idx ON public.submissions_1632495600_1649775600 USING btree (contest_id, length);


--
-- Name: submissions_1632495600_1649775600_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1632495600_1649775600_contest_id_point_idx ON public.submissions_1632495600_1649775600 USING btree (contest_id, point);


--
-- Name: submissions_1632495600_1649775600_epoch_second_updated_at_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1632495600_1649775600_epoch_second_updated_at_idx ON public.submissions_1632495600_1649775600 USING btree (epoch_second, updated_at);


--
-- Name: submissions_1632495600_1649775600_execution_time_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1632495600_1649775600_execution_time_length_idx ON public.submissions_1632495600_1649775600 USING btree (execution_time, length);


--
-- Name: submissions_1632495600_1649775600_execution_time_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1632495600_1649775600_execution_time_point_idx ON public.submissions_1632495600_1649775600 USING btree (execution_time, point);


--
-- Name: submissions_1632495600_1649775600_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_1632495600_1649775600_id_epoch_second_idx ON public.submissions_1632495600_1649775600 USING btree (id, epoch_second);


--
-- Name: submissions_1632495600_1649775600_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1632495600_1649775600_language_epoch_second_idx ON public.submissions_1632495600_1649775600 USING btree (language, epoch_second);


--
-- Name: submissions_1632495600_1649775600_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1632495600_1649775600_language_execution_time_idx ON public.submissions_1632495600_1649775600 USING btree (language, execution_time);


--
-- Name: submissions_1632495600_1649775600_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1632495600_1649775600_language_length_idx ON public.submissions_1632495600_1649775600 USING btree (language, length);


--
-- Name: submissions_1632495600_1649775600_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1632495600_1649775600_language_point_idx ON public.submissions_1632495600_1649775600 USING btree (language, point);


--
-- Name: submissions_1632495600_1649775600_length_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1632495600_1649775600_length_epoch_second_idx ON public.submissions_1632495600_1649775600 USING btree (length, epoch_second);


--
-- Name: submissions_1632495600_1649775600_length_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1632495600_1649775600_length_execution_time_idx ON public.submissions_1632495600_1649775600 USING btree (length, execution_time);


--
-- Name: submissions_1632495600_1649775600_length_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1632495600_1649775600_length_point_idx ON public.submissions_1632495600_1649775600 USING btree (length, point);


--
-- Name: submissions_1632495600_1649775600_point_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1632495600_1649775600_point_epoch_second_idx ON public.submissions_1632495600_1649775600 USING btree (point, epoch_second);


--
-- Name: submissions_1632495600_1649775600_point_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1632495600_1649775600_point_execution_time_idx ON public.submissions_1632495600_1649775600 USING btree (point, execution_time);


--
-- Name: submissions_1632495600_1649775600_point_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1632495600_1649775600_point_length_idx ON public.submissions_1632495600_1649775600 USING btree (point, length);


--
-- Name: submissions_1632495600_1649775600_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1632495600_1649775600_problem_id_epoch_second_idx ON public.submissions_1632495600_1649775600 USING btree (problem_id, epoch_second);


--
-- Name: submissions_1632495600_1649775600_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1632495600_1649775600_problem_id_execution_time_idx ON public.submissions_1632495600_1649775600 USING btree (problem_id, execution_time);


--
-- Name: submissions_1632495600_1649775600_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1632495600_1649775600_problem_id_length_idx ON public.submissions_1632495600_1649775600 USING btree (problem_id, length);


--
-- Name: submissions_1632495600_1649775600_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1632495600_1649775600_problem_id_point_idx ON public.submissions_1632495600_1649775600 USING btree (problem_id, point);


--
-- Name: submissions_1632495600_1649775600_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1632495600_1649775600_result_epoch_second_idx ON public.submissions_1632495600_1649775600 USING btree (result, epoch_second);


--
-- Name: submissions_1632495600_1649775600_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1632495600_1649775600_result_execution_time_idx ON public.submissions_1632495600_1649775600 USING btree (result, execution_time);


--
-- Name: submissions_1632495600_1649775600_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1632495600_1649775600_result_length_idx ON public.submissions_1632495600_1649775600 USING btree (result, length);


--
-- Name: submissions_1632495600_1649775600_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1632495600_1649775600_result_point_idx ON public.submissions_1632495600_1649775600 USING btree (result, point);


--
-- Name: submissions_1632495600_1649775600_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1632495600_1649775600_user_id_epoch_second_idx ON public.submissions_1632495600_1649775600 USING btree (user_id, epoch_second);


--
-- Name: submissions_1632495600_1649775600_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1632495600_1649775600_user_id_execution_time_idx ON public.submissions_1632495600_1649775600 USING btree (user_id, execution_time);


--
-- Name: submissions_1632495600_1649775600_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1632495600_1649775600_user_id_length_idx ON public.submissions_1632495600_1649775600 USING btree (user_id, length);


--
-- Name: submissions_1632495600_1649775600_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1632495600_1649775600_user_id_point_idx ON public.submissions_1632495600_1649775600 USING btree (user_id, point);


--
-- Name: submissions_1632495600_16497756_execution_time_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1632495600_16497756_execution_time_epoch_second_idx ON public.submissions_1632495600_1649775600 USING btree (execution_time, epoch_second);


--
-- Name: submissions_1649775600_1667055600_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1649775600_1667055600_contest_id_epoch_second_idx ON public.submissions_1649775600_1667055600 USING btree (contest_id, epoch_second);


--
-- Name: submissions_1649775600_1667055600_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1649775600_1667055600_contest_id_execution_time_idx ON public.submissions_1649775600_1667055600 USING btree (contest_id, execution_time);


--
-- Name: submissions_1649775600_1667055600_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1649775600_1667055600_contest_id_length_idx ON public.submissions_1649775600_1667055600 USING btree (contest_id, length);


--
-- Name: submissions_1649775600_1667055600_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1649775600_1667055600_contest_id_point_idx ON public.submissions_1649775600_1667055600 USING btree (contest_id, point);


--
-- Name: submissions_1649775600_1667055600_epoch_second_updated_at_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1649775600_1667055600_epoch_second_updated_at_idx ON public.submissions_1649775600_1667055600 USING btree (epoch_second, updated_at);


--
-- Name: submissions_1649775600_1667055600_execution_time_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1649775600_1667055600_execution_time_length_idx ON public.submissions_1649775600_1667055600 USING btree (execution_time, length);


--
-- Name: submissions_1649775600_1667055600_execution_time_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1649775600_1667055600_execution_time_point_idx ON public.submissions_1649775600_1667055600 USING btree (execution_time, point);


--
-- Name: submissions_1649775600_1667055600_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_1649775600_1667055600_id_epoch_second_idx ON public.submissions_1649775600_1667055600 USING btree (id, epoch_second);


--
-- Name: submissions_1649775600_1667055600_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1649775600_1667055600_language_epoch_second_idx ON public.submissions_1649775600_1667055600 USING btree (language, epoch_second);


--
-- Name: submissions_1649775600_1667055600_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1649775600_1667055600_language_execution_time_idx ON public.submissions_1649775600_1667055600 USING btree (language, execution_time);


--
-- Name: submissions_1649775600_1667055600_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1649775600_1667055600_language_length_idx ON public.submissions_1649775600_1667055600 USING btree (language, length);


--
-- Name: submissions_1649775600_1667055600_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1649775600_1667055600_language_point_idx ON public.submissions_1649775600_1667055600 USING btree (language, point);


--
-- Name: submissions_1649775600_1667055600_length_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1649775600_1667055600_length_epoch_second_idx ON public.submissions_1649775600_1667055600 USING btree (length, epoch_second);


--
-- Name: submissions_1649775600_1667055600_length_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1649775600_1667055600_length_execution_time_idx ON public.submissions_1649775600_1667055600 USING btree (length, execution_time);


--
-- Name: submissions_1649775600_1667055600_length_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1649775600_1667055600_length_point_idx ON public.submissions_1649775600_1667055600 USING btree (length, point);


--
-- Name: submissions_1649775600_1667055600_point_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1649775600_1667055600_point_epoch_second_idx ON public.submissions_1649775600_1667055600 USING btree (point, epoch_second);


--
-- Name: submissions_1649775600_1667055600_point_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1649775600_1667055600_point_execution_time_idx ON public.submissions_1649775600_1667055600 USING btree (point, execution_time);


--
-- Name: submissions_1649775600_1667055600_point_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1649775600_1667055600_point_length_idx ON public.submissions_1649775600_1667055600 USING btree (point, length);


--
-- Name: submissions_1649775600_1667055600_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1649775600_1667055600_problem_id_epoch_second_idx ON public.submissions_1649775600_1667055600 USING btree (problem_id, epoch_second);


--
-- Name: submissions_1649775600_1667055600_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1649775600_1667055600_problem_id_execution_time_idx ON public.submissions_1649775600_1667055600 USING btree (problem_id, execution_time);


--
-- Name: submissions_1649775600_1667055600_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1649775600_1667055600_problem_id_length_idx ON public.submissions_1649775600_1667055600 USING btree (problem_id, length);


--
-- Name: submissions_1649775600_1667055600_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1649775600_1667055600_problem_id_point_idx ON public.submissions_1649775600_1667055600 USING btree (problem_id, point);


--
-- Name: submissions_1649775600_1667055600_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1649775600_1667055600_result_epoch_second_idx ON public.submissions_1649775600_1667055600 USING btree (result, epoch_second);


--
-- Name: submissions_1649775600_1667055600_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1649775600_1667055600_result_execution_time_idx ON public.submissions_1649775600_1667055600 USING btree (result, execution_time);


--
-- Name: submissions_1649775600_1667055600_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1649775600_1667055600_result_length_idx ON public.submissions_1649775600_1667055600 USING btree (result, length);


--
-- Name: submissions_1649775600_1667055600_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1649775600_1667055600_result_point_idx ON public.submissions_1649775600_1667055600 USING btree (result, point);


--
-- Name: submissions_1649775600_1667055600_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1649775600_1667055600_user_id_epoch_second_idx ON public.submissions_1649775600_1667055600 USING btree (user_id, epoch_second);


--
-- Name: submissions_1649775600_1667055600_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1649775600_1667055600_user_id_execution_time_idx ON public.submissions_1649775600_1667055600 USING btree (user_id, execution_time);


--
-- Name: submissions_1649775600_1667055600_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1649775600_1667055600_user_id_length_idx ON public.submissions_1649775600_1667055600 USING btree (user_id, length);


--
-- Name: submissions_1649775600_1667055600_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1649775600_1667055600_user_id_point_idx ON public.submissions_1649775600_1667055600 USING btree (user_id, point);


--
-- Name: submissions_1649775600_16670556_execution_time_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1649775600_16670556_execution_time_epoch_second_idx ON public.submissions_1649775600_1667055600 USING btree (execution_time, epoch_second);


--
-- Name: submissions_1667055600_1684335600_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1667055600_1684335600_contest_id_epoch_second_idx ON public.submissions_1667055600_1684335600 USING btree (contest_id, epoch_second);


--
-- Name: submissions_1667055600_1684335600_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1667055600_1684335600_contest_id_execution_time_idx ON public.submissions_1667055600_1684335600 USING btree (contest_id, execution_time);


--
-- Name: submissions_1667055600_1684335600_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1667055600_1684335600_contest_id_length_idx ON public.submissions_1667055600_1684335600 USING btree (contest_id, length);


--
-- Name: submissions_1667055600_1684335600_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1667055600_1684335600_contest_id_point_idx ON public.submissions_1667055600_1684335600 USING btree (contest_id, point);


--
-- Name: submissions_1667055600_1684335600_epoch_second_updated_at_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1667055600_1684335600_epoch_second_updated_at_idx ON public.submissions_1667055600_1684335600 USING btree (epoch_second, updated_at);


--
-- Name: submissions_1667055600_1684335600_execution_time_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1667055600_1684335600_execution_time_length_idx ON public.submissions_1667055600_1684335600 USING btree (execution_time, length);


--
-- Name: submissions_1667055600_1684335600_execution_time_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1667055600_1684335600_execution_time_point_idx ON public.submissions_1667055600_1684335600 USING btree (execution_time, point);


--
-- Name: submissions_1667055600_1684335600_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_1667055600_1684335600_id_epoch_second_idx ON public.submissions_1667055600_1684335600 USING btree (id, epoch_second);


--
-- Name: submissions_1667055600_1684335600_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1667055600_1684335600_language_epoch_second_idx ON public.submissions_1667055600_1684335600 USING btree (language, epoch_second);


--
-- Name: submissions_1667055600_1684335600_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1667055600_1684335600_language_execution_time_idx ON public.submissions_1667055600_1684335600 USING btree (language, execution_time);


--
-- Name: submissions_1667055600_1684335600_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1667055600_1684335600_language_length_idx ON public.submissions_1667055600_1684335600 USING btree (language, length);


--
-- Name: submissions_1667055600_1684335600_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1667055600_1684335600_language_point_idx ON public.submissions_1667055600_1684335600 USING btree (language, point);


--
-- Name: submissions_1667055600_1684335600_length_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1667055600_1684335600_length_epoch_second_idx ON public.submissions_1667055600_1684335600 USING btree (length, epoch_second);


--
-- Name: submissions_1667055600_1684335600_length_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1667055600_1684335600_length_execution_time_idx ON public.submissions_1667055600_1684335600 USING btree (length, execution_time);


--
-- Name: submissions_1667055600_1684335600_length_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1667055600_1684335600_length_point_idx ON public.submissions_1667055600_1684335600 USING btree (length, point);


--
-- Name: submissions_1667055600_1684335600_point_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1667055600_1684335600_point_epoch_second_idx ON public.submissions_1667055600_1684335600 USING btree (point, epoch_second);


--
-- Name: submissions_1667055600_1684335600_point_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1667055600_1684335600_point_execution_time_idx ON public.submissions_1667055600_1684335600 USING btree (point, execution_time);


--
-- Name: submissions_1667055600_1684335600_point_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1667055600_1684335600_point_length_idx ON public.submissions_1667055600_1684335600 USING btree (point, length);


--
-- Name: submissions_1667055600_1684335600_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1667055600_1684335600_problem_id_epoch_second_idx ON public.submissions_1667055600_1684335600 USING btree (problem_id, epoch_second);


--
-- Name: submissions_1667055600_1684335600_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1667055600_1684335600_problem_id_execution_time_idx ON public.submissions_1667055600_1684335600 USING btree (problem_id, execution_time);


--
-- Name: submissions_1667055600_1684335600_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1667055600_1684335600_problem_id_length_idx ON public.submissions_1667055600_1684335600 USING btree (problem_id, length);


--
-- Name: submissions_1667055600_1684335600_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1667055600_1684335600_problem_id_point_idx ON public.submissions_1667055600_1684335600 USING btree (problem_id, point);


--
-- Name: submissions_1667055600_1684335600_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1667055600_1684335600_result_epoch_second_idx ON public.submissions_1667055600_1684335600 USING btree (result, epoch_second);


--
-- Name: submissions_1667055600_1684335600_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1667055600_1684335600_result_execution_time_idx ON public.submissions_1667055600_1684335600 USING btree (result, execution_time);


--
-- Name: submissions_1667055600_1684335600_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1667055600_1684335600_result_length_idx ON public.submissions_1667055600_1684335600 USING btree (result, length);


--
-- Name: submissions_1667055600_1684335600_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1667055600_1684335600_result_point_idx ON public.submissions_1667055600_1684335600 USING btree (result, point);


--
-- Name: submissions_1667055600_1684335600_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1667055600_1684335600_user_id_epoch_second_idx ON public.submissions_1667055600_1684335600 USING btree (user_id, epoch_second);


--
-- Name: submissions_1667055600_1684335600_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1667055600_1684335600_user_id_execution_time_idx ON public.submissions_1667055600_1684335600 USING btree (user_id, execution_time);


--
-- Name: submissions_1667055600_1684335600_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1667055600_1684335600_user_id_length_idx ON public.submissions_1667055600_1684335600 USING btree (user_id, length);


--
-- Name: submissions_1667055600_1684335600_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1667055600_1684335600_user_id_point_idx ON public.submissions_1667055600_1684335600 USING btree (user_id, point);


--
-- Name: submissions_1667055600_16843356_execution_time_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1667055600_16843356_execution_time_epoch_second_idx ON public.submissions_1667055600_1684335600 USING btree (execution_time, epoch_second);


--
-- Name: submissions_1684335600_1701615600_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1684335600_1701615600_contest_id_epoch_second_idx ON public.submissions_1684335600_1701615600 USING btree (contest_id, epoch_second);


--
-- Name: submissions_1684335600_1701615600_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1684335600_1701615600_contest_id_execution_time_idx ON public.submissions_1684335600_1701615600 USING btree (contest_id, execution_time);


--
-- Name: submissions_1684335600_1701615600_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1684335600_1701615600_contest_id_length_idx ON public.submissions_1684335600_1701615600 USING btree (contest_id, length);


--
-- Name: submissions_1684335600_1701615600_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1684335600_1701615600_contest_id_point_idx ON public.submissions_1684335600_1701615600 USING btree (contest_id, point);


--
-- Name: submissions_1684335600_1701615600_epoch_second_updated_at_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1684335600_1701615600_epoch_second_updated_at_idx ON public.submissions_1684335600_1701615600 USING btree (epoch_second, updated_at);


--
-- Name: submissions_1684335600_1701615600_execution_time_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1684335600_1701615600_execution_time_length_idx ON public.submissions_1684335600_1701615600 USING btree (execution_time, length);


--
-- Name: submissions_1684335600_1701615600_execution_time_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1684335600_1701615600_execution_time_point_idx ON public.submissions_1684335600_1701615600 USING btree (execution_time, point);


--
-- Name: submissions_1684335600_1701615600_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_1684335600_1701615600_id_epoch_second_idx ON public.submissions_1684335600_1701615600 USING btree (id, epoch_second);


--
-- Name: submissions_1684335600_1701615600_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1684335600_1701615600_language_epoch_second_idx ON public.submissions_1684335600_1701615600 USING btree (language, epoch_second);


--
-- Name: submissions_1684335600_1701615600_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1684335600_1701615600_language_execution_time_idx ON public.submissions_1684335600_1701615600 USING btree (language, execution_time);


--
-- Name: submissions_1684335600_1701615600_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1684335600_1701615600_language_length_idx ON public.submissions_1684335600_1701615600 USING btree (language, length);


--
-- Name: submissions_1684335600_1701615600_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1684335600_1701615600_language_point_idx ON public.submissions_1684335600_1701615600 USING btree (language, point);


--
-- Name: submissions_1684335600_1701615600_length_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1684335600_1701615600_length_epoch_second_idx ON public.submissions_1684335600_1701615600 USING btree (length, epoch_second);


--
-- Name: submissions_1684335600_1701615600_length_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1684335600_1701615600_length_execution_time_idx ON public.submissions_1684335600_1701615600 USING btree (length, execution_time);


--
-- Name: submissions_1684335600_1701615600_length_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1684335600_1701615600_length_point_idx ON public.submissions_1684335600_1701615600 USING btree (length, point);


--
-- Name: submissions_1684335600_1701615600_point_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1684335600_1701615600_point_epoch_second_idx ON public.submissions_1684335600_1701615600 USING btree (point, epoch_second);


--
-- Name: submissions_1684335600_1701615600_point_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1684335600_1701615600_point_execution_time_idx ON public.submissions_1684335600_1701615600 USING btree (point, execution_time);


--
-- Name: submissions_1684335600_1701615600_point_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1684335600_1701615600_point_length_idx ON public.submissions_1684335600_1701615600 USING btree (point, length);


--
-- Name: submissions_1684335600_1701615600_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1684335600_1701615600_problem_id_epoch_second_idx ON public.submissions_1684335600_1701615600 USING btree (problem_id, epoch_second);


--
-- Name: submissions_1684335600_1701615600_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1684335600_1701615600_problem_id_execution_time_idx ON public.submissions_1684335600_1701615600 USING btree (problem_id, execution_time);


--
-- Name: submissions_1684335600_1701615600_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1684335600_1701615600_problem_id_length_idx ON public.submissions_1684335600_1701615600 USING btree (problem_id, length);


--
-- Name: submissions_1684335600_1701615600_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1684335600_1701615600_problem_id_point_idx ON public.submissions_1684335600_1701615600 USING btree (problem_id, point);


--
-- Name: submissions_1684335600_1701615600_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1684335600_1701615600_result_epoch_second_idx ON public.submissions_1684335600_1701615600 USING btree (result, epoch_second);


--
-- Name: submissions_1684335600_1701615600_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1684335600_1701615600_result_execution_time_idx ON public.submissions_1684335600_1701615600 USING btree (result, execution_time);


--
-- Name: submissions_1684335600_1701615600_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1684335600_1701615600_result_length_idx ON public.submissions_1684335600_1701615600 USING btree (result, length);


--
-- Name: submissions_1684335600_1701615600_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1684335600_1701615600_result_point_idx ON public.submissions_1684335600_1701615600 USING btree (result, point);


--
-- Name: submissions_1684335600_1701615600_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1684335600_1701615600_user_id_epoch_second_idx ON public.submissions_1684335600_1701615600 USING btree (user_id, epoch_second);


--
-- Name: submissions_1684335600_1701615600_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1684335600_1701615600_user_id_execution_time_idx ON public.submissions_1684335600_1701615600 USING btree (user_id, execution_time);


--
-- Name: submissions_1684335600_1701615600_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1684335600_1701615600_user_id_length_idx ON public.submissions_1684335600_1701615600 USING btree (user_id, length);


--
-- Name: submissions_1684335600_1701615600_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1684335600_1701615600_user_id_point_idx ON public.submissions_1684335600_1701615600 USING btree (user_id, point);


--
-- Name: submissions_1684335600_17016156_execution_time_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1684335600_17016156_execution_time_epoch_second_idx ON public.submissions_1684335600_1701615600 USING btree (execution_time, epoch_second);


--
-- Name: submissions_1701615600_1718895600_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1701615600_1718895600_contest_id_epoch_second_idx ON public.submissions_1701615600_1718895600 USING btree (contest_id, epoch_second);


--
-- Name: submissions_1701615600_1718895600_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1701615600_1718895600_contest_id_execution_time_idx ON public.submissions_1701615600_1718895600 USING btree (contest_id, execution_time);


--
-- Name: submissions_1701615600_1718895600_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1701615600_1718895600_contest_id_length_idx ON public.submissions_1701615600_1718895600 USING btree (contest_id, length);


--
-- Name: submissions_1701615600_1718895600_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1701615600_1718895600_contest_id_point_idx ON public.submissions_1701615600_1718895600 USING btree (contest_id, point);


--
-- Name: submissions_1701615600_1718895600_epoch_second_updated_at_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1701615600_1718895600_epoch_second_updated_at_idx ON public.submissions_1701615600_1718895600 USING btree (epoch_second, updated_at);


--
-- Name: submissions_1701615600_1718895600_execution_time_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1701615600_1718895600_execution_time_length_idx ON public.submissions_1701615600_1718895600 USING btree (execution_time, length);


--
-- Name: submissions_1701615600_1718895600_execution_time_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1701615600_1718895600_execution_time_point_idx ON public.submissions_1701615600_1718895600 USING btree (execution_time, point);


--
-- Name: submissions_1701615600_1718895600_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_1701615600_1718895600_id_epoch_second_idx ON public.submissions_1701615600_1718895600 USING btree (id, epoch_second);


--
-- Name: submissions_1701615600_1718895600_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1701615600_1718895600_language_epoch_second_idx ON public.submissions_1701615600_1718895600 USING btree (language, epoch_second);


--
-- Name: submissions_1701615600_1718895600_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1701615600_1718895600_language_execution_time_idx ON public.submissions_1701615600_1718895600 USING btree (language, execution_time);


--
-- Name: submissions_1701615600_1718895600_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1701615600_1718895600_language_length_idx ON public.submissions_1701615600_1718895600 USING btree (language, length);


--
-- Name: submissions_1701615600_1718895600_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1701615600_1718895600_language_point_idx ON public.submissions_1701615600_1718895600 USING btree (language, point);


--
-- Name: submissions_1701615600_1718895600_length_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1701615600_1718895600_length_epoch_second_idx ON public.submissions_1701615600_1718895600 USING btree (length, epoch_second);


--
-- Name: submissions_1701615600_1718895600_length_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1701615600_1718895600_length_execution_time_idx ON public.submissions_1701615600_1718895600 USING btree (length, execution_time);


--
-- Name: submissions_1701615600_1718895600_length_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1701615600_1718895600_length_point_idx ON public.submissions_1701615600_1718895600 USING btree (length, point);


--
-- Name: submissions_1701615600_1718895600_point_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1701615600_1718895600_point_epoch_second_idx ON public.submissions_1701615600_1718895600 USING btree (point, epoch_second);


--
-- Name: submissions_1701615600_1718895600_point_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1701615600_1718895600_point_execution_time_idx ON public.submissions_1701615600_1718895600 USING btree (point, execution_time);


--
-- Name: submissions_1701615600_1718895600_point_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1701615600_1718895600_point_length_idx ON public.submissions_1701615600_1718895600 USING btree (point, length);


--
-- Name: submissions_1701615600_1718895600_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1701615600_1718895600_problem_id_epoch_second_idx ON public.submissions_1701615600_1718895600 USING btree (problem_id, epoch_second);


--
-- Name: submissions_1701615600_1718895600_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1701615600_1718895600_problem_id_execution_time_idx ON public.submissions_1701615600_1718895600 USING btree (problem_id, execution_time);


--
-- Name: submissions_1701615600_1718895600_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1701615600_1718895600_problem_id_length_idx ON public.submissions_1701615600_1718895600 USING btree (problem_id, length);


--
-- Name: submissions_1701615600_1718895600_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1701615600_1718895600_problem_id_point_idx ON public.submissions_1701615600_1718895600 USING btree (problem_id, point);


--
-- Name: submissions_1701615600_1718895600_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1701615600_1718895600_result_epoch_second_idx ON public.submissions_1701615600_1718895600 USING btree (result, epoch_second);


--
-- Name: submissions_1701615600_1718895600_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1701615600_1718895600_result_execution_time_idx ON public.submissions_1701615600_1718895600 USING btree (result, execution_time);


--
-- Name: submissions_1701615600_1718895600_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1701615600_1718895600_result_length_idx ON public.submissions_1701615600_1718895600 USING btree (result, length);


--
-- Name: submissions_1701615600_1718895600_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1701615600_1718895600_result_point_idx ON public.submissions_1701615600_1718895600 USING btree (result, point);


--
-- Name: submissions_1701615600_1718895600_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1701615600_1718895600_user_id_epoch_second_idx ON public.submissions_1701615600_1718895600 USING btree (user_id, epoch_second);


--
-- Name: submissions_1701615600_1718895600_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1701615600_1718895600_user_id_execution_time_idx ON public.submissions_1701615600_1718895600 USING btree (user_id, execution_time);


--
-- Name: submissions_1701615600_1718895600_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1701615600_1718895600_user_id_length_idx ON public.submissions_1701615600_1718895600 USING btree (user_id, length);


--
-- Name: submissions_1701615600_1718895600_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1701615600_1718895600_user_id_point_idx ON public.submissions_1701615600_1718895600 USING btree (user_id, point);


--
-- Name: submissions_1701615600_17188956_execution_time_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1701615600_17188956_execution_time_epoch_second_idx ON public.submissions_1701615600_1718895600 USING btree (execution_time, epoch_second);


--
-- Name: submissions_1718895600_1736175600_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1718895600_1736175600_contest_id_epoch_second_idx ON public.submissions_1718895600_1736175600 USING btree (contest_id, epoch_second);


--
-- Name: submissions_1718895600_1736175600_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1718895600_1736175600_contest_id_execution_time_idx ON public.submissions_1718895600_1736175600 USING btree (contest_id, execution_time);


--
-- Name: submissions_1718895600_1736175600_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1718895600_1736175600_contest_id_length_idx ON public.submissions_1718895600_1736175600 USING btree (contest_id, length);


--
-- Name: submissions_1718895600_1736175600_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1718895600_1736175600_contest_id_point_idx ON public.submissions_1718895600_1736175600 USING btree (contest_id, point);


--
-- Name: submissions_1718895600_1736175600_epoch_second_updated_at_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1718895600_1736175600_epoch_second_updated_at_idx ON public.submissions_1718895600_1736175600 USING btree (epoch_second, updated_at);


--
-- Name: submissions_1718895600_1736175600_execution_time_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1718895600_1736175600_execution_time_length_idx ON public.submissions_1718895600_1736175600 USING btree (execution_time, length);


--
-- Name: submissions_1718895600_1736175600_execution_time_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1718895600_1736175600_execution_time_point_idx ON public.submissions_1718895600_1736175600 USING btree (execution_time, point);


--
-- Name: submissions_1718895600_1736175600_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_1718895600_1736175600_id_epoch_second_idx ON public.submissions_1718895600_1736175600 USING btree (id, epoch_second);


--
-- Name: submissions_1718895600_1736175600_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1718895600_1736175600_language_epoch_second_idx ON public.submissions_1718895600_1736175600 USING btree (language, epoch_second);


--
-- Name: submissions_1718895600_1736175600_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1718895600_1736175600_language_execution_time_idx ON public.submissions_1718895600_1736175600 USING btree (language, execution_time);


--
-- Name: submissions_1718895600_1736175600_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1718895600_1736175600_language_length_idx ON public.submissions_1718895600_1736175600 USING btree (language, length);


--
-- Name: submissions_1718895600_1736175600_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1718895600_1736175600_language_point_idx ON public.submissions_1718895600_1736175600 USING btree (language, point);


--
-- Name: submissions_1718895600_1736175600_length_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1718895600_1736175600_length_epoch_second_idx ON public.submissions_1718895600_1736175600 USING btree (length, epoch_second);


--
-- Name: submissions_1718895600_1736175600_length_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1718895600_1736175600_length_execution_time_idx ON public.submissions_1718895600_1736175600 USING btree (length, execution_time);


--
-- Name: submissions_1718895600_1736175600_length_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1718895600_1736175600_length_point_idx ON public.submissions_1718895600_1736175600 USING btree (length, point);


--
-- Name: submissions_1718895600_1736175600_point_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1718895600_1736175600_point_epoch_second_idx ON public.submissions_1718895600_1736175600 USING btree (point, epoch_second);


--
-- Name: submissions_1718895600_1736175600_point_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1718895600_1736175600_point_execution_time_idx ON public.submissions_1718895600_1736175600 USING btree (point, execution_time);


--
-- Name: submissions_1718895600_1736175600_point_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1718895600_1736175600_point_length_idx ON public.submissions_1718895600_1736175600 USING btree (point, length);


--
-- Name: submissions_1718895600_1736175600_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1718895600_1736175600_problem_id_epoch_second_idx ON public.submissions_1718895600_1736175600 USING btree (problem_id, epoch_second);


--
-- Name: submissions_1718895600_1736175600_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1718895600_1736175600_problem_id_execution_time_idx ON public.submissions_1718895600_1736175600 USING btree (problem_id, execution_time);


--
-- Name: submissions_1718895600_1736175600_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1718895600_1736175600_problem_id_length_idx ON public.submissions_1718895600_1736175600 USING btree (problem_id, length);


--
-- Name: submissions_1718895600_1736175600_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1718895600_1736175600_problem_id_point_idx ON public.submissions_1718895600_1736175600 USING btree (problem_id, point);


--
-- Name: submissions_1718895600_1736175600_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1718895600_1736175600_result_epoch_second_idx ON public.submissions_1718895600_1736175600 USING btree (result, epoch_second);


--
-- Name: submissions_1718895600_1736175600_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1718895600_1736175600_result_execution_time_idx ON public.submissions_1718895600_1736175600 USING btree (result, execution_time);


--
-- Name: submissions_1718895600_1736175600_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1718895600_1736175600_result_length_idx ON public.submissions_1718895600_1736175600 USING btree (result, length);


--
-- Name: submissions_1718895600_1736175600_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1718895600_1736175600_result_point_idx ON public.submissions_1718895600_1736175600 USING btree (result, point);


--
-- Name: submissions_1718895600_1736175600_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1718895600_1736175600_user_id_epoch_second_idx ON public.submissions_1718895600_1736175600 USING btree (user_id, epoch_second);


--
-- Name: submissions_1718895600_1736175600_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1718895600_1736175600_user_id_execution_time_idx ON public.submissions_1718895600_1736175600 USING btree (user_id, execution_time);


--
-- Name: submissions_1718895600_1736175600_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1718895600_1736175600_user_id_length_idx ON public.submissions_1718895600_1736175600 USING btree (user_id, length);


--
-- Name: submissions_1718895600_1736175600_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1718895600_1736175600_user_id_point_idx ON public.submissions_1718895600_1736175600 USING btree (user_id, point);


--
-- Name: submissions_1718895600_17361756_execution_time_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1718895600_17361756_execution_time_epoch_second_idx ON public.submissions_1718895600_1736175600 USING btree (execution_time, epoch_second);


--
-- Name: submissions_1736175600_1753455600_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1736175600_1753455600_contest_id_epoch_second_idx ON public.submissions_1736175600_1753455600 USING btree (contest_id, epoch_second);


--
-- Name: submissions_1736175600_1753455600_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1736175600_1753455600_contest_id_execution_time_idx ON public.submissions_1736175600_1753455600 USING btree (contest_id, execution_time);


--
-- Name: submissions_1736175600_1753455600_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1736175600_1753455600_contest_id_length_idx ON public.submissions_1736175600_1753455600 USING btree (contest_id, length);


--
-- Name: submissions_1736175600_1753455600_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1736175600_1753455600_contest_id_point_idx ON public.submissions_1736175600_1753455600 USING btree (contest_id, point);


--
-- Name: submissions_1736175600_1753455600_epoch_second_updated_at_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1736175600_1753455600_epoch_second_updated_at_idx ON public.submissions_1736175600_1753455600 USING btree (epoch_second, updated_at);


--
-- Name: submissions_1736175600_1753455600_execution_time_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1736175600_1753455600_execution_time_length_idx ON public.submissions_1736175600_1753455600 USING btree (execution_time, length);


--
-- Name: submissions_1736175600_1753455600_execution_time_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1736175600_1753455600_execution_time_point_idx ON public.submissions_1736175600_1753455600 USING btree (execution_time, point);


--
-- Name: submissions_1736175600_1753455600_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_1736175600_1753455600_id_epoch_second_idx ON public.submissions_1736175600_1753455600 USING btree (id, epoch_second);


--
-- Name: submissions_1736175600_1753455600_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1736175600_1753455600_language_epoch_second_idx ON public.submissions_1736175600_1753455600 USING btree (language, epoch_second);


--
-- Name: submissions_1736175600_1753455600_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1736175600_1753455600_language_execution_time_idx ON public.submissions_1736175600_1753455600 USING btree (language, execution_time);


--
-- Name: submissions_1736175600_1753455600_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1736175600_1753455600_language_length_idx ON public.submissions_1736175600_1753455600 USING btree (language, length);


--
-- Name: submissions_1736175600_1753455600_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1736175600_1753455600_language_point_idx ON public.submissions_1736175600_1753455600 USING btree (language, point);


--
-- Name: submissions_1736175600_1753455600_length_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1736175600_1753455600_length_epoch_second_idx ON public.submissions_1736175600_1753455600 USING btree (length, epoch_second);


--
-- Name: submissions_1736175600_1753455600_length_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1736175600_1753455600_length_execution_time_idx ON public.submissions_1736175600_1753455600 USING btree (length, execution_time);


--
-- Name: submissions_1736175600_1753455600_length_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1736175600_1753455600_length_point_idx ON public.submissions_1736175600_1753455600 USING btree (length, point);


--
-- Name: submissions_1736175600_1753455600_point_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1736175600_1753455600_point_epoch_second_idx ON public.submissions_1736175600_1753455600 USING btree (point, epoch_second);


--
-- Name: submissions_1736175600_1753455600_point_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1736175600_1753455600_point_execution_time_idx ON public.submissions_1736175600_1753455600 USING btree (point, execution_time);


--
-- Name: submissions_1736175600_1753455600_point_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1736175600_1753455600_point_length_idx ON public.submissions_1736175600_1753455600 USING btree (point, length);


--
-- Name: submissions_1736175600_1753455600_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1736175600_1753455600_problem_id_epoch_second_idx ON public.submissions_1736175600_1753455600 USING btree (problem_id, epoch_second);


--
-- Name: submissions_1736175600_1753455600_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1736175600_1753455600_problem_id_execution_time_idx ON public.submissions_1736175600_1753455600 USING btree (problem_id, execution_time);


--
-- Name: submissions_1736175600_1753455600_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1736175600_1753455600_problem_id_length_idx ON public.submissions_1736175600_1753455600 USING btree (problem_id, length);


--
-- Name: submissions_1736175600_1753455600_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1736175600_1753455600_problem_id_point_idx ON public.submissions_1736175600_1753455600 USING btree (problem_id, point);


--
-- Name: submissions_1736175600_1753455600_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1736175600_1753455600_result_epoch_second_idx ON public.submissions_1736175600_1753455600 USING btree (result, epoch_second);


--
-- Name: submissions_1736175600_1753455600_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1736175600_1753455600_result_execution_time_idx ON public.submissions_1736175600_1753455600 USING btree (result, execution_time);


--
-- Name: submissions_1736175600_1753455600_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1736175600_1753455600_result_length_idx ON public.submissions_1736175600_1753455600 USING btree (result, length);


--
-- Name: submissions_1736175600_1753455600_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1736175600_1753455600_result_point_idx ON public.submissions_1736175600_1753455600 USING btree (result, point);


--
-- Name: submissions_1736175600_1753455600_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1736175600_1753455600_user_id_epoch_second_idx ON public.submissions_1736175600_1753455600 USING btree (user_id, epoch_second);


--
-- Name: submissions_1736175600_1753455600_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1736175600_1753455600_user_id_execution_time_idx ON public.submissions_1736175600_1753455600 USING btree (user_id, execution_time);


--
-- Name: submissions_1736175600_1753455600_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1736175600_1753455600_user_id_length_idx ON public.submissions_1736175600_1753455600 USING btree (user_id, length);


--
-- Name: submissions_1736175600_1753455600_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1736175600_1753455600_user_id_point_idx ON public.submissions_1736175600_1753455600 USING btree (user_id, point);


--
-- Name: submissions_1736175600_17534556_execution_time_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1736175600_17534556_execution_time_epoch_second_idx ON public.submissions_1736175600_1753455600 USING btree (execution_time, epoch_second);


--
-- Name: submissions_1753455600_1770735600_contest_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1753455600_1770735600_contest_id_epoch_second_idx ON public.submissions_1753455600_1770735600 USING btree (contest_id, epoch_second);


--
-- Name: submissions_1753455600_1770735600_contest_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1753455600_1770735600_contest_id_execution_time_idx ON public.submissions_1753455600_1770735600 USING btree (contest_id, execution_time);


--
-- Name: submissions_1753455600_1770735600_contest_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1753455600_1770735600_contest_id_length_idx ON public.submissions_1753455600_1770735600 USING btree (contest_id, length);


--
-- Name: submissions_1753455600_1770735600_contest_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1753455600_1770735600_contest_id_point_idx ON public.submissions_1753455600_1770735600 USING btree (contest_id, point);


--
-- Name: submissions_1753455600_1770735600_epoch_second_updated_at_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1753455600_1770735600_epoch_second_updated_at_idx ON public.submissions_1753455600_1770735600 USING btree (epoch_second, updated_at);


--
-- Name: submissions_1753455600_1770735600_execution_time_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1753455600_1770735600_execution_time_length_idx ON public.submissions_1753455600_1770735600 USING btree (execution_time, length);


--
-- Name: submissions_1753455600_1770735600_execution_time_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1753455600_1770735600_execution_time_point_idx ON public.submissions_1753455600_1770735600 USING btree (execution_time, point);


--
-- Name: submissions_1753455600_1770735600_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX submissions_1753455600_1770735600_id_epoch_second_idx ON public.submissions_1753455600_1770735600 USING btree (id, epoch_second);


--
-- Name: submissions_1753455600_1770735600_language_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1753455600_1770735600_language_epoch_second_idx ON public.submissions_1753455600_1770735600 USING btree (language, epoch_second);


--
-- Name: submissions_1753455600_1770735600_language_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1753455600_1770735600_language_execution_time_idx ON public.submissions_1753455600_1770735600 USING btree (language, execution_time);


--
-- Name: submissions_1753455600_1770735600_language_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1753455600_1770735600_language_length_idx ON public.submissions_1753455600_1770735600 USING btree (language, length);


--
-- Name: submissions_1753455600_1770735600_language_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1753455600_1770735600_language_point_idx ON public.submissions_1753455600_1770735600 USING btree (language, point);


--
-- Name: submissions_1753455600_1770735600_length_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1753455600_1770735600_length_epoch_second_idx ON public.submissions_1753455600_1770735600 USING btree (length, epoch_second);


--
-- Name: submissions_1753455600_1770735600_length_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1753455600_1770735600_length_execution_time_idx ON public.submissions_1753455600_1770735600 USING btree (length, execution_time);


--
-- Name: submissions_1753455600_1770735600_length_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1753455600_1770735600_length_point_idx ON public.submissions_1753455600_1770735600 USING btree (length, point);


--
-- Name: submissions_1753455600_1770735600_point_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1753455600_1770735600_point_epoch_second_idx ON public.submissions_1753455600_1770735600 USING btree (point, epoch_second);


--
-- Name: submissions_1753455600_1770735600_point_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1753455600_1770735600_point_execution_time_idx ON public.submissions_1753455600_1770735600 USING btree (point, execution_time);


--
-- Name: submissions_1753455600_1770735600_point_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1753455600_1770735600_point_length_idx ON public.submissions_1753455600_1770735600 USING btree (point, length);


--
-- Name: submissions_1753455600_1770735600_problem_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1753455600_1770735600_problem_id_epoch_second_idx ON public.submissions_1753455600_1770735600 USING btree (problem_id, epoch_second);


--
-- Name: submissions_1753455600_1770735600_problem_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1753455600_1770735600_problem_id_execution_time_idx ON public.submissions_1753455600_1770735600 USING btree (problem_id, execution_time);


--
-- Name: submissions_1753455600_1770735600_problem_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1753455600_1770735600_problem_id_length_idx ON public.submissions_1753455600_1770735600 USING btree (problem_id, length);


--
-- Name: submissions_1753455600_1770735600_problem_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1753455600_1770735600_problem_id_point_idx ON public.submissions_1753455600_1770735600 USING btree (problem_id, point);


--
-- Name: submissions_1753455600_1770735600_result_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1753455600_1770735600_result_epoch_second_idx ON public.submissions_1753455600_1770735600 USING btree (result, epoch_second);


--
-- Name: submissions_1753455600_1770735600_result_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1753455600_1770735600_result_execution_time_idx ON public.submissions_1753455600_1770735600 USING btree (result, execution_time);


--
-- Name: submissions_1753455600_1770735600_result_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1753455600_1770735600_result_length_idx ON public.submissions_1753455600_1770735600 USING btree (result, length);


--
-- Name: submissions_1753455600_1770735600_result_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1753455600_1770735600_result_point_idx ON public.submissions_1753455600_1770735600 USING btree (result, point);


--
-- Name: submissions_1753455600_1770735600_user_id_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1753455600_1770735600_user_id_epoch_second_idx ON public.submissions_1753455600_1770735600 USING btree (user_id, epoch_second);


--
-- Name: submissions_1753455600_1770735600_user_id_execution_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1753455600_1770735600_user_id_execution_time_idx ON public.submissions_1753455600_1770735600 USING btree (user_id, execution_time);


--
-- Name: submissions_1753455600_1770735600_user_id_length_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1753455600_1770735600_user_id_length_idx ON public.submissions_1753455600_1770735600 USING btree (user_id, length);


--
-- Name: submissions_1753455600_1770735600_user_id_point_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1753455600_1770735600_user_id_point_idx ON public.submissions_1753455600_1770735600 USING btree (user_id, point);


--
-- Name: submissions_1753455600_17707356_execution_time_epoch_second_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX submissions_1753455600_17707356_execution_time_epoch_second_idx ON public.submissions_1753455600_1770735600 USING btree (execution_time, epoch_second);


--
-- Name: submissions_0_1304175600_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_0_1304175600_contest_id_epoch_second_idx;


--
-- Name: submissions_0_1304175600_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_0_1304175600_contest_id_execution_time_idx;


--
-- Name: submissions_0_1304175600_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_0_1304175600_contest_id_length_idx;


--
-- Name: submissions_0_1304175600_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_0_1304175600_contest_id_point_idx;


--
-- Name: submissions_0_1304175600_epoch_second_updated_at_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_updated_at_index ATTACH PARTITION public.submissions_0_1304175600_epoch_second_updated_at_idx;


--
-- Name: submissions_0_1304175600_execution_time_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_epoch_second_index ATTACH PARTITION public.submissions_0_1304175600_execution_time_epoch_second_idx;


--
-- Name: submissions_0_1304175600_execution_time_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_length_index ATTACH PARTITION public.submissions_0_1304175600_execution_time_length_idx;


--
-- Name: submissions_0_1304175600_execution_time_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_point_index ATTACH PARTITION public.submissions_0_1304175600_execution_time_point_idx;


--
-- Name: submissions_0_1304175600_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_id_epoch_second_unique ATTACH PARTITION public.submissions_0_1304175600_id_epoch_second_idx;


--
-- Name: submissions_0_1304175600_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_0_1304175600_language_epoch_second_idx;


--
-- Name: submissions_0_1304175600_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_0_1304175600_language_execution_time_idx;


--
-- Name: submissions_0_1304175600_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_0_1304175600_language_length_idx;


--
-- Name: submissions_0_1304175600_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_0_1304175600_language_point_idx;


--
-- Name: submissions_0_1304175600_length_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_epoch_second_index ATTACH PARTITION public.submissions_0_1304175600_length_epoch_second_idx;


--
-- Name: submissions_0_1304175600_length_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_execution_time_index ATTACH PARTITION public.submissions_0_1304175600_length_execution_time_idx;


--
-- Name: submissions_0_1304175600_length_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_point_index ATTACH PARTITION public.submissions_0_1304175600_length_point_idx;


--
-- Name: submissions_0_1304175600_point_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_epoch_second_index ATTACH PARTITION public.submissions_0_1304175600_point_epoch_second_idx;


--
-- Name: submissions_0_1304175600_point_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_execution_time_index ATTACH PARTITION public.submissions_0_1304175600_point_execution_time_idx;


--
-- Name: submissions_0_1304175600_point_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_length_index ATTACH PARTITION public.submissions_0_1304175600_point_length_idx;


--
-- Name: submissions_0_1304175600_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_0_1304175600_problem_id_epoch_second_idx;


--
-- Name: submissions_0_1304175600_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_0_1304175600_problem_id_execution_time_idx;


--
-- Name: submissions_0_1304175600_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_0_1304175600_problem_id_length_idx;


--
-- Name: submissions_0_1304175600_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_0_1304175600_problem_id_point_idx;


--
-- Name: submissions_0_1304175600_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_0_1304175600_result_epoch_second_idx;


--
-- Name: submissions_0_1304175600_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_0_1304175600_result_execution_time_idx;


--
-- Name: submissions_0_1304175600_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_0_1304175600_result_length_idx;


--
-- Name: submissions_0_1304175600_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_0_1304175600_result_point_idx;


--
-- Name: submissions_0_1304175600_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_0_1304175600_user_id_epoch_second_idx;


--
-- Name: submissions_0_1304175600_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_0_1304175600_user_id_execution_time_idx;


--
-- Name: submissions_0_1304175600_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_0_1304175600_user_id_length_idx;


--
-- Name: submissions_0_1304175600_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_0_1304175600_user_id_point_idx;


--
-- Name: submissions_1304175600_1321455600_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_1304175600_1321455600_contest_id_epoch_second_idx;


--
-- Name: submissions_1304175600_1321455600_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_1304175600_1321455600_contest_id_execution_time_idx;


--
-- Name: submissions_1304175600_1321455600_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_1304175600_1321455600_contest_id_length_idx;


--
-- Name: submissions_1304175600_1321455600_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_1304175600_1321455600_contest_id_point_idx;


--
-- Name: submissions_1304175600_1321455600_epoch_second_updated_at_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_updated_at_index ATTACH PARTITION public.submissions_1304175600_1321455600_epoch_second_updated_at_idx;


--
-- Name: submissions_1304175600_1321455600_execution_time_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_length_index ATTACH PARTITION public.submissions_1304175600_1321455600_execution_time_length_idx;


--
-- Name: submissions_1304175600_1321455600_execution_time_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_point_index ATTACH PARTITION public.submissions_1304175600_1321455600_execution_time_point_idx;


--
-- Name: submissions_1304175600_1321455600_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_id_epoch_second_unique ATTACH PARTITION public.submissions_1304175600_1321455600_id_epoch_second_idx;


--
-- Name: submissions_1304175600_1321455600_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_1304175600_1321455600_language_epoch_second_idx;


--
-- Name: submissions_1304175600_1321455600_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_1304175600_1321455600_language_execution_time_idx;


--
-- Name: submissions_1304175600_1321455600_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_1304175600_1321455600_language_length_idx;


--
-- Name: submissions_1304175600_1321455600_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_1304175600_1321455600_language_point_idx;


--
-- Name: submissions_1304175600_1321455600_length_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_epoch_second_index ATTACH PARTITION public.submissions_1304175600_1321455600_length_epoch_second_idx;


--
-- Name: submissions_1304175600_1321455600_length_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_execution_time_index ATTACH PARTITION public.submissions_1304175600_1321455600_length_execution_time_idx;


--
-- Name: submissions_1304175600_1321455600_length_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_point_index ATTACH PARTITION public.submissions_1304175600_1321455600_length_point_idx;


--
-- Name: submissions_1304175600_1321455600_point_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_epoch_second_index ATTACH PARTITION public.submissions_1304175600_1321455600_point_epoch_second_idx;


--
-- Name: submissions_1304175600_1321455600_point_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_execution_time_index ATTACH PARTITION public.submissions_1304175600_1321455600_point_execution_time_idx;


--
-- Name: submissions_1304175600_1321455600_point_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_length_index ATTACH PARTITION public.submissions_1304175600_1321455600_point_length_idx;


--
-- Name: submissions_1304175600_1321455600_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_1304175600_1321455600_problem_id_epoch_second_idx;


--
-- Name: submissions_1304175600_1321455600_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_1304175600_1321455600_problem_id_execution_time_idx;


--
-- Name: submissions_1304175600_1321455600_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_1304175600_1321455600_problem_id_length_idx;


--
-- Name: submissions_1304175600_1321455600_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_1304175600_1321455600_problem_id_point_idx;


--
-- Name: submissions_1304175600_1321455600_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_1304175600_1321455600_result_epoch_second_idx;


--
-- Name: submissions_1304175600_1321455600_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_1304175600_1321455600_result_execution_time_idx;


--
-- Name: submissions_1304175600_1321455600_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_1304175600_1321455600_result_length_idx;


--
-- Name: submissions_1304175600_1321455600_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_1304175600_1321455600_result_point_idx;


--
-- Name: submissions_1304175600_1321455600_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_1304175600_1321455600_user_id_epoch_second_idx;


--
-- Name: submissions_1304175600_1321455600_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_1304175600_1321455600_user_id_execution_time_idx;


--
-- Name: submissions_1304175600_1321455600_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_1304175600_1321455600_user_id_length_idx;


--
-- Name: submissions_1304175600_1321455600_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_1304175600_1321455600_user_id_point_idx;


--
-- Name: submissions_1304175600_13214556_execution_time_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_epoch_second_index ATTACH PARTITION public.submissions_1304175600_13214556_execution_time_epoch_second_idx;


--
-- Name: submissions_1321455600_1338735600_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_1321455600_1338735600_contest_id_epoch_second_idx;


--
-- Name: submissions_1321455600_1338735600_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_1321455600_1338735600_contest_id_execution_time_idx;


--
-- Name: submissions_1321455600_1338735600_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_1321455600_1338735600_contest_id_length_idx;


--
-- Name: submissions_1321455600_1338735600_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_1321455600_1338735600_contest_id_point_idx;


--
-- Name: submissions_1321455600_1338735600_epoch_second_updated_at_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_updated_at_index ATTACH PARTITION public.submissions_1321455600_1338735600_epoch_second_updated_at_idx;


--
-- Name: submissions_1321455600_1338735600_execution_time_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_length_index ATTACH PARTITION public.submissions_1321455600_1338735600_execution_time_length_idx;


--
-- Name: submissions_1321455600_1338735600_execution_time_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_point_index ATTACH PARTITION public.submissions_1321455600_1338735600_execution_time_point_idx;


--
-- Name: submissions_1321455600_1338735600_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_id_epoch_second_unique ATTACH PARTITION public.submissions_1321455600_1338735600_id_epoch_second_idx;


--
-- Name: submissions_1321455600_1338735600_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_1321455600_1338735600_language_epoch_second_idx;


--
-- Name: submissions_1321455600_1338735600_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_1321455600_1338735600_language_execution_time_idx;


--
-- Name: submissions_1321455600_1338735600_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_1321455600_1338735600_language_length_idx;


--
-- Name: submissions_1321455600_1338735600_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_1321455600_1338735600_language_point_idx;


--
-- Name: submissions_1321455600_1338735600_length_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_epoch_second_index ATTACH PARTITION public.submissions_1321455600_1338735600_length_epoch_second_idx;


--
-- Name: submissions_1321455600_1338735600_length_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_execution_time_index ATTACH PARTITION public.submissions_1321455600_1338735600_length_execution_time_idx;


--
-- Name: submissions_1321455600_1338735600_length_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_point_index ATTACH PARTITION public.submissions_1321455600_1338735600_length_point_idx;


--
-- Name: submissions_1321455600_1338735600_point_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_epoch_second_index ATTACH PARTITION public.submissions_1321455600_1338735600_point_epoch_second_idx;


--
-- Name: submissions_1321455600_1338735600_point_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_execution_time_index ATTACH PARTITION public.submissions_1321455600_1338735600_point_execution_time_idx;


--
-- Name: submissions_1321455600_1338735600_point_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_length_index ATTACH PARTITION public.submissions_1321455600_1338735600_point_length_idx;


--
-- Name: submissions_1321455600_1338735600_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_1321455600_1338735600_problem_id_epoch_second_idx;


--
-- Name: submissions_1321455600_1338735600_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_1321455600_1338735600_problem_id_execution_time_idx;


--
-- Name: submissions_1321455600_1338735600_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_1321455600_1338735600_problem_id_length_idx;


--
-- Name: submissions_1321455600_1338735600_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_1321455600_1338735600_problem_id_point_idx;


--
-- Name: submissions_1321455600_1338735600_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_1321455600_1338735600_result_epoch_second_idx;


--
-- Name: submissions_1321455600_1338735600_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_1321455600_1338735600_result_execution_time_idx;


--
-- Name: submissions_1321455600_1338735600_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_1321455600_1338735600_result_length_idx;


--
-- Name: submissions_1321455600_1338735600_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_1321455600_1338735600_result_point_idx;


--
-- Name: submissions_1321455600_1338735600_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_1321455600_1338735600_user_id_epoch_second_idx;


--
-- Name: submissions_1321455600_1338735600_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_1321455600_1338735600_user_id_execution_time_idx;


--
-- Name: submissions_1321455600_1338735600_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_1321455600_1338735600_user_id_length_idx;


--
-- Name: submissions_1321455600_1338735600_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_1321455600_1338735600_user_id_point_idx;


--
-- Name: submissions_1321455600_13387356_execution_time_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_epoch_second_index ATTACH PARTITION public.submissions_1321455600_13387356_execution_time_epoch_second_idx;


--
-- Name: submissions_1338735600_1356015600_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_1338735600_1356015600_contest_id_epoch_second_idx;


--
-- Name: submissions_1338735600_1356015600_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_1338735600_1356015600_contest_id_execution_time_idx;


--
-- Name: submissions_1338735600_1356015600_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_1338735600_1356015600_contest_id_length_idx;


--
-- Name: submissions_1338735600_1356015600_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_1338735600_1356015600_contest_id_point_idx;


--
-- Name: submissions_1338735600_1356015600_epoch_second_updated_at_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_updated_at_index ATTACH PARTITION public.submissions_1338735600_1356015600_epoch_second_updated_at_idx;


--
-- Name: submissions_1338735600_1356015600_execution_time_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_length_index ATTACH PARTITION public.submissions_1338735600_1356015600_execution_time_length_idx;


--
-- Name: submissions_1338735600_1356015600_execution_time_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_point_index ATTACH PARTITION public.submissions_1338735600_1356015600_execution_time_point_idx;


--
-- Name: submissions_1338735600_1356015600_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_id_epoch_second_unique ATTACH PARTITION public.submissions_1338735600_1356015600_id_epoch_second_idx;


--
-- Name: submissions_1338735600_1356015600_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_1338735600_1356015600_language_epoch_second_idx;


--
-- Name: submissions_1338735600_1356015600_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_1338735600_1356015600_language_execution_time_idx;


--
-- Name: submissions_1338735600_1356015600_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_1338735600_1356015600_language_length_idx;


--
-- Name: submissions_1338735600_1356015600_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_1338735600_1356015600_language_point_idx;


--
-- Name: submissions_1338735600_1356015600_length_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_epoch_second_index ATTACH PARTITION public.submissions_1338735600_1356015600_length_epoch_second_idx;


--
-- Name: submissions_1338735600_1356015600_length_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_execution_time_index ATTACH PARTITION public.submissions_1338735600_1356015600_length_execution_time_idx;


--
-- Name: submissions_1338735600_1356015600_length_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_point_index ATTACH PARTITION public.submissions_1338735600_1356015600_length_point_idx;


--
-- Name: submissions_1338735600_1356015600_point_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_epoch_second_index ATTACH PARTITION public.submissions_1338735600_1356015600_point_epoch_second_idx;


--
-- Name: submissions_1338735600_1356015600_point_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_execution_time_index ATTACH PARTITION public.submissions_1338735600_1356015600_point_execution_time_idx;


--
-- Name: submissions_1338735600_1356015600_point_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_length_index ATTACH PARTITION public.submissions_1338735600_1356015600_point_length_idx;


--
-- Name: submissions_1338735600_1356015600_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_1338735600_1356015600_problem_id_epoch_second_idx;


--
-- Name: submissions_1338735600_1356015600_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_1338735600_1356015600_problem_id_execution_time_idx;


--
-- Name: submissions_1338735600_1356015600_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_1338735600_1356015600_problem_id_length_idx;


--
-- Name: submissions_1338735600_1356015600_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_1338735600_1356015600_problem_id_point_idx;


--
-- Name: submissions_1338735600_1356015600_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_1338735600_1356015600_result_epoch_second_idx;


--
-- Name: submissions_1338735600_1356015600_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_1338735600_1356015600_result_execution_time_idx;


--
-- Name: submissions_1338735600_1356015600_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_1338735600_1356015600_result_length_idx;


--
-- Name: submissions_1338735600_1356015600_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_1338735600_1356015600_result_point_idx;


--
-- Name: submissions_1338735600_1356015600_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_1338735600_1356015600_user_id_epoch_second_idx;


--
-- Name: submissions_1338735600_1356015600_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_1338735600_1356015600_user_id_execution_time_idx;


--
-- Name: submissions_1338735600_1356015600_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_1338735600_1356015600_user_id_length_idx;


--
-- Name: submissions_1338735600_1356015600_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_1338735600_1356015600_user_id_point_idx;


--
-- Name: submissions_1338735600_13560156_execution_time_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_epoch_second_index ATTACH PARTITION public.submissions_1338735600_13560156_execution_time_epoch_second_idx;


--
-- Name: submissions_1356015600_1373295600_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_1356015600_1373295600_contest_id_epoch_second_idx;


--
-- Name: submissions_1356015600_1373295600_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_1356015600_1373295600_contest_id_execution_time_idx;


--
-- Name: submissions_1356015600_1373295600_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_1356015600_1373295600_contest_id_length_idx;


--
-- Name: submissions_1356015600_1373295600_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_1356015600_1373295600_contest_id_point_idx;


--
-- Name: submissions_1356015600_1373295600_epoch_second_updated_at_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_updated_at_index ATTACH PARTITION public.submissions_1356015600_1373295600_epoch_second_updated_at_idx;


--
-- Name: submissions_1356015600_1373295600_execution_time_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_length_index ATTACH PARTITION public.submissions_1356015600_1373295600_execution_time_length_idx;


--
-- Name: submissions_1356015600_1373295600_execution_time_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_point_index ATTACH PARTITION public.submissions_1356015600_1373295600_execution_time_point_idx;


--
-- Name: submissions_1356015600_1373295600_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_id_epoch_second_unique ATTACH PARTITION public.submissions_1356015600_1373295600_id_epoch_second_idx;


--
-- Name: submissions_1356015600_1373295600_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_1356015600_1373295600_language_epoch_second_idx;


--
-- Name: submissions_1356015600_1373295600_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_1356015600_1373295600_language_execution_time_idx;


--
-- Name: submissions_1356015600_1373295600_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_1356015600_1373295600_language_length_idx;


--
-- Name: submissions_1356015600_1373295600_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_1356015600_1373295600_language_point_idx;


--
-- Name: submissions_1356015600_1373295600_length_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_epoch_second_index ATTACH PARTITION public.submissions_1356015600_1373295600_length_epoch_second_idx;


--
-- Name: submissions_1356015600_1373295600_length_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_execution_time_index ATTACH PARTITION public.submissions_1356015600_1373295600_length_execution_time_idx;


--
-- Name: submissions_1356015600_1373295600_length_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_point_index ATTACH PARTITION public.submissions_1356015600_1373295600_length_point_idx;


--
-- Name: submissions_1356015600_1373295600_point_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_epoch_second_index ATTACH PARTITION public.submissions_1356015600_1373295600_point_epoch_second_idx;


--
-- Name: submissions_1356015600_1373295600_point_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_execution_time_index ATTACH PARTITION public.submissions_1356015600_1373295600_point_execution_time_idx;


--
-- Name: submissions_1356015600_1373295600_point_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_length_index ATTACH PARTITION public.submissions_1356015600_1373295600_point_length_idx;


--
-- Name: submissions_1356015600_1373295600_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_1356015600_1373295600_problem_id_epoch_second_idx;


--
-- Name: submissions_1356015600_1373295600_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_1356015600_1373295600_problem_id_execution_time_idx;


--
-- Name: submissions_1356015600_1373295600_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_1356015600_1373295600_problem_id_length_idx;


--
-- Name: submissions_1356015600_1373295600_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_1356015600_1373295600_problem_id_point_idx;


--
-- Name: submissions_1356015600_1373295600_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_1356015600_1373295600_result_epoch_second_idx;


--
-- Name: submissions_1356015600_1373295600_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_1356015600_1373295600_result_execution_time_idx;


--
-- Name: submissions_1356015600_1373295600_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_1356015600_1373295600_result_length_idx;


--
-- Name: submissions_1356015600_1373295600_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_1356015600_1373295600_result_point_idx;


--
-- Name: submissions_1356015600_1373295600_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_1356015600_1373295600_user_id_epoch_second_idx;


--
-- Name: submissions_1356015600_1373295600_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_1356015600_1373295600_user_id_execution_time_idx;


--
-- Name: submissions_1356015600_1373295600_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_1356015600_1373295600_user_id_length_idx;


--
-- Name: submissions_1356015600_1373295600_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_1356015600_1373295600_user_id_point_idx;


--
-- Name: submissions_1356015600_13732956_execution_time_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_epoch_second_index ATTACH PARTITION public.submissions_1356015600_13732956_execution_time_epoch_second_idx;


--
-- Name: submissions_1373295600_1390575600_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_1373295600_1390575600_contest_id_epoch_second_idx;


--
-- Name: submissions_1373295600_1390575600_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_1373295600_1390575600_contest_id_execution_time_idx;


--
-- Name: submissions_1373295600_1390575600_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_1373295600_1390575600_contest_id_length_idx;


--
-- Name: submissions_1373295600_1390575600_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_1373295600_1390575600_contest_id_point_idx;


--
-- Name: submissions_1373295600_1390575600_epoch_second_updated_at_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_updated_at_index ATTACH PARTITION public.submissions_1373295600_1390575600_epoch_second_updated_at_idx;


--
-- Name: submissions_1373295600_1390575600_execution_time_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_length_index ATTACH PARTITION public.submissions_1373295600_1390575600_execution_time_length_idx;


--
-- Name: submissions_1373295600_1390575600_execution_time_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_point_index ATTACH PARTITION public.submissions_1373295600_1390575600_execution_time_point_idx;


--
-- Name: submissions_1373295600_1390575600_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_id_epoch_second_unique ATTACH PARTITION public.submissions_1373295600_1390575600_id_epoch_second_idx;


--
-- Name: submissions_1373295600_1390575600_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_1373295600_1390575600_language_epoch_second_idx;


--
-- Name: submissions_1373295600_1390575600_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_1373295600_1390575600_language_execution_time_idx;


--
-- Name: submissions_1373295600_1390575600_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_1373295600_1390575600_language_length_idx;


--
-- Name: submissions_1373295600_1390575600_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_1373295600_1390575600_language_point_idx;


--
-- Name: submissions_1373295600_1390575600_length_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_epoch_second_index ATTACH PARTITION public.submissions_1373295600_1390575600_length_epoch_second_idx;


--
-- Name: submissions_1373295600_1390575600_length_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_execution_time_index ATTACH PARTITION public.submissions_1373295600_1390575600_length_execution_time_idx;


--
-- Name: submissions_1373295600_1390575600_length_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_point_index ATTACH PARTITION public.submissions_1373295600_1390575600_length_point_idx;


--
-- Name: submissions_1373295600_1390575600_point_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_epoch_second_index ATTACH PARTITION public.submissions_1373295600_1390575600_point_epoch_second_idx;


--
-- Name: submissions_1373295600_1390575600_point_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_execution_time_index ATTACH PARTITION public.submissions_1373295600_1390575600_point_execution_time_idx;


--
-- Name: submissions_1373295600_1390575600_point_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_length_index ATTACH PARTITION public.submissions_1373295600_1390575600_point_length_idx;


--
-- Name: submissions_1373295600_1390575600_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_1373295600_1390575600_problem_id_epoch_second_idx;


--
-- Name: submissions_1373295600_1390575600_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_1373295600_1390575600_problem_id_execution_time_idx;


--
-- Name: submissions_1373295600_1390575600_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_1373295600_1390575600_problem_id_length_idx;


--
-- Name: submissions_1373295600_1390575600_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_1373295600_1390575600_problem_id_point_idx;


--
-- Name: submissions_1373295600_1390575600_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_1373295600_1390575600_result_epoch_second_idx;


--
-- Name: submissions_1373295600_1390575600_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_1373295600_1390575600_result_execution_time_idx;


--
-- Name: submissions_1373295600_1390575600_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_1373295600_1390575600_result_length_idx;


--
-- Name: submissions_1373295600_1390575600_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_1373295600_1390575600_result_point_idx;


--
-- Name: submissions_1373295600_1390575600_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_1373295600_1390575600_user_id_epoch_second_idx;


--
-- Name: submissions_1373295600_1390575600_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_1373295600_1390575600_user_id_execution_time_idx;


--
-- Name: submissions_1373295600_1390575600_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_1373295600_1390575600_user_id_length_idx;


--
-- Name: submissions_1373295600_1390575600_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_1373295600_1390575600_user_id_point_idx;


--
-- Name: submissions_1373295600_13905756_execution_time_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_epoch_second_index ATTACH PARTITION public.submissions_1373295600_13905756_execution_time_epoch_second_idx;


--
-- Name: submissions_1390575600_1407855600_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_1390575600_1407855600_contest_id_epoch_second_idx;


--
-- Name: submissions_1390575600_1407855600_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_1390575600_1407855600_contest_id_execution_time_idx;


--
-- Name: submissions_1390575600_1407855600_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_1390575600_1407855600_contest_id_length_idx;


--
-- Name: submissions_1390575600_1407855600_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_1390575600_1407855600_contest_id_point_idx;


--
-- Name: submissions_1390575600_1407855600_epoch_second_updated_at_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_updated_at_index ATTACH PARTITION public.submissions_1390575600_1407855600_epoch_second_updated_at_idx;


--
-- Name: submissions_1390575600_1407855600_execution_time_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_length_index ATTACH PARTITION public.submissions_1390575600_1407855600_execution_time_length_idx;


--
-- Name: submissions_1390575600_1407855600_execution_time_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_point_index ATTACH PARTITION public.submissions_1390575600_1407855600_execution_time_point_idx;


--
-- Name: submissions_1390575600_1407855600_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_id_epoch_second_unique ATTACH PARTITION public.submissions_1390575600_1407855600_id_epoch_second_idx;


--
-- Name: submissions_1390575600_1407855600_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_1390575600_1407855600_language_epoch_second_idx;


--
-- Name: submissions_1390575600_1407855600_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_1390575600_1407855600_language_execution_time_idx;


--
-- Name: submissions_1390575600_1407855600_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_1390575600_1407855600_language_length_idx;


--
-- Name: submissions_1390575600_1407855600_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_1390575600_1407855600_language_point_idx;


--
-- Name: submissions_1390575600_1407855600_length_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_epoch_second_index ATTACH PARTITION public.submissions_1390575600_1407855600_length_epoch_second_idx;


--
-- Name: submissions_1390575600_1407855600_length_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_execution_time_index ATTACH PARTITION public.submissions_1390575600_1407855600_length_execution_time_idx;


--
-- Name: submissions_1390575600_1407855600_length_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_point_index ATTACH PARTITION public.submissions_1390575600_1407855600_length_point_idx;


--
-- Name: submissions_1390575600_1407855600_point_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_epoch_second_index ATTACH PARTITION public.submissions_1390575600_1407855600_point_epoch_second_idx;


--
-- Name: submissions_1390575600_1407855600_point_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_execution_time_index ATTACH PARTITION public.submissions_1390575600_1407855600_point_execution_time_idx;


--
-- Name: submissions_1390575600_1407855600_point_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_length_index ATTACH PARTITION public.submissions_1390575600_1407855600_point_length_idx;


--
-- Name: submissions_1390575600_1407855600_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_1390575600_1407855600_problem_id_epoch_second_idx;


--
-- Name: submissions_1390575600_1407855600_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_1390575600_1407855600_problem_id_execution_time_idx;


--
-- Name: submissions_1390575600_1407855600_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_1390575600_1407855600_problem_id_length_idx;


--
-- Name: submissions_1390575600_1407855600_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_1390575600_1407855600_problem_id_point_idx;


--
-- Name: submissions_1390575600_1407855600_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_1390575600_1407855600_result_epoch_second_idx;


--
-- Name: submissions_1390575600_1407855600_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_1390575600_1407855600_result_execution_time_idx;


--
-- Name: submissions_1390575600_1407855600_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_1390575600_1407855600_result_length_idx;


--
-- Name: submissions_1390575600_1407855600_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_1390575600_1407855600_result_point_idx;


--
-- Name: submissions_1390575600_1407855600_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_1390575600_1407855600_user_id_epoch_second_idx;


--
-- Name: submissions_1390575600_1407855600_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_1390575600_1407855600_user_id_execution_time_idx;


--
-- Name: submissions_1390575600_1407855600_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_1390575600_1407855600_user_id_length_idx;


--
-- Name: submissions_1390575600_1407855600_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_1390575600_1407855600_user_id_point_idx;


--
-- Name: submissions_1390575600_14078556_execution_time_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_epoch_second_index ATTACH PARTITION public.submissions_1390575600_14078556_execution_time_epoch_second_idx;


--
-- Name: submissions_1407855600_1425135600_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_1407855600_1425135600_contest_id_epoch_second_idx;


--
-- Name: submissions_1407855600_1425135600_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_1407855600_1425135600_contest_id_execution_time_idx;


--
-- Name: submissions_1407855600_1425135600_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_1407855600_1425135600_contest_id_length_idx;


--
-- Name: submissions_1407855600_1425135600_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_1407855600_1425135600_contest_id_point_idx;


--
-- Name: submissions_1407855600_1425135600_epoch_second_updated_at_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_updated_at_index ATTACH PARTITION public.submissions_1407855600_1425135600_epoch_second_updated_at_idx;


--
-- Name: submissions_1407855600_1425135600_execution_time_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_length_index ATTACH PARTITION public.submissions_1407855600_1425135600_execution_time_length_idx;


--
-- Name: submissions_1407855600_1425135600_execution_time_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_point_index ATTACH PARTITION public.submissions_1407855600_1425135600_execution_time_point_idx;


--
-- Name: submissions_1407855600_1425135600_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_id_epoch_second_unique ATTACH PARTITION public.submissions_1407855600_1425135600_id_epoch_second_idx;


--
-- Name: submissions_1407855600_1425135600_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_1407855600_1425135600_language_epoch_second_idx;


--
-- Name: submissions_1407855600_1425135600_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_1407855600_1425135600_language_execution_time_idx;


--
-- Name: submissions_1407855600_1425135600_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_1407855600_1425135600_language_length_idx;


--
-- Name: submissions_1407855600_1425135600_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_1407855600_1425135600_language_point_idx;


--
-- Name: submissions_1407855600_1425135600_length_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_epoch_second_index ATTACH PARTITION public.submissions_1407855600_1425135600_length_epoch_second_idx;


--
-- Name: submissions_1407855600_1425135600_length_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_execution_time_index ATTACH PARTITION public.submissions_1407855600_1425135600_length_execution_time_idx;


--
-- Name: submissions_1407855600_1425135600_length_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_point_index ATTACH PARTITION public.submissions_1407855600_1425135600_length_point_idx;


--
-- Name: submissions_1407855600_1425135600_point_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_epoch_second_index ATTACH PARTITION public.submissions_1407855600_1425135600_point_epoch_second_idx;


--
-- Name: submissions_1407855600_1425135600_point_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_execution_time_index ATTACH PARTITION public.submissions_1407855600_1425135600_point_execution_time_idx;


--
-- Name: submissions_1407855600_1425135600_point_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_length_index ATTACH PARTITION public.submissions_1407855600_1425135600_point_length_idx;


--
-- Name: submissions_1407855600_1425135600_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_1407855600_1425135600_problem_id_epoch_second_idx;


--
-- Name: submissions_1407855600_1425135600_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_1407855600_1425135600_problem_id_execution_time_idx;


--
-- Name: submissions_1407855600_1425135600_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_1407855600_1425135600_problem_id_length_idx;


--
-- Name: submissions_1407855600_1425135600_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_1407855600_1425135600_problem_id_point_idx;


--
-- Name: submissions_1407855600_1425135600_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_1407855600_1425135600_result_epoch_second_idx;


--
-- Name: submissions_1407855600_1425135600_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_1407855600_1425135600_result_execution_time_idx;


--
-- Name: submissions_1407855600_1425135600_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_1407855600_1425135600_result_length_idx;


--
-- Name: submissions_1407855600_1425135600_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_1407855600_1425135600_result_point_idx;


--
-- Name: submissions_1407855600_1425135600_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_1407855600_1425135600_user_id_epoch_second_idx;


--
-- Name: submissions_1407855600_1425135600_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_1407855600_1425135600_user_id_execution_time_idx;


--
-- Name: submissions_1407855600_1425135600_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_1407855600_1425135600_user_id_length_idx;


--
-- Name: submissions_1407855600_1425135600_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_1407855600_1425135600_user_id_point_idx;


--
-- Name: submissions_1407855600_14251356_execution_time_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_epoch_second_index ATTACH PARTITION public.submissions_1407855600_14251356_execution_time_epoch_second_idx;


--
-- Name: submissions_1425135600_1442415600_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_1425135600_1442415600_contest_id_epoch_second_idx;


--
-- Name: submissions_1425135600_1442415600_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_1425135600_1442415600_contest_id_execution_time_idx;


--
-- Name: submissions_1425135600_1442415600_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_1425135600_1442415600_contest_id_length_idx;


--
-- Name: submissions_1425135600_1442415600_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_1425135600_1442415600_contest_id_point_idx;


--
-- Name: submissions_1425135600_1442415600_epoch_second_updated_at_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_updated_at_index ATTACH PARTITION public.submissions_1425135600_1442415600_epoch_second_updated_at_idx;


--
-- Name: submissions_1425135600_1442415600_execution_time_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_length_index ATTACH PARTITION public.submissions_1425135600_1442415600_execution_time_length_idx;


--
-- Name: submissions_1425135600_1442415600_execution_time_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_point_index ATTACH PARTITION public.submissions_1425135600_1442415600_execution_time_point_idx;


--
-- Name: submissions_1425135600_1442415600_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_id_epoch_second_unique ATTACH PARTITION public.submissions_1425135600_1442415600_id_epoch_second_idx;


--
-- Name: submissions_1425135600_1442415600_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_1425135600_1442415600_language_epoch_second_idx;


--
-- Name: submissions_1425135600_1442415600_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_1425135600_1442415600_language_execution_time_idx;


--
-- Name: submissions_1425135600_1442415600_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_1425135600_1442415600_language_length_idx;


--
-- Name: submissions_1425135600_1442415600_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_1425135600_1442415600_language_point_idx;


--
-- Name: submissions_1425135600_1442415600_length_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_epoch_second_index ATTACH PARTITION public.submissions_1425135600_1442415600_length_epoch_second_idx;


--
-- Name: submissions_1425135600_1442415600_length_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_execution_time_index ATTACH PARTITION public.submissions_1425135600_1442415600_length_execution_time_idx;


--
-- Name: submissions_1425135600_1442415600_length_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_point_index ATTACH PARTITION public.submissions_1425135600_1442415600_length_point_idx;


--
-- Name: submissions_1425135600_1442415600_point_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_epoch_second_index ATTACH PARTITION public.submissions_1425135600_1442415600_point_epoch_second_idx;


--
-- Name: submissions_1425135600_1442415600_point_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_execution_time_index ATTACH PARTITION public.submissions_1425135600_1442415600_point_execution_time_idx;


--
-- Name: submissions_1425135600_1442415600_point_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_length_index ATTACH PARTITION public.submissions_1425135600_1442415600_point_length_idx;


--
-- Name: submissions_1425135600_1442415600_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_1425135600_1442415600_problem_id_epoch_second_idx;


--
-- Name: submissions_1425135600_1442415600_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_1425135600_1442415600_problem_id_execution_time_idx;


--
-- Name: submissions_1425135600_1442415600_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_1425135600_1442415600_problem_id_length_idx;


--
-- Name: submissions_1425135600_1442415600_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_1425135600_1442415600_problem_id_point_idx;


--
-- Name: submissions_1425135600_1442415600_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_1425135600_1442415600_result_epoch_second_idx;


--
-- Name: submissions_1425135600_1442415600_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_1425135600_1442415600_result_execution_time_idx;


--
-- Name: submissions_1425135600_1442415600_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_1425135600_1442415600_result_length_idx;


--
-- Name: submissions_1425135600_1442415600_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_1425135600_1442415600_result_point_idx;


--
-- Name: submissions_1425135600_1442415600_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_1425135600_1442415600_user_id_epoch_second_idx;


--
-- Name: submissions_1425135600_1442415600_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_1425135600_1442415600_user_id_execution_time_idx;


--
-- Name: submissions_1425135600_1442415600_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_1425135600_1442415600_user_id_length_idx;


--
-- Name: submissions_1425135600_1442415600_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_1425135600_1442415600_user_id_point_idx;


--
-- Name: submissions_1425135600_14424156_execution_time_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_epoch_second_index ATTACH PARTITION public.submissions_1425135600_14424156_execution_time_epoch_second_idx;


--
-- Name: submissions_1442415600_1459695600_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_1442415600_1459695600_contest_id_epoch_second_idx;


--
-- Name: submissions_1442415600_1459695600_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_1442415600_1459695600_contest_id_execution_time_idx;


--
-- Name: submissions_1442415600_1459695600_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_1442415600_1459695600_contest_id_length_idx;


--
-- Name: submissions_1442415600_1459695600_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_1442415600_1459695600_contest_id_point_idx;


--
-- Name: submissions_1442415600_1459695600_epoch_second_updated_at_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_updated_at_index ATTACH PARTITION public.submissions_1442415600_1459695600_epoch_second_updated_at_idx;


--
-- Name: submissions_1442415600_1459695600_execution_time_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_length_index ATTACH PARTITION public.submissions_1442415600_1459695600_execution_time_length_idx;


--
-- Name: submissions_1442415600_1459695600_execution_time_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_point_index ATTACH PARTITION public.submissions_1442415600_1459695600_execution_time_point_idx;


--
-- Name: submissions_1442415600_1459695600_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_id_epoch_second_unique ATTACH PARTITION public.submissions_1442415600_1459695600_id_epoch_second_idx;


--
-- Name: submissions_1442415600_1459695600_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_1442415600_1459695600_language_epoch_second_idx;


--
-- Name: submissions_1442415600_1459695600_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_1442415600_1459695600_language_execution_time_idx;


--
-- Name: submissions_1442415600_1459695600_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_1442415600_1459695600_language_length_idx;


--
-- Name: submissions_1442415600_1459695600_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_1442415600_1459695600_language_point_idx;


--
-- Name: submissions_1442415600_1459695600_length_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_epoch_second_index ATTACH PARTITION public.submissions_1442415600_1459695600_length_epoch_second_idx;


--
-- Name: submissions_1442415600_1459695600_length_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_execution_time_index ATTACH PARTITION public.submissions_1442415600_1459695600_length_execution_time_idx;


--
-- Name: submissions_1442415600_1459695600_length_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_point_index ATTACH PARTITION public.submissions_1442415600_1459695600_length_point_idx;


--
-- Name: submissions_1442415600_1459695600_point_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_epoch_second_index ATTACH PARTITION public.submissions_1442415600_1459695600_point_epoch_second_idx;


--
-- Name: submissions_1442415600_1459695600_point_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_execution_time_index ATTACH PARTITION public.submissions_1442415600_1459695600_point_execution_time_idx;


--
-- Name: submissions_1442415600_1459695600_point_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_length_index ATTACH PARTITION public.submissions_1442415600_1459695600_point_length_idx;


--
-- Name: submissions_1442415600_1459695600_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_1442415600_1459695600_problem_id_epoch_second_idx;


--
-- Name: submissions_1442415600_1459695600_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_1442415600_1459695600_problem_id_execution_time_idx;


--
-- Name: submissions_1442415600_1459695600_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_1442415600_1459695600_problem_id_length_idx;


--
-- Name: submissions_1442415600_1459695600_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_1442415600_1459695600_problem_id_point_idx;


--
-- Name: submissions_1442415600_1459695600_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_1442415600_1459695600_result_epoch_second_idx;


--
-- Name: submissions_1442415600_1459695600_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_1442415600_1459695600_result_execution_time_idx;


--
-- Name: submissions_1442415600_1459695600_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_1442415600_1459695600_result_length_idx;


--
-- Name: submissions_1442415600_1459695600_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_1442415600_1459695600_result_point_idx;


--
-- Name: submissions_1442415600_1459695600_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_1442415600_1459695600_user_id_epoch_second_idx;


--
-- Name: submissions_1442415600_1459695600_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_1442415600_1459695600_user_id_execution_time_idx;


--
-- Name: submissions_1442415600_1459695600_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_1442415600_1459695600_user_id_length_idx;


--
-- Name: submissions_1442415600_1459695600_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_1442415600_1459695600_user_id_point_idx;


--
-- Name: submissions_1442415600_14596956_execution_time_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_epoch_second_index ATTACH PARTITION public.submissions_1442415600_14596956_execution_time_epoch_second_idx;


--
-- Name: submissions_1459695600_1476975600_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_1459695600_1476975600_contest_id_epoch_second_idx;


--
-- Name: submissions_1459695600_1476975600_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_1459695600_1476975600_contest_id_execution_time_idx;


--
-- Name: submissions_1459695600_1476975600_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_1459695600_1476975600_contest_id_length_idx;


--
-- Name: submissions_1459695600_1476975600_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_1459695600_1476975600_contest_id_point_idx;


--
-- Name: submissions_1459695600_1476975600_epoch_second_updated_at_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_updated_at_index ATTACH PARTITION public.submissions_1459695600_1476975600_epoch_second_updated_at_idx;


--
-- Name: submissions_1459695600_1476975600_execution_time_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_length_index ATTACH PARTITION public.submissions_1459695600_1476975600_execution_time_length_idx;


--
-- Name: submissions_1459695600_1476975600_execution_time_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_point_index ATTACH PARTITION public.submissions_1459695600_1476975600_execution_time_point_idx;


--
-- Name: submissions_1459695600_1476975600_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_id_epoch_second_unique ATTACH PARTITION public.submissions_1459695600_1476975600_id_epoch_second_idx;


--
-- Name: submissions_1459695600_1476975600_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_1459695600_1476975600_language_epoch_second_idx;


--
-- Name: submissions_1459695600_1476975600_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_1459695600_1476975600_language_execution_time_idx;


--
-- Name: submissions_1459695600_1476975600_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_1459695600_1476975600_language_length_idx;


--
-- Name: submissions_1459695600_1476975600_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_1459695600_1476975600_language_point_idx;


--
-- Name: submissions_1459695600_1476975600_length_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_epoch_second_index ATTACH PARTITION public.submissions_1459695600_1476975600_length_epoch_second_idx;


--
-- Name: submissions_1459695600_1476975600_length_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_execution_time_index ATTACH PARTITION public.submissions_1459695600_1476975600_length_execution_time_idx;


--
-- Name: submissions_1459695600_1476975600_length_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_point_index ATTACH PARTITION public.submissions_1459695600_1476975600_length_point_idx;


--
-- Name: submissions_1459695600_1476975600_point_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_epoch_second_index ATTACH PARTITION public.submissions_1459695600_1476975600_point_epoch_second_idx;


--
-- Name: submissions_1459695600_1476975600_point_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_execution_time_index ATTACH PARTITION public.submissions_1459695600_1476975600_point_execution_time_idx;


--
-- Name: submissions_1459695600_1476975600_point_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_length_index ATTACH PARTITION public.submissions_1459695600_1476975600_point_length_idx;


--
-- Name: submissions_1459695600_1476975600_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_1459695600_1476975600_problem_id_epoch_second_idx;


--
-- Name: submissions_1459695600_1476975600_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_1459695600_1476975600_problem_id_execution_time_idx;


--
-- Name: submissions_1459695600_1476975600_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_1459695600_1476975600_problem_id_length_idx;


--
-- Name: submissions_1459695600_1476975600_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_1459695600_1476975600_problem_id_point_idx;


--
-- Name: submissions_1459695600_1476975600_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_1459695600_1476975600_result_epoch_second_idx;


--
-- Name: submissions_1459695600_1476975600_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_1459695600_1476975600_result_execution_time_idx;


--
-- Name: submissions_1459695600_1476975600_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_1459695600_1476975600_result_length_idx;


--
-- Name: submissions_1459695600_1476975600_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_1459695600_1476975600_result_point_idx;


--
-- Name: submissions_1459695600_1476975600_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_1459695600_1476975600_user_id_epoch_second_idx;


--
-- Name: submissions_1459695600_1476975600_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_1459695600_1476975600_user_id_execution_time_idx;


--
-- Name: submissions_1459695600_1476975600_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_1459695600_1476975600_user_id_length_idx;


--
-- Name: submissions_1459695600_1476975600_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_1459695600_1476975600_user_id_point_idx;


--
-- Name: submissions_1459695600_14769756_execution_time_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_epoch_second_index ATTACH PARTITION public.submissions_1459695600_14769756_execution_time_epoch_second_idx;


--
-- Name: submissions_1476975600_1494255600_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_1476975600_1494255600_contest_id_epoch_second_idx;


--
-- Name: submissions_1476975600_1494255600_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_1476975600_1494255600_contest_id_execution_time_idx;


--
-- Name: submissions_1476975600_1494255600_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_1476975600_1494255600_contest_id_length_idx;


--
-- Name: submissions_1476975600_1494255600_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_1476975600_1494255600_contest_id_point_idx;


--
-- Name: submissions_1476975600_1494255600_epoch_second_updated_at_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_updated_at_index ATTACH PARTITION public.submissions_1476975600_1494255600_epoch_second_updated_at_idx;


--
-- Name: submissions_1476975600_1494255600_execution_time_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_length_index ATTACH PARTITION public.submissions_1476975600_1494255600_execution_time_length_idx;


--
-- Name: submissions_1476975600_1494255600_execution_time_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_point_index ATTACH PARTITION public.submissions_1476975600_1494255600_execution_time_point_idx;


--
-- Name: submissions_1476975600_1494255600_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_id_epoch_second_unique ATTACH PARTITION public.submissions_1476975600_1494255600_id_epoch_second_idx;


--
-- Name: submissions_1476975600_1494255600_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_1476975600_1494255600_language_epoch_second_idx;


--
-- Name: submissions_1476975600_1494255600_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_1476975600_1494255600_language_execution_time_idx;


--
-- Name: submissions_1476975600_1494255600_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_1476975600_1494255600_language_length_idx;


--
-- Name: submissions_1476975600_1494255600_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_1476975600_1494255600_language_point_idx;


--
-- Name: submissions_1476975600_1494255600_length_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_epoch_second_index ATTACH PARTITION public.submissions_1476975600_1494255600_length_epoch_second_idx;


--
-- Name: submissions_1476975600_1494255600_length_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_execution_time_index ATTACH PARTITION public.submissions_1476975600_1494255600_length_execution_time_idx;


--
-- Name: submissions_1476975600_1494255600_length_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_point_index ATTACH PARTITION public.submissions_1476975600_1494255600_length_point_idx;


--
-- Name: submissions_1476975600_1494255600_point_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_epoch_second_index ATTACH PARTITION public.submissions_1476975600_1494255600_point_epoch_second_idx;


--
-- Name: submissions_1476975600_1494255600_point_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_execution_time_index ATTACH PARTITION public.submissions_1476975600_1494255600_point_execution_time_idx;


--
-- Name: submissions_1476975600_1494255600_point_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_length_index ATTACH PARTITION public.submissions_1476975600_1494255600_point_length_idx;


--
-- Name: submissions_1476975600_1494255600_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_1476975600_1494255600_problem_id_epoch_second_idx;


--
-- Name: submissions_1476975600_1494255600_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_1476975600_1494255600_problem_id_execution_time_idx;


--
-- Name: submissions_1476975600_1494255600_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_1476975600_1494255600_problem_id_length_idx;


--
-- Name: submissions_1476975600_1494255600_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_1476975600_1494255600_problem_id_point_idx;


--
-- Name: submissions_1476975600_1494255600_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_1476975600_1494255600_result_epoch_second_idx;


--
-- Name: submissions_1476975600_1494255600_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_1476975600_1494255600_result_execution_time_idx;


--
-- Name: submissions_1476975600_1494255600_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_1476975600_1494255600_result_length_idx;


--
-- Name: submissions_1476975600_1494255600_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_1476975600_1494255600_result_point_idx;


--
-- Name: submissions_1476975600_1494255600_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_1476975600_1494255600_user_id_epoch_second_idx;


--
-- Name: submissions_1476975600_1494255600_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_1476975600_1494255600_user_id_execution_time_idx;


--
-- Name: submissions_1476975600_1494255600_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_1476975600_1494255600_user_id_length_idx;


--
-- Name: submissions_1476975600_1494255600_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_1476975600_1494255600_user_id_point_idx;


--
-- Name: submissions_1476975600_14942556_execution_time_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_epoch_second_index ATTACH PARTITION public.submissions_1476975600_14942556_execution_time_epoch_second_idx;


--
-- Name: submissions_1494255600_1511535600_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_1494255600_1511535600_contest_id_epoch_second_idx;


--
-- Name: submissions_1494255600_1511535600_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_1494255600_1511535600_contest_id_execution_time_idx;


--
-- Name: submissions_1494255600_1511535600_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_1494255600_1511535600_contest_id_length_idx;


--
-- Name: submissions_1494255600_1511535600_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_1494255600_1511535600_contest_id_point_idx;


--
-- Name: submissions_1494255600_1511535600_epoch_second_updated_at_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_updated_at_index ATTACH PARTITION public.submissions_1494255600_1511535600_epoch_second_updated_at_idx;


--
-- Name: submissions_1494255600_1511535600_execution_time_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_length_index ATTACH PARTITION public.submissions_1494255600_1511535600_execution_time_length_idx;


--
-- Name: submissions_1494255600_1511535600_execution_time_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_point_index ATTACH PARTITION public.submissions_1494255600_1511535600_execution_time_point_idx;


--
-- Name: submissions_1494255600_1511535600_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_id_epoch_second_unique ATTACH PARTITION public.submissions_1494255600_1511535600_id_epoch_second_idx;


--
-- Name: submissions_1494255600_1511535600_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_1494255600_1511535600_language_epoch_second_idx;


--
-- Name: submissions_1494255600_1511535600_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_1494255600_1511535600_language_execution_time_idx;


--
-- Name: submissions_1494255600_1511535600_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_1494255600_1511535600_language_length_idx;


--
-- Name: submissions_1494255600_1511535600_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_1494255600_1511535600_language_point_idx;


--
-- Name: submissions_1494255600_1511535600_length_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_epoch_second_index ATTACH PARTITION public.submissions_1494255600_1511535600_length_epoch_second_idx;


--
-- Name: submissions_1494255600_1511535600_length_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_execution_time_index ATTACH PARTITION public.submissions_1494255600_1511535600_length_execution_time_idx;


--
-- Name: submissions_1494255600_1511535600_length_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_point_index ATTACH PARTITION public.submissions_1494255600_1511535600_length_point_idx;


--
-- Name: submissions_1494255600_1511535600_point_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_epoch_second_index ATTACH PARTITION public.submissions_1494255600_1511535600_point_epoch_second_idx;


--
-- Name: submissions_1494255600_1511535600_point_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_execution_time_index ATTACH PARTITION public.submissions_1494255600_1511535600_point_execution_time_idx;


--
-- Name: submissions_1494255600_1511535600_point_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_length_index ATTACH PARTITION public.submissions_1494255600_1511535600_point_length_idx;


--
-- Name: submissions_1494255600_1511535600_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_1494255600_1511535600_problem_id_epoch_second_idx;


--
-- Name: submissions_1494255600_1511535600_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_1494255600_1511535600_problem_id_execution_time_idx;


--
-- Name: submissions_1494255600_1511535600_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_1494255600_1511535600_problem_id_length_idx;


--
-- Name: submissions_1494255600_1511535600_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_1494255600_1511535600_problem_id_point_idx;


--
-- Name: submissions_1494255600_1511535600_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_1494255600_1511535600_result_epoch_second_idx;


--
-- Name: submissions_1494255600_1511535600_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_1494255600_1511535600_result_execution_time_idx;


--
-- Name: submissions_1494255600_1511535600_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_1494255600_1511535600_result_length_idx;


--
-- Name: submissions_1494255600_1511535600_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_1494255600_1511535600_result_point_idx;


--
-- Name: submissions_1494255600_1511535600_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_1494255600_1511535600_user_id_epoch_second_idx;


--
-- Name: submissions_1494255600_1511535600_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_1494255600_1511535600_user_id_execution_time_idx;


--
-- Name: submissions_1494255600_1511535600_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_1494255600_1511535600_user_id_length_idx;


--
-- Name: submissions_1494255600_1511535600_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_1494255600_1511535600_user_id_point_idx;


--
-- Name: submissions_1494255600_15115356_execution_time_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_epoch_second_index ATTACH PARTITION public.submissions_1494255600_15115356_execution_time_epoch_second_idx;


--
-- Name: submissions_1511535600_1528815600_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_1511535600_1528815600_contest_id_epoch_second_idx;


--
-- Name: submissions_1511535600_1528815600_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_1511535600_1528815600_contest_id_execution_time_idx;


--
-- Name: submissions_1511535600_1528815600_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_1511535600_1528815600_contest_id_length_idx;


--
-- Name: submissions_1511535600_1528815600_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_1511535600_1528815600_contest_id_point_idx;


--
-- Name: submissions_1511535600_1528815600_epoch_second_updated_at_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_updated_at_index ATTACH PARTITION public.submissions_1511535600_1528815600_epoch_second_updated_at_idx;


--
-- Name: submissions_1511535600_1528815600_execution_time_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_length_index ATTACH PARTITION public.submissions_1511535600_1528815600_execution_time_length_idx;


--
-- Name: submissions_1511535600_1528815600_execution_time_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_point_index ATTACH PARTITION public.submissions_1511535600_1528815600_execution_time_point_idx;


--
-- Name: submissions_1511535600_1528815600_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_id_epoch_second_unique ATTACH PARTITION public.submissions_1511535600_1528815600_id_epoch_second_idx;


--
-- Name: submissions_1511535600_1528815600_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_1511535600_1528815600_language_epoch_second_idx;


--
-- Name: submissions_1511535600_1528815600_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_1511535600_1528815600_language_execution_time_idx;


--
-- Name: submissions_1511535600_1528815600_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_1511535600_1528815600_language_length_idx;


--
-- Name: submissions_1511535600_1528815600_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_1511535600_1528815600_language_point_idx;


--
-- Name: submissions_1511535600_1528815600_length_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_epoch_second_index ATTACH PARTITION public.submissions_1511535600_1528815600_length_epoch_second_idx;


--
-- Name: submissions_1511535600_1528815600_length_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_execution_time_index ATTACH PARTITION public.submissions_1511535600_1528815600_length_execution_time_idx;


--
-- Name: submissions_1511535600_1528815600_length_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_point_index ATTACH PARTITION public.submissions_1511535600_1528815600_length_point_idx;


--
-- Name: submissions_1511535600_1528815600_point_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_epoch_second_index ATTACH PARTITION public.submissions_1511535600_1528815600_point_epoch_second_idx;


--
-- Name: submissions_1511535600_1528815600_point_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_execution_time_index ATTACH PARTITION public.submissions_1511535600_1528815600_point_execution_time_idx;


--
-- Name: submissions_1511535600_1528815600_point_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_length_index ATTACH PARTITION public.submissions_1511535600_1528815600_point_length_idx;


--
-- Name: submissions_1511535600_1528815600_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_1511535600_1528815600_problem_id_epoch_second_idx;


--
-- Name: submissions_1511535600_1528815600_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_1511535600_1528815600_problem_id_execution_time_idx;


--
-- Name: submissions_1511535600_1528815600_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_1511535600_1528815600_problem_id_length_idx;


--
-- Name: submissions_1511535600_1528815600_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_1511535600_1528815600_problem_id_point_idx;


--
-- Name: submissions_1511535600_1528815600_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_1511535600_1528815600_result_epoch_second_idx;


--
-- Name: submissions_1511535600_1528815600_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_1511535600_1528815600_result_execution_time_idx;


--
-- Name: submissions_1511535600_1528815600_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_1511535600_1528815600_result_length_idx;


--
-- Name: submissions_1511535600_1528815600_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_1511535600_1528815600_result_point_idx;


--
-- Name: submissions_1511535600_1528815600_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_1511535600_1528815600_user_id_epoch_second_idx;


--
-- Name: submissions_1511535600_1528815600_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_1511535600_1528815600_user_id_execution_time_idx;


--
-- Name: submissions_1511535600_1528815600_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_1511535600_1528815600_user_id_length_idx;


--
-- Name: submissions_1511535600_1528815600_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_1511535600_1528815600_user_id_point_idx;


--
-- Name: submissions_1511535600_15288156_execution_time_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_epoch_second_index ATTACH PARTITION public.submissions_1511535600_15288156_execution_time_epoch_second_idx;


--
-- Name: submissions_1528815600_1546095600_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_1528815600_1546095600_contest_id_epoch_second_idx;


--
-- Name: submissions_1528815600_1546095600_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_1528815600_1546095600_contest_id_execution_time_idx;


--
-- Name: submissions_1528815600_1546095600_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_1528815600_1546095600_contest_id_length_idx;


--
-- Name: submissions_1528815600_1546095600_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_1528815600_1546095600_contest_id_point_idx;


--
-- Name: submissions_1528815600_1546095600_epoch_second_updated_at_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_updated_at_index ATTACH PARTITION public.submissions_1528815600_1546095600_epoch_second_updated_at_idx;


--
-- Name: submissions_1528815600_1546095600_execution_time_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_length_index ATTACH PARTITION public.submissions_1528815600_1546095600_execution_time_length_idx;


--
-- Name: submissions_1528815600_1546095600_execution_time_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_point_index ATTACH PARTITION public.submissions_1528815600_1546095600_execution_time_point_idx;


--
-- Name: submissions_1528815600_1546095600_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_id_epoch_second_unique ATTACH PARTITION public.submissions_1528815600_1546095600_id_epoch_second_idx;


--
-- Name: submissions_1528815600_1546095600_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_1528815600_1546095600_language_epoch_second_idx;


--
-- Name: submissions_1528815600_1546095600_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_1528815600_1546095600_language_execution_time_idx;


--
-- Name: submissions_1528815600_1546095600_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_1528815600_1546095600_language_length_idx;


--
-- Name: submissions_1528815600_1546095600_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_1528815600_1546095600_language_point_idx;


--
-- Name: submissions_1528815600_1546095600_length_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_epoch_second_index ATTACH PARTITION public.submissions_1528815600_1546095600_length_epoch_second_idx;


--
-- Name: submissions_1528815600_1546095600_length_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_execution_time_index ATTACH PARTITION public.submissions_1528815600_1546095600_length_execution_time_idx;


--
-- Name: submissions_1528815600_1546095600_length_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_point_index ATTACH PARTITION public.submissions_1528815600_1546095600_length_point_idx;


--
-- Name: submissions_1528815600_1546095600_point_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_epoch_second_index ATTACH PARTITION public.submissions_1528815600_1546095600_point_epoch_second_idx;


--
-- Name: submissions_1528815600_1546095600_point_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_execution_time_index ATTACH PARTITION public.submissions_1528815600_1546095600_point_execution_time_idx;


--
-- Name: submissions_1528815600_1546095600_point_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_length_index ATTACH PARTITION public.submissions_1528815600_1546095600_point_length_idx;


--
-- Name: submissions_1528815600_1546095600_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_1528815600_1546095600_problem_id_epoch_second_idx;


--
-- Name: submissions_1528815600_1546095600_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_1528815600_1546095600_problem_id_execution_time_idx;


--
-- Name: submissions_1528815600_1546095600_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_1528815600_1546095600_problem_id_length_idx;


--
-- Name: submissions_1528815600_1546095600_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_1528815600_1546095600_problem_id_point_idx;


--
-- Name: submissions_1528815600_1546095600_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_1528815600_1546095600_result_epoch_second_idx;


--
-- Name: submissions_1528815600_1546095600_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_1528815600_1546095600_result_execution_time_idx;


--
-- Name: submissions_1528815600_1546095600_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_1528815600_1546095600_result_length_idx;


--
-- Name: submissions_1528815600_1546095600_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_1528815600_1546095600_result_point_idx;


--
-- Name: submissions_1528815600_1546095600_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_1528815600_1546095600_user_id_epoch_second_idx;


--
-- Name: submissions_1528815600_1546095600_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_1528815600_1546095600_user_id_execution_time_idx;


--
-- Name: submissions_1528815600_1546095600_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_1528815600_1546095600_user_id_length_idx;


--
-- Name: submissions_1528815600_1546095600_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_1528815600_1546095600_user_id_point_idx;


--
-- Name: submissions_1528815600_15460956_execution_time_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_epoch_second_index ATTACH PARTITION public.submissions_1528815600_15460956_execution_time_epoch_second_idx;


--
-- Name: submissions_1546095600_1563375600_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_1546095600_1563375600_contest_id_epoch_second_idx;


--
-- Name: submissions_1546095600_1563375600_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_1546095600_1563375600_contest_id_execution_time_idx;


--
-- Name: submissions_1546095600_1563375600_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_1546095600_1563375600_contest_id_length_idx;


--
-- Name: submissions_1546095600_1563375600_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_1546095600_1563375600_contest_id_point_idx;


--
-- Name: submissions_1546095600_1563375600_epoch_second_updated_at_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_updated_at_index ATTACH PARTITION public.submissions_1546095600_1563375600_epoch_second_updated_at_idx;


--
-- Name: submissions_1546095600_1563375600_execution_time_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_length_index ATTACH PARTITION public.submissions_1546095600_1563375600_execution_time_length_idx;


--
-- Name: submissions_1546095600_1563375600_execution_time_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_point_index ATTACH PARTITION public.submissions_1546095600_1563375600_execution_time_point_idx;


--
-- Name: submissions_1546095600_1563375600_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_id_epoch_second_unique ATTACH PARTITION public.submissions_1546095600_1563375600_id_epoch_second_idx;


--
-- Name: submissions_1546095600_1563375600_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_1546095600_1563375600_language_epoch_second_idx;


--
-- Name: submissions_1546095600_1563375600_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_1546095600_1563375600_language_execution_time_idx;


--
-- Name: submissions_1546095600_1563375600_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_1546095600_1563375600_language_length_idx;


--
-- Name: submissions_1546095600_1563375600_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_1546095600_1563375600_language_point_idx;


--
-- Name: submissions_1546095600_1563375600_length_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_epoch_second_index ATTACH PARTITION public.submissions_1546095600_1563375600_length_epoch_second_idx;


--
-- Name: submissions_1546095600_1563375600_length_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_execution_time_index ATTACH PARTITION public.submissions_1546095600_1563375600_length_execution_time_idx;


--
-- Name: submissions_1546095600_1563375600_length_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_point_index ATTACH PARTITION public.submissions_1546095600_1563375600_length_point_idx;


--
-- Name: submissions_1546095600_1563375600_point_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_epoch_second_index ATTACH PARTITION public.submissions_1546095600_1563375600_point_epoch_second_idx;


--
-- Name: submissions_1546095600_1563375600_point_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_execution_time_index ATTACH PARTITION public.submissions_1546095600_1563375600_point_execution_time_idx;


--
-- Name: submissions_1546095600_1563375600_point_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_length_index ATTACH PARTITION public.submissions_1546095600_1563375600_point_length_idx;


--
-- Name: submissions_1546095600_1563375600_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_1546095600_1563375600_problem_id_epoch_second_idx;


--
-- Name: submissions_1546095600_1563375600_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_1546095600_1563375600_problem_id_execution_time_idx;


--
-- Name: submissions_1546095600_1563375600_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_1546095600_1563375600_problem_id_length_idx;


--
-- Name: submissions_1546095600_1563375600_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_1546095600_1563375600_problem_id_point_idx;


--
-- Name: submissions_1546095600_1563375600_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_1546095600_1563375600_result_epoch_second_idx;


--
-- Name: submissions_1546095600_1563375600_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_1546095600_1563375600_result_execution_time_idx;


--
-- Name: submissions_1546095600_1563375600_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_1546095600_1563375600_result_length_idx;


--
-- Name: submissions_1546095600_1563375600_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_1546095600_1563375600_result_point_idx;


--
-- Name: submissions_1546095600_1563375600_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_1546095600_1563375600_user_id_epoch_second_idx;


--
-- Name: submissions_1546095600_1563375600_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_1546095600_1563375600_user_id_execution_time_idx;


--
-- Name: submissions_1546095600_1563375600_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_1546095600_1563375600_user_id_length_idx;


--
-- Name: submissions_1546095600_1563375600_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_1546095600_1563375600_user_id_point_idx;


--
-- Name: submissions_1546095600_15633756_execution_time_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_epoch_second_index ATTACH PARTITION public.submissions_1546095600_15633756_execution_time_epoch_second_idx;


--
-- Name: submissions_1563375600_1580655600_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_1563375600_1580655600_contest_id_epoch_second_idx;


--
-- Name: submissions_1563375600_1580655600_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_1563375600_1580655600_contest_id_execution_time_idx;


--
-- Name: submissions_1563375600_1580655600_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_1563375600_1580655600_contest_id_length_idx;


--
-- Name: submissions_1563375600_1580655600_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_1563375600_1580655600_contest_id_point_idx;


--
-- Name: submissions_1563375600_1580655600_epoch_second_updated_at_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_updated_at_index ATTACH PARTITION public.submissions_1563375600_1580655600_epoch_second_updated_at_idx;


--
-- Name: submissions_1563375600_1580655600_execution_time_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_length_index ATTACH PARTITION public.submissions_1563375600_1580655600_execution_time_length_idx;


--
-- Name: submissions_1563375600_1580655600_execution_time_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_point_index ATTACH PARTITION public.submissions_1563375600_1580655600_execution_time_point_idx;


--
-- Name: submissions_1563375600_1580655600_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_id_epoch_second_unique ATTACH PARTITION public.submissions_1563375600_1580655600_id_epoch_second_idx;


--
-- Name: submissions_1563375600_1580655600_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_1563375600_1580655600_language_epoch_second_idx;


--
-- Name: submissions_1563375600_1580655600_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_1563375600_1580655600_language_execution_time_idx;


--
-- Name: submissions_1563375600_1580655600_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_1563375600_1580655600_language_length_idx;


--
-- Name: submissions_1563375600_1580655600_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_1563375600_1580655600_language_point_idx;


--
-- Name: submissions_1563375600_1580655600_length_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_epoch_second_index ATTACH PARTITION public.submissions_1563375600_1580655600_length_epoch_second_idx;


--
-- Name: submissions_1563375600_1580655600_length_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_execution_time_index ATTACH PARTITION public.submissions_1563375600_1580655600_length_execution_time_idx;


--
-- Name: submissions_1563375600_1580655600_length_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_point_index ATTACH PARTITION public.submissions_1563375600_1580655600_length_point_idx;


--
-- Name: submissions_1563375600_1580655600_point_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_epoch_second_index ATTACH PARTITION public.submissions_1563375600_1580655600_point_epoch_second_idx;


--
-- Name: submissions_1563375600_1580655600_point_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_execution_time_index ATTACH PARTITION public.submissions_1563375600_1580655600_point_execution_time_idx;


--
-- Name: submissions_1563375600_1580655600_point_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_length_index ATTACH PARTITION public.submissions_1563375600_1580655600_point_length_idx;


--
-- Name: submissions_1563375600_1580655600_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_1563375600_1580655600_problem_id_epoch_second_idx;


--
-- Name: submissions_1563375600_1580655600_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_1563375600_1580655600_problem_id_execution_time_idx;


--
-- Name: submissions_1563375600_1580655600_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_1563375600_1580655600_problem_id_length_idx;


--
-- Name: submissions_1563375600_1580655600_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_1563375600_1580655600_problem_id_point_idx;


--
-- Name: submissions_1563375600_1580655600_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_1563375600_1580655600_result_epoch_second_idx;


--
-- Name: submissions_1563375600_1580655600_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_1563375600_1580655600_result_execution_time_idx;


--
-- Name: submissions_1563375600_1580655600_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_1563375600_1580655600_result_length_idx;


--
-- Name: submissions_1563375600_1580655600_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_1563375600_1580655600_result_point_idx;


--
-- Name: submissions_1563375600_1580655600_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_1563375600_1580655600_user_id_epoch_second_idx;


--
-- Name: submissions_1563375600_1580655600_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_1563375600_1580655600_user_id_execution_time_idx;


--
-- Name: submissions_1563375600_1580655600_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_1563375600_1580655600_user_id_length_idx;


--
-- Name: submissions_1563375600_1580655600_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_1563375600_1580655600_user_id_point_idx;


--
-- Name: submissions_1563375600_15806556_execution_time_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_epoch_second_index ATTACH PARTITION public.submissions_1563375600_15806556_execution_time_epoch_second_idx;


--
-- Name: submissions_1580655600_1597935600_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_1580655600_1597935600_contest_id_epoch_second_idx;


--
-- Name: submissions_1580655600_1597935600_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_1580655600_1597935600_contest_id_execution_time_idx;


--
-- Name: submissions_1580655600_1597935600_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_1580655600_1597935600_contest_id_length_idx;


--
-- Name: submissions_1580655600_1597935600_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_1580655600_1597935600_contest_id_point_idx;


--
-- Name: submissions_1580655600_1597935600_epoch_second_updated_at_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_updated_at_index ATTACH PARTITION public.submissions_1580655600_1597935600_epoch_second_updated_at_idx;


--
-- Name: submissions_1580655600_1597935600_execution_time_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_length_index ATTACH PARTITION public.submissions_1580655600_1597935600_execution_time_length_idx;


--
-- Name: submissions_1580655600_1597935600_execution_time_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_point_index ATTACH PARTITION public.submissions_1580655600_1597935600_execution_time_point_idx;


--
-- Name: submissions_1580655600_1597935600_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_id_epoch_second_unique ATTACH PARTITION public.submissions_1580655600_1597935600_id_epoch_second_idx;


--
-- Name: submissions_1580655600_1597935600_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_1580655600_1597935600_language_epoch_second_idx;


--
-- Name: submissions_1580655600_1597935600_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_1580655600_1597935600_language_execution_time_idx;


--
-- Name: submissions_1580655600_1597935600_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_1580655600_1597935600_language_length_idx;


--
-- Name: submissions_1580655600_1597935600_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_1580655600_1597935600_language_point_idx;


--
-- Name: submissions_1580655600_1597935600_length_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_epoch_second_index ATTACH PARTITION public.submissions_1580655600_1597935600_length_epoch_second_idx;


--
-- Name: submissions_1580655600_1597935600_length_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_execution_time_index ATTACH PARTITION public.submissions_1580655600_1597935600_length_execution_time_idx;


--
-- Name: submissions_1580655600_1597935600_length_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_point_index ATTACH PARTITION public.submissions_1580655600_1597935600_length_point_idx;


--
-- Name: submissions_1580655600_1597935600_point_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_epoch_second_index ATTACH PARTITION public.submissions_1580655600_1597935600_point_epoch_second_idx;


--
-- Name: submissions_1580655600_1597935600_point_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_execution_time_index ATTACH PARTITION public.submissions_1580655600_1597935600_point_execution_time_idx;


--
-- Name: submissions_1580655600_1597935600_point_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_length_index ATTACH PARTITION public.submissions_1580655600_1597935600_point_length_idx;


--
-- Name: submissions_1580655600_1597935600_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_1580655600_1597935600_problem_id_epoch_second_idx;


--
-- Name: submissions_1580655600_1597935600_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_1580655600_1597935600_problem_id_execution_time_idx;


--
-- Name: submissions_1580655600_1597935600_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_1580655600_1597935600_problem_id_length_idx;


--
-- Name: submissions_1580655600_1597935600_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_1580655600_1597935600_problem_id_point_idx;


--
-- Name: submissions_1580655600_1597935600_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_1580655600_1597935600_result_epoch_second_idx;


--
-- Name: submissions_1580655600_1597935600_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_1580655600_1597935600_result_execution_time_idx;


--
-- Name: submissions_1580655600_1597935600_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_1580655600_1597935600_result_length_idx;


--
-- Name: submissions_1580655600_1597935600_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_1580655600_1597935600_result_point_idx;


--
-- Name: submissions_1580655600_1597935600_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_1580655600_1597935600_user_id_epoch_second_idx;


--
-- Name: submissions_1580655600_1597935600_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_1580655600_1597935600_user_id_execution_time_idx;


--
-- Name: submissions_1580655600_1597935600_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_1580655600_1597935600_user_id_length_idx;


--
-- Name: submissions_1580655600_1597935600_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_1580655600_1597935600_user_id_point_idx;


--
-- Name: submissions_1580655600_15979356_execution_time_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_epoch_second_index ATTACH PARTITION public.submissions_1580655600_15979356_execution_time_epoch_second_idx;


--
-- Name: submissions_1597935600_1615215600_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_1597935600_1615215600_contest_id_epoch_second_idx;


--
-- Name: submissions_1597935600_1615215600_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_1597935600_1615215600_contest_id_execution_time_idx;


--
-- Name: submissions_1597935600_1615215600_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_1597935600_1615215600_contest_id_length_idx;


--
-- Name: submissions_1597935600_1615215600_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_1597935600_1615215600_contest_id_point_idx;


--
-- Name: submissions_1597935600_1615215600_epoch_second_updated_at_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_updated_at_index ATTACH PARTITION public.submissions_1597935600_1615215600_epoch_second_updated_at_idx;


--
-- Name: submissions_1597935600_1615215600_execution_time_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_length_index ATTACH PARTITION public.submissions_1597935600_1615215600_execution_time_length_idx;


--
-- Name: submissions_1597935600_1615215600_execution_time_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_point_index ATTACH PARTITION public.submissions_1597935600_1615215600_execution_time_point_idx;


--
-- Name: submissions_1597935600_1615215600_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_id_epoch_second_unique ATTACH PARTITION public.submissions_1597935600_1615215600_id_epoch_second_idx;


--
-- Name: submissions_1597935600_1615215600_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_1597935600_1615215600_language_epoch_second_idx;


--
-- Name: submissions_1597935600_1615215600_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_1597935600_1615215600_language_execution_time_idx;


--
-- Name: submissions_1597935600_1615215600_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_1597935600_1615215600_language_length_idx;


--
-- Name: submissions_1597935600_1615215600_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_1597935600_1615215600_language_point_idx;


--
-- Name: submissions_1597935600_1615215600_length_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_epoch_second_index ATTACH PARTITION public.submissions_1597935600_1615215600_length_epoch_second_idx;


--
-- Name: submissions_1597935600_1615215600_length_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_execution_time_index ATTACH PARTITION public.submissions_1597935600_1615215600_length_execution_time_idx;


--
-- Name: submissions_1597935600_1615215600_length_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_point_index ATTACH PARTITION public.submissions_1597935600_1615215600_length_point_idx;


--
-- Name: submissions_1597935600_1615215600_point_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_epoch_second_index ATTACH PARTITION public.submissions_1597935600_1615215600_point_epoch_second_idx;


--
-- Name: submissions_1597935600_1615215600_point_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_execution_time_index ATTACH PARTITION public.submissions_1597935600_1615215600_point_execution_time_idx;


--
-- Name: submissions_1597935600_1615215600_point_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_length_index ATTACH PARTITION public.submissions_1597935600_1615215600_point_length_idx;


--
-- Name: submissions_1597935600_1615215600_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_1597935600_1615215600_problem_id_epoch_second_idx;


--
-- Name: submissions_1597935600_1615215600_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_1597935600_1615215600_problem_id_execution_time_idx;


--
-- Name: submissions_1597935600_1615215600_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_1597935600_1615215600_problem_id_length_idx;


--
-- Name: submissions_1597935600_1615215600_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_1597935600_1615215600_problem_id_point_idx;


--
-- Name: submissions_1597935600_1615215600_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_1597935600_1615215600_result_epoch_second_idx;


--
-- Name: submissions_1597935600_1615215600_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_1597935600_1615215600_result_execution_time_idx;


--
-- Name: submissions_1597935600_1615215600_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_1597935600_1615215600_result_length_idx;


--
-- Name: submissions_1597935600_1615215600_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_1597935600_1615215600_result_point_idx;


--
-- Name: submissions_1597935600_1615215600_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_1597935600_1615215600_user_id_epoch_second_idx;


--
-- Name: submissions_1597935600_1615215600_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_1597935600_1615215600_user_id_execution_time_idx;


--
-- Name: submissions_1597935600_1615215600_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_1597935600_1615215600_user_id_length_idx;


--
-- Name: submissions_1597935600_1615215600_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_1597935600_1615215600_user_id_point_idx;


--
-- Name: submissions_1597935600_16152156_execution_time_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_epoch_second_index ATTACH PARTITION public.submissions_1597935600_16152156_execution_time_epoch_second_idx;


--
-- Name: submissions_1615215600_1632495600_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_1615215600_1632495600_contest_id_epoch_second_idx;


--
-- Name: submissions_1615215600_1632495600_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_1615215600_1632495600_contest_id_execution_time_idx;


--
-- Name: submissions_1615215600_1632495600_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_1615215600_1632495600_contest_id_length_idx;


--
-- Name: submissions_1615215600_1632495600_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_1615215600_1632495600_contest_id_point_idx;


--
-- Name: submissions_1615215600_1632495600_epoch_second_updated_at_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_updated_at_index ATTACH PARTITION public.submissions_1615215600_1632495600_epoch_second_updated_at_idx;


--
-- Name: submissions_1615215600_1632495600_execution_time_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_length_index ATTACH PARTITION public.submissions_1615215600_1632495600_execution_time_length_idx;


--
-- Name: submissions_1615215600_1632495600_execution_time_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_point_index ATTACH PARTITION public.submissions_1615215600_1632495600_execution_time_point_idx;


--
-- Name: submissions_1615215600_1632495600_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_id_epoch_second_unique ATTACH PARTITION public.submissions_1615215600_1632495600_id_epoch_second_idx;


--
-- Name: submissions_1615215600_1632495600_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_1615215600_1632495600_language_epoch_second_idx;


--
-- Name: submissions_1615215600_1632495600_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_1615215600_1632495600_language_execution_time_idx;


--
-- Name: submissions_1615215600_1632495600_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_1615215600_1632495600_language_length_idx;


--
-- Name: submissions_1615215600_1632495600_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_1615215600_1632495600_language_point_idx;


--
-- Name: submissions_1615215600_1632495600_length_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_epoch_second_index ATTACH PARTITION public.submissions_1615215600_1632495600_length_epoch_second_idx;


--
-- Name: submissions_1615215600_1632495600_length_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_execution_time_index ATTACH PARTITION public.submissions_1615215600_1632495600_length_execution_time_idx;


--
-- Name: submissions_1615215600_1632495600_length_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_point_index ATTACH PARTITION public.submissions_1615215600_1632495600_length_point_idx;


--
-- Name: submissions_1615215600_1632495600_point_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_epoch_second_index ATTACH PARTITION public.submissions_1615215600_1632495600_point_epoch_second_idx;


--
-- Name: submissions_1615215600_1632495600_point_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_execution_time_index ATTACH PARTITION public.submissions_1615215600_1632495600_point_execution_time_idx;


--
-- Name: submissions_1615215600_1632495600_point_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_length_index ATTACH PARTITION public.submissions_1615215600_1632495600_point_length_idx;


--
-- Name: submissions_1615215600_1632495600_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_1615215600_1632495600_problem_id_epoch_second_idx;


--
-- Name: submissions_1615215600_1632495600_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_1615215600_1632495600_problem_id_execution_time_idx;


--
-- Name: submissions_1615215600_1632495600_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_1615215600_1632495600_problem_id_length_idx;


--
-- Name: submissions_1615215600_1632495600_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_1615215600_1632495600_problem_id_point_idx;


--
-- Name: submissions_1615215600_1632495600_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_1615215600_1632495600_result_epoch_second_idx;


--
-- Name: submissions_1615215600_1632495600_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_1615215600_1632495600_result_execution_time_idx;


--
-- Name: submissions_1615215600_1632495600_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_1615215600_1632495600_result_length_idx;


--
-- Name: submissions_1615215600_1632495600_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_1615215600_1632495600_result_point_idx;


--
-- Name: submissions_1615215600_1632495600_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_1615215600_1632495600_user_id_epoch_second_idx;


--
-- Name: submissions_1615215600_1632495600_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_1615215600_1632495600_user_id_execution_time_idx;


--
-- Name: submissions_1615215600_1632495600_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_1615215600_1632495600_user_id_length_idx;


--
-- Name: submissions_1615215600_1632495600_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_1615215600_1632495600_user_id_point_idx;


--
-- Name: submissions_1615215600_16324956_execution_time_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_epoch_second_index ATTACH PARTITION public.submissions_1615215600_16324956_execution_time_epoch_second_idx;


--
-- Name: submissions_1632495600_1649775600_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_1632495600_1649775600_contest_id_epoch_second_idx;


--
-- Name: submissions_1632495600_1649775600_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_1632495600_1649775600_contest_id_execution_time_idx;


--
-- Name: submissions_1632495600_1649775600_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_1632495600_1649775600_contest_id_length_idx;


--
-- Name: submissions_1632495600_1649775600_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_1632495600_1649775600_contest_id_point_idx;


--
-- Name: submissions_1632495600_1649775600_epoch_second_updated_at_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_updated_at_index ATTACH PARTITION public.submissions_1632495600_1649775600_epoch_second_updated_at_idx;


--
-- Name: submissions_1632495600_1649775600_execution_time_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_length_index ATTACH PARTITION public.submissions_1632495600_1649775600_execution_time_length_idx;


--
-- Name: submissions_1632495600_1649775600_execution_time_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_point_index ATTACH PARTITION public.submissions_1632495600_1649775600_execution_time_point_idx;


--
-- Name: submissions_1632495600_1649775600_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_id_epoch_second_unique ATTACH PARTITION public.submissions_1632495600_1649775600_id_epoch_second_idx;


--
-- Name: submissions_1632495600_1649775600_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_1632495600_1649775600_language_epoch_second_idx;


--
-- Name: submissions_1632495600_1649775600_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_1632495600_1649775600_language_execution_time_idx;


--
-- Name: submissions_1632495600_1649775600_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_1632495600_1649775600_language_length_idx;


--
-- Name: submissions_1632495600_1649775600_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_1632495600_1649775600_language_point_idx;


--
-- Name: submissions_1632495600_1649775600_length_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_epoch_second_index ATTACH PARTITION public.submissions_1632495600_1649775600_length_epoch_second_idx;


--
-- Name: submissions_1632495600_1649775600_length_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_execution_time_index ATTACH PARTITION public.submissions_1632495600_1649775600_length_execution_time_idx;


--
-- Name: submissions_1632495600_1649775600_length_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_point_index ATTACH PARTITION public.submissions_1632495600_1649775600_length_point_idx;


--
-- Name: submissions_1632495600_1649775600_point_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_epoch_second_index ATTACH PARTITION public.submissions_1632495600_1649775600_point_epoch_second_idx;


--
-- Name: submissions_1632495600_1649775600_point_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_execution_time_index ATTACH PARTITION public.submissions_1632495600_1649775600_point_execution_time_idx;


--
-- Name: submissions_1632495600_1649775600_point_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_length_index ATTACH PARTITION public.submissions_1632495600_1649775600_point_length_idx;


--
-- Name: submissions_1632495600_1649775600_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_1632495600_1649775600_problem_id_epoch_second_idx;


--
-- Name: submissions_1632495600_1649775600_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_1632495600_1649775600_problem_id_execution_time_idx;


--
-- Name: submissions_1632495600_1649775600_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_1632495600_1649775600_problem_id_length_idx;


--
-- Name: submissions_1632495600_1649775600_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_1632495600_1649775600_problem_id_point_idx;


--
-- Name: submissions_1632495600_1649775600_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_1632495600_1649775600_result_epoch_second_idx;


--
-- Name: submissions_1632495600_1649775600_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_1632495600_1649775600_result_execution_time_idx;


--
-- Name: submissions_1632495600_1649775600_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_1632495600_1649775600_result_length_idx;


--
-- Name: submissions_1632495600_1649775600_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_1632495600_1649775600_result_point_idx;


--
-- Name: submissions_1632495600_1649775600_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_1632495600_1649775600_user_id_epoch_second_idx;


--
-- Name: submissions_1632495600_1649775600_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_1632495600_1649775600_user_id_execution_time_idx;


--
-- Name: submissions_1632495600_1649775600_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_1632495600_1649775600_user_id_length_idx;


--
-- Name: submissions_1632495600_1649775600_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_1632495600_1649775600_user_id_point_idx;


--
-- Name: submissions_1632495600_16497756_execution_time_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_epoch_second_index ATTACH PARTITION public.submissions_1632495600_16497756_execution_time_epoch_second_idx;


--
-- Name: submissions_1649775600_1667055600_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_1649775600_1667055600_contest_id_epoch_second_idx;


--
-- Name: submissions_1649775600_1667055600_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_1649775600_1667055600_contest_id_execution_time_idx;


--
-- Name: submissions_1649775600_1667055600_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_1649775600_1667055600_contest_id_length_idx;


--
-- Name: submissions_1649775600_1667055600_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_1649775600_1667055600_contest_id_point_idx;


--
-- Name: submissions_1649775600_1667055600_epoch_second_updated_at_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_updated_at_index ATTACH PARTITION public.submissions_1649775600_1667055600_epoch_second_updated_at_idx;


--
-- Name: submissions_1649775600_1667055600_execution_time_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_length_index ATTACH PARTITION public.submissions_1649775600_1667055600_execution_time_length_idx;


--
-- Name: submissions_1649775600_1667055600_execution_time_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_point_index ATTACH PARTITION public.submissions_1649775600_1667055600_execution_time_point_idx;


--
-- Name: submissions_1649775600_1667055600_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_id_epoch_second_unique ATTACH PARTITION public.submissions_1649775600_1667055600_id_epoch_second_idx;


--
-- Name: submissions_1649775600_1667055600_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_1649775600_1667055600_language_epoch_second_idx;


--
-- Name: submissions_1649775600_1667055600_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_1649775600_1667055600_language_execution_time_idx;


--
-- Name: submissions_1649775600_1667055600_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_1649775600_1667055600_language_length_idx;


--
-- Name: submissions_1649775600_1667055600_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_1649775600_1667055600_language_point_idx;


--
-- Name: submissions_1649775600_1667055600_length_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_epoch_second_index ATTACH PARTITION public.submissions_1649775600_1667055600_length_epoch_second_idx;


--
-- Name: submissions_1649775600_1667055600_length_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_execution_time_index ATTACH PARTITION public.submissions_1649775600_1667055600_length_execution_time_idx;


--
-- Name: submissions_1649775600_1667055600_length_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_point_index ATTACH PARTITION public.submissions_1649775600_1667055600_length_point_idx;


--
-- Name: submissions_1649775600_1667055600_point_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_epoch_second_index ATTACH PARTITION public.submissions_1649775600_1667055600_point_epoch_second_idx;


--
-- Name: submissions_1649775600_1667055600_point_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_execution_time_index ATTACH PARTITION public.submissions_1649775600_1667055600_point_execution_time_idx;


--
-- Name: submissions_1649775600_1667055600_point_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_length_index ATTACH PARTITION public.submissions_1649775600_1667055600_point_length_idx;


--
-- Name: submissions_1649775600_1667055600_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_1649775600_1667055600_problem_id_epoch_second_idx;


--
-- Name: submissions_1649775600_1667055600_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_1649775600_1667055600_problem_id_execution_time_idx;


--
-- Name: submissions_1649775600_1667055600_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_1649775600_1667055600_problem_id_length_idx;


--
-- Name: submissions_1649775600_1667055600_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_1649775600_1667055600_problem_id_point_idx;


--
-- Name: submissions_1649775600_1667055600_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_1649775600_1667055600_result_epoch_second_idx;


--
-- Name: submissions_1649775600_1667055600_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_1649775600_1667055600_result_execution_time_idx;


--
-- Name: submissions_1649775600_1667055600_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_1649775600_1667055600_result_length_idx;


--
-- Name: submissions_1649775600_1667055600_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_1649775600_1667055600_result_point_idx;


--
-- Name: submissions_1649775600_1667055600_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_1649775600_1667055600_user_id_epoch_second_idx;


--
-- Name: submissions_1649775600_1667055600_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_1649775600_1667055600_user_id_execution_time_idx;


--
-- Name: submissions_1649775600_1667055600_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_1649775600_1667055600_user_id_length_idx;


--
-- Name: submissions_1649775600_1667055600_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_1649775600_1667055600_user_id_point_idx;


--
-- Name: submissions_1649775600_16670556_execution_time_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_epoch_second_index ATTACH PARTITION public.submissions_1649775600_16670556_execution_time_epoch_second_idx;


--
-- Name: submissions_1667055600_1684335600_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_1667055600_1684335600_contest_id_epoch_second_idx;


--
-- Name: submissions_1667055600_1684335600_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_1667055600_1684335600_contest_id_execution_time_idx;


--
-- Name: submissions_1667055600_1684335600_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_1667055600_1684335600_contest_id_length_idx;


--
-- Name: submissions_1667055600_1684335600_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_1667055600_1684335600_contest_id_point_idx;


--
-- Name: submissions_1667055600_1684335600_epoch_second_updated_at_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_updated_at_index ATTACH PARTITION public.submissions_1667055600_1684335600_epoch_second_updated_at_idx;


--
-- Name: submissions_1667055600_1684335600_execution_time_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_length_index ATTACH PARTITION public.submissions_1667055600_1684335600_execution_time_length_idx;


--
-- Name: submissions_1667055600_1684335600_execution_time_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_point_index ATTACH PARTITION public.submissions_1667055600_1684335600_execution_time_point_idx;


--
-- Name: submissions_1667055600_1684335600_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_id_epoch_second_unique ATTACH PARTITION public.submissions_1667055600_1684335600_id_epoch_second_idx;


--
-- Name: submissions_1667055600_1684335600_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_1667055600_1684335600_language_epoch_second_idx;


--
-- Name: submissions_1667055600_1684335600_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_1667055600_1684335600_language_execution_time_idx;


--
-- Name: submissions_1667055600_1684335600_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_1667055600_1684335600_language_length_idx;


--
-- Name: submissions_1667055600_1684335600_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_1667055600_1684335600_language_point_idx;


--
-- Name: submissions_1667055600_1684335600_length_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_epoch_second_index ATTACH PARTITION public.submissions_1667055600_1684335600_length_epoch_second_idx;


--
-- Name: submissions_1667055600_1684335600_length_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_execution_time_index ATTACH PARTITION public.submissions_1667055600_1684335600_length_execution_time_idx;


--
-- Name: submissions_1667055600_1684335600_length_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_point_index ATTACH PARTITION public.submissions_1667055600_1684335600_length_point_idx;


--
-- Name: submissions_1667055600_1684335600_point_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_epoch_second_index ATTACH PARTITION public.submissions_1667055600_1684335600_point_epoch_second_idx;


--
-- Name: submissions_1667055600_1684335600_point_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_execution_time_index ATTACH PARTITION public.submissions_1667055600_1684335600_point_execution_time_idx;


--
-- Name: submissions_1667055600_1684335600_point_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_length_index ATTACH PARTITION public.submissions_1667055600_1684335600_point_length_idx;


--
-- Name: submissions_1667055600_1684335600_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_1667055600_1684335600_problem_id_epoch_second_idx;


--
-- Name: submissions_1667055600_1684335600_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_1667055600_1684335600_problem_id_execution_time_idx;


--
-- Name: submissions_1667055600_1684335600_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_1667055600_1684335600_problem_id_length_idx;


--
-- Name: submissions_1667055600_1684335600_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_1667055600_1684335600_problem_id_point_idx;


--
-- Name: submissions_1667055600_1684335600_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_1667055600_1684335600_result_epoch_second_idx;


--
-- Name: submissions_1667055600_1684335600_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_1667055600_1684335600_result_execution_time_idx;


--
-- Name: submissions_1667055600_1684335600_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_1667055600_1684335600_result_length_idx;


--
-- Name: submissions_1667055600_1684335600_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_1667055600_1684335600_result_point_idx;


--
-- Name: submissions_1667055600_1684335600_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_1667055600_1684335600_user_id_epoch_second_idx;


--
-- Name: submissions_1667055600_1684335600_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_1667055600_1684335600_user_id_execution_time_idx;


--
-- Name: submissions_1667055600_1684335600_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_1667055600_1684335600_user_id_length_idx;


--
-- Name: submissions_1667055600_1684335600_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_1667055600_1684335600_user_id_point_idx;


--
-- Name: submissions_1667055600_16843356_execution_time_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_epoch_second_index ATTACH PARTITION public.submissions_1667055600_16843356_execution_time_epoch_second_idx;


--
-- Name: submissions_1684335600_1701615600_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_1684335600_1701615600_contest_id_epoch_second_idx;


--
-- Name: submissions_1684335600_1701615600_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_1684335600_1701615600_contest_id_execution_time_idx;


--
-- Name: submissions_1684335600_1701615600_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_1684335600_1701615600_contest_id_length_idx;


--
-- Name: submissions_1684335600_1701615600_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_1684335600_1701615600_contest_id_point_idx;


--
-- Name: submissions_1684335600_1701615600_epoch_second_updated_at_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_updated_at_index ATTACH PARTITION public.submissions_1684335600_1701615600_epoch_second_updated_at_idx;


--
-- Name: submissions_1684335600_1701615600_execution_time_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_length_index ATTACH PARTITION public.submissions_1684335600_1701615600_execution_time_length_idx;


--
-- Name: submissions_1684335600_1701615600_execution_time_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_point_index ATTACH PARTITION public.submissions_1684335600_1701615600_execution_time_point_idx;


--
-- Name: submissions_1684335600_1701615600_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_id_epoch_second_unique ATTACH PARTITION public.submissions_1684335600_1701615600_id_epoch_second_idx;


--
-- Name: submissions_1684335600_1701615600_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_1684335600_1701615600_language_epoch_second_idx;


--
-- Name: submissions_1684335600_1701615600_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_1684335600_1701615600_language_execution_time_idx;


--
-- Name: submissions_1684335600_1701615600_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_1684335600_1701615600_language_length_idx;


--
-- Name: submissions_1684335600_1701615600_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_1684335600_1701615600_language_point_idx;


--
-- Name: submissions_1684335600_1701615600_length_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_epoch_second_index ATTACH PARTITION public.submissions_1684335600_1701615600_length_epoch_second_idx;


--
-- Name: submissions_1684335600_1701615600_length_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_execution_time_index ATTACH PARTITION public.submissions_1684335600_1701615600_length_execution_time_idx;


--
-- Name: submissions_1684335600_1701615600_length_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_point_index ATTACH PARTITION public.submissions_1684335600_1701615600_length_point_idx;


--
-- Name: submissions_1684335600_1701615600_point_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_epoch_second_index ATTACH PARTITION public.submissions_1684335600_1701615600_point_epoch_second_idx;


--
-- Name: submissions_1684335600_1701615600_point_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_execution_time_index ATTACH PARTITION public.submissions_1684335600_1701615600_point_execution_time_idx;


--
-- Name: submissions_1684335600_1701615600_point_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_length_index ATTACH PARTITION public.submissions_1684335600_1701615600_point_length_idx;


--
-- Name: submissions_1684335600_1701615600_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_1684335600_1701615600_problem_id_epoch_second_idx;


--
-- Name: submissions_1684335600_1701615600_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_1684335600_1701615600_problem_id_execution_time_idx;


--
-- Name: submissions_1684335600_1701615600_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_1684335600_1701615600_problem_id_length_idx;


--
-- Name: submissions_1684335600_1701615600_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_1684335600_1701615600_problem_id_point_idx;


--
-- Name: submissions_1684335600_1701615600_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_1684335600_1701615600_result_epoch_second_idx;


--
-- Name: submissions_1684335600_1701615600_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_1684335600_1701615600_result_execution_time_idx;


--
-- Name: submissions_1684335600_1701615600_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_1684335600_1701615600_result_length_idx;


--
-- Name: submissions_1684335600_1701615600_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_1684335600_1701615600_result_point_idx;


--
-- Name: submissions_1684335600_1701615600_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_1684335600_1701615600_user_id_epoch_second_idx;


--
-- Name: submissions_1684335600_1701615600_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_1684335600_1701615600_user_id_execution_time_idx;


--
-- Name: submissions_1684335600_1701615600_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_1684335600_1701615600_user_id_length_idx;


--
-- Name: submissions_1684335600_1701615600_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_1684335600_1701615600_user_id_point_idx;


--
-- Name: submissions_1684335600_17016156_execution_time_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_epoch_second_index ATTACH PARTITION public.submissions_1684335600_17016156_execution_time_epoch_second_idx;


--
-- Name: submissions_1701615600_1718895600_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_1701615600_1718895600_contest_id_epoch_second_idx;


--
-- Name: submissions_1701615600_1718895600_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_1701615600_1718895600_contest_id_execution_time_idx;


--
-- Name: submissions_1701615600_1718895600_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_1701615600_1718895600_contest_id_length_idx;


--
-- Name: submissions_1701615600_1718895600_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_1701615600_1718895600_contest_id_point_idx;


--
-- Name: submissions_1701615600_1718895600_epoch_second_updated_at_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_updated_at_index ATTACH PARTITION public.submissions_1701615600_1718895600_epoch_second_updated_at_idx;


--
-- Name: submissions_1701615600_1718895600_execution_time_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_length_index ATTACH PARTITION public.submissions_1701615600_1718895600_execution_time_length_idx;


--
-- Name: submissions_1701615600_1718895600_execution_time_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_point_index ATTACH PARTITION public.submissions_1701615600_1718895600_execution_time_point_idx;


--
-- Name: submissions_1701615600_1718895600_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_id_epoch_second_unique ATTACH PARTITION public.submissions_1701615600_1718895600_id_epoch_second_idx;


--
-- Name: submissions_1701615600_1718895600_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_1701615600_1718895600_language_epoch_second_idx;


--
-- Name: submissions_1701615600_1718895600_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_1701615600_1718895600_language_execution_time_idx;


--
-- Name: submissions_1701615600_1718895600_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_1701615600_1718895600_language_length_idx;


--
-- Name: submissions_1701615600_1718895600_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_1701615600_1718895600_language_point_idx;


--
-- Name: submissions_1701615600_1718895600_length_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_epoch_second_index ATTACH PARTITION public.submissions_1701615600_1718895600_length_epoch_second_idx;


--
-- Name: submissions_1701615600_1718895600_length_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_execution_time_index ATTACH PARTITION public.submissions_1701615600_1718895600_length_execution_time_idx;


--
-- Name: submissions_1701615600_1718895600_length_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_point_index ATTACH PARTITION public.submissions_1701615600_1718895600_length_point_idx;


--
-- Name: submissions_1701615600_1718895600_point_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_epoch_second_index ATTACH PARTITION public.submissions_1701615600_1718895600_point_epoch_second_idx;


--
-- Name: submissions_1701615600_1718895600_point_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_execution_time_index ATTACH PARTITION public.submissions_1701615600_1718895600_point_execution_time_idx;


--
-- Name: submissions_1701615600_1718895600_point_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_length_index ATTACH PARTITION public.submissions_1701615600_1718895600_point_length_idx;


--
-- Name: submissions_1701615600_1718895600_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_1701615600_1718895600_problem_id_epoch_second_idx;


--
-- Name: submissions_1701615600_1718895600_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_1701615600_1718895600_problem_id_execution_time_idx;


--
-- Name: submissions_1701615600_1718895600_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_1701615600_1718895600_problem_id_length_idx;


--
-- Name: submissions_1701615600_1718895600_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_1701615600_1718895600_problem_id_point_idx;


--
-- Name: submissions_1701615600_1718895600_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_1701615600_1718895600_result_epoch_second_idx;


--
-- Name: submissions_1701615600_1718895600_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_1701615600_1718895600_result_execution_time_idx;


--
-- Name: submissions_1701615600_1718895600_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_1701615600_1718895600_result_length_idx;


--
-- Name: submissions_1701615600_1718895600_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_1701615600_1718895600_result_point_idx;


--
-- Name: submissions_1701615600_1718895600_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_1701615600_1718895600_user_id_epoch_second_idx;


--
-- Name: submissions_1701615600_1718895600_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_1701615600_1718895600_user_id_execution_time_idx;


--
-- Name: submissions_1701615600_1718895600_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_1701615600_1718895600_user_id_length_idx;


--
-- Name: submissions_1701615600_1718895600_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_1701615600_1718895600_user_id_point_idx;


--
-- Name: submissions_1701615600_17188956_execution_time_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_epoch_second_index ATTACH PARTITION public.submissions_1701615600_17188956_execution_time_epoch_second_idx;


--
-- Name: submissions_1718895600_1736175600_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_1718895600_1736175600_contest_id_epoch_second_idx;


--
-- Name: submissions_1718895600_1736175600_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_1718895600_1736175600_contest_id_execution_time_idx;


--
-- Name: submissions_1718895600_1736175600_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_1718895600_1736175600_contest_id_length_idx;


--
-- Name: submissions_1718895600_1736175600_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_1718895600_1736175600_contest_id_point_idx;


--
-- Name: submissions_1718895600_1736175600_epoch_second_updated_at_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_updated_at_index ATTACH PARTITION public.submissions_1718895600_1736175600_epoch_second_updated_at_idx;


--
-- Name: submissions_1718895600_1736175600_execution_time_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_length_index ATTACH PARTITION public.submissions_1718895600_1736175600_execution_time_length_idx;


--
-- Name: submissions_1718895600_1736175600_execution_time_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_point_index ATTACH PARTITION public.submissions_1718895600_1736175600_execution_time_point_idx;


--
-- Name: submissions_1718895600_1736175600_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_id_epoch_second_unique ATTACH PARTITION public.submissions_1718895600_1736175600_id_epoch_second_idx;


--
-- Name: submissions_1718895600_1736175600_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_1718895600_1736175600_language_epoch_second_idx;


--
-- Name: submissions_1718895600_1736175600_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_1718895600_1736175600_language_execution_time_idx;


--
-- Name: submissions_1718895600_1736175600_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_1718895600_1736175600_language_length_idx;


--
-- Name: submissions_1718895600_1736175600_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_1718895600_1736175600_language_point_idx;


--
-- Name: submissions_1718895600_1736175600_length_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_epoch_second_index ATTACH PARTITION public.submissions_1718895600_1736175600_length_epoch_second_idx;


--
-- Name: submissions_1718895600_1736175600_length_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_execution_time_index ATTACH PARTITION public.submissions_1718895600_1736175600_length_execution_time_idx;


--
-- Name: submissions_1718895600_1736175600_length_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_point_index ATTACH PARTITION public.submissions_1718895600_1736175600_length_point_idx;


--
-- Name: submissions_1718895600_1736175600_point_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_epoch_second_index ATTACH PARTITION public.submissions_1718895600_1736175600_point_epoch_second_idx;


--
-- Name: submissions_1718895600_1736175600_point_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_execution_time_index ATTACH PARTITION public.submissions_1718895600_1736175600_point_execution_time_idx;


--
-- Name: submissions_1718895600_1736175600_point_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_length_index ATTACH PARTITION public.submissions_1718895600_1736175600_point_length_idx;


--
-- Name: submissions_1718895600_1736175600_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_1718895600_1736175600_problem_id_epoch_second_idx;


--
-- Name: submissions_1718895600_1736175600_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_1718895600_1736175600_problem_id_execution_time_idx;


--
-- Name: submissions_1718895600_1736175600_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_1718895600_1736175600_problem_id_length_idx;


--
-- Name: submissions_1718895600_1736175600_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_1718895600_1736175600_problem_id_point_idx;


--
-- Name: submissions_1718895600_1736175600_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_1718895600_1736175600_result_epoch_second_idx;


--
-- Name: submissions_1718895600_1736175600_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_1718895600_1736175600_result_execution_time_idx;


--
-- Name: submissions_1718895600_1736175600_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_1718895600_1736175600_result_length_idx;


--
-- Name: submissions_1718895600_1736175600_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_1718895600_1736175600_result_point_idx;


--
-- Name: submissions_1718895600_1736175600_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_1718895600_1736175600_user_id_epoch_second_idx;


--
-- Name: submissions_1718895600_1736175600_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_1718895600_1736175600_user_id_execution_time_idx;


--
-- Name: submissions_1718895600_1736175600_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_1718895600_1736175600_user_id_length_idx;


--
-- Name: submissions_1718895600_1736175600_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_1718895600_1736175600_user_id_point_idx;


--
-- Name: submissions_1718895600_17361756_execution_time_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_epoch_second_index ATTACH PARTITION public.submissions_1718895600_17361756_execution_time_epoch_second_idx;


--
-- Name: submissions_1736175600_1753455600_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_1736175600_1753455600_contest_id_epoch_second_idx;


--
-- Name: submissions_1736175600_1753455600_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_1736175600_1753455600_contest_id_execution_time_idx;


--
-- Name: submissions_1736175600_1753455600_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_1736175600_1753455600_contest_id_length_idx;


--
-- Name: submissions_1736175600_1753455600_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_1736175600_1753455600_contest_id_point_idx;


--
-- Name: submissions_1736175600_1753455600_epoch_second_updated_at_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_updated_at_index ATTACH PARTITION public.submissions_1736175600_1753455600_epoch_second_updated_at_idx;


--
-- Name: submissions_1736175600_1753455600_execution_time_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_length_index ATTACH PARTITION public.submissions_1736175600_1753455600_execution_time_length_idx;


--
-- Name: submissions_1736175600_1753455600_execution_time_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_point_index ATTACH PARTITION public.submissions_1736175600_1753455600_execution_time_point_idx;


--
-- Name: submissions_1736175600_1753455600_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_id_epoch_second_unique ATTACH PARTITION public.submissions_1736175600_1753455600_id_epoch_second_idx;


--
-- Name: submissions_1736175600_1753455600_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_1736175600_1753455600_language_epoch_second_idx;


--
-- Name: submissions_1736175600_1753455600_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_1736175600_1753455600_language_execution_time_idx;


--
-- Name: submissions_1736175600_1753455600_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_1736175600_1753455600_language_length_idx;


--
-- Name: submissions_1736175600_1753455600_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_1736175600_1753455600_language_point_idx;


--
-- Name: submissions_1736175600_1753455600_length_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_epoch_second_index ATTACH PARTITION public.submissions_1736175600_1753455600_length_epoch_second_idx;


--
-- Name: submissions_1736175600_1753455600_length_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_execution_time_index ATTACH PARTITION public.submissions_1736175600_1753455600_length_execution_time_idx;


--
-- Name: submissions_1736175600_1753455600_length_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_point_index ATTACH PARTITION public.submissions_1736175600_1753455600_length_point_idx;


--
-- Name: submissions_1736175600_1753455600_point_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_epoch_second_index ATTACH PARTITION public.submissions_1736175600_1753455600_point_epoch_second_idx;


--
-- Name: submissions_1736175600_1753455600_point_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_execution_time_index ATTACH PARTITION public.submissions_1736175600_1753455600_point_execution_time_idx;


--
-- Name: submissions_1736175600_1753455600_point_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_length_index ATTACH PARTITION public.submissions_1736175600_1753455600_point_length_idx;


--
-- Name: submissions_1736175600_1753455600_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_1736175600_1753455600_problem_id_epoch_second_idx;


--
-- Name: submissions_1736175600_1753455600_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_1736175600_1753455600_problem_id_execution_time_idx;


--
-- Name: submissions_1736175600_1753455600_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_1736175600_1753455600_problem_id_length_idx;


--
-- Name: submissions_1736175600_1753455600_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_1736175600_1753455600_problem_id_point_idx;


--
-- Name: submissions_1736175600_1753455600_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_1736175600_1753455600_result_epoch_second_idx;


--
-- Name: submissions_1736175600_1753455600_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_1736175600_1753455600_result_execution_time_idx;


--
-- Name: submissions_1736175600_1753455600_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_1736175600_1753455600_result_length_idx;


--
-- Name: submissions_1736175600_1753455600_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_1736175600_1753455600_result_point_idx;


--
-- Name: submissions_1736175600_1753455600_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_1736175600_1753455600_user_id_epoch_second_idx;


--
-- Name: submissions_1736175600_1753455600_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_1736175600_1753455600_user_id_execution_time_idx;


--
-- Name: submissions_1736175600_1753455600_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_1736175600_1753455600_user_id_length_idx;


--
-- Name: submissions_1736175600_1753455600_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_1736175600_1753455600_user_id_point_idx;


--
-- Name: submissions_1736175600_17534556_execution_time_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_epoch_second_index ATTACH PARTITION public.submissions_1736175600_17534556_execution_time_epoch_second_idx;


--
-- Name: submissions_1753455600_1770735600_contest_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_epoch_second_index ATTACH PARTITION public.submissions_1753455600_1770735600_contest_id_epoch_second_idx;


--
-- Name: submissions_1753455600_1770735600_contest_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_execution_time_index ATTACH PARTITION public.submissions_1753455600_1770735600_contest_id_execution_time_idx;


--
-- Name: submissions_1753455600_1770735600_contest_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_length_index ATTACH PARTITION public.submissions_1753455600_1770735600_contest_id_length_idx;


--
-- Name: submissions_1753455600_1770735600_contest_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_contest_id_point_index ATTACH PARTITION public.submissions_1753455600_1770735600_contest_id_point_idx;


--
-- Name: submissions_1753455600_1770735600_epoch_second_updated_at_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_updated_at_index ATTACH PARTITION public.submissions_1753455600_1770735600_epoch_second_updated_at_idx;


--
-- Name: submissions_1753455600_1770735600_execution_time_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_length_index ATTACH PARTITION public.submissions_1753455600_1770735600_execution_time_length_idx;


--
-- Name: submissions_1753455600_1770735600_execution_time_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_point_index ATTACH PARTITION public.submissions_1753455600_1770735600_execution_time_point_idx;


--
-- Name: submissions_1753455600_1770735600_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_id_epoch_second_unique ATTACH PARTITION public.submissions_1753455600_1770735600_id_epoch_second_idx;


--
-- Name: submissions_1753455600_1770735600_language_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_epoch_second_index ATTACH PARTITION public.submissions_1753455600_1770735600_language_epoch_second_idx;


--
-- Name: submissions_1753455600_1770735600_language_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_execution_time_index ATTACH PARTITION public.submissions_1753455600_1770735600_language_execution_time_idx;


--
-- Name: submissions_1753455600_1770735600_language_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_length_index ATTACH PARTITION public.submissions_1753455600_1770735600_language_length_idx;


--
-- Name: submissions_1753455600_1770735600_language_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_language_point_index ATTACH PARTITION public.submissions_1753455600_1770735600_language_point_idx;


--
-- Name: submissions_1753455600_1770735600_length_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_epoch_second_index ATTACH PARTITION public.submissions_1753455600_1770735600_length_epoch_second_idx;


--
-- Name: submissions_1753455600_1770735600_length_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_execution_time_index ATTACH PARTITION public.submissions_1753455600_1770735600_length_execution_time_idx;


--
-- Name: submissions_1753455600_1770735600_length_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_length_point_index ATTACH PARTITION public.submissions_1753455600_1770735600_length_point_idx;


--
-- Name: submissions_1753455600_1770735600_point_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_epoch_second_index ATTACH PARTITION public.submissions_1753455600_1770735600_point_epoch_second_idx;


--
-- Name: submissions_1753455600_1770735600_point_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_execution_time_index ATTACH PARTITION public.submissions_1753455600_1770735600_point_execution_time_idx;


--
-- Name: submissions_1753455600_1770735600_point_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_point_length_index ATTACH PARTITION public.submissions_1753455600_1770735600_point_length_idx;


--
-- Name: submissions_1753455600_1770735600_problem_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_epoch_second_index ATTACH PARTITION public.submissions_1753455600_1770735600_problem_id_epoch_second_idx;


--
-- Name: submissions_1753455600_1770735600_problem_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_execution_time_index ATTACH PARTITION public.submissions_1753455600_1770735600_problem_id_execution_time_idx;


--
-- Name: submissions_1753455600_1770735600_problem_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_length_index ATTACH PARTITION public.submissions_1753455600_1770735600_problem_id_length_idx;


--
-- Name: submissions_1753455600_1770735600_problem_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_problem_id_point_index ATTACH PARTITION public.submissions_1753455600_1770735600_problem_id_point_idx;


--
-- Name: submissions_1753455600_1770735600_result_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_epoch_second_index ATTACH PARTITION public.submissions_1753455600_1770735600_result_epoch_second_idx;


--
-- Name: submissions_1753455600_1770735600_result_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_execution_time_index ATTACH PARTITION public.submissions_1753455600_1770735600_result_execution_time_idx;


--
-- Name: submissions_1753455600_1770735600_result_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_length_index ATTACH PARTITION public.submissions_1753455600_1770735600_result_length_idx;


--
-- Name: submissions_1753455600_1770735600_result_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_result_point_index ATTACH PARTITION public.submissions_1753455600_1770735600_result_point_idx;


--
-- Name: submissions_1753455600_1770735600_user_id_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_epoch_second_index ATTACH PARTITION public.submissions_1753455600_1770735600_user_id_epoch_second_idx;


--
-- Name: submissions_1753455600_1770735600_user_id_execution_time_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_execution_time_index ATTACH PARTITION public.submissions_1753455600_1770735600_user_id_execution_time_idx;


--
-- Name: submissions_1753455600_1770735600_user_id_length_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_length_index ATTACH PARTITION public.submissions_1753455600_1770735600_user_id_length_idx;


--
-- Name: submissions_1753455600_1770735600_user_id_point_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_user_id_point_index ATTACH PARTITION public.submissions_1753455600_1770735600_user_id_point_idx;


--
-- Name: submissions_1753455600_17707356_execution_time_epoch_second_idx; Type: INDEX ATTACH; Schema: public; Owner: -
--

ALTER INDEX public.submissions_execution_time_epoch_second_index ATTACH PARTITION public.submissions_1753455600_17707356_execution_time_epoch_second_idx;


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
    ('20240831033102');
