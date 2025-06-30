write ahead log
就是一个存储前的log系统,用于数据的安全和崩溃一致性

Wiki:
 A write ahead log is an append-only auxiliary disk-resident structure used for crash and transaction recovery. The changes are first recorded in the log, which must be written to [stable storage](https://en.wikipedia.org/wiki/Stable_storage "Stable storage"), before the changes are written to the database.[[2]](https://en.wikipedia.org/wiki/Write-ahead_logging#cite_note-2)
 
Google AI:
Here's a more detailed explanation:

- **Purpose:**
    
    WAL ensures that even if a system crashes, the database can be recovered to a consistent state. This is achieved by writing all changes to a log file first, before they are written to the database itself. 
    
- **How it works:**
    
    When a database operation (like an update or insert) is performed, the corresponding log record is written to the WAL file. Only after the log record is safely written to disk, the actual data changes are applied to the database files. 
    
- **Benefits:**
    
    - **Durability:** Ensures that committed transactions are not lost in the event of a crash. 
    - **Atomicity:** Guarantees that all changes within a transaction are either fully applied or not at all. 
    - **Recovery:** Enables point-in-time recovery and fast crash recovery by replaying the log. 
    - **Efficiency:** Reduces the number of writes to the main data files, potentially improving performance.