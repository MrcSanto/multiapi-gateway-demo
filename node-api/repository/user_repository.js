const pool = require('../db/conn');

/**
 * UserRepository encapsulates all database interactions related to the
 * users table. Each method returns either the retrieved data or metadata
 * about the executed operation. Errors are propagated for callers to handle.
 */
class UserRepository {

  async getUsers() {
    const result = await pool.query('SELECT * FROM users ORDER BY created_at DESC');
    return result.rows;
  }

  async getUserById(id) {
    const result = await pool.query('SELECT * FROM users WHERE id = $1', [id]);
    return result.rows[0] || null;
  }

  async getUserByEmail(email) {
    const result = await pool.query('SELECT * FROM users WHERE email = $1', [email]);
    return result.rows[0] || null;
  }

  async getUserByUserName(username) {
    const result = await pool.query('SELECT * FROM users WHERE username = $1', [username]);
    return result.rows[0] || null;
  }

  async createUser(user) {
    const text =
      'INSERT INTO users(name, email, username, password) VALUES ($1, $2, $3, $4) RETURNING id';
    const values = [user.name, user.email, user.username, user.password];
    const result = await pool.query(text, values);
    return result.rows[0].id;
  }

  async updateUser(id, user) {
    const text =
      'UPDATE users SET name = $1, email = $2, username = $3, updated_at = NOW() WHERE id = $4';
    const values = [user.name, user.email, user.username, id];
    const result = await pool.query(text, values);
    return result.rowCount;
  }

  async deleteUser(id) {
    const result = await pool.query('DELETE FROM users WHERE id = $1', [id]);
    return result.rowCount;
  }
}

module.exports = UserRepository;