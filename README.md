# paper
A serverless NoSQL database focussed on scalability & flexibility

<h3> Goal of this project</h3>
<p>The goal of Paper is to provide a way to optimize usage of data inside a database protocol. Through ground design of internal workings of storing data a attempt is made to retrieve data effeciently and quick. </p>

<h4> Functionality goals </h4>
<p>One of the major design changes, which will involve significant logic and functionality, is the implementation of the 'one global source of truth' principle. The concept is straightforward: tables have keys to make entries identifiable, and entries are related to one another to retrieve additional data. With this principle, all primary keys will be unique across all tables. This allows data retrieval to be based on a single source of truth—the unique primary key. However, this presents a challenge when managing relationships between tables, and addressing this will be one of the key challenges of the project.</p>
