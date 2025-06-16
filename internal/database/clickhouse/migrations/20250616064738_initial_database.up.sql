CREATE TABLE project_settings
(
    `project_id` String COMMENT 'Identifier for the project to which these settings apply.',
    `retention_days` UInt32 COMMENT 'Number of days to retain raw event data for this project.',
    `last_updated` DateTime('UTC') DEFAULT now() COMMENT 'Timestamp of the last update to these settings, used by ReplacingMergeTree.'
)
ENGINE = ReplacingMergeTree(last_updated)
PRIMARY KEY (project_id)
ORDER BY (project_id);

CREATE DICTIONARY project_retention_dictionary
(
    `project_id` String,
    `retention_days` UInt32
)
PRIMARY KEY project_id
SOURCE(CLICKHOUSE(TABLE 'project_settings'))
LAYOUT(FLAT())
LIFETIME(MIN 300 MAX 360);

-- Function to get the TTL for events based on project settings
CREATE FUNCTION get_event_ttl
    AS (project_id, default_ttl_days) -> dictGetOrDefault(project_retention_dictionary, 'retention_days', project_id, default_ttl_days);

CREATE TABLE raw_events
(
    -- Core Identifiers
    `project_id` String COMMENT 'Identifier for the project to which this event belongs.',
    `event_id` UUID DEFAULT generateUUIDv4() COMMENT 'Unique identifier for each event.',
    `event_timestamp` DateTime64(3, 'UTC') COMMENT 'High-precision timestamp of when the event occurred on the client-side.',
    `ingestion_timestamp` DateTime64(3, 'UTC') DEFAULT now() COMMENT 'Timestamp of when the event was received and recorded by the system.',
    `event_name` LowCardinality(String) COMMENT 'Name of the event, e.g., "page_view", "click", "form_submit".',
    `source` Enum8('client' = 1, 'server' = 2) COMMENT 'Source of the event, client-side or server-side.',

    -- Visitor & Session Information
    `visitor_fingerprint` String COMMENT 'A unique and persistent identifier for the visitor, generated via fingerprinting.',
    `session_id` String COMMENT 'Identifier for a single user session, helps group events from the same browsing session.',

    -- Page & URL Information
    `url` String COMMENT 'The full URL where the event occurred, including query parameters and fragments.',
    `url_path` String COMMENT 'The path part of the URL (e.g., "/pricing"), useful for aggregation.',
    `url_host` LowCardinality(String) COMMENT 'The hostname from the URL (e.g., "example.com").',
    `url_query` String COMMENT 'The query string part of the URL (e.g., "?id=123").',
    `referrer_url` String COMMENT 'The full URL of the referring page.',
    `referrer_host` LowCardinality(String) COMMENT 'The hostname of the referring URL (e.g., "google.com", "t.co").',

    -- Marketing & Attribution (UTM parameters)
    `utm_source` LowCardinality(Nullable(String)) COMMENT 'UTM source parameter, identifies the advertiser, site, publication, etc.',
    `utm_medium` LowCardinality(Nullable(String)) COMMENT 'UTM medium parameter, identifies the advertising or marketing medium (e.g., "cpc", "email").',
    `utm_campaign` LowCardinality(Nullable(String)) COMMENT 'UTM campaign parameter, identifies a specific promotion or strategic campaign.',
    `utm_term` LowCardinality(Nullable(String)) COMMENT 'UTM term parameter, identifies paid search keywords.',
    `utm_content` LowCardinality(Nullable(String)) COMMENT 'UTM content parameter, used to differentiate ads or links that point to the same URL.',

    -- A/B Testing
    `ab_test_name` LowCardinality(Nullable(String)) COMMENT 'Name of the A/B test the user is part of.',
    `ab_test_variant` LowCardinality(Nullable(String)) COMMENT 'The specific variant of the A/B test shown to the user.',

    -- Geographical & IP-based Location
    `country_code` LowCardinality(String) COMMENT 'ISO 3166-1 alpha-2 country code derived from the visitor''s IP address.',
    `region_name` String COMMENT 'Name of the region or state derived from the visitor''s IP address.',
    `city_name` String COMMENT 'Name of the city derived from the visitor''s IP address.',

    -- Privacy & Security Flags
    `is_vpn` UInt8 COMMENT 'Flag indicating if the visitor is using a VPN (1 for true, 0 for false).',
    `vpn_provider` LowCardinality(Nullable(String)) COMMENT 'Name of the VPN provider if applicable.',
    `is_proxy` UInt8 COMMENT 'Flag indicating if the visitor is using a proxy (1 for true, 0 for false).',
    `proxy_provider` LowCardinality(Nullable(String)) COMMENT 'Name of the proxy provider if applicable.',
    `is_tor_node` UInt8 COMMENT 'Flag indicating if the visitor is using the Tor network (1 for true, 0 for false).',
    `is_bot` UInt8 COMMENT 'Flag indicating if the visitor is a known bot (1 for true, 0 for false).',
    `bot_name` LowCardinality(Nullable(String)) COMMENT 'Name of the bot if applicable (e.g., "Googlebot", "Bingbot").',

    -- Client/Device Information
    `user_agent` String COMMENT 'The full user-agent string of the client''s browser.',
    `browser_name` LowCardinality(String) COMMENT 'Name of the browser (e.g., "Chrome", "Firefox").',
    `browser_version` String COMMENT 'Version of the browser (e.g., "108.0.0").',
    `os_name` LowCardinality(String) COMMENT 'Name of the operating system (e.g., "Windows", "macOS").',
    `os_version` String COMMENT 'Version of the operating system.',
    `device_type` LowCardinality(String) COMMENT 'Type of device used (e.g., "desktop", "mobile", "tablet").',
    `screen_width` Nullable(UInt16) COMMENT 'Width of the device screen in pixels.',
    `screen_height` Nullable(UInt16) COMMENT 'Height of the device screen in pixels.',

    -- Performance Metrics (Core Web Vitals & others)
    `page_load_time_ms` UInt32 COMMENT 'Total time taken for the page to load in milliseconds.',
    `time_on_page_s` UInt16 COMMENT 'Time the user spent on the page in seconds.',
    `first_contentful_paint_ms` UInt32 COMMENT 'First Contentful Paint (FCP) metric in milliseconds, a Core Web Vital.',
    `largest_contentful_paint_ms` UInt32 COMMENT 'Largest Contentful Paint (LCP) metric in milliseconds, a Core Web Vital.',

    -- Custom Data Payload
    `custom_properties` Map(String, String) COMMENT 'A key-value map for sending any custom event-specific data.',

    -- TTL (Time to Live) for raw events
    `retention_days` UInt32 MATERIALIZED get_event_ttl(project_id, 365) COMMENT 'Materialized column determining event TTL based on project settings, defaults to 365 days.'
)
ENGINE = MergeTree()
PARTITION BY toYYYYMM(event_timestamp)
ORDER BY (project_id, url_host, event_name, event_timestamp, visitor_fingerprint)
TTL toDateTime(event_timestamp) + INTERVAL retention_days DAY;

