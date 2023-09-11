DROP TABLE IF EXISTS termEntry;
DROP TABLE IF EXISTS terms; 
DROP TABLE IF EXISTS docs;

create table docs (
    docName varchar(100) not null primary key
);

CREATE INDEX docIndex 
ON docs USING hash(docName);

create table terms (
    termindex serial,
    containingCount integer,
    termName varchar(15) primary key not null
); 

CREATE INDEX termIndex 
ON terms USING hash(termName);

create table termEntry (
    termName varchar(15) references terms(termName),
    docName varchar(100) references docs(docName),
    tfScore FLOAT
);