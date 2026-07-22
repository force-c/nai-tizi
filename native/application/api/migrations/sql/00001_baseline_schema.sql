-- +goose Up
SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

-- Name: biz_attachment; Type: TABLE; Schema: public; Owner: -

CREATE TABLE public.biz_attachment (
    id bigint NOT NULL,
    env_id bigint NOT NULL,
    file_name character varying(255) NOT NULL,
    file_key character varying(512) NOT NULL,
    file_size bigint NOT NULL,
    file_type character varying(128),
    file_ext character varying(32),
    business_type character varying(64),
    business_id character varying(64),
    business_field character varying(64),
    is_public boolean DEFAULT false,
    access_url character varying(1024),
    metadata jsonb,
    status smallint DEFAULT 0,
    expire_time timestamp with time zone,
    create_by bigint,
    create_time timestamp with time zone,
    update_time timestamp with time zone
);


-- Name: m_role_api_permission; Type: TABLE; Schema: public; Owner: -

CREATE TABLE public.m_role_api_permission (
    id bigint NOT NULL,
    role_id bigint NOT NULL,
    permission_id bigint NOT NULL,
    create_by bigint,
    update_by bigint,
    created_time timestamp with time zone,
    updated_time timestamp with time zone
);


-- Name: m_role_menu; Type: TABLE; Schema: public; Owner: -

CREATE TABLE public.m_role_menu (
    id bigint NOT NULL,
    role_id bigint NOT NULL,
    menu_id bigint NOT NULL,
    create_by bigint,
    update_by bigint,
    created_time timestamp with time zone,
    updated_time timestamp with time zone
);


-- Name: m_user_api_permission; Type: TABLE; Schema: public; Owner: -

CREATE TABLE public.m_user_api_permission (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    permission_id bigint NOT NULL,
    create_by bigint,
    update_by bigint,
    created_time timestamp with time zone,
    updated_time timestamp with time zone
);


-- Name: m_user_role; Type: TABLE; Schema: public; Owner: -

CREATE TABLE public.m_user_role (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    role_id bigint NOT NULL,
    create_by bigint,
    update_by bigint,
    created_time timestamp with time zone,
    updated_time timestamp with time zone
);


-- Name: s_api_permission; Type: TABLE; Schema: public; Owner: -

CREATE TABLE public.s_api_permission (
    id bigint NOT NULL,
    parent_id bigint DEFAULT 0,
    module character varying(64) NOT NULL,
    code character varying(128) NOT NULL,
    name character varying(64) NOT NULL,
    node_type smallint DEFAULT 2 NOT NULL,
    action character varying(32) DEFAULT '*'::character varying NOT NULL,
    method character varying(16),
    path character varying(255),
    sort bigint DEFAULT 0,
    status smallint DEFAULT 0,
    remark character varying(500),
    create_by bigint,
    update_by bigint,
    created_time timestamp with time zone,
    updated_time timestamp with time zone
);


-- Name: s_auth_client; Type: TABLE; Schema: public; Owner: -

CREATE TABLE public.s_auth_client (
    client_id character varying(64) NOT NULL,
    grant_type character varying(255),
    device_type character varying(32),
    status smallint DEFAULT 0,
    timeout bigint DEFAULT 604800,
    active_timeout bigint DEFAULT 1800,
    remark character varying(500),
    create_by bigint,
    created_time timestamp with time zone,
    update_by bigint,
    updated_time timestamp with time zone
);


-- Name: s_config; Type: TABLE; Schema: public; Owner: -

CREATE TABLE public.s_config (
    id bigint NOT NULL,
    name character varying(128) NOT NULL,
    code character varying(128) NOT NULL,
    data jsonb,
    remark character varying(500),
    create_by bigint,
    created_time timestamp with time zone,
    update_by bigint,
    updated_time timestamp with time zone
);


-- Name: s_dict_data; Type: TABLE; Schema: public; Owner: -

