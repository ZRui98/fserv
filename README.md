# fserv - file server in go
Simple little file server with file name shortening and jwt authentication tokens.
## Installation
set up postgres db following `schema.sql`  
add all update sql scripts  
edit env.sh to change env keys and db urls(or use the docker-compose file for testing locally)  
go build  
## Run with:
`. env.sh && ./fserv`
## Todo:
- [x] File viewing  
- [x] Private files
- [x] Edit file properties after upload
- [x] View txt files without download
- [ ] Albums/grouping of files  
- [ ] User Roles  
- [ ] Admin panel  
