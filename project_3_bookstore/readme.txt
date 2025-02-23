Used:
    Mysql
    GORM
    JSON
    GORILLA MUX


Structure:
    PKG:
        config: to connect to database
        controllers: to process data from database and req from users and response
        models: structs to model data
        routes: routing
        utils: marshal unmarshal json


routes:
    Controller functions             route             method 
    GETBOOKS                         /book/            GET
    CREATE BOOK                      /book/            POST
    GET BOOK BY ID                   /book/{bookid}    GET
    UPDATE BOOK                      /book/{bookid}    PUT
    DELETE BOOKS                     /book/{bookid}    DELETE
