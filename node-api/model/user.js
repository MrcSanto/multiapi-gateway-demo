/**
 * User entity definition. While Node.js does not require explicit models to
 * interact with a relational database, defining a class helps clarify the
 * structure of the data our application handles. This model mirrors the
 * structure of the users table used in the Go implementation.
 */
class User {
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