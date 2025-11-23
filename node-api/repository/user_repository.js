const pool = require('../db/conn');

/**
 * UserRepository encapsulates all database interactions related to the
 * users table. Each method returns either the retrieved data or metadata
 * about the executed operation. Errors are propagated for callers to handle.
 */
class UserRepository {
  /**
   * Retrieves all users ordered by most recent creation.
   *
   * @returns {Promise<Array>} An array of user records.
   */
  async getUsers() {
    const result = await pool.query('SELECT * FROM users ORDER BY created_at DESC');
    return result.rows;
  }

  /**
   * Retrieves a single user by their unique identifier.
   *
   * @param {number} id - The ID of the user to fetch.
   * @returns {Promise<Object|null>} The user record or null if not found.
   */
  async getUserById(id) {
    const result = await pool.query('SELECT * FROM users WHERE id = $1', [id]);
    return result.rows[0] || null;
  }

  /**
   * Retrieves a user by their email. Returns null if no user exists with
   * that email.
   *
   * @param {string} email - The email to search for.
   * @returns {Promise<Object|null>}
   */
  async getUserByEmail(email) {
    const result = await pool.query('SELECT * FROM users WHERE email = $1', [email]);
    return result.rows[0] || null;
  }

  /**
   * Retrieves a user by their username. Returns null if no user exists.
   *
   * @param {string} username - The username to search for.
   * @returns {Promise<Object|null>}
   */
  async getUserByUserName(username) {
    const result = await pool.query('SELECT * FROM users WHERE username = $1', [username]);
    return result.rows[0] || null;
  }

  /**
   * Creates a new user record in the database.
   *
   * @param {Object} user - The user data to insert.
   * @param {string} user.name - The name of the user.
   * @param {string} user.email - The email of the user.
   * @param {string} user.username - The username of the user.
   * @param {string} user.password - The hashed password of the user.
   * @returns {Promise<number>} The ID of the newly created user.
   */
  async createUser(user) {
    const text =
      'INSERT INTO users(name, email, username, password) VALUES ($1, $2, $3, $4) RETURNING id';
    const values = [user.name, user.email, user.username, user.password];
    const result = await pool.query(text, values);
    return result.rows[0].id;
  }

  /**
   * Updates an existing user record by ID. Only the name, email and username
   * are updated; password changes are not handled here.
   *
   * @param {number} id - The ID of the user to update.
   * @param {Object} user - The new user data.
   * @param {string} user.name - The updated name.
   * @param {string} user.email - The updated email.
   * @param {string} user.username - The updated username.
   * @returns {Promise<number>} Number of rows affected (0 if not found).
   */
  async updateUser(id, user) {
    const text =
      'UPDATE users SET name = $1, email = $2, username = $3, updated_at = NOW() WHERE id = $4';
    const values = [user.name, user.email, user.username, id];
    const result = await pool.query(text, values);
    return result.rowCount;
  }

  /**
   * Deletes a user record from the database by ID.
   *
   * @param {number} id - The ID of the user to delete.
   * @returns {Promise<number>} Number of rows deleted (0 if not found).
   */
  async deleteUser(id) {
    const result = await pool.query('DELETE FROM users WHERE id = $1', [id]);
    return result.rowCount;
  }
}

module.exports = UserRepository;