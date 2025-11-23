const { Pool } = require('pg');

/**
 * Creates a new PostgreSQL connection pool. Connection parameters are
 * configured via environment variables. The pool is reused across the
 * application to avoid creating a new connection for every query.
 */
const pool = new Pool({
  host: process.env.DB_HOST,
  port: process.env.DB_PORT,
  user: process.env.DB_USER,
  password: process.env.DB_PASSWORD,
  database: process.env.DB_NAME,
});

// Set the search path to the configured schema on every new client connection.
// This ensures queries run against the correct schema (e.g., node_app).
const schema = process.env.DB_SCHEMA_NODE || process.env.DB_SCHEMA;
pool.on('connect', (client) => {
  if (schema) {
    client.query(`SET search_path TO ${schema}`)
      .catch((err) => console.error('Failed to set search_path', err));
  }
});

module.exports = pool;