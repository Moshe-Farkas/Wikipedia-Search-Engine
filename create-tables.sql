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

insert into terms(termName)
VALUES
    ('term 1'),
    ('term 2'),
    ('term 3'),
    ('term 4'),
    ('term 5'),
    ('term 6'),
    ('term 7'),
    ('term 8'),
    ('term 9'),
    ('term 10'),
    ('term 11'),
    ('term 12');

insert into docs(docName)
VALUES
    ('doc 1'),
    ('doc 2'),
    ('doc 3'),
    ('doc 4'),
    ('doc 5'),
    ('doc 6'),
    ('doc 7'),
    ('doc 8'),
    ('doc 9'),
    ('doc 10');
          

insert into termEntry(termName, docName, tfScore)
VALUES
    ('term 1', 'doc 7', 0.999),
    ('term 1', 'doc 8', 4.999),
    ('term 1', 'doc 3', 0.459),

    ('term 3', 'doc 3', 3.349),
    ('term 3', 'doc 1', 0.955),
    ('term 3', 'doc 2', 3.999),
    ('term 3', 'doc 4', 0.123),

     ('term 5', 'doc 1', 0.123),
     ('term 5', 'doc 2', 0.123),
     ('term 5', 'doc 3', 0.123),
     ('term 5', 'doc 4', 0.123),
     ('term 5', 'doc 7', 0.123),
     ('term 5', 'doc 8', 0.123),
     ('term 5', 'doc 9', 0.123),
     ('term 5', 'doc 10', 0.123);


          
-- INSERT INTO docs(docName)
-- VALUES
-- ('doc 1'),
-- ('doc 2'),
-- ('doc 3'),
-- ('doc 4'),
-- ('doc 5'),
-- ('doc 6'),
-- ('doc 7'),
-- ('doc 8');


-- create table terms (
--     termIndex int primary key not null,
--     termName varchar(15) not null,
--     containingDoc varchar(50) references docs(docName)
-- );

-- CREATE INDEX termLookupIndex
-- on terms using hash(termName);

-- insert into terms(termIndex, termName, containingDoc)
-- VALUES 
--     (0, 'word1', 'doc 1'),
--     (1, 'word2', 'doc 2'),
--     (2, 'word3', 'doc 3'),
--     (3, 'word4', 'doc 4'),
--     (4, 'word5', 'doc 5'),
--     (5, 'word6', 'doc 6'),
--     (6, 'word7', 'doc 7'),
--     (7, 'word8', 'doc 8');


