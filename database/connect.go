package database

import (
	"co-msa-checker/utils"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

var DbConnection *sql.DB

func Connect() {
	var (
		err          error
		dbUrl        string //database url
		sqlstatement string //SQL statement to exec
	)
	// -------------------------
	// .env loading
	// -------------------------

	// Read dbUrls from env
	dbUrl, err = utils.ReadEnv("DATABASE_URL")
	if err != nil {
		log.Fatalf("FATAL:\terror setting database url: %v", err)
	}
	log.Printf("DbConnection url set")

	// -------------------------
	// DB pool connection
	// -------------------------
	DbConnection, err = sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("FATAL:\tcan't connect to db: %v", err)
	}
	log.Printf("Connection string set")

	// Try ping db to check for availability
	err = DbConnection.Ping()
	if err != nil {
		log.Fatalf("FATAL:\tcan't ping database %v", err)
	}
	log.Printf("DbConnection correctly pinged")

	// -------------------------
	// Init DB if not already done
	// -------------------------

	// users table
	sqlstatement = `create table if not exists users
						(
							id           uuid    default gen_random_uuid() not null
								primary key,
							username     varchar                           not null,
							password     varchar                           not null,
							manager_role boolean default false             not null
						);
						
						comment on table users is 'Users table for auth';
						
						comment on column users.password is 'hashed pwd';
						
						comment on column users.manager_role is 'if user is a manager';

`
	_, err = DbConnection.Exec(sqlstatement)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error creating Document table: %v", err))
	}

	// infolist table
	sqlstatement = `create table if not exists infolist
						(
							id        uuid      default gen_random_uuid() not null
								constraint infolist_pk
									primary key,
							timestamp timestamp default now()              not null,
							note      varchar                              not null,
							status    boolean   default false,
							operator  varchar                              not null,
							priority  varchar                              not null,
							msa_id    uuid                                 not null
								constraint infolist_msa_id_fk
									references msa
									on update cascade on delete cascade
						);
						
						create unique index if not exists infolist_id_uindex
							on infolist (id);

`

	_, err = DbConnection.Exec(sqlstatement)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error creating Document table: %v", err))
	}

	// updates table
	sqlstatement = `create table if not exists updates
						(
							id         uuid      default gen_random_uuid() not null
						constraint updates_pk
						primary key,
						info_id    uuid                                 not null
						constraint updates_infolist_id_fk
						references infolist,
						timestamp  timestamp default now()              not null,
						note       varchar                              not null,
						operator   varchar                              not null,
						deprecated boolean   default false              not null
						);
					
						comment on table updates is 'Contain info updates (linked to infolist)';
					
						create unique index if not exists updates_id_uindex
						on updates (id);
`
	_, err = DbConnection.Exec(sqlstatement)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error creating Document table: %v", err))
	}

	// msa table
	sqlstatement = `create table if not exists msa
						(
							id        uuid default gen_random_uuid() not null
								constraint msa_pk
									primary key,
							radiocode varchar                         not null,
							plate     varchar						not null,
							note      text
						);
						
						create unique index if not exists msa_id_uindex
							on msa (id);
`

	_, err = DbConnection.Exec(sqlstatement)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error creating Document table: %v", err))
	}
}
