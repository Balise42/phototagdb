create table protobufs(filename string primary key, protobuf blob);
create table labels(mid string primary key, description string);
create table imagelabels(filename string, mid string, score number, constraint imagelabelspk primary key (filename, mid), constraint imagelabelsfk foreign key (mid) references labels(mid));