CREATE TABLE public.s_dict_data (
    id bigint NOT NULL,
    parent_id bigint DEFAULT 0,
    sort bigint DEFAULT 0,
    dict_label character varying(128),
    dict_value character varying(128),
    dict_type character varying(128),
    is_default boolean DEFAULT false,
    status smallint DEFAULT 0,
    remark character varying(500),
    create_by bigint,
    created_time timestamp with time zone,
    update_by bigint,
    updated_time timestamp with time zone
);


-- Name: s_login_log; Type: TABLE; Schema: public; Owner: -

CREATE TABLE public.s_login_log (
    id bigint NOT NULL,
    user_name character varying(64),
    ipaddr character varying(64),
    login_location character varying(128),
    browser character varying(64),
    os character varying(64),
    status smallint DEFAULT 0,
    msg character varying(500),
    login_time timestamp with time zone,
    client_id character varying(64)
);


-- Name: s_menu; Type: TABLE; Schema: public; Owner: -

CREATE TABLE public.s_menu (
    id bigint NOT NULL,
    menu_name character varying(64) NOT NULL,
    parent_id bigint DEFAULT 0,
    sort bigint DEFAULT 0,
    path character varying(255),
    component character varying(255),
    query character varying(255),
    is_frame smallint DEFAULT 0,
    is_cache smallint DEFAULT 0,
    menu_type smallint NOT NULL,
    visible smallint DEFAULT 0,
    status smallint DEFAULT 0,
    perms character varying(255),
    icon character varying(64),
    remark character varying(500),
    create_by bigint,
    update_by bigint,
    created_time timestamp with time zone,
    updated_time timestamp with time zone
);


-- Name: s_oper_log; Type: TABLE; Schema: public; Owner: -

CREATE TABLE public.s_oper_log (
    id bigint NOT NULL,
    title character varying(128),
    business_type character varying(32),
    method character varying(255),
    request_method character varying(16),
    device_type character varying(32),
    oper_name character varying(64),
    oper_url character varying(1024),
    oper_ip character varying(64),
    oper_location character varying(128),
    oper_param text,
    json_result text,
    status character(1),
    error_msg text,
    oper_time timestamp with time zone,
    cost_time bigint,
    user_agent character varying(512)
);


-- Name: s_org; Type: TABLE; Schema: public; Owner: -

CREATE TABLE public.s_org (
    id bigint NOT NULL,
    parent_id bigint DEFAULT 0,
    ancestors character varying(512),
    org_name character varying(128) NOT NULL,
    org_code character varying(64),
    org_type character varying(32) DEFAULT 'company'::character varying,
    leader character varying(64),
    phone character varying(32),
    email character varying(128),
    status smallint DEFAULT 0,
    sort bigint DEFAULT 0,
    remark character varying(500),
    create_by bigint,
    update_by bigint,
    created_time timestamp with time zone,
    updated_time timestamp with time zone
);


-- Name: s_role; Type: TABLE; Schema: public; Owner: -

CREATE TABLE public.s_role (
    id bigint NOT NULL,
    role_key character varying(64) NOT NULL,
    role_name character varying(64) NOT NULL,
    sort bigint DEFAULT 0,
    status smallint DEFAULT 0,
    data_scope smallint DEFAULT 1,
    is_system boolean DEFAULT false,
    remark character varying(500),
    create_by bigint,
    update_by bigint,
    created_time timestamp with time zone,
    updated_time timestamp with time zone
);


-- Name: s_storage_env; Type: TABLE; Schema: public; Owner: -

CREATE TABLE public.s_storage_env (
    id bigint NOT NULL,
    name character varying(128) NOT NULL,
    code character varying(64) NOT NULL,
    storage_type character varying(32),
    is_default boolean DEFAULT false,
    status smallint DEFAULT 0,
    config jsonb NOT NULL,
    remark character varying(500),
    create_by bigint,
    created_time timestamp with time zone,
    update_by bigint,
    updated_time timestamp with time zone
);


