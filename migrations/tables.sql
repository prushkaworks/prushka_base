CREATE TABLE IF NOT EXISTS "prushka_beta.User" (
	"id" serial NOT NULL,
	"name" varchar(255) NOT NULL,
	"email" varchar(255) NOT NULL,
	"password" varchar(255) NOT NULL,
	CONSTRAINT "User_pk" PRIMARY KEY ("id")
) WITH (
  OIDS=FALSE
);



CREATE TABLE IF NOT EXISTS "prushka_beta.Privilege" (
	"id" serial NOT NULL,
	"name" varchar(255) NOT NULL,
	CONSTRAINT "Privilege_pk" PRIMARY KEY ("id")
) WITH (
  OIDS=FALSE
);



CREATE TABLE IF NOT EXISTS "prushka_beta.UserPrivilege" (
	"privilege_id" integer NOT NULL,
	"user_id" integer NOT NULL,
	"workspace_id" integer NOT NULL
) WITH (
  OIDS=FALSE
);



CREATE TABLE IF NOT EXISTS "prushka_beta.Workspace" (
	"id" serial NOT NULL,
	"name" varchar(255) NOT NULL,
	"date_created" timestamp NOT NULL,
	CONSTRAINT "Workspace_pk" PRIMARY KEY ("id")
) WITH (
  OIDS=FALSE
);



CREATE TABLE IF NOT EXISTS "prushka_beta.Desk" (
	"id" serial NOT NULL,
	"name" varchar(255) NOT NULL,
	"date_created" timestamp NOT NULL,
	"workspace_id" integer NOT NULL,
	CONSTRAINT "Desk_pk" PRIMARY KEY ("id")
) WITH (
  OIDS=FALSE
);



CREATE TABLE IF NOT EXISTS "prushka_beta.Column" (
	"id" serial NOT NULL,
	"name" varchar(255) NOT NULL,
	"desk_id" integer NOT NULL,
	CONSTRAINT "Column_pk" PRIMARY KEY ("id")
) WITH (
  OIDS=FALSE
);



CREATE TABLE IF NOT EXISTS "prushka_beta.Card" (
	"id" serial NOT NULL,
	"name" varchar(255) NOT NULL,
	"description" TEXT NOT NULL,
	"date_created" timestamp NOT NULL,
	"date_expiration" timestamp NOT NULL,
	"is_done" BOOLEAN NOT NULL,
	"column_id" integer NOT NULL,
	"assigned" integer NOT NULL,
	"creator" integer NOT NULL,
	CONSTRAINT "Card_pk" PRIMARY KEY ("id")
) WITH (
  OIDS=FALSE
);



CREATE TABLE IF NOT EXISTS "prushka_beta.Label" (
	"id" serial NOT NULL,
	"name" varchar(255) NOT NULL,
	"color" integer NOT NULL,
	CONSTRAINT "Label_pk" PRIMARY KEY ("id")
) WITH (
  OIDS=FALSE
);



CREATE TABLE IF NOT EXISTS "prushka_beta.Attachment" (
	"id" serial NOT NULL,
	"path" varchar(255) NOT NULL,
	"card_id" integer NOT NULL,
	CONSTRAINT "Attachment_pk" PRIMARY KEY ("id")
) WITH (
  OIDS=FALSE
);



CREATE TABLE IF NOT EXISTS "prushka_beta.CardsLabel" (
	"card_id" integer NOT NULL,
	"label_id" integer NOT NULL
) WITH (
  OIDS=FALSE
);
