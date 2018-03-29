# README #

This README would normally document whatever steps are necessary to get your application up and running.

### How do I get set up? ###
 - use **go get**
    
   >$ go get github.com/aaa59891/mosi_demo_go
    
 - use **git clone & govendor**
   
   >$ git clone https://github.com/aaa59891/mosi_demo_go.git  
   >$ cd mosi_demo_go  
   >$ govendor sync
   
### Configuration setting: ###

>Using environment variable: **MOSI_GO** to set it
  
>This repo has two envs:  
>1. **MOSI_GO=prod**  
>2. **MOSI_GO=dev**  

>you can set all share config in **/configs/default.toml**,

>then set the different fields in **/configs/prod.toml** OR **/configs/dev.toml** 
  
### Database ###

>This repo uses **MySQL**, version: **at least 5.7**
  
>Can set MySQL configuration in **any env** you want
  
>database setting example:
>>[database]
>>>driveName = "mysql"  
>>>host = "localhost"  
>>>port = 5400  
>>>account = "account"  
>>>password = "password"  
>>>databaseName = "dbName" 

>**There are SQL init scripts in /sqlScripts/initScript, you can create/insert into your MySQL**
>> Init scripts includes accounts: 
>>>admin:admin  
>>>default:default  

>**Or comment routers/routes.go line 19 to disable authorization**

### Session ###
>This repo uses redis as default session 

### How do I run this repo? ###

>1. If you have not installed MySQL or Redis, run the docker scripts below  
>2. Set the config  
>3. Run all the sqlScripts/initScripts
>4. Run the command below depends on the env  

There are two commands in different env:

* Stage env:
```
$ MOSI_GO=dev gin -i go run main.go
```
* Prod  env:
```
$ MOSI_GO=prod go run main.go
```

### Docker Script Example ###
>MySQL
```
$ docker container run --name mysql_mosi -d -e MYSQL_USER=account -e MYSQL_PASSWORD=password -e MYSQL_DATABASE=dbName -p 5400:3306 --restart always -v /your/path/to/store:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=rootPassword mysql
```

>Redis
```
$ docker container run -d --name redis -p 6379:6379 -v /your/path/to/store:/data --restart always redis redis-server --appendonly yes
```

### GO i18n ###
```
$ go get -u github.com/nicksnyder/go-i18n/goi18n
```
Translation data format: json
```
"ID":{
    "other": "VALUE"
}
```
* Add new key
>Add data to **i18n/en-us.all.json** file

* Add new language
>Add new language file(**language.all.json**) to i18n/

>Translate all
```
$ goi18n *all.json
```
>set all the untranslated file data than run
```
$ goi18n *all.json *untranslated.json
```

### Build ###
Please check build.sh is executable, if not run the command:
```
$ chmod +x build.sh
```
Set the GOOS and GOARCH env depend on what os you want to deploy  
GOOS and GOARCH reference:
>https://www.digitalocean.com/community/tutorials/how-to-build-go-executables-for-multiple-platforms-on-ubuntu-16-04

Ubuntu example:
```
$ GOOS=linux GOARCH=amd64 ./build
```

Windows example:
```
$ GOOS=windows GOARCH=386 ./build
```

Then you can find your deployment files in export/*ProjectName*/
Just move this folder to the machine and run the binary file!