CREATE TABLE visitor_directory
(
    `project_id` String COMMENT 'Identifier for the project to which this visitor belongs.',
    `visitor_fingerprint` String COMMENT 'The unique and persistent identifier for the visitor.',

    `first_seen_timestamp` DateTime('UTC') COMMENT 'Timestamp of when the visitor was first seen.',
    `last_seen_timestamp` DateTime('UTC') COMMENT 'Timestamp of when the visitor was last seen, used by ReplacingMergeTree for updates.',
    `total_sessions` UInt32 COMMENT 'The total number of sessions recorded for this visitor.',
    `total_events` UInt64 COMMENT 'The total number of events recorded for this visitor.',

    `initial_referrer_host` LowCardinality(String) COMMENT 'The hostname of the first-ever referrer for this visitor.',
    `initial_utm_campaign` LowCardinality(String) COMMENT 'The first UTM campaign that brought this visitor to the site.',

    `last_known_country_code` LowCardinality(String) COMMENT 'The most recent country code associated with the visitor.',
    `last_known_device_type` LowCardinality(String) COMMENT 'The most recent device type used by the visitor.',

    `custom_user_properties` Map(String, String) COMMENT 'A key-value map for storing custom properties about the user that persist across sessions.'
)
ENGINE = ReplacingMergeTree(last_seen_timestamp) -- Keeps only the latest version of a visitor row
PRIMARY KEY (project_id, visitor_fingerprint)
ORDER BY (project_id, visitor_fingerprint);
