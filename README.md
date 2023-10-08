# Restaurant Management System

This is a Restaurant Management System built in Go using the Gin framework, designed to help restaurants manage various aspects of their operations efficiently. It provides a set of APIs for managing users, food items, invoices, menus, orders, tables, and more.

## Features

- **User Management:** Allows restaurant staff to sign up and log in.
- **Food Management:** Add, view, and update food items on the menu.
- **Invoice Generation:** Generate and manage invoices for orders.
- **Menu Management:** Create and update menus for different meal options.
- **Order Processing:** Create, view, and update customer orders.
- **Table Management:** Manage restaurant tables, their availability, and reservations.
- **Order Item Tracking:** Keep track of individual items within an order.

## API Endpoints

Here are the main API endpoints available in this project:

- **Users:**
  - GET `/users`: Get a list of users.
  - GET `/users/:user_id`: Get user details by ID.
  - POST `/users/signup`: Sign up as a new user.
  - POST `/users/login`: Log in as an existing user.

- **Foods:**
  - GET `/foods`: Get a list of food items.
  - GET `/foods/:food_id`: Get food item details by ID.
  - POST `/foods`: Create a new food item.
  - PATCH `/foods/:food_id`: Update a food item.

- **Invoices:**
  - GET `/invoices`: Get a list of invoices.
  - GET `/invoices/:invoice_id`: Get invoice details by ID.
  - POST `/invoices`: Create a new invoice.
  - PATCH `/invoices/:invoice_id`: Update an invoice.

- **Menus:**
  - GET `/menus`: Get a list of menus.
  - GET `/menus/:menu_id`: Get menu details by ID.
  - POST `/menus`: Create a new menu.
  - PATCH `/menus/:menu_id`: Update a menu.

- **Order Items:**
  - GET `/orderItems`: Get a list of order items.
  - GET `/orderItems/:orderItem_id`: Get order item details by ID.
  - GET `/orderItems-order/:order_id`: Get order items by order.
  - POST `/orderItems`: Create a new order item.
  - PATCH `/orderItems/:orderItem_id`: Update an order item.

- **Orders:**
  - GET `/orders`: Get a list of orders.
  - GET `/orders/:order_id`: Get order details by ID.
  - POST `/orders`: Create a new order.
  - PATCH `/orders/:order_id`: Update an order.

- **Tables:**
  - GET `/tables`: Get a list of tables.
  - GET `/tables/:table_id`: Get table details by ID.
  - POST `/table`: Create a new table.
  - PATCH `/table/:table_id`: Update a table.

## Usage

To run this project locally, follow these steps:

1. Clone the repository:

2. Install the required dependencies:

3. Set up a MongoDB database and configure the connection in the project.

4. Run the application:

5. The API endpoints can be accessed via `http://localhost:8000`.

## Contributing

Contributions to this project are welcome! Feel free to open issues or submit pull requests to improve or extend the functionality.

## License

This project is licensed under the [MIT License](LICENSE).
