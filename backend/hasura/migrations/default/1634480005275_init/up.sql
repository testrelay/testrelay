SET check_function_bodies = false;
CREATE FUNCTION public.set_current_timestamp_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
  _new record;
BEGIN
  _new := NEW;
  _new."updated_at" = NOW();
  RETURN _new;
END;
$$;
CREATE TABLE public.assignment_events (
    id integer NOT NULL,
    user_id integer,
    assignment_id integer NOT NULL,
    meta jsonb DEFAULT jsonb_build_object() NOT NULL,
    event_type character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);
CREATE SEQUENCE public.assignment_events_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.assignment_events_id_seq OWNED BY public.assignment_events.id;
CREATE TABLE public.assignment_status (
    value text NOT NULL
);
CREATE TABLE public.assignment_users (
    id integer NOT NULL,
    assignment_id integer NOT NULL,
    user_id integer NOT NULL
);
CREATE SEQUENCE public.assignment_users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.assignment_users_id_seq OWNED BY public.assignment_users.id;
CREATE TABLE public.assignments (
    id integer NOT NULL,
    test_id integer,
    candidate_name character varying NOT NULL,
    recruiter_id integer NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    time_limit integer NOT NULL,
    candidate_email character varying NOT NULL,
    status character varying DEFAULT 'sending'::character varying NOT NULL,
    candidate_id integer,
    invite_code uuid DEFAULT gen_random_uuid(),
    test_day_chosen date,
    test_time_chosen time without time zone,
    github_repo_url character varying,
    test_timezone_chosen character varying,
    choose_until date NOT NULL,
    step_arn character varying
);
CREATE TABLE public.business_users (
    id integer NOT NULL,
    business_id integer NOT NULL,
    user_id integer NOT NULL,
    user_type character varying NOT NULL
);
CREATE SEQUENCE public.business_users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.business_users_id_seq OWNED BY public.business_users.id;
CREATE TABLE public.businesses (
    id integer NOT NULL,
    name character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    setup boolean DEFAULT false NOT NULL,
    github_installation_id character varying,
    creator_id integer NOT NULL
);
CREATE SEQUENCE public.businesses_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.businesses_id_seq OWNED BY public.businesses.id;
CREATE SEQUENCE public.candidates_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.candidates_id_seq OWNED BY public.assignments.id;
CREATE TABLE public.languages (
    id integer NOT NULL,
    name character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);
CREATE SEQUENCE public.languages_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.languages_id_seq OWNED BY public.languages.id;
CREATE TABLE public.test_languages (
    test_id integer NOT NULL,
    language_id integer NOT NULL,
    id integer NOT NULL
);
CREATE SEQUENCE public.test_languages_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.test_languages_id_seq OWNED BY public.test_languages.id;
CREATE TABLE public.tests (
    id integer NOT NULL,
    name character varying NOT NULL,
    business_id integer NOT NULL,
    user_id integer NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    time_limit integer NOT NULL,
    test_window integer NOT NULL,
    github_repo character varying,
    zip character varying
);
CREATE SEQUENCE public.tests_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.tests_id_seq OWNED BY public.tests.id;
CREATE TABLE public.users (
    id integer NOT NULL,
    email character varying NOT NULL,
    auth_id character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    github_username character varying,
    github_access_token character varying
);
CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;
ALTER TABLE ONLY public.assignment_events ALTER COLUMN id SET DEFAULT nextval('public.assignment_events_id_seq'::regclass);
ALTER TABLE ONLY public.assignment_users ALTER COLUMN id SET DEFAULT nextval('public.assignment_users_id_seq'::regclass);
ALTER TABLE ONLY public.assignments ALTER COLUMN id SET DEFAULT nextval('public.candidates_id_seq'::regclass);
ALTER TABLE ONLY public.business_users ALTER COLUMN id SET DEFAULT nextval('public.business_users_id_seq'::regclass);
ALTER TABLE ONLY public.businesses ALTER COLUMN id SET DEFAULT nextval('public.businesses_id_seq'::regclass);
ALTER TABLE ONLY public.languages ALTER COLUMN id SET DEFAULT nextval('public.languages_id_seq'::regclass);
ALTER TABLE ONLY public.test_languages ALTER COLUMN id SET DEFAULT nextval('public.test_languages_id_seq'::regclass);
ALTER TABLE ONLY public.tests ALTER COLUMN id SET DEFAULT nextval('public.tests_id_seq'::regclass);
ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);
ALTER TABLE ONLY public.assignment_events
    ADD CONSTRAINT assignment_events_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.assignment_users
    ADD CONSTRAINT assignment_users_assignment_id_user_id_key UNIQUE (assignment_id, user_id);
ALTER TABLE ONLY public.assignment_users
    ADD CONSTRAINT assignment_users_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.business_users
    ADD CONSTRAINT business_users_business_id_user_id_user_type_key UNIQUE (business_id, user_id, user_type);
