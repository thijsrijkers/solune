# Solune - A NoSQL Database

**Solune** is a high-performance, NoSQL database designed with an emphasis on scalability and flexibility. By focusing on efficient internal data storage and retrieval, Solune aims to provide rapid access to data while scaling seamlessly across a wide range of use cases.

## Why Go?

Solune is built with the Go programming language for several key reasons:

- **Concurrency**: Go’s goroutines and channels make it ideal for handling multiple concurrent operations, a critical feature for databases that need to manage large volumes of requests simultaneously.
- **Performance**: Go is known for its fast execution times, making it a great choice for performance-sensitive applications like databases.
- **Simplicity and Readability**: Go is a relatively simple language, making it easy to maintain and extend Solune as the project evolves.
- **Strong Ecosystem**: Go has a rich ecosystem, providing great libraries and tools that help with building robust, high-performance applications.

## In-Memory Data Storage Priority

One of the core design decisions behind Solune is the use of **in-memory data storage**. This choice prioritizes speed and efficiency in data retrieval. Here’s why:

- **Faster Access**: Accessing data in memory is significantly faster than querying disk-based storage, which makes Solune ideal for high-throughput applications where speed is critical.
- **Reduced Latency**: By storing data in memory, Solune reduces the need for costly disk I/O operations, which translates to lower latency and faster response times.
- **Scalability**: In-memory databases can scale easily as the entire database is stored in RAM, and systems with larger amounts of memory can handle growing data volumes without the need for complex disk-based scaling solutions.

While in-memory storage does have certain trade-offs (e.g., limited by system memory), this design choice aligns with Solune focus on speed and flexibility, especially for applications with high performance requirements.

## Getting Started

To run the Solune server locally, follow these steps:

### 1. Clone the repository:

```bash
git clone https://github.com/thijsrijkers/Solune.git
cd Solune
```

### 2. Launch the server:

Inside the `source` folder:

```bash
go run .
```

This will start the Solune server locally, allowing you to interact with the database.

A port will be open through the TCP protocol so you can connect with the database with:

```bash
 telnet localhost 9000  
```

### 3. Command Format
The command follows this format:
```bash
 instruction:=<action>|store=<store_name>|key=<key_of_entry>|data=<data_to_store>
```

Where:
- **`instruction`**: Specifies the action to be performed. The possible actions are:
  - **`get`**: Retrieve the data associated with the given key.
  - **`set`**: Store the provided data under the given key.

- **`store`** (optional): The name of the store or the storage container where the data is to be saved or retrieved from. This is required for both **`get`** and **`set`** actions.

- **`key`** (optional): The unique identifier used to access or save the data within the specified store. If the instruction is **`get`**, the `key` is required to specify which entry to retrieve. If the instruction is **`set`**, the `key` is required to specify the entry under which the data will be stored.

- **`data`** (optional): The data to be stored in the store under the specified key. This is only required for the **`set`** action to define what data you want to save.

##### Example Commands:

1. **Set Data**
   ```bash
   instruction:=set|store=user_data|key=user123|data={"name": "John Doe", "age": 30}
   ```

- This command stores the data `{"name": "John Doe", "age": 30}` in the `user_data` store under the key `user123`.

2. **Get Data**
   ```bash
   instruction:=get|store=user_data|key=user123
   ```
- This command retrieves the data associated with the key user123 from the user_data store.

3. **Get Data Without Key**
   ```bash
   instruction:=get|store=system_config
   ```
- This command retrieves all data from the system_config store without specifying a key. This could be used if the store is designed to return all entries or a default entry.

### 4. Unit Testing
To run the unit tests, navigate to the source folder and run:

```bash
go test ./unit...
```

This will execute the unit tests and display the results in your terminal.
