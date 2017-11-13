create table protobufs(filename string primary key, protobuf blob);
create table labels(mid string primary key, description string);
create table imagelabels(filename string, mid string, score number, constraint imagelabelspk primary key (filename, mid), constraint imagelabelsfk foreign key (mid) references labels(mid));
create table colors(filename string, color string, amount number, constraint colorspk primary key (filename, color));
create table texts(filename string, text string, constraint textspk primary key (filename, text))
