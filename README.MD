Follows on a Designing Data Intensive Applications.
Might be considered a workbook code.
What is done:
1. Get and set key value pairs into a file. Populate script for synthetic data.
2. Get key is slow for a db file 1.5gb+ with 30 millions of entries
3. Interface is CLI with one off run to write or read from a file.
4. In-memory hash index using byte offset - involved having state, interface now is http server running in memory and client is one-off CLI command. 
5. Server takes a lot of memory 3x size of database file, because of the go map
overhead. Takes 13 seconds to start to run an index scan.
6. Improved and drastically reduced overhead from 4.5 gb to 900mb when stabilized. Custom key value in memory store and key hashing - ints are cheaper then strings.
7. Using stnd http lib server - fast enough for now. request take under 0.1 second. 
8. Implementation of compaction for key deduplication. We keeping only latest record
for the key. 
9. 


