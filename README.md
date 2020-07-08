# fserv - file server in go
Simple little file server with file name shortening and jwt authentication tokens.
## Installation
set up postgres db following `shema.sql`  
edit env.sh to change env keys and db urls(or use the docker-compose file for testing locally)  
go build  
## Run with:
`. env.sh && ./fserv`
## Todo:
- [x] File viewing  
- [ ] Albums/grouping of files  
- [ ] User Roles  
- [ ] Admin panel  
