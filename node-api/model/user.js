/**
 * User entity definition. While Node.js does not require explicit models to
 * interact with a relational database, defining a class helps clarify the
 * structure of the data our application handles. This model mirrors the
 * structure of the users table used in the Go implementation.
 */
class User {
  /**
   * Creates a new User instance.
   *
   * @param {Object} params
   * @param {number} [params.id] - Unique identifier of the user.
   * @param {string} params.name - Full name of the user.
   * @param {string} params.email - Email address (must be unique).
   * @param {string} params.username - Login username (must be unique).
   * @param {string} [params.password] - Hashed password. Optional when sending user data back to clients.
   * @param {Date} [params.created_at] - Timestamp when the user was created.
   * @param {Date} [params.updated_at] - Timestamp of the last user update.
   */
  constructor({ id, name, email, username, password, created_at, updated_at }) {
    this.id = id;
    this.name = name;
    this.email = email;
    this.username = username;
    this.password = password;
    this.created_at = created_at;
    this.updated_at = updated_at;
  }
}

module.exports = User;