-- Name: s_user; Type: TABLE; Schema: public; Owner: -

CREATE TABLE public.s_user (
    id bigint NOT NULL,
    user_name character varying(64) NOT NULL,
    nick_name character varying(64),
    user_type smallint DEFAULT 0,
    org_id bigint DEFAULT 0,
    email character varying(128),
    phonenumber character varying(32),
    sex smallint DEFAULT 2,
    avatar character varying(512),
    password character varying(255),
    status smallint DEFAULT 0,
    sort bigint DEFAULT 0,
    login_ip character varying(64),
    login_date bigint,
    open_id character varying(128),
    union_id character varying(128),
    remark character varying(500),
    create_by bigint,
    update_by bigint,
    created_time timestamp with time zone,
    updated_time timestamp with time zone
);


-- Name: biz_attachment biz_attachment_pkey; Type: CONSTRAINT; Schema: public; Owner: -

ALTER TABLE ONLY public.biz_attachment
    ADD CONSTRAINT biz_attachment_pkey PRIMARY KEY (id);


-- Name: m_role_api_permission m_role_api_permission_pkey; Type: CONSTRAINT; Schema: public; Owner: -

ALTER TABLE ONLY public.m_role_api_permission
    ADD CONSTRAINT m_role_api_permission_pkey PRIMARY KEY (id);


-- Name: m_role_menu m_role_menu_pkey; Type: CONSTRAINT; Schema: public; Owner: -

ALTER TABLE ONLY public.m_role_menu
    ADD CONSTRAINT m_role_menu_pkey PRIMARY KEY (id);


-- Name: m_user_api_permission m_user_api_permission_pkey; Type: CONSTRAINT; Schema: public; Owner: -

ALTER TABLE ONLY public.m_user_api_permission
    ADD CONSTRAINT m_user_api_permission_pkey PRIMARY KEY (id);


-- Name: m_user_role m_user_role_pkey; Type: CONSTRAINT; Schema: public; Owner: -

ALTER TABLE ONLY public.m_user_role
    ADD CONSTRAINT m_user_role_pkey PRIMARY KEY (id);


-- Name: s_api_permission s_api_permission_pkey; Type: CONSTRAINT; Schema: public; Owner: -

ALTER TABLE ONLY public.s_api_permission
    ADD CONSTRAINT s_api_permission_pkey PRIMARY KEY (id);


-- Name: s_auth_client s_auth_client_pkey; Type: CONSTRAINT; Schema: public; Owner: -

ALTER TABLE ONLY public.s_auth_client
    ADD CONSTRAINT s_auth_client_pkey PRIMARY KEY (client_id);


-- Name: s_config s_config_pkey; Type: CONSTRAINT; Schema: public; Owner: -

ALTER TABLE ONLY public.s_config
    ADD CONSTRAINT s_config_pkey PRIMARY KEY (id);


-- Name: s_dict_data s_dict_data_pkey; Type: CONSTRAINT; Schema: public; Owner: -

ALTER TABLE ONLY public.s_dict_data
    ADD CONSTRAINT s_dict_data_pkey PRIMARY KEY (id);


-- Name: s_login_log s_login_log_pkey; Type: CONSTRAINT; Schema: public; Owner: -

ALTER TABLE ONLY public.s_login_log
    ADD CONSTRAINT s_login_log_pkey PRIMARY KEY (id);


-- Name: s_menu s_menu_pkey; Type: CONSTRAINT; Schema: public; Owner: -

ALTER TABLE ONLY public.s_menu
    ADD CONSTRAINT s_menu_pkey PRIMARY KEY (id);


-- Name: s_oper_log s_oper_log_pkey; Type: CONSTRAINT; Schema: public; Owner: -

ALTER TABLE ONLY public.s_oper_log
    ADD CONSTRAINT s_oper_log_pkey PRIMARY KEY (id);


