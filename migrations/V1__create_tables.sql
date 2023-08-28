-- Create new positions table
create table positions (
	positionid uuid,
	profileid uuid,
	shortOrLong VARCHAR,
	shareName VARCHAR,
	shareAmount double precision,
	shareStartPrice double precision,
	shareEndPrice double precision,
	stopLoss double precision,
	takeProfit double precision,
	openedTime timestamp,
	closedTime timestamp,
	primary key (positionid)
);