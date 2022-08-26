# On Boarding Portal Accounts Service #

Onboarding Portal Accounts service is used to register account for employees and store their information. 
This service uses NATS to communicate with Notification service to send emails

## Requirements ##
1. go (1.13 or above)
2. Postgres Database
3. JWT token
4. NATS Server 2.1 or above

## Documentation ##
1. [Api Documentation on Swagger](openapi.yaml)
2. [Sample Configuration using toml](example.toml)

## Environment Variables => ##
```
//Postgres database
export AUTH_DB_HOST="localhost"
export AUTH_DB_PORT="5432"
export AUTH_DB_USER="Postgres"
export AUTH_DB_PASS="Postgres"
export AUTH_DB_IDEAL_CONNECTIONS = "50"

//admin account credentials configuration
export AUTH_ADMIN_EMAIL="admin@xyz.com"
export AUTH_ADMIN_PASS="admin"

//Service related settings
export ENVIRONMENT="development"
export BUILD_IMAGE="image"
export HOST_ADDRESS="" //this is the url that will be used with token for invite based registration
export AUTH_SERVICE_PORT = ""

//NATS related settings
export NATS_URL="localhost:4222"
export NATS_TOKEN="auth token" // either token or user name/password can be used for authentication
export NATS_USERNAME="user"
export NATS_PASSWORD="password"

//jwt Configuration
export PRIVATE_KEY="same among all the backend services"

//Career Portal
export CAREER_PORTAL_FRONTENDADDRESS = ""
export CAREER_PORTAL_USERNAME = ""
export CAREER_PORTAL_PASSWORD = ""

//Google 
export GOOGLE_CLIENT_ID = ""
export GOOGLE_CLIENT_KEY = ""
export GOOGLE_SCOPES = ""
export GOOGLE_REDIRECT = ""

```

## How to run the app ##


### 1. Database configuration ###

i. Create extension -> 

`create extension dblink;`

ii. Run the Sql query to create view internlist -> 

```
create view internslist as SELECT name,email,level,team,department,account_status,role,reviewerid,slug,scores, total_scores,created_at,profile_img_link,form_filled_at
	FROM   dblink('dbname=? user=? password=?','select m.userid,m.slug,f.scores,f.total_scores,f.created_at from interns_modules_registries as m left join interns_feedback_forms as f on f.module_id = m.id')
	AS     tb2(userid text, slug text, scores numeric,total_scores numeric,form_filled_at timestamptz)
	right join(
		select a2.id,a2.name,a2.email,a2.level,a2.account_status , a2.role,a2.created_at, r2.reviewerid, pcd.team,pcd.department,pi2.profile_img_link from accounts as a2
		left join reviewers as r2 on  a2.id = r2.userid 
		left join personal_infos  as pi2 on a2.id = pi2.userid
		left join pand_c_data as pcd on a2.id = pcd.userid
	) 
	as acc on acc.email = tb2.userid
```

### 2. Configuration using environment variables ###

```
i.    Export above all environment variables
ii.   Build the app or binary -> command -> `$ go install`
iii.  Run the app or binary -> command -> `$GOPATH/bin/accounts --conf=environment`
```

### 3. Configuration using TOML file ###

```
i.    Create a configuration toml file
ii.   Build the app or binary -> command -> `$ go install`
iii.  Run the app or binary -> command -> `$GOPATH/bin/accounts --conf=toml --file=<path of toml file>`
```


### 4. Run the Test Case ###

```
i.    Run the Test case -> go test ./... -p 1 -v -coverprofile=coverage.out
ii.   To see the Test Coverage Output in html page -> go tool cover -html=coverage.out 
```


`Note :- for any help regarding flags, run this command '$GOPATH/bin/accounts --help'`
