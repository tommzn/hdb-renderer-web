CREATE TABLE "datasources" (
    "id"    		VARCHAR2(36) NOT NULL,		
	"type"  		VARCHAR2(32),
	"event_type"  	VARCHAR2(32),
	"name"  		VARCHAR2(64),
	PRIMARY KEY ("id")
);

CREATE TABLE "datasource_config" (
    "datasource_id" VARCHAR2(36) NOT NULL,		
	"key"           VARCHAR2(32),
	"value"         VARCHAR2(256),
	PRIMARY KEY ("datasource_id", "key")
);