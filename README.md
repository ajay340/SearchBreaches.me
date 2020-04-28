# [SearchBreaches.me](https://searchbreaches.me)

<p align="center">
    <img src="https://img.shields.io/badge/Made%20with-Golang-blue.svg?logo=Go" />
    <img src="https://img.shields.io/badge/Database-Postgres-blue.svg?logo=postgresql" />
    <img src="https://img.shields.io/badge/Using-AWS-orange.svg?logo=amazon" />
    <br>
    Created by Ajay Shah, Harshal Patel, and Eric Gargiulo
    <br>
    <img src="https://github.com/ajay340/SearchBreaches.me/blob/master/media/search.gif" />
</p>

## About
This web application provides a search engine for viewing data breaches, recommending similar data breaches, and downloading a pdf report of the breach.

This project was presented at Petershiem Academic Exposition

## Get Started
Install GoLang

Create and add your postgresql database credentials in `postgres_config.json`, following the format of the `postgres_config_example.json`

cd to src folder
```
go run main.go
```

## Data Sources

#### [Dataset provided by Kaggle](https://www.kaggle.com/alukosayoenoch/cyber-security-breaches-data)

#### [Wikipedia'a List of Data Breaches](https://en.wikipedia.org/wiki/List_of_data_breaches)

#### [HaveIBeenPwned](https://haveibeenpwned.com/)


Group Roles: 
### Ajay Shah:
- Lead developer. Most familiar with GoLang, a strong resource for github and all GoLang coding to other roles. Group Leader and coordinated other jobs. 
- Created system design architecture and modified throughout and maintained integrity of our model through evolving changes in project.
- Emphasis on Test Driven Development, Github Actions, and addressing supplemental issues.
- Database ORM Engineer. SQLite and Postgresql work for database queries and overall assistance with all other roles. 
- Database Migration Engineer - Moving data from sqlite3 to postgresql
- AWS Configuration Engineer - Configured AWS RDS Postgresql, EC2, and Code Deploy
- Recommendation Algorithm designer and implementer 
- CI/CD configuration manager
- Backend and Web Framework Engineer. Utilizing the Echo framework
- Authentication Engineer
### Eric Gargiulo:
- PDF generation engineer
- Testing and Quality Assurance engineer
- Add populator from Wikipedia and HaveIBeenPwned

### Harshal Patel:
- Cleaned and modified original kaggle data.
- Created and implemented all UI/UX design elements: 
    - Functional Navigation with Search, User, Logout. 
    - Frontend AutoComplete Search
    - Webpage Design, all pages
- Created all frontend webpages: Login, Register, Register Success, User, Search Not Found, Error Pages, Index, and About pages.
- Continaul work on representing backend work as features in frontend for users.