-- Name: s_org s_org_pkey; Type: CONSTRAINT; Schema: public; Owner: -

ALTER TABLE ONLY public.s_org
    ADD CONSTRAINT s_org_pkey PRIMARY KEY (id);


-- Name: s_role s_role_pkey; Type: CONSTRAINT; Schema: public; Owner: -

ALTER TABLE ONLY public.s_role
    ADD CONSTRAINT s_role_pkey PRIMARY KEY (id);


-- Name: s_storage_env s_storage_env_pkey; Type: CONSTRAINT; Schema: public; Owner: -

ALTER TABLE ONLY public.s_storage_env
    ADD CONSTRAINT s_storage_env_pkey PRIMARY KEY (id);


-- Name: s_user s_user_pkey; Type: CONSTRAINT; Schema: public; Owner: -

ALTER TABLE ONLY public.s_user
    ADD CONSTRAINT s_user_pkey PRIMARY KEY (id);


-- Name: idx_biz_attachment_business_id; Type: INDEX; Schema: public; Owner: -

CREATE INDEX idx_biz_attachment_business_id ON public.biz_attachment USING btree (business_id);


-- Name: idx_biz_attachment_business_type; Type: INDEX; Schema: public; Owner: -

CREATE INDEX idx_biz_attachment_business_type ON public.biz_attachment USING btree (business_type);


-- Name: idx_biz_attachment_env_id; Type: INDEX; Schema: public; Owner: -

CREATE INDEX idx_biz_attachment_env_id ON public.biz_attachment USING btree (env_id);


-- Name: idx_biz_attachment_file_key; Type: INDEX; Schema: public; Owner: -

CREATE INDEX idx_biz_attachment_file_key ON public.biz_attachment USING btree (file_key);


-- Name: idx_role_api_permission; Type: INDEX; Schema: public; Owner: -

CREATE UNIQUE INDEX idx_role_api_permission ON public.m_role_api_permission USING btree (role_id, permission_id);


-- Name: idx_role_menu; Type: INDEX; Schema: public; Owner: -

CREATE INDEX idx_role_menu ON public.m_role_menu USING btree (role_id, menu_id);


-- Name: idx_s_api_permission_code; Type: INDEX; Schema: public; Owner: -

CREATE UNIQUE INDEX idx_s_api_permission_code ON public.s_api_permission USING btree (code);


-- Name: idx_s_api_permission_module; Type: INDEX; Schema: public; Owner: -

CREATE INDEX idx_s_api_permission_module ON public.s_api_permission USING btree (module);


-- Name: idx_s_api_permission_parent_id; Type: INDEX; Schema: public; Owner: -

CREATE INDEX idx_s_api_permission_parent_id ON public.s_api_permission USING btree (parent_id);


-- Name: idx_s_config_code; Type: INDEX; Schema: public; Owner: -

CREATE INDEX idx_s_config_code ON public.s_config USING btree (code);


-- Name: idx_s_dict_data_dict_type; Type: INDEX; Schema: public; Owner: -

CREATE INDEX idx_s_dict_data_dict_type ON public.s_dict_data USING btree (dict_type);


-- Name: idx_s_dict_data_parent_id; Type: INDEX; Schema: public; Owner: -

CREATE INDEX idx_s_dict_data_parent_id ON public.s_dict_data USING btree (parent_id);


-- Name: idx_s_login_log_login_time; Type: INDEX; Schema: public; Owner: -

CREATE INDEX idx_s_login_log_login_time ON public.s_login_log USING btree (login_time);


-- Name: idx_s_login_log_user_name; Type: INDEX; Schema: public; Owner: -

CREATE INDEX idx_s_login_log_user_name ON public.s_login_log USING btree (user_name);


-- Name: idx_s_menu_parent_id; Type: INDEX; Schema: public; Owner: -

CREATE INDEX idx_s_menu_parent_id ON public.s_menu USING btree (parent_id);


