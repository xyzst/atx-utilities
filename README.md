# city of austin, tx utilities

## find-city-council-district

given a single comma separated values (CSV) file with the following data format (no header):

```text
address,city,state,zip_code
```

will then transform each line with the likely council district number, council district website, and confidence score:

```text
address,city,state,zip_code,district,district_url,confidence_score
```

### Use Case

#### Given

- csv w/o header line
- a mix of perfect/complete entries, misspelled addresses, and partial data

```text
200 congress ave, austin, tx, 78701
2713 e 2nd st, austin, tx, 78702
3112 Windsor Rd, austin, tx,
1300 s mopac expy,,,78746
5808 burnet rd, austin, ,78756
1319 Rosewood Ave,,,
3600 presidentiaL rd,,,
8557 Reserch Blv,,,
4001 S Lamr Bld,,,
13429 N US 183,,,
```

#### When

- processing input csv w/ app

#### Then

- will produce `output.csv`
- will reasonably determine the appropriate council district associated with that address
- will add header line

```text
address,city,state,zip_code,district,district_url,confidence_score
200 CONGRESS AVE,AUSTIN,TX,78701,9,http://www.austintexas.gov/department/district-9,98.890
2713 E 2ND ST,AUSTIN,TX,78702,3,http://www.austintexas.gov/department/district-3,99.130
3112 WINDSOR RD,AUSTIN,TX,,10,http://www.austintexas.gov/department/district-10,98.670
1300 S MOPAC EXPY,,,78746,8,http://www.austintexas.gov/department/district-8,99.210
5808 BURNET RD,AUSTIN,,78756,7,http://www.austintexas.gov/department/district-7,98.470
1319 ROSEWOOD AVE,,,,1,http://www.austintexas.gov/department/district-1,88.000
3600 PRESIDENTIAL RD,,,,2,http://www.austintexas.gov/department/district-2,97.560
8557 RESERCH BLV,,,,4,http://www.austintexas.gov/department/district-4,95.480
4001 S LAMR BLD,,,,5,http://www.austintexas.gov/department/district-5,81.180
13429 N US 183,,,,6,http://www.austintexas.gov/department/district-6,93.650
```

### run w/ docker (go not needed!)

```shell
git clone git@github.com:xyzst/atx-utilities.git
cd atx-utilities
cp /path/to/your/your_csv_here.csv .
cp /path/to/your/another_csv_here.csv .

# update environment variable 'TARGET_CSV' to reflect the csv you wish to process.
# In this case, I want to process "your_csv_here.csv" which was copied in the previous commands to current working directory
docker run -e TARGET_CSV='rethink.csv' -v $PWD/output/o.csv:/usr/src/atx-utilities/output.csv --rm -it $(docker build -q .)

# after your TARGET_CSV is processed in the container, the output will be found in
# the output directory
more output/o.csv
```

### run w/o pants (requires go)

```text
$ go run ./... /full/path/to/addresses.csv         
```