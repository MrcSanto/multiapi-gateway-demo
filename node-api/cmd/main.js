const express = require('express');
const dotenv = require('dotenv');
const userController = require('../controller/user_controller');

// Load environment variables from .env file if present
dotenv.config();

const app = express();

// Middleware to parse JSON bodies
app.use(express.json());

// Root route for a friendly message
app.get('/', (req, res) => {
  res.status(200).send('API em Node funcionando!');
});

// Healthcheck route to identify language/framework
app.get('/healthcheck', (req, res) => {
  res.status(200).json({ status: true, api_lang: 'node', framework: 'express' });
});

// User routes
app.get('/users', userController.getUsers);
app.get('/users/:id', userController.getUserById);
app.post('/users', userController.createUser);
app.put('/users/:id', userController.updateUser);
app.delete('/users/:id', userController.deleteUser);

// Determine the port to listen on
const port = 8000;

// Start the server
app.listen(port, () => {
  console.log(`Server listening on port ${port}`);
});