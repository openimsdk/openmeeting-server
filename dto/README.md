# Data Transfer Object
## Usage

1. Define the struct of an object using as to transfer data from different layers. like
configuration etc.

> We used to put configuration struct in this directory

2. Define the function using to convert from dto to vo(request).
or define the function using to convert from po(table) to dto and vice versa.

> Like we use this function to convert struct from repository to 
> object which used to act as a passing data transfer to the client.

