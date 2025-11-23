const { Pool } = require('pg');

const pool = new Pool({
  host: process.env.DB_HOST,
  port: process.env.DB_PORT,
  user: process.env.DB_USER,
  password: process.env.DB_PASSWORD,
  database: process.env.DB_NAME,
});


const schema = process.env.DB_SCHEMA_NODE || process.env.DB_SCHEMA;
pool.on('connect', (client) => {
  if (schema) {
    client.query(`SET search_path TO ${schema}`)
      .catch((err) => console.error('Failed to set search_path', err));
  }
});

module.exports = pool;