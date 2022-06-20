SET CLIENT_ENCODING TO 'UTF8';

CREATE TABLE public.logs 
(
    id bigserial NOT NULL,
    record_time timestamp without time zone NOT NULL DEFAULT Now(),
    log_time timestamp without time zone,
    service text,
    source text,
    category text,
    level text,
    session text,
    info text,
    url text,
    http_type text,
    json_body jsonb,
    PRIMARY KEY (id)
);

COMMENT ON COLUMN public.logs.log_time IS 'время отправки записи';
COMMENT ON COLUMN public.logs.service IS 'какой сервис записал информацию: имя сервиса';
COMMENT ON COLUMN public.logs.source IS 'откуда пришли данные в рамках сервиса, например: DEMO, STAGE, PROD';
COMMENT ON COLUMN public.logs.category IS 'произвольная информация для возможности группировки и фильтрации';
COMMENT ON COLUMN public.logs.level IS 'информация, ошибка, предупреждение и т.п., например: INFO, ERROR, WARNING';
COMMENT ON COLUMN public.logs.session IS 'сквозной ID для возможности отслеживания записей в рамках одного запроса';
COMMENT ON COLUMN public.logs.url IS 'http url';
COMMENT ON COLUMN public.logs.http_type IS 'тип HTTP запроса POST, GET, PUT, PATCH, DELETE';
COMMENT ON COLUMN public.logs.json_body IS 'тело HTTP запроса и т.п.';

CREATE INDEX idx_logs_date ON public.logs (record_time);
CREATE INDEX idx_logs_log_time ON public.logs (log_time);
CREATE INDEX idx_logs_service ON public.logs (service);
CREATE INDEX idx_logs_source ON public.logs (source);
CREATE INDEX idx_logs_category ON public.logs (category);
CREATE INDEX idx_logs_level ON public.logs (level);
CREATE INDEX idx_logs_session ON public.logs (session);
CREATE INDEX idx_logs_info ON public.logs (info);
CREATE INDEX idx_logs_url ON public.logs (url);
CREATE INDEX idx_logs_http_type ON public.logs (http_type);
CREATE INDEX idx_logs_json_body ON public.logs (json_body);

CREATE TABLE public.http_headers 
(
    record_id bigint NOT NULL,
    header_name text NOT NULL,
    header_value text,
    PRIMARY KEY (record_id, header_name),
    CONSTRAINT fk_http_headers_log FOREIGN KEY (record_id) REFERENCES public.logs (id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE CASCADE
);

COMMENT ON TABLE public.http_headers IS 'http заголовки';

CREATE TABLE public.properties
(
    record_id bigint NOT NULL,
    p_name text NOT NULL,
    p_value text,
    CONSTRAINT properties_pkey PRIMARY KEY (record_id, p_name),
    CONSTRAINT fk_properties_log FOREIGN KEY (record_id) REFERENCES public.logs (id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE CASCADE
);

COMMENT ON TABLE public.properties IS 'произвольные параметры';