-- Name: idx_s_oper_log_oper_time; Type: INDEX; Schema: public; Owner: -

CREATE INDEX idx_s_oper_log_oper_time ON public.s_oper_log USING btree (oper_time);


-- Name: idx_s_org_org_code; Type: INDEX; Schema: public; Owner: -

CREATE UNIQUE INDEX idx_s_org_org_code ON public.s_org USING btree (org_code);


-- Name: idx_s_org_parent_id; Type: INDEX; Schema: public; Owner: -

CREATE INDEX idx_s_org_parent_id ON public.s_org USING btree (parent_id);


-- Name: idx_s_role_role_key; Type: INDEX; Schema: public; Owner: -

CREATE UNIQUE INDEX idx_s_role_role_key ON public.s_role USING btree (role_key);


-- Name: idx_s_storage_env_env_code; Type: INDEX; Schema: public; Owner: -

CREATE UNIQUE INDEX idx_s_storage_env_env_code ON public.s_storage_env USING btree (code);


-- Name: idx_s_user_email; Type: INDEX; Schema: public; Owner: -

CREATE UNIQUE INDEX idx_s_user_email ON public.s_user USING btree (email) WHERE ((email)::text <> ''::text);


-- Name: idx_s_user_open_id; Type: INDEX; Schema: public; Owner: -

CREATE UNIQUE INDEX idx_s_user_open_id ON public.s_user USING btree (open_id) WHERE ((open_id)::text <> ''::text);


-- Name: idx_s_user_org_id; Type: INDEX; Schema: public; Owner: -

CREATE INDEX idx_s_user_org_id ON public.s_user USING btree (org_id);


-- Name: idx_s_user_phonenumber; Type: INDEX; Schema: public; Owner: -

CREATE UNIQUE INDEX idx_s_user_phonenumber ON public.s_user USING btree (phonenumber) WHERE ((phonenumber)::text <> ''::text);


-- Name: idx_s_user_union_id; Type: INDEX; Schema: public; Owner: -

CREATE UNIQUE INDEX idx_s_user_union_id ON public.s_user USING btree (union_id) WHERE ((union_id)::text <> ''::text);


-- Name: idx_s_user_user_name; Type: INDEX; Schema: public; Owner: -

CREATE UNIQUE INDEX idx_s_user_user_name ON public.s_user USING btree (user_name);


-- Name: idx_user_api_permission; Type: INDEX; Schema: public; Owner: -

CREATE UNIQUE INDEX idx_user_api_permission ON public.m_user_api_permission USING btree (user_id, permission_id);


-- Name: idx_user_role; Type: INDEX; Schema: public; Owner: -

CREATE UNIQUE INDEX idx_user_role ON public.m_user_role USING btree (user_id, role_id);

-- +goose Down
DROP TABLE IF EXISTS public.s_user CASCADE;
DROP TABLE IF EXISTS public.s_storage_env CASCADE;
DROP TABLE IF EXISTS public.s_role CASCADE;
DROP TABLE IF EXISTS public.s_org CASCADE;
DROP TABLE IF EXISTS public.s_oper_log CASCADE;
DROP TABLE IF EXISTS public.s_menu CASCADE;
DROP TABLE IF EXISTS public.s_login_log CASCADE;
DROP TABLE IF EXISTS public.s_dict_data CASCADE;
DROP TABLE IF EXISTS public.s_config CASCADE;
DROP TABLE IF EXISTS public.s_auth_client CASCADE;
DROP TABLE IF EXISTS public.s_api_permission CASCADE;
DROP TABLE IF EXISTS public.m_user_role CASCADE;
DROP TABLE IF EXISTS public.m_user_api_permission CASCADE;
DROP TABLE IF EXISTS public.m_role_menu CASCADE;
DROP TABLE IF EXISTS public.m_role_api_permission CASCADE;
DROP TABLE IF EXISTS public.biz_attachment CASCADE;