ALTER TABLE ONLY public.business_users
    ADD CONSTRAINT business_users_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.businesses
    ADD CONSTRAINT businesses_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.assignment_status
    ADD CONSTRAINT candidate_status_pkey PRIMARY KEY (value);
ALTER TABLE ONLY public.assignments
    ADD CONSTRAINT candidates_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.languages
    ADD CONSTRAINT languages_name_key UNIQUE (name);
ALTER TABLE ONLY public.languages
    ADD CONSTRAINT languages_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.test_languages
    ADD CONSTRAINT test_languages_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.test_languages
    ADD CONSTRAINT test_languages_test_id_language_id_key UNIQUE (test_id, language_id);
ALTER TABLE ONLY public.tests
    ADD CONSTRAINT tests_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_auth_id_key UNIQUE (auth_id);
ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);
CREATE TRIGGER set_public_businesses_updated_at BEFORE UPDATE ON public.businesses FOR EACH ROW EXECUTE FUNCTION public.set_current_timestamp_updated_at();
COMMENT ON TRIGGER set_public_businesses_updated_at ON public.businesses IS 'trigger to set value of column "updated_at" to current timestamp on row update';
CREATE TRIGGER set_public_candidates_updated_at BEFORE UPDATE ON public.assignments FOR EACH ROW EXECUTE FUNCTION public.set_current_timestamp_updated_at();
COMMENT ON TRIGGER set_public_candidates_updated_at ON public.assignments IS 'trigger to set value of column "updated_at" to current timestamp on row update';
CREATE TRIGGER set_public_languages_updated_at BEFORE UPDATE ON public.languages FOR EACH ROW EXECUTE FUNCTION public.set_current_timestamp_updated_at();
COMMENT ON TRIGGER set_public_languages_updated_at ON public.languages IS 'trigger to set value of column "updated_at" to current timestamp on row update';
CREATE TRIGGER set_public_tests_updated_at BEFORE UPDATE ON public.tests FOR EACH ROW EXECUTE FUNCTION public.set_current_timestamp_updated_at();
COMMENT ON TRIGGER set_public_tests_updated_at ON public.tests IS 'trigger to set value of column "updated_at" to current timestamp on row update';
CREATE TRIGGER set_public_users_updated_at BEFORE UPDATE ON public.users FOR EACH ROW EXECUTE FUNCTION public.set_current_timestamp_updated_at();
COMMENT ON TRIGGER set_public_users_updated_at ON public.users IS 'trigger to set value of column "updated_at" to current timestamp on row update';
ALTER TABLE ONLY public.assignment_events
    ADD CONSTRAINT assignment_events_assignment_id_fkey FOREIGN KEY (assignment_id) REFERENCES public.assignments(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.assignment_events
    ADD CONSTRAINT assignment_events_event_type_fkey FOREIGN KEY (event_type) REFERENCES public.assignment_status(value) ON UPDATE RESTRICT ON DELETE RESTRICT;
ALTER TABLE ONLY public.assignment_events
    ADD CONSTRAINT assignment_events_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.assignment_users
    ADD CONSTRAINT assignment_users_assignment_id_fkey FOREIGN KEY (assignment_id) REFERENCES public.assignments(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.assignment_users
    ADD CONSTRAINT assignment_users_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.assignments
    ADD CONSTRAINT assignments_candidate_id_fkey FOREIGN KEY (candidate_id) REFERENCES public.users(id) ON UPDATE SET NULL ON DELETE SET NULL;
ALTER TABLE ONLY public.business_users
    ADD CONSTRAINT business_users_business_id_fkey FOREIGN KEY (business_id) REFERENCES public.businesses(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.business_users
    ADD CONSTRAINT business_users_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.businesses
    ADD CONSTRAINT businesses_creator_id_fkey FOREIGN KEY (creator_id) REFERENCES public.users(id) ON UPDATE RESTRICT ON DELETE RESTRICT;
ALTER TABLE ONLY public.assignments
    ADD CONSTRAINT candidates_status_fkey FOREIGN KEY (status) REFERENCES public.assignment_status(value) ON UPDATE RESTRICT ON DELETE RESTRICT;
ALTER TABLE ONLY public.assignments
    ADD CONSTRAINT candidates_test_id_fkey FOREIGN KEY (test_id) REFERENCES public.tests(id) ON UPDATE RESTRICT ON DELETE SET NULL;
ALTER TABLE ONLY public.assignments
    ADD CONSTRAINT candidates_user_id_fkey FOREIGN KEY (recruiter_id) REFERENCES public.users(id) ON UPDATE RESTRICT ON DELETE CASCADE;
ALTER TABLE ONLY public.test_languages
    ADD CONSTRAINT test_languages_language_id_fkey FOREIGN KEY (language_id) REFERENCES public.languages(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.test_languages
    ADD CONSTRAINT test_languages_test_id_fkey FOREIGN KEY (test_id) REFERENCES public.tests(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.tests
    ADD CONSTRAINT tests_business_id_fkey FOREIGN KEY (business_id) REFERENCES public.businesses(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.tests
    ADD CONSTRAINT tests_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